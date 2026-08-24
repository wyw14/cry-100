package journal_test

import (
	"github.com/wyw14/cry-100/internal/journal"
	"github.com/wyw14/cry-100/internal/model"
	"testing"
	"time"
)

func TestRecoveryAppliesIncidentAliasesBeforeMissions(t *testing.T) {
	now := time.Now().UTC()
	primary := model.NewIncident("north-ridge", now)
	child := model.NewIncident("pine-valley", now)
	mission := model.Mission{ID: model.StableMissionID(model.NewOperationID()), IncidentID: child.ID, AssigneeID: "crew-a", Kind: "ground", Stage: model.MissionActive, CreatedAt: now, UpdatedAt: now}
	snapshot := model.Snapshot{CreatedAt: now, Aliases: []model.IncidentAlias{{ChildID: child.ID, PrimaryID: primary.ID, MergedAt: now}}, Incidents: []model.Incident{primary, child}, Missions: []model.Mission{mission}}
	recovery := journal.NewRecovery()
	recovery.Apply(snapshot)
	missions := recovery.Missions()
	if len(missions) != 1 || missions[0].IncidentID != primary.ID {
		t.Fatalf("mission owner=%v expected=%s", missions, primary.ID)
	}
	for _, value := range recovery.Incidents() {
		if value.ID == child.ID && value.IsActive() {
			t.Fatal("child incident became active")
		}
	}
}
