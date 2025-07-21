package usecase

import (
	"context"
	"go-svc-metrics/internal/storage"
	"go-svc-metrics/models"
)

type MetricUseCase struct {
	metricRepo storage.MetricRepository
}

func NewMetricUseCase(metricRepo storage.MetricRepository) *MetricUseCase {
	return &MetricUseCase{metricRepo: metricRepo}
}

func (u *MetricUseCase) UpdateMetrics(metrics []models.Metrics) ([]models.Metrics, error) {
	return u.metricRepo.UpdateMetrics(metrics)
}

func (u *MetricUseCase) GetAllMetrics() ([]models.Metrics, error) {
	return u.metricRepo.GetAllMetrics()
}

func (u *MetricUseCase) GetMetric(metric models.Metrics) (models.Metrics, error) {
	return u.metricRepo.GetMetric(metric)
}

func (u *MetricUseCase) Ping() error {
	return u.metricRepo.Ping()
}

func (u *MetricUseCase) Close() error {
	return u.metricRepo.Close()
}

func (u *MetricUseCase) DumpMetricsByInterval(ctx context.Context) error {
	return u.metricRepo.DumpMetricsByInterval(ctx)
}
