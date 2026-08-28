package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"bridge-trajectory/domain"
)

type bridgeRequest struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Length     float64 `json:"length"`
	DeckHeight float64 `json:"deck_height"`
}

type scenarioRequest struct {
	ID          string  `json:"id"`
	BridgeID    string  `json:"bridge_id"`
	WindSpeed   float64 `json:"wind_speed"`
	Amplitude   float64 `json:"amplitude"`
	Step        float64 `json:"step"`
	Duration    float64 `json:"duration"`
	Description string  `json:"description"`
}

type trajectoryRequest struct {
	BridgeID   string          `json:"bridge_id"`
	ScenarioID string          `json:"scenario_id"`
	WindSpeed  float64         `json:"wind_speed"`
	Amplitude  float64         `json:"amplitude"`
	Step       float64         `json:"step"`
	Duration   float64         `json:"duration"`
	View       domain.ViewMode `json:"view"`
}

func (s *Server) handleHealth(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		methodAllowed(writer, http.MethodGet)
		return
	}
	writeJSON(writer, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleBridges(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		term := request.URL.Query().Get("q")
		bridges, err := s.catalog.SearchBridges(term)
		if err != nil {
			writeError(writer, http.StatusInternalServerError, err)
			return
		}
		writeJSON(writer, http.StatusOK, bridges)
	case http.MethodPost:
		var input bridgeRequest
		if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
			writeError(writer, http.StatusBadRequest, fmt.Errorf("invalid bridge: %w", err))
			return
		}
		bridge, err := s.bridge.RegisterBridge(input.ID, input.Name, input.Length, input.DeckHeight)
		if err != nil {
			writeError(writer, http.StatusBadRequest, err)
			return
		}
		writeJSON(writer, http.StatusCreated, bridge)
	default:
		methodAllowed(writer, http.MethodGet, http.MethodPost)
	}
}

func (s *Server) handleScenarios(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		items, err := s.catalog.ScenariosFor(request.URL.Query().Get("bridge_id"))
		if err != nil {
			writeError(writer, http.StatusInternalServerError, err)
			return
		}
		writeJSON(writer, http.StatusOK, items)
	case http.MethodPost:
		var input scenarioRequest
		if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
			writeError(writer, http.StatusBadRequest, err)
			return
		}
		item, err := s.bridge.SaveScenario(input.ID, input.BridgeID, input.WindSpeed, input.Amplitude, input.Step, input.Duration, input.Description)
		if err != nil {
			writeError(writer, http.StatusBadRequest, err)
			return
		}
		writeJSON(writer, http.StatusCreated, item)
	default:
		methodAllowed(writer, http.MethodGet, http.MethodPost)
	}
}

func (s *Server) handleTrajectories(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		filter := domain.TrajectoryFilter{BridgeID: request.URL.Query().Get("bridge_id"), ScenarioID: request.URL.Query().Get("scenario_id"), Status: request.URL.Query().Get("status"), Limit: 20}
		items, err := s.trajectory.List(filter)
		if err != nil {
			writeError(writer, http.StatusInternalServerError, err)
			return
		}
		writeJSON(writer, http.StatusOK, items)
	case http.MethodPost:
		var input trajectoryRequest
		if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
			writeError(writer, http.StatusBadRequest, err)
			return
		}
		item, err := s.trajectory.Calculate(domain.CalculationRequest{BridgeID: input.BridgeID, ScenarioID: input.ScenarioID, WindSpeed: input.WindSpeed, Amplitude: input.Amplitude, Step: input.Step, Duration: input.Duration, View: input.View})
		if err != nil {
			writeError(writer, http.StatusBadRequest, err)
			return
		}
		writeJSON(writer, http.StatusCreated, item)
	default:
		methodAllowed(writer, http.MethodGet, http.MethodPost)
	}
}

func joinMethods(methods []string) string {
	result := ""
	for index, method := range methods {
		if index > 0 {
			result += ", "
		}
		result += method
	}
	return result
}
