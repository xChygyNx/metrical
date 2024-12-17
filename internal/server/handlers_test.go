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
		resBody     types.Metrics
	}
	tests := []struct {
		name    string
		url     string
		value   float64
		reqBody types.Metrics
		want    want
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
		})
	}
}
