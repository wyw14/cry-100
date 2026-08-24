package crew

import "sync"

type Crew struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	SectorID  string `json:"sector_id"`
	Available bool   `json:"available"`
	Progress  string `json:"progress"`
}
type Registry struct {
	mu    sync.RWMutex
	crews map[string]Crew
}

func NewRegistry() *Registry      { return &Registry{crews: map[string]Crew{}} }
func (r *Registry) Add(crew Crew) { r.mu.Lock(); defer r.mu.Unlock(); r.crews[crew.ID] = crew }
func (r *Registry) List() []Crew {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v := make([]Crew, 0, len(r.crews))
	for _, c := range r.crews {
		v = append(v, c)
	}
	return v
}
