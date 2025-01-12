package types

import (
	"net/http"
	"strings"
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
