package service

import (
	"context"
	"go-svc-metrics/internal/domain"
	errors2 "go-svc-metrics/internal/utils/custom_errors"
	"go-svc-metrics/models"
	"strconv"
)

type MetricService struct {
	metricRepo domain.MetricRepo
}

func NewMetricService(metricRepo domain.MetricRepo) *MetricService {
	return &MetricService{metricRepo: metricRepo}
}

func (m *MetricService) UpdateMetric(ctx context.Context, metricType, metricName, metricValue string) error {
	metric := models.Metrics{
		ID:    metricName,
		MType: metricType,
	}
	switch metricType {
	case models.Counter:
		metricValue, err := strconv.Atoi(metricValue)
		if err != nil {
			return errors2.ErrInvalidMetricValue
		}
		value := int64(metricValue)
		metric.Delta = &value
		_, err = m.metricRepo.UpdateMetrics(ctx, []models.Metrics{metric})
		if err != nil {
			return errors2.ErrInvalidCounterOperation
		}
	case models.Gauge:
		metricValue, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return errors2.ErrInvalidMetricValue
		}
		metric.Value = &metricValue
		_, err = m.metricRepo.UpdateMetrics(ctx, []models.Metrics{metric})
		if err != nil {
			return errors2.ErrInvalidCGaugeOperation
		}
	default:
		return errors2.ErrInvalidMetricVType
	}
	return nil
}

func (m *MetricService) UpdateMetrics(ctx context.Context, metrics []models.Metrics) ([]models.Metrics, error) {
	return m.metricRepo.UpdateMetrics(ctx, metrics)
}

func (m *MetricService) GetAllMetrics(ctx context.Context) ([]models.Metrics, error) {
	return m.metricRepo.GetAllMetrics(ctx)
}

func (m *MetricService) GetMetricValue(ctx context.Context, metricType, metricName string) (string, error) {
	var value string
	if metricType != models.Counter && metricType != models.Gauge {
		return value, errors2.ErrInvalidMetricVType
	}
	metric, err := m.metricRepo.GetMetric(ctx, models.Metrics{MType: metricType, ID: metricName})
	if err != nil {
		return value, err
	}

	switch metric.MType {
	case models.Counter:
		value = strconv.Itoa(int(*metric.Delta))
	case models.Gauge:
		value = strconv.FormatFloat(*metric.Value, 'f', -1, 64)
	}
	return value, nil
}

func (m *MetricService) GetMetric(ctx context.Context, metric models.Metrics) (models.Metrics, error) {
	return m.metricRepo.GetMetric(ctx, metric)
}

func (m *MetricService) Ping() error {
	return m.metricRepo.Ping()
}

func (m *MetricService) Close() error {
	return m.metricRepo.Close()
}

func (m *MetricService) DumpMetricsByInterval(ctx context.Context) error {
	return m.metricRepo.DumpMetricsByInterval(ctx)
}
