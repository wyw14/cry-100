package dispatch

import (
	"github.com/wyw14/cry-100/internal/journal"
	"github.com/wyw14/cry-100/internal/model"
)

func Restore(recovery *journal.Recovery, planner *Planner, snapshot model.Snapshot) {
	planner.mu.Lock()
	defer planner.mu.Unlock()
	planner.missions = make(map[model.MissionID]model.Mission, len(snapshot.Missions))
	for _, mission := range recovery.Missions() {
		planner.missions[mission.ID] = mission
	}
}
