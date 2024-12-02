package main

import (
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
		contentType string
	}
	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "Incorrect gauge path",
			url:  "/update/other/value/12.3456",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Correct gauge path",
			url:  "/update/gauge/someMetric/100.123",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain",
			},
		},
		{
			name: "Correct counter path",
			url:  "/update/counter/someMetric/100",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain",
			},
		},
		{
			name: "Incorrect gauge value",
			url:  "/update/gauge/someMetric/none",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Incorrect counter value",
			url:  "/update/counter/someMetric/none",
			want: want{
				code:        http.StatusBadRequest,
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
			_, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			assert.Equal(t, test.want.contentType, result.Header.Get("Content-Type"))
		})
	}
}
