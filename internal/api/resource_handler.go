package api

import "net/http"

func (s *Server) resources(w http.ResponseWriter, r *http.Request) {
	corridor, ok := s.Runtime.Resources.Corridor("ridge-water")
	writeJSON(w, http.StatusOK, map[string]any{"corridors": []any{map[string]any{"value": corridor, "present": ok, "active_owners": s.Runtime.Resources.ActiveOwner("ridge-water")}}, "leases": s.Runtime.Resources.Leases(), "crews": s.Runtime.Crews.List(), "sectors": s.Runtime.Grid.List()})
}
