package agent

import (
	"fmt"
	"os"
	"strconv"
)

type config struct {
	PollInterval   int
	ReportInterval int
	HostAddr       HostPort
}

func GetConfig() (*config, error) {
	config := new(config)
	agentConfig := parseFlag()
	pollInterval, ok := os.LookupEnv("POLL_INTERVAL")
	if ok {
		res, err := strconv.Atoi(pollInterval)
		if err != nil {
			return nil, fmt.Errorf("incorrect value of environment variable POLL_INTERVAL")
		}
		config.PollInterval = res
	} else {
		config.PollInterval = agentConfig.PollInterval
	}

	reportInterval, ok := os.LookupEnv("POLL_INTERVAL")
	if ok {
		res, err := strconv.Atoi(reportInterval)
		if err != nil {
			return nil, fmt.Errorf("incorrect value of environment variable REPORT_INTERVAL")
		}
		config.ReportInterval = res
	} else {
		config.ReportInterval = agentConfig.ReportInterval
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
