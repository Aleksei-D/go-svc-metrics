package database

import (
	"database/sql"
	"github.com/jackc/pgerrcode"
	"go-svc-metrics/internal/utils"
	"go-svc-metrics/models"
	"time"
)

type RetryWrapperStorage struct {
	*Storage
	attempts uint
}

func NewRetryWrapperStorage(db *sql.DB, attempts uint) (*RetryWrapperStorage, error) {
	storage, err := NewDatabaseStorage(db)
	if err != nil {
		return nil, err
	}

	return &RetryWrapperStorage{
		Storage:  storage,
		attempts: attempts,
	}, nil
}

func (r *RetryWrapperStorage) UpdateMetrics(metrics []models.Metrics) ([]models.Metrics, error) {
	var resultMetrics []models.Metrics
	delay := utils.GetDelay()
	for i := 0; i <= int(r.attempts); i++ {
		upsertMetrics, err := r.Storage.UpdateMetrics(metrics)
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

func (r *RetryWrapperStorage) GetValue(metric models.Metrics) (models.Metrics, error) {
	var resultMetric models.Metrics
	delay := utils.GetDelay()
	for i := 0; i <= int(r.attempts); i++ {
		fetchMetric, err := r.Storage.GetValue(metric)
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

func (r *RetryWrapperStorage) GetAllMetrics() ([]models.Metrics, error) {
	var resultMetrics []models.Metrics
	delay := utils.GetDelay()
	for i := 0; i <= int(r.attempts); i++ {
		fetchMetrics, err := r.Storage.GetAllMetrics()
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

func (r *RetryWrapperStorage) Ping() error {
	delay := utils.GetDelay()
	for i := 0; i <= int(r.attempts); i++ {
		err := r.Storage.Ping()
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
