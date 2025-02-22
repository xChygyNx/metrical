package agent

import (
	"log"

	"github.com/sethgrid/pester"

	myconfig "github.com/xChygyNx/metrical/internal/agent/config"
)

const (
	sendCloseChannelMsg = "Send response report in closed channel"
)

type Worker struct {
	client           *pester.Client
	collectJobCh     <-chan struct{}
	doneCh           <-chan struct{}
	reportJobCh      <-chan map[string]float64
	collectOutCh     *myChannel[map[string]float64]
	responseReportCh *myChannel[bool]
	metrics          [2]metrics
	config           *myconfig.AgentConfig
	id               int
}

func newWorker(id int, client *pester.Client, collectJobCh <-chan struct{}, reportJobCh <-chan map[string]float64,
	collectOutCh *myChannel[map[string]float64], responseReportCh *myChannel[bool], doneCh <-chan struct{},
	config *myconfig.AgentConfig) *Worker {
	return &Worker{
		id:               id,
		client:           client,
		collectJobCh:     collectJobCh,
		reportJobCh:      reportJobCh,
		collectOutCh:     collectOutCh,
		responseReportCh: responseReportCh,
		doneCh:           doneCh,
		config:           config,
		metrics:          [2]metrics{NewMemStatsMetrics(), NewGoPsUtilMetrics()},
	}
}

func (w *Worker) collectMetrics() map[string]float64 {
	result := make(map[string]float64)
	for i := range w.metrics {
		metric := w.metrics[i]
		collectingMetrics, err := metric.collectMetrics()
		if err != nil {
			continue
		}
		for k, v := range collectingMetrics {
			result[k] = v
		}
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
			result := w.collectMetrics()
			err := w.collectOutCh.send(result)
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
					log.Println(sendCloseChannelMsg)
					continueWork = false
				}
				break
			}

			err = SendCounter(w.client, int(pollCount), w.config)
			if err != nil {
				log.Printf("error in send counter: %v\n", err)
				err := w.responseReportCh.send(false)
				if err != nil {
					log.Println(sendCloseChannelMsg)
					continueWork = false
				}
				break
			}

			err = BatchSendGauge(w.client, metrics, w.config)
			if err != nil {
				log.Printf("error in batch send gauge: %v\n", err)
				err := w.responseReportCh.send(false)
				if err != nil {
					log.Println(sendCloseChannelMsg)
					continueWork = false
				}
				break
			}

			err = BatchSendCounter(w.client, int(pollCount), w.config)
			if err != nil {
				log.Printf("error in batch send counter: %v\n", err)
				err := w.responseReportCh.send(false)
				if err != nil {
					log.Println(sendCloseChannelMsg)
					continueWork = false
				}
				break
			}

			err = w.responseReportCh.send(true)
			if err != nil {
				log.Println(sendCloseChannelMsg)
				continueWork = false
			}
		case <-w.doneCh:
			continueWork = false
		}
	}
}
