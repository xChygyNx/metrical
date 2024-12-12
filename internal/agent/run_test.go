package agent

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
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

	memStats := prepareStatsForSend(&memStat)
	for _, stat := range stats {
		msg = "MemStat not contain stat " + stat
		_, ok := memStats[stat]
		assert.True(t, ok, msg)
	}
}
