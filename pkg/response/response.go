package response

import (
	"encoding/json"
	"net/http"
)

type JSONResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func JSON(w http.ResponseWriter, code int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(JSONResponse{
		Status:  http.StatusText(code),
		Message: message,
		Data:    data,
	})
}

func Error(w http.ResponseWriter, code int, message string) {
	JSON(w, code, message, nil)
}
