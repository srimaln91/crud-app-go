package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/srimaln91/crud-app-go/core/entities"
	"github.com/srimaln91/crud-app-go/core/interfaces"
	"github.com/srimaln91/crud-app-go/http/request"
	"github.com/srimaln91/crud-app-go/http/response"
)

type handler struct {
	logger          interfaces.Logger
	eventRepository interfaces.EventRepository
}

const URL_PARAM_ID = "id"
const MESSAGE_RECORD_NOT_EXIST = "record does not exist"

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
		response.GenerateInvalidRequestError().Write(rw)
		return
	}

	defer r.Body.Close()

	event.ID = uuid.New().String()

	err = h.eventRepository.Add(r.Context(), event)
	if err != nil {
		response.GenerateInternalServerError().Write(rw)
		return
	}

	response.New(
		response.ACTION_SUCCESS,
		http.StatusCreated,
		response.WithData(event),
	).Write(rw)
}

func (h *handler) GetEvent(rw http.ResponseWriter, r *http.Request) {

	eventID, err := getURLParam(r, URL_PARAM_ID)
	if err != nil {
		response.GenerateInvalidRequestError().Write(rw)
		return
	}

	entry, err := h.eventRepository.Get(r.Context(), eventID)
	if err != nil {
		response.GenerateInternalServerError().Write(rw)
		return
	}

	if entry == nil {
		response.New(
			response.ACTION_ERROR,
			http.StatusNotFound,
			response.WithMessage(MESSAGE_RECORD_NOT_EXIST),
		).Write(rw)
		return
	}

	response.New(
		response.ACTION_SUCCESS,
		http.StatusOK,
		response.WithData(entry),
	).Write(rw)
}

func (h *handler) UpdateEvent(rw http.ResponseWriter, r *http.Request) {

	eventID, err := getURLParam(r, URL_PARAM_ID)
	if err != nil {
		response.GenerateInvalidRequestError().Write(rw)
		return
	}

	var event entities.Event
	err = json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		response.GenerateInvalidRequestError().Write(rw)
		return
	}

	event.ID = eventID
	defer r.Body.Close()

	recordExist, err := h.eventRepository.Update(r.Context(), event.ID, event)
	if err != nil {
		response.GenerateInternalServerError().Write(rw)
		return
	}

	if !recordExist {
		r := response.New(
			response.ACTION_ERROR,
			http.StatusNotFound,
			response.WithMessage(MESSAGE_RECORD_NOT_EXIST),
		)
		r.Write(rw)
		return
	}

	response.New(
		response.ACTION_SUCCESS,
		http.StatusOK,
		response.WithData(event),
	).Write(rw)
}

func (h *handler) GetAllEvents(rw http.ResponseWriter, r *http.Request) {
	entries, err := h.eventRepository.GetAll(r.Context())
	if err != nil {
		response.GenerateInternalServerError().Write(rw)
		return
	}

	response.New(
		response.ACTION_SUCCESS,
		http.StatusOK,
		response.WithData(entries),
	).Write(rw)
}

func (h *handler) DeleteEvent(rw http.ResponseWriter, r *http.Request) {

	eventID, err := getURLParam(r, URL_PARAM_ID)
	if err != nil {
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	resourceExist, err := h.eventRepository.Remove(r.Context(), eventID)
	if err != nil {
		response.GenerateInternalServerError().Write(rw)
		return
	}

	if !resourceExist {
		response.New(
			response.ACTION_ERROR,
			http.StatusNotFound,
			response.WithMessage(MESSAGE_RECORD_NOT_EXIST),
		).Write(rw)
		return
	}

	response.New(
		response.ACTION_SUCCESS,
		http.StatusOK,
	).Write(rw)
}

func (h *handler) AddEventBatch(rw http.ResponseWriter, r *http.Request) {
	var batchRequest request.EventBatch
	err := json.NewDecoder(r.Body).Decode(&batchRequest)
	if err != nil {
		response.GenerateInvalidRequestError().Write(rw)
		return
	}

	defer r.Body.Close()

	events := make([]entities.Event, 0)

	for _, record := range batchRequest.Records {
		for _, event := range record.Event {
			events = append(events, entities.Event{
				ID:          uuid.New().String(),
				TransId:     record.TransID,
				TransTms:    record.TransTms,
				RcNum:       record.RcNum,
				ClientId:    record.ClientID,
				EventCnt:    event.EventCnt,
				LocationCd:  event.LocationCd,
				LocationId1: event.LocationID1,
				LocationId2: event.LocationID2,
				AddrNbr:     event.AddrNbr,
			})
		}
	}

	err = h.eventRepository.InsertBatch(r.Context(), events)
	if err != nil {
		response.GenerateInternalServerError().Write(rw)
		return
	}

	response.New(
		response.ACTION_SUCCESS,
		http.StatusCreated,
	).Write(rw)
}

func getURLParam(r *http.Request, parameter string) (string, error) {
	urlParams := mux.Vars(r)
	eventID, ok := urlParams[parameter]
	if !ok {
		return "", errURLParamDoesNotExist
	}

	return eventID, nil
}
