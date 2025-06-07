package handlers

import (
	"go-svc-metrics/internal/storage"
	"go-svc-metrics/models"
	"net/http"
	"strconv"
)

const (
	metricTypePath    = "metricType"
	metricNamePath    = "metricName"
	metricValuePath   = "metricValue"
	MetricHandlerPath = "POST /update/{metricType}/{metricName}/{metricValue}"
)

type MetricHandler struct {
	Storage storage.Repositories
}

func (m *MetricHandler) Serve(res http.ResponseWriter, req *http.Request) {
	metricType := req.PathValue(metricTypePath)
	metricNameFromPath := req.PathValue(metricNamePath)
	metricValueFromPath := req.PathValue(metricValuePath)
	switch metricType {
	case models.Counter:
		metricValue, err := strconv.Atoi(metricValueFromPath)
		if err != nil {
			http.Error(res, "invalid metric value", http.StatusBadRequest)
			return
		}
		m.Storage.UpdateCounter(metricNameFromPath, int64(metricValue))
	case models.Gauge:
		metricValue, err := strconv.ParseFloat(metricValueFromPath, 64)
		if err != nil {
			http.Error(res, "invalid metric value", http.StatusBadRequest)
			return
		}
		m.Storage.UpdateGauge(metricNameFromPath, metricValue)
	default:
		http.Error(res, "invalid metric type", http.StatusBadRequest)
		return
	}
}
