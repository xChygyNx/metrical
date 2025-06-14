package agent

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

// Config структура для хранения параметров агента, считанных
// из командной строки.
type Config struct {
	HostPort       HostPort // хост и порт для отправки собранных метрик в формате "host:port".
	PollInterval   int      // интервал времени для сбора метрик в секундах.
	ReportInterval int      // интервал времени для отправки метрик на сервер в секундах.
}

// HostPort структура для хранения адреса сервера, куда будут отправляться собранные метрики.
type HostPort struct {
	Host string // хост сервера.
	Port int    // порт сервера.
}

// String представляет данные из структуры HostPort в текстовом формате "host:port".
func (hp *HostPort) String() string {
	return fmt.Sprintf("%s:%d", hp.Host, hp.Port)
}

// Set считывет данные из строки формата "host:port" в структуру HostPort.
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

func parseFlag() *Config {
	agentConfig := new(Config)
	defaultPollInterval := 2
	defaultReportInterval := 10
	pollInterval := flag.Int("p", defaultPollInterval, "Interval of collect metrics in seconds")
	reportInterval := flag.Int("r", defaultReportInterval, "Interval of send metrics on server in seconds")

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
