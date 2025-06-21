package main

import (
	"go-svc-metrics/internal/config"
	"go-svc-metrics/internal/logger"
	"go-svc-metrics/internal/server"
	"go-svc-metrics/internal/storage"
	"net/http"
)

func main() {
	configServe := config.GetServerConfig()
	err := logger.Initialize(configServe.LogLevel)
	if err != nil {
		panic("cannot initialize zap")
	}

	metricStorage := storage.InitMemStorage()
	r := server.GetMetricRouter(metricStorage)
	err = http.ListenAndServe(configServe.GetServeAddress(), r)
	if err != nil {
		panic(err)
	}
}
