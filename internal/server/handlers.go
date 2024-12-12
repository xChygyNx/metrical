package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const (
	GAUGE         = "gauge"
	COUNTER       = "counter"
	InternalError = "Internal error"
)

func parseGaugeMetricValue(value string) (num float64, err error) {
	num, err = strconv.ParseFloat(value, 64)
	return
}

func parseCounterMetricValue(value string) (num int64, err error) {
	num, err = strconv.ParseInt(value, 10, 64)
	return
}

func saveMetricValue(mType, mName, value string, storage *memStorage) (err error) {
	switch mType {
	case GAUGE:
		var num float64
		num, err = parseGaugeMetricValue(value)
		if err != nil {
			return
		}
		storage.SetGauge(mName, num)
	case COUNTER:
		var num int64
		num, err = parseCounterMetricValue(value)
		if err != nil {
			return
		}
		storage.SetCounter(mName, num)
	}
	return
}

func SaveMetricHandle(storage *memStorage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-type", "text/plain")

		metricType := req.PathValue("mType")
		if metricType != GAUGE && metricType != COUNTER {
			errorMsg := "Unknown metric type, must be gauge or counter, got " + metricType
			http.Error(res, errorMsg, http.StatusBadRequest)
			return
		}

		metricName := req.PathValue("metric")
		metricValue := req.PathValue("value")

		err := saveMetricValue(metricType, metricName, metricValue, storage)
		if err != nil {
			errorMsg := fmt.Sprintf("Value of metric must be numeric, got %s, err: %v\n", metricValue, err)
			http.Error(res, errorMsg, http.StatusBadRequest)
			return
		}

		res.WriteHeader(http.StatusOK)
		_, err = res.Write([]byte("OK"))
		if err != nil {
			log.Printf("Error of write data in http.ResponseWriter: %v\n", err)
			http.Error(res, InternalError, http.StatusInternalServerError)
			return
		}
	}
}

func getMetricValue(mType, mName string, storage *memStorage) (num interface{}, ok bool) {
	switch mType {
	case GAUGE:
		num, ok = storage.GetGauge(mName)
		return
	case COUNTER:
		num, ok = storage.GetCounter(mName)
		return
	}
	return
}

func GetMetricHandle(storage *memStorage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-type", "text/plain")
		metricType := req.PathValue("mType")
		if metricType != GAUGE && metricType != COUNTER {
			errorMsg := "Unknown metric type, must be gauge or counter, got " + metricType
			http.Error(res, errorMsg, http.StatusBadRequest)
			return
		}

		metricName := req.PathValue("metric")
		valueInterface, ok := getMetricValue(metricType, metricName, storage)
		if !ok {
			http.Error(res, "Metric "+metricName+" not set", http.StatusNotFound)
			return
		}

		switch metricValue := valueInterface.(type) {
		case int64:
			_, err := res.Write([]byte(strconv.FormatInt(metricValue, 10)))
			if err != nil {
				log.Printf("Error in format integer from receive data: %v\n", err)
				http.Error(res, InternalError, http.StatusInternalServerError)
				return
			}
		case float64:
			_, err := res.Write([]byte(strconv.FormatFloat(metricValue, 'f', -1, 64)))
			if err != nil {
				log.Printf("Error in format float from receive data: %v\n", err)
				http.Error(res, InternalError, http.StatusInternalServerError)
				return
			}
		}

		res.WriteHeader(http.StatusOK)
	}
}

func ListMetricHandle(storage *memStorage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-type", "text/plain")

		metricsInfo := map[string]map[string]string{
			"Gauges":   storage.GetGauges(),
			"Counters": storage.GetCounters(),
		}
		metricInfoStr, err := json.Marshal(metricsInfo)
		if err != nil {
			log.Printf("Error in serialize of metrics storage: %v\n", err)
			http.Error(res, InternalError, http.StatusInternalServerError)
			return
		}

		res.WriteHeader(http.StatusOK)
		_, err = res.Write(metricInfoStr)
		if err != nil {
			log.Printf("Error of write data in http.ResponseWriter: %v\n", err)
			http.Error(res, InternalError, http.StatusInternalServerError)
			return
		}
	}
}
