package server

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
)

func checkHashSum(resp *http.Request) error {
	hashSum := resp.Header.Values(encodingHeader)[0]

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error in read response body: %w", err)
	}
	defer func() {
		err = resp.Body.Close()
	}()
	bodyHashSum := sha256.Sum256(body)

	if hashSum != string(bodyHashSum[:]) {
		return fmt.Errorf("error, didn't match hash sum: \n"+
			"hash sum in header: %s\nbody hash sum: %s", hashSum, string(bodyHashSum[:]))
	}
	return nil
}
