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
	return nil
}
