package service

import (
	"path/filepath"
	"testing"

	"bridge-trajectory/domain"
	"bridge-trajectory/store"
)

func TestWorkflowCaptureCalculatePersist(t *testing.T) {
	database, err := store.Open(filepath.Join(t.TempDir(), "capture.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	workflow, err := RunCaptureWorkflow(database, "bridge-a", "scenario-a")
	if err != nil {
		t.Fatal(err)
	}
	if workflow.Bridge.ID != "bridge-a" || len(workflow.Trajectory.Points) < 2 {
		t.Fatalf("workflow did not produce trajectory")
	}
	loaded, err := database.GetTrajectory(workflow.Trajectory.ID)
	if err != nil || len(loaded.Points) != len(workflow.Trajectory.Points) {
		t.Fatalf("trajectory not persisted: %v", err)
	}
}

func TestWorkflowQueryAndProject(t *testing.T) {
	database, err := store.Open(filepath.Join(t.TempDir(), "query.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	if _, err := RunCaptureWorkflow(database, "bridge-b", "scenario-b"); err != nil {
		t.Fatal(err)
	}
	points, err := RunQueryWorkflow(database, "bridge-b", domain.ViewSide)
	if err != nil || len(points) == 0 {
		t.Fatalf("query workflow failed: %v", err)
	}
	if points[0].B < 30 {
		t.Fatalf("side projection did not use height")
	}
}

func TestWorkflowRecovery(t *testing.T) {
	path := filepath.Join(t.TempDir(), "recovery.db")
	database, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	workflow, err := RunCaptureWorkflow(database, "bridge-c", "scenario-c")
	if err != nil {
		t.Fatal(err)
	}
	if err := database.Close(); err != nil {
		t.Fatal(err)
	}
	recovered, err := RunRecoveryWorkflow(path, workflow.Trajectory.ID)
	if err != nil || len(recovered.Points) != len(workflow.Trajectory.Points) {
		t.Fatalf("recovery failed: %v", err)
	}
}
