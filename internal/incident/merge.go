package incident

import (
	"fmt"
	"github.com/wyw14/cry-100/internal/model"
	"time"
)

type MergeResult struct {
	Primary model.Incident
	Alias   model.IncidentAlias
}

func (s *Service) Merge(primaryID, childID model.IncidentID, at time.Time) (MergeResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	primaryID, childID = s.resolve(primaryID), s.resolve(childID)
	primary, ok := s.incidents[primaryID]
	if !ok {
		return MergeResult{}, fmt.Errorf("primary incident not found")
	}
	child, ok := s.incidents[childID]
	if !ok {
		return MergeResult{}, fmt.Errorf("child incident not found")
	}
	if primaryID == childID || !primary.IsActive() || !child.IsActive() {
		return MergeResult{}, fmt.Errorf("incidents cannot be merged")
	}
	primary.ObservationIDs = append(primary.ObservationIDs, child.ObservationIDs...)
	primary.MissionIDs = append(primary.MissionIDs, child.MissionIDs...)
	primary.UpdatedAt = at
	child.PrimaryID, child.Stage, child.UpdatedAt = primaryID, model.StageClosed, at
	s.incidents[primaryID], s.incidents[childID] = primary, child
	s.aliases[childID] = primaryID
	alias := model.IncidentAlias{ChildID: childID, PrimaryID: primaryID, MergedAt: at}
	s.record("incident.merged", primaryID, alias)
	return MergeResult{Primary: primary, Alias: alias}, nil
}
