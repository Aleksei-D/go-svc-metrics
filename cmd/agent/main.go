package main

import (
	"go-svc-metrics/internal/agent"
	"go-svc-metrics/internal/logger"
	"go-svc-metrics/internal/utils/helpers"
)

var buildVersion, buildDate, buildCommit string

func main() {
	helpers.PrintBuildVersion(buildVersion, buildDate, buildCommit)
	metricAgent, err := agent.NewMetricUpdater()
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	err = metricAgent.Run()
	if err != nil {
		logger.Log.Fatal(err.Error())
	}
}
