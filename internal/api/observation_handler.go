package api

import (
	"github.com/wyw14/cry-100/internal/model"
	"net/http"
	"time"
)

type observationRequest struct {
	SensorID    string    `json:"sensor_id"`
	SectorID    string    `json:"sector_id"`
	EventTime   time.Time `json:"event_time"`
	Temperature float64   `json:"temperature_c"`
	Smoke       float64   `json:"smoke_index"`
	Hotspot     bool      `json:"image_hotspot"`
}

func (s *Server) observations(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"observations": s.Runtime.Streams.All()})
}
func (s *Server) createObservation(w http.ResponseWriter, r *http.Request) {
	var request observationRequest
	if err := decode(r, &request); err != nil {
		writeJSON(w, 400, map[string]string{"error": err.Error()})
		return
	}
	id := model.SensorID(request.SensorID)
	observation, err := s.Runtime.Sensors.Ingest(id, request.SectorID, request.EventTime, request.Temperature, request.Smoke, request.Hotspot)
	if err != nil {
		writeJSON(w, 409, map[string]string{"error": err.Error()})
		return
	}
	s.Runtime.Streams.Add(observation)
	candidate, ok := s.Runtime.Detector.Analyze(observation)
	if ok {
		_, _, _ = s.Runtime.Incidents.Promote(candidate, observation, time.Now().UTC())
	}
	writeJSON(w, http.StatusAccepted, observation)
}
