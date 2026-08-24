package radio

import (
	"testing"
	"time"

	"github.com/wyw14/cry-100/internal/model"
)

// TestMigrateCarriesPendingAckThread reproduces the duplicate-dispatch chain:
// a crew's task is dispatched under a child fire, the fires merge before the
// radio ack arrives, and the ack must then settle on the primary thread so the
// backup query sees it instead of re-dispatching a spare crew.
func TestMigrateCarriesPendingAckThread(t *testing.T) {
	now := time.Date(2026, 8, 24, 9, 0, 0, 0, time.UTC)
	child := model.IncidentID("incident-child")
	primary := model.IncidentID("incident-primary")
	mission := model.Mission{ID: model.MissionID("mission-a"), IncidentID: child}

	publisher := NewPublisher()
	if _, err := publisher.Send(model.OperationID("op-a"), mission, "Proceed to north gate", now); err != nil {
		t.Fatalf("send: %v", err)
	}

	// Merge happens before the field ack returns.
	publisher.Migrate(child, primary)

	// The crew confirms arrival over the radio after the merge.
	if err := publisher.Ack(mission.ID, "crew-a", now.Add(time.Minute)); err != nil {
		t.Fatalf("ack: %v", err)
	}

	message, ok := publisher.Get(mission.ID)
	if !ok {
		t.Fatal("radio message not found")
	}
	if message.AckedAt.IsZero() {
		t.Fatal("ack not applied to migrated thread")
	}
	if message.IncidentID != primary {
		t.Fatalf("ack thread stayed on sealed child %q, want primary %q", message.IncidentID, primary)
	}

	// The backup query scans the primary fire; it must see this acknowledgement
	// and therefore must not raise a duplicate dispatch for a spare crew.
	var seen bool
	for _, msg := range publisher.List() {
		if msg.IncidentID == child {
			t.Fatalf("radio thread leaked onto sealed child %q", child)
		}
		if msg.IncidentID == primary && msg.MissionID == mission.ID && !msg.AckedAt.IsZero() {
			seen = true
		}
	}
	if !seen {
		t.Fatal("backup query on primary saw no ack")
	}
}
