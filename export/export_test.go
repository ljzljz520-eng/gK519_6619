package export

import (
	"strings"
	"testing"

	"bridge-trajectory/analysis"
	"bridge-trajectory/domain"
)

func TestTrajectoryExports(t *testing.T) {
	record := domain.TrajectoryRecord{ID: "t1", BridgeID: "b", ScenarioID: "s", View: domain.ViewSide, Status: "ready", Points: []domain.TrajectoryPoint{{Time: 0, Z: 10}, {Time: 1, Z: 11}}}
	csv, err := TrajectoryCSV(record)
	if err != nil || !strings.Contains(csv, "trajectory_id") || !strings.Contains(csv, "t1") {
		t.Fatalf("csv export failed: %v", err)
	}
	report := analysis.Report{TrajectoryID: record.ID, PointCount: 2, RiskScore: 20}
	markdown := Markdown(record, report)
	if !strings.Contains(markdown, "Risk score") {
		t.Fatal("markdown export failed")
	}
}
