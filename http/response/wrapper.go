package response

import (
	"encoding/json"
	"net/http"
)

type wrapper struct {
	Data    interface{} `json:"data"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
}

const (
	HTTP_SUCCESS string = "OK"
	HTTP_ERROR   string = "ERROR"
)

func WriteSuccessResponse(rw http.ResponseWriter, data interface{}, statusCode int) {

	wrappedResponse := wrapper{
		Data:   data,
		Status: HTTP_SUCCESS,
	}
	responseData, err := json.Marshal(wrappedResponse)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	rw.Write(responseData)
}

func WriteErrorResponse(rw http.ResponseWriter, err error, statusCode int) {

	wrappedResponse := wrapper{
		Data:    nil,
		Status:  HTTP_ERROR,
		Message: err.Error(),
	}

	responseData, err := json.Marshal(wrappedResponse)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	rw.Write(responseData)
}
