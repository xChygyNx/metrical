package agent

import (
	"fmt"
	"math/rand"
	"runtime"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type metrics interface {
	collectMetrics() (map[string]float64, error)
}

type MemStatsMetrics struct {
	store *runtime.MemStats
}

type GoPsUtilMetrics struct {
	store map[string]float64
}

func NewMemStatsMetrics() *MemStatsMetrics {
	return &MemStatsMetrics{}
}

func NewGoPsUtilMetrics() *GoPsUtilMetrics {
	return &GoPsUtilMetrics{
		store: make(map[string]float64),
	}
}

func (msm *MemStatsMetrics) prepareMemStatsForSend() map[string]float64 {
	result := make(map[string]float64)

	result["Alloc"] = float64(msm.store.Alloc)
	result["BuckHashSys"] = float64(msm.store.BuckHashSys)
	result["Frees"] = float64(msm.store.Frees)
	result["GCCPUFraction"] = float64(msm.store.GCCPUFraction)
	result["GCSys"] = float64(msm.store.GCSys)
	result["HeapAlloc"] = float64(msm.store.HeapAlloc)
	result["HeapIdle"] = float64(msm.store.HeapIdle)
	result["HeapInuse"] = float64(msm.store.HeapInuse)
	result["HeapObjects"] = float64(msm.store.HeapObjects)
	result["HeapReleased"] = float64(msm.store.HeapReleased)
	result["HeapSys"] = float64(msm.store.HeapSys)
	result["LastGC"] = float64(msm.store.LastGC)
	result["Lookups"] = float64(msm.store.Lookups)
	result["MCacheInuse"] = float64(msm.store.MCacheInuse)
	result["MCacheSys"] = float64(msm.store.MCacheSys)
	result["MSpanInuse"] = float64(msm.store.MSpanInuse)
	result["MSpanSys"] = float64(msm.store.MSpanSys)
	result["Mallocs"] = float64(msm.store.Mallocs)
	result["NextGC"] = float64(msm.store.NextGC)
	result["NumForcedGC"] = float64(msm.store.NumForcedGC)
	result["NumGC"] = float64(msm.store.NumGC)
	result["OtherSys"] = float64(msm.store.OtherSys)
	result["PauseTotalNs"] = float64(msm.store.PauseTotalNs)
	result["StackInuse"] = float64(msm.store.StackInuse)
	result["StackSys"] = float64(msm.store.StackSys)
	result["Sys"] = float64(msm.store.Sys)
	result["TotalAlloc"] = float64(msm.store.TotalAlloc)
	result["RandomValue"] = rand.Float64()

	return result
}

func (msm *MemStatsMetrics) collectMetrics() (map[string]float64, error) {
	runtime.ReadMemStats(msm.store)
	result := msm.prepareMemStatsForSend()

	return result, nil
}

func (gpum *GoPsUtilMetrics) collectMetrics() (map[string]float64, error) {
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
	for i := range infoStats {
		key := fmt.Sprintf("CPUutilization%d", i+1)
		result[key] = float64(infoStats[i].CPU)
	}

	return result, nil
}
