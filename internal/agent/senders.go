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
	"runtime"
	"strings"

	"github.com/sethgrid/pester"
	"github.com/xChygyNx/metrical/internal/server/types"

	"github.com/xChygyNx/metrical/internal/agent/config"
)

const (
	contentType          = "Content-Type"
	contentTypeValue     = "application/json"
	contentEncoding      = "Content-Encoding"
	contentEncodingValue = "gzip"
	hashHeader           = "HashSHA256"
	countGaugeMetrics    = 28
	responseStatusMsg    = "response Status: "
	responseHeadersMsg   = "response Headers: "
	responseBodyMsg      = "response Body: "
)

func getFuncName() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	fullFuncName := runtime.FuncForPC(pc).Name()
	funcName := fullFuncName[strings.LastIndex(fullFuncName, ".")+1:]
	return funcName
}

func hashSendData(data []byte) string {
	hashSum := sha256.Sum256(data)
	dataHash := base64.StdEncoding.EncodeToString(hashSum[:])
	return dataHash
}

func SendGauge(client *pester.Client, sendInfo map[string]float64, agentConfig *config.AgentConfig) (err error) {
	iterationLogic := func(attr string, value float64) (err error) {
		hostAddr := agentConfig.HostAddr
		urlString := "http://" + hostAddr.String() + "/update"

		sendJSON := types.Metrics{
			ID:    attr,
			MType: "gauge",
			Value: &value,
		}
		jsonString, err := json.Marshal(sendJSON)
		if err != nil {
			return fmt.Errorf("error in serialize json for send gauge metric in %s: %w", getFuncName(), err)
		}

		compressJSON, err := compress(jsonString)
		if err != nil {
			return fmt.Errorf("error in compress gauge metrics: %w", err)
		}

		req, err := http.NewRequest(http.MethodPost, urlString, bytes.NewBuffer(compressJSON))
		if err != nil {
			return fmt.Errorf("failed to create http Request in SendGauge: %w", err)
		}
		req.Header.Set(contentType, contentTypeValue)
		req.Header.Set(contentEncoding, contentEncodingValue)
		if agentConfig.Sha256Key != "" {
			hash := hashSendData(jsonString)
			log.Printf("Agent hash data in SendGauge: %s", jsonString)
			log.Printf("Agent set hash in header in SendGauge: %s", hash)
			req.Header.Set(hashHeader, hash)
		}

		resp, err := client.Do(req)
		if err != nil && !errors.Is(err, io.EOF) {
			return fmt.Errorf("failed to send http Request by http Client in %s: %w", getFuncName(), err)
		}

		log.Println(responseStatusMsg, resp.Status)
		log.Println(responseHeadersMsg, resp.Header)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error in read response body in SendGauge: %w", err)
		}

		defer func() {
			err = resp.Body.Close()
		}()
		err = checkHashSum(resp)
		if err != nil {
			return fmt.Errorf("error CheckHashSum in SendGauge: %w", err)
		}

		if len(resp.Header.Values(contentEncoding)) > 0 &&
			resp.Header.Values(contentEncoding)[0] == contentEncodingValue {
			body, err = decompress(body)
			if err != nil {
				return fmt.Errorf("error of decompress response body in SendGauge: %w", err)
			}
		}
		log.Println(responseBodyMsg, string(body))
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

func SendCounter(client *pester.Client, pollCount int, agentConfig *config.AgentConfig) (err error) {
	hostAddr := agentConfig.HostAddr
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
		return fmt.Errorf("failed to create http Request in SendCounter: %w", err)
	}
	req.Header.Set(contentType, contentTypeValue)
	req.Header.Set(contentEncoding, contentEncodingValue)
	if agentConfig.Sha256Key != "" {
		hash := hashSendData(jsonString)
		log.Printf("Agent hash data in SendCounter: %s", jsonString)
		log.Printf("Agent set hash in header in SendCounter: %s", hash)
		req.Header.Set(hashHeader, hash)
	}

	resp, err := client.Do(req)
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("failed to send http Request by http Client in SendCounter: %w", err)
	}

	log.Println(responseStatusMsg, resp.Status)
	log.Println(responseHeadersMsg, resp.Header)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error in read response body in SendCounter: %w", err)
	}

	defer func() {
		err = resp.Body.Close()
	}()
	err = checkHashSum(resp)
	if err != nil {
		return fmt.Errorf("error CheckHashSum in SendCounter: %w", err)
	}

	if len(resp.Header.Values(contentEncoding)) > 0 &&
		resp.Header.Values(contentEncoding)[0] == contentEncodingValue {
		body, err = decompress(body)
		if err != nil {
			return fmt.Errorf("error of decompress response body in SendCounter: %w", err)
		}
	}
	log.Println(responseBodyMsg, string(body))
	return
}

func BatchSendGauge(client *pester.Client, sendInfo map[string]float64, agentConfig *config.AgentConfig) (err error) {
	hostAddr := agentConfig.HostAddr
	sendData := make([]types.Metrics, 0, countGaugeMetrics)
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
		return fmt.Errorf("error in serialize json for batch send gauge metric: %w", err)
	}

	compressJSON, err := compress(jsonString)
	if err != nil {
		return fmt.Errorf("error in compress batch gauge metrics: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, urlString, bytes.NewBuffer(compressJSON))
	if err != nil {
		return fmt.Errorf("failed to create http Request in BatchSendGauge: %w", err)
	}

	req.Header.Set(contentType, contentTypeValue)
	req.Header.Set(contentEncoding, contentEncodingValue)
	if agentConfig.Sha256Key != "" {
		hash := hashSendData(jsonString)
		log.Printf("Agent hash data in BatchSendGauge: %s", jsonString)
		log.Printf("Agent set hash in header in BatchSendGauge: %s", hash)
		req.Header.Set(hashHeader, hash)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send http Request by http Client in BatchSendGauge: %w", err)
	}

	log.Println(responseStatusMsg, resp.Status)
	log.Println(responseHeadersMsg, resp.Header)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error in read response body in BatchSendGauge: %w", err)
	}
	log.Println(responseBodyMsg, string(body))
	err = resp.Body.Close()
	if err != nil {
		_, err = io.Copy(os.Stdout, bytes.NewReader([]byte(err.Error())))
		if err != nil {
			return fmt.Errorf("error in copy text of error in stdout: %w", err)
		}
	}

	err = checkHashSum(resp)
	if err != nil {
		return fmt.Errorf("error CheckHashSum in BatchSendGauge: %w", err)
	}

	if len(resp.Header.Values(contentEncoding)) > 0 &&
		resp.Header.Values(contentEncoding)[0] == contentEncodingValue {
		body, err = decompress(body)
		if err != nil {
			return fmt.Errorf("error of decompress response body in BatchSendGauge: %w", err)
		}
	}
	log.Println(responseBodyMsg, string(body))

	return
}

func BatchSendCounter(client *pester.Client, pollCount int, agentConfig *config.AgentConfig) (err error) {
	hostAddr := agentConfig.HostAddr
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
		return fmt.Errorf("failed to create http Request in BatchSendCounter: %w", err)
	}

	req.Header.Set(contentType, contentTypeValue)
	req.Header.Set(contentEncoding, contentEncodingValue)
	if agentConfig.Sha256Key != "" {
		hash := hashSendData(jsonString)
		log.Printf("Agent hash data in BatchSendGauge: %s", jsonString)
		log.Printf("Agent set hash in header in BatchSendGauge: %s", hash)
		req.Header.Set(hashHeader, hash)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send http Request by http Client in %s: %w", getFuncName(), err)
	}
	defer func() {
		err = resp.Body.Close()
	}()

	log.Println(responseStatusMsg, resp.Status)
	log.Println(responseHeadersMsg, resp.Header)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error in read response body in BatchSendCounter: %w", err)
	}

	defer func() {
		err = resp.Body.Close()
	}()
	err = checkHashSum(resp)
	if err != nil {
		return fmt.Errorf("error CheckHashSum in BatchSendCounter: %w", err)
	}

	if len(resp.Header.Values(contentEncoding)) > 0 &&
		resp.Header.Values(contentEncoding)[0] == contentEncodingValue {
		body, err = decompress(body)
		if err != nil {
			return fmt.Errorf("error of decompress response body in BatchSendCounter: %w", err)
		}
	}

	log.Println(responseBodyMsg, string(body))
	return
}
