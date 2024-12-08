package server

import (
	"os"
)

type config struct {
	HostAddr HostPort
}

func GetConfig() (*config, error) {
	config := new(config)
	serverConfig := parseFlag()

	hostAddr, ok := os.LookupEnv("ADDRESS")
	if ok {
		err := config.HostAddr.Set(hostAddr)
		if err != nil {
			return nil, err
		}
	} else {
		err := config.HostAddr.Set(serverConfig.String())
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}
