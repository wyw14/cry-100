package alert

import (
	"errors"

	"github.com/wyw14/cry-100/internal/model"
)

type Gateway struct {
	available bool
}

func NewGateway(available ...bool) *Gateway {
	value := true
	if len(available) > 0 {
		value = available[0]
	}
	return &Gateway{available: value}
}
func (g *Gateway) Withdraw(id model.IncidentID) error {
	if !g.available {
		return ErrBrokerUnavailable
	}
	return nil
}

var ErrBrokerUnavailable = errors.New("broker unavailable")
