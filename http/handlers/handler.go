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
	logger         interfaces.Logger
	taskRepository interfaces.TaskRepository
}

const URL_PARAM_ID = "id"
const MESSAGE_RECORD_NOT_EXIST = "record does not exist"

var errURLParamDoesNotExist = errors.New("url parameter does not exist")

func NewHttpHandler(taskRepository interfaces.TaskRepository, logger interfaces.Logger) *handler {
	return &handler{
		taskRepository: taskRepository,
		logger:         logger,
	}
}

func (h *handler) AddTask(rw http.ResponseWriter, r *http.Request) {
	var task entities.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		h.logger.Error(r.Context(), err.Error())
		response.GenerateInvalidRequestError().Write(rw)
		return
	}

	defer r.Body.Close()

	task.ID = uuid.New().String()

	err = h.taskRepository.Add(r.Context(), task)
	if err != nil {
		h.logger.Error(r.Context(), err.Error())
		response.GenerateInternalServerError().Write(rw)
		return
	}

	response.New(
		response.ACTION_SUCCESS,
		http.StatusCreated,
		response.WithData(task),
	).Write(rw)
}

func (h *handler) GetTask(rw http.ResponseWriter, r *http.Request) {

	taskID, err := getURLParam(r, URL_PARAM_ID)
	if err != nil {
		response.GenerateInvalidRequestError().Write(rw)
		return
	}

	entry, err := h.taskRepository.Get(r.Context(), taskID)
	if err != nil {
		h.logger.Error(r.Context(), err.Error())
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

func (h *handler) UpdateTask(rw http.ResponseWriter, r *http.Request) {

	taskID, err := getURLParam(r, URL_PARAM_ID)
	if err != nil {
		response.GenerateInvalidRequestError().Write(rw)
		return
	}

	var task entities.Task
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		response.GenerateInvalidRequestError().Write(rw)
		return
	}

	task.ID = taskID
	defer r.Body.Close()

	recordExist, err := h.taskRepository.Update(r.Context(), task.ID, task)
	if err != nil {
		h.logger.Error(r.Context(), err.Error())
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
		response.WithData(task),
	).Write(rw)
}

func (h *handler) GetAllTasks(rw http.ResponseWriter, r *http.Request) {
	entries, err := h.taskRepository.GetAll(r.Context())
	if err != nil {
		h.logger.Error(r.Context(), err.Error())
		response.GenerateInternalServerError().Write(rw)
		return
	}

	response.New(
		response.ACTION_SUCCESS,
		http.StatusOK,
		response.WithData(entries),
	).Write(rw)
}

func (h *handler) DeleteTask(rw http.ResponseWriter, r *http.Request) {

	taskID, err := getURLParam(r, URL_PARAM_ID)
	if err != nil {
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	resourceExist, err := h.taskRepository.Remove(r.Context(), taskID)
	if err != nil {
		h.logger.Error(r.Context(), err.Error())
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

func getURLParam(r *http.Request, parameter string) (string, error) {
	urlParams := mux.Vars(r)
	eventID, ok := urlParams[parameter]
	if !ok {
		return "", errURLParamDoesNotExist
	}

	return eventID, nil
}
