package incident

import (
	"fmt"
	"github.com/wyw14/cry-100/internal/model"
	"time"
)

func CanTransition(from, to model.IncidentStage) bool {
	switch from {
	case model.StageCandidate:
		return to == model.StageConfirmed || to == model.StageClosed
	case model.StageConfirmed:
		return to == model.StageAttacking || to == model.StageContained
	case model.StageAttacking:
		return to == model.StageContained || to == model.StageClosed
	case model.StageContained:
		return to == model.StageClosed
	case model.StageClosed:
		return false
	}
	return false
}
func (s *Service) Transition(id model.IncidentID, to model.IncidentStage, at time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	id = s.resolve(id)
	incident, ok := s.incidents[id]
	if !ok {
		return fmt.Errorf("incident not found")
	}
	if !CanTransition(incident.Stage, to) {
		return fmt.Errorf("cannot transition %s to %s", incident.Stage, to)
	}
	incident.Stage, incident.UpdatedAt = to, at
	if to == model.StageContained {
		incident.ContainedAt = at
	}
	if to == model.StageClosed {
		incident.ClosedAt = at
	}
	s.incidents[id] = incident
	s.record("incident.transition", id, incident)
	return nil
}
func (s *Service) AddMission(id model.IncidentID, mission model.Mission) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	id = s.resolve(id)
	incident, ok := s.incidents[id]
	if !ok {
		return fmt.Errorf("incident not found")
	}
	incident.MissionIDs = append(incident.MissionIDs, mission.ID)
	incident.UpdatedAt = mission.UpdatedAt
	s.incidents[id] = incident
	return nil
}
