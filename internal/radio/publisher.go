package radio

import (
	"fmt"
	"github.com/wyw14/cry-100/internal/model"
	"sync"
	"time"
)

type Publisher struct {
	mu       sync.RWMutex
	messages map[model.MissionID]model.RadioMessage
}

func NewPublisher() *Publisher {
	return &Publisher{messages: map[model.MissionID]model.RadioMessage{}}
}
func (p *Publisher) Send(operation model.OperationID, mission model.Mission, body string, now time.Time) (model.RadioMessage, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	message := model.RadioMessage{ID: model.NewMessageID(), MissionID: mission.ID, IncidentID: mission.IncidentID, Body: body, SentAt: now}
	p.messages[mission.ID] = message
	return message, nil
}
func (p *Publisher) Ack(missionID model.MissionID, by string, now time.Time) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	message, ok := p.messages[missionID]
	if !ok {
		return fmt.Errorf("radio message not found")
	}
	message.AckedAt, message.AckBy = now, by
	p.messages[missionID] = message
	return nil
}
func (p *Publisher) Get(missionID model.MissionID) (model.RadioMessage, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	value, ok := p.messages[missionID]
	return value, ok
}
func (p *Publisher) Migrate(child, primary model.IncidentID) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for id, message := range p.messages {
		if message.IncidentID == child && !message.AckedAt.IsZero() {
			message.IncidentID = primary
			p.messages[id] = message
		}
	}
}
func (p *Publisher) List() []model.RadioMessage {
	p.mu.RLock()
	defer p.mu.RUnlock()
	values := make([]model.RadioMessage, 0, len(p.messages))
	for _, value := range p.messages {
		values = append(values, value)
	}
	return values
}
