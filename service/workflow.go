package service

import (
	"fmt"
	"time"

	"bridge-trajectory/calc"
	"bridge-trajectory/domain"
	"bridge-trajectory/store"
)

type Workflow struct {
	Bridge     domain.BridgeRecord
	Scenario   domain.WindScenarioRecord
	Trajectory domain.TrajectoryRecord
	View       domain.ViewMode
}

func RunCaptureWorkflow(database *store.Store, bridgeID, scenarioID string) (Workflow, error) {
	bridgeService := NewBridgeService(database)
	trajectoryService := NewTrajectoryService(database)
	bridge, err := bridgeService.RegisterBridge(bridgeID, "Harbor Span", 420, 38)
	if err != nil {
		return Workflow{}, err
	}
	scenario, err := bridgeService.SaveScenario(scenarioID, bridge.ID, 22, 1.8, 0.5, 5, "inspection wind")
	if err != nil {
		return Workflow{}, err
	}
	trajectory, err := trajectoryService.Calculate(domain.CalculationRequest{BridgeID: bridge.ID, ScenarioID: scenario.ID, WindSpeed: scenario.WindSpeed, Amplitude: scenario.Amplitude, Step: scenario.Step, Duration: scenario.Duration, View: domain.ViewTop})
	if err != nil {
		return Workflow{}, err
	}
	return Workflow{Bridge: bridge, Scenario: scenario, Trajectory: trajectory, View: domain.ViewTop}, nil
}

func RunQueryWorkflow(database *store.Store, bridgeID string, mode domain.ViewMode) ([]domain.ProjectedPoint, error) {
	trajectoryService := NewTrajectoryService(database)
	items, err := trajectoryService.List(domain.TrajectoryFilter{BridgeID: bridgeID, Limit: 20})
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("no trajectory for bridge %q", bridgeID)
	}
	projection, err := projectLatest(items, mode)
	if err != nil {
		return nil, err
	}
	return projection.Points, nil
}

func RunRecoveryWorkflow(path string, trajectoryID string) (domain.TrajectoryRecord, error) {
	database, err := store.Open(path)
	if err != nil {
		return domain.TrajectoryRecord{}, err
	}
	if err := database.Close(); err != nil {
		return domain.TrajectoryRecord{}, err
	}
	database, err = store.Open(path)
	if err != nil {
		return domain.TrajectoryRecord{}, err
	}
	defer database.Close()
	service := NewTrajectoryService(database)
	record, err := service.Get(trajectoryID)
	if err != nil {
		return domain.TrajectoryRecord{}, err
	}
	if record.CreatedUnix == 0 {
		record.CreatedUnix = time.Unix(0, 0).Unix()
	}
	return record, nil
}

func projectLatest(items []domain.TrajectoryRecord, mode domain.ViewMode) (domain.Projection, error) {
	latest := items[0]
	for _, item := range items[1:] {
		if item.CreatedUnix > latest.CreatedUnix {
			latest = item
		}
	}
	return project(latest, mode)
}

func project(record domain.TrajectoryRecord, mode domain.ViewMode) (domain.Projection, error) {
	if mode == "" {
		mode = record.View
	}
	return calc.Project(record.Points, mode)
}
