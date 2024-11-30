package senders

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func SendGauge(client *http.Client, sendInfo []uint8) error {
	var mapInfo map[string]float64
	fmt.Println(string(sendInfo))
	err := json.Unmarshal([]byte(sendInfo), &mapInfo)
	if err != nil {
		return err
	}
	for attr, value := range mapInfo {
		urlString := "http://localhost/update/gauge/" + attr + "/" + strconv.FormatFloat(value, 'f', -1, 64) + ":8080"
		req, err := http.NewRequest(http.MethodPost, urlString, bytes.NewBuffer([]byte(sendInfo)))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "text/plain")
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
		err = resp.Body.Close()
		if err != nil {
			io.Copy(os.Stdout, bytes.NewReader([]byte(err.Error())))
		}
	}
	return nil
}

func SendCounter(client *http.Client, pollCount int) error {
	counterPath := "http://localhost/update/counter/PollCount/" + strconv.Itoa(pollCount) + ":8080"
	req, err := http.NewRequest(http.MethodPost, counterPath, bytes.NewBuffer([]byte(strconv.Itoa(pollCount))))
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	return nil
}
