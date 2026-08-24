package api

import (
	"net/http"
	"time"

	"github.com/wyw14/cry-100/internal/model"
)

type perimeterRequest struct {
	IncidentID  model.IncidentID `json:"incident_id"`
	Boundary    []model.Point    `json:"boundary"`
	NoFlyZones  []string         `json:"no_fly_zones"`
	WindBearing int              `json:"wind_bearing"`
}

func (s *Server) perimeters(w http.ResponseWriter, r *http.Request) {
	respond(w, http.StatusOK, map[string]any{"perimeters": s.Runtime.Perimeters.All()})
}

func (s *Server) publishPerimeter(w http.ResponseWriter, r *http.Request) {
	var request perimeterRequest
	if err := decode(r, &request); err != nil {
		fail(w, http.StatusBadRequest, err)
		return
	}
	value := s.Runtime.Perimeters.Publish(request.IncidentID, request.Boundary, request.NoFlyZones, request.WindBearing, time.Now().UTC())
	respond(w, http.StatusCreated, value)
}
