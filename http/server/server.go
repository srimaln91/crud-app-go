package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/srimaln91/crud-app-go/container"
	"github.com/srimaln91/crud-app-go/core/interfaces"
)

type server struct {
	httpServer http.Server
	logger     interfaces.Logger
}

// Init initializes the server
func Start(address string, ctr *container.Container) (server, error) {

	// initialize the router
	handler, _ := newRouter(ctr)

	httpServer := &http.Server{
		Addr: address,

		// good practice to set timeouts to avoid Slowloris attacks
		WriteTimeout: time.Second * 60,
		ReadTimeout:  time.Second * 60,
		IdleTimeout:  time.Second * 60,

		// pass our instance of gorilla/mux in
		Handler: handler,
	}

	// run our server in a goroutine so that it doesn't block
	go func() {

		err := httpServer.ListenAndServe()
		if err != nil {
			ctr.Logger.Fatal(context.TODO(), err.Error())
		}
	}()

	// srv.server = server
	ctr.Logger.Info(context.Background(), fmt.Sprintf("HTTP server listening on %s", address), "functional_path", "http.server.Init")

	return server{
		httpServer: *httpServer,
	}, nil

}

// ShutDown releases all http connections gracefully and shut down the server
func (srv *server) ShutDown(ctx context.Context) {

	// srv.logger.Warn("http.server.ShutDown", "Stopping HTTP Server")
	srv.httpServer.SetKeepAlivesEnabled(false)

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	err := srv.httpServer.Shutdown(ctx)
	if err != nil {
		srv.logger.Fatal(ctx, "Unable to stop HTTP server", "functional_path", "http.server.ShutDown")

	}

}
