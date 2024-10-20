package utils

import (
	"encoding/json"
	"net/http"
)

type httpError struct {
	Message string `json:"message"`
}

func NewHttpError(w http.ResponseWriter, err error) {
	httpError := httpError{
		Message: err.Error(),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)
	json.NewEncoder(w).Encode(httpError)
}
