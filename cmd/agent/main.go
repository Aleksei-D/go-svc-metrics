package main

import (
	"go-svc-metrics/internal/agent"
	"go-svc-metrics/internal/logger"
)

func main() {
	app, err := agent.GetNewMetricUpdater()
	if err != nil {
		logger.Log.Fatal(err.Error())
	}
	app.MetricProcessing()
}
