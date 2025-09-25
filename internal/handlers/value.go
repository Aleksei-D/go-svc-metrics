package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"go-svc-metrics/internal/service"
	errors2 "go-svc-metrics/internal/utils/errors"
	"go-svc-metrics/models"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ValueHandlers struct {
	metricService *service.MetricService
}

// NewValueHandlers создает и возвращает новый ValueHandlers.
func NewValueHandlers(metricService *service.MetricService) *ValueHandlers {
	return &ValueHandlers{metricService: metricService}
}

// GetMetricValue обработка ендпоинта GET /value/{metricType}/{metricName}.
// Возвращает значение метрики.
//
// Example:
//
//	http://localhost:8080/value/counter/CounterMetric
//
// Output:
//
// 4
func (m *ValueHandlers) GetMetricValue(res http.ResponseWriter, req *http.Request) {
	metricTypeFromPath := chi.URLParam(req, MetricTypePath)
	metricNameFromPath := chi.URLParam(req, MetricNamePath)
	value, err := m.metricService.GetMetricValue(req.Context(), metricTypeFromPath, metricNameFromPath)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			http.Error(res, err.Error(), http.StatusNotFound)
		case errors.Is(err, errors2.ErrInvalidMetricVType):
			http.Error(res, err.Error(), http.StatusBadRequest)
		}
		return
	}

	_, err = res.Write([]byte(value))
	if err != nil {
		http.Error(res, "invalid value", http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
}

// GetMetric обработка ендпоинта POST /value/ .
// Возвращает значение метрики.
//
// Example:
//
//	http://localhost:8080/value/
//
// Input:
//
//	{
//	    "id": "GaugeMetric",
//	    "type": "gauge"
//	}
//
// Output:
//
//	{
//	    "id": "GaugeMetric",
//	    "type": "gauge",
//	    "value": "1.0"
//	}
func (m *ValueHandlers) GetMetric(res http.ResponseWriter, req *http.Request) {
	var metricReq models.Metrics
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &metricReq); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	metric, err := m.metricService.GetMetric(req.Context(), metricReq)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}

	jsonData, err := json.Marshal(metric)
	if err != nil {
		http.Error(res, "invalid marshaling", http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(jsonData)
}
