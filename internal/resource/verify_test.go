package resource_test

import (
	"bytes"
	"encoding/json"
	"github.com/wyw14/cry-100/internal/api"
	"github.com/wyw14/cry-100/internal/app"
	"github.com/wyw14/cry-100/internal/model"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestConcurrentMissionsKeepWaterCorridorExclusive(t *testing.T) {
	runtime, err := app.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	fire := runtime.Incidents.Create("north-ridge", time.Now().UTC())
	server := httptest.NewServer(api.NewServer(runtime).Router())
	defer server.Close()
	start := make(chan struct{})
	var wg sync.WaitGroup
	statuses := make(chan int, 2)
	for index := 0; index < 2; index++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			<-start
			request := model.DispatchRequest{OperationID: model.NewOperationID(), IncidentID: fire.ID, AssigneeID: "heli-" + string(rune('a'+index)), Kind: "ground", RallyPoint: "ridge", CorridorID: "ridge-water"}
			data, _ := json.Marshal(request)
			response, err := http.Post(server.URL+"/api/dispatch", "application/json", bytes.NewReader(data))
			if err != nil {
				statuses <- 0
				return
			}
			defer response.Body.Close()
			statuses <- response.StatusCode
		}(index)
	}
	close(start)
	wg.Wait()
	close(statuses)
	accepted := 0
	for status := range statuses {
		if status == http.StatusAccepted {
			accepted++
		}
	}
	if accepted != 1 {
		t.Fatalf("accepted requests=%d", accepted)
	}
	if runtime.Resources.ActiveOwner("ridge-water") != 1 {
		t.Fatalf("owners=%d", runtime.Resources.ActiveOwner("ridge-water"))
	}
}
