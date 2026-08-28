package analysis

import (
	"testing"

	"bridge-trajectory/domain"
)

func TestAnalyzeBuildsWarningsAndPeaks(t *testing.T) {
	record := domain.TrajectoryRecord{ID: "t", BridgeID: "b", ScenarioID: "s", View: domain.ViewTop, Points: []domain.TrajectoryPoint{{Time: 0, X: 0, Y: 0, Z: 1}, {Time: 1, X: 1, Y: 4, Z: 2}, {Time: 2, X: 2, Y: 0, Z: 3}}}
	report, err := Analyze(record)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Peaks) != 1 || len(report.Warnings) == 0 || report.RiskScore <= 0 {
		t.Fatalf("analysis output incomplete")
	}
}

func TestInterpolationAndSmoothing(t *testing.T) {
	points := []domain.TrajectoryPoint{{Time: 0, X: 0}, {Time: 2, X: 2}}
	point, err := Interpolate(points, 1)
	if err != nil || point.X != 1 {
		t.Fatalf("interpolation failed")
	}
	smoothed := Smooth([]domain.TrajectoryPoint{{Time: 0, Y: 0}, {Time: 1, Y: 2}, {Time: 2, Y: 0}}, 1)
	if len(smoothed) != 3 || smoothed[1].Y != 2.0/3.0 {
		t.Fatalf("smoothing failed")
	}
}

func TestCalibrationRoundTrip(t *testing.T) {
	calibration, err := NewCalibration(1, 2, 2)
	if err != nil {
		t.Fatal(err)
	}
	points := []domain.TrajectoryPoint{{Y: 3, Z: 6}}
	corrected := calibration.Apply(points)
	original := calibration.Reverse(corrected)
	if original[0] != points[0] || !calibration.Equal(calibration) {
		t.Fatalf("calibration round trip failed")
	}
}
