package handlers

import (
	"fmt"
	"go-svc-metrics/internal/config"
	"go-svc-metrics/internal/domain/local"
	"go-svc-metrics/internal/service"
	"net/http"
	"net/http/httptest"
)

// ExampleGetMetrics пример использования хендлера UpdateMetric.
// пример обращения к ендпоинту / .
func ExampleGetMetrics() {
	_ = config.InitDefaultEnv()
	configServe, _ := config.InitConfig()
	metricRepo, _ := local.NewMetricLocalRepository(configServe)

	metricService := service.NewMetricService(metricRepo)
	metricHandlers := NewCommonHandlers(metricService)

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	metricHandlers.GetMetrics(rr, req)
	fmt.Println(rr.Body.String())
}
