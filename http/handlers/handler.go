package handlers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/srimaln91/crud-app-go/core/entities"
	"github.com/srimaln91/crud-app-go/core/interfaces"
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
	//read body

	//validate struct
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

	var responseWriter bytes.Buffer
	err = json.NewEncoder(bufio.NewWriter(&responseWriter)).Encode(event)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	rw.Write(responseWriter.Bytes())
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

	bytes, err := json.Marshal(entry)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(bytes)
}

func (h *handler) GetAllEvents(rw http.ResponseWriter, r *http.Request) {
	entries, err := h.eventRepository.GetAll()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(entries)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(bytes)
}
