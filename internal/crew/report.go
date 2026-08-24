package crew

import (
	"errors"
	"time"

	"github.com/wyw14/cry-100/internal/model"
)

type FieldReport struct {
	ReportID    string
	IncidentID  model.IncidentID
	Summary     string
	ImageName   string
	ImageBytes  int
	SubmittedAt time.Time
}

func BuildFieldReport(proof model.FieldProof, image []byte) FieldReport {
	return FieldReport{
		ReportID: proof.ReportID, IncidentID: proof.IncidentID,
		Summary: proof.Summary, ImageName: proof.ImageName,
		ImageBytes: len(image), SubmittedAt: proof.SubmittedAt,
	}
}

func (r FieldReport) Validate() error {
	if r.ReportID == "" || r.ImageName == "" || r.ImageBytes == 0 {
		return errors.New("field report requires identity and image")
	}
	return nil
}

func (r FieldReport) Proof() model.FieldProof {
	return model.FieldProof{
		IncidentID: r.IncidentID, ReportID: r.ReportID, Summary: r.Summary,
		ImageName: r.ImageName, SubmittedAt: r.SubmittedAt,
	}
}
