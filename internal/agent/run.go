package agent

import (
	"errors"
	config2 "github.com/xChygyNx/metrical/internal/agent/config"
	"log"
	"sync"
	"time"

	"github.com/sethgrid/pester"
)

const (
	countRetries = 3
)

type myChannel[T map[string]float64 | bool] struct {
	C      chan T
	closed bool
	once   sync.Once
	mutex  sync.Mutex
}

func newMyChannel[T map[string]float64 | bool](capacity int) *myChannel[T] {
	return &myChannel[T]{
		C: make(chan T, capacity),
	}
}

func (mc *myChannel[T]) close() {
	mc.once.Do(func() {
		mc.closed = true
		close(mc.C)
	})
}

func (mc *myChannel[T]) send(data T) error {
	if !mc.closed {
		mc.mutex.Lock()
		defer mc.mutex.Unlock()
		mc.C <- data
		return nil
	}
	return errors.New("channel is closed")
}

func (mc *myChannel[T]) get() (T, bool) {
	data, ok := <-mc.C
	return data, ok
}

func getRetryClient() *pester.Client {
	client := pester.New()
	client.MaxRetries = countRetries
	client.Backoff = pester.ExponentialBackoff
	return client
}

func Run() error {
	var pollCount int
	var sendInfo map[string]float64

	config, err := config2.GetConfig()
	if err != nil {
		return err
	}
	pollTicker := time.NewTicker(time.Duration(config.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(config.ReportInterval) * time.Second)
	client := getRetryClient()
	collectJobCh := make(chan struct{})
	doneCh := make(chan struct{})
	reportJobCh := make(chan map[string]float64)
	collectOutCh := newMyChannel[map[string]float64](1)
	responseReportCh := newMyChannel[bool](1)

	defer close(doneCh)

	for i := 0; i < config.RateLimit; i++ {
		worker := newWorker(i+1, client, collectJobCh, reportJobCh, collectOutCh, responseReportCh, doneCh, config)
		go worker.collectSendMetrics()
	}
	interrupt := false
	for !interrupt {
		select {
		case <-pollTicker.C:
			collectJobCh <- struct{}{}
			sendInfo = <-reportJobCh
			pollCount++
			sendInfo["PollCounter"] = float64(pollCount)
		case <-reportTicker.C:
			reportJobCh <- sendInfo

			successSend, ok := responseReportCh.get()
			if !ok {
				log.Println("Response report Channel is closed")
				interrupt = true
				break
			}
			if successSend {
				sendInfo = make(map[string]float64)
				pollCount = 0
			}
		}
	}
	return nil
}
