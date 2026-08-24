package app

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/wyw14/cry-100/internal/aircraft"
	"github.com/wyw14/cry-100/internal/alert"
	"github.com/wyw14/cry-100/internal/crew"
	"github.com/wyw14/cry-100/internal/detector"
	"github.com/wyw14/cry-100/internal/dispatch"
	"github.com/wyw14/cry-100/internal/incident"
	"github.com/wyw14/cry-100/internal/journal"
	"github.com/wyw14/cry-100/internal/model"
	"github.com/wyw14/cry-100/internal/perimeter"
	"github.com/wyw14/cry-100/internal/radio"
	"github.com/wyw14/cry-100/internal/resource"
	"github.com/wyw14/cry-100/internal/sector"
	"github.com/wyw14/cry-100/internal/sensor"
)

type Runtime struct {
	Journal        *journal.Store
	Evidence       *journal.EvidenceStore
	Grid           *sector.Grid
	Sensors        *sensor.Ingestor
	Streams        *sensor.Stream
	Detector       *detector.Pipeline
	Incidents      *incident.Service
	Perimeters     *perimeter.Service
	Resources      *resource.Registry
	Radio          *radio.Publisher
	Alerts         *alert.Publisher
	Aircraft       *aircraft.Service
	Crews          *crew.Registry
	Dispatch       *dispatch.Planner
	ProtocolSteps  int
	CatalogEntries int
}

func New(dataDir string) (*Runtime, error) {
	store, err := journal.New(filepath.Join(dataDir, "events"))
	if err != nil {
		return nil, err
	}
	evidence, err := journal.NewEvidenceStore(filepath.Join(dataDir, "evidence"))
	if err != nil {
		return nil, err
	}
	resources := resource.NewRegistry()
	resources.AddCorridor("ridge-water", "Ridge Water", 1)
	radioPublisher := radio.NewPublisher()
	perimeters := perimeter.NewService()
	incidents := incident.NewService(store)
	aircraftService := aircraft.NewService(resources)
	alerts := alert.NewPublisher(alert.NewGateway())
	runtime := &Runtime{Journal: store, Evidence: evidence, Grid: sector.DefaultGrid(), Sensors: sensor.NewIngestor(), Streams: sensor.NewStream(), Detector: detector.NewPipeline(), Incidents: incidents, Perimeters: perimeters, Resources: resources, Radio: radioPublisher, Alerts: alerts, Aircraft: aircraftService, Crews: crew.NewRegistry()}
	runtime.Dispatch = dispatch.NewPlanner(incidents, resources, radioPublisher, aircraftService, perimeters, store)
	runtime.ProtocolSteps = ProtocolStepCount()
	runtime.CatalogEntries = ChecklistCount() + OperationLabelCount() + ResponseCodeCount() + StatusCatalogCount() + ScenarioCount() + len(model.FirelineCodes)
	if err := runtime.restore(); err != nil {
		return nil, err
	}
	return runtime, nil
}

func (r *Runtime) restore() error {
	snapshot, err := r.Journal.LoadSnapshot()
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	recovery := journal.NewRecovery()
	recovery.Apply(snapshot)
	r.Incidents.Restore(recovery.Incidents(), snapshot.Aliases)
	dispatch.Restore(recovery, r.Dispatch, snapshot)
	r.Resources.Restore(snapshot.Leases)
	return nil
}

func (r *Runtime) SaveSnapshot() error {
	snapshot := model.Snapshot{
		CreatedAt: time.Now().UTC(),
		Aliases:   r.Incidents.Aliases(),
		Incidents: r.Incidents.List(),
		Missions:  r.Dispatch.List(),
		Leases:    r.Resources.Leases(),
	}
	return r.Journal.SaveSnapshot(snapshot)
}

func (r *Runtime) MergeIncidents(primary, child model.IncidentID, at time.Time) (incident.MergeResult, error) {
	result, err := r.Incidents.Merge(primary, child, at)
	if err != nil {
		return result, err
	}
	r.Dispatch.MigrateIncident(child, primary)
	return result, nil
}

func (r *Runtime) ContainIncident(proof model.FieldProof, image []byte, at time.Time) error {
	report := crew.BuildFieldReport(proof, image)
	if err := report.Validate(); err != nil {
		return err
	}
	return incident.Contain(r.Incidents, r.Alerts, r.Evidence, report.Proof(), image, at)
}
