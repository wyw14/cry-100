package aircraft_test

import (
	"github.com/wyw14/cry-100/internal/aircraft"
	"github.com/wyw14/cry-100/internal/dispatch"
	"github.com/wyw14/cry-100/internal/incident"
	"github.com/wyw14/cry-100/internal/model"
	"github.com/wyw14/cry-100/internal/perimeter"
	"github.com/wyw14/cry-100/internal/radio"
	"github.com/wyw14/cry-100/internal/resource"
	"testing"
	"time"
)

func TestNewFlightUsesCurrentPerimeterRevision(t *testing.T) {
	now := time.Now().UTC()
	incidents := incident.NewService(nil)
	fire := incidents.Create("north-ridge", now)
	perimeters := perimeter.NewService()
	perimeters.Publish(fire.ID, []model.Point{{Lat: 1, Lon: 1}, {Lat: 2, Lon: 1}, {Lat: 2, Lon: 2}}, nil, 10, now)
	current := perimeters.Publish(fire.ID, []model.Point{{Lat: 2, Lon: 2}, {Lat: 3, Lon: 2}, {Lat: 3, Lon: 3}}, []string{"smoke-east"}, 90, now.Add(time.Minute))
	resources := resource.NewRegistry()
	air := aircraft.NewService(resources)
	planner := dispatch.NewPlanner(incidents, resources, radio.NewPublisher(), air, perimeters)
	result, err := planner.Plan(model.DispatchRequest{OperationID: model.NewOperationID(), IncidentID: fire.ID, AssigneeID: "heli-1", Kind: "air", RallyPoint: "ridge"}, now.Add(2*time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	var plan model.FlightPlan
	for _, value := range air.Plans() {
		if value.MissionID == result.MissionID {
			plan = value
		}
	}
	if plan.MissionID == "" {
		t.Fatal("flight plan missing")
	}
	if plan.RevisionID != current.ID || !plan.RevisionTime.Equal(current.ApprovedAt) {
		t.Fatalf("flight revision=%s expected=%s", plan.RevisionID, current.ID)
	}
}
