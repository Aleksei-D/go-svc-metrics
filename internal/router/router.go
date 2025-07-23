package router

import (
	"github.com/go-chi/chi/v5"
	"go-svc-metrics/internal/config"
	"go-svc-metrics/internal/handlers"
	"go-svc-metrics/internal/middleware"
	"go-svc-metrics/internal/storage"
	"go-svc-metrics/internal/usecase"
)

func GetMetricRouter(metricRepository storage.MetricRepository, config *config.Config) chi.Router {
	metricUseCase := usecase.NewMetricUseCase(metricRepository)
	metricHandler := handlers.MetricHandler{MetricUseCase: metricUseCase}

	r := chi.NewRouter()
	cryptoMiddleware := middleware.CryptoMiddleware{Config: config}
	r.Use(cryptoMiddleware.GetCryptoMiddleware)
	r.Use(middleware.CompressMiddleware)
	r.Use(middleware.LoggingMiddleware)

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
