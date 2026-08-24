package sensor

import (
	"fmt"
	"github.com/wyw14/cry-100/internal/model"
	"sync"
	"time"
)

type Ingestor struct {
	mu       sync.RWMutex
	profiles map[model.SensorID]model.SensorProfile
	sequence map[model.SensorID]uint64
}

func NewIngestor() *Ingestor {
	return &Ingestor{profiles: map[model.SensorID]model.SensorProfile{}, sequence: map[model.SensorID]uint64{}}
}
func (i *Ingestor) Register(sensorID model.SensorID, mode string, calibration model.Calibration) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.profiles[sensorID] = model.SensorProfile{SensorID: sensorID, Mode: mode, Epoch: 1, Revision: calibration.Revision, Calibration: &calibration}
}
func (i *Ingestor) SetMode(sensorID model.SensorID, mode string, calibration *model.Calibration) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	profile, ok := i.profiles[sensorID]
	if !ok {
		return fmt.Errorf("unknown sensor")
	}
	updated, err := activateProfile(profile, mode, calibration)
	if err != nil {
		return err
	}
	i.profiles[sensorID] = updated
	return nil
}
func (i *Ingestor) Profile(sensorID model.SensorID) (model.SensorProfile, bool) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	profile, ok := i.profiles[sensorID]
	return profile, ok
}
func (i *Ingestor) Ingest(sensorID model.SensorID, sector string, eventAt time.Time, temperature, smoke float64, hotspot bool) (model.Observation, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	profile, ok := i.profiles[sensorID]
	if !ok || !profile.Ready() {
		return model.Observation{}, fmt.Errorf("sensor calibration unavailable")
	}
	i.sequence[sensorID]++
	observation := model.Observation{ID: model.NewObservationID(), SensorID: sensorID, SectorID: sector, EventTime: eventAt, ReceivedAt: time.Now().UTC(), TemperatureC: temperature, SmokeIndex: smoke, ImageHotspot: hotspot, Sequence: i.sequence[sensorID], CollectionEpoch: profile.Epoch}
	return profile.Calibration.Apply(observation), nil
}
