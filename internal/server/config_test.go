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
			config := parseFlag()
			assert.Equal(t, config.HostPort.Host, test.want.host)
			assert.Equal(t, config.HostPort.Port, test.want.port)
		})
	}
}
