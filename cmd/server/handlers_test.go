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

func TestStatusGaugeHandler(t *testing.T) {
	type want struct {
		code		int
		response	string
		contentType	string
	}
	tests := []struct {
		name	string
		url			string
		want	want
	} {
		{
			name: "Correct path",
			url: "gauge/someMetric/100.123",
			want: want{
				code: http.StatusOK,
				response: fmt.Sprintf(`{"status":"%s"}`, http.StatusText(http.StatusOK)),
				contentType: "text/plain",
			},
		},
		{
			name: "Absence metric Name",
			url: "gauge/100.123",
			want: want{
				code: http.StatusNotFound,
				response: fmt.Sprintf(`{"status":"%s"}`, http.StatusText(http.StatusNotFound)),
				contentType: "text/plain",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T){
			request := httptest.NewRequest(http.MethodPost, "/update/" + test.url, nil)
			w := httptest.NewRecorder()
			gaugeHandle(w, request)

			result := w.Result()


			assert.Equal(t, result.StatusCode, test.want.code)
			defer result.Body.Close()
			resBody, err := io.ReadAll(result.Body)

			require.NoError(t, err)
			if string(resBody) != "" {
				assert.JSONEq(t, string(resBody), test.want.response)
			}
				assert.Equal(t, test.want.contentType, result.Header.Get("Content-Type"))
		})
	}
}