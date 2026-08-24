package service

import (
	"path/filepath"
	"testing"

	"bridge-trajectory/domain"
	"bridge-trajectory/store"
)

func TestBridgeTrajectoriesKeepSnapshots(t *testing.T) {
	database, err := store.Open(filepath.Join(t.TempDir(), "bug.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	bridgeService := NewBridgeService(database)
	bridge, err := bridgeService.RegisterBridge("bridge-bug", "Snapshot Bridge", 240, 22)
	if err != nil {
		t.Fatal(err)
	}
	scenario, err := bridgeService.SaveScenario("scenario-bug", bridge.ID, 30, 2, 0.5, 4, "snapshot")
	if err != nil {
		t.Fatal(err)
	}
	trajectoryService := NewTrajectoryService(database)
	first, err := trajectoryService.Calculate(domain.CalculationRequest{BridgeID: bridge.ID, ScenarioID: scenario.ID, WindSpeed: 30, Amplitude: 2, Step: 0.5, Duration: 4, View: domain.ViewTop})
	if err != nil {
		t.Fatal(err)
	}
	snapshot := append([]domain.TrajectoryPoint(nil), first.Points...)
	if _, err := trajectoryService.Calculate(domain.CalculationRequest{BridgeID: bridge.ID, ScenarioID: scenario.ID, WindSpeed: 30, Amplitude: 2, Step: 1, Duration: 4, View: domain.ViewTop}); err != nil {
		t.Fatal(err)
	}
	if len(first.Points) != len(snapshot) {
		t.Fatalf("first trajectory changed point count")
	}
	for index := range snapshot {
		if first.Points[index] != snapshot[index] {
			t.Fatalf("snapshot point %d was rewritten", index)
		}
	}
}
