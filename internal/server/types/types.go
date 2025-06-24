// Package types определяет структуры для работы сервера сбора метрик, а также их
// методы для удобной работы с ними.
package types

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
)

type gauge float64

type counter int64

// MemStorage структура для хранения gauges и counters метрик.
type MemStorage struct {
	Gauges   map[string]gauge   `json:"gauges"`   // метрики типа gauge.
	Counters map[string]counter `json:"counters"` // метрики типа counter.
}

// SyncInfo структура хранящая подключение к базе данных и файлу где следует хранить
// собранные метрики, а также атрибут определяющий, следует ли немедленно сохранять
// присланные метрики в файл или БД. Так же хранит доверенную подсеть, из которой
// можно обращаться к серверу
type SyncInfo struct {
	DB                *sql.DB // подключение к PostgreSQL базе данных для хранения собранных метрик.
	FileMetricStorage string  // файл для хранения собранных метрик.
	TrustedSubnet     string  // Доверенная подсеть.
	SyncFileRecord    bool    // атрибут определяющий, следует ли немедленно сохранять присланные метрики
	// в файл или БД.
}

// GetMemStorage возвращает хранилище метрик.
func GetMemStorage() *MemStorage {
	instance := new(MemStorage)
	instance.Gauges = map[string]gauge{}
	instance.Counters = map[string]counter{}
	return instance
}

// SetGauge перезаписывает значение заданной метрики типа gauge в хранилище.
func (ms *MemStorage) SetGauge(mName string, mValue float64) {
	ms.Gauges[mName] = gauge(mValue)
}

// SetCounter увеличивает значение заданной метрики типа counter в хранилище на заданную величину.
func (ms *MemStorage) SetCounter(mName string, mValue int64) {
	ms.Counters[mName] += counter(mValue)
}

// GetGauge возвращает значение заданной метрики типа gauge из хранилища.
func (ms *MemStorage) GetGauge(mName string) (float64, bool) {
	metric, ok := ms.Gauges[mName]
	return float64(metric), ok
}

// GetCounter возвращает значение заданной метрики типа counter из хранилища.
func (ms *MemStorage) GetCounter(mName string) (int64, bool) {
	metric, ok := ms.Counters[mName]
	return int64(metric), ok
}

// GetGauges возвращает значение всех метрик типа gauge из хранилища.
func (ms *MemStorage) GetGauges() map[string]string {
	gauges := make(map[string]string)
	for k, v := range ms.Gauges {
		gauges[k] = strconv.FormatFloat(float64(v), 'f', -1, 64)
	}
	return gauges
}

// GetCounters возвращает значение всех метрик типа counter из хранилища.
func (ms *MemStorage) GetCounters() map[string]string {
	counters := make(map[string]string)
	for k, v := range ms.Counters {
		counters[k] = strconv.FormatInt(int64(v), 10)
	}
	return counters
}

// SetGauges заменяет значение нескольких переданных метрик типа gauge в хранилище.
func (ms *MemStorage) SetGauges(data map[string]float64) {
	for k, v := range data {
		ms.Gauges[k] = gauge(v)
	}
}

// SetCounters заменяет значение нескольких переданных метрик типа counter в хранилище.
func (ms *MemStorage) SetCounters(data map[string]float64) {
	for k, v := range data {
		ms.Counters[k] += counter(v)
	}
}

// Структуры для логирования расширенной информации о http ответе.
type (
	// ResponseData структура с информацией о http ответе.
	ResponseData struct {
		Status int // http статус ответа сервера.
		Size   int // размер ответа сервера.
	}

	// LoggingResponseWriter структура для логирования данных об пришедшем http ответе.
	LoggingResponseWriter struct {
		http.ResponseWriter
		ResponseData *ResponseData // данные из тела ответа.
	}
)

// Write записывает переданные данные в тело ответа и сохраняет саккумуированный размер записанных данных.
func (lrw *LoggingResponseWriter) Write(b []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(b)
	if err != nil {
		err = fmt.Errorf("error in method Write of loggingResponseWrirer: %w", err)
	}
	lrw.ResponseData.Size += size
	return size, err
}

// WriteHeader записывает http код в заголовок ответа и сохраняет его у себя для записи в лог.
func (lrw *LoggingResponseWriter) WriteHeader(statusCode int) {
	lrw.ResponseWriter.WriteHeader(statusCode)
	lrw.ResponseData.Status = statusCode
}
