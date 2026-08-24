package incident_test

import (
	"github.com/wyw14/cry-100/internal/incident"
	"github.com/wyw14/cry-100/internal/journal"
	"github.com/wyw14/cry-100/internal/model"
	"testing"
	"time"
)

func TestContainmentWaitsForDurableFieldProof(t *testing.T) {
	now := time.Now().UTC()
	service := incident.NewService(nil)
	fire := service.Create("north-ridge", now)
	if err := service.Transition(fire.ID, model.StageConfirmed, now); err != nil {
		t.Fatal(err)
	}
	if err := service.Transition(fire.ID, model.StageAttacking, now); err != nil {
		t.Fatal(err)
	}
	evidence, err := journal.NewEvidenceStore(t.TempDir(), true)
	if err != nil {
		t.Fatal(err)
	}
	proof := model.FieldProof{IncidentID: fire.ID, ReportID: "report-1", ImageName: "line.jpg", SubmittedAt: now}
	if err = incident.Contain(service, nil, evidence, proof, []byte("image"), now.Add(time.Minute)); err == nil {
		t.Fatal("containment unexpectedly succeeded")
	}
	current, _ := service.Get(fire.ID)
	if current.Stage != model.StageAttacking {
		t.Fatalf("stage=%s", current.Stage)
	}
}
