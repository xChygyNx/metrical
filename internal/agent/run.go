package agent

import (
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/sethgrid/pester"
)

const (
	countRetries = 3
)

func getRetryClient() *pester.Client {
	client := pester.New()
	client.MaxRetries = countRetries
	client.Backoff = pester.ExponentialBackoff
	return client
}

func prepareStatsForSend(stats *runtime.MemStats) map[string]float64 {
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

func Run() error {
	var pollCount int
	var memStats runtime.MemStats

	config, err := GetConfig()
	if err != nil {
		return err
	}
	pollTicker := time.NewTicker(time.Duration(config.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(config.ReportInterval) * time.Second)
	for {
		select {
		case <-pollTicker.C:
			runtime.ReadMemStats(&memStats)
			pollCount++
		case <-reportTicker.C:
			sendInfo := prepareStatsForSend(&memStats)
			client := getRetryClient()

			err = SendGauge(client, sendInfo, config.HostAddr)
			if err != nil {
				log.Printf("error in send gauge: %v\n", err)
				continue
			}

			err = SendCounter(client, pollCount, config.HostAddr)
			if err != nil {
				log.Printf("error in send counter: %v\n", err)
				continue
			}

			//err = BatchSendGauge(client, sendInfo, config)
			//if err != nil {
			//	log.Printf("error in batch send gauge: %v\n", err)
			//	continue
			//}
			//
			//err = BatchSendCounter(client, pollCount, config)
			//if err != nil {
			//	log.Printf("error in batch send counter: %v\n", err)
			//	continue
			//}

			pollCount = 0
		}
	}
}
