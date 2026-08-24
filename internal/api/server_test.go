package api_test

import (
	"bytes"
	"encoding/json"
	"github.com/wyw14/cry-100/internal/api"
	"github.com/wyw14/cry-100/internal/app"
	"github.com/wyw14/cry-100/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPublicRoutesRespond(t *testing.T) {
	runtime, err := app.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(api.NewServer(runtime).Router())
	defer server.Close()
	for _, path := range []string{"/healthz", "/api/observations", "/api/incidents", "/api/dispatch", "/api/resources", "/api/alerts"} {
		response, err := http.Get(server.URL + path)
		if err != nil {
			t.Fatalf("%s: %v", path, err)
		}
		response.Body.Close()
		if response.StatusCode != http.StatusOK {
			t.Fatalf("%s status=%d", path, response.StatusCode)
		}
	}
}

func TestIncidentAndDispatchWritePaths(t *testing.T) {
	runtime, err := app.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(api.NewServer(runtime).Router())
	defer server.Close()
	response, err := http.Post(server.URL+"/api/incidents", "application/json", bytes.NewBufferString(`{"sector_id":"north-ridge"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusCreated {
		t.Fatalf("incident status=%d", response.StatusCode)
	}
	var incident model.Incident
	if err := json.NewDecoder(response.Body).Decode(&incident); err != nil {
		t.Fatal(err)
	}
	request := model.DispatchRequest{OperationID: model.NewOperationID(), IncidentID: incident.ID, AssigneeID: "crew-a", Kind: "ground", RallyPoint: "north gate"}
	payload, _ := json.Marshal(request)
	dispatchResponse, err := http.Post(server.URL+"/api/dispatch", "application/json", bytes.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}
	defer dispatchResponse.Body.Close()
	if dispatchResponse.StatusCode != http.StatusAccepted {
		t.Fatalf("dispatch status=%d", dispatchResponse.StatusCode)
	}
	if len(runtime.Dispatch.List()) != 1 {
		t.Fatalf("missions=%d", len(runtime.Dispatch.List()))
	}
}
