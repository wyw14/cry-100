package incident

import (
	"fmt"
	"github.com/wyw14/cry-100/internal/alert"
	"github.com/wyw14/cry-100/internal/journal"
	"github.com/wyw14/cry-100/internal/model"
	"time"
)

func Contain(s *Service, alerts *alert.Publisher, evidence *journal.EvidenceStore, proof model.FieldProof, image []byte, now time.Time) error {
	if evidence == nil {
		return fmt.Errorf("evidence store is required")
	}
	if err := evidence.Save(proof, image); err != nil {
		return err
	}
	if !evidence.Has(proof) {
		return fmt.Errorf("field proof is not durable")
	}
	if alerts != nil {
		_ = alerts.Withdraw(proof.IncidentID)
	}
	return s.Transition(proof.IncidentID, model.StageContained, now)
}
