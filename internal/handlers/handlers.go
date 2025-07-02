package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"go-svc-metrics/internal/storage"
	"go-svc-metrics/models"
	"io"
	"net/http"
	"strconv"
)

const (
	MetricTypePath  = "metricType"
	MetricNamePath  = "metricName"
	MetricValuePath = "metricValue"
)

type MetricHandler struct {
	Storage storage.Repositories
}

func (m *MetricHandler) UpdateMetric(res http.ResponseWriter, req *http.Request) {
	metricType := chi.URLParam(req, MetricTypePath)
	metricNameFromPath := chi.URLParam(req, MetricNamePath)
	metricValueFromPath := chi.URLParam(req, MetricValuePath)
	metric := models.Metrics{
		ID:    metricNameFromPath,
		MType: metricType,
	}
	switch metricType {
	case models.Counter:
		metricValue, err := strconv.Atoi(metricValueFromPath)
		if err != nil {
			http.Error(res, "invalid metric value", http.StatusBadRequest)
			return
		}
		value := int64(metricValue)
		metric.Delta = &value
		_, err = m.Storage.UpdateMetric(metric)
		if err != nil {
			http.Error(res, "invalid counter operation", http.StatusBadRequest)
			return
		}
		res.WriteHeader(http.StatusOK)
	case models.Gauge:
		metricValue, err := strconv.ParseFloat(metricValueFromPath, 64)
		if err != nil {
			http.Error(res, "invalid metric value", http.StatusBadRequest)
			return
		}
		metric.Value = &metricValue
		_, err = m.Storage.UpdateMetric(metric)
		if err != nil {
			http.Error(res, "invalid gauge operation", http.StatusBadRequest)
			return
		}
		res.WriteHeader(http.StatusOK)
	default:
		http.Error(res, "invalid metric type", http.StatusBadRequest)
		return
	}
}

func (m *MetricHandler) GetMetricValue(res http.ResponseWriter, req *http.Request) {
	var value string
	metricTypeFromPath := chi.URLParam(req, MetricTypePath)
	metricNameFromPath := chi.URLParam(req, MetricNamePath)
	if metricTypeFromPath != models.Counter && metricTypeFromPath != models.Gauge {
		http.Error(res, "invalid metric type", http.StatusBadRequest)
		return
	}

	metric, ok := m.Storage.GetValue(models.Metrics{
		ID:    metricNameFromPath,
		MType: metricTypeFromPath,
	})
	if !ok {
		http.Error(res, "invalid metric type", http.StatusNotFound)
		return
	}
	switch metric.MType {
	case models.Counter:
		value = strconv.Itoa(int(*metric.Delta))
	case models.Gauge:
		value = strconv.FormatFloat(*metric.Value, 'f', -1, 64)
	}

	res.WriteHeader(http.StatusOK)
	_, err := res.Write([]byte(value))
	if err != nil {
		http.Error(res, "invalid value", http.StatusBadRequest)
		return
	}
}

func (m *MetricHandler) GetMetrics(res http.ResponseWriter, _ *http.Request) {
	metrics := m.Storage.GetAllMetrics()
	jsonString, err := json.Marshal(metrics)
	if err != nil {
		http.Error(res, "invalid marshaling", http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "text/html")
	res.WriteHeader(http.StatusOK)
	res.Write(jsonString)
}

func (m *MetricHandler) V2UpdateMetric(res http.ResponseWriter, req *http.Request) {
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

	updatedMetric, err := m.Storage.UpdateMetric(metric)
	if err != nil {
		http.Error(res, "invalid counter operation", http.StatusBadRequest)
		return
	}

	jsonData, err := json.Marshal(updatedMetric)
	if err != nil {
		http.Error(res, "invalid marshaling", http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(jsonData)
}

func (m *MetricHandler) GetMetric(res http.ResponseWriter, req *http.Request) {
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

	metric, ok := m.Storage.GetValue(metricReq)
	if !ok {
		http.Error(res, "invalid counter operation", http.StatusNotFound)
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

func (m *MetricHandler) GetPing(res http.ResponseWriter, _ *http.Request) {
	if ok := m.Storage.Ping(); !ok {
		http.Error(res, "Not connect to DB", http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}
