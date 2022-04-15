package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/srimaln91/crud-app-go/container"
	"github.com/srimaln91/crud-app-go/externals/repositories"
	"github.com/srimaln91/crud-app-go/http/handlers"
)

func newRouter(ctr *container.Container) (http.Handler, error) {
	eventRepository := repositories.NewEventRepository(ctr.DBAdapter, ctr.Logger)

	httpHandler := handlers.NewHttpHandler(eventRepository, ctr.Logger)

	r := mux.NewRouter()

	r.HandleFunc("/rest/events", httpHandler.AddEvent).Methods(http.MethodPost)
	r.HandleFunc("/rest/events", httpHandler.GetAllEvents).Methods(http.MethodGet)
	r.HandleFunc("/rest/events/{id}", httpHandler.GetEvent).Methods(http.MethodGet)
	r.HandleFunc("/rest/events/{id}", httpHandler.UpdateEvent).Methods(http.MethodPut)
	r.HandleFunc("/rest/events/{id}", httpHandler.DeleteEvent).Methods(http.MethodDelete)

	return r, nil
}
