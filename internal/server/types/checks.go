// Package types определяет структуры для работы сервера сбора метрик, а также их
// методы для удобной работы с ними
package types

import (
	"net/http"
	"strings"
)

// IsAcceptEncoding проверка того, поддерживает ли отправитель ответы сжатые алгоритмом gzip
func IsAcceptEncoding(headers http.Header) bool {
	values := headers.Values("Accept-Encoding")
	for _, value := range values {
		if strings.Contains(value, "gzip") {
			return true
		}
	}
	return false
}

// IsContentEncoding проверка того, сжат ли пришедший запрос алгоритмом gzip
func IsContentEncoding(headers http.Header) bool {
	values := headers.Values("Content-Encoding")
	for _, value := range values {
		if strings.Contains(value, "gzip") {
			return true
		}
	}
	return false
}

// IsCompressData проверка типа контента для определения необходимости его сжатия
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
