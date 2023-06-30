package bootstrap

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/srimaln91/crud-app-go/config"
	"github.com/srimaln91/crud-app-go/container"
	"github.com/srimaln91/crud-app-go/core/interfaces"
	"github.com/srimaln91/crud-app-go/http/server"
)

var logger interfaces.Logger
var httpServer *server.Server

func Start() {
	configFilePath := flag.String("config", "config.yaml", "config file path")
	flag.Parse()
	cfg, err := config.Parse(*configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// Init container
	ctr, err := container.Init(cfg)
	if err != nil {
		log.Fatal(err)
	}

	logger = ctr.Logger

	// Initialize and start HTTP server
	httpServer, err = server.Start(fmt.Sprintf("%s:%d", "0.0.0.0", cfg.HTTP.Port), ctr)
	if err != nil {
		logger.Fatal(context.TODO(), err.Error())
	}

	// Listen for term signals
	c := make(chan os.Signal, 1)

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	// Block until we receive our signal
	signal := <-c

	logger.Info(context.TODO(), fmt.Sprintf("received signal: %s", signal))
	logger.Info(context.TODO(), "shutting down the service...")

	// Destruct other respouces and stop the service
	Destruct(ctr)

	// Exit with non zero error code
	os.Exit(0)
}

func Destruct(ctr *container.Container) {

	// Shutdown Http server
	// create a deadline of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	httpServer.ShutDown(ctx)

	<-ctx.Done()

	// Close active/idle DB connections
	err := ctr.DBAdapter.Close()
	if err != nil {
		logger.Error(context.TODO(), err.Error())
	}

	logger.Info(context.TODO(), "service shutted down gracefully.")
}
