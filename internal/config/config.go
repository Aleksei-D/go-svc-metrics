package config

import (
	"flag"
	"os"
	"time"
)

func GetServerConfig() *Config {
	serverFlagSet := flag.NewFlagSet("Server", flag.ExitOnError)
	serverAddr := serverFlagSet.String("a", defaultServerAddr, "input endpoint")
	err := serverFlagSet.Parse(os.Args[1:])
	if err != nil {
		panic(err)
	}
	return &Config{
		serverAddr: *serverAddr,
	}
}

func GetAgentConfig() *Config {
	agentFlagSet := flag.NewFlagSet("Agent", flag.ExitOnError)
	serverAddr := agentFlagSet.String("a", defaultServerAddr, "input endpoint")
	reportInterval := agentFlagSet.Int("r", reportInterval, "input reportInterval")
	pollInterval := agentFlagSet.Int64("p", pollInterval, "input pollInterval")
	err := agentFlagSet.Parse(os.Args[1:])
	if err != nil {
		panic(err)
	}
	return &Config{
		serverAddr:     *serverAddr,
		reportInterval: time.Duration(*reportInterval) * time.Second,
		pollInterval:   time.Duration(*pollInterval) * time.Second,
	}
}

type Config struct {
	serverAddr     string
	reportInterval time.Duration
	pollInterval   time.Duration
}

func (s Config) GetServeAddress() string {
	return s.serverAddr
}

func (s Config) GetPollInterval() time.Duration {
	return s.pollInterval
}

func (s Config) GetReportInterval() time.Duration {
	return s.reportInterval
}
