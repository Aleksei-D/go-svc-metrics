// модуль handlers релизуют хендлеры сервера по сбору метрик.
package handlers

import (
	"encoding/json"
	"go-svc-metrics/internal/service"
	errors2 "go-svc-metrics/internal/utils/errors"
	"go-svc-metrics/models"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Константы для работы с квер строкой.
const (
	MetricTypePath  = "metricType"
	MetricNamePath  = "metricName"
	MetricValuePath = "metricValue"
)

// MetricHandler хранит слой сервиса.
type UpdateHandlers struct {
	metricService *service.MetricService
}

// NewMetricHandler создает и возвращает новый MetricHandler.
func NewUpdateHandlers(metricService *service.MetricService) *UpdateHandlers {
	return &UpdateHandlers{metricService: metricService}
}

// V2UpdateMetric обработка ендпоинта POST /update/ .
// Возвращает значение метрики.
//
// Example:
//
//	http://localhost:8080/update/
//
// Input:
//
//	{
//	    "id": "GaugeMetric",
//	    "type": "gauge",
//	    "value": "1.0"
//	}
//
// Output:
//
//	{
//	    "id": "GaugeMetric",
//	    "type": "gauge",
//	    "value": "1.0"
//	}
func (m *UpdateHandlers) V2UpdateMetric(res http.ResponseWriter, req *http.Request) {
	var metric models.Metrics

	buf, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf, &metric); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	updatedMetrics, err := m.metricService.UpdateMetrics(req.Context(), []models.Metrics{metric})
	if err != nil {
		http.Error(res, errors2.ErrInvalidCounterOperation.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(updatedMetrics[0])
	if err != nil {
		http.Error(res, "invalid marshaling", http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(jsonData)
}

// UpdateBatchMetrics обработка ендпоинта POST /updates/ .
// Возвращает значение метрики.
//
// Example:
//
//	http://localhost:8080/updates/
//
// Input:
//
//	[
//	  {
//	    "id": "GaugeMetric",
//	    "type": "gauge",
//	    "value": "1.0"
//	  },
//	  {
//	    "id": "CounterMetric",
//	    "type": "counter",
//	    "delta": "4"
//	  }
//	]
//
// Output:
//
//	[
//	  {
//	    "id": "GaugeMetric",
//	    "type": "gauge",
//	    "value": "1.0"
//	  },
//	  {
//	    "id": "CounterMetric",
//	    "type": "counter",
//	    "delta": "4"
//	  }
//	]
func (m *UpdateHandlers) UpdateBatchMetrics(res http.ResponseWriter, req *http.Request) {
	var metrics []models.Metrics

	buf, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf, &metrics); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	updatedMetrics, err := m.metricService.UpdateMetrics(req.Context(), metrics)
	if err != nil {
		http.Error(res, "invalid counter operation", http.StatusBadRequest)
		return
	}

	jsonData, err := json.Marshal(updatedMetrics)
	if err != nil {
		http.Error(res, "invalid marshaling", http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(jsonData)
}

// UpdateMetric обработка ендпоинта POST /update/{metricType}/{metricName}/{metricValue}.
// Записывает в репозиторий метрику из квери.
//
// Example:
//
//	http://localhost:8080/update/counter/CounterMetric/4
func (m *UpdateHandlers) UpdateMetric(res http.ResponseWriter, req *http.Request) {
	metricType := chi.URLParam(req, MetricTypePath)
	metricNameFromPath := chi.URLParam(req, MetricNamePath)
	metricValueFromPath := chi.URLParam(req, MetricValuePath)
	err := m.metricService.UpdateMetric(req.Context(), metricType, metricNameFromPath, metricValueFromPath)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	res.WriteHeader(http.StatusOK)
}
