package domain

import (
	"testing"
	"time"
)

func TestBridgeAndScenarioValidation(t *testing.T) {
	bridge := NewBridge("b1", "North", 100, 20, time.Unix(10, 0))
	if err := bridge.Validate(); err != nil {
		t.Fatal(err)
	}
	scenario := NewScenario("s1", bridge.ID, 10, 2, 0.5, 4, "steady")
	if err := scenario.Validate(); err != nil {
		t.Fatal(err)
	}
	if err := (BridgeRecord{}).Validate(); err == nil {
		t.Fatal("expected invalid bridge")
	}
}

func TestClonePointsIsIndependent(t *testing.T) {
	original := TrajectoryRecord{ID: "t1", BridgeID: "b1", ScenarioID: "s1", View: ViewTop, Points: []TrajectoryPoint{{X: 1, Y: 2, Z: 3}}}
	clone := original.Clone()
	clone.Points[0].Y = 99
	if original.Points[0].Y == clone.Points[0].Y {
		t.Fatal("clone shares points")
	}
}
