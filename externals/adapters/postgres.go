package adapters

import (
	"database/sql"

	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type PostgresConfig struct {
	Database           string
	User               string
	Password           string
	Host               string
	Port               int
	PoolSize           int
	MaxIdleConnections int
	ConnMaxLifeTime    time.Duration
}

// Adapter is used to communicate with a Postgres database.
// New creates a new Postgres adapter instance.
func NewPostgresDB(config PostgresConfig) (*sql.DB, error) {

	connString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
		config.User, config.Password, config.Database, config.Host, config.Port)

	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	// Pool configurations
	db.SetMaxOpenConns(config.PoolSize)
	db.SetMaxIdleConns(config.MaxIdleConnections)
	db.SetConnMaxLifetime(config.ConnMaxLifeTime)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
