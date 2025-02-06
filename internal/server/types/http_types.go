package types

import (
	"bytes"
	"compress/gzip"
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

type gzipWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func NewGzipWriter(w http.ResponseWriter) *gzipWriter {
	return &gzipWriter{
		ResponseWriter: w,
		Writer:         gzip.NewWriter(w),
	}
}

func (gw *gzipWriter) Write(b []byte) (int, error) {
	numRead, err := gw.Writer.Write(b)
	if err != nil {
		return 0, fmt.Errorf("error in write of GzipWritrer: %w", err)
	}
	return numRead, nil
}

func (gw *gzipWriter) WriteHeader(statusCode int) {
	if statusCode < http.StatusMultipleChoices {
		gw.Header().Set("Content-Encoding", "gzip")
	}
	gw.ResponseWriter.WriteHeader(statusCode)
}

func (gw *gzipWriter) Close() error {
	err := gw.Writer.Close()
	if err != nil {
		err = fmt.Errorf("error in close of gzipWriter.Writer: %w", err)
	}
	return err
}

type GzipReader struct {
	io.ReadCloser
	Reader *gzip.Reader
}

func NewGzipReader(r io.ReadCloser) (*GzipReader, error) {
	bodyRec, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("error in read body of request in GzipReader: %w", err)
	}
	gr, err := gzip.NewReader(bytes.NewBuffer(bodyRec))
	if err != nil {
		return nil, fmt.Errorf("error in create GzipReader: %w", err)
	}

	return &GzipReader{
		ReadCloser: r,
		Reader:     gr,
	}, nil
}

func (gr *GzipReader) Read(p []byte) (int, error) {
	n, err := gr.Reader.Read(p)
	if err != nil {
		return 0, fmt.Errorf("error in read of Gzip Reader: %w", err)
	}
	return n, nil
}

func (gr *GzipReader) Close() error {
	err := gr.ReadCloser.Close()
	if err != nil {
		return fmt.Errorf("error in close GzipReader: %w", err)
	}
	err = gr.Reader.Close()
	if err != nil {
		err = fmt.Errorf("error in close GzipReader: %w", err)
	}
	return err
}
