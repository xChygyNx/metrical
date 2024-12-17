package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/xChygyNx/metrical/internal/server/types"
)

const (
	GAUGE         = "gauge"
	COUNTER       = "counter"
	InternalError = "Internal error"
	contentType   = "application/json"
)

func SaveMetricHandle(storage *types.MemStorage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-type", contentType)

		bodyByte, err := io.ReadAll(req.Body)
		defer func() {
			err = req.Body.Close()
		}()
		if err != nil {
			error_msg := "error in read response body: " + err.Error()
			http.Error(res, error_msg, http.StatusInternalServerError)
			return
		}
		metricData := types.Metrics{}
		var responseData types.Metrics
		if req.Header.Get("Content-Type") == contentType {
			requestDecoder := json.NewDecoder(bytes.NewBuffer(bodyByte))
			err = requestDecoder.Decode(&metricData)
			if err != nil {
				errorMsg := "error in decode response body: " + err.Error()
				log.Println(errorMsg)
				http.Error(res, errorMsg, http.StatusInternalServerError)
				return
			}
			metricName := metricData.ID

			switch metricData.MType {
			case GAUGE:
				storage.SetGauge(metricName, *metricData.Value)
				value, ok := storage.GetGauge(metricData.ID)
				if !ok {
					errorMsg := fmt.Sprintf("Value gauge metric %s don't saved", metricData.ID)
					log.Println(errorMsg)
					http.Error(res, errorMsg, http.StatusInternalServerError)
					return
				}
				responseData = types.Metrics{
					ID:    metricData.ID,
					MType: metricData.MType,
					Value: &value,
				}
			case COUNTER:
				storage.SetCounter(metricName, *metricData.Delta)
				value, ok := storage.GetCounter(metricData.ID)
				if !ok {
					errorMsg := fmt.Sprintf("Value counter metric %s don't saved", metricData.ID)
					log.Println(errorMsg)
					http.Error(res, errorMsg, http.StatusInternalServerError)
					return
				}
				responseData = types.Metrics{
					ID:    metricData.ID,
					MType: metricData.MType,
					Delta: &value,
				}
			default:
				errorMsg := "Unknown metric type, must be gauge or counter, got " + metricData.MType
				http.Error(res, errorMsg, http.StatusBadRequest)
				return
			}
		} else {
			http.Error(res, "incorrect Content-Type of request. Must be application/json", http.StatusBadRequest)
			return
		}

		encodedResponseData, err := json.Marshal(responseData)
		if err != nil {
			errorMsg := fmt.Errorf("error in serialize response for send by server: %w", err)
			log.Println(errorMsg)
			http.Error(res, errorMsg.Error(), http.StatusInternalServerError)
			return
		}

		res.WriteHeader(http.StatusOK)
		_, err = res.Write(encodedResponseData)
		if err != nil {
			log.Printf("Error of write data in http.ResponseWriter: %v\n", err)
			http.Error(res, InternalError, http.StatusInternalServerError)
			return
		}
	}
}

func getMetricValue(mType, mName string, storage *types.MemStorage) (num interface{}, ok bool) {
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

func GetMetricHandle(storage *types.MemStorage) http.HandlerFunc {
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

func ListMetricHandle(storage *types.MemStorage) http.HandlerFunc {
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
