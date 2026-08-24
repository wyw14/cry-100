package journal

import (
	"github.com/wyw14/cry-100/internal/model"
	"time"
)

func AppendValue(store *Store, kind, aggregate string, value any) (model.Event, error) {
	event, err := model.NewEvent(kind, aggregate, value, time.Now().UTC())
	if err != nil {
		return event, err
	}
	return event, store.Append(event)
}
