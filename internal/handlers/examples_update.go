package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-svc-metrics/internal/config"
	"go-svc-metrics/internal/domain/local"
	"go-svc-metrics/internal/service"
	"go-svc-metrics/internal/utils/helpers"
	"go-svc-metrics/models"
	"net/http"
	"net/http/httptest"
)

// ExampleUpdateMetric пример использования хендлера UpdateMetric.
// пример обращения к ендпоинту update/{metricType}/{metricName}/{metricValue}.
func ExampleUpdateMetric() {
	_ = config.InitDefaultEnv()
	configServe, _ := config.InitConfig()
	metricRepo, _ := local.NewMetricLocalRepository(configServe)

	metricService := service.NewMetricService(metricRepo)
	metricHandlers := NewUpdateHandlers(metricService)

	path := fmt.Sprintf("/update/%s/%s/%s", models.Counter, helpers.ValidCounterID, helpers.ValidCounterSTRING)
	req, _ := http.NewRequest(http.MethodGet, path, nil)
	rr := httptest.NewRecorder()
	metricHandlers.UpdateMetric(rr, req)
	fmt.Println(rr.Body.String())
}

// ExampleV2UpdateMetric пример использования хендлера UpdateMetric.
// пример обращения к ендпоинту /update/ .
func ExampleV2UpdateMetric() {
	_ = config.InitDefaultEnv()
	configServe, _ := config.InitConfig()
	metricRepo, _ := local.NewMetricLocalRepository(configServe)

	metricService := service.NewMetricService(metricRepo)
	metricHandlers := NewUpdateHandlers(metricService)

	metricJSON, _ := json.Marshal(helpers.GaugeMetric)

	req, _ := http.NewRequest(http.MethodPost, "/update/", bytes.NewBuffer(metricJSON))
	rr := httptest.NewRecorder()
	metricHandlers.UpdateMetric(rr, req)
	fmt.Println(rr.Body.String())
}

// ExampleUpdateBatchMetrics пример использования хендлера UpdateBatchMetrics.
// пример обращения к ендпоинту /updates/ .
func ExampleUpdateBatchMetrics() {
	_ = config.InitDefaultEnv()
	configServe, _ := config.InitConfig()
	metricRepo, _ := local.NewMetricLocalRepository(configServe)

	metricService := service.NewMetricService(metricRepo)
	metricHandlers := NewUpdateHandlers(metricService)

	metricJSON, _ := json.Marshal([]models.Metrics{helpers.CounterMetric, helpers.CounterMetric})

	req, _ := http.NewRequest(http.MethodPost, "/updates/", bytes.NewBuffer(metricJSON))
	rr := httptest.NewRecorder()
	metricHandlers.UpdateMetric(rr, req)
	fmt.Println(rr.Body.String())
}

// ExampleUpdateMetric пример использования хендлера UpdateMetric.
// пример обращения к ендпоинту /update/ .
func ExampleGetMetric() {
	_ = config.InitDefaultEnv()
	configServe, _ := config.InitConfig()
	metricRepo, _ := local.NewMetricLocalRepository(configServe)

	metricService := service.NewMetricService(metricRepo)
	metricHandlers := NewUpdateHandlers(metricService)

	metricJSON, _ := json.Marshal(helpers.CounterMetricRequest)

	req, _ := http.NewRequest(http.MethodPost, "/value/", bytes.NewBuffer(metricJSON))
	rr := httptest.NewRecorder()
	metricHandlers.UpdateMetric(rr, req)
	fmt.Println(rr.Body.String())
}
