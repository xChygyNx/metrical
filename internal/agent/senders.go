package agent

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/sethgrid/pester"
	"github.com/xChygyNx/metrical/internal/server/types"
)

const (
	checkSumErrorMsg     = "error in checkHashSum:"
	contentType          = "Content-Type"
	contentTypeValue     = "application/json"
	contentEncoding      = "Content-Encoding"
	contentEncodingValue = "gzip"
	countGaugeMetrics    = 28
	encodingHeader       = "HashSHA256"
	errorMsgWildcard     = "%s %w"
	responseStatusMsg    = "response Status: "
	responseHeadersMsg   = "response Headers: "
	responseBodyMsg      = "response Body: "
)

func SendGauge(client *pester.Client, sendInfo map[string]float64, config *config) (err error) {
	iterationLogic := func(attr string, value float64) (err error) {
		hostAddr := config.HostAddr
		urlString := "http://" + hostAddr.String() + "/update"

		sendJSON := types.Metrics{
			ID:    attr,
			MType: "gauge",
			Value: &value,
		}
		jsonString, err := json.Marshal(sendJSON)
		if err != nil {
			return fmt.Errorf("error in serialize json for send gauge metric: %w", err)
		}

		compressJSON, err := compress(jsonString)
		if err != nil {
			return fmt.Errorf("error in compress gauge metrics: %w", err)
		}

		req, err := http.NewRequest(http.MethodPost, urlString, bytes.NewBuffer(compressJSON))
		if err != nil {
			return fmt.Errorf("failed to create http Request: %w", err)
		}
		req.Header.Set(contentType, contentTypeValue)
		req.Header.Set(contentEncoding, contentEncodingValue)
		if config.Sha256Key != "" {
			hashSum := sha256.Sum256(compressJSON)
			hashSumStr := base64.StdEncoding.EncodeToString(hashSum[:])

			req.Header.Set(encodingHeader, hashSumStr)
		}
		resp, err := client.Do(req)
		if err != nil && !errors.Is(err, io.EOF) {
			return fmt.Errorf("failed to send http Request by http Client: %w", err)
		}
		defer func() {
			err = resp.Body.Close()
		}()

		log.Println(responseStatusMsg, resp.Status)
		log.Println(responseHeadersMsg, resp.Header)
		if len(resp.Header.Values(encodingHeader)) != 0 {
			err = checkHashSum(resp)
			if err != nil {
				return fmt.Errorf(errorMsgWildcard, checkSumErrorMsg, err)
			}
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error in read response body: %w", err)
		}
		log.Println(responseBodyMsg, string(body))
		err = resp.Body.Close()
		if err != nil {
			_, err = io.Copy(os.Stdout, bytes.NewReader([]byte(err.Error())))
			if err != nil {
				return fmt.Errorf("error in copy text of error in stdout: %w", err)
			}
		}
		return
	}
	for attr, value := range sendInfo {
		err = iterationLogic(attr, value)
		if err != nil {
			return fmt.Errorf("error in SendGauge: %w", err)
		}
	}
	return
}

func SendCounter(client *pester.Client, pollCount int, config *config) (err error) {
	hostAddr := config.HostAddr
	counterPath := "http://" + hostAddr.String() + "/update"
	pollCount64 := int64(pollCount)
	sendJSON := types.Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: &pollCount64,
	}
	jsonString, err := json.Marshal(sendJSON)
	if err != nil {
		return fmt.Errorf("error in serialize json for counter metric: %w", err)
	}
	compressJSON, err := compress(jsonString)
	if err != nil {
		return fmt.Errorf("error in compress counter metrics: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, counterPath, bytes.NewBuffer(compressJSON))
	if err != nil {
		return
	}
	req.Header.Set(contentType, contentTypeValue)
	req.Header.Set(contentEncoding, contentEncodingValue)
	if config.Sha256Key != "" {
		hashSum := sha256.Sum256(compressJSON)
		hashSumSumStr := base64.StdEncoding.EncodeToString(hashSum[:])

		req.Header.Set(encodingHeader, hashSumSumStr)
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		err = resp.Body.Close()
	}()

	log.Println(responseStatusMsg, resp.Status)
	log.Println(responseHeadersMsg, resp.Header)
	if len(resp.Header.Values(encodingHeader)) != 0 {
		err = checkHashSum(resp)
		if err != nil {
			return fmt.Errorf(errorMsgWildcard, checkSumErrorMsg, err)
		}
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	log.Println(responseBodyMsg, string(body))
	return
}

func BatchSendGauge(client *pester.Client, sendInfo map[string]float64, config *config) (err error) {
	sendData := make([]types.Metrics, 0, countGaugeMetrics)
	hostAddr := config.HostAddr
	urlString := "http://" + hostAddr.String() + "/updates/"

	for attr, value := range sendInfo {
		metricInfo := types.Metrics{
			ID:    attr,
			MType: "gauge",
			Value: &value,
		}
		sendData = append(sendData, metricInfo)
	}

	jsonString, err := json.Marshal(sendData)
	if err != nil {
		return fmt.Errorf("error in serialize json for send gauge metric: %w", err)
	}

	compressJSON, err := compress(jsonString)
	if err != nil {
		return fmt.Errorf("error in compress gauge metrics: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, urlString, bytes.NewBuffer(compressJSON))
	if err != nil {
		return fmt.Errorf("failed to create http Request: %w", err)
	}
	req.Header.Set(contentType, contentTypeValue)
	req.Header.Set(contentEncoding, contentEncodingValue)
	if config.Sha256Key != "" {
		hashSum := sha256.Sum256(compressJSON)
		hashSumSumStr := base64.StdEncoding.EncodeToString(hashSum[:])

		req.Header.Set(encodingHeader, hashSumSumStr)
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send http Request by http Client: %w", err)
	}
	defer func() {
		err = resp.Body.Close()
	}()

	log.Println(responseStatusMsg, resp.Status)
	log.Println(responseHeadersMsg, resp.Header)
	if len(resp.Header.Values(encodingHeader)) != 0 {
		err = checkHashSum(resp)
		if err != nil {
			return fmt.Errorf(errorMsgWildcard, checkSumErrorMsg, err)
		}
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error in read response body: %w", err)
	}
	log.Println(responseBodyMsg, string(body))
	err = resp.Body.Close()
	if err != nil {
		_, err = io.Copy(os.Stdout, bytes.NewReader([]byte(err.Error())))
		if err != nil {
			return fmt.Errorf("error in copy text of error in stdout: %w", err)
		}
	}

	return
}

func BatchSendCounter(client *pester.Client, pollCount int, config *config) (err error) {
	hostAddr := config.HostAddr
	counterPath := "http://" + hostAddr.String() + "/updates/"
	pollCount64 := int64(pollCount)
	sendData := make([]types.Metrics, 0, 1)
	metricInfo := types.Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: &pollCount64,
	}

	sendData = append(sendData, metricInfo)
	jsonString, err := json.Marshal(sendData)
	if err != nil {
		return fmt.Errorf("error in serialize json for counter metric: %w", err)
	}
	compressJSON, err := compress(jsonString)
	if err != nil {
		return fmt.Errorf("error in compress counter metrics: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, counterPath, bytes.NewBuffer(compressJSON))
	if err != nil {
		return
	}
	req.Header.Set(contentType, contentTypeValue)
	req.Header.Set(contentEncoding, contentEncodingValue)
	if config.Sha256Key != "" {
		hashSum := sha256.Sum256(compressJSON)
		hashSumSumStr := base64.StdEncoding.EncodeToString(hashSum[:])

		req.Header.Set(encodingHeader, hashSumSumStr)
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		err = resp.Body.Close()
	}()

	log.Println(responseStatusMsg, resp.Status)
	log.Println(responseHeadersMsg, resp.Header)
	if len(resp.Header.Values(encodingHeader)) != 0 {
		err = checkHashSum(resp)
		if err != nil {
			return fmt.Errorf(errorMsgWildcard, checkSumErrorMsg, err)
		}
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	log.Println(responseBodyMsg, string(body))
	return
}
