package model

import (
	"encoding/json"
	"time"
)

type Event struct {
	ID          string          `json:"id"`
	Kind        string          `json:"kind"`
	AggregateID string          `json:"aggregate_id"`
	OccurredAt  time.Time       `json:"occurred_at"`
	Payload     json.RawMessage `json:"payload"`
}

func NewEvent(kind, aggregateID string, value any, now time.Time) (Event, error) {
	payload, err := json.Marshal(value)
	if err != nil {
		return Event{}, err
	}
	return Event{ID: NewMessageID().String(), Kind: kind, AggregateID: aggregateID, OccurredAt: now, Payload: payload}, nil
}
func (id MessageID) String() string  { return string(id) }
func (id IncidentID) String() string { return string(id) }
func (id MissionID) String() string  { return string(id) }

type OperationRecord struct {
	OperationID OperationID    `json:"operation_id"`
	Result      DispatchResult `json:"result"`
	CommittedAt time.Time      `json:"committed_at"`
}
type Snapshot struct {
	CreatedAt time.Time       `json:"created_at"`
	Aliases   []IncidentAlias `json:"aliases"`
	Incidents []Incident      `json:"incidents"`
	Missions  []Mission       `json:"missions"`
	Leases    []Lease         `json:"leases"`
}
