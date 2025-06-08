package main

import (
	"go-svc-metrics/internal/config"
	"go-svc-metrics/internal/server"
	"go-svc-metrics/internal/storage"
	"net/http"
)

func main() {
	metricStorage := storage.InitMemStorage()
	configServe := config.GetServerConfig()
	r := server.GetMetricRouter(metricStorage)
	err := http.ListenAndServe(configServe.GetServeAddress(), r)
	if err != nil {
		panic(err)
	}
}
