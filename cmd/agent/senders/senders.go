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

func SendGauge(client *http.Client, sendInfo []byte) error {
	var mapInfo map[string]string
	err := json.Unmarshal(sendInfo, &mapInfo)
	if err != nil {
		return err
	}
	for attr, value := range mapInfo {
		urlString := "/update/gauge/" + attr + "/" + value
		req, err := http.NewRequest(http.MethodPost, urlString, bytes.NewBuffer([]byte(sendInfo)))
		if err != nil {
			return err
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
		err = resp.Body.Close()
		if err != nil {
			io.Copy(os.Stdout, bytes.NewReader([]byte(err.Error())))
		}
	}
	return nil
}

func SendCounter(client *http.Client, pollCount int) error {
	counterPath := "/update/counter/PollCount/" + strconv.Itoa(pollCount)
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
