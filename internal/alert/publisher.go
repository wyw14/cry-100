package alert

import (
	"fmt"
	"github.com/wyw14/cry-100/internal/model"
	"sync"
	"time"
)

type Channel interface{ Withdraw(model.IncidentID) error }
type Publisher struct {
	mu      sync.RWMutex
	alerts  map[model.IncidentID]model.PublicAlert
	channel Channel
	now     func() time.Time
}

func NewPublisher(channel Channel) *Publisher {
	return &Publisher{alerts: map[model.IncidentID]model.PublicAlert{}, channel: channel, now: time.Now}
}
func (p *Publisher) Publish(incidentID model.IncidentID, area, message string) model.PublicAlert {
	p.mu.Lock()
	defer p.mu.Unlock()
	alert := model.PublicAlert{ID: model.NewMessageID().String(), IncidentID: incidentID, Area: area, Message: message, Active: true, UpdatedAt: p.now().UTC()}
	p.alerts[incidentID] = alert
	return alert
}
func (p *Publisher) Get(incidentID model.IncidentID) (model.PublicAlert, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	value, ok := p.alerts[incidentID]
	return value, ok
}
func (p *Publisher) Withdraw(incidentID model.IncidentID) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	value, ok := p.alerts[incidentID]
	if !ok {
		return fmt.Errorf("alert not found")
	}
	if p.channel != nil {
		if err := p.channel.Withdraw(incidentID); err != nil {
			value.Pending, value.LastError, value.UpdatedAt = true, err.Error(), p.now().UTC()
			p.alerts[incidentID] = value
			return err
		}
	}
	value.Active, value.Pending, value.LastError, value.UpdatedAt = false, false, "", p.now().UTC()
	p.alerts[incidentID] = value
	return nil
}
func (p *Publisher) List() []model.PublicAlert {
	p.mu.RLock()
	defer p.mu.RUnlock()
	values := make([]model.PublicAlert, 0, len(p.alerts))
	for _, value := range p.alerts {
		values = append(values, value)
	}
	return values
}
