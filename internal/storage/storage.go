package storage

import (
	"strconv"
	"sync"
)

type MemStorage struct {
	Metrics map[string]string
	mutex   sync.Mutex
}

type Repositories interface {
	UpdateGauge(metricName, value string) error
	UpdateCounter(metricName string, delta int64) error
	GetValue(metricName string) (string, bool)
	GetAllMetrics() map[string]string
}

func InitMemStorage() *MemStorage {
	return &MemStorage{
		Metrics: make(map[string]string),
	}
}

func (m *MemStorage) UpdateCounter(metricName string, delta int64) error {
	m.mutex.Lock()
	value, ok := m.Metrics[metricName]
	if !ok {
		m.Metrics[metricName] = strconv.FormatInt(delta, 10)
	} else {
		valueInt, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		m.Metrics[metricName] = strconv.FormatInt(delta+valueInt, 10)
	}
	m.mutex.Unlock()
	return nil
}

func (m *MemStorage) UpdateGauge(metricName, value string) error {
	m.mutex.Lock()
	m.Metrics[metricName] = value
	m.mutex.Unlock()
	return nil
}

func (m *MemStorage) GetAllMetrics() map[string]string {
	return m.Metrics
}

func (m *MemStorage) GetValue(metricName string) (string, bool) {
	value, ok := m.Metrics[metricName]
	return value, ok
}
