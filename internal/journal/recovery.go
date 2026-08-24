package journal

import (
	"github.com/wyw14/cry-100/internal/model"
	"sort"
)

type Recovery struct {
	incidents map[model.IncidentID]model.Incident
	missions  map[model.MissionID]model.Mission
	aliases   map[model.IncidentID]model.IncidentID
}

func NewRecovery() *Recovery {
	return &Recovery{incidents: map[model.IncidentID]model.Incident{}, missions: map[model.MissionID]model.Mission{}, aliases: map[model.IncidentID]model.IncidentID{}}
}
func (r *Recovery) Apply(snapshot model.Snapshot) {
	r.aliases = map[model.IncidentID]model.IncidentID{}
	for _, alias := range snapshot.Aliases {
		r.aliases[alias.ChildID] = r.resolve(alias.PrimaryID)
	}
	for _, incident := range snapshot.Incidents {
		incident.ID = r.resolve(incident.ID)
		if primary := r.aliases[incident.ID]; primary != "" {
			incident.ID = primary
		}
		r.incidents[incident.ID] = incident
	}
	for _, mission := range snapshot.Missions {
		mission.IncidentID = r.resolve(mission.IncidentID)
		r.missions[mission.ID] = mission
	}
}
func (r *Recovery) resolve(id model.IncidentID) model.IncidentID {
	seen := map[model.IncidentID]bool{}
	for {
		next, ok := r.aliases[id]
		if !ok || seen[id] {
			return id
		}
		seen[id] = true
		id = next
	}
}
func (r *Recovery) Incidents() []model.Incident {
	result := make([]model.Incident, 0, len(r.incidents))
	for _, incident := range r.incidents {
		result = append(result, incident)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result
}
func (r *Recovery) Missions() []model.Mission {
	result := make([]model.Mission, 0, len(r.missions))
	for _, mission := range r.missions {
		result = append(result, mission)
	}
	return result
}
