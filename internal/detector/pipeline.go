package detector

import (
	"github.com/wyw14/cry-100/internal/model"
	"sync"
)

type Pipeline struct {
	mu         sync.Mutex
	confidence float64
	last       map[model.SensorID]model.Observation
}

func NewPipeline() *Pipeline {
	return &Pipeline{confidence: 0.65, last: map[model.SensorID]model.Observation{}}
}
func (p *Pipeline) Analyze(observation model.Observation) (model.Candidate, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.last[observation.SensorID] = observation
	if !observation.IndicatesFire() {
		return model.Candidate{}, false
	}
	confidence := 0.55
	if observation.ImageHotspot {
		confidence += 0.25
	}
	if observation.SmokeIndex > 0.8 {
		confidence += 0.15
	}
	return model.Candidate{ObservationID: observation.ID, SectorID: observation.SectorID, EventTime: observation.EventTime, CollectionEpoch: observation.CollectionEpoch, IncidentGeneration: observation.IncidentGeneration, Confidence: confidence}, confidence >= p.confidence
}
