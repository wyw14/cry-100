package incident

import (
	"testing"
	"time"

	"github.com/wyw14/cry-100/internal/journal"
	"github.com/wyw14/cry-100/internal/model"
)

// Verifies: when evidence write fails, the incident must NOT advance to contained.
func TestContainDoesNotAdvanceWhenEvidenceFails(t *testing.T) {
	dir := t.TempDir()
	store, err := journal.NewEvidenceStore(dir, true) // read-only -> Save fails
	if err != nil {
		t.Fatal(err)
	}
	svc := NewService(nil)
	incident := svc.Create("north-ridge", time.Now().UTC())
	incident.Stage = model.StageAttacking
	svc.incidents[incident.ID] = incident

	proof := model.FieldProof{IncidentID: incident.ID, ReportID: "rep-1", ImageName: "img-1.bin", SubmittedAt: time.Now().UTC()}
	err = Contain(svc, nil, store, proof, []byte{1, 2, 3}, time.Now().UTC())
	if err == nil {
		t.Fatal("expected evidence write failure to surface as error")
	}
	got, _ := svc.Get(incident.ID)
	if got.Stage != model.StageAttacking {
		t.Fatalf("incident advanced to %s after evidence failure; must remain attacking", got.Stage)
	}
}

// Verifies: when evidence is durable, the incident advances to contained.
func TestContainAdvancesWhenEvidenceDurable(t *testing.T) {
	dir := t.TempDir()
	store, err := journal.NewEvidenceStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := NewService(nil)
	incident := svc.Create("north-ridge", time.Now().UTC())
	incident.Stage = model.StageAttacking
	svc.incidents[incident.ID] = incident

	proof := model.FieldProof{IncidentID: incident.ID, ReportID: "rep-2", ImageName: "img-2.bin", SubmittedAt: time.Now().UTC()}
	if err := Contain(svc, nil, store, proof, []byte{1, 2, 3}, time.Now().UTC()); err != nil {
		t.Fatalf("expected durable evidence to contain incident: %v", err)
	}
	got, _ := svc.Get(incident.ID)
	if got.Stage != model.StageContained {
		t.Fatalf("incident stage=%s; want contained", got.Stage)
	}
	if !store.Has(proof) {
		t.Fatal("durable evidence missing after successful contain")
	}
}
