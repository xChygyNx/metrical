package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFlags(t *testing.T) {
	type want struct {
		host string
		port int
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "Default values",
			want: want{
				port: 8080,
				host: "localhost",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hostPort := parseFlag()
			assert.Equal(t, hostPort.Host, test.want.host)
			assert.Equal(t, hostPort.Port, test.want.port)
		})
	}
}
