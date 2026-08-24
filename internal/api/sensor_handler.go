package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/wyw14/cry-100/internal/model"
)

type sensorModeRequest struct {
	Mode        string  `json:"mode"`
	Revision    uint64  `json:"revision"`
	Temperature float64 `json:"temperature_offset"`
	SmokeScale  float64 `json:"smoke_scale"`
}

func (s *Server) setSensorMode(w http.ResponseWriter, r *http.Request) {
	var request sensorModeRequest
	if err := decode(r, &request); err != nil {
		fail(w, http.StatusBadRequest, err)
		return
	}
	id := model.SensorID(chi.URLParam(r, "sensorID"))
	calibration := &model.Calibration{Revision: request.Revision, Mode: request.Mode, Temperature: request.Temperature, SmokeScale: request.SmokeScale, PublishedAt: time.Now().UTC()}
	if err := s.Runtime.Sensors.SetMode(id, request.Mode, calibration); err != nil {
		fail(w, http.StatusConflict, err)
		return
	}
	value, _ := s.Runtime.Sensors.Profile(id)
	respond(w, http.StatusOK, value)
}
