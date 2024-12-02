package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func TestPrepareStatsForSend(t *testing.T) {
	var memStat runtime.MemStats
	var msg string

	runtime.ReadMemStats(&memStat)
	stats := []string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys",
		"HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased",
		"HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys",
		"MSpanInuse", "Mallocs", "NextGC", "NumForcedGC", "NumGC",
		"OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys",
		"TotalAlloc", "RandomValue"}

	memStats := prepareStatsForSend(memStat)
	for _, stat := range stats {
		msg = fmt.Sprintf("MemStat not contain stat %s", stat)
		_, ok := memStats[stat]
		assert.True(t, ok, msg)

	}
}
