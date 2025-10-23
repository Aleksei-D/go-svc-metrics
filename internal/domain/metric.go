package domain

import (
	"context"
	"go-svc-metrics/internal/config"
	"go-svc-metrics/internal/datasource"
	"go-svc-metrics/internal/domain/local"
	"go-svc-metrics/internal/domain/postgres"
	"go-svc-metrics/internal/logger"
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

func NewRepo(cfg *config.Config) (MetricRepo, error) {
	var metricRepo MetricRepo
	db, err := datasource.NewDatabase(*cfg.DatabaseDsn)
	if err != nil {
		logger.Log.Warn("not connected db")
	}

	switch db {
	case nil:
		localRepo, err := local.NewMetricLocalRepository(cfg)
		if err != nil {
			return metricRepo, err
		}
		metricRepo = localRepo
	default:
		postgresRepo := postgres.NewMetricRepository(db)
		metricRepo = postgresRepo
	}
	return metricRepo, nil
}
