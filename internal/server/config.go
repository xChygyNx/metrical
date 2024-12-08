package server

import (
	"os"
)

type config struct {
	hostAddr HostPort
}

func getConfig() (*config, error) {
	config := new(config)
	serverConfig := parseFlag()

	hostAddr, ok := os.LookupEnv("ADDRESS")
	if ok {
		err := config.hostAddr.Set(hostAddr)
		if err != nil {
			return nil, err
		}
	} else {
		err := config.hostAddr.Set(serverConfig.String())
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}
