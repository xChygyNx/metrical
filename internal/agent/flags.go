package agent

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type AgentConfig struct {
	PollInterval   int
	ReportInterval int
	HostPort       HostPort
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
	if len(hostPort) != 2 {
		error_msg := "must be value like <Host>:<Port>, got " + value
		return errors.New(error_msg)
	}
	port, err := strconv.Atoi(hostPort[1])
	if err != nil {
		return err
	}
	hp.Host = hostPort[0]
	hp.Port = port
	return nil
}

func parseFlag() *AgentConfig {
	agentConfig := new(AgentConfig)
	pollInterval := flag.Int("p", 2, "Interval of collect metrics in seconds")
	reportInterval := flag.Int("r", 10, "Interval of send metrics on server in seconds")

	hostPort := new(HostPort)
	flag.Var(hostPort, "a", "Net address host:port")

	flag.Parse()
	agentConfig.PollInterval = *pollInterval
	agentConfig.ReportInterval = *reportInterval

	if hostPort.Host == "" && hostPort.Port == 0 {
		hostPort.Host = "localhost"
		hostPort.Port = 8080
	}
	agentConfig.HostPort = *hostPort
	return agentConfig
}
