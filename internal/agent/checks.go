package agent

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
)

func CheckHashSum(resp *http.Response) (err error) {
	hashSum := resp.Header.Values(encodingHeader)[0]
	log.Printf("Hash sum from Response Header in agent: %s\n", hashSum)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error in read response body1: %w", err)
	}
	defer func() {
		err = resp.Body.Close()
	}()

	bodyHashSum := sha256.Sum256(body)
	hashSumStr := base64.StdEncoding.EncodeToString(bodyHashSum[:])
	log.Printf("Hash sum from Response Body in agent: %s\n", hashSumStr)

	if hashSum != hashSumStr {
		return fmt.Errorf("error, didn't match hash sum in agent: \n"+
			"hash sum in header: %s\nbody hash sum: %s", hashSum, hashSumStr)
	}
	return nil
}
