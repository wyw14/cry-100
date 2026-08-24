package alert_test

import (
	"github.com/wyw14/cry-100/internal/alert"
	"github.com/wyw14/cry-100/internal/incident"
	"github.com/wyw14/cry-100/internal/journal"
	"github.com/wyw14/cry-100/internal/model"
	"testing"
	"time"
)

func TestAlertWithdrawalFailureRemainsVisible(t *testing.T) {
	now := time.Now().UTC()
	service := incident.NewService(nil)
	fire := service.Create("north-ridge", now)
	if err := service.Transition(fire.ID, model.StageConfirmed, now); err != nil {
		t.Fatal(err)
	}
	if err := service.Transition(fire.ID, model.StageAttacking, now); err != nil {
		t.Fatal(err)
	}
	gateway := alert.NewGateway(false)
	publisher := alert.NewPublisher(gateway)
	publisher.Publish(fire.ID, "village", "evacuate")
	evidence, _ := journal.NewEvidenceStore(t.TempDir())
	proof := model.FieldProof{IncidentID: fire.ID, ReportID: "r1", ImageName: "proof.jpg", SubmittedAt: now}
	if err := incident.Contain(service, publisher, evidence, proof, []byte("image"), now.Add(time.Minute)); err == nil {
		t.Fatal("withdrawal failure hidden")
	}
	current, _ := service.Get(fire.ID)
	if current.Stage == model.StageContained {
		t.Fatal("incident reported contained")
	}
	state, _ := publisher.Get(fire.ID)
	if !state.Active || !state.Pending || state.LastError == "" {
		t.Fatalf("alert=%+v", state)
	}
}
