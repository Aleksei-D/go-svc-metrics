package storage

import (
	"bufio"
	"encoding/json"
	"go-svc-metrics/models"
	"os"
	"time"
)

type Producer struct {
	file          *os.File
	storage       *MemStorage
	storeInterval time.Duration
}

func NewProducer(filename string, storage *MemStorage, storeInterval time.Duration) (*Producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{file: file, storage: storage, storeInterval: storeInterval}, nil
}

func (p *Producer) Close() error {
	return p.file.Close()
}

func (p *Producer) WriteMetric(metric *models.Metrics) error {
	data, err := json.Marshal(&metric)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = p.file.Write(data)
	return err
}

func (p *Producer) DumpMetrics() error {
	metricCopy := make(map[string]models.Metrics)
	p.storage.mutex.Lock()
	for k, v := range p.storage.Metrics {
		metricCopy[k] = v
	}
	p.storage.mutex.Unlock()

	for _, metric := range metricCopy {
		err := p.WriteMetric(&metric)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Producer) DumpMetricsByInterval() error {
	storeIntervalTicker := time.NewTicker(p.storeInterval)
	defer storeIntervalTicker.Stop()
	for range storeIntervalTicker.C {
		return p.DumpMetrics()
	}
	return nil
}

type Consumer struct {
	file    *os.File
	scanner *bufio.Scanner
}

func NewConsumer(filename string) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func (c *Consumer) RestoreMetrics(storage *MemStorage) error {
	storage.mutex.Lock()
	for c.scanner.Scan() {
		data := c.scanner.Bytes()
		metric := models.Metrics{}
		err := json.Unmarshal(data, &metric)
		if err != nil {
			return err
		}

		storage.Metrics[metric.ID] = metric
	}
	storage.mutex.Unlock()

	if err := c.scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (c *Consumer) Close() error {
	return c.file.Close()
}
