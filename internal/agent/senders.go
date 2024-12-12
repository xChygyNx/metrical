package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func SendGauge(client *http.Client, sendInfo []uint8, hostAddr HostPort) (err error) {
	var mapInfo map[string]float64
	err = json.Unmarshal(sendInfo, &mapInfo)
	if err != nil {
		return
	}
	iterationLogic := func(attr string, value float64) (err error) {
		urlString := "http://" + hostAddr.String() + "/update/gauge/" + attr + "/" +
			strconv.FormatFloat(value, 'f', -1, 64)
		req, err := http.NewRequest(http.MethodPost, urlString, bytes.NewBuffer(sendInfo))
		if err != nil {
			return fmt.Errorf("failed to create http Request: %w", err)
		}
		req.Header.Set("Content-Type", "text/plain")
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
	for attr, value := range mapInfo {
		err = iterationLogic(attr, value)
		if err != nil {
			return fmt.Errorf("error in SendGauge: %w", err)
		}
	}
	return
}

func SendCounter(client *http.Client, pollCount int, hostAddr HostPort) (err error) {
	counterPath := "http://" + hostAddr.String() + "/update/counter/PollCount/" + strconv.Itoa(pollCount)
	req, err := http.NewRequest(http.MethodPost, counterPath, bytes.NewBufferString(strconv.Itoa(pollCount)))
	if err != nil {
		return
	}
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
