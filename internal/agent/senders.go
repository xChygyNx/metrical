package agent

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/sethgrid/pester"

	"github.com/xChygyNx/metrical/internal/server/types"
)

const (
	contentType          = "Content-Type"
	contentTypeValue     = "application/json"
	contentEncoding      = "Content-Encoding"
	contentEncodingValue = "gzip"
	countGaugeMetrics    = 28
	getIPErrorText       = "error in get own IP address: %w"
	realIPHeader         = "X-Real-IP"
	responseStatusMsg    = "response Status: "
	responseHeadersMsg   = "response Headers: "
	responseBodyMsg      = "response Body: "
)

// GetOutboundIP возвращает IP адрес хоста, на котором запущем агент.
func GetOutboundIP() (ip net.IP, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, fmt.Errorf("error in set UDP connection: %w", err)
	}
	defer func() {
		err = conn.Close()
	}()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}

// SendGauge отправляет по одной собранные метрики типа gauge на сервер. При ошибке отправки
// повторяет отправку заданное в системеколичество раз (для смены данного параметра требуется
// обратиться к администратору системы).
//
// Deprecated: используйте BatchSendGauge.
func SendGauge(client *pester.Client, sendInfo map[string]float64, config *Config) (err error) {
	hostAddr := config.HostPort
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

		//compre:q nil {
		//	return fmt.Errorf("error in compress gauge metrics: %w", err)
		//}

		req, err := http.NewRequest(http.MethodPost, urlString, bytes.NewBuffer(jsonString))
		if err != nil {
			return fmt.Errorf("failed to create http Request: %w", err)
		}
		req.Header.Set(contentType, contentTypeValue)
		req.Header.Set(contentEncoding, contentEncodingValue)
		ip, err := GetOutboundIP()
		if err != nil {
			return fmt.Errorf(getIPErrorText, err)
		}
		req.Header.Set(realIPHeader, ip.String())
		resp, err := client.Do(req)
		if err != nil && !errors.Is(err, io.EOF) {
			return fmt.Errorf("failed to send http Request by http Client: %w", err)
		}
		defer func() {
			err = resp.Body.Close()
		}()

		log.Println(responseStatusMsg, resp.Status)
		log.Println(responseHeadersMsg, resp.Header)
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

// SendCounter отправляет по одной собранные метрики типа counter на сервер. При ошибке отправки
// повторяет отправку заданное в системеколичество раз (для смены данного параметра требуется
// обратиться к администратору системы).
//
// Deprecated: используйте BatchSendCounter.
func SendCounter(client *pester.Client, pollCount int, config *Config) (err error) {
	hostAddr := config.HostPort
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
	// compressJSON, err := compress(jsonString)
	// if err != nil {
	//	return fmt.Errorf("error in compress counter metrics: %w", err)
	// }
	req, err := http.NewRequest(http.MethodPost, counterPath, bytes.NewBuffer(jsonString))
	if err != nil {
		return
	}
	req.Header.Set(contentType, contentTypeValue)
	req.Header.Set(contentEncoding, contentEncodingValue)
	ip, err := GetOutboundIP()
	if err != nil {
		return fmt.Errorf(getIPErrorText, err)
	}
	req.Header.Set(realIPHeader, ip.String())
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		err = resp.Body.Close()
	}()

	log.Println(responseStatusMsg, resp.Status)
	log.Println(responseHeadersMsg, resp.Header)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	log.Println(responseBodyMsg, string(body))
	return
}

// BatchSendGauge собирает метрики системы тапа gauge в один JSON и отправляет их на сервер. При ошибке отправки
// повторяет отправку заданное в системеколичество раз (для смены данного параметра требуется
// обратиться к администратору системы).
func BatchSendGauge(client *pester.Client, sendInfo map[string]float64, config *Config) (err error) {
	hostAddr := config.HostPort
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
		return fmt.Errorf("error in serialize json for send gauge metric: %w", err)
	}

	if config.RSAPublicKey != "" {
		jsonString, err = encodeDataRSA(jsonString, config.RSAPublicKey)
		if err != nil {
			return fmt.Errorf("error in encode send data: %w", err)
		}
		fmt.Printf("Cipher data: %s", string(jsonString))
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
	ip, err := GetOutboundIP()
	if err != nil {
		return fmt.Errorf(getIPErrorText, err)
	}
	req.Header.Set(realIPHeader, ip.String())
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send http Request by http Client: %w", err)
	}
	defer func() {
		err = resp.Body.Close()
	}()

	log.Println(responseStatusMsg, resp.Status)
	log.Println(responseHeadersMsg, resp.Header)
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

// BatchSendCounter собирает метрики системы тапа counter в один JSON и отправляет их на сервер. При ошибке отправки
// повторяет отправку заданное в системеколичество раз (для смены данного параметра требуется
// обратиться к администратору системы).
func BatchSendCounter(client *pester.Client, pollCount int, config *Config) (err error) {
	hostAddr := config.HostPort
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
	if config.RSAPublicKey != "" {
		jsonString, err = encodeDataRSA(jsonString, config.RSAPublicKey)
		if err != nil {
			return fmt.Errorf("error in encode send data: %w", err)
		}
		fmt.Printf("Cipher data: %s", string(jsonString))
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
	ip, err := GetOutboundIP()
	if err != nil {
		return fmt.Errorf(getIPErrorText, err)
	}
	req.Header.Set(realIPHeader, ip.String())
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		err = resp.Body.Close()
	}()

	log.Println(responseStatusMsg, resp.Status)
	log.Println(responseHeadersMsg, resp.Header)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	log.Println(responseBodyMsg, string(body))
	return
}
