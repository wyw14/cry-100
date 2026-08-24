package aircraft

import (
	"github.com/wyw14/cry-100/internal/model"
	"github.com/wyw14/cry-100/internal/resource"
	"sync"
	"time"
)

type Service struct {
	mu        sync.Mutex
	plans     map[model.MissionID]model.FlightPlan
	resources *resource.Registry
}

func NewService(resources *resource.Registry) *Service {
	return &Service{plans: map[model.MissionID]model.FlightPlan{}, resources: resources}
}
func (s *Service) Create(mission model.Mission, revision model.PerimeterRevision, at time.Time) (model.FlightPlan, error) {
	plan, err := Plan(mission, revision)
	if err != nil {
		return plan, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.plans[mission.ID] = plan
	return plan, nil
}
func (s *Service) Cancel(mission model.Mission, at time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.plans, mission.ID)
	if mission.LeaseID != "" {
		return s.resources.Release(mission.LeaseID, mission.ID, at)
	}
	return nil
}
func (s *Service) Plans() []model.FlightPlan {
	s.mu.Lock()
	defer s.mu.Unlock()
	v := make([]model.FlightPlan, 0, len(s.plans))
	for _, p := range s.plans {
		v = append(v, p)
	}
	return v
}
