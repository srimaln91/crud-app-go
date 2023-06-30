package adapters

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteConfig struct {
	DatabasePath string
}

// Adapter is used to communicate with a Postgres database.
// New creates a new Postgres adapter instance.
func NewSQLiteDB(config SQLiteConfig) (*sql.DB, error) {

	db, err := sql.Open("sqlite3", config.DatabasePath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
