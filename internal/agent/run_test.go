package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepareStatsForSend(t *testing.T) {
	var msg string

	metricCollector := NewMemStatsMetrics()

	result, err := metricCollector.collectMetrics()
	assert.Nil(t, err)
	stats := []string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys",
		"HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased",
		"HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys",
		"MSpanInuse", "Mallocs", "NextGC", "NumForcedGC", "NumGC",
		"OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys",
		"TotalAlloc", "RandomValue"}

	for _, stat := range stats {
		msg = "MemStat not contain stat " + stat
		_, ok := result[stat]
		assert.True(t, ok, msg)
	}
}
