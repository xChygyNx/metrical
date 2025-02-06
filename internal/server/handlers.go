package server

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/xChygyNx/metrical/internal/server/types"
)

const (
	GAUGE   = "gauge"
	COUNTER = "counter"

	contentEncoding        = "Content-Encoding"
	contentEncodingValue   = "gzip"
	contentType            = "Content-type"
	countGaugeMetrics      = 28
	internalServerErrorMsg = "Internal server error"
	encodingHeader         = "HashSHA256"
	errorMsgWildcard       = "%s %w"
	jsonContentType        = "application/json"
	notMatchedHashSumMsg   = "Didn't match hash sums "
	retryDBWriteCount      = 4
	retryFileWriteCount    = 4
	textContentType        = "text/plain"
	writeHandlerErrorMsg   = "error of write data in http.ResponseWriter:"
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
			log.Println(errorMsg)
			http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
			return
		}

		defer func() {
			err := db.Close()
			if err != nil {
				errorMsg := fmt.Errorf("can't close connection with DB videos: %w", err)
				log.Println(errorMsg)
				http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
				return
			}
		}()

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		err = db.PingContext(ctx)
		if err != nil {
			errorMsg := fmt.Errorf("can't connect to DB videos: %w", err)
			log.Println(errorMsg)
			http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
			return
		}

		res.WriteHeader(http.StatusOK)
	}
}

func SaveMetricHandleOld(storage *types.MemStorage, handlerConf *types.HandlerConf) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set(contentType, textContentType)

		if handlerConf.Sha256Key != "" {
			err := checkHashSum(req)
			if err != nil {
				http.Error(res, notMatchedHashSumMsg, http.StatusBadRequest)
				return
			}
		}

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
		if handlerConf.SyncFileRecord {
			err = retryFileWrite(handlerConf.FileMetricStorage, storage, retryFileWriteCount)
			if err != nil {
				errorMsg := fmt.Errorf("failed to write metrics in file: %w", err).Error()
				log.Println(errorMsg)
				http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
				return
			}
		}

		res.WriteHeader(http.StatusOK)
		_, err = res.Write([]byte("OK"))
		if err != nil {
			errorMsg := fmt.Errorf(errorMsgWildcard, writeHandlerErrorMsg, err).Error()
			log.Println(errorMsg)
			http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
			return
		}
	}
}

func SaveMetricHandle(storage *types.MemStorage, handlerConf *types.HandlerConf) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log.Println("Enter SaveMetricHandle")
		res.Header().Set(contentType, jsonContentType)

		log.Printf("HandlerConfig.SHA256 in SaveMetricHandle: |%s|\n", handlerConf.Sha256Key)

		if handlerConf.Sha256Key != "" {
			err := checkHashSum(req)
			if err != nil {
				http.Error(res, notMatchedHashSumMsg, http.StatusBadRequest)
				return
			}
		}

		bodyByte, err := io.ReadAll(req.Body)
		if err != nil {
			errorMsg := "error in read response body5: " + err.Error()
			fmt.Println(errorMsg)
			http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
			return
		}
		defer func() {
			err := req.Body.Close()
			errorMsg := err.Error()
			log.Println(errorMsg)
			http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
		}()

		var metricData types.Metrics
		log.Printf("req.Body type: %T\n", req.Body)
		log.Printf("BodyByte: %s\n", string(bodyByte))
		res.Header().WriteSubset(os.Stdout, nil)
		err = json.Unmarshal(bodyByte, &metricData)
		log.Printf("Unmarshalling MetricData: %v\n", metricData)
		if err != nil {
			errorMsg := fmt.Errorf("error in unmarshaling bodyByte in SaveMetricHandle: %w", err).Error()
			log.Println(errorMsg)
			http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
			return
		}
		metricName := metricData.ID
		var responseData types.Metrics
		log.Printf("Metric name: %s, Metrice type: %s\n", metricName, metricData.MType)
		switch metricData.MType {
		case GAUGE:
			log.Println("Set Gauge")
			storage.SetGauge(metricName, *metricData.Value)
			log.Println("Set Gauge Complete")
			value, ok := storage.GetGauge(metricData.ID)
			log.Printf("Get Gauge Complete. Value: %v, ok: %v\n", value, ok)
			if !ok {
				errorMsg := fmt.Sprintf("Value gauge metric %s don't saved", metricData.ID)
				log.Println(errorMsg)
				http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
				return
			}
			responseData = types.Metrics{
				ID:    metricData.ID,
				MType: metricData.MType,
				Value: &value,
			}
			log.Printf("responseGata: %v\n", responseData)
		case COUNTER:
			storage.SetCounter(metricName, *metricData.Delta)
			value, ok := storage.GetCounter(metricData.ID)
			if !ok {
				errorMsg := fmt.Sprintf("Value counter metric %s don't saved", metricData.ID)
				log.Println(errorMsg)
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
		log.Printf("handlerConf.DB: %v\nhandlerConf.SyncFileRecord: %v\n", handlerConf.DB, handlerConf.SyncFileRecord)
		if handlerConf.DB != nil {
			err = retryDBWrite(handlerConf.DB, storage, retryDBWriteCount)
			if err != nil && err.Error() != "sql: transaction has already been committed or rolled back" {
				errorMsg := fmt.Errorf("failed to write metrics in DB: %w", err).Error()
				log.Println(errorMsg)
				http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
				return
			}
		} else if handlerConf.SyncFileRecord {
			err = retryFileWrite(handlerConf.FileMetricStorage, storage, retryFileWriteCount)
			if err != nil {
				errorMsg := fmt.Errorf("failed to write metrics in file: %w", err).Error()
				http.Error(res, errorMsg, http.StatusBadRequest)
				return
			}
		}
		log.Println("Marshall responseData")

		encodedResponseData, err := json.Marshal(responseData)
		if err != nil {
			errorMsg := fmt.Errorf("error in serialize response for send by server: %w", err).Error()
			log.Println(errorMsg)
			http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
			return
		}

		log.Println("Check header on IsAcceptEncoding")
		if types.IsAcceptEncoding(req.Header) {
			res.Header().Set(contentEncoding, contentEncodingValue)
			if err != nil {
				errorMsg := "error in compress response data: " + err.Error()
				fmt.Println(errorMsg)
				http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
				return
			}
		}

		log.Printf("handlerConf.Sha256Key: %v\n", handlerConf.Sha256Key)
		if handlerConf.Sha256Key != "" {
			hashSum := sha256.Sum256(encodedResponseData)
			hashSumStr := base64.StdEncoding.EncodeToString(hashSum[:])
			res.Header().Set(encodingHeader, hashSumStr)
		}

		log.Println("res.WriteHeader")
		res.WriteHeader(http.StatusOK)
		log.Println("res.WriteHeader Complete")
		log.Printf("encodedResponseData: %v\n", encodedResponseData)
		_, err = res.Write(encodedResponseData)
		log.Println("res.Write Complete")
		if err != nil {
			errorMsg := fmt.Errorf(errorMsgWildcard, writeHandlerErrorMsg, err).Error()
			log.Println(errorMsg)
			http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
			return
		}
	}
}

func SaveBatchMetricHandle(storage *types.MemStorage, handlerConf *types.HandlerConf) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set(contentType, jsonContentType)

		if handlerConf.Sha256Key != "" {
			err := checkHashSum(req)
			if err != nil {
				http.Error(res, notMatchedHashSumMsg, http.StatusBadRequest)
				return
			}
		}

		bodyByte, err := io.ReadAll(req.Body)

		if err != nil {
			errorMsg := "error in read response body6: " + err.Error()
			log.Println(errorMsg)
			http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
			return
		}
		defer func() {
			err := req.Body.Close()
			errorMsg := err.Error()
			log.Println(errorMsg)
			http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
		}()

		metricsData := make([]types.Metrics, 0, countGaugeMetrics)
		err = json.Unmarshal(bodyByte, &metricsData)
		if err != nil {
			errorMsg := fmt.Errorf("error in unmarshaling bodyByte in SaveBatchMetricHandle: %w", err).Error()
			log.Println(errorMsg)
			http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
			return
		}
		log.Printf("Unmarshalling metricsData: %v\n", metricsData)

		for _, metricData := range metricsData {
			metricName := metricData.ID
			switch metricData.MType {
			case GAUGE:
				storage.SetGauge(metricName, *metricData.Value)
				_, ok := storage.GetGauge(metricData.ID)
				if !ok {
					errorMsg := fmt.Sprintf("Value gauge metric %s don't saved", metricData.ID)
					log.Println(errorMsg)
					http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
					return
				}
			case COUNTER:
				storage.SetCounter(metricName, *metricData.Delta)
				_, ok := storage.GetCounter(metricData.ID)
				if !ok {
					errorMsg := fmt.Sprintf("Value counter metric %s don't saved", metricData.ID)
					log.Println(errorMsg)
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

		if handlerConf.DB != nil {
			err = retryDBWrite(handlerConf.DB, storage, retryDBWriteCount)
			if err != nil && err.Error() != "sql: transaction has already been committed or rolled back" {
				errorMsg := fmt.Errorf("failed to write metrics in DB: %w", err).Error()
				log.Println(errorMsg)
				http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
				return
			}
		} else if handlerConf.SyncFileRecord {
			err = retryFileWrite(handlerConf.FileMetricStorage, storage, retryFileWriteCount)
			if err != nil {
				errorMsg := fmt.Errorf("failed to write metrics in file: %w", err).Error()
				http.Error(res, errorMsg, http.StatusBadRequest)
				return
			}
		}

		encodedResponseData, err := json.Marshal(metricsData)
		if err != nil {
			errorMsg := fmt.Errorf("error in serialize response for send by server: %w", err).Error()
			log.Println(errorMsg)
			http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
			return
		}

		if types.IsAcceptEncoding(req.Header) {
			res.Header().Set(contentEncoding, contentEncodingValue)
			if err != nil {
				errorMsg := "error in compress response data: " + err.Error()
				fmt.Println(errorMsg)
				http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
				return
			}
		}

		if handlerConf.Sha256Key != "" {
			hashSum := sha256.Sum256(encodedResponseData)
			hashSumStr := base64.StdEncoding.EncodeToString(hashSum[:])
			req.Header.Set(encodingHeader, hashSumStr)
		}

		res.WriteHeader(http.StatusOK)
		_, err = res.Write(encodedResponseData)
		if err != nil {
			errorMsg := fmt.Errorf(errorMsgWildcard, writeHandlerErrorMsg, err).Error()
			log.Println(errorMsg)
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
				log.Println(errorMsg)
				http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
				return
			}
		case float64:
			_, err := res.Write([]byte(strconv.FormatFloat(metricValue, 'f', -1, 64)))
			if err != nil {
				errorMsg := fmt.Errorf("error in format float from receive data: %w", err).Error()
				log.Println(errorMsg)
				http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
				return
			}
		}

		res.WriteHeader(http.StatusOK)
	}
}

func GetJSONMetricHandle(storage *types.MemStorage, handlerConf *types.HandlerConf) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set(contentType, jsonContentType)
		if handlerConf.Sha256Key != "" {
			err := checkHashSum(req)
			if err != nil {
				http.Error(res, notMatchedHashSumMsg, http.StatusBadRequest)
				return
			}
		}
		bodyByte, err := io.ReadAll(req.Body)

		if err != nil {
			errorMsg := fmt.Errorf("error in read response body7: %w", err).Error()
			log.Println(errorMsg)
			http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
			return
		}
		defer func() {
			err := req.Body.Close()
			errorMsg := err.Error()
			log.Println(errorMsg)
			http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
		}()
		var reqJSON types.Metrics

		requestDecoder := json.NewDecoder(bytes.NewBuffer(bodyByte))
		err = requestDecoder.Decode(&reqJSON)
		if err != nil {
			errorMsg := fmt.Errorf("error in  decode response body: %w", err).Error()
			log.Println(errorMsg)
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
			log.Println(errorMsg)
			http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
			return
		}

		if handlerConf.Sha256Key != "" {
			hashSum := sha256.Sum256(responseData)
			hashSumStr := base64.StdEncoding.EncodeToString(hashSum[:])
			req.Header.Set(encodingHeader, hashSumStr)
		}

		res.WriteHeader(http.StatusOK)
		_, err = res.Write(responseData)
		if err != nil {
			errorMsg := fmt.Errorf("error in write body of response: %w", err).Error()
			log.Println(errorMsg)
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
			log.Println(errorMsg)
			http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
			return
		}

		res.WriteHeader(http.StatusOK)
		_, err = res.Write(metricInfoStr)
		if err != nil {
			errorMsg := fmt.Errorf(errorMsgWildcard, writeHandlerErrorMsg, err).Error()
			log.Println(errorMsg)
			http.Error(res, internalServerErrorMsg, http.StatusInternalServerError)
			return
		}
	}
}

func GzipHandler(internal http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		resWriter := w

		log.Printf("Request headers in GzipHandler: %s\n", req.Header)
		log.Printf("Internal in GzipHandler: %T\n", internal)
		//res, _ := io.ReadAll(req.Body)
		//defer func() {
		//	err := req.Body.Close()
		//	errorMsg := err.Error()
		//	log.Println(errorMsg)
		//	http.Error(w, internalServerErrorMsg, http.StatusInternalServerError)
		//}()
		//
		//readCloser := io.NopCloser(bytes.NewBuffer(res))

		if types.IsCompressData(req.Header) && types.IsContentEncoding(req.Header) {
			//log.Printf("BodyByte in middle: %v\n", string(req.Body))
			log.Println("NewGzipReader")
			gzipReader, err := types.NewGzipReader(req.Body)
			log.Println("Complete NewGzipReader")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			log.Println("Set gzipReader in req.Body")
			req.Body = gzipReader
			defer func() {
				err := gzipReader.Close()
				if err != nil {
					errorMsg := fmt.Errorf("error in close gzipReader: %w", err).Error()
					log.Println(errorMsg)
					http.Error(w, internalServerErrorMsg, http.StatusInternalServerError)
					return
				}
			}()
		}
		//log.Printf("type of req.Body in GzipHandler: %T\n", req.Body)
		//bodyInfo, err := io.ReadAll(req.Body)
		//if err != nil {
		//	log.Printf("err in read body in GzipHandler: %v\n", err.Error())
		//}
		//log.Printf("Read body in GzipReader after replacment request Body: %v\n", bodyInfo)

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
		log.Printf("request Body type before exit GzipReader: %T\n", req.Body)
		internal.ServeHTTP(resWriter, req)
		log.Println("Internal ServeHTTP complete")
	})
}
