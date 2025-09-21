package datasource

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func NewDatabase(databaseURI string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURI)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := upMigrations(db); err != nil {
		return nil, err
	}
	return db, nil
}
