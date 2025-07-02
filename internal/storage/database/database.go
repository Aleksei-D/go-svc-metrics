package database

import (
	"context"
	"database/sql"
	"go-svc-metrics/models"
)

type Storage struct {
	*sql.DB
}

func NewDatabaseStorage(db *sql.DB) (*Storage, error) {
	return &Storage{db}, nil
}

func (d *Storage) Ping() bool {
	if err := d.DB.Ping(); err != nil {
		return false
	}
	return true
}

func (d *Storage) UpdateMetric(_ models.Metrics) (models.Metrics, error) {
	return models.Metrics{}, nil
}

func (d *Storage) GetValue(_ models.Metrics) (models.Metrics, bool) {
	return models.Metrics{}, false
}

func (d *Storage) GetAllMetrics() map[string]models.Metrics {
	metrics := make(map[string]models.Metrics)
	return metrics
}

func (d *Storage) Close() error { return d.DB.Close() }

func (d *Storage) DumpMetricsByInterval(_ context.Context) error {
	return nil
}
