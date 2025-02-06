package agent

import (
	"bytes"

	//"crypto/sha256"
	//"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	encodingHeader       = "HashSHA256"
	responseStatusMsg    = "response Status: "
	responseHeadersMsg   = "response Headers: "
	responseBodyMsg      = "response Body: "
)

func SendGauge(client *pester.Client, sendInfo map[string]float64, hostAddr HostPort) (err error) {
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
		log.Printf("JsonString: %s\n", jsonString)

		compressJSON, err := compress(jsonString)
		if err != nil {
			return fmt.Errorf("error in compress gauge metrics: %w", err)
		}

		log.Printf("compressJson string: %v\n", string(compressJSON))
		req, err := http.NewRequest(http.MethodPost, urlString, bytes.NewBuffer(compressJSON))
		if err != nil {
			return fmt.Errorf("failed to create http Request: %w", err)
		}
		req.Header.Set(contentType, contentTypeValue)
		req.Header.Set(contentEncoding, contentEncodingValue)

		//if config.Sha256Key != "" {
		//	hashSum := sha256.Sum256(compressJSON)
		//	hashSumStr := base64.StdEncoding.EncodeToString(hashSum[:])
		//
		//	req.Header.Set(encodingHeader, hashSumStr)

		//}
		log.Printf("Agent send request with headers: %s\n", req.Header)
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to send http Request by http Client: %w", err)
		}
		defer func() {
			err = resp.Body.Close()
		}()

		log.Println(responseStatusMsg, resp.Status)
		log.Println(responseHeadersMsg, "1", resp.Header)
		//if len(resp.Header.Values(encodingHeader)) != 0 {
		//	err = checkHashSum(resp)
		//	if err != nil {
		//		return fmt.Errorf(errorMsgWildcard, checkSumErrorMsg, err)
		//	}
		//}
		body := resp.Body
		log.Printf("resp.Body type: %T\n", body)
		log.Printf("resp.Headers: %s\n", resp.Header)
		if types.IsContentEncoding(resp.Header) {
			body, err = types.NewGzipReader(resp.Body)
			if err != nil {
				return fmt.Errorf("error in create GzipReader: %w", err)
			}
		}
		bodyData, err := io.ReadAll(body)
		if err != nil {
			return fmt.Errorf("error in read response body2: %w", err)
		}
		defer func() {
			err = body.Close()
		}()

		log.Println(responseBodyMsg, string(bodyData))
		err = body.Close()
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

func SendCounter(client *pester.Client, pollCount int, hostAddr HostPort) (err error) {
	counterPath := "http://" + hostAddr.String() + "/update"
	pollCount64 := int64(pollCount)
	sendJSON := types.Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: &pollCount64,
	}
	jsonString, err := json.Marshal(sendJSON)
	if err != nil {
		return fmt.Errorf("error in serialize json for send counter metric: %w", err)
	}

	compressJSON, err := compress(jsonString)
	if err != nil {
		return fmt.Errorf("error in compress gauge metrics: %w", err)
	}
	log.Printf("compressJson: %v\n", compressJSON)

	req, err := http.NewRequest(http.MethodPost, counterPath, bytes.NewBuffer(compressJSON))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set(contentEncoding, contentEncodingValue)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		err = resp.Body.Close()
	}()

	log.Println(responseStatusMsg, resp.Status)
	log.Println(responseHeadersMsg, resp.Header)
	//if len(resp.Header.Values(encodingHeader)) != 0 {
	//	err = checkHashSum(resp)
	//	if err != nil {
	//		return fmt.Errorf(errorMsgWildcard, checkSumErrorMsg, err)
	//	}
	//}
	body := resp.Body
	if types.IsContentEncoding(resp.Header) {
		body, err = types.NewGzipReader(resp.Body)
		if err != nil {
			return fmt.Errorf("error in create GzipReader: %w", err)
		}
	}
	bodyData, err := io.ReadAll(body)
	if err != nil {
		return
	}
	defer func() {
		err = body.Close()
	}()
	log.Println("response Body:", string(bodyData))
	return
}
