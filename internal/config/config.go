package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
	"os"
	"time"
)

func GetServerConfig() *Config {
	newConfig := initConfig()
	serverFlagSet := flag.NewFlagSet("Server", flag.ExitOnError)
	serverAddr := serverFlagSet.String("a", defaultServerAddr, "input endpoint")
	err := serverFlagSet.Parse(os.Args[1:])
	if err != nil {
		panic(err)
	}
	if newConfig.ServerAddr == "" {
		newConfig.ServerAddr = *serverAddr
	}
	return newConfig
}

func GetAgentConfig() *Config {
	newConfig := initConfig()
	agentFlagSet := flag.NewFlagSet("Agent", flag.ExitOnError)
	serverAddr := agentFlagSet.String("a", defaultServerAddr, "input endpoint")
	reportInterval := agentFlagSet.Int("r", reportInterval, "input reportInterval")
	pollInterval := agentFlagSet.Int("p", pollInterval, "input pollInterval")
	err := agentFlagSet.Parse(os.Args[1:])
	if err != nil {
		panic(err)
	}
	if newConfig.ServerAddr == "" {
		newConfig.ServerAddr = *serverAddr
	}
	if newConfig.ReportInterval == 0 {
		newConfig.ReportInterval = *reportInterval
	}
	if newConfig.PollInterval == 0 {
		newConfig.PollInterval = *pollInterval
	}
	return newConfig
}

func initConfig() *Config {
	var newConfig Config
	err := env.Parse(&newConfig)
	if err != nil {
		log.Fatal(err)
	}
	return &newConfig
}

type Config struct {
	ServerAddr     string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

func (s Config) GetServeAddress() string {
	return s.ServerAddr
}

func (s Config) GetPollInterval() time.Duration {
	return time.Duration(s.PollInterval) * time.Second
}

func (s Config) GetReportInterval() time.Duration {
	return time.Duration(s.ReportInterval) * time.Second
}
