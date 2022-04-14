package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/srimaln91/crud-app-go/config"
	"github.com/srimaln91/crud-app-go/container"
	"github.com/srimaln91/crud-app-go/core/interfaces"
	"github.com/srimaln91/crud-app-go/externals/adapters"
	repositiories "github.com/srimaln91/crud-app-go/externals/repositories"
	"github.com/srimaln91/crud-app-go/http/server"
	"github.com/srimaln91/crud-app-go/logger"
)

var dbAdapter *sql.DB
var logAdapter interfaces.Logger

func Start() {
	cfg, err := config.Parse("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	logAdapter, err = logger.NewLogger(cfg.Logger.Level)
	if err != nil {
		log.Fatal(err)
	}

	dbConfig := adapters.PostgresConfig{
		Host:               cfg.Database.Host,
		Database:           cfg.Database.Name,
		User:               cfg.Database.User,
		Password:           cfg.Database.Password,
		Port:               cfg.Database.Port,
		PoolSize:           cfg.Database.PoolSize,
		MaxIdleConnections: cfg.Database.MaxIdleConnections,
		ConnMaxLifeTime:    time.Duration(cfg.Database.ConnMaxLifeTime) * time.Millisecond,
	}

	dbAdapter, err = adapters.NewPostgresDB(dbConfig)
	if err != nil {
		logAdapter.Fatal(context.Background(), err.Error())
	}

	eventRepository := repositiories.NewEventRepository(dbAdapter, logAdapter)
	// Build container
	ctr := &container.Container{
		DBAdapter:       dbAdapter,
		EventRepository: eventRepository,
		Logger:          logAdapter,
	}

	httpServer, err := server.Start(fmt.Sprintf("%s:%d", "0.0.0.0", cfg.HTTP.Port), ctr)
	if err != nil {
		logAdapter.Fatal(context.Background(), err.Error())
	}

	// Listen for term signals
	c := make(chan os.Signal, 1)

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	// Block until we receive our signal
	signal := <-c

	logAdapter.Info(context.Background(), fmt.Sprintf("received signal: %s", signal))

	// Shutdown Http server
	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	httpServer.ShutDown(ctx)

	<-ctx.Done()

	// Destruct other respouces and stop the service
	Destruct()

	// Exit with non zero error code
	os.Exit(0)
}

func Destruct() {
	err := dbAdapter.Close()
	if err != nil {
		logAdapter.Error(context.Background(), err.Error())
	}
}
