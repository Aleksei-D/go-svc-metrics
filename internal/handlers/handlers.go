package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"go-svc-metrics/internal/storage"
	"go-svc-metrics/models"
	"net/http"
	"strconv"
)

const (
	MetricTypePath          = "metricType"
	MetricNamePath          = "metricName"
	MetricValuePath         = "metricValue"
	UpdateMetricHandlerPath = "/update/{metricType}/{metricName}/{metricValue}"
	GetMetricHandlerPath    = "/value/{metricType}/{metricName}"
)

type MetricHandler struct {
	Storage storage.Repositories
}

func (m *MetricHandler) UpdateMetric(res http.ResponseWriter, req *http.Request) {
	metricType := chi.URLParam(req, MetricTypePath)
	metricNameFromPath := chi.URLParam(req, MetricNamePath)
	metricValueFromPath := chi.URLParam(req, MetricValuePath)
	switch metricType {
	case models.Counter:
		metricValue, err := strconv.Atoi(metricValueFromPath)
		if err != nil {
			http.Error(res, "invalid metric value", http.StatusBadRequest)
			return
		}
		err = m.Storage.UpdateCounter(metricNameFromPath, int64(metricValue))
		if err != nil {
			http.Error(res, "invalid counter operation", http.StatusBadRequest)
			return
		}
	case models.Gauge:
		_, err := strconv.ParseFloat(metricValueFromPath, 64)
		if err != nil {
			http.Error(res, "invalid metric value", http.StatusBadRequest)
			return
		}
		err = m.Storage.UpdateGauge(metricNameFromPath, metricValueFromPath)
		if err != nil {
			http.Error(res, "invalid gauge operation", http.StatusBadRequest)
			return
		}
	default:
		http.Error(res, "invalid metric type", http.StatusBadRequest)
		return
	}
}

func (m *MetricHandler) GetMetricByName(res http.ResponseWriter, req *http.Request) {
	metricTypeFromPath := chi.URLParam(req, MetricTypePath)
	metricNameFromPath := chi.URLParam(req, MetricNamePath)
	if metricTypeFromPath != models.Counter && metricTypeFromPath != models.Gauge {
		http.Error(res, "invalid metric type", http.StatusBadRequest)
		return
	}
	value, ok := m.Storage.GetValue(metricNameFromPath)
	if !ok {
		http.Error(res, "invalid metric type", http.StatusNotFound)
		return
	}

	res.WriteHeader(http.StatusOK)
	_, err := res.Write([]byte(value))
	if err != nil {
		http.Error(res, "invalid value", http.StatusBadRequest)
		return
	}
}

func (m *MetricHandler) GetMetrics(res http.ResponseWriter, req *http.Request) {
	metrics := m.Storage.GetAllMetrics()
	jsonString, err := json.Marshal(metrics)
	if err != nil {
		http.Error(res, "invalid marshaling", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(jsonString)
}
