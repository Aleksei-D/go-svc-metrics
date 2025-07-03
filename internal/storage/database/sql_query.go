package database

const (
	CreateGaugeTableSQL = "CREATE TABLE IF NOT EXISTS gauge_table " +
		"(name_id TEXT PRIMARY KEY NOT NULL, value DOUBLE PRECISION NOT NULL)"
	CreateCounterTableSQL = "CREATE TABLE IF NOT EXISTS counter_table " +
		"(name_id TEXT PRIMARY KEY NOT NULL, delta INTEGER NOT NULL)"
)
