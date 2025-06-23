package storage

import (
	"go-svc-metrics/models"
	"sync"
)

type MemStorage struct {
	Metrics map[string]models.Metrics
	mutex   sync.Mutex
}

type Repositories interface {
	UpdateMetric(metric models.Metrics) (models.Metrics, error)
	GetValue(metric models.Metrics) (models.Metrics, bool)
	GetAllMetrics() map[string]models.Metrics
}

func InitMemStorage() *MemStorage {
	return &MemStorage{
		Metrics: make(map[string]models.Metrics),
	}
}

func (m *MemStorage) UpdateMetric(metricToUpdate models.Metrics) (models.Metrics, error) {
	m.mutex.Lock()
	switch metricToUpdate.MType {
	case models.Gauge:
		m.Metrics[metricToUpdate.ID] = metricToUpdate
	case models.Counter:
		_, ok := m.Metrics[metricToUpdate.ID]
		if !ok {
			m.Metrics[metricToUpdate.ID] = metricToUpdate
		} else {
			*m.Metrics[metricToUpdate.ID].Delta += *metricToUpdate.Delta
		}
		metricToUpdate.Delta = m.Metrics[metricToUpdate.ID].Delta
	}
	m.mutex.Unlock()
	return metricToUpdate, nil
}

func (m *MemStorage) GetAllMetrics() map[string]models.Metrics {
	return m.Metrics
}

func (m *MemStorage) GetValue(metric models.Metrics) (models.Metrics, bool) {
	value, ok := m.Metrics[metric.ID]
	return value, ok
}
