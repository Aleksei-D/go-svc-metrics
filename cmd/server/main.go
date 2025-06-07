package main

import (
	"go-svc-metrics/internal/config"
	"go-svc-metrics/models"
	"net/http"
	"strconv"
)

const (
	metricTypePath  = "metricType"
	metricNamePath  = "metricName"
	metricValuePath = "metricValue"
)

type MemStorage struct {
	CounterMetrics map[string]int64
	GaugeMetrics   map[string]float64
}

var memStorage MemStorage

func CounterMetricProcessing(metricName string, value int64) {
	_, ok := memStorage.CounterMetrics[metricName]
	if !ok {
		memStorage.CounterMetrics[metricName] = value
	} else {
		memStorage.CounterMetrics[metricName] += value
	}
}

func GaugeMetricProcessing(metricName string, value float64) {
	_, ok := memStorage.GaugeMetrics[metricName]
	if !ok {
		memStorage.GaugeMetrics[metricName] = value
	} else {
		memStorage.GaugeMetrics[metricName] += value
	}
}

func metricsPage(res http.ResponseWriter, req *http.Request) {
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
		CounterMetricProcessing(metricNameFromPath, int64(metricValue))
	case models.Gauge:
		metricValue, err := strconv.ParseFloat(metricValueFromPath, 64)
		if err != nil {
			http.Error(res, "invalid metric value", http.StatusBadRequest)
			return
		}
		GaugeMetricProcessing(metricNameFromPath, metricValue)
	default:
		http.Error(res, "invalid metric type", http.StatusBadRequest)
		return
	}
}

func main() {
	serveConfig := config.GetServeConfig()
	memStorage = MemStorage{
		CounterMetrics: make(map[string]int64),
		GaugeMetrics:   make(map[string]float64),
	}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /update/{metricType}/{metricName}/{metricValue}", metricsPage)

	err := http.ListenAndServe(serveConfig.GetServeAddress(), mux)
	if err != nil {
		panic(err)
	}
}
