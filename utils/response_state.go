package utilis

import (
	"encoding/json"
	"net/http"
)

type CommonResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ResponseState struct {
	StatusCode int
	Message    string
}

func (rs ResponseState) WriteToResponse(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rs.StatusCode)
	return json.NewEncoder(w).Encode(CommonResponse{
		Message: rs.Message,
		Data:    data,
	})
}
