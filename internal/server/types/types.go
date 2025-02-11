package types

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
)

type gauge float64

type counter int64

type MemStorage struct {
	Gauges   map[string]gauge   `json:"gauges"`
	Counters map[string]counter `json:"counters"`
}

type HandlerConfig struct {
	DB                *sql.DB
	FileMetricStorage string
	Sha256Key         string
	SyncFileRecord    bool
}

func (s *HandlerConfig) CheckBDConnection() error {
	err := s.DB.PingContext(context.Background())
	return fmt.Errorf("DB is unreachable: %w", err)
}

func GetMemStorage() *MemStorage {
	instance := new(MemStorage)
	instance.Gauges = map[string]gauge{}
	instance.Counters = map[string]counter{}
	return instance
}

func (ms *MemStorage) SetGauge(mName string, mValue float64) {
	ms.Gauges[mName] = gauge(mValue)
}

func (ms *MemStorage) SetCounter(mName string, mValue int64) {
	ms.Counters[mName] += counter(mValue)
}

func (ms *MemStorage) GetGauge(mName string) (float64, bool) {
	metric, ok := ms.Gauges[mName]
	return float64(metric), ok
}

func (ms *MemStorage) GetCounter(mName string) (int64, bool) {
	metric, ok := ms.Counters[mName]
	return int64(metric), ok
}

func (ms *MemStorage) GetGauges() map[string]string {
	gauges := make(map[string]string)
	for k, v := range ms.Gauges {
		gauges[k] = strconv.FormatFloat(float64(v), 'f', -1, 64)
	}
	return gauges
}

func (ms *MemStorage) GetCounters() map[string]string {
	counters := make(map[string]string)
	for k, v := range ms.Counters {
		counters[k] = strconv.FormatInt(int64(v), 10)
	}
	return counters
}

func (ms *MemStorage) SetGauges(data map[string]float64) {
	for k, v := range data {
		ms.Gauges[k] = gauge(v)
	}
}

func (ms *MemStorage) SetCounters(data map[string]float64) {
	for k, v := range data {
		ms.Counters[k] += counter(v)
	}
}

type (
	ResponseData struct {
		Status int
		Size   int
	}

	LoggingResponseWriter struct {
		http.ResponseWriter
		ResponseData *ResponseData
	}
)

func (lrw *LoggingResponseWriter) Write(b []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(b)
	if err != nil {
		err = fmt.Errorf("error in method Write of loggingResponseWrirer: %w", err)
	}
	lrw.ResponseData.Size += size
	return size, err
}

func (lrw *LoggingResponseWriter) WriteHeader(statusCode int) {
	lrw.ResponseWriter.WriteHeader(statusCode)
	lrw.ResponseData.Status = statusCode
}
