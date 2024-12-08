package agent

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func SendGauge(client *http.Client, sendInfo []uint8, hostAddr HostPort) (err error) {
	var mapInfo map[string]float64
	err = json.Unmarshal([]byte(sendInfo), &mapInfo)
	if err != nil {
		return
	}
	for attr, value := range mapInfo {
		urlString := "http://" + hostAddr.String() + "/update/gauge/" + attr + "/" + strconv.FormatFloat(value, 'f', -1, 64)
		req, err := http.NewRequest(http.MethodPost, urlString, bytes.NewBuffer([]byte(sendInfo)))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "text/plain")
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer func() {
			err = resp.Body.Close()
		}()

		log.Println("response Status:", resp.Status)
		log.Println("response Headers:", resp.Header)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		log.Println("response Body:", string(body))
		err = resp.Body.Close()
		if err != nil {
			_, err = io.Copy(os.Stdout, bytes.NewReader([]byte(err.Error())))
			if err != nil {
				return err
			}
		}
	}
	return
}

func SendCounter(client *http.Client, pollCount int, hostAddr HostPort) (err error) {
	counterPath := "http://" + hostAddr.String() + "/update/counter/PollCount/" + strconv.Itoa(pollCount)
	req, err := http.NewRequest(http.MethodPost, counterPath, bytes.NewBuffer([]byte(strconv.Itoa(pollCount))))
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
