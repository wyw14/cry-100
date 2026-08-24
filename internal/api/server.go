package api

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/wyw14/cry-100/internal/app"
	"net/http"
)

type Server struct{ Runtime *app.Runtime }

func NewServer(runtime *app.Runtime) *Server { return &Server{Runtime: runtime} }
func (s *Server) Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/healthz", s.health)
	r.Route("/api", func(r chi.Router) {
		r.Get("/observations", s.observations)
		r.Post("/observations", s.createObservation)
		r.Get("/incidents", s.incidents)
		r.Post("/incidents", s.createIncident)
		r.Post("/incidents/{incidentID}/close", s.closeIncident)
		r.Post("/incidents/{incidentID}/merge", s.mergeIncident)
		r.Post("/incidents/{incidentID}/contain", s.containIncident)
		r.Get("/perimeters", s.perimeters)
		r.Post("/perimeters", s.publishPerimeter)
		r.Get("/dispatch", s.dispatches)
		r.Post("/dispatch", s.dispatch)
		r.Post("/dispatch/{missionID}/cancel", s.cancelDispatch)
		r.Post("/dispatch/{missionID}/ack", s.ackDispatch)
		r.Get("/resources", s.resources)
		r.Get("/alerts", s.alerts)
		r.Post("/alerts", s.publishAlert)
		r.Post("/alerts/{incidentID}/withdraw", s.withdrawAlert)
		r.Post("/sensors/{sensorID}/mode", s.setSensorMode)
	})
	return r
}
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func decode(r *http.Request, value any) error { return json.NewDecoder(r.Body).Decode(value) }
