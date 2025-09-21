package router

import (
	"go-svc-metrics/internal/config"
	"go-svc-metrics/internal/handlers"
	middleware2 "go-svc-metrics/internal/middleware"
	"go-svc-metrics/internal/service"

	"github.com/go-chi/chi/v5"
)

func NewRouter(metricService *service.MetricService, config *config.Config) chi.Router {
	metricHandler := handlers.NewMetricHandler(metricService)

	r := chi.NewRouter()
	cryptoMiddleware := middleware2.CryptoMiddleware{Config: config}
	r.Use(cryptoMiddleware.GetCryptoMiddleware)
	r.Use(middleware2.CompressMiddleware)
	r.Use(middleware2.LoggingMiddleware)

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
