package agent

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

// GetConfig возвращает структуру config в которой заданы такие параметры
// работы агента, как адрес хоста и порт куда будут отправляться собранные метрики,
// интервалы времени для сбора и отправки метрик.
//
// Приоритет источников для задания параметров агента:
// 1) Переменные окружения
// 2) Аргументы командной строки.
func GetConfig() (*Config, error) {
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

	reportInterval, ok := os.LookupEnv("POLL_INTERVAL")
	if ok {
		res, err := strconv.Atoi(reportInterval)
		if err != nil {
			errorMsg := fmt.Sprintf("Incorrect value of environment variable REPORT_INTERVAL: %v\n", err)
			return nil, errors.New(errorMsg)
		}
		config.ReportInterval = res
	}

	hostAddr, ok := os.LookupEnv("ADDRESS")
	if ok {
		err := config.HostPort.Set(hostAddr)
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}
