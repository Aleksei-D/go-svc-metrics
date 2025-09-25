// модуль handlers релизуют хендлеры сервера по сбору метрик.
package handlers

import (
	"encoding/json"
	"go-svc-metrics/internal/service"
	"net/http"
)

// MetricHandler хранит слой сервиса.
type CommonHandlers struct {
	metricService *service.MetricService
}

// NewCommonHandlers создает и возвращает новый MetricHandler.
func NewCommonHandlers(metricService *service.MetricService) *CommonHandlers {
	return &CommonHandlers{metricService: metricService}
}

// GetMetrics обработка ендпоинта GET / .
// Возвращает значение метрики.
//
// Example:
//
//	http://localhost:8080/
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
func (m *CommonHandlers) GetMetrics(res http.ResponseWriter, req *http.Request) {
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

// GetPing проверяет подключение к БД.
func (m *CommonHandlers) GetPing(res http.ResponseWriter, req *http.Request) {
	err := m.metricService.Ping()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}
