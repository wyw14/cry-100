package model

import "time"

type MissionStage string

const (
	MissionPlanned      MissionStage = "planned"
	MissionReserved     MissionStage = "reserved"
	MissionAssigned     MissionStage = "assigned"
	MissionEnroute      MissionStage = "enroute"
	MissionActive       MissionStage = "active"
	MissionCompleted    MissionStage = "completed"
	MissionCancelled    MissionStage = "cancelled"
	MissionCompensating MissionStage = "compensating"
)

type Mission struct {
	ID                MissionID    `json:"id"`
	OperationID       OperationID  `json:"operation_id"`
	IncidentID        IncidentID   `json:"incident_id"`
	AssigneeID        string       `json:"assignee_id"`
	Kind              string       `json:"kind"`
	Stage             MissionStage `json:"stage"`
	RallyPoint        string       `json:"rally_point"`
	CorridorID        string       `json:"corridor_id,omitempty"`
	LeaseID           LeaseID      `json:"lease_id,omitempty"`
	PerimeterRevision RevisionID   `json:"perimeter_revision,omitempty"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at"`
}

func (m Mission) Finished() bool {
	switch m.Stage {
	case MissionAssigned, MissionEnroute, MissionActive:
		return false
	case MissionCompleted, MissionCancelled:
		return true
	default:
		return false
	}
}

type DispatchRequest struct {
	OperationID OperationID `json:"operation_id"`
	IncidentID  IncidentID  `json:"incident_id"`
	AssigneeID  string      `json:"assignee_id"`
	Kind        string      `json:"kind"`
	RallyPoint  string      `json:"rally_point"`
	CorridorID  string      `json:"corridor_id,omitempty"`
}

type DispatchResult struct {
	MissionID MissionID `json:"mission_id"`
	MessageID MessageID `json:"message_id"`
	Recovered bool      `json:"recovered"`
}

type RadioMessage struct {
	ID         MessageID  `json:"id"`
	MissionID  MissionID  `json:"mission_id"`
	IncidentID IncidentID `json:"incident_id"`
	Body       string     `json:"body"`
	SentAt     time.Time  `json:"sent_at"`
	AckedAt    time.Time  `json:"acked_at,omitempty"`
	AckBy      string     `json:"ack_by,omitempty"`
}

type FlightPlan struct {
	MissionID        MissionID  `json:"mission_id"`
	IncidentID       IncidentID `json:"incident_id"`
	RevisionID       RevisionID `json:"revision_id"`
	RevisionTime     time.Time  `json:"revision_time"`
	Waypoints        []Point    `json:"waypoints"`
	AvoidedCorridors []string   `json:"avoided_corridors"`
	Frozen           bool       `json:"frozen"`
}
