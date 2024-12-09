package server

import (
	"errors"
	"fmt"
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
			errorMsg := fmt.Sprintf("Addres must be like <host>:<port>, got %s\n", serverConfig.String())
			return nil, errors.New(errorMsg)
		}
	}

	return config, nil
}
