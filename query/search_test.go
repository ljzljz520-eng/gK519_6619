package query

import (
	"path/filepath"
	"testing"

	"bridge-trajectory/domain"
	"bridge-trajectory/service"
	"bridge-trajectory/store"
)

func TestCatalogSearchAndMetrics(t *testing.T) {
	database, err := store.Open(filepath.Join(t.TempDir(), "catalog.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	workflow, err := service.RunCaptureWorkflow(database, "golden", "wind")
	if err != nil {
		t.Fatal(err)
	}
	catalog := NewCatalog(database)
	bridges, err := catalog.SearchBridges("GOLD")
	if err != nil || len(bridges) != 1 {
		t.Fatalf("bridge search failed: %v", err)
	}
	items, err := catalog.SearchTrajectories(domain.TrajectoryFilter{BridgeID: workflow.Bridge.ID})
	if err != nil || len(items) != 1 {
		t.Fatalf("trajectory search failed: %v", err)
	}
	metrics, err := Measure(items[0])
	if err != nil || metrics.PointCount < 2 {
		t.Fatalf("metrics failed: %v", err)
	}
}
