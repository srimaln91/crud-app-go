package container

import (
	"database/sql"

	"github.com/srimaln91/crud-app-go/config"
	"github.com/srimaln91/crud-app-go/core/interfaces"
	"github.com/srimaln91/crud-app-go/externals/adapters"
	"github.com/srimaln91/crud-app-go/externals/repositories"
	"github.com/srimaln91/crud-app-go/log"
)

type Container struct {
	DBAdapter      *sql.DB
	TaskRepository interfaces.TaskRepository
	Logger         interfaces.Logger
}

func Init(cfg *config.AppConfig) (*Container, error) {
	// Resolve logger
	logAdapter, err := resolveLogger(cfg.Logger.Level)
	if err != nil {
		return nil, err
	}

	// resolve DB Adapter
	db, err := resolveDatabase(cfg.Database)
	if err != nil {
		return nil, err
	}

	// Resolve repositories and return container
	return &Container{
		DBAdapter:      db,
		Logger:         logAdapter,
		TaskRepository: repositories.NewTaskRepository(db, logAdapter),
	}, nil
}

func resolveLogger(level log.Level) (interfaces.Logger, error) {
	logAdapter, err := log.NewLogger(level)
	if err != nil {
		return nil, err
	}

	return logAdapter, nil
}

func resolveDatabase(cfg config.DBConfig) (*sql.DB, error) {
	dbConfig := adapters.SQLiteConfig{
		DatabasePath: cfg.Path,
	}

	dbAdapter, err := adapters.NewSQLiteDB(dbConfig)
	if err != nil {
		return nil, err
	}

	return dbAdapter, nil
}
