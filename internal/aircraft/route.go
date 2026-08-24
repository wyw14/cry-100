package aircraft

import (
	"fmt"
	"github.com/wyw14/cry-100/internal/model"
	"github.com/wyw14/cry-100/internal/perimeter"
)

func Plan(mission model.Mission, revision model.PerimeterRevision) (model.FlightPlan, error) {
	if mission.IncidentID != revision.IncidentID {
		return model.FlightPlan{}, fmt.Errorf("incident mismatch")
	}
	if len(revision.Boundary) == 0 {
		return model.FlightPlan{}, fmt.Errorf("boundary missing")
	}
	if !perimeter.Avoids(revision, mission.CorridorID) {
		return model.FlightPlan{}, fmt.Errorf("flight enters no-fly corridor")
	}
	return model.FlightPlan{MissionID: mission.ID, IncidentID: mission.IncidentID, RevisionID: revision.ID, RevisionTime: revision.ApprovedAt, Waypoints: append([]model.Point(nil), revision.Boundary...), AvoidedCorridors: append([]string(nil), revision.NoFlyZones...), Frozen: false}, nil
}
