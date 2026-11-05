package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func getArgsWithoutDashes() map[string]bool {
	args := make(map[string]bool, 0)
	for _, elem := range os.Args {
		if strings.HasPrefix(elem, "-") {
			elem = strings.TrimLeft(elem, "-")
			args[elem] = true
		}
	}
	return args
}

func parseConfigFromJSON(configFile string, config *Config) (*Config, error) {
	fileData, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("error in read file %s: %w", configFile, err)
	}
	tmpConfig := &TmpConfig{}
	err = json.Unmarshal(fileData, tmpConfig)
	if err != nil {
		return nil, fmt.Errorf("error in unmarshal json data: %w", err)
	}
	args := getArgsWithoutDashes()

	if _, ok := args["a"]; !ok {
		err = config.HostPort.Set(tmpConfig.HostPort)
		if err != nil {
			return nil, fmt.Errorf("error in parse HostPort from JSON: %w", err)
		}
	}
	if _, ok := args["p"]; !ok {
		config.PollInterval = tmpConfig.PollInterval
	}
	if _, ok := args["r"]; !ok {
		config.ReportInterval = tmpConfig.ReportInterval
	}
	if _, ok := args["g"]; !ok {
		config.GRPCPort = tmpConfig.GRPCPort
	}
	if _, ok := args["crypto-key"]; !ok {
		config.RSAPublicKey = tmpConfig.RSAPublicKey
	}

	return config, nil
}

// GetConfig возвращает структуру config в которой заданы такие параметры
// работы агента, как адрес хоста и порт куда будут отправляться собранные метрики,
// интервалы времени для сбора и отправки метрик.
//
// Приоритет источников для задания параметров агента:
// 1) Переменные окружения
// 2) Аргументы командной строки.
func GetConfig() (config *Config, err error) {
	config = parseFlag()

	configFileEnv, ok := os.LookupEnv("CONFIG")
	configFileArg := config.ConfigFile
	if ok {
		config, err = parseConfigFromJSON(configFileEnv, config)
		if err != nil {
			return nil, fmt.Errorf("error parse config from JSON from env: %w", err)
		}
	} else if configFileArg != "" {
		config, err = parseConfigFromJSON(configFileArg, config)
		if err != nil {
			return nil, fmt.Errorf("error parse config from JSON from args: %w", err)
		}
	}

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

	gRPCPort, ok := os.LookupEnv("GRPC_PORT")
	if ok {
		config.GRPCPort = gRPCPort
	}

	hostAddr, ok := os.LookupEnv("ADDRESS")
	if ok {
		err := config.HostPort.Set(hostAddr)
		if err != nil {
			return nil, err
		}
	}

	publicKey, ok := os.LookupEnv("CRYPTO_KEY")
	if ok {
		config.RSAPublicKey = publicKey
	}

	return config, nil
}
