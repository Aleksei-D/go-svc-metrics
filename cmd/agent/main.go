package main

import "go-svc-metrics/internal/agent"

func main() {
	app := agent.GetNewMetricUpdater()
	app.MetricProcessing()
}
