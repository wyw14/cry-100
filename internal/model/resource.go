package model

import "time"

type Point struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}
type PerimeterRevision struct {
	ID          RevisionID `json:"id"`
	IncidentID  IncidentID `json:"incident_id"`
	Number      uint64     `json:"number"`
	Boundary    []Point    `json:"boundary"`
	NoFlyZones  []string   `json:"no_fly_zones"`
	WindBearing int        `json:"wind_bearing"`
	ApprovedAt  time.Time  `json:"approved_at"`
}
type Corridor struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
	Owners   int    `json:"owners"`
}
type Lease struct {
	ID         LeaseID   `json:"id"`
	CorridorID string    `json:"corridor_id"`
	MissionID  MissionID `json:"mission_id"`
	Owner      string    `json:"owner"`
	Token      uint64    `json:"fencing_token"`
	AcquiredAt time.Time `json:"acquired_at"`
	ReleasedAt time.Time `json:"released_at,omitempty"`
}

func (l Lease) Active() bool { return l.ReleasedAt.IsZero() }

type FieldProof struct {
	IncidentID  IncidentID `json:"incident_id"`
	ReportID    string     `json:"report_id"`
	Summary     string     `json:"summary"`
	ImageName   string     `json:"image_name"`
	SubmittedAt time.Time  `json:"submitted_at"`
}
type PublicAlert struct {
	ID         string     `json:"id"`
	IncidentID IncidentID `json:"incident_id"`
	Area       string     `json:"area"`
	Message    string     `json:"message"`
	Active     bool       `json:"active"`
	Pending    bool       `json:"pending_withdrawal"`
	UpdatedAt  time.Time  `json:"updated_at"`
	LastError  string     `json:"last_error,omitempty"`
}
