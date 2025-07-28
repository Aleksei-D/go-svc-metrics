package local

import (
	"errors"
	"go-svc-metrics/internal/config"
	"go-svc-metrics/internal/utils/delay"
	"go-svc-metrics/models"
	"os"
	"time"
)

type RetryWrapperMetricLocalRepository struct {
	*MetricLocalRepository
	retryErrors error
	attempts    uint
}

func NewRetryWrapperLocalStorage(config *config.Config, attempts uint) (*RetryWrapperMetricLocalRepository, error) {
	storage, err := NewMetricLocalRepository(config)
	if err != nil {
		return nil, err
	}

	return &RetryWrapperMetricLocalRepository{
		MetricLocalRepository: storage,
		attempts:              attempts,
		retryErrors:           errors.Join([]error{os.ErrDeadlineExceeded, os.ErrNoDeadline, os.ErrInvalid}...),
	}, nil
}

func (r *RetryWrapperMetricLocalRepository) UpdateMetrics(metrics []models.Metrics) ([]models.Metrics, error) {
	var resultMetrics []models.Metrics
	delay := delay.NewDelay()
	for i := 0; i <= int(r.attempts); i++ {
		upsertMetrics, err := r.MetricLocalRepository.UpdateMetrics(metrics)
		resultMetrics = append(resultMetrics, upsertMetrics...)
		if errors.Is(err, r.retryErrors) {
			time.Sleep(delay())
			continue
		}

		if err != nil {
			return resultMetrics, err
		}
		return resultMetrics, nil
	}
	return resultMetrics, nil
}

func (r *RetryWrapperMetricLocalRepository) GetMetric(metric models.Metrics) (models.Metrics, error) {
	var resultMetric models.Metrics
	delay := delay.NewDelay()
	for i := 0; i <= int(r.attempts); i++ {
		fetchMetric, err := r.MetricLocalRepository.GetMetric(metric)
		resultMetric = fetchMetric
		if errors.Is(err, r.retryErrors) {
			time.Sleep(delay())
			continue
		}

		if err != nil {
			return resultMetric, err
		}
		return resultMetric, nil
	}
	return resultMetric, nil
}

func (r *RetryWrapperMetricLocalRepository) GetAllMetrics() ([]models.Metrics, error) {
	var resultMetrics []models.Metrics
	delay := delay.NewDelay()
	for i := 0; i <= int(r.attempts); i++ {
		fetchMetrics, err := r.MetricLocalRepository.GetAllMetrics()
		resultMetrics = append(resultMetrics, fetchMetrics...)
		if errors.Is(err, r.retryErrors) {
			time.Sleep(delay())
			continue
		}

		if err != nil {
			return resultMetrics, err
		}
		return resultMetrics, nil
	}
	return resultMetrics, nil
}
