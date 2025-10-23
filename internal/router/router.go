// Модуль router реализует роутеры.
package router

import (
	"crypto/rsa"
	"go-svc-metrics/internal/config"
	"go-svc-metrics/internal/handlers"
	middleware2 "go-svc-metrics/internal/middleware"
	"go-svc-metrics/internal/service"

	"net/http/pprof"

	"github.com/go-chi/chi/v5"
)

// NewRouter возвращает роутеры для сервера.
func NewRouter(metricService *service.MetricService, config *config.Config, privateKey *rsa.PrivateKey) chi.Router {
	updateHandlers := handlers.NewUpdateHandlers(metricService)
	valueHandlers := handlers.NewValueHandlers(metricService)
	commonHandlers := handlers.NewCommonHandlers(metricService)

	r := chi.NewRouter()
	r.Get("/", commonHandlers.GetMetrics)
	r.Get("/ping", commonHandlers.GetPing)
	r.Route("/update", func(r chi.Router) {
		cryptoMiddleware := middleware2.CryptoRSAMiddleware{PrivateKey: privateKey}
		r.Use(cryptoMiddleware.GetCryptoRSAMiddleware)
		r.Use(middleware2.CompressMiddleware)
		r.Use(middleware2.LoggingMiddleware)
		r.Post("/{metricType}/{metricName}/{metricValue}", updateHandlers.UpdateMetric)
		r.Post("/", updateHandlers.V2UpdateMetric)
	})
	r.Route("/value", func(r chi.Router) {
		r.Post("/", valueHandlers.GetMetric)
	})
	r.Route("/updates", func(r chi.Router) {
		cryptoMiddleware := middleware2.CryptoRSAMiddleware{PrivateKey: privateKey}
		r.Use(cryptoMiddleware.GetCryptoRSAMiddleware)
		r.Use(middleware2.CompressMiddleware)
		r.Use(middleware2.LoggingMiddleware)
		r.Post("/", updateHandlers.UpdateBatchMetrics)
	})

	r.Route("/debug/pprof", func(r chi.Router) {
		r.Get("/", pprof.Index)
		r.Get("/cmdline", pprof.Cmdline)
		r.Get("/profile", pprof.Profile)
		r.Get("/symbol", pprof.Symbol)
		r.Get("/trace", pprof.Trace)
		r.Get("/{cmd}", pprof.Index)
	})
	return r
}
