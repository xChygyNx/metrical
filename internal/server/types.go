package server

import (
	"fmt"
	"net/http"
	"strconv"
)

type gauge float64

type counter int64

type memStorage struct {
	Gauges   map[string]gauge   `json:"gauges"`
	Counters map[string]counter `json:"counters"`
}

func GetMemStorage() *memStorage {
	instance := new(memStorage)
	instance.Gauges = map[string]gauge{}
	instance.Counters = map[string]counter{}
	return instance
}

func (ms *memStorage) SetGauge(mName string, mValue float64) {
	ms.Gauges[mName] = gauge(mValue)
}

func (ms *memStorage) SetCounter(mName string, mValue int64) {
	ms.Counters[mName] += counter(mValue)
}

func (ms *memStorage) GetGauge(mName string) (float64, bool) {
	metric, ok := ms.Gauges[mName]
	return float64(metric), ok
}

func (ms *memStorage) GetCounter(mName string) (int64, bool) {
	metric, ok := ms.Counters[mName]
	return int64(metric), ok
}

func (ms *memStorage) GetGauges() map[string]string {
	gauges := make(map[string]string)
	for k, v := range ms.Gauges {
		gauges[k] = strconv.FormatFloat(float64(v), 'f', -1, 64)
	}
	return gauges
}

func (ms *memStorage) GetCounters() map[string]string {
	counters := make(map[string]string)
	for k, v := range ms.Counters {
		counters[k] = strconv.FormatInt(int64(v), 10)
	}
	return counters
}

func (ms *memStorage) SetGauges(data map[string]float64) {
	for k, v := range data {
		ms.Gauges[k] = gauge(v)
	}
}

func (ms *memStorage) SetCounters(data map[string]float64) {
	for k, v := range data {
		ms.Counters[k] += counter(v)
	}
}

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(b)
	lrw.responseData.size += size
	return size, fmt.Errorf("error in method Write of loggingResponseWrirer: %w\n", err)
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.ResponseWriter.WriteHeader(statusCode)
	lrw.responseData.status = statusCode
}
