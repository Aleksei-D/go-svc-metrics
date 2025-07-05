package server

import (
	"github.com/go-chi/chi/v5"
	"go-svc-metrics/internal/handlers"
	"go-svc-metrics/internal/middleware"
	"go-svc-metrics/internal/storage"
)

func GetMetricRouter(storage storage.Repositories) chi.Router {
	metricHandler := handlers.MetricHandler{Storage: storage}
	r := chi.NewRouter()
	r.Use(middleware.CompressMiddleware, middleware.LoggingMiddleware)
	r.Get("/", metricHandler.GetMetrics)
	r.Get("/ping", metricHandler.GetPing)
	r.Route("/update", func(r chi.Router) {
		r.Post("/{metricType}/{metricName}/{metricValue}", metricHandler.UpdateMetric)
		r.Post("/", metricHandler.V2UpdateMetric)
	})
	r.Route("/value", func(r chi.Router) {
		r.Get("/{metricType}/{metricName}", metricHandler.GetMetricValue)
		r.Post("/", metricHandler.GetMetric)
	})
	r.Route("/updates", func(r chi.Router) {
		r.Post("/", metricHandler.UpdateBatchMetrics)
	})
	return r
}
