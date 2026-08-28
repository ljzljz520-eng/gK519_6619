package store

import (
	"path/filepath"
	"testing"

	"bridge-trajectory/domain"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bridge.db")
	database, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	bridge := domain.NewBridge("bridge-1", "River", 500, 42, fixedTime())
	if err := database.PutBridge(bridge); err != nil {
		t.Fatal(err)
	}
	scenario := domain.NewScenario("scenario-1", bridge.ID, 25, 2, 0.5, 4, "gust")
	if err := database.PutScenario(scenario); err != nil {
		t.Fatal(err)
	}
	trajectory := domain.TrajectoryRecord{ID: "trajectory-1", BridgeID: bridge.ID, ScenarioID: scenario.ID, View: domain.ViewTop, Status: "ready", Points: []domain.TrajectoryPoint{{Time: 0, X: 0, Y: 0, Z: 42}}}
	if err := database.PutTrajectory(trajectory); err != nil {
		t.Fatal(err)
	}
	if err := database.PutViewPreference(domain.ViewPreferenceRecord{UserID: "operator", Mode: domain.ViewSide}); err != nil {
		t.Fatal(err)
	}
	if err := database.Close(); err != nil {
		t.Fatal(err)
	}
	database, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	got, err := database.GetTrajectory(trajectory.ID)
	if err != nil || len(got.Points) != 1 || got.Points[0].Z != 42 {
		t.Fatalf("persistence lost: %v", err)
	}
	view, err := database.GetViewPreference("operator")
	if err != nil || view.Mode != domain.ViewSide {
		t.Fatalf("view persistence failed: %v", err)
	}
}

func TestStoreFiltersTrajectoryRecords(t *testing.T) {
	database, err := Open(filepath.Join(t.TempDir(), "filter.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	for index, status := range []string{"ready", "archived"} {
		record := domain.TrajectoryRecord{ID: string(rune('a' + index)), BridgeID: "b", ScenarioID: "s", View: domain.ViewTop, Status: status, Points: []domain.TrajectoryPoint{{Time: float64(index)}}}
		if err := database.PutTrajectory(record); err != nil {
			t.Fatal(err)
		}
	}
	items, err := database.ListTrajectories(domain.TrajectoryFilter{BridgeID: "b", Status: "ready"})
	if err != nil || len(items) != 1 {
		t.Fatalf("filter failed: %v", err)
	}
}
