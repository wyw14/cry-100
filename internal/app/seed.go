package app

import (
	"github.com/wyw14/cry-100/internal/crew"
	"github.com/wyw14/cry-100/internal/model"
	"time"
)

func (r *Runtime) Seed() {
	sensorID := model.SensorID("00000000-0000-4000-8000-000000000100")
	r.Sensors.Register(sensorID, "visible", model.Calibration{Revision: 1, Mode: "visible", SmokeScale: 1, PublishedAt: time.Now().UTC()})
	r.Crews.Add(structCrew("alpha", "Alpha Crew", "north-ridge"))
}
func structCrew(id, name, sector string) crew.Crew {
	return crew.Crew{ID: id, Name: name, SectorID: sector, Available: true}
}
