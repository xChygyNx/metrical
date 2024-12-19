package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetGauge(t *testing.T) {
	tests := []struct {
		name   string
		metric string
		set    float64
		want   gauge
	}{
		{
			name:   "Set Gauge metric",
			metric: "some_metric",
			set:    10.789,
			want:   10.789,
		},
		{
			name:   "Set Gauge metric again",
			metric: "some_metric",
			set:    32.6017,
			want:   32.6017,
		},
		{
			name:   "Set other Gauge metric",
			metric: "other_metric",
			set:    73.08,
			want:   73.08,
		},
		{
			name:   "Set other Gauge metric again",
			metric: "other_metric",
			set:    23,
			want:   23,
		},
	}
	storage := GetMemStorage()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage.SetGauge(test.metric, test.set)
			assert.Equal(t, test.want, storage.Gauges[test.metric])
		})
	}
}

func TestGetGauge(t *testing.T) {
	storage := GetMemStorage()
	existMetric := "exist_metric"
	notExistMetric := "not_exist_metric"
	val := 10.789
	storage.Gauges[existMetric] = gauge(val)
	tests := []struct {
		name   string
		metcic string
		want   float64
		ok     bool
	}{
		{
			name:   "Get exists Gauge metric",
			metcic: existMetric,
			want:   val,
			ok:     true,
		},
		{
			name:   "Get not exists Gauge metric",
			metcic: notExistMetric,
			want:   0,
			ok:     false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gaugeVal, ok := storage.GetGauge(test.metcic)
			assert.Equal(t, test.want, gaugeVal)
			assert.Equal(t, test.ok, ok)
		})
	}
}

func TestSetCounter(t *testing.T) {
	tests := []struct {
		name   string
		metric string
		set    int64
		want   counter
	}{
		{
			name:   "Set Counter metric",
			metric: "some_metric",
			set:    15,
			want:   15,
		},
		{
			name:   "Set Counter metric again",
			metric: "some_metric",
			set:    17,
			want:   32,
		},
		{
			name:   "Set other Counter metric",
			metric: "other_metric",
			set:    35,
			want:   35,
		},
		{
			name:   "Set other Counter metric again",
			metric: "other_metric",
			set:    23,
			want:   58,
		},
	}
	storage := GetMemStorage()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage.SetCounter(test.metric, test.set)
			assert.Equal(t, test.want, storage.Counters[test.metric])
		})
	}
}

func TestGetCounter(t *testing.T) {
	storage := GetMemStorage()
	existMetric := "exist_metric"
	notExistMetric := "not_exist_metric"
	var val int64 = 10
	storage.Counters[existMetric] = counter(val)
	tests := []struct {
		name   string
		metcic string
		want   int64
		ok     bool
	}{
		{
			name:   "Get exists Gauge metric",
			metcic: existMetric,
			want:   val,
			ok:     true,
		},
		{
			name:   "Get not exists Gauge metric",
			metcic: notExistMetric,
			want:   0,
			ok:     false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			counterVal, ok := storage.GetCounter(test.metcic)
			assert.Equal(t, test.want, counterVal)
			assert.Equal(t, test.ok, ok)
		})
	}
}
