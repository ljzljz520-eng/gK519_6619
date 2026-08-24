package analysis

import (
	"fmt"
	"math"
	"sort"

	"bridge-trajectory/domain"
)

type Comparison struct {
	FirstID       string  `json:"first_id"`
	SecondID      string  `json:"second_id"`
	LateralDelta  float64 `json:"lateral_delta"`
	HeightDelta   float64 `json:"height_delta"`
	DistanceDelta float64 `json:"distance_delta"`
	Trend         string  `json:"trend"`
}

func Compare(first, second Report) Comparison {
	trend := "stable"
	if second.RiskScore > first.RiskScore+1 {
		trend = "increasing"
	} else if second.RiskScore < first.RiskScore-1 {
		trend = "decreasing"
	}
	return Comparison{FirstID: first.TrajectoryID, SecondID: second.TrajectoryID, LateralDelta: second.MaxLateral - first.MaxLateral, HeightDelta: second.MaxVertical - first.MaxVertical, DistanceDelta: second.Distance - first.Distance, Trend: trend}
}

func Interpolate(points []domain.TrajectoryPoint, target float64) (domain.TrajectoryPoint, error) {
	if len(points) == 0 {
		return domain.TrajectoryPoint{}, fmt.Errorf("cannot interpolate empty trajectory")
	}
	if target <= points[0].Time {
		return points[0], nil
	}
	for index := 1; index < len(points); index++ {
		if target <= points[index].Time {
			previous := points[index-1]
			current := points[index]
			fraction := (target - previous.Time) / (current.Time - previous.Time)
			return domain.TrajectoryPoint{Time: target, X: lerp(previous.X, current.X, fraction), Y: lerp(previous.Y, current.Y, fraction), Z: lerp(previous.Z, current.Z, fraction)}, nil
		}
	}
	return points[len(points)-1], nil
}

func Smooth(points []domain.TrajectoryPoint, radius int) []domain.TrajectoryPoint {
	if radius < 1 {
		return append([]domain.TrajectoryPoint(nil), points...)
	}
	result := make([]domain.TrajectoryPoint, len(points))
	for index, point := range points {
		start := index - radius
		if start < 0 {
			start = 0
		}
		end := index + radius
		if end >= len(points) {
			end = len(points) - 1
		}
		count := float64(end - start + 1)
		result[index] = average(points[start:end+1], point.Time, count)
	}
	return result
}

func Rank(reports []Report) []Report {
	result := append([]Report(nil), reports...)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].RiskScore == result[j].RiskScore {
			return result[i].TrajectoryID < result[j].TrajectoryID
		}
		return result[i].RiskScore > result[j].RiskScore

	})
	return result
}

func RiskBand(score float64) string {
	score = math.Max(0, math.Min(score, 100))
	if score >= 75 {
		return "critical"
	}
	if score >= 45 {
		return "watch"
	}
	return "normal"
}

func Segmentize(points []domain.TrajectoryPoint) []Segment {
	segments := make([]Segment, 0, len(points)-1)
	for index := 1; index < len(points); index++ {
		segments = append(segments, Segment{Start: points[index-1], End: points[index], Span: distance(points[index-1], points[index])})
	}
	return segments
}

func lerp(first, second, fraction float64) float64 { return first + (second-first)*fraction }

func average(points []domain.TrajectoryPoint, timeValue, count float64) domain.TrajectoryPoint {
	result := domain.TrajectoryPoint{Time: timeValue}
	for _, point := range points {
		result.X += point.X / count
		result.Y += point.Y / count
		result.Z += point.Z / count
	}
	return result
}
