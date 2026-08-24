package model

import "time"

type IncidentStage string

const (
	StageCandidate IncidentStage = "candidate"
	StageConfirmed IncidentStage = "confirmed"
	StageAttacking IncidentStage = "attacking"
	StageContained IncidentStage = "contained"
	StageClosed    IncidentStage = "closed"
)

type Incident struct {
	ID             IncidentID      `json:"id"`
	SectorID       string          `json:"sector_id"`
	Stage          IncidentStage   `json:"stage"`
	Generation     uint64          `json:"generation"`
	PrimaryID      IncidentID      `json:"primary_id,omitempty"`
	OpenedAt       time.Time       `json:"opened_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	ContainedAt    time.Time       `json:"contained_at,omitempty"`
	ClosedAt       time.Time       `json:"closed_at,omitempty"`
	ObservationIDs []ObservationID `json:"observation_ids"`
	MissionIDs     []MissionID     `json:"mission_ids"`
	ArchivedCount  int             `json:"archived_observation_count"`
}

func NewIncident(sectorID string, at time.Time) Incident {
	return Incident{ID: NewIncidentID(), SectorID: sectorID, Stage: StageCandidate, Generation: 1, OpenedAt: at, UpdatedAt: at}
}

func (i Incident) IsTerminal() bool { return i.Stage == StageClosed }
func (i Incident) IsActive() bool {
	return i.Stage == StageCandidate || i.Stage == StageConfirmed || i.Stage == StageAttacking
}

type Candidate struct {
	ObservationID      ObservationID `json:"observation_id"`
	SectorID           string        `json:"sector_id"`
	EventTime          time.Time     `json:"event_time"`
	CollectionEpoch    uint64        `json:"collection_epoch"`
	IncidentGeneration uint64        `json:"incident_generation"`
	Confidence         float64       `json:"confidence"`
}

type IncidentAlias struct {
	ChildID   IncidentID `json:"child_id"`
	PrimaryID IncidentID `json:"primary_id"`
	MergedAt  time.Time  `json:"merged_at"`
}
