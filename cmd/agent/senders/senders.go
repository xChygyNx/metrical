package senders

import (
	"net/http"
	"bytes"
	"fmt"
	"io"
	"strconv"
)

func SendGauge (client *http.Client, sendInfo []byte) error {
	req, err := http.NewRequest(http.MethodPost, "/update/gauge/all/123", bytes.NewBuffer([]byte(sendInfo)))
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
	return nil
}

func SendCounter (client *http.Client, pollCount int) error {
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