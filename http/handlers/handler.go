package handlers

import (
	"encoding/json"
	"errors"
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

const URL_PARAM_ID = "id"

var errURLParamDoesNotExist = errors.New("url parameter does not exist")

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

	err = h.eventRepository.Add(r.Context(), event)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteSuccessResponse(rw, event, http.StatusCreated)
}

func (h *handler) GetEvent(rw http.ResponseWriter, r *http.Request) {

	eventID, err := getURLParam(r, URL_PARAM_ID)
	if err != nil {
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	entry, err := h.eventRepository.Get(r.Context(), eventID)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteSuccessResponse(rw, entry, http.StatusCreated)
}

func (h *handler) UpdateEvent(rw http.ResponseWriter, r *http.Request) {

	eventID, err := getURLParam(r, URL_PARAM_ID)
	if err != nil {
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	var event entities.Event
	err = json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	event.ID = eventID
	defer r.Body.Close()

	err = h.eventRepository.Update(r.Context(), event.ID, event)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteSuccessResponse(rw, event, http.StatusOK)
}

func (h *handler) GetAllEvents(rw http.ResponseWriter, r *http.Request) {
	entries, err := h.eventRepository.GetAll(r.Context())
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteSuccessResponse(rw, entries, http.StatusOK)
}

func (h *handler) DeleteEvent(rw http.ResponseWriter, r *http.Request) {

	eventID, err := getURLParam(r, URL_PARAM_ID)
	if err != nil {
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	err = h.eventRepository.Remove(r.Context(), eventID)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteSuccessResponse(rw, nil, http.StatusOK)
}

func getURLParam(r *http.Request, parameter string) (string, error) {
	urlParams := mux.Vars(r)
	eventID, ok := urlParams[parameter]
	if !ok {
		return "", errURLParamDoesNotExist
	}

	return eventID, nil
}
