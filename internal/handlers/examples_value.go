package handlers

import (
	"fmt"
	"go-svc-metrics/internal/config"
	"go-svc-metrics/internal/domain/local"
	"go-svc-metrics/internal/service"
	"go-svc-metrics/internal/utils/helpers"
	"go-svc-metrics/models"
	"net/http"
	"net/http/httptest"
)

// ExampleGetMetricValue пример использования хендлера UpdateMetric.
// пример обращения к ендпоинту /value/{metricType}/{metricName}.
func ExampleGetMetricValue() {
	_ = config.InitDefaultEnv()
	configServe, _ := config.InitConfig()
	metricRepo, _ := local.NewMetricLocalRepository(configServe)

	metricService := service.NewMetricService(metricRepo)
	valueHandlers := NewValueHandlers(metricService)

	path := fmt.Sprintf("/value/%s/%s", models.Counter, helpers.ValidCounterID)
	req, _ := http.NewRequest(http.MethodPost, path, nil)
	rr := httptest.NewRecorder()
	valueHandlers.GetMetricValue(rr, req)
	fmt.Println(rr.Body.String())
}
