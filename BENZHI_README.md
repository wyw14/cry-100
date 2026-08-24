# FireLine

FireLine is a Go HTTP service for forest-fire observations, incident state, dispatch coordination, water-corridor leases, radio acknowledgements, public alerts, and local event recovery.

The service vendors its Go dependencies for offline builds. Build it with `go build -mod=vendor ./...` and start it with `go run -mod=vendor ./cmd/fireline`. It listens on `127.0.0.1:19700` by default; override the address with `FIRELINE_ADDR` and the local event directory with `FIRELINE_DATA_DIR`.

Health is available at `/healthz`. JSON APIs are rooted at `/api/observations`, `/api/incidents`, `/api/dispatch`, `/api/resources`, and `/api/alerts`.
