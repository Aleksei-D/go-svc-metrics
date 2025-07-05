package database

const CreateTableSQL = "CREATE TABLE IF NOT EXISTS metric_table " +
	"(name_id TEXT PRIMARY KEY NOT NULL, value DOUBLE PRECISION DEFAULT 0, delta INTEGER DEFAULT 0)"
