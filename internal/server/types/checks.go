package types

import (
	"net/http"
	"strings"
)

func IsGzipAccepted(headers http.Header) bool {
	values := headers.Values("Content-Encoding")
	for _, value := range values {
		if strings.Contains(value, "gzip") {
			return true
		}
	}
	return false
}
