package types

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
)

// Metrics структура для хранения значений метрики.
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

// NewGzipWriter оборачивает ResponseWriter в "записыватель" использующий gzip сжатие.
func NewGzipWriter(w http.ResponseWriter) *gzipWriter {
	return &gzipWriter{
		ResponseWriter: w,
		Writer:         gzip.NewWriter(w),
	}
}

// Write сжимает и записывает данные в тело ответа.
func (gw *gzipWriter) Write(b []byte) (int, error) {
	numRead, err := gw.Writer.Write(b)
	if err != nil {
		return 0, fmt.Errorf("error in write of GzipWritrer: %w", err)
	}
	return numRead, nil
}

// WriteHeader записывает в заголовок ответа Content-Encoding значение gzip.
func (gw *gzipWriter) WriteHeader(statusCode int) {
	if statusCode < http.StatusMultipleChoices {
		gw.Header().Set("Content-Encoding", "gzip")
	}
	gw.ResponseWriter.WriteHeader(statusCode)
}

// Close закрывает "записыватель" gzipWriter.
func (gw *gzipWriter) Close() error {
	err := gw.Writer.Close()
	if err != nil {
		err = fmt.Errorf("error in close of gzipWriter.Writer: %w", err)
	}
	return err
}

type gzipReader struct {
	io.ReadCloser
	reader *gzip.Reader
}

// NewGzipReader возвращает "считыватель" способный читать сжатые с помощью алгоритма gzip данные.
func NewGzipReader(r io.ReadCloser) (*gzipReader, error) {
	bodyRec, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("error in read body of request in gzipReader: %w", err)
	}
	gr, err := gzip.NewReader(bytes.NewBuffer(bodyRec))
	if err != nil {
		return nil, fmt.Errorf("error in create gzipReader: %w", err)
	}

	return &gzipReader{
		ReadCloser: r,
		reader:     gr,
	}, nil
}

// Read читает переданные сжатые данные.
func (gr *gzipReader) Read(p []byte) (int, error) {
	return gr.reader.Read(p)
}

// Close закрывает "считыватель" сжатых с помощью алгоритма gzip данных.
func (gr *gzipReader) Close() error {
	err := gr.ReadCloser.Close()
	if err != nil {
		return fmt.Errorf("error in close gzipReader: %w", err)
	}
	err = gr.reader.Close()
	if err != nil {
		err = fmt.Errorf("error in close gzipReader: %w", err)
	}
	return err
}
