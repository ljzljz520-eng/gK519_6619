package main

import (
	"flag"
	"fmt"
	"log"

	"bridge-trajectory/domain"
	"bridge-trajectory/export"
	"bridge-trajectory/render"
	"bridge-trajectory/service"
	"bridge-trajectory/store"
)

func main() {
	path := flag.String("db", "bridge-trajectory.db", "database path")
	view := flag.String("view", string(domain.ViewTop), "top or side")
	flag.Parse()
	database, err := store.Open(*path)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	if err := runDemo(database, domain.ViewMode(*view)); err != nil {
		log.Fatal(err)
	}
}

func runDemo(database *store.Store, mode domain.ViewMode) error {
	bridgeService := service.NewBridgeService(database)
	bridge, err := bridgeService.RegisterBridge("demo-bridge", "Demo Bridge", 300, 25)
	if err != nil {
		if existing, getErr := database.GetBridge("demo-bridge"); getErr == nil {
			bridge = existing
		} else {
			return err
		}
	}
	scenario, err := bridgeService.SaveScenario("demo-wind", bridge.ID, 18, 1.2, 0.5, 4, "command line sample")
	if err != nil {
		if existing, getErr := database.GetScenario("demo-wind"); getErr == nil {
			scenario = existing
		} else {
			return err
		}
	}
	trajectoryService := service.NewTrajectoryService(database)
	trajectory, err := trajectoryService.Calculate(domain.CalculationRequest{BridgeID: bridge.ID, ScenarioID: scenario.ID, WindSpeed: scenario.WindSpeed, Amplitude: scenario.Amplitude, Step: scenario.Step, Duration: scenario.Duration, View: domain.DefaultView(mode)})
	if err != nil {
		return err
	}
	projection, err := projectForCLI(trajectory)
	if err != nil {
		return err
	}
	fmt.Println(render.Legend(projection.Mode))
	fmt.Println(render.ASCII(projection, 72, 16))
	fmt.Printf("trajectory=%s points=%d\n", trajectory.ID, len(trajectory.Points))
	report, err := service.NewAnalysisService(database).Analyze(trajectory.ID)
	if err != nil {
		return err
	}
	fmt.Println(export.Markdown(trajectory, report))
	return nil
}

func projectForCLI(record domain.TrajectoryRecord) (domain.Projection, error) {
	if record.View == domain.ViewSide {
		return projectSide(record)
	}
	return projectTop(record)
}

func projectTop(record domain.TrajectoryRecord) (domain.Projection, error) {
	points := make([]domain.ProjectedPoint, 0, len(record.Points))
	for index, point := range record.Points {
		points = append(points, domain.ProjectedPoint{Label: fmt.Sprintf("p%d", index), A: point.X, B: point.Y})
	}
	return domain.Projection{Mode: domain.ViewTop, Points: points}, nil
}

func projectSide(record domain.TrajectoryRecord) (domain.Projection, error) {
	points := make([]domain.ProjectedPoint, 0, len(record.Points))
	for index, point := range record.Points {
		points = append(points, domain.ProjectedPoint{Label: fmt.Sprintf("p%d", index), A: point.X, B: point.Z})
	}
	return domain.Projection{Mode: domain.ViewSide, Points: points}, nil
}
