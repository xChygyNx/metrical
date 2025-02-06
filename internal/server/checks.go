package server

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
)

func checkHashSum(resp *http.Request) (err error) {
	if len(resp.Header.Values(encodingHeader)) == 0 {
		return nil
	}
	hashSum := resp.Header.Values(encodingHeader)[0]
	log.Printf("Hash sum from Response Header on server: %s\n", hashSum)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error on read response body4: %w", err)
	}
	defer func() {
		err = resp.Body.Close()
	}()

	bodyHashSum := sha256.Sum256(body)
	hashSumStr := base64.StdEncoding.EncodeToString(bodyHashSum[:])
	log.Printf("Hash sum from Response Body on server: %s\n", hashSumStr)
	if hashSum != hashSumStr {
		return fmt.Errorf("error, didn't match hash sum on server: \n"+
			"hash sum in header: %s\nbody hash sum: %s", hashSum, hashSumStr)
	}
	return nil
}
