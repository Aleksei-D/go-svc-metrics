package main

import (
	"go-svc-metrics/internal/config"
	"go-svc-metrics/internal/handlers"
	"go-svc-metrics/internal/storage"
	"net/http"
)

func main() {
	serveConfig := config.GetServeConfig()
	memStorage := storage.InitMemStorage()
	metricHandler := handlers.MetricHandler{Storage: memStorage}
	mux := http.NewServeMux()
	mux.HandleFunc(handlers.MetricHandlerPath, metricHandler.Serve)

	err := http.ListenAndServe(serveConfig.GetServeAddress(), mux)
	if err != nil {
		panic(err)
	}
}
