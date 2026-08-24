package sensor

import (
	"fmt"

	"github.com/wyw14/cry-100/internal/model"
)

func activateProfile(profile model.SensorProfile, mode string, calibration *model.Calibration) (model.SensorProfile, error) {
	if calibration == nil || calibration.Mode != mode {
		return model.SensorProfile{}, fmt.Errorf("calibration for mode is not ready")
	}
	profile.Mode = mode
	profile.Epoch++
	profile.Revision = calibration.Revision
	profile.Calibration = calibration
	return profile, nil
}
