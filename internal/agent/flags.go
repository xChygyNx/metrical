package agent

import (
	"errors"
	"flag"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

type AgentConfig struct {
	Sha256Key      string
	HostPort       HostPort
	PollInterval   int
	ReportInterval int
	RateLimit      int8
}

type HostPort struct {
	Host string
	Port int
}

func (hp *HostPort) String() string {
	return fmt.Sprintf("%s:%d", hp.Host, hp.Port)
}

func (hp *HostPort) Set(value string) error {
	hostPort := strings.Split(value, ":")
	numPartsHostPort := 2
	if len(hostPort) != numPartsHostPort {
		errorMsg := "must be value like <Host>:<Port>, got " + value
		return errors.New(errorMsg)
	}
	port, err := strconv.Atoi(hostPort[1])
	if err != nil {
		return fmt.Errorf("error in Atoi port value: %w", err)
	}
	hp.Host = hostPort[0]
	hp.Port = port
	return nil
}

func parseFlag() *AgentConfig {
	agentConfig := new(AgentConfig)
	const defaultPollInterval = 2
	const defaultReportInterval = 10
	const defaultCryptoKey = ""
	pollInterval := flag.Int("p", defaultPollInterval, "Interval of collect metrics in seconds")
	reportInterval := flag.Int("r", defaultReportInterval, "Interval of send metrics on server in seconds")
	cryptoKey := flag.String("k", defaultCryptoKey, "Crypto key for encoding send data")
	rateLimit := flag.Int("l", runtime.NumCPU(), "Count of workers for collect metrics")

	hostPort := new(HostPort)
	flag.Var(hostPort, "a", "Net address host:port")

	flag.Parse()
	agentConfig.PollInterval = *pollInterval
	agentConfig.ReportInterval = *reportInterval
	agentConfig.Sha256Key = *cryptoKey
	agentConfig.RateLimit = int8(*rateLimit)

	if hostPort.Host == "" && hostPort.Port == 0 {
		hostPort.Host = "localhost"
		hostPort.Port = 8080
	}
	agentConfig.HostPort = *hostPort
	return agentConfig
}
