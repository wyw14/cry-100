package incident_test

import (
	"github.com/wyw14/cry-100/internal/incident"
	"github.com/wyw14/cry-100/internal/model"
	"testing"
	"time"
)

func TestHistoricalObservationCannotReopenClosedIncident(t *testing.T) {
	now := time.Now().UTC()
	service := incident.NewService(nil)
	current := service.Create("north-ridge", now.Add(-time.Hour))
	if err := service.Close(current.ID, now); err != nil {
		t.Fatal(err)
	}
	observation := model.Observation{ID: model.NewObservationID(), SectorID: "north-ridge", EventTime: now.Add(-30 * time.Minute), ReceivedAt: now.Add(time.Minute), CollectionEpoch: 1, IncidentGeneration: current.Generation, TemperatureC: 80}
	candidate := model.Candidate{ObservationID: observation.ID, SectorID: observation.SectorID, EventTime: observation.EventTime, CollectionEpoch: 1, IncidentGeneration: current.Generation, Confidence: .9}
	result, active, err := service.Promote(candidate, observation, now.Add(time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if active || result.Stage != model.StageClosed {
		t.Fatalf("historical observation changed terminal incident: active=%v stage=%s", active, result.Stage)
	}
	if result.ArchivedCount != 1 {
		t.Fatalf("archived count=%d", result.ArchivedCount)
	}
}
