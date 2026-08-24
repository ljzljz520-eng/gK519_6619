package service

import (
	"path/filepath"
	"testing"

	"bridge-trajectory/store"
)

func TestAnalysisServiceReviewsStoredTrajectory(t *testing.T) {
	database, err := store.Open(filepath.Join(t.TempDir(), "analysis.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	workflow, err := RunCaptureWorkflow(database, "analysis-bridge", "analysis-scenario")
	if err != nil {
		t.Fatal(err)
	}
	report, err := NewAnalysisService(database).Analyze(workflow.Trajectory.ID)
	if err != nil || report.PointCount < 2 {
		t.Fatalf("analysis failed: %v", err)
	}
	band, err := NewAnalysisService(database).Review(workflow.Trajectory.ID)
	if err != nil || band == "" {
		t.Fatalf("review failed: %v", err)
	}
}
