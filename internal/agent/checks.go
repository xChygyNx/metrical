package agent

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
)

func checkHashSum(resp *http.Response) error {
	hashSum := resp.Header.Values(encodingHeader)[0]

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error in read response body: %w", err)
	}
	bodyHashSum := sha256.Sum256(body)

	if hashSum != string(bodyHashSum[:]) {
		return fmt.Errorf("error, didn't match hash sum: \n"+
			"hash sum in header: %s\nbody hash sum: %s", hashSum, string(bodyHashSum[:]))
	}
	return nil
}
