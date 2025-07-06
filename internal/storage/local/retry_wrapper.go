package local

import (
	"errors"
	"go-svc-metrics/internal/config"
	"go-svc-metrics/internal/utils"
	"go-svc-metrics/models"
	"os"
	"time"
)

type RetryWrapperLocalStorage struct {
	*Storage
	retryErrors error
	attempts    uint
}

func NewRetryWrapperLocalStorage(config *config.Config, attempts uint) (*RetryWrapperLocalStorage, error) {
	storage, err := NewLocalStorage(config)
	if err != nil {
		return nil, err
	}

	return &RetryWrapperLocalStorage{
		Storage:     storage,
		attempts:    attempts,
		retryErrors: errors.Join([]error{os.ErrDeadlineExceeded, os.ErrNoDeadline, os.ErrInvalid}...),
	}, nil
}

func (r *RetryWrapperLocalStorage) UpdateMetrics(metrics []models.Metrics) ([]models.Metrics, error) {
	var resultMetrics []models.Metrics
	delay := utils.GetDelay()
	for i := 0; i <= int(r.attempts); i++ {
		upsertMetrics, err := r.Storage.UpdateMetrics(metrics)
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

func (r *RetryWrapperLocalStorage) GetValue(metric models.Metrics) (models.Metrics, error) {
	var resultMetric models.Metrics
	delay := utils.GetDelay()
	for i := 0; i <= int(r.attempts); i++ {
		fetchMetric, err := r.Storage.GetValue(metric)
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

func (r *RetryWrapperLocalStorage) GetAllMetrics() ([]models.Metrics, error) {
	var resultMetrics []models.Metrics
	delay := utils.GetDelay()
	for i := 0; i <= int(r.attempts); i++ {
		fetchMetrics, err := r.Storage.GetAllMetrics()
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

func (r *RetryWrapperLocalStorage) Ping() error {
	delay := utils.GetDelay()
	for i := 0; i <= int(r.attempts); i++ {
		err := r.Storage.Ping()
		if errors.Is(err, r.retryErrors) {
			time.Sleep(delay())
			continue
		}

		if err != nil {
			return err
		}
		return err
	}
	return nil
}
