package server

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/xChygyNx/metrical/internal/server/types"
)

const (
	GAUGE                  = "gauge"
	COUNTER                = "counter"
	internalServerErrorMsg = "Internal server error"
	jsonContentType        = "application/json"
	textContentType        = "text/plain"
	contentType            = "Content-type"
	countGaugeMetrics      = 28
	writeHandlerErrorMsg   = "error of write data in http.ResponseWriter:"
	errorMsgWildcard       = "%s %w"
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

func pingDBHandle(dBAddress string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		db, err := sql.Open("pgx", dBAddress)
		if err != nil {
			errorMsg := fmt.Errorf("can't connect to DB videos: %w", err)
			fmt.Println(errorMsg)
			http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
			return
		}

		defer func() {
			err := db.Close()
			if err != nil {
				errorMsg := fmt.Errorf("can't close connection with DB videos: %w", err)
				fmt.Println(errorMsg)
				http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
				return
			}
		}()

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		err = db.PingContext(ctx)
		if err != nil {
			errorMsg := fmt.Errorf("can't connect to DB videos: %w", err)
			fmt.Println(errorMsg)
			http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
			return
		}

		res.WriteHeader(http.StatusOK)
	}
}

func SaveMetricHandleOld(storage *types.MemStorage, syncInfo *types.SyncInfo) http.HandlerFunc {
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
		if syncInfo.SyncFileRecord {
			err = writeMetricStorageFile(syncInfo.FileMetricStorage, storage)
			if err != nil {
				errorMsg := fmt.Errorf("failed to write metrics in file: %w", err).Error()
				fmt.Println(errorMsg)
				http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
				return
			}
		}

		res.WriteHeader(http.StatusOK)
		_, err = res.Write([]byte("OK"))
		if err != nil {
			errorMsg := fmt.Errorf(errorMsgWildcard, writeHandlerErrorMsg, err).Error()
			fmt.Println(errorMsg)
			http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
			return
		}
	}
}

func SaveMetricHandle(storage *types.MemStorage, syncInfo *types.SyncInfo) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set(contentType, jsonContentType)

		bodyByte, err := io.ReadAll(req.Body)
		defer func() {
			err = req.Body.Close()
		}()
		if err != nil {
			errorMsg := "error in read response body: " + err.Error()
			fmt.Println(errorMsg)
			http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
			return
		}
		var metricData types.Metrics

		err = json.Unmarshal(bodyByte, &metricData)
		metricName := metricData.ID
		var responseData types.Metrics
		switch metricData.MType {
		case GAUGE:
			storage.SetGauge(metricName, *metricData.Value)
			value, ok := storage.GetGauge(metricData.ID)
			if !ok {
				errorMsg := fmt.Sprintf("Value gauge metric %s don't saved", metricData.ID)
				fmt.Println(errorMsg)
				http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
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
				fmt.Println(errorMsg)
				http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
				return
			}
			responseData = types.Metrics{
				ID:    metricData.ID,
				MType: metricData.MType,
				Delta: &value,
			}
		default:
			bodyStr := string(bodyByte)
			errorMsg := "Unknown metric type, must be gauge or counter, got |" + metricData.MType +
				"|\n" + bodyStr
			http.Error(res, errorMsg, http.StatusBadRequest)
			return
		}

		if syncInfo.DB != nil {
			err = writeMetricStorageDB(syncInfo.DB, storage)
			if err != nil && err.Error() != "sql: transaction has already been committed or rolled back" {
				errorMsg := fmt.Errorf("failed to write metrics in DB: %w", err).Error()
				fmt.Println(errorMsg)
				http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
				return
			}
		} else if syncInfo.SyncFileRecord {
			err = writeMetricStorageFile(syncInfo.FileMetricStorage, storage)
			if err != nil {
				errorMsg := fmt.Errorf("failed to write metrics in file: %w", err).Error()
				http.Error(res, errorMsg, http.StatusBadRequest)
				return
			}
		}

		encodedResponseData, err := json.Marshal(responseData)
		if err != nil {
			errorMsg := fmt.Errorf("error in serialize response for send by server: %w", err).Error()
			fmt.Println(errorMsg)
			http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
			return
		}

		res.WriteHeader(http.StatusOK)
		_, err = res.Write(encodedResponseData)
		if err != nil {
			errorMsg := fmt.Errorf(errorMsgWildcard, writeHandlerErrorMsg, err).Error()
			fmt.Println(errorMsg)
			http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
			return
		}
	}
}

func SaveBatchMetricHandle(storage *types.MemStorage, syncInfo *types.SyncInfo) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set(contentType, jsonContentType)

		bodyByte, err := io.ReadAll(req.Body)
		defer func() {
			err = req.Body.Close()
		}()
		if err != nil {
			errorMsg := "error in read response body: " + err.Error()
			fmt.Println(errorMsg)
			http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
			return
		}
		metricsData := make([]types.Metrics, 0, countGaugeMetrics)

		err = json.Unmarshal(bodyByte, &metricsData)
		fmt.Printf("Unmarshalling metricsData: %v\n", metricsData)

		for _, metricData := range metricsData {
			metricName := metricData.ID
			switch metricData.MType {
			case GAUGE:
				storage.SetGauge(metricName, *metricData.Value)
				_, ok := storage.GetGauge(metricData.ID)
				if !ok {
					errorMsg := fmt.Sprintf("Value gauge metric %s don't saved", metricData.ID)
					fmt.Println(errorMsg)
					http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
					return
				}
			case COUNTER:
				storage.SetCounter(metricName, *metricData.Delta)
				_, ok := storage.GetCounter(metricData.ID)
				if !ok {
					errorMsg := fmt.Sprintf("Value counter metric %s don't saved", metricData.ID)
					fmt.Println(errorMsg)
					http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
					return
				}
			default:
				bodyStr := string(bodyByte)
				errorMsg := "Unknown metric type, must be gauge or counter, got |" + metricData.MType +
					"|\n" + bodyStr
				http.Error(res, errorMsg, http.StatusBadRequest)
				return
			}
		}

		if syncInfo.DB != nil {
			err = writeMetricStorageDB(syncInfo.DB, storage)
			if err != nil && err.Error() != "sql: transaction has already been committed or rolled back" {
				errorMsg := fmt.Errorf("failed to write metrics in DB: %w", err).Error()
				fmt.Println(errorMsg)
				http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
				return
			}
		} else if syncInfo.SyncFileRecord {
			err = writeMetricStorageFile(syncInfo.FileMetricStorage, storage)
			if err != nil {
				errorMsg := fmt.Errorf("failed to write metrics in file: %w", err).Error()
				http.Error(res, errorMsg, http.StatusBadRequest)
				return
			}
		}

		encodedResponseData, err := json.Marshal(metricsData)
		if err != nil {
			errorMsg := fmt.Errorf("error in serialize response for send by server: %w", err).Error()
			fmt.Println(errorMsg)
			http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
			return
		}

		res.WriteHeader(http.StatusOK)
		_, err = res.Write(encodedResponseData)
		if err != nil {
			errorMsg := fmt.Errorf(errorMsgWildcard, writeHandlerErrorMsg, err).Error()
			fmt.Println(errorMsg)
			http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
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
				errorMsg := fmt.Errorf("error in format integer from receive data: %w", err).Error()
				fmt.Println(errorMsg)
				http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
				return
			}
		case float64:
			_, err := res.Write([]byte(strconv.FormatFloat(metricValue, 'f', -1, 64)))
			if err != nil {
				errorMsg := fmt.Errorf("error in format float from receive data: %w", err).Error()
				fmt.Println(errorMsg)
				http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
				return
			}
		}

		res.WriteHeader(http.StatusOK)
	}
}

func GetJSONMetricHandle(storage *types.MemStorage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set(contentType, jsonContentType)
		bodyByte, err := io.ReadAll(req.Body)
		defer func() {
			err = req.Body.Close()
		}()
		if err != nil {
			errorMsg := fmt.Errorf("error in read response body: %w", err).Error()
			fmt.Println(errorMsg)
			http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
			return
		}
		var reqJSON types.Metrics

		requestDecoder := json.NewDecoder(bytes.NewBuffer(bodyByte))
		err = requestDecoder.Decode(&reqJSON)
		if err != nil {
			errorMsg := fmt.Errorf("error in  decode response body: %w", err).Error()
			fmt.Println(errorMsg)
			http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
			return
		}
		mType := reqJSON.MType
		switch mType {
		case GAUGE:
			value, ok := storage.GetGauge(reqJSON.ID)
			if !ok {
				errorMsg := fmt.Sprintf("Gauge metric %s don't saved", reqJSON.ID)
				http.Error(res, errorMsg, http.StatusNotFound)
				return
			}
			reqJSON.Value = &value
		case COUNTER:
			delta, ok := storage.GetCounter(reqJSON.ID)
			if !ok {
				errorMsg := fmt.Sprintf("Counter metric %s don't saved", reqJSON.ID)
				http.Error(res, errorMsg, http.StatusNotFound)
				return
			}
			reqJSON.Delta = &delta
		}
		responseData, err := json.Marshal(reqJSON)
		if err != nil {
			errorMsg := fmt.Errorf("error in serialize response for send by server: %w", err).Error()
			fmt.Println(errorMsg)
			http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusOK)
		_, err = res.Write(responseData)
		if err != nil {
			errorMsg := fmt.Errorf("error in write body of response: %w", err).Error()
			fmt.Println(errorMsg)
			http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
			return
		}
	}
}

func ListMetricHandle(storage *types.MemStorage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add(contentType, "text/html")

		metricsInfo := map[string]map[string]string{
			"Gauges":   storage.GetGauges(),
			"Counters": storage.GetCounters(),
		}
		metricInfoStr, err := json.Marshal(metricsInfo)
		if err != nil {
			errorMsg := fmt.Errorf("error in serialize of metrics storage: %w", err).Error()
			fmt.Println(errorMsg)
			http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
			return
		}

		res.WriteHeader(http.StatusOK)
		_, err = res.Write(metricInfoStr)
		if err != nil {
			errorMsg := fmt.Errorf(errorMsgWildcard, writeHandlerErrorMsg, err).Error()
			fmt.Println(errorMsg)
			http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
			return
		}
	}
}

func GzipHandler(internal http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		resWriter := w

		if types.IsCompressData(req.Header) && types.IsContentEncoding(req.Header) {
			gzipReader, err := types.NewGzipReader(req.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			req.Body = gzipReader
			defer func() {
				err := gzipReader.Close()
				if err != nil {
					errorMsg := fmt.Errorf("error in close gzipReader: %w", err).Error()
					fmt.Println(errorMsg)
					http.Error(w, internalServerErrorMsg, http.StatusInternalServerError)
					return
				}
			}()
		}
		if types.IsAcceptEncoding(req.Header) {
			writer := types.NewGzipWriter(w)
			resWriter = writer
			defer func() {
				err := writer.Close()
				if err != nil {
					errorMsg := fmt.Errorf("error in close gzipWriter: %w", err)
					http.Error(w, errorMsg.Error(), http.StatusInternalServerError)
					return
				}
			}()
		}
		internal.ServeHTTP(resWriter, req)
	})
}
