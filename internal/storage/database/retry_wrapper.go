package database

import (
	"database/sql"
	"github.com/jackc/pgerrcode"
	"go-svc-metrics/internal/utils/delay"
	"go-svc-metrics/models"
	"time"
)

func NewRetryMetricRepository(db *sql.DB, attempts uint) (*RetryMetricRepository, error) {
	metricRepository, err := NewMetricRepository(db)
	if err != nil {
		return nil, err
	}

	return &RetryMetricRepository{MetricDatabaseRepository: metricRepository, attempts: attempts}, nil
}

type RetryMetricRepository struct {
	*MetricDatabaseRepository
	attempts uint
}

func (r *RetryMetricRepository) UpdateMetrics(metrics []models.Metrics) ([]models.Metrics, error) {
	var resultMetrics []models.Metrics
	delay := delay.NewDelay()
	for i := 0; i <= int(r.attempts); i++ {
		upsertMetrics, err := r.MetricDatabaseRepository.UpdateMetrics(metrics)
		resultMetrics = append(resultMetrics, upsertMetrics...)
		if err != nil {
			if pgerrcode.IsConnectionException(err.Error()) || pgerrcode.IsIntegrityConstraintViolation(err.Error()) {
				time.Sleep(delay())
				continue
			}
			return resultMetrics, err
		}
		return resultMetrics, err
	}
	return resultMetrics, nil
}

func (r *RetryMetricRepository) GetValue(metric models.Metrics) (models.Metrics, error) {
	var resultMetric models.Metrics
	delay := delay.NewDelay()
	for i := 0; i <= int(r.attempts); i++ {
		fetchMetric, err := r.MetricDatabaseRepository.GetMetric(metric)
		resultMetric = fetchMetric
		if err != nil {
			if pgerrcode.IsConnectionException(err.Error()) || pgerrcode.IsIntegrityConstraintViolation(err.Error()) {
				time.Sleep(delay())
				continue
			}
			return resultMetric, err
		}
		return resultMetric, err
	}
	return resultMetric, nil
}

func (r *RetryMetricRepository) GetAllMetrics() ([]models.Metrics, error) {
	var resultMetrics []models.Metrics
	delay := delay.NewDelay()
	for i := 0; i <= int(r.attempts); i++ {
		fetchMetrics, err := r.MetricDatabaseRepository.GetAllMetrics()
		resultMetrics = append(resultMetrics, fetchMetrics...)
		if err != nil {
			if pgerrcode.IsConnectionException(err.Error()) || pgerrcode.IsIntegrityConstraintViolation(err.Error()) {
				time.Sleep(delay())
				continue
			}
			return resultMetrics, err
		}
		return resultMetrics, err
	}
	return resultMetrics, nil
}

func (r *RetryMetricRepository) Ping() error {
	delay := delay.NewDelay()
	for i := 0; i <= int(r.attempts); i++ {
		err := r.MetricDatabaseRepository.Ping()
		if err != nil {
			if pgerrcode.IsConnectionException(err.Error()) {
				time.Sleep(delay())
				continue
			}
			return err
		}
		return err
	}
	return nil
}
