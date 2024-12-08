package agent

import (
	"fmt"
	"os"
	"strconv"
)

type config struct {
	pollInterval   int
	reportInterval int
	hostAddr       HostPort
}

func getConfig() (*config, error) {
	config := new(config)
	agentConfig := parseFlag()
	pollInterval, ok := os.LookupEnv("POLL_INTERVAL")
	if ok {
		res, err := strconv.Atoi(pollInterval)
		if err != nil {
			return nil, fmt.Errorf("incorrect value of environment variable POLL_INTERVAL")
		}
		config.pollInterval = res
	} else {
		config.pollInterval = agentConfig.PollInterval
	}

	reportInterval, ok := os.LookupEnv("POLL_INTERVAL")
	if ok {
		res, err := strconv.Atoi(reportInterval)
		if err != nil {
			return nil, fmt.Errorf("incorrect value of environment variable REPORT_INTERVAL")
		}
		config.reportInterval = res
	} else {
		config.reportInterval = agentConfig.ReportInterval
	}

	hostAddr, ok := os.LookupEnv("ADDRESS")
	if ok {
		err := config.hostAddr.Set(hostAddr)
		if err != nil {
			return nil, err
		}
	} else {
		config.hostAddr = agentConfig.HostPort
	}

	return config, nil
}
