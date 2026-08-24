package store

import (
	"path/filepath"
	"testing"

	"bridge-trajectory/domain"
)

func TestStoreCountsAndDeletion(t *testing.T) {
	database, err := Open(filepath.Join(t.TempDir(), "maintenance.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	if err := database.PutBridge(domain.NewBridge("b", "Bridge", 100, 10, fixedTime())); err != nil {
		t.Fatal(err)
	}
	if err := database.PutTrajectory(domain.TrajectoryRecord{ID: "t", BridgeID: "b", ScenarioID: "s", View: domain.ViewTop, Status: "ready", Points: []domain.TrajectoryPoint{{Time: 0, X: 0, Z: 10}}}); err != nil {
		t.Fatal(err)
	}
	counts, err := database.Counts()
	if err != nil || counts.Bridges != 1 || counts.Trajectories != 1 {
		t.Fatalf("counts failed: %v", err)
	}
	if err := database.VerifyBuckets(); err != nil {
		t.Fatal(err)
	}
	if err := database.DeleteTrajectory("t"); err != nil {
		t.Fatal(err)
	}
	if _, err := database.GetTrajectory("t"); err != ErrNotFound {
		t.Fatalf("expected deleted trajectory")
	}
}
