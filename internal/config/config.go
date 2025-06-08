package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
	"os"
	"time"
)

func GetServerConfig() *Config {
	var newConfig Config
	err := env.Parse(&newConfig)
	if err != nil {
		log.Fatal(err)
	}
	serverFlagSet := flag.NewFlagSet("Server", flag.ExitOnError)
	serverAddr := serverFlagSet.String("a", defaultServerAddr, "input endpoint")
	err = serverFlagSet.Parse(os.Args[1:])
	if err != nil {
		panic(err)
	}
	if newConfig.serverAddr == "" {
		newConfig.serverAddr = *serverAddr
	}
	return &newConfig
}

func GetAgentConfig() *Config {
	var newConfig Config
	err := env.Parse(&newConfig)
	if err != nil {
		log.Fatal(err)
	}
	agentFlagSet := flag.NewFlagSet("Agent", flag.ExitOnError)
	serverAddr := agentFlagSet.String("a", defaultServerAddr, "input endpoint")
	reportInterval := agentFlagSet.Int("r", reportInterval, "input reportInterval")
	pollInterval := agentFlagSet.Int("p", pollInterval, "input pollInterval")
	err = agentFlagSet.Parse(os.Args[1:])
	if err != nil {
		panic(err)
	}
	if newConfig.serverAddr == "" {
		newConfig.serverAddr = *serverAddr
	}
	if newConfig.reportInterval == 0 {
		newConfig.reportInterval = *reportInterval
	}
	if newConfig.pollInterval == 0 {
		newConfig.pollInterval = *pollInterval
	}
	return &newConfig
}

type Config struct {
	serverAddr     string `env:"ADDRESS"`
	reportInterval int    `env:"REPORT_INTERVAL"`
	pollInterval   int    `env:"POLL_INTERVAL"`
}

func (s Config) GetServeAddress() string {
	return s.serverAddr
}

func (s Config) GetPollInterval() time.Duration {
	return time.Duration(s.pollInterval) * time.Second
}

func (s Config) GetReportInterval() time.Duration {
	return time.Duration(s.reportInterval) * time.Second
}
