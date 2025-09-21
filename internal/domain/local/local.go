package local

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"go-svc-metrics/internal/config"
	"go-svc-metrics/models"
	"os"
	"sync"
	"time"
)

type MetricLocalRepository struct {
	Metrics       map[string]models.Metrics
	mutex         sync.Mutex
	file          *os.File
	storeInterval time.Duration
	scanner       *bufio.Scanner
}

func NewMetricLocalRepository(config *config.Config) (*MetricLocalRepository, error) {
	file, err := os.OpenFile(*config.FileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	localStorage := MetricLocalRepository{
		Metrics:       make(map[string]models.Metrics),
		file:          file,
		scanner:       bufio.NewScanner(file),
		storeInterval: config.GetStoreInterval(),
	}
	if *config.Restore {
		err = localStorage.RestoreMetrics()
		if err != nil {
			return nil, err
		}
	}
	return &localStorage, nil
}

func (m *MetricLocalRepository) UpdateMetrics(_ context.Context, metricsToUpdate []models.Metrics) ([]models.Metrics, error) {
	for _, metricToUpdate := range metricsToUpdate {
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
	}
	return metricsToUpdate, nil
}

func (m *MetricLocalRepository) GetAllMetrics(_ context.Context) ([]models.Metrics, error) {
	metrics := make([]models.Metrics, 0)
	for _, metric := range m.Metrics {
		metrics = append(metrics, metric)
	}
	return metrics, nil
}

func (m *MetricLocalRepository) GetMetric(_ context.Context, metric models.Metrics) (models.Metrics, error) {
	m.mutex.Lock()
	value, ok := m.Metrics[metric.ID]
	m.mutex.Unlock()
	if !ok {
		return models.Metrics{}, fmt.Errorf("metric not found")
	}
	return value, nil
}

func (m *MetricLocalRepository) Ping() error {
	return fmt.Errorf("is local storage (file)")
}

func (m *MetricLocalRepository) Close() error {
	return m.file.Close()
}

func (m *MetricLocalRepository) RestoreMetrics() error {
	m.mutex.Lock()
	for m.scanner.Scan() {
		data := m.scanner.Bytes()
		metric := models.Metrics{}
		err := json.Unmarshal(data, &metric)
		if err != nil {
			return err
		}

		m.Metrics[metric.ID] = metric
	}
	m.mutex.Unlock()

	if err := m.scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (m *MetricLocalRepository) WriteMetric(metric *models.Metrics) error {
	data, err := json.Marshal(&metric)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = m.file.Write(data)
	return err
}

func (m *MetricLocalRepository) DumpMetricsByInterval(ctx context.Context) error {
	storeIntervalTicker := time.NewTicker(m.storeInterval)
	defer storeIntervalTicker.Stop()
	select {
	case <-storeIntervalTicker.C:
		return m.DumpMetrics()
	case <-ctx.Done():
		return m.DumpMetrics()
	}
}

func (m *MetricLocalRepository) DumpMetrics() error {
	metricCopy := make(map[string]models.Metrics)
	m.mutex.Lock()
	for k, v := range m.Metrics {
		metricCopy[k] = v
	}
	m.mutex.Unlock()

	for _, metric := range metricCopy {
		err := m.WriteMetric(&metric)
		if err != nil {
			return err
		}
	}
	return nil
}
