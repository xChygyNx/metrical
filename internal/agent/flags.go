package agent

import (
	"flag"
)

func parseFlag() *config {
	agentConfig := new(config)
	const defaultPollInterval = 2
	const defaultReportInterval = 10
	const defaultCryptoKey = ""
	pollInterval := flag.Int("p", defaultPollInterval, "Interval of collect metrics in seconds")
	reportInterval := flag.Int("r", defaultReportInterval, "Interval of send metrics on server in seconds")
	cryptoKey := flag.String("k", defaultCryptoKey, "Crypto key for encoding send data")

	hostPort := new(HostPort)
	flag.Var(hostPort, "a", "Net address host:port")

	flag.Parse()
	agentConfig.PollInterval = *pollInterval
	agentConfig.ReportInterval = *reportInterval
	agentConfig.Sha256Key = *cryptoKey

	if hostPort.Host == "" && hostPort.Port == 0 {
		hostPort.Host = "localhost"
		hostPort.Port = 8080
	}
	agentConfig.HostAddr = *hostPort
	return agentConfig
}
