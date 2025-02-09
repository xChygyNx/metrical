package agent

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type config struct {
	Sha256Key      string
	HostAddr       HostPort
	PollInterval   int
	ReportInterval int
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

func GetConfig() (*config, error) {
	config := parseFlag()
	pollInterval, ok := os.LookupEnv("POLL_INTERVAL")
	if ok {
		res, err := strconv.Atoi(pollInterval)
		if err != nil {
			errorMsg := fmt.Sprintf("Incorrect value of environment variable POLL_INTERVAL: %v\n", err)
			return nil, errors.New(errorMsg)
		}
		config.PollInterval = res
	}

	reportInterval, ok := os.LookupEnv("REPORT_INTERVAL")
	if ok {
		res, err := strconv.Atoi(reportInterval)
		if err != nil {
			errorMsg := fmt.Sprintf("Incorrect value of environment variable REPORT_INTERVAL: %v\n", err)
			return nil, errors.New(errorMsg)
		}
		config.ReportInterval = res
	}

	cryptoKey, ok := os.LookupEnv("KEY")
	if ok {
		config.Sha256Key = cryptoKey
	}

	hostAddr, ok := os.LookupEnv("ADDRESS")
	if ok {
		err := config.HostAddr.Set(hostAddr)
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}
