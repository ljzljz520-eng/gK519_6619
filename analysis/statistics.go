package analysis

import (
	"math"
	"sort"

	"bridge-trajectory/domain"
)

type Statistics struct {
	MeanLateral      float64 `json:"mean_lateral"`
	LateralVariance  float64 `json:"lateral_variance"`
	MeanHeight       float64 `json:"mean_height"`
	HeightVariance   float64 `json:"height_variance"`
	AverageStep      float64 `json:"average_step"`
	MinimumStep      float64 `json:"minimum_step"`
	MaximumStep      float64 `json:"maximum_step"`
	PositiveSwingPct float64 `json:"positive_swing_pct"`
}

func Summarize(points []domain.TrajectoryPoint) Statistics {
	if len(points) == 0 {
		return Statistics{}
	}
	meanY, meanZ := 0.0, 0.0
	positive := 0
	for _, point := range points {
		meanY += point.Y
		meanZ += point.Z
		if point.Y > 0 {
			positive++
		}
	}
	meanY /= float64(len(points))
	meanZ /= float64(len(points))
	varianceY, varianceZ := 0.0, 0.0
	for _, point := range points {
		varianceY += math.Pow(point.Y-meanY, 2)
		varianceZ += math.Pow(point.Z-meanZ, 2)
	}
	varianceY /= float64(len(points))
	varianceZ /= float64(len(points))
	minimum, maximum, total := 0.0, 0.0, 0.0
	if len(points) > 1 {
		minimum = points[1].Time - points[0].Time
		maximum = minimum
		for index := 1; index < len(points); index++ {
			step := points[index].Time - points[index-1].Time
			total += step
			if step < minimum {
				minimum = step
			}
			if step > maximum {
				maximum = step
			}
		}
	}
	averageStep := 0.0
	if len(points) > 1 {
		averageStep = total / float64(len(points)-1)
	}
	return Statistics{MeanLateral: meanY, LateralVariance: varianceY, MeanHeight: meanZ, HeightVariance: varianceZ, AverageStep: averageStep, MinimumStep: minimum, MaximumStep: maximum, PositiveSwingPct: float64(positive) / float64(len(points)) * 100}
}

func ClassifyMotion(stats Statistics) string {
	if stats.LateralVariance > 4 || stats.HeightVariance > 4 {
		return "volatile"
	}
	if stats.LateralVariance > 0.5 || stats.HeightVariance > 0.5 {
		return "oscillating"
	}
	return "steady"
}

func Percentile(values []float64, percentile float64) float64 {
	if len(values) == 0 {
		return 0
	}
	copyOf := append([]float64(nil), values...)
	sort.Float64s(copyOf)
	if percentile <= 0 {
		return copyOf[0]
	}
	if percentile >= 100 {
		return copyOf[len(copyOf)-1]
	}
	position := percentile / 100 * float64(len(copyOf)-1)
	lower := int(math.Floor(position))
	upper := int(math.Ceil(position))
	if lower == upper {
		return copyOf[lower]
	}
	return copyOf[lower] + (copyOf[upper]-copyOf[lower])*(position-float64(lower))
}

func LateralPercentiles(points []domain.TrajectoryPoint) (float64, float64, float64) {
	values := make([]float64, 0, len(points))
	for _, point := range points {
		values = append(values, math.Abs(point.Y))
	}
	return Percentile(values, 50), Percentile(values, 90), Percentile(values, 99)
}

func Trend(points []domain.TrajectoryPoint) string {
	if len(points) < 2 {
		return "insufficient"
	}
	first := points[0].Y
	last := points[len(points)-1].Y
	delta := last - first
	if math.Abs(delta) < 0.01 {
		return "flat"
	}
	if delta > 0 {
		return "rising"
	}
	return "falling"
}

func Window(points []domain.TrajectoryPoint, start, end float64) []domain.TrajectoryPoint {
	if end < start {
		start, end = end, start
	}
	window := make([]domain.TrajectoryPoint, 0)
	for _, point := range points {
		if point.Time >= start && point.Time <= end {
			window = append(window, point)
		}
	}
	return window
}

func Normalize(points []domain.TrajectoryPoint) []domain.TrajectoryPoint {
	if len(points) == 0 {
		return nil
	}
	maxValue := 0.0
	for _, point := range points {
		maxValue = math.Max(maxValue, math.Abs(point.Y))
	}
	if maxValue == 0 {
		return append([]domain.TrajectoryPoint(nil), points...)
	}
	result := make([]domain.TrajectoryPoint, len(points))
	for index, point := range points {
		point.Y /= maxValue
		result[index] = point
	}
	return result
}

func Correlate(first, second []domain.TrajectoryPoint) float64 {
	count := len(first)
	if len(second) < count {
		count = len(second)
	}
	if count < 2 {
		return 0
	}
	meanFirst, meanSecond := 0.0, 0.0
	for index := 0; index < count; index++ {
		meanFirst += first[index].Y
		meanSecond += second[index].Y
	}
	meanFirst /= float64(count)
	meanSecond /= float64(count)
	numerator, firstSum, secondSum := 0.0, 0.0, 0.0
	for index := 0; index < count; index++ {
		left := first[index].Y - meanFirst
		right := second[index].Y - meanSecond
		numerator += left * right
		firstSum += left * left
		secondSum += right * right
	}
	if firstSum == 0 || secondSum == 0 {
		return 0
	}
	return numerator / math.Sqrt(firstSum*secondSum)
}
