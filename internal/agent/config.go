package agent

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type config struct {
	Sha256Key      string
	HostAddr       HostPort
	PollInterval   int
	ReportInterval int
}

func GetConfig() (*config, error) {
	config := &config{}
	agentConfig := parseFlag()
	pollInterval, ok := os.LookupEnv("POLL_INTERVAL")
	if ok {
		res, err := strconv.Atoi(pollInterval)
		if err != nil {
			errorMsg := fmt.Sprintf("Incorrect value of environment variable POLL_INTERVAL: %v\n", err)
			return nil, errors.New(errorMsg)
		}
		config.PollInterval = res
	} else {
		config.PollInterval = agentConfig.PollInterval
	}

	reportInterval, ok := os.LookupEnv("POLL_INTERVAL")
	if ok {
		res, err := strconv.Atoi(reportInterval)
		if err != nil {
			errorMsg := fmt.Sprintf("Incorrect value of environment variable REPORT_INTERVAL: %v\n", err)
			return nil, errors.New(errorMsg)
		}
		config.ReportInterval = res
	} else {
		config.ReportInterval = agentConfig.ReportInterval
	}

	cryptoKey, ok := os.LookupEnv("KEY")
	if ok {
		config.Sha256Key = cryptoKey
	} else {
		config.Sha256Key = agentConfig.Sha256Key
	}

	hostAddr, ok := os.LookupEnv("ADDRESS")
	if ok {
		err := config.HostAddr.Set(hostAddr)
		if err != nil {
			return nil, err
		}
	} else {
		config.HostAddr = agentConfig.HostPort
	}

	return config, nil
}
