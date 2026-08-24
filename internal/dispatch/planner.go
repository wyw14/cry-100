package dispatch

import (
	"fmt"
	"github.com/wyw14/cry-100/internal/aircraft"
	"github.com/wyw14/cry-100/internal/incident"
	"github.com/wyw14/cry-100/internal/model"
	"github.com/wyw14/cry-100/internal/perimeter"
	"github.com/wyw14/cry-100/internal/radio"
	"github.com/wyw14/cry-100/internal/resource"
	"sync"
	"time"
)

type Planner struct {
	mu         sync.Mutex
	missions   map[model.MissionID]model.Mission
	operations map[model.OperationID]model.DispatchResult
	incidents  *incident.Service
	resources  *resource.Registry
	radio      *radio.Publisher
	aircraft   *aircraft.Service
	perimeters *perimeter.Service
	recorder   ResultRecorder
}

type ResultRecorder interface {
	RecordOperation(model.OperationRecord) error
}

func NewPlanner(incidents *incident.Service, resources *resource.Registry, radioPublisher *radio.Publisher, aircraftService *aircraft.Service, perimeters *perimeter.Service, recorders ...ResultRecorder) *Planner {
	planner := &Planner{missions: map[model.MissionID]model.Mission{}, operations: map[model.OperationID]model.DispatchResult{}, incidents: incidents, resources: resources, radio: radioPublisher, aircraft: aircraftService, perimeters: perimeters}
	if len(recorders) > 0 {
		planner.recorder = recorders[0]
	}
	return planner
}
func (p *Planner) Plan(request model.DispatchRequest, now time.Time) (model.DispatchResult, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	request.OperationID = operationIdentity(request)
	if result, ok := p.operations[request.OperationID]; ok {
		result.Recovered = true
		return result, nil
	}
	stableID := model.StableMissionID(request.OperationID)
	if existing, ok := p.missions[stableID]; ok {
		if message, sent := p.radio.Get(existing.ID); sent {
			result := model.DispatchResult{MissionID: existing.ID, MessageID: message.ID, Recovered: true}
			p.operations[request.OperationID] = result
			return result, nil
		}
	}
	inc, ok := p.incidents.Get(request.IncidentID)
	if !ok {
		return model.DispatchResult{}, fmt.Errorf("incident not found")
	}
	if err := validateAdmission(request, inc); err != nil {
		return model.DispatchResult{}, err
	}
	mission := model.Mission{ID: stableID, OperationID: request.OperationID, IncidentID: inc.ID, AssigneeID: request.AssigneeID, Kind: request.Kind, Stage: model.MissionPlanned, RallyPoint: request.RallyPoint, CorridorID: request.CorridorID, CreatedAt: now, UpdatedAt: now}
	if request.CorridorID != "" {
		lease, err := p.resources.Acquire(request.CorridorID, mission.ID, request.AssigneeID, now)
		if err != nil {
			return model.DispatchResult{}, err
		}
		mission.LeaseID = lease.ID
		mission.Stage = model.MissionReserved
	}
	if request.Kind == "air" {
		revision, ok := p.perimeters.Current(inc.ID)
		if !ok {
			return model.DispatchResult{}, fmt.Errorf("perimeter not found")
		}
		if _, err := p.aircraft.Create(mission, revision, now); err != nil {
			return model.DispatchResult{}, err
		}
	}
	if err := p.incidents.AddMission(inc.ID, mission); err != nil {
		return model.DispatchResult{}, err
	}
	message, err := p.radio.Send(request.OperationID, mission, "Proceed to "+request.RallyPoint, now)
	if err != nil {
		return model.DispatchResult{}, err
	}
	result := model.DispatchResult{MissionID: mission.ID, MessageID: message.ID}
	p.missions[mission.ID] = mission
	if p.recorder != nil {
		record := model.OperationRecord{OperationID: request.OperationID, Result: result, CommittedAt: now}
		if err := p.recorder.RecordOperation(record); err != nil {
			return model.DispatchResult{}, err
		}
	}
	p.operations[request.OperationID] = result
	return result, nil
}
func (p *Planner) Get(id model.MissionID) (model.Mission, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	v, ok := p.missions[id]
	return v, ok
}
func (p *Planner) List() []model.Mission {
	p.mu.Lock()
	defer p.mu.Unlock()
	v := make([]model.Mission, 0, len(p.missions))
	for _, m := range p.missions {
		v = append(v, m)
	}
	return v
}
