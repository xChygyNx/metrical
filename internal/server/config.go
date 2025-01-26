package server

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/xChygyNx/metrical/internal/server/types"
)

const (
	sqlCreateGaugeTableCmd = "CREATE TABLE IF NOT EXISTS gauges (" +
		"	metric_name		varchar(100) PRIMARY KEY," +
		"	value			double precision" +
		");"
	sqlCreateCounterTableCmd = "CREATE TABLE IF NOT EXISTS counters (" +
		"	metric_name		varchar(100) PRIMARY KEY," +
		"	value			integer" +
		");"
)

type HostPort struct {
	Host string
	Port int
}

type Config struct {
	FileStoragePath string
	DBAddress       string
	Sha256Key       string
	HostPort        HostPort
	StoreInterval   int
	Restore         bool
}

func (hp *HostPort) String() string {
	return fmt.Sprintf("%s:%d", hp.Host, hp.Port)
}

func (conf *Config) String() string {
	return fmt.Sprintf(
		"StoreInterval: %d sec\n"+
			"FileStoragePath: %s\n"+
			"Restore: %t\n"+
			"Host: %s:%d\n"+
			"DBAddress:%s",
		conf.StoreInterval, conf.FileStoragePath, conf.Restore, conf.HostPort.Host, conf.HostPort.Port, conf.DBAddress)
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
	config := &Config{}
	defaultStoreInterval := 300
	defaultCryptoKey := ""

	flag.Var(&config.HostPort, "a", "Net address host:port")
	flag.IntVar(&config.StoreInterval, "i", defaultStoreInterval, "Time period for store metrics in the file")
	flag.StringVar(&config.FileStoragePath, "f", "", "File path for store metrics")
	flag.BoolVar(&config.Restore, "r", true, "Define should or not load store data from file before start")
	flag.StringVar(&config.DBAddress, "d", "", "Address of connecting to Data Base")
	flag.StringVar(&config.Sha256Key, "k", defaultCryptoKey, "Crypto key for encoding send data")
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

	dBAddress, ok := os.LookupEnv("DATABASE_DSN")
	if ok {
		config.DBAddress = dBAddress
	}

	return config, nil
}

func createMetricDB(connectInfo string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connectInfo)
	if err != nil {
		return nil, fmt.Errorf("error in create Metric DB: %w", err)
	}

	err = db.PingContext(context.Background())
	if err != nil {
		return nil, fmt.Errorf("DB is unreachable: %w", err)
	}

	ctx := context.Background()
	_, err = db.ExecContext(ctx, sqlCreateGaugeTableCmd)
	if err != nil {
		return nil, fmt.Errorf("error in create gauges table: %w", err)
	}
	_, err = db.ExecContext(ctx, sqlCreateCounterTableCmd)
	if err != nil {
		return nil, fmt.Errorf("error in create counters table: %w", err)
	}

	return db, nil
}

func GetHandlerConf(conf *Config) (*types.HandlerConf, error) {
	var db *sql.DB
	var err error
	if conf.DBAddress != "" {
		db, err = createMetricDB(conf.DBAddress)
		if err != nil {
			return nil, fmt.Errorf("error in create Metric Data Base: %w", err)
		}
	}
	return &types.HandlerConf{
		DB:                db,
		FileMetricStorage: conf.FileStoragePath,
		Sha256Key:         conf.Sha256Key,
		SyncFileRecord:    conf.StoreInterval == 0,
	}, nil
}
