package types

import (
	"log"
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
	log.Println("Enter in IsContentEncoding")
	log.Printf("Request Headers in IsContentEncoding: \n %v \n", headers)
	values := headers.Values("Content-Encoding")
	for _, value := range values {
		if strings.Contains(value, "gzip") {
			log.Println("IsContentEncoding return true")
			return true
		}
	}
	log.Println("IsContentEncoding return false")
	return false
}

func IsCompressData(headers http.Header) bool {
	log.Println("Enter in IsCompressData")
	values := headers.Values("Content-Type")
	for _, value := range values {
		if strings.Contains(value, "application/json") ||
			strings.Contains(value, "text/html") {
			log.Println("IsCompressData return true")
			return true
		}
	}
	log.Println("IsCompressData return false")
	return false
}
