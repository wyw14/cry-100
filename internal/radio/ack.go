package radio

import (
	"github.com/wyw14/cry-100/internal/model"
	"time"
)

func (p *Publisher) ApplyAck(missionID model.MissionID, responder string) error {
	return p.Ack(missionID, responder, time.Now().UTC())
}
