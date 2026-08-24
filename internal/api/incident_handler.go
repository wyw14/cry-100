package api

import (
	"encoding/base64"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/wyw14/cry-100/internal/model"
)

type incidentRequest struct {
	SectorID string `json:"sector_id"`
}

type mergeRequest struct {
	ChildID model.IncidentID `json:"child_id"`
}

type containmentRequest struct {
	ReportID string `json:"report_id"`
	Summary  string `json:"summary"`
	Image    string `json:"image_base64"`
	Name     string `json:"image_name"`
}

func (s *Server) incidents(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"incidents": s.Runtime.Incidents.List()})
}
func (s *Server) createIncident(w http.ResponseWriter, r *http.Request) {
	var request incidentRequest
	if err := decode(r, &request); err != nil {
		writeJSON(w, 400, map[string]string{"error": err.Error()})
		return
	}
	if request.SectorID == "" {
		request.SectorID = "north-ridge"
	}
	incident := s.Runtime.Incidents.Create(request.SectorID, time.Now().UTC())
	writeJSON(w, http.StatusCreated, incident)
}
func (s *Server) closeIncident(w http.ResponseWriter, r *http.Request) {
	id := model.IncidentID(chi.URLParam(r, "incidentID"))
	if err := s.Runtime.Incidents.Close(id, time.Now().UTC()); err != nil {
		fail(w, http.StatusConflict, err)
		return
	}
	value, _ := s.Runtime.Incidents.Get(id)
	respond(w, http.StatusOK, value)
}

func (s *Server) mergeIncident(w http.ResponseWriter, r *http.Request) {
	var request mergeRequest
	if err := decode(r, &request); err != nil {
		fail(w, http.StatusBadRequest, err)
		return
	}
	result, err := s.Runtime.MergeIncidents(model.IncidentID(chi.URLParam(r, "incidentID")), request.ChildID, time.Now().UTC())
	if err != nil {
		fail(w, http.StatusConflict, err)
		return
	}
	respond(w, http.StatusOK, result)
}

func (s *Server) containIncident(w http.ResponseWriter, r *http.Request) {
	var request containmentRequest
	if err := decode(r, &request); err != nil {
		fail(w, http.StatusBadRequest, err)
		return
	}
	image, err := base64.StdEncoding.DecodeString(request.Image)
	if err != nil {
		fail(w, http.StatusBadRequest, err)
		return
	}
	proof := model.FieldProof{IncidentID: model.IncidentID(chi.URLParam(r, "incidentID")), ReportID: request.ReportID, Summary: request.Summary, ImageName: request.Name, SubmittedAt: time.Now().UTC()}
	if err := s.Runtime.ContainIncident(proof, image, time.Now().UTC()); err != nil {
		fail(w, http.StatusConflict, err)
		return
	}
	value, _ := s.Runtime.Incidents.Get(proof.IncidentID)
	respond(w, http.StatusOK, value)
}
