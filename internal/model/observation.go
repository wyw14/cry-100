package model

import "time"

type Observation struct {
	ID                 ObservationID `json:"id"`
	SensorID           SensorID      `json:"sensor_id"`
	SectorID           string        `json:"sector_id"`
	EventTime          time.Time     `json:"event_time"`
	ReceivedAt         time.Time     `json:"received_at"`
	TemperatureC       float64       `json:"temperature_c"`
	SmokeIndex         float64       `json:"smoke_index"`
	ImageHotspot       bool          `json:"image_hotspot"`
	Sequence           uint64        `json:"sequence"`
	CollectionEpoch    uint64        `json:"collection_epoch"`
	IncidentGeneration uint64        `json:"incident_generation"`
	Archived           bool          `json:"archived"`
}

func (o Observation) IndicatesFire() bool {
	return o.TemperatureC >= 55 || o.SmokeIndex >= 0.72 || o.ImageHotspot
}

type Calibration struct {
	Revision    uint64    `json:"revision"`
	Mode        string    `json:"mode"`
	Temperature float64   `json:"temperature_offset"`
	SmokeScale  float64   `json:"smoke_scale"`
	PublishedAt time.Time `json:"published_at"`
}

func (c Calibration) Apply(observation Observation) Observation {
	observation.TemperatureC += c.Temperature
	if c.SmokeScale != 0 {
		observation.SmokeIndex *= c.SmokeScale
	}
	return observation
}

type SensorProfile struct {
	SensorID    SensorID     `json:"sensor_id"`
	Mode        string       `json:"mode"`
	Epoch       uint64       `json:"epoch"`
	Revision    uint64       `json:"revision"`
	Calibration *Calibration `json:"calibration"`
}

func (p SensorProfile) Ready() bool {
	return p.Mode != "" && p.Calibration != nil && p.Calibration.Mode == p.Mode
}
