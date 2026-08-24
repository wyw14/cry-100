package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/wyw14/cry-100/internal/model"
)

func (s *Server) dispatches(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"missions": s.Runtime.Dispatch.List(), "flight_plans": s.Runtime.Aircraft.Plans(), "radio_messages": s.Runtime.Radio.List()})
}

type acknowledgementRequest struct {
	Responder string `json:"responder"`
}

func (s *Server) cancelDispatch(w http.ResponseWriter, r *http.Request) {
	id := model.MissionID(chi.URLParam(r, "missionID"))
	if err := s.Runtime.Dispatch.Cancel(id, time.Now().UTC()); err != nil {
		fail(w, http.StatusConflict, err)
		return
	}
	value, _ := s.Runtime.Dispatch.Get(id)
	respond(w, http.StatusOK, value)
}

func (s *Server) ackDispatch(w http.ResponseWriter, r *http.Request) {
	var request acknowledgementRequest
	if err := decode(r, &request); err != nil {
		fail(w, http.StatusBadRequest, err)
		return
	}
	id := model.MissionID(chi.URLParam(r, "missionID"))
	if err := s.Runtime.Dispatch.Acknowledge(id, request.Responder); err != nil {
		fail(w, http.StatusConflict, err)
		return
	}
	value, _ := s.Runtime.Radio.Get(id)
	respond(w, http.StatusOK, value)
}
func (s *Server) dispatch(w http.ResponseWriter, r *http.Request) {
	var request model.DispatchRequest
	if err := decode(r, &request); err != nil {
		writeJSON(w, 400, map[string]string{"error": err.Error()})
		return
	}
	if request.OperationID == "" {
		request.OperationID = model.NewOperationID()
	}
	result, err := s.Runtime.Dispatch.Plan(request, time.Now().UTC())
	if err != nil {
		writeJSON(w, 409, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusAccepted, result)
}
