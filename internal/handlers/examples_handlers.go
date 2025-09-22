package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-svc-metrics/internal/config"
	"go-svc-metrics/internal/domain/local"
	"go-svc-metrics/internal/service"
	"go-svc-metrics/models"
	"net/http"
	"net/http/httptest"
)

var (
	counterID     string           = "CounterMetric"
	validCounter  string           = "4"
	counterValue  int64            = 4
	gaugeValue    float64          = 1.0
	gaugeID       string           = "GaugeMetric"
	gaugeMetric   models.Metrics   = models.Metrics{ID: gaugeID, MType: models.Gauge, Value: &gaugeValue}
	counterMetric models.Metrics   = models.Metrics{ID: counterID, MType: models.Counter, Delta: &counterValue}
	metrics       []models.Metrics = []models.Metrics{gaugeMetric, counterMetric}
)

// ExampleUpdateMetric пример использования хендлера UpdateMetric.
// пример обращения к ендпоинту update/{metricType}/{metricName}/{metricValue}.
func ExampleUpdateMetric() {
	_ = config.InitDefaultEnv()
	configServe, _ := config.InitConfig()
	metricRepo, _ := local.NewMetricLocalRepository(configServe)

	metricService := service.NewMetricService(metricRepo)
	metricHandlers := NewMetricHandler(metricService)

	path := fmt.Sprintf("/update/%s/%s/%s", models.Counter, counterID, validCounter)
	req, _ := http.NewRequest(http.MethodGet, path, nil)
	rr := httptest.NewRecorder()
	metricHandlers.UpdateMetric(rr, req)
	fmt.Println(rr.Body.String())
}

// ExampleGetMetricValue пример использования хендлера UpdateMetric.
// пример обращения к ендпоинту /value/{metricType}/{metricName}.
func ExampleGetMetricValue() {
	_ = config.InitDefaultEnv()
	configServe, _ := config.InitConfig()
	metricRepo, _ := local.NewMetricLocalRepository(configServe)

	metricService := service.NewMetricService(metricRepo)
	metricHandlers := NewMetricHandler(metricService)

	path := fmt.Sprintf("/value/%s/%s", models.Counter, counterID)
	req, _ := http.NewRequest(http.MethodPost, path, nil)
	rr := httptest.NewRecorder()
	metricHandlers.GetMetricValue(rr, req)
	fmt.Println(rr.Body.String())
}

// ExampleGetMetrics пример использования хендлера UpdateMetric.
// пример обращения к ендпоинту / .
func ExampleGetMetrics() {
	_ = config.InitDefaultEnv()
	configServe, _ := config.InitConfig()
	metricRepo, _ := local.NewMetricLocalRepository(configServe)

	metricService := service.NewMetricService(metricRepo)
	metricHandlers := NewMetricHandler(metricService)

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	metricHandlers.GetMetricValue(rr, req)
	fmt.Println(rr.Body.String())
}

// ExampleV2UpdateMetric пример использования хендлера UpdateMetric.
// пример обращения к ендпоинту /update/ .
func ExampleV2UpdateMetric() {
	_ = config.InitDefaultEnv()
	configServe, _ := config.InitConfig()
	metricRepo, _ := local.NewMetricLocalRepository(configServe)

	metricService := service.NewMetricService(metricRepo)
	metricHandlers := NewMetricHandler(metricService)

	metricJSON, _ := json.Marshal(gaugeMetric)

	req, _ := http.NewRequest(http.MethodPost, "/update/", bytes.NewBuffer(metricJSON))
	rr := httptest.NewRecorder()
	metricHandlers.UpdateMetric(rr, req)
	fmt.Println(rr.Body.String())
}

// ExampleGetMetric пример использования хендлера GetMetric.
// пример обращения к ендпоинту /update/ .
func ExampleGetMetric() {
	_ = config.InitDefaultEnv()
	configServe, _ := config.InitConfig()
	metricRepo, _ := local.NewMetricLocalRepository(configServe)

	metricService := service.NewMetricService(metricRepo)
	metricHandlers := NewMetricHandler(metricService)

	metricJSON, _ := json.Marshal(models.Metrics{ID: counterID, MType: models.Counter})

	req, _ := http.NewRequest(http.MethodPost, "/value/", bytes.NewBuffer(metricJSON))
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
	metricHandlers := NewMetricHandler(metricService)

	metricJSON, _ := json.Marshal(metrics)

	req, _ := http.NewRequest(http.MethodPost, "/updates/", bytes.NewBuffer(metricJSON))
	rr := httptest.NewRecorder()
	metricHandlers.UpdateMetric(rr, req)
	fmt.Println(rr.Body.String())
}
