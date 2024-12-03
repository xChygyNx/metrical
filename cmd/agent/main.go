package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

func prepareStatsForSend(stats runtime.MemStats) map[string]float64 {
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

func main() {
	var pollCount int
	var memStats runtime.MemStats

	agentConfig := parseFlag()
	pollInterval := agentConfig.PollInterval
	reportInterval := agentConfig.ReportInterval

	timeReport := time.Now()
	for {
		time.Sleep(time.Second * time.Duration(pollInterval))
		runtime.ReadMemStats(&memStats)
		pollCount += 1
		if time.Now().After(timeReport.Add(time.Duration(reportInterval) * time.Second)) {

			sendInfo, err := json.Marshal(prepareStatsForSend(memStats))
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Printf("Type sendInfo: %T\n", sendInfo)
			client := &http.Client{}
			err = SendGauge(client, sendInfo)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				continue
			}

			err = SendCounter(client, pollCount)
			if err != nil {
				fmt.Println(err)
				continue
			}
			pollCount = 0
			timeReport = time.Now()
		}
	}
}
