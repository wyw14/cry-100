package dispatch

import (
	"testing"
	"time"

	"github.com/wyw14/cry-100/internal/aircraft"
	"github.com/wyw14/cry-100/internal/incident"
	"github.com/wyw14/cry-100/internal/model"
	"github.com/wyw14/cry-100/internal/perimeter"
	"github.com/wyw14/cry-100/internal/radio"
	"github.com/wyw14/cry-100/internal/resource"
)

func newPlanner(t *testing.T) (*Planner, *incident.Service) {
	t.Helper()
	resources := resource.NewRegistry()
	resources.AddCorridor("ridge-water", "Ridge Water", 1)
	radioPublisher := radio.NewPublisher()
	perimeters := perimeter.NewService()
	incidents := incident.NewService(nil)
	aircraftService := aircraft.NewService(resources)
	return NewPlanner(incidents, resources, radioPublisher, aircraftService, perimeters), incidents
}

func TestCancelReleasesCorridorLease(t *testing.T) {
	p, incidents := newPlanner(t)
	now := time.Now().UTC()
	inc := incidents.Create("north-ridge", now)

	// Publish a perimeter so the air mission's flight plan can be created.
	p.perimeters.Publish(inc.ID, []model.Point{{Lat: 1, Lon: 1}}, nil, 0, now)

	first, err := p.Plan(model.DispatchRequest{
		OperationID: model.NewOperationID(),
		IncidentID:  inc.ID,
		AssigneeID:  "helo-1",
		Kind:        "air",
		RallyPoint:  "north gate",
		CorridorID:  "ridge-water",
	}, now)
	if err != nil {
		t.Fatalf("first plan: %v", err)
	}

	if got := p.resources.ActiveOwner("ridge-water"); got != 1 {
		t.Fatalf("owners after acquire = %d, want 1", got)
	}

	if err := p.Cancel(first.MissionID, now); err != nil {
		t.Fatalf("cancel: %v", err)
	}

	if got := p.resources.ActiveOwner("ridge-water"); got != 0 {
		t.Fatalf("owners after cancel = %d, want 0", got)
	}
	if plans := p.aircraft.Plans(); len(plans) != 0 {
		t.Fatalf("flight plans after cancel = %d, want 0", len(plans))
	}
	mission, ok := p.Get(first.MissionID)
	if !ok || mission.Stage != model.MissionCancelled {
		t.Fatalf("mission stage = %q, want cancelled", mission.Stage)
	}

	// A follow-on aircraft must now find an available water source.
	second, err := p.Plan(model.DispatchRequest{
		OperationID: model.NewOperationID(),
		IncidentID:  inc.ID,
		AssigneeID:  "helo-2",
		Kind:        "air",
		RallyPoint:  "north gate",
		CorridorID:  "ridge-water",
	}, now.Add(time.Second))
	if err != nil {
		t.Fatalf("follow-on plan after cancel: %v", err)
	}
	if second.MissionID == "" || second.MissionID == first.MissionID {
		t.Fatalf("follow-on mission id = %q", second.MissionID)
	}
}

func TestCancelIsIdempotent(t *testing.T) {
	p, incidents := newPlanner(t)
	now := time.Now().UTC()
	inc := incidents.Create("north-ridge", now)
	p.perimeters.Publish(inc.ID, []model.Point{{Lat: 1, Lon: 1}}, nil, 0, now)

	result, err := p.Plan(model.DispatchRequest{
		OperationID: model.NewOperationID(),
		IncidentID:  inc.ID,
		AssigneeID:  "helo-1",
		Kind:        "air",
		RallyPoint:  "north gate",
		CorridorID:  "ridge-water",
	}, now)
	if err != nil {
		t.Fatalf("plan: %v", err)
	}

	if err := p.Cancel(result.MissionID, now); err != nil {
		t.Fatalf("first cancel: %v", err)
	}
	if err := p.Cancel(result.MissionID, now); err != nil {
		t.Fatalf("second cancel: %v", err)
	}

	if got := p.resources.ActiveOwner("ridge-water"); got != 0 {
		t.Fatalf("owners after repeat cancel = %d, want 0", got)
	}
}
