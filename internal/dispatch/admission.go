package dispatch

import (
	"fmt"

	"github.com/wyw14/cry-100/internal/model"
)

func validateAdmission(request model.DispatchRequest, incident model.Incident) error {
	if incident.ID == "" {
		return fmt.Errorf("incident identity is required")
	}
	if request.AssigneeID == "" || request.RallyPoint == "" {
		return fmt.Errorf("assignee and rally point are required")
	}
	if request.Kind != "ground" && request.Kind != "air" {
		return fmt.Errorf("unsupported mission kind %q", request.Kind)
	}
	return nil
}
