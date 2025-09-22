package domain

import (
	"context"
	"go-svc-metrics/models"
)

// MetricRepo интерфейс работы с репозиторием.
type MetricRepo interface {
	UpdateMetrics(ctx context.Context, metrics []models.Metrics) ([]models.Metrics, error)
	GetMetric(ctx context.Context, metric models.Metrics) (models.Metrics, error)
	GetAllMetrics(ctx context.Context) ([]models.Metrics, error)
	Ping() error
	Close() error
	DumpMetricsByInterval(ctx context.Context) error
}
