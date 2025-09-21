package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"go-svc-metrics/internal/service"
	errors2 "go-svc-metrics/internal/utils/errors"
	"go-svc-metrics/models"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const (
	MetricTypePath  = "metricType"
	MetricNamePath  = "metricName"
	MetricValuePath = "metricValue"
)

type MetricHandler struct {
	metricService *service.MetricService
}

func NewMetricHandler(metricService *service.MetricService) *MetricHandler {
	return &MetricHandler{metricService: metricService}
}

func (m *MetricHandler) UpdateMetric(res http.ResponseWriter, req *http.Request) {
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

func (m *MetricHandler) GetMetricValue(res http.ResponseWriter, req *http.Request) {
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

func (m *MetricHandler) GetMetrics(res http.ResponseWriter, req *http.Request) {
	metrics, err := m.metricService.GetAllMetrics(req.Context())
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

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

func (m *MetricHandler) UpdateBatchMetrics(res http.ResponseWriter, req *http.Request) {
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

func (m *MetricHandler) GetPing(res http.ResponseWriter, req *http.Request) {
	err := m.metricService.Ping()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}
