package types

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
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

func CheckHashSum(resp *http.Request, bodyByte []uint8) (err error) {
	if len(resp.Header.Values(hashHeader)) > 0 {
		headerHashSum := resp.Header.Values(hashHeader)[0]
		log.Printf("Sever response hashSum from header: %v\n", headerHashSum)

		bodyHashSumBytes := sha256.Sum256(bodyByte)
		bodyHashSumStr := base64.StdEncoding.EncodeToString(bodyHashSumBytes[:])

		log.Printf("Hash sum of response body and hash sum from header in server\n"+
			"%s\n%s", headerHashSum, bodyHashSumStr)

		if headerHashSum != bodyHashSumStr {
			return errors.New("not match hash sum of response body and hash sum from header in agent")
		}
	}
	return
}
