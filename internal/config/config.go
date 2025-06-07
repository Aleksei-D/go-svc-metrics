package config

import (
	"fmt"
)

func GetServeConfig() *ServeConfig {
	return &ServeConfig{
		port: serverPort,
	}
}

type ServeConfig struct {
	port string
}

func (s ServeConfig) GetServeAddress() string {
	return fmt.Sprintf(serverAddressTemplate, s.port)
}
