package api

import (
	"net/http"
	"time"
)

type envelope struct {
	Data      any       `json:"data,omitempty"`
	Error     string    `json:"error,omitempty"`
	RequestAt time.Time `json:"request_at"`
}

func respond(w http.ResponseWriter, status int, data any) {
	writeJSON(w, status, envelope{Data: data, RequestAt: time.Now().UTC()})
}
func fail(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, envelope{Error: err.Error(), RequestAt: time.Now().UTC()})
}
