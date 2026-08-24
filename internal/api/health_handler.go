package api

import "net/http"

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "service": "fireline", "protocol_steps": s.Runtime.ProtocolSteps, "catalog_entries": s.Runtime.CatalogEntries})
}
