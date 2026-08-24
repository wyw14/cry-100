package incident_test

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

func TestIncidentMergePreservesMissionAcknowledgement(t *testing.T) {
	now := time.Now().UTC()
	incidents := incident.NewService(nil)
	primary := incidents.Create("north-ridge", now)
	child := incidents.Create("pine-valley", now)
	resources := resource.NewRegistry()
	messages := radio.NewPublisher()
	planner := dispatch.NewPlanner(incidents, resources, messages, aircraft.NewService(resources), perimeter.NewService())
	result, err := planner.Plan(model.DispatchRequest{OperationID: model.NewOperationID(), IncidentID: child.ID, AssigneeID: "crew-a", Kind: "ground", RallyPoint: "pass"}, now)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = incidents.Merge(primary.ID, child.ID, now.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	planner.MigrateIncident(child.ID, primary.ID)
	if err = planner.Acknowledge(result.MissionID, "crew-a"); err != nil {
		t.Fatal(err)
	}
	message, ok := messages.Get(result.MissionID)
	if !ok || message.IncidentID != primary.ID {
		t.Fatalf("ack owner=%s expected=%s", message.IncidentID, primary.ID)
	}
	if message.AckedAt.IsZero() {
		t.Fatal("acknowledged mission is still pending")
	}
}
