package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/xChygyNx/metrical/internal/server/types"
)

func TestSetGaugeMetricHandler(t *testing.T) {
	type want struct {
		contentType string
		code        int
	}
	tests := []struct {
		reqBody types.Metrics
		name    string
		url     string
		want    want
		value   float64
	}{
		{
			name:  "Incorrect metric type",
			url:   "/update",
			value: 15.135,
			reqBody: types.Metrics{
				ID:    "someMetric",
				MType: "someType",
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "application/json",
			},
		},
		{
			name:  "Correct gauge data",
			url:   "/update",
			value: 15.135,
			reqBody: types.Metrics{
				ID:    "Alloc",
				MType: "gauge",
			},
			want: want{
				code:        http.StatusOK,
				contentType: "application/json",
			},
		},
	}
	storage := types.GetMemStorage()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.reqBody.Value = &test.value
			encodeData, _ := json.Marshal(test.reqBody)
			request := httptest.NewRequest(http.MethodPost, test.url, bytes.NewBuffer(encodeData))
			request.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			handler := SaveMetricHandle(storage)
			handler(w, request)
			result := w.Result()

			require.Equal(t, test.want.code, result.StatusCode)
			if result.StatusCode == http.StatusOK {
				defer func() {
					err := result.Body.Close()
					require.NoError(t, err)
				}()
				encodedBody, err := io.ReadAll(result.Body)
				require.NoError(t, err)

				var resultData types.Metrics
				err = json.Unmarshal(encodedBody, &resultData)
				require.NoError(t, err)
				assert.True(t, reflect.DeepEqual(resultData, test.reqBody))

				assert.Equal(t, test.want.contentType, result.Header.Get("Content-Type"))
			}
		})
	}
}

func TestStatusMetricHandler(t *testing.T) {
	type want struct {
		contentType string
		code        int
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
	storage := types.GetMemStorage()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.url, http.NoBody)
			for k, v := range test.pathValues {
				request.SetPathValue(k, v)
			}
			w := httptest.NewRecorder()
			handler := SaveMetricHandleOld(storage)
			handler(w, request)
			result := w.Result()

			require.Equal(t, test.want.code, result.StatusCode)
			defer func() {
				err := result.Body.Close()
				require.NoError(t, err)
			}()
			_, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			assert.Equal(t, test.want.contentType, result.Header.Get("Content-Type"))
		})
	}
}
