package types

import (
	"fmt"
	"io"
	"net/http"
)

type Metrics struct {
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	ID    string   `json:"id"`
	MType string   `json:"type"`
}

type GzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (gw GzipWriter) Write(b []byte) (int, error) {
	numRead, err := gw.Writer.Write(b)
	if err != nil {
		return 0, fmt.Errorf("error in write of GzipWritrer: %w", err)
	}
	return numRead, nil
}
