package dispatch

import (
	"github.com/wyw14/cry-100/internal/model"
)

func (p *Planner) MigrateIncident(child, primary model.IncidentID) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for id, mission := range p.missions {
		if mission.IncidentID == child {
			mission.IncidentID = primary
			p.missions[id] = mission
		}
	}
	p.radio.Migrate(child, primary)
}
func (p *Planner) Acknowledge(id model.MissionID, by string) error { return p.radio.ApplyAck(id, by) }
