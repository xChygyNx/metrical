package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetGauge(t *testing.T) {
	tests := []struct {
		name 	string
		metcic  string
		set		float64
		want 	gauge
	}{
		{
			name: "Set Gauge metric",
			metcic:  "some_metric",
			set: 10.789,
			want: 10.789,
		},
		{
			name: "Set Gauge metric again",
			metcic:  "some_metric",
			set: 32.6017,
			want: 32.6017,
		},
		{
			name: "Set other Gauge metric",
			metcic:  "other_metric",
			set: 73.08,
			want: 73.08,
		},
		{
			name: "Set other Gauge metric again",
			metcic:  "other_metric",
			set: 23,
			want: 23,
		},
	}
	storage := GetMemStorage()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage.SetGauge(test.metcic, test.set)
			assert.Equal(t, test.want, storage.Gauges[test.metcic])
		})
	}
}

func TestGetGauge(t *testing.T) {
	storage := GetMemStorage()
	exist_metric := "exist_metric"
	not_exist_metric := "not_exist_metric"
	val := 10.789
	storage.Gauges[exist_metric] = gauge(val)
	tests := []struct {
		name 	string
		metcic  string
		want 	float64
		ok		bool
	}{
		{
			name: "Get exists Gauge metric",
			metcic:  exist_metric,
			want: val,
			ok: true,
		},
		{
			name: "Get not exists Gauge metric",
			metcic:  not_exist_metric,
			want: 0,
			ok: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gauge_val, ok := storage.GetGauge(test.metcic)
			assert.Equal(t, test.want, gauge_val)
			assert.Equal(t, test.ok, ok)
		})
	}
}

func TestSetCounter(t *testing.T) {
	tests := []struct {
		name 	string
		metcic  string
		set		int64
		want 	counter
	}{
		{
			name: "Set Counter metric",
			metcic:  "some_metric",
			set: 15,
			want: 15,
		},
		{
			name: "Set Counter metric again",
			metcic:  "some_metric",
			set: 17,
			want: 32,
		},
		{
			name: "Set other Counter metric",
			metcic:  "other_metric",
			set: 35,
			want: 35,
		},
		{
			name: "Set other Counter metric again",
			metcic:  "other_metric",
			set: 23,
			want: 58,
		},
	}
	storage := GetMemStorage()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage.SetConunter(test.metcic, test.set)
			assert.Equal(t, test.want, storage.Counters[test.metcic])
		})
	}
}

func TestGetCounter(t *testing.T) {
	storage := GetMemStorage()
	exist_metric := "exist_metric"
	not_exist_metric := "not_exist_metric"
	var val int64 = 10
	storage.Counters[exist_metric] = counter(val)
	tests := []struct {
		name 	string
		metcic  string
		want 	int64
		ok		bool
	}{
		{
			name: "Get exists Gauge metric",
			metcic:  exist_metric,
			want: val,
			ok: true,
		},
		{
			name: "Get not exists Gauge metric",
			metcic:  not_exist_metric,
			want: 0,
			ok: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gauge_val, ok := storage.GetCounter(test.metcic)
			assert.Equal(t, test.want, gauge_val)
			assert.Equal(t, test.ok, ok)
		})
	}
}

