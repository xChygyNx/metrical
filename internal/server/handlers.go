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
	GAUGE           = "gauge"
	COUNTER         = "counter"
	InternalError   = "Internal error"
	jsonContentType = "application/json"
	textContentType = "text/plain"
	contentType     = "Content-type"
)

func parseGaugeMetricValue(value string) (num float64, err error) {
	num, err = strconv.ParseFloat(value, 64)
	return
}

func parseCounterMetricValue(value string) (num int64, err error) {
	num, err = strconv.ParseInt(value, 10, 64)
	return
}

func saveMetricValue(mType, mName, value string, storage *types.MemStorage) (err error) {
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

func SaveMetricHandleOld(storage *types.MemStorage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set(contentType, textContentType)

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

func SaveMetricHandle(storage *types.MemStorage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set(contentType, jsonContentType)

		bodyByte, err := io.ReadAll(req.Body)
		defer func() {
			err = req.Body.Close()
		}()
		if err != nil {
			errorMsg := "error in read response body: " + err.Error()
			http.Error(res, errorMsg, http.StatusInternalServerError)
			return
		}
		var metricData types.Metrics

		requestDecoder := json.NewDecoder(bytes.NewBuffer(bodyByte))
		err = requestDecoder.Decode(&metricData)
		if err != nil {
			errorMsg := "error in decode response body: " + err.Error()
			log.Println(errorMsg)
			http.Error(res, errorMsg, http.StatusInternalServerError)
			return
		}
		metricName := metricData.ID
		var responseData types.Metrics
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
		res.Header().Set(contentType, textContentType)
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
		res.Header().Set(contentType, textContentType)

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
