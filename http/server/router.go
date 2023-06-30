package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/srimaln91/crud-app-go/container"
	"github.com/srimaln91/crud-app-go/externals/repositories"
	"github.com/srimaln91/crud-app-go/http/handlers"
)

func newRouter(ctr *container.Container) (http.Handler, error) {
	eventRepository := repositories.NewTaskRepository(ctr.DBAdapter, ctr.Logger)

	httpHandler := handlers.NewHttpHandler(eventRepository, ctr.Logger)

	r := mux.NewRouter()

	r.HandleFunc("/api/tasks", httpHandler.AddTask).Methods(http.MethodPost)
	r.HandleFunc("/api/tasks", httpHandler.GetAllTasks).Methods(http.MethodGet)
	r.HandleFunc("/api/tasks/{id}", httpHandler.GetTask).Methods(http.MethodGet)
	r.HandleFunc("/api/tasks/{id}", httpHandler.UpdateTask).Methods(http.MethodPut)
	r.HandleFunc("/api/tasks/{id}", httpHandler.DeleteTask).Methods(http.MethodDelete)

	return r, nil
}
