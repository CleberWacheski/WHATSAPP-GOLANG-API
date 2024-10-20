package utils

import (
	"encoding/json"
	"net/http"
)

func HttpJsonResponse(w http.ResponseWriter, jsonData interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jsonData)
}
