package detector_test

import (
	"github.com/wyw14/cry-100/internal/detector"
	"github.com/wyw14/cry-100/internal/model"
	"github.com/wyw14/cry-100/internal/sensor"
	"testing"
	"time"
)

func TestModeActivationPublishesCalibrationAtomically(t *testing.T) {
	now := time.Now().UTC()
	id := model.SensorID(model.NewMessageID().String())
	ingestor := sensor.NewIngestor()
	visible := model.Calibration{Revision: 1, Mode: "visible", SmokeScale: 1, PublishedAt: now}
	ingestor.Register(id, "visible", visible)
	combined := model.Calibration{Revision: 2, Mode: "combined", SmokeScale: 1, PublishedAt: now}
	if err := ingestor.SetMode(id, "combined", &combined); err != nil {
		t.Fatal(err)
	}
	observation, err := ingestor.Ingest(id, "north-ridge", now, 70, .8, true)
	if err != nil {
		t.Fatal(err)
	}
	candidate, ok := detector.NewPipeline().Analyze(observation)
	if !ok || candidate.CollectionEpoch != 2 {
		t.Fatalf("candidate=%+v ok=%v", candidate, ok)
	}
}
