package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/xChygyNx/metrical/internal/server/types"
)

const (
	contentType   = "application/json"
	contentEncode = "gzip"
)

func SendGauge(client *http.Client, sendInfo map[string]float64, hostAddr HostPort) (err error) {
	iterationLogic := func(attr string, value float64) (err error) {
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
		req.Header.Set("Content-Type", contentType)
		req.Header.Set("Content-Encoding", contentEncode)
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to send http Request by http Client: %w", err)
		}
		defer func() {
			err = resp.Body.Close()
		}()

		log.Println("response Status:", resp.Status)
		log.Println("response Headers:", resp.Header)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error in read response body: %w", err)
		}
		log.Println("response Body:", string(body))
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

func SendCounter(client *http.Client, pollCount int, hostAddr HostPort) (err error) {
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
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Content-Encoding", contentEncode)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		err = resp.Body.Close()
	}()

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	log.Println("response Body:", string(body))
	return
}

func BatchSendGauge(client *http.Client, sendInfo map[string]float64, hostAddr HostPort) (err error) {
	var sendData []types.Metrics
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
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Content-Encoding", contentEncode)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send http Request by http Client: %w", err)
	}
	defer func() {
		err = resp.Body.Close()
	}()

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error in read response body: %w", err)
	}
	log.Println("response Body:", string(body))
	err = resp.Body.Close()
	if err != nil {
		_, err = io.Copy(os.Stdout, bytes.NewReader([]byte(err.Error())))
		if err != nil {
			return fmt.Errorf("error in copy text of error in stdout: %w", err)
		}
	}

	return
}

func BatchSendCounter(client *http.Client, pollCount int, hostAddr HostPort) (err error) {
	counterPath := "http://" + hostAddr.String() + "/updates/"
	pollCount64 := int64(pollCount)
	var sendData []types.Metrics
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
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Content-Encoding", contentEncode)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		err = resp.Body.Close()
	}()

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	log.Println("response Body:", string(body))
	return
}
