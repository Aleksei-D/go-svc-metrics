package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"os"
	"time"
)

func GetServerConfig() (*Config, error) {
	newConfig, err := initConfig()
	if err != nil {
		return nil, err
	}

	serverFlagSet := flag.NewFlagSet("Server", flag.ExitOnError)
	serverAddr := serverFlagSet.String("a", defaultServerAddr, "input endpoint")
	logLevel := serverFlagSet.String("l", logLevelDefault, "log level")
	storeInterval := serverFlagSet.Int("i", StoreIntervalDefault, "store interval")
	fileStoragePath := serverFlagSet.String("f", FileStoragePathDefault, "file storage path")
	restore := serverFlagSet.Bool("r", restoreDefault, "log level")
	err = serverFlagSet.Parse(os.Args[1:])
	if err != nil {
		return nil, err
	}
	if newConfig.ServerAddr == nil {
		newConfig.ServerAddr = serverAddr
	}
	if newConfig.LogLevel == nil {
		newConfig.LogLevel = logLevel
	}
	if newConfig.StoreInterval == nil {
		newConfig.StoreInterval = storeInterval
	}
	if newConfig.FileStoragePath == nil {
		newConfig.FileStoragePath = fileStoragePath
	}
	if newConfig.Restore == nil {
		newConfig.Restore = restore
	}
	return newConfig, nil
}

func GetAgentConfig() (*Config, error) {
	newConfig, err := initConfig()
	if err != nil {
		return nil, err
	}

	agentFlagSet := flag.NewFlagSet("Agent", flag.ExitOnError)
	serverAddr := agentFlagSet.String("a", defaultServerAddr, "input endpoint")
	reportInterval := agentFlagSet.Int("r", reportInterval, "input reportInterval")
	pollInterval := agentFlagSet.Int("p", pollInterval, "input pollInterval")
	err = agentFlagSet.Parse(os.Args[1:])
	if err != nil {
		panic(err)
	}
	if newConfig.ServerAddr == nil {
		newConfig.ServerAddr = serverAddr
	}
	if newConfig.ReportInterval == nil {
		newConfig.ReportInterval = reportInterval
	}
	if newConfig.PollInterval == nil {
		newConfig.PollInterval = pollInterval
	}
	return newConfig, nil
}

func initConfig() (*Config, error) {
	var newConfig Config
	err := env.Parse(&newConfig)
	if err != nil {
		return nil, err
	}
	return &newConfig, err
}

type Config struct {
	ServerAddr      *string `env:"ADDRESS"`
	ReportInterval  *int    `env:"REPORT_INTERVAL"`
	PollInterval    *int    `env:"POLL_INTERVAL"`
	LogLevel        *string `env:"LOG_LEVEL"`
	StoreInterval   *int    `env:"STORE_INTERVAL"`
	FileStoragePath *string `env:"FILE_STORAGE_PATH"`
	Restore         *bool   `env:"RESTORE"`
}

func (s Config) GetServeAddress() string {
	return *s.ServerAddr
}

func (s Config) GetPollInterval() time.Duration {
	return time.Duration(*s.PollInterval) * time.Second
}

func (s Config) GetReportInterval() time.Duration {
	return time.Duration(*s.ReportInterval) * time.Second
}

func (s Config) GetStoreInterval() time.Duration {
	return time.Duration(*s.StoreInterval) * time.Second
}
