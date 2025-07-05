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
	_, err := db.Exec(CreateTableSQL)
	if err != nil {
		return nil, err
	}
	return &Storage{db}, nil
}

func (d *Storage) Ping() bool {
	if err := d.DB.Ping(); err != nil {
		return false
	}
	return true
}

func (d *Storage) UpdateMetrics(metrics []models.Metrics) ([]models.Metrics, error) {
	tx, err := d.DB.Begin()
	if err != nil {
		return metrics, err
	}

	stmt, err := tx.Prepare("INSERT INTO metric_table as t1 (name_id, delta, value) VALUES ($1, $2, $3) " +
		"ON CONFLICT (name_id) DO UPDATE SET delta = t1.delta + EXCLUDED.delta, value = $3 RETURNING delta, value")
	if err != nil {
		return metrics, err
	}

	for _, metric := range metrics {
		var delta sql.NullInt64
		var value sql.NullFloat64

		row := stmt.QueryRow(metric.ID, metric.Delta, metric.Value)
		err = row.Scan(&delta, &value)
		if err != nil {
			return metrics, err
		}

		if delta.Valid {
			*metric.Delta = delta.Int64
		}
		if value.Valid {
			*metric.Value = value.Float64
		}
	}
	err = tx.Commit()
	if err != nil {
		return metrics, err
	}
	return metrics, nil
}

func (d *Storage) GetValue(metric models.Metrics) (models.Metrics, bool) {
	var delta sql.NullInt64
	var value sql.NullFloat64
	row := d.DB.QueryRow("SELECT delta, value FROM metric_table WHERE name_id = $1", metric.ID)
	err := row.Scan(&delta, &value)
	if err != nil {
		return metric, false
	}

	if delta.Valid {
		metric.Delta = &delta.Int64
	}
	if value.Valid {
		metric.Value = &value.Float64
	}
	return metric, true
}

func (d *Storage) GetAllMetrics() ([]models.Metrics, error) {
	metrics := make([]models.Metrics, 0)
	rows, err := d.DB.Query("SELECT name_id, delta, value FROM metric_table")
	if err != nil {
		return metrics, err
	}
	if err := rows.Err(); err != nil {
		return metrics, err
	}
	defer rows.Close()

	for rows.Next() {
		var metric models.Metrics
		err := rows.Scan(&metric.ID, &metric.Delta, &metric.Value)
		if err != nil {
			return metrics, err
		}
		metrics = append(metrics, metric)
	}
	return metrics, nil
}

func (d *Storage) Close() error { return d.DB.Close() }

func (d *Storage) DumpMetricsByInterval(_ context.Context) error {
	return nil
}
