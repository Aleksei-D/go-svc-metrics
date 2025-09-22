// Модуль представляет реализацию интерфейса для работы с посгрес.
package postgres

import (
	"context"
	"database/sql"
	"go-svc-metrics/models"
)

type PostgresMetricRepository struct {
	db *sql.DB
}

func NewMetricRepository(db *sql.DB) *PostgresMetricRepository {
	return &PostgresMetricRepository{db: db}
}

func (m *PostgresMetricRepository) Ping() error {
	return m.db.Ping()
}

func (m *PostgresMetricRepository) Close() error { return m.db.Close() }

func (m *PostgresMetricRepository) DumpMetricsByInterval(_ context.Context) error {
	return nil
}

func (m *PostgresMetricRepository) UpdateMetrics(ctx context.Context, metrics []models.Metrics) ([]models.Metrics, error) {
	tx, err := m.db.Begin()
	if err != nil {
		return metrics, err
	}
	query := `INSERT INTO metric_table as t1 (name_id, type, delta, value) VALUES ($1, $2, $3, $4) 
    ON CONFLICT (name_id) DO UPDATE SET delta = t1.delta + EXCLUDED.delta, value = $4 RETURNING delta, value`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return metrics, err
	}

	for _, metric := range metrics {
		var delta sql.NullInt64
		var value sql.NullFloat64

		row := stmt.QueryRowContext(ctx, metric.ID, metric.MType, metric.Delta, metric.Value)
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

func (m *PostgresMetricRepository) GetMetric(ctx context.Context, metric models.Metrics) (models.Metrics, error) {
	var delta sql.NullInt64
	var value sql.NullFloat64
	query := `SELECT delta, value FROM metric_table WHERE name_id = $1 and type = $2`
	row := m.db.QueryRowContext(ctx, query, metric.ID, metric.MType)
	err := row.Scan(&delta, &value)
	if err != nil {
		return metric, err
	}

	if delta.Valid {
		metric.Delta = &delta.Int64
	}
	if value.Valid {
		metric.Value = &value.Float64
	}
	return metric, nil
}

func (m *PostgresMetricRepository) GetAllMetrics(ctx context.Context) ([]models.Metrics, error) {
	metrics := make([]models.Metrics, 0)
	rows, err := m.db.QueryContext(ctx, "SELECT name_id, type, delta, value FROM metric_table")
	if err != nil {
		return metrics, err
	}
	if err := rows.Err(); err != nil {
		return metrics, err
	}
	defer rows.Close()

	for rows.Next() {
		var delta sql.NullInt64
		var value sql.NullFloat64
		var metric models.Metrics
		err := rows.Scan(&metric.ID, &metric.MType, &delta, &value)
		if delta.Valid {
			metric.Delta = &delta.Int64
		}
		if value.Valid {
			metric.Value = &value.Float64
		}
		if err != nil {
			return metrics, err
		}
		metrics = append(metrics, metric)
	}
	return metrics, nil
}
