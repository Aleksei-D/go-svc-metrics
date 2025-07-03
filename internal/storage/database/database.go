package database

import (
	"context"
	"database/sql"
	"fmt"
	"go-svc-metrics/models"
)

type Storage struct {
	*sql.DB
}

func NewDatabaseStorage(db *sql.DB) (*Storage, error) {
	_, err := db.Exec(CreateGaugeTableSQL)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(CreateCounterTableSQL)
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

func (d *Storage) UpdateMetric(metric models.Metrics) (models.Metrics, error) {
	switch metric.MType {
	case models.Gauge:
		return d.updateGaugeMetric(metric)
	case models.Counter:
		return d.updateCounterMetric(metric)
	}
	return models.Metrics{}, fmt.Errorf("metric type not supported")
}

func (d *Storage) updateCounterMetric(metric models.Metrics) (models.Metrics, error) {
	var delta int64
	tx, err := d.DB.Begin()
	if err != nil {
		return metric, err
	}
	_, err = tx.Exec("INSERT INTO counter_table (name_id, delta) VALUES($1, $2)", metric.ID, metric.Delta)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return metric, err
		}
	} else {
		err = tx.Commit()
		if err != nil {
			return metric, err
		}
		return metric, nil
	}
	row := d.DB.QueryRow(
		"UPDATE counter_table SET delta = delta + $1 WHERE name_id = $2 RETURNING delta",
		metric.Delta,
		metric.ID,
	)
	err = row.Scan(&delta)
	if err != nil {
		return metric, nil
	}
	metric.Delta = &delta
	return metric, nil
}

func (d *Storage) updateGaugeMetric(metric models.Metrics) (models.Metrics, error) {
	tx, err := d.DB.Begin()
	if err != nil {
		return metric, err
	}
	_, err = tx.Exec("INSERT INTO gauge_table (name_id, value) VALUES($1, $2)", metric.ID, metric.Value)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return metric, err
		}
	} else {
		err = tx.Commit()
		if err != nil {
			return metric, err
		}
		return metric, nil
	}

	_, err = tx.Exec("UPDATE gauge_table SET value = $1 WHERE name_id = $2", metric.Value, metric.ID)
	if err != nil {
		return metric, tx.Rollback()
	}
	return metric, tx.Commit()
}

func (d *Storage) GetValue(metric models.Metrics) (models.Metrics, bool) {
	switch metric.MType {
	case models.Gauge:
		return d.getGaugeMetric(metric)
	case models.Counter:
		return d.getCounterMetric(metric)
	}
	return models.Metrics{}, false
}

func (d *Storage) getGaugeMetric(metric models.Metrics) (models.Metrics, bool) {
	row := d.DB.QueryRow("SELECT value FROM gauge_table WHERE name_id = $1", metric.ID)
	err := row.Scan(&metric.Value)
	if err != nil {
		return metric, false
	}
	return metric, true
}

func (d *Storage) getCounterMetric(metric models.Metrics) (models.Metrics, bool) {
	row := d.DB.QueryRow("SELECT delta FROM counter_table WHERE name_id = $1", metric.ID)
	err := row.Scan(&metric.Delta)
	if err != nil {
		return metric, false
	}
	return metric, true
}

func (d *Storage) GetAllMetrics() (map[string]models.Metrics, error) {
	metrics := make(map[string]models.Metrics)
	gaugeCounter, err := d.getAllGaugeMetrics()
	if err != nil {
		return metrics, err
	}
	for _, metric := range gaugeCounter {
		metrics[metric.ID] = metric
	}

	counterMetrics, err := d.getAllCounterMetrics()
	if err != nil {
		return metrics, err
	}
	for _, metric := range counterMetrics {
		metrics[metric.ID] = metric
	}

	return metrics, nil
}

func (d *Storage) getAllGaugeMetrics() ([]models.Metrics, error) {
	metrics := make([]models.Metrics, 0)
	rowsGauge, err := d.DB.Query("SELECT name_id, value FROM gauge_table")
	if err != nil {
		return metrics, err
	}
	if err := rowsGauge.Err(); err != nil {
		return metrics, err
	}

	for rowsGauge.Next() {
		var metric models.Metrics
		err := rowsGauge.Scan(&metric.ID, &metric.Delta)
		if err != nil {
			return metrics, err
		}
		metrics = append(metrics, metric)
	}
	return metrics, err
}

func (d *Storage) getAllCounterMetrics() ([]models.Metrics, error) {
	metrics := make([]models.Metrics, 0)
	rowsCounter, err := d.DB.Query("SELECT name_id, delta FROM counter_table")
	if err != nil {
		return metrics, err
	}
	if err := rowsCounter.Err(); err != nil {
		return metrics, err
	}

	for rowsCounter.Next() {
		var metric models.Metrics
		err := rowsCounter.Scan(&metric.ID, &metric.Delta)
		if err != nil {
			return metrics, err
		}
		metrics = append(metrics, metric)
	}
	return metrics, err
}

func (d *Storage) Close() error { return d.DB.Close() }

func (d *Storage) DumpMetricsByInterval(_ context.Context) error {
	return nil
}
