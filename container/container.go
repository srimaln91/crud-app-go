package container

import (
	"database/sql"
	"time"

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
	dbConfig := adapters.PostgresConfig{
		Database:           cfg.Database,
		Host:               cfg.Host,
		User:               cfg.User,
		Password:           cfg.Password,
		Port:               cfg.Port,
		PoolSize:           cfg.PoolSize,
		MaxIdleConnections: cfg.MaxIdleConnections,
		ConnMaxLifeTime:    time.Second * time.Duration(cfg.ConnMaxLifeTime),
	}

	dbAdapter, err := adapters.NewPostgresDB(dbConfig)
	if err != nil {
		return nil, err
	}

	return dbAdapter, nil
}
