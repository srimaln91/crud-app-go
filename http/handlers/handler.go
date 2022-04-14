package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/srimaln91/crud-app-go/core/entities"
	"github.com/srimaln91/crud-app-go/core/interfaces"
	"github.com/srimaln91/crud-app-go/http/response"
)

type handler struct {
	logger          interfaces.Logger
	eventRepository interfaces.EventRepository
}

func NewHttpHandler(eventRepository interfaces.EventRepository, logger interfaces.Logger) *handler {
	return &handler{
		eventRepository: eventRepository,
		logger:          logger,
	}
}

func (h *handler) AddEvent(rw http.ResponseWriter, r *http.Request) {
	var event entities.Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	defer r.Body.Close()

	event.ID = uuid.New().String()

	err = h.eventRepository.Add(event)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteSuccessResponse(rw, event, http.StatusCreated)
}

func (h *handler) GetEvent(rw http.ResponseWriter, r *http.Request) {

	urlParams := mux.Vars(r)
	eventID, ok := urlParams["id"]
	if !ok {
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	entry, err := h.eventRepository.Get(eventID)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteSuccessResponse(rw, entry, http.StatusCreated)
}

func (h *handler) GetAllEvents(rw http.ResponseWriter, r *http.Request) {
	entries, err := h.eventRepository.GetAll()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteSuccessResponse(rw, entries, http.StatusOK)
}
