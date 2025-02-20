package agent

import (
	"fmt"
	"github.com/sethgrid/pester"
	"log"
	"math/rand"
	"runtime"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func prepareMemStatsForSend(stats *runtime.MemStats) map[string]float64 {
	result := make(map[string]float64)

	result["Alloc"] = float64(stats.Alloc)
	result["BuckHashSys"] = float64(stats.BuckHashSys)
	result["Frees"] = float64(stats.Frees)
	result["GCCPUFraction"] = float64(stats.GCCPUFraction)
	result["GCSys"] = float64(stats.GCSys)
	result["HeapAlloc"] = float64(stats.HeapAlloc)
	result["HeapIdle"] = float64(stats.HeapIdle)
	result["HeapInuse"] = float64(stats.HeapInuse)
	result["HeapObjects"] = float64(stats.HeapObjects)
	result["HeapReleased"] = float64(stats.HeapReleased)
	result["HeapSys"] = float64(stats.HeapSys)
	result["LastGC"] = float64(stats.LastGC)
	result["Lookups"] = float64(stats.Lookups)
	result["MCacheInuse"] = float64(stats.MCacheInuse)
	result["MCacheSys"] = float64(stats.MCacheSys)
	result["MSpanInuse"] = float64(stats.MSpanInuse)
	result["MSpanSys"] = float64(stats.MSpanSys)
	result["Mallocs"] = float64(stats.Mallocs)
	result["NextGC"] = float64(stats.NextGC)
	result["NumForcedGC"] = float64(stats.NumForcedGC)
	result["NumGC"] = float64(stats.NumGC)
	result["OtherSys"] = float64(stats.OtherSys)
	result["PauseTotalNs"] = float64(stats.PauseTotalNs)
	result["StackInuse"] = float64(stats.StackInuse)
	result["StackSys"] = float64(stats.StackSys)
	result["Sys"] = float64(stats.Sys)
	result["TotalAlloc"] = float64(stats.TotalAlloc)
	result["RandomValue"] = rand.Float64()

	return result
}

func getGoPsutilStats() (map[string]float64, error) {
	result := make(map[string]float64)

	memData, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("error in collect gopsutil metrics: %w", err)
	}
	result["TotalMemory"] = float64(memData.Total)
	result["FreeMemory"] = float64(memData.Available)

	infoStats, err := cpu.Info()
	if err != nil {
		return nil, fmt.Errorf("error in get cpuStats: %w", err)
	}
	for i, cpuStats := range infoStats {
		key := fmt.Sprintf("CPUutilization%d", i+1)
		result[key] = float64(cpuStats.CPU)
	}

	return result, nil
}

type Worker struct {
	client           *pester.Client
	collectJobCh     <-chan struct{}
	doneCh           <-chan struct{}
	reportJobCh      <-chan map[string]float64
	collectOutCh     *myChannel[map[string]float64]
	responseReportCh *myChannel[bool]
	config           *config
	id               int
}

func newWorker(id int, client *pester.Client, collectJobCh <-chan struct{}, reportJobCh <-chan map[string]float64,
	collectOutCh *myChannel[map[string]float64], responseReportCh *myChannel[bool], doneCh <-chan struct{},
	config *config) *Worker {
	return &Worker{
		id:               id,
		client:           client,
		collectJobCh:     collectJobCh,
		reportJobCh:      reportJobCh,
		collectOutCh:     collectOutCh,
		responseReportCh: responseReportCh,
		doneCh:           doneCh,
		config:           config,
	}
}

func (w *Worker) unionStats(stat1 map[string]float64, stat2 map[string]float64) map[string]float64 {
	result := make(map[string]float64)
	for k, v := range stat1 {
		result[k] = v
	}
	for k, v := range stat2 {
		result[k] = v
	}
	return result
}

func (w *Worker) collectSendMetrics() {
	defer w.collectOutCh.close()
	defer w.responseReportCh.close()

	continueWork := true
	for continueWork {
		select {
		case <-w.collectJobCh:
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)
			runtimeMetrics := prepareMemStatsForSend(&memStats)
			gopsutilMetrics, err := getGoPsutilStats()
			if err != nil {
				gopsutilMetrics = map[string]float64{}
			}
			result := w.unionStats(runtimeMetrics, gopsutilMetrics)
			err = w.collectOutCh.send(result)
			if err != nil {
				log.Println("Send metrics in closed channel")
				continueWork = false
			}
		case metrics := <-w.reportJobCh:
			pollCount, ok := metrics["pollCount"]
			if !ok {
				log.Println("Collecting metrics not contain pollCount")
				err := w.responseReportCh.send(false)
				if err != nil {
					log.Println("Send metrics in closed channel")
					continueWork = false
				}
				break
			}
			delete(metrics, "pollCount")

			err := SendGauge(w.client, metrics, w.config)
			if err != nil {
				log.Printf("error in send gauge: %v\n", err)
				err := w.responseReportCh.send(false)
				if err != nil {
					log.Println("Send response report in closed channel")
					continueWork = false
				}
				break
			}

			err = SendCounter(w.client, int(pollCount), w.config)
			if err != nil {
				log.Printf("error in send counter: %v\n", err)
				err := w.responseReportCh.send(false)
				if err != nil {
					log.Println("Send response report in closed channel")
					continueWork = false
				}
				break
			}

			err = BatchSendGauge(w.client, metrics, w.config)
			if err != nil {
				log.Printf("error in batch send gauge: %v\n", err)
				err := w.responseReportCh.send(false)
				if err != nil {
					log.Println("Send response report in closed channel")
					continueWork = false
				}
				break
			}

			err = BatchSendCounter(w.client, int(pollCount), w.config)
			if err != nil {
				log.Printf("error in batch send counter: %v\n", err)
				err := w.responseReportCh.send(false)
				if err != nil {
					log.Println("Send response report in closed channel")
					continueWork = false
				}
				break
			}

			err = w.responseReportCh.send(true)
			if err != nil {
				log.Println("Send response report in closed channel")
				continueWork = false
			}
		case <-w.doneCh:
			continueWork = false
		}
	}
}
