package server

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type HostPort struct {
	Host string
	Port int
}

func (hp *HostPort) String() string {
	return fmt.Sprintf("%s:%d", hp.Host, hp.Port)
}

func (hp *HostPort) Set(value string) error {
	hostPort := strings.Split(value, ":")
	numHostPortParts := 2
	if len(hostPort) != numHostPortParts {
		errorMsg := "must be value like <Host>:<Port>, got " + value
		return errors.New(errorMsg)
	}
	port, err := strconv.Atoi(hostPort[1])
	if err != nil {
		return fmt.Errorf("error in Atoi value of port: %w", err)
	}
	hp.Host = hostPort[0]
	hp.Port = port
	return nil
}

func parseFlag() *HostPort {
	hostPort := new(HostPort)
	flag.Var(hostPort, "a", "Net address host:port")
	flag.Parse()
	if hostPort.Host == "" && hostPort.Port == 0 {
		hostPort.Host = "localhost"
		hostPort.Port = 8080
	}
	return hostPort
}
