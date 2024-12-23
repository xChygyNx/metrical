package types

import (
	"net/http"
	"strings"
)

func IsContentEncoding(headers http.Header) bool {
	values := headers.Values("Content-Encoding")
	for _, value := range values {
		if strings.Contains(value, "gzip") {
			return true
		}
	}
	return false
}

func IsApplicationJSON(headers http.Header) bool {
	values := headers.Values("Content-Type")
	for _, value := range values {
		if strings.Contains(value, "application/json") {
			return true
		}
	}
	return false
}
