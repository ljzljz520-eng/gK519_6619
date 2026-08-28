package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"bridge-trajectory/query"
	"bridge-trajectory/service"
	"bridge-trajectory/store"
)

type Server struct {
	database   *store.Store
	bridge     *service.BridgeService
	trajectory *service.TrajectoryService
	catalog    *query.Catalog
	http       *http.Server
}

func NewServer(database *store.Store) *Server {
	server := &Server{database: database, bridge: service.NewBridgeService(database), trajectory: service.NewTrajectoryService(database), catalog: query.NewCatalog(database)}
	server.http = &http.Server{Handler: server.routes(), ReadHeaderTimeout: 2 * time.Second}
	return server
}

func (s *Server) Handler() http.Handler { return s.http.Handler }

func (s *Server) ListenAndServe(address string) error {
	if address == "" {
		address = ":8080"
	}
	s.http.Addr = address
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error { return s.http.Shutdown(ctx) }

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/bridges", s.handleBridges)
	mux.HandleFunc("/scenarios", s.handleScenarios)
	mux.HandleFunc("/trajectories", s.handleTrajectories)
	mux.HandleFunc("/analysis", s.handleAnalysis)
	mux.HandleFunc("/views", s.handleViews)
	mux.HandleFunc("/events", s.handleEvents)
	return requestLog(mux)
}

func requestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodOptions {
			writer.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(writer, request)
	})
}

func writeJSON(writer http.ResponseWriter, status int, value any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	if err := json.NewEncoder(writer).Encode(value); err != nil {
		return
	}
}

func writeError(writer http.ResponseWriter, status int, err error) {
	writeJSON(writer, status, map[string]string{"error": err.Error()})
}

func methodAllowed(writer http.ResponseWriter, methods ...string) {
	writer.Header().Set("Allow", joinMethods(methods))
	writeError(writer, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
}
