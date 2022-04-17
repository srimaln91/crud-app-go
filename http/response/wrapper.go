package response

import (
	"encoding/json"
	"net/http"
)

type wrapper struct {
	Data       interface{} `json:"data,omitempty"`
	Status     string      `json:"status,omitempty"`
	Message    string      `json:"message,omitempty"`
	HttpStatus int         `json:"-"`
}

type responseStatus string

const (
	ACTION_SUCCESS responseStatus = "SUCCESS"
	ACTION_ERROR   responseStatus = "ERROR"
)

func New(status responseStatus, httpStatus int, options ...func(resp *wrapper)) *wrapper {
	resp := &wrapper{
		Status:     string(status),
		Data:       nil,
		HttpStatus: httpStatus,
	}

	// apply options
	for _, op := range options {
		op(resp)
	}

	return resp
}

func WithMessage(message string) func(resp *wrapper) {
	return func(resp *wrapper) {
		resp.Message = message
	}
}

func WithData(data interface{}) func(resp *wrapper) {
	return func(resp *wrapper) {
		resp.Data = data
	}
}

func (resp *wrapper) Write(rw http.ResponseWriter) {
	responseData, err := json.Marshal(resp)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(resp.HttpStatus)
	rw.Write(responseData)
}

func GenerateInternalServerError() *wrapper {
	return &wrapper{
		HttpStatus: http.StatusInternalServerError,
		Status:     string(ACTION_ERROR),
		Message:    "internal server error",
	}
}

func GenerateInvalidRequestError() *wrapper {
	return &wrapper{
		HttpStatus: http.StatusNotAcceptable,
		Status:     string(ACTION_ERROR),
		Message:    "invalid request",
	}
}
