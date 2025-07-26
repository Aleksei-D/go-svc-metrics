package main

import (
	"go-svc-metrics/internal/agent"
	"go-svc-metrics/internal/logger"
)

func main() {
	metricAgent, err := agent.NewMetricUpdater()
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	err = metricAgent.Run()
	if err != nil {
		logger.Log.Fatal(err.Error())
	}
}
