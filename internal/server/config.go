package server

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/xChygyNx/metrical/internal/server/types"
)

type HostPort struct {
	Host string
	Port int
}

type Config struct {
	FileStoragePath string
	HostPort        HostPort
	StoreInterval   int
	Restore         bool
}

var StorageFile = "metrics.json"

func (hp *HostPort) String() string {
	return fmt.Sprintf("%s:%d", hp.Host, hp.Port)
}

func (conf *Config) String() string {
	return fmt.Sprintf(
		"StoreInterval: %d sec\n"+
			"FileStoragePath: %s\n"+
			"Restore: %t\n"+
			"Host: %s:%d",
		conf.StoreInterval, conf.FileStoragePath, conf.Restore, conf.HostPort.Host, conf.HostPort.Port)
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

func parseFlag() *Config {
	config := new(Config)
	hostPort := new(HostPort)
	flag.Var(hostPort, "a", "Net address host:port")
	defaultStoreInterval := 300
	flag.IntVar(&config.StoreInterval, "i", defaultStoreInterval, "Time period for store metrics in the file")
	flag.StringVar(&config.FileStoragePath, "f", StorageFile, "File path for store metrics")
	flag.BoolVar(&config.Restore, "r", true, "Define should or not load store data from file before start")
	flag.Parse()
	if config.HostPort.Host == "" && config.HostPort.Port == 0 {
		config.HostPort.Host = "localhost"
		config.HostPort.Port = 8080
	}
	return config
}

func GetConfig() (*Config, error) {
	config := parseFlag()

	hostAddr, ok := os.LookupEnv("ADDRESS")
	if ok {
		err := config.HostPort.Set(hostAddr)
		if err != nil {
			return nil, err
		}
	}

	storeInterval, ok := os.LookupEnv("STORE_INTERVAL")
	if ok {
		interval, err := strconv.Atoi(storeInterval)
		if err != nil {
			return nil, fmt.Errorf(
				"environment variable STORE_INTERVAL must be numerical, got %s: %w", storeInterval, err)
		}
		config.StoreInterval = interval
	}

	filePath, ok := os.LookupEnv("FILE_STORAGE_PATH")
	if ok {
		config.FileStoragePath = filePath
	}

	restore, ok := os.LookupEnv("RESTORE")
	if ok {
		restoreBool, err := strconv.ParseBool(restore)
		if err != nil {
			return nil, fmt.Errorf(
				"environment variable RESTORE must be bool, got %s: %w", restore, err)
		}
		config.Restore = restoreBool
	}

	return config, nil
}

func GetSyncInfo(conf Config) types.SyncInfo {
	return types.SyncInfo{
		FileMetricStorage: conf.FileStoragePath,
		SyncFileRecord:    conf.StoreInterval == 0,
	}
}
