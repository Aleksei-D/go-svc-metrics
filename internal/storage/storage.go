package storage

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"go-svc-metrics/internal/config"
	"go-svc-metrics/internal/logger"
	"go-svc-metrics/internal/storage/database"
	"go-svc-metrics/internal/storage/local"
	"go-svc-metrics/models"
)

const attemptsDefault uint = 3

type MetricRepository interface {
	UpdateMetrics(metrics []models.Metrics) ([]models.Metrics, error)
	GetMetric(metric models.Metrics) (models.Metrics, error)
	GetAllMetrics() ([]models.Metrics, error)
	Ping() error
	Close() error
	DumpMetricsByInterval(ctx context.Context) error
}

func NewMetricRepository(config *config.Config) (MetricRepository, error) {
	db, err := getDBConnect(*config.DatabaseDsn)
	if err != nil {
		logger.Log.Info(err.Error())
		return local.NewRetryWrapperLocalStorage(config, attemptsDefault)
	}
	return database.NewRetryMetricRepository(db, attemptsDefault)
}

func getDBConnect(databaseDsn string) (*sql.DB, error) {
	if databaseDsn != "" {
		db, err := sql.Open("postgres", databaseDsn)
		if err != nil {
			return nil, err
		}

		if err := db.Ping(); err != nil {
			return nil, err
		}
		return db, nil
	}
	return nil, fmt.Errorf("database connection is empty")
}
