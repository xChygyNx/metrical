package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatusMetricHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "Incorrect path",
			url:  "/update/other/value",
			want: want{
				code:        http.StatusNotFound,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Correct path",
			url:  "/update/gauge/someMetric/100.123",
			want: want{
				code: http.StatusOK,
				response: fmt.Sprintf(`{"status":"%s", "metric":"%s", "value":"%s"}`,
					http.StatusText(http.StatusOK), "someMetric", "100.123"),
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Absence metric Name",
			url:  "/update/gauge/100.123",
			want: want{
				code:        http.StatusNotFound,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Too much long path",
			url:  "/update/gauge/someMetric/100.123/needlessInformation",
			want: want{
				code:        http.StatusInternalServerError,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Incorrect path",
			url:  "/something/other/value",
			want: want{
				code:        http.StatusInternalServerError,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Correct path",
			url:  "/update/counter/someMetric/100",
			want: want{
				code:        http.StatusOK,
				response:    fmt.Sprintf(`{"status":"%s"}`, http.StatusText(http.StatusOK)),
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Absence metric Name",
			url:  "/update/counter/100; charset=utf-8",
			want: want{
				code:        http.StatusNotFound,
				response:    "",
				contentType: "text/plain",
			},
		},
		{
			name: "Too much long path",
			url:  "/update/counter/someMetric/100/needlessInformation",
			want: want{
				code:        http.StatusBadRequest,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.url, nil)
			w := httptest.NewRecorder()
			MetricHandle(w, request)

			result := w.Result()

			require.Equal(t, test.want.code, result.StatusCode)
			defer result.Body.Close()
			resBody, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			if string(resBody) != "" {
				assert.JSONEq(t, test.want.response, string(resBody))
			}
			assert.Equal(t, test.want.contentType, result.Header.Get("Content-Type"))
		})
	}
}
