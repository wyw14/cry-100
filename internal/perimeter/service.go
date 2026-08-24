package perimeter

import (
	"github.com/wyw14/cry-100/internal/model"
	"sync"
	"time"
)

type Service struct {
	mu        sync.RWMutex
	revisions map[model.IncidentID]model.PerimeterRevision
}

func NewService() *Service {
	return &Service{revisions: map[model.IncidentID]model.PerimeterRevision{}}
}
func (s *Service) Publish(incident model.IncidentID, boundary []model.Point, zones []string, wind int, at time.Time) model.PerimeterRevision {
	s.mu.Lock()
	defer s.mu.Unlock()
	previous := s.revisions[incident]
	revision := model.PerimeterRevision{ID: model.NewRevisionID(), IncidentID: incident, Number: previous.Number + 1, Boundary: append([]model.Point(nil), boundary...), NoFlyZones: append([]string(nil), zones...), WindBearing: wind, ApprovedAt: at}
	if previous.ID == "" {
		s.revisions[incident] = revision
	}
	return revision
}
func (s *Service) Current(incident model.IncidentID) (model.PerimeterRevision, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.revisions[incident]
	return value, ok
}
func (s *Service) All() []model.PerimeterRevision {
	s.mu.RLock()
	defer s.mu.RUnlock()
	values := make([]model.PerimeterRevision, 0, len(s.revisions))
	for _, v := range s.revisions {
		values = append(values, v)
	}
	return values
}
