package types

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func emptyRequest() (*http.Request, error) {
	res, err := http.NewRequest(http.MethodGet, "localhost:8080", bytes.NewBuffer(make([]byte, 0)))
	if err != nil {
		err = fmt.Errorf("error in create request: %w", err)
		return nil, err
	}
	return res, nil
}

func notCorrectContentAcceptEncodingHeaderRequest() (*http.Request, error) {
	res, err := http.NewRequest(http.MethodGet, "localhost:8080", bytes.NewBuffer(make([]byte, 0)))
	if err != nil {
		err = fmt.Errorf("error in create request: %w", err)
		return nil, err
	}
	res.Header.Add("Accept-Encoding", "zip")
	res.Header.Add("Content-Encoding", "zip")
	return res, nil
}

func correctContentAcceptEncodingHeaderRequest() (*http.Request, error) {
	res, err := http.NewRequest(http.MethodGet, "localhost:8080", bytes.NewBuffer(make([]byte, 0)))
	if err != nil {
		err = fmt.Errorf("error in create request: %w", err)
		return nil, err
	}
	res.Header.Add("Accept-Encoding", "gzip")
	res.Header.Add("Content-Encoding", "gzip")
	return res, nil
}

func plainTextContentTypeRequest() (*http.Request, error) {
	res, err := http.NewRequest(http.MethodGet, "localhost:8080", bytes.NewBuffer(make([]byte, 0)))
	if err != nil {
		err = fmt.Errorf("error in create request: %w", err)
		return nil, err
	}
	res.Header.Add("Content-Type", "plain/text")
	return res, nil
}

func textHTMLContentTypeRequest() (*http.Request, error) {
	res, err := http.NewRequest(http.MethodGet, "localhost:8080", bytes.NewBuffer(make([]byte, 0)))
	if err != nil {
		err = fmt.Errorf("error in create request: %w", err)
		return nil, err
	}
	res.Header.Add("Content-Type", "text/html")
	return res, nil
}

func applicationJSONContentTypeRequest() (*http.Request, error) {
	res, err := http.NewRequest(http.MethodGet, "localhost:8080", bytes.NewBuffer(make([]byte, 0)))
	if err != nil {
		err = fmt.Errorf("error in create request: %w", err)
		return nil, err
	}
	res.Header.Add("Content-Type", "application/json")
	return res, nil
}

func TestIsAcceptEncoding(t *testing.T) {
	tests := []struct {
		name            string
		requestCallable func() (*http.Request, error)
		want            bool
	}{
		{
			name:            "Request without Headers",
			requestCallable: emptyRequest,
			want:            false,
		},
		{
			name:            "Request with incorrect Header Accept-Encoding",
			requestCallable: notCorrectContentAcceptEncodingHeaderRequest,
			want:            false,
		},
		{
			name:            "Request with correct Header Accept-Encoding",
			requestCallable: correctContentAcceptEncodingHeaderRequest,
			want:            true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request, err := test.requestCallable()
			assert.Nil(t, err)
			assert.Equal(t, IsAcceptEncoding(request.Header), test.want)
		})
	}
}

func TestIsContentEncoding(t *testing.T) {
	tests := []struct {
		name            string
		requestCallable func() (*http.Request, error)
		want            bool
	}{
		{
			name:            "Request without Headers",
			requestCallable: emptyRequest,
			want:            false,
		},
		{
			name:            "Request with incorrect Header Content-Encoding",
			requestCallable: notCorrectContentAcceptEncodingHeaderRequest,
			want:            false,
		},
		{
			name:            "Request with correct Header Content-Encoding",
			requestCallable: correctContentAcceptEncodingHeaderRequest,
			want:            true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request, err := test.requestCallable()
			assert.Nil(t, err)
			assert.Equal(t, IsContentEncoding(request.Header), test.want)
		})
	}
}

func TestIsCompressData(t *testing.T) {
	tests := []struct {
		name            string
		requestCallable func() (*http.Request, error)
		want            bool
	}{
		{
			name:            "Request without Headers",
			requestCallable: emptyRequest,
			want:            false,
		},
		{
			name:            "Request with Header Content-Type plain/text",
			requestCallable: plainTextContentTypeRequest,
			want:            false,
		},
		{
			name:            "Request with Header Content-Type text/html",
			requestCallable: textHTMLContentTypeRequest,
			want:            true,
		},
		{
			name:            "Request with Header Content-Type application/json",
			requestCallable: applicationJSONContentTypeRequest,
			want:            true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request, err := test.requestCallable()
			assert.Nil(t, err)
			assert.Equal(t, IsCompressData(request.Header), test.want)
		})
	}
}
