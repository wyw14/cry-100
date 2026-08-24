package dispatch

import (
	"fmt"
	"github.com/wyw14/cry-100/internal/model"
	"time"
)

func (p *Planner) Cancel(id model.MissionID, now time.Time) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	mission, ok := p.missions[id]
	if !ok {
		return fmt.Errorf("mission not found")
	}
	if mission.Stage == model.MissionCompleted {
		return nil
	}
	if mission.Stage != model.MissionCancelled {
		mission.Stage = model.MissionCancelled
		mission.UpdatedAt = now
	}
	if err := p.releaseMissionResources(mission, now); err != nil {
		p.missions[id] = mission
		return err
	}
	p.missions[id] = mission
	return nil
}

// releaseMissionResources tears down the corridor lease and flight plan booked
// during Plan so the cancelled mission holds no inventory. Every step is
// idempotent: a missing lease or plan is not an error, and releasing an
// already-released lease is a no-op, so cancelling twice converges to the
// same state even when an earlier cancel left a lease behind.
func (p *Planner) releaseMissionResources(mission model.Mission, now time.Time) error {
	if mission.LeaseID != "" {
		if err := p.resources.Release(mission.LeaseID, mission.ID, now); err != nil {
			return err
		}
	}
	if mission.Kind == "air" {
		_ = p.aircraft.Cancel(mission, now)
	}
	return nil
}
