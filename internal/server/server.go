package server

import (
	"github.com/go-chi/chi/v5"
	"go-svc-metrics/internal/handlers"
	"go-svc-metrics/internal/logger"
	"go-svc-metrics/internal/storage"
)

func GetMetricRouter(storage storage.Repositories) chi.Router {
	metricHandler := handlers.MetricHandler{Storage: storage}
	r := chi.NewRouter()
	r.Use(logger.LoggingMiddleware)
	r.Post(handlers.UpdateMetricHandlerPath, metricHandler.UpdateMetric)
	r.Get(handlers.GetMetricHandlerPath, metricHandler.GetMetricByName)
	r.Get("/", metricHandler.GetMetrics)
	return r
}
