package agent

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
)

func checkHashSum(resp *http.Response) (err error) {
	if len(resp.Header.Values(hashHeader)) > 0 {
		headerHashSum := resp.Header.Values(hashHeader)[0]
		log.Printf("Agent response hashSum from header: %v\n", headerHashSum)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error in read response Body in agent: %w", err)
		}
		defer func() {
			err = resp.Body.Close()
		}()

		if len(resp.Header.Values(contentEncoding)) > 0 &&
			resp.Header.Values(contentEncoding)[0] == contentEncodingValue {
			body, err = decompress(body)
		}

		bodyHashSumBytes := sha256.Sum256(body)
		bodyHashSumStr := base64.StdEncoding.EncodeToString(bodyHashSumBytes[:])

		log.Printf("Agent response hashSum from body: %v\n", bodyHashSumStr)
		if headerHashSum != bodyHashSumStr {
			return fmt.Errorf("not match hash sum of response body and hash sum from header in agent\n")
		}
	}
	return
}
