package agent

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func SendGauge(client *http.Client, sendInfo []uint8, hostAddr HostPort) error {
	var mapInfo map[string]float64
	err := json.Unmarshal([]byte(sendInfo), &mapInfo)
	if err != nil {
		return err
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
		defer resp.Body.Close()

		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
		err = resp.Body.Close()
		if err != nil {
			_, err = io.Copy(os.Stdout, bytes.NewReader([]byte(err.Error())))
			if err != nil {
				return errors.New("can't output error in Stdout")
			}
		}
	}
	return nil
}

func SendCounter(client *http.Client, pollCount int, hostAddr HostPort) error {
	counterPath := "http://" + hostAddr.String() + "/update/counter/PollCount/" + strconv.Itoa(pollCount)
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
