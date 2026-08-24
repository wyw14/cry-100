package perimeter

import "github.com/wyw14/cry-100/internal/model"

func Avoids(revision model.PerimeterRevision, corridor string) bool {
	for _, zone := range revision.NoFlyZones {
		if zone == corridor {
			return false
		}
	}
	return true
}
