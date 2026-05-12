package response

import (
	"encoding/json"
	"net/http"
)

// Resp is a standard response format
type Resp struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// JSON writes a JSON response with a status code
func JSON(w http.ResponseWriter, status int, success bool, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := Resp{
		Success: success,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(resp)
}
