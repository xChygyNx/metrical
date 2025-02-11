package types

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	hashHeader = "HashSHA256"
)

func IsAcceptEncoding(headers http.Header) bool {
	values := headers.Values("Accept-Encoding")
	for _, value := range values {
		if strings.Contains(value, "gzip") {
			return true
		}
	}
	return false
}

func IsContentEncoding(headers http.Header) bool {
	values := headers.Values("Content-Encoding")
	for _, value := range values {
		if strings.Contains(value, "gzip") {
			return true
		}
	}
	return false
}

func IsCompressData(headers http.Header) bool {
	values := headers.Values("Content-Type")
	for _, value := range values {
		if strings.Contains(value, "application/json") ||
			strings.Contains(value, "text/html") {
			return true
		}
	}
	return false
}

func CheckHashSum(resp *http.Request) (err error) {
	if len(resp.Header.Values(hashHeader)) > 0 {
		headerHashSum := resp.Header.Values(hashHeader)[0]
		log.Printf("Sever response hashSum from header: %v\n", headerHashSum)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error in read response Body in server: %w", err)
		}
		defer func() {
			err = resp.Body.Close()
		}()

		bodyHashSumBytes := sha256.Sum256(body)
		bodyHashSumStr := base64.StdEncoding.EncodeToString(bodyHashSumBytes[:])

		log.Printf("Server response hashSum from body: %v\n", bodyHashSumStr)
		if headerHashSum != bodyHashSumStr {
			log.Printf("not match hash sum of response body and hash sum from header in server\n"+
				"%s\n%s", headerHashSum, bodyHashSumStr)
			return errors.New("not match hash sum of response body and hash sum from header in agent")
		}
	}
	return
}
