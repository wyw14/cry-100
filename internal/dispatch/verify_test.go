package dispatch_test

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

func TestCancelledDropReleasesOwnedWaterSlot(t *testing.T) {
	now := time.Now().UTC()
	incidents := incident.NewService(nil)
	fire := incidents.Create("north-ridge", now)
	resources := resource.NewRegistry()
	resources.AddCorridor("water", "Water", 1)
	planner := dispatch.NewPlanner(incidents, resources, radio.NewPublisher(), aircraft.NewService(resources), perimeter.NewService())
	result, err := planner.Plan(model.DispatchRequest{OperationID: model.NewOperationID(), IncidentID: fire.ID, AssigneeID: "heli-a", Kind: "ground", RallyPoint: "ridge", CorridorID: "water"}, now)
	if err != nil {
		t.Fatal(err)
	}
	if resources.ActiveOwner("water") != 1 {
		t.Fatal("lease not reserved")
	}
	if err = planner.Cancel(result.MissionID, now.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	if resources.ActiveOwner("water") != 0 {
		t.Fatalf("water owners=%d", resources.ActiveOwner("water"))
	}
	mission, _ := planner.Get(result.MissionID)
	if mission.Stage != model.MissionCancelled {
		t.Fatalf("stage=%s", mission.Stage)
	}
}
