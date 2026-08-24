package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wyw14/cry-100/internal/model"
)

type alertRequest struct {
	IncidentID model.IncidentID `json:"incident_id"`
	Area       string           `json:"area"`
	Message    string           `json:"message"`
}

func (s *Server) alerts(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"alerts": s.Runtime.Alerts.List()})
}

func (s *Server) publishAlert(w http.ResponseWriter, r *http.Request) {
	var request alertRequest
	if err := decode(r, &request); err != nil {
		fail(w, http.StatusBadRequest, err)
		return
	}
	respond(w, http.StatusCreated, s.Runtime.Alerts.Publish(request.IncidentID, request.Area, request.Message))
}

func (s *Server) withdrawAlert(w http.ResponseWriter, r *http.Request) {
	id := model.IncidentID(chi.URLParam(r, "incidentID"))
	if err := s.Runtime.Alerts.Withdraw(id); err != nil {
		fail(w, http.StatusConflict, err)
		return
	}
	value, _ := s.Runtime.Alerts.Get(id)
	respond(w, http.StatusOK, value)
}
