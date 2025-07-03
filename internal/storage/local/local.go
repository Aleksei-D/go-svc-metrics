package local

import (
	"bufio"
	"context"
	"encoding/json"
	"go-svc-metrics/internal/config"
	"go-svc-metrics/models"
	"os"
	"sync"
	"time"
)

type Storage struct {
	Metrics       map[string]models.Metrics
	mutex         sync.Mutex
	file          *os.File
	storeInterval time.Duration
	scanner       *bufio.Scanner
}

func NewLocalStorage(config *config.Config) (*Storage, error) {
	file, err := os.OpenFile(*config.FileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	localStorage := Storage{
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

func (l *Storage) UpdateMetric(metricToUpdate models.Metrics) (models.Metrics, error) {
	l.mutex.Lock()
	switch metricToUpdate.MType {
	case models.Gauge:
		l.Metrics[metricToUpdate.ID] = metricToUpdate
	case models.Counter:
		_, ok := l.Metrics[metricToUpdate.ID]
		if !ok {
			l.Metrics[metricToUpdate.ID] = metricToUpdate
		} else {
			*l.Metrics[metricToUpdate.ID].Delta += *metricToUpdate.Delta
		}
		metricToUpdate.Delta = l.Metrics[metricToUpdate.ID].Delta
	}
	l.mutex.Unlock()
	return metricToUpdate, nil
}

func (l *Storage) GetAllMetrics() (map[string]models.Metrics, error) {
	return l.Metrics, nil
}

func (l *Storage) GetValue(metric models.Metrics) (models.Metrics, bool) {
	value, ok := l.Metrics[metric.ID]
	return value, ok
}

func (l *Storage) Ping() bool {
	return false
}

func (l *Storage) Close() error {
	return l.file.Close()
}

func (l *Storage) RestoreMetrics() error {
	l.mutex.Lock()
	for l.scanner.Scan() {
		data := l.scanner.Bytes()
		metric := models.Metrics{}
		err := json.Unmarshal(data, &metric)
		if err != nil {
			return err
		}

		l.Metrics[metric.ID] = metric
	}
	l.mutex.Unlock()

	if err := l.scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (l *Storage) WriteMetric(metric *models.Metrics) error {
	data, err := json.Marshal(&metric)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = l.file.Write(data)
	return err
}

func (l *Storage) DumpMetrics() error {
	metricCopy := make(map[string]models.Metrics)
	l.mutex.Lock()
	for k, v := range l.Metrics {
		metricCopy[k] = v
	}
	l.mutex.Unlock()

	for _, metric := range metricCopy {
		err := l.WriteMetric(&metric)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *Storage) DumpMetricsByInterval(ctx context.Context) error {
	storeIntervalTicker := time.NewTicker(l.storeInterval)
	defer storeIntervalTicker.Stop()
	select {
	case <-storeIntervalTicker.C:
		return l.DumpMetrics()
	case <-ctx.Done():
		return l.DumpMetrics()
	}
}
