package sensor

import (
	"github.com/wyw14/cry-100/internal/model"
	"sync"
)

type Stream struct {
	mu     sync.RWMutex
	values []model.Observation
}

func NewStream() *Stream { return &Stream{} }
func (s *Stream) Add(observation model.Observation) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values = append(s.values, observation)
}
func (s *Stream) All() []model.Observation {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]model.Observation(nil), s.values...)
}
