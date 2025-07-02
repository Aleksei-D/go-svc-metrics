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

type Repositories interface {
	UpdateMetric(metric models.Metrics) (models.Metrics, error)
	GetValue(metric models.Metrics) (models.Metrics, bool)
	GetAllMetrics() map[string]models.Metrics
	Ping() bool
	Close() error
	DumpMetricsByInterval(ctx context.Context) error
}

func InitRepositories(config *config.Config) (Repositories, error) {
	db, err := getDBConnect(*config.DatabaseDsn)
	if err != nil {
		logger.Log.Info(err.Error())
		return local.NewLocalStorage(config)
	}
	return database.NewDatabaseStorage(db)
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
