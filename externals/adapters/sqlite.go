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

	db.Exec(`CREATE TABLE IF NOT EXISTS events (
		event_id uuid NOT NULL,
		addr_nbr varchar(255) NULL,
		client_id varchar(255) NULL,
		event_cnt int4 NOT NULL,
		location_cd varchar(255) NULL,
		location_id1 varchar(255) NULL,
		location_id2 varchar(255) NULL,
		rc_num varchar(255) NULL,
		trans_id varchar(255) NULL,
		trans_tms varchar(255) NULL,
		CONSTRAINT events_pkey PRIMARY KEY (event_id)
	);`)

	return db, nil
}
