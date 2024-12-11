package server

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
		name       string
		url        string
		pathValues map[string]string
		want       want
	}{
		{
			name: "Incorrect gauge path",
			url:  "/update/other/metric/12.3456",
			pathValues: map[string]string{
				"mType":  "other",
				"metric": "metric",
				"value":  "12.3456",
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Correct gauge path",
			url:  "/update/gauge/Mallocs/100.123",
			pathValues: map[string]string{
				"mType":  "gauge",
				"metric": "Mallocs",
				"value":  "100.123",
			},
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain",
			},
		},
		{
			name: "Inorrect gauge metric value",
			url:  "/update/gauge/Mallocs/none",
			pathValues: map[string]string{
				"mType":  "gauge",
				"metric": "Mallocs",
				"value":  "none",
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Correct counter path",
			url:  "/update/counter/PollCount/100",
			pathValues: map[string]string{
				"mType":  "counter",
				"metric": "PollCount",
				"value":  "100",
			},
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain",
			},
		},
		{
			name: "Incorrect counter value",
			url:  "/update/counter/PollCount/none",
			pathValues: map[string]string{
				"mType":  "counter",
				"metric": "PollCount",
				"value":  "none",
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	storage := GetMemStorage()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.url, nil)
			for k, v := range test.pathValues {
				request.SetPathValue(k, v)
			}
			w := httptest.NewRecorder()
			handler := SaveMetricHandle(storage)
			handler(w, request)
			result := w.Result()

			require.Equal(t, test.want.code, result.StatusCode)
			defer result.Body.Close()
			_, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			assert.Equal(t, test.want.contentType, result.Header.Get("Content-Type"))
		})
	}
}
