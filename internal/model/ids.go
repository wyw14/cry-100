package model

import "github.com/google/uuid"

type ObservationID string
type IncidentID string
type MissionID string
type RevisionID string
type LeaseID string
type MessageID string
type OperationID string
type SensorID string

func NewObservationID() ObservationID { return ObservationID(uuid.NewString()) }
func NewIncidentID() IncidentID       { return IncidentID(uuid.NewString()) }
func NewRevisionID() RevisionID       { return RevisionID(uuid.NewString()) }
func NewLeaseID() LeaseID             { return LeaseID(uuid.NewString()) }
func NewMessageID() MessageID         { return MessageID(uuid.NewString()) }
func NewOperationID() OperationID     { return OperationID(uuid.NewString()) }

func StableMissionID(operation OperationID) MissionID {
	return MissionID(uuid.NewSHA1(uuid.NameSpaceOID, []byte(operation)).String())
}
