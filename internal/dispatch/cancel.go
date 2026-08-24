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
	if mission.Finished() {
		return nil
	}
	mission.Stage = model.MissionCancelled
	mission.UpdatedAt = now
	p.missions[id] = mission
	if mission.Kind == "air" {
		if err := p.aircraft.Cancel(mission, now); err != nil {
			mission.Stage = model.MissionCompensating
			p.missions[id] = mission
			return err
		}
	} else if mission.LeaseID != "" {
		if err := p.resources.Release(mission.LeaseID, mission.ID, now); err != nil {
			mission.Stage = model.MissionCompensating
			p.missions[id] = mission
			return err
		}
	}
	return nil
}
