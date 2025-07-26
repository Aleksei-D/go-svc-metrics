package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"os"
	"strconv"
	"time"
)

func GetServerConfig() (*Config, error) {
	newConfig, err := InitConfig()
	if err != nil {
		return nil, err
	}

	serverFlagSet := flag.NewFlagSet("Server", flag.ExitOnError)
	serverAddr := serverFlagSet.String("a", defaultServerAddr, "input endpoint")
	logLevel := serverFlagSet.String("w", logLevelDefault, "log level")
	storeInterval := serverFlagSet.Int("i", StoreIntervalDefault, "store interval")
	fileStoragePath := serverFlagSet.String("f", FileStoragePathDefault, "file storage path")
	restore := serverFlagSet.Bool("r", restoreDefault, "log level")
	databaseDsn := serverFlagSet.String("d", "", "Database DSN")
	key := serverFlagSet.String("k", secretKeyDefault, "sha key")
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
	if newConfig.DatabaseDsn == nil {
		newConfig.DatabaseDsn = databaseDsn
	}
	if newConfig.Key == nil {
		newConfig.Key = key
	}
	return newConfig, nil
}

func GetAgentConfig() (*Config, error) {
	newConfig, err := InitConfig()
	if err != nil {
		return nil, err
	}

	agentFlagSet := flag.NewFlagSet("Agent", flag.ExitOnError)
	serverAddr := agentFlagSet.String("a", defaultServerAddr, "input endpoint")
	reportInterval := agentFlagSet.Int("r", reportInterval, "input reportInterval")
	pollInterval := agentFlagSet.Int("p", pollInterval, "input pollInterval")
	key := agentFlagSet.String("k", secretKeyDefault, "sha key")
	rateLimit := agentFlagSet.Uint("l", defaultRateLimit, "rate limit")
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
	if newConfig.Key == nil {
		newConfig.Key = key
	}
	if newConfig.RateLimit == nil {
		newConfig.RateLimit = rateLimit
	}
	return newConfig, nil
}

func InitConfig() (*Config, error) {
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
	DatabaseDsn     *string `env:"DATABASE_DSN"`
	Key             *string `env:"KEY"`
	RateLimit       *uint   `env:"RATE_LIMIT"`
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

func InitDefaultEnv() error {
	envDefaults := map[string]string{
		"ADDRESS":           defaultServerAddr,
		"LOG_LEVEL":         logLevelDefault,
		"STORE_INTERVAL":    strconv.Itoa(StoreIntervalDefault),
		"FILE_STORAGE_PATH": FileStoragePathDefault,
		"RESTORE":           "false",
	}
	for k, v := range envDefaults {
		err := os.Setenv(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
