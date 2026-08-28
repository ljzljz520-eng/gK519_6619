package query

import (
	"fmt"
	"math"

	"bridge-trajectory/calc"
	"bridge-trajectory/domain"
)

type TrajectoryMetrics struct {
	PointCount int     `json:"point_count"`
	Duration   float64 `json:"duration"`
	MinLateral float64 `json:"min_lateral"`
	MaxLateral float64 `json:"max_lateral"`
	MinHeight  float64 `json:"min_height"`
	MaxHeight  float64 `json:"max_height"`
	Span       float64 `json:"span"`
}

func Measure(record domain.TrajectoryRecord) (TrajectoryMetrics, error) {
	if err := record.Validate(); err != nil {
		return TrajectoryMetrics{}, err
	}
	minY, maxY, minZ, maxZ := (&calc.Sampler{}).SampleEnvelope(record.Points)
	duration := record.Points[len(record.Points)-1].Time
	if duration < 0 || math.IsNaN(duration) {
		return TrajectoryMetrics{}, fmt.Errorf("trajectory duration is invalid")
	}
	return TrajectoryMetrics{PointCount: len(record.Points), Duration: duration, MinLateral: minY, MaxLateral: maxY, MinHeight: minZ, MaxHeight: maxZ, Span: maxY - minY}, nil
}

func Compare(a, b domain.TrajectoryRecord) (float64, error) {
	first, err := Measure(a)
	if err != nil {
		return 0, err
	}
	second, err := Measure(b)
	if err != nil {
		return 0, err
	}
	if first.PointCount == 0 || second.PointCount == 0 {
		return 0, nil
	}
	return math.Abs(first.Span - second.Span), nil
}

func Format(metrics TrajectoryMetrics) string {
	return fmt.Sprintf("points=%d duration=%.2f lateral=[%.3f,%.3f] height=[%.3f,%.3f]", metrics.PointCount, metrics.Duration, metrics.MinLateral, metrics.MaxLateral, metrics.MinHeight, metrics.MaxHeight)
}
