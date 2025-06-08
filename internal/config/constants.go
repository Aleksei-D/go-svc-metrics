package config

import "time"

const (
	defaultServerAddr = "localhost:8080"
	pollInterval      = 2 * time.Second
	reportInterval    = 10 * time.Second
)
