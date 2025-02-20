package agent

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

type config struct {
	Sha256Key      string
	HostAddr       HostPort
	PollInterval   int
	ReportInterval int
	RateLimit      int
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

func parseFlag() *config {
	agentConfig := new(config)
	const defaultPollInterval = 2
	const defaultReportInterval = 10
	const defaultCryptoKey = ""
	defaultRateLimit := runtime.NumCPU()
	pollInterval := flag.Int("p", defaultPollInterval, "Interval of collect metrics in seconds")
	reportInterval := flag.Int("r", defaultReportInterval, "Interval of send metrics on server in seconds")
	cryptoKey := flag.String("k", defaultCryptoKey, "Crypto key for encoding send data")
	rateLimit := flag.Int("l", defaultRateLimit, "Number of agent threads")

	hostPort := new(HostPort)
	flag.Var(hostPort, "a", "Net address host:port")

	flag.Parse()
	agentConfig.PollInterval = *pollInterval
	agentConfig.ReportInterval = *reportInterval
	agentConfig.Sha256Key = *cryptoKey
	agentConfig.RateLimit = *rateLimit

	if hostPort.Host == "" && hostPort.Port == 0 {
		hostPort.Host = "localhost"
		hostPort.Port = 8080
	}
	agentConfig.HostAddr = *hostPort
	return agentConfig
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

	rateLimit, ok := os.LookupEnv("RATE_LIMIT")
	if ok {
		res, err := strconv.Atoi(rateLimit)
		if err != nil {
			errorMsg := fmt.Sprintf("Incorrect value of environment variable RATE_LIMIT: %v\n", err)
			return nil, errors.New(errorMsg)
		}
		config.RateLimit = res
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
