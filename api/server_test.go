package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"bridge-trajectory/store"
)

func TestServerWorkflowEndpoints(t *testing.T) {
	database, err := store.Open(filepath.Join(t.TempDir(), "api.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	server := NewServer(database)
	bridgeBody := bytes.NewBufferString(`{"id":"api-bridge","name":"API Bridge","length":180,"deck_height":18}`)
	request := httptest.NewRequest(http.MethodPost, "/bridges", bridgeBody)
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("bridge endpoint status %d", response.Code)
	}
	scenarioBody := bytes.NewBufferString(`{"id":"api-scenario","bridge_id":"api-bridge","wind_speed":12,"amplitude":1,"step":1,"duration":3}`)
	request = httptest.NewRequest(http.MethodPost, "/scenarios", scenarioBody)
	response = httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("scenario endpoint status %d", response.Code)
	}
	trajectoryBody := bytes.NewBufferString(`{"bridge_id":"api-bridge","scenario_id":"api-scenario","wind_speed":12,"amplitude":1,"step":1,"duration":3,"view":"side"}`)
	request = httptest.NewRequest(http.MethodPost, "/trajectories", trajectoryBody)
	response = httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("trajectory endpoint status %d", response.Code)
	}
	var payload map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil || payload["status"] != "ready" {
		t.Fatalf("trajectory response invalid: %v", err)
	}
}
