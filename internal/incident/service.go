package incident

import (
	"fmt"
	"github.com/wyw14/cry-100/internal/journal"
	"github.com/wyw14/cry-100/internal/model"
	"sync"
	"time"
)

type Service struct {
	mu        sync.RWMutex
	incidents map[model.IncidentID]model.Incident
	aliases   map[model.IncidentID]model.IncidentID
	journal   *journal.Store
}

func NewService(store *journal.Store) *Service {
	return &Service{incidents: map[model.IncidentID]model.Incident{}, aliases: map[model.IncidentID]model.IncidentID{}, journal: store}
}
func (s *Service) Create(sector string, now time.Time) model.Incident {
	s.mu.Lock()
	defer s.mu.Unlock()
	incident := model.NewIncident(sector, now)
	s.incidents[incident.ID] = incident
	s.record("incident.created", incident.ID, incident)
	return incident
}
func (s *Service) Get(id model.IncidentID) (model.Incident, bool) {
	id = s.Resolve(id)
	s.mu.RLock()
	defer s.mu.RUnlock()
	incident, ok := s.incidents[id]
	return incident, ok
}
func (s *Service) List() []model.Incident {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]model.Incident, 0, len(s.incidents))
	for _, incident := range s.incidents {
		result = append(result, incident)
	}
	return result
}
func (s *Service) Resolve(id model.IncidentID) model.IncidentID {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.resolve(id)
}
func (s *Service) resolve(id model.IncidentID) model.IncidentID {
	for {
		primary, ok := s.aliases[id]
		if !ok || primary == id {
			return id
		}
		id = primary
	}
}
func (s *Service) Close(id model.IncidentID, at time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	id = s.resolve(id)
	incident, ok := s.incidents[id]
	if !ok {
		return fmt.Errorf("incident not found")
	}
	incident.Stage, incident.ClosedAt, incident.UpdatedAt = model.StageClosed, at, at
	s.incidents[id] = incident
	s.record("incident.closed", id, incident)
	return nil
}
func (s *Service) Promote(candidate model.Candidate, observation model.Observation, now time.Time) (model.Incident, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, current := range s.incidents {
		if current.SectorID != candidate.SectorID {
			continue
		}
		if current.IsTerminal() && !observation.ReceivedAt.After(current.ClosedAt) {
			current.ArchivedCount++
			s.incidents[id] = current
			return current, false, nil
		}
		if current.IsActive() {
			current.Stage = model.StageAttacking
			current.UpdatedAt = now
			current.ObservationIDs = append(current.ObservationIDs, observation.ID)
			s.incidents[id] = current
			return current, true, nil
		}
	}
	incident := model.NewIncident(candidate.SectorID, now)
	incident.Stage = model.StageConfirmed
	incident.ObservationIDs = append(incident.ObservationIDs, observation.ID)
	s.incidents[incident.ID] = incident
	s.record("incident.confirmed", incident.ID, incident)
	return incident, true, nil
}
func (s *Service) record(kind string, id model.IncidentID, value any) {
	if s.journal != nil {
		_, _ = journal.AppendValue(s.journal, kind, id.String(), value)
	}
}

func (s *Service) Restore(incidents []model.Incident, aliases []model.IncidentAlias) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.incidents = make(map[model.IncidentID]model.Incident, len(incidents))
	s.aliases = make(map[model.IncidentID]model.IncidentID, len(aliases))
	for _, alias := range aliases {
		s.aliases[alias.ChildID] = alias.PrimaryID
	}
	for _, value := range incidents {
		s.incidents[value.ID] = value
	}
}

func (s *Service) Aliases() []model.IncidentAlias {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]model.IncidentAlias, 0, len(s.aliases))
	for child, primary := range s.aliases {
		result = append(result, model.IncidentAlias{ChildID: child, PrimaryID: primary})
	}
	return result
}
