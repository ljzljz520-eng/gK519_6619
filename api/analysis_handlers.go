package api

import (
	"encoding/json"
	"net/http"

	"bridge-trajectory/domain"
	"bridge-trajectory/service"
)

type viewRequest struct {
	UserID string          `json:"user_id"`
	Mode   domain.ViewMode `json:"mode"`
}

type eventRequest struct {
	ID      string `json:"id"`
	Kind    string `json:"kind"`
	Subject string `json:"subject"`
	Detail  string `json:"detail"`
	AtUnix  int64  `json:"at_unix"`
}

func (s *Server) handleAnalysis(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		methodAllowed(writer, http.MethodGet)
		return
	}
	id := request.URL.Query().Get("trajectory_id")
	if id == "" {
		writeError(writer, http.StatusBadRequest, errMissing("trajectory_id"))
		return
	}
	report, err := service.NewAnalysisService(s.database).Analyze(id)
	if err != nil {
		writeError(writer, http.StatusNotFound, err)
		return
	}
	writeJSON(writer, http.StatusOK, report)
}

func (s *Server) handleViews(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		userID := request.URL.Query().Get("user_id")
		mode, err := s.bridge.GetView(userID)
		if err != nil {
			writeError(writer, http.StatusBadRequest, err)
			return
		}
		writeJSON(writer, http.StatusOK, map[string]any{"user_id": userID, "mode": mode})
	case http.MethodPost:
		var input viewRequest
		if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
			writeError(writer, http.StatusBadRequest, err)
			return
		}
		preference, err := s.bridge.SetView(input.UserID, input.Mode)
		if err != nil {
			writeError(writer, http.StatusBadRequest, err)
			return
		}
		writeJSON(writer, http.StatusCreated, preference)
	default:
		methodAllowed(writer, http.MethodGet, http.MethodPost)
	}
}

func (s *Server) handleEvents(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		items, err := s.catalog.EventsFor(request.URL.Query().Get("subject"))
		if err != nil {
			writeError(writer, http.StatusInternalServerError, err)
			return
		}
		writeJSON(writer, http.StatusOK, items)
	case http.MethodPost:
		var input eventRequest
		if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
			writeError(writer, http.StatusBadRequest, err)
			return
		}
		if err := s.catalog.StoreEvent(domain.AuditEvent{ID: input.ID, Kind: input.Kind, Subject: input.Subject, Detail: input.Detail, AtUnix: input.AtUnix}); err != nil {
			writeError(writer, http.StatusBadRequest, err)
			return
		}
		writeJSON(writer, http.StatusCreated, input)
	default:
		methodAllowed(writer, http.MethodGet, http.MethodPost)
	}
}

func errMissing(name string) error { return &missingField{name: name} }

type missingField struct{ name string }

func (e *missingField) Error() string { return e.name + " is required" }
