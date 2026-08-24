package resource

import (
	"fmt"
	"github.com/wyw14/cry-100/internal/model"
	"sync"
	"time"
)

type Registry struct {
	mu        sync.Mutex
	corridors map[string]model.Corridor
	leases    map[model.LeaseID]model.Lease
	nextToken uint64
}

func NewRegistry() *Registry {
	return &Registry{corridors: map[string]model.Corridor{}, leases: map[model.LeaseID]model.Lease{}}
}
func (r *Registry) AddCorridor(id, name string, capacity int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.corridors[id] = model.Corridor{ID: id, Name: name, Capacity: capacity}
}
func (r *Registry) Corridor(id string) (model.Corridor, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	value, ok := r.corridors[id]
	return value, ok
}
func (r *Registry) Acquire(corridor string, mission model.MissionID, owner string, now time.Time) (model.Lease, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	value, ok := r.corridors[corridor]
	if !ok {
		return model.Lease{}, fmt.Errorf("corridor not found")
	}
	if value.Owners >= value.Capacity {
		return model.Lease{}, fmt.Errorf("corridor %s is occupied", corridor)
	}
	r.nextToken++
	lease := model.Lease{ID: model.NewLeaseID(), CorridorID: corridor, MissionID: mission, Owner: owner, Token: r.nextToken, AcquiredAt: now}
	value.Owners++
	r.corridors[corridor] = value
	r.leases[lease.ID] = lease
	return lease, nil
}
func (r *Registry) Release(id model.LeaseID, mission model.MissionID, now time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	lease, ok := r.leases[id]
	if !ok {
		return fmt.Errorf("lease not found")
	}
	if lease.MissionID != mission {
		return fmt.Errorf("lease owner mismatch")
	}
	if !lease.ReleasedAt.IsZero() {
		return nil
	}
	lease.ReleasedAt = now
	r.leases[id] = lease
	value := r.corridors[lease.CorridorID]
	if value.Owners > 0 {
		value.Owners--
	}
	r.corridors[lease.CorridorID] = value
	return nil
}
func (r *Registry) Leases() []model.Lease {
	r.mu.Lock()
	defer r.mu.Unlock()
	values := make([]model.Lease, 0, len(r.leases))
	for _, value := range r.leases {
		values = append(values, value)
	}
	return values
}
func (r *Registry) ActiveOwner(corridor string) int {
	value, _ := r.Corridor(corridor)
	return value.Owners
}

func (r *Registry) Restore(leases []model.Lease) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, lease := range leases {
		r.leases[lease.ID] = lease
		if lease.Token > r.nextToken {
			r.nextToken = lease.Token
		}
		if lease.Active() {
			corridor := r.corridors[lease.CorridorID]
			corridor.Owners++
			r.corridors[lease.CorridorID] = corridor
		}
	}
}
