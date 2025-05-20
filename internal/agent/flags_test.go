package agent

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFlagWithoutArgs(t *testing.T) {
	os.Args = append(os.Args, "-p=10")
	os.Args = append(os.Args, "-r=30")
	os.Args = append(os.Args, "-a=google.com:1234")
	t.Run("Read command line args with args", func(t *testing.T) {
		conf := parseFlag()
		assert.Equal(t, conf.PollInterval, 10)
		assert.Equal(t, conf.ReportInterval, 30)
		assert.Equal(t, conf.HostPort.Port, 1234)
		assert.Equal(t, conf.HostPort.Host, "google.com")
	})
}
