package analysis

import (
	"fmt"
	"math"

	"bridge-trajectory/domain"
)

type Report struct {
	TrajectoryID     string     `json:"trajectory_id"`
	PointCount       int        `json:"point_count"`
	Duration         float64    `json:"duration"`
	Distance         float64    `json:"distance"`
	MaxSpeed         float64    `json:"max_speed"`
	MaxLateral       float64    `json:"max_lateral"`
	MaxVertical      float64    `json:"max_vertical"`
	Energy           float64    `json:"energy"`
	RiskScore        float64    `json:"risk_score"`
	Peaks            []Peak     `json:"peaks"`
	Warnings         []string   `json:"warnings"`
	Statistics       Statistics `json:"statistics"`
	Motion           string     `json:"motion"`
	MedianLateral    float64    `json:"median_lateral"`
	NinetiethLateral float64    `json:"ninetieth_lateral"`
	Trend            string     `json:"trend"`
}

type Peak struct {
	Index int     `json:"index"`
	Time  float64 `json:"time"`
	Value float64 `json:"value"`
}

type Segment struct {
	Start domain.TrajectoryPoint
	End   domain.TrajectoryPoint
	Span  float64
}

func Analyze(record domain.TrajectoryRecord) (Report, error) {
	if err := record.Validate(); err != nil {
		return Report{}, err
	}
	if !ValidateContinuity(record.Points) {
		return Report{}, fmt.Errorf("trajectory time is not monotonic")
	}
	report := Report{TrajectoryID: record.ID, PointCount: len(record.Points), Duration: record.Points[len(record.Points)-1].Time}
	report.Distance = PathLength(record.Points)
	report.MaxSpeed = MaxVelocity(record.Points)
	report.MaxLateral = maxAbsY(record.Points)
	report.MaxVertical = maxAbsZ(record.Points)
	report.Energy = MotionEnergy(record.Points)
	report.Peaks = DetectPeaks(record.Points, 0.01)
	report.Warnings = BuildWarnings(report)
	report.RiskScore = Risk(report)
	report.Statistics = Summarize(record.Points)
	report.Motion = ClassifyMotion(report.Statistics)
	report.MedianLateral, report.NinetiethLateral, _ = LateralPercentiles(record.Points)
	report.Trend = Trend(record.Points)
	return report, nil
}

func PathLength(points []domain.TrajectoryPoint) float64 {
	total := 0.0
	for index := 1; index < len(points); index++ {
		dx := points[index].X - points[index-1].X
		dy := points[index].Y - points[index-1].Y
		dz := points[index].Z - points[index-1].Z
		total += math.Sqrt(dx*dx + dy*dy + dz*dz)
	}
	return total
}

func MaxVelocity(points []domain.TrajectoryPoint) float64 {
	maxVelocity := 0.0
	for index := 1; index < len(points); index++ {
		delta := points[index].Time - points[index-1].Time
		if delta <= 0 {
			continue
		}
		dx := points[index].X - points[index-1].X
		dy := points[index].Y - points[index-1].Y
		dz := points[index].Z - points[index-1].Z
		velocity := math.Sqrt(dx*dx+dy*dy+dz*dz) / delta
		if velocity > maxVelocity {
			maxVelocity = velocity
		}
	}
	return maxVelocity
}

func MotionEnergy(points []domain.TrajectoryPoint) float64 {
	energy := 0.0
	for index := 1; index < len(points); index++ {
		delta := points[index].Time - points[index-1].Time
		if delta <= 0 {
			continue
		}
		velocity := distance(points[index-1], points[index]) / delta
		energy += 0.5 * velocity * velocity * delta
	}
	return energy
}

func DetectPeaks(points []domain.TrajectoryPoint, threshold float64) []Peak {
	if threshold < 0 {
		threshold = 0
	}
	peaks := make([]Peak, 0)
	for index := 1; index < len(points)-1; index++ {
		value := math.Abs(points[index].Y)
		before := math.Abs(points[index-1].Y)
		after := math.Abs(points[index+1].Y)
		if value >= threshold && value >= before && value >= after && (value > before || value > after) {
			peaks = append(peaks, Peak{Index: index, Time: points[index].Time, Value: points[index].Y})
		}
	}
	return peaks
}

func BuildWarnings(report Report) []string {
	warnings := make([]string, 0, 3)
	if report.MaxLateral > 3 {
		warnings = append(warnings, "lateral displacement exceeds inspection threshold")
	}
	if report.MaxVertical > 50 {
		warnings = append(warnings, "vertical displacement exceeds deck envelope")
	}
	if report.MaxSpeed > 80 {
		warnings = append(warnings, "motion speed requires structural review")
	}
	if report.PointCount < 3 {
		warnings = append(warnings, "sampling resolution is low")
	}
	return warnings
}

func Risk(report Report) float64 {
	score := 0.0
	score += math.Min(report.MaxLateral/5, 1) * 35
	score += math.Min(report.MaxVertical/10, 1) * 25
	score += math.Min(report.MaxSpeed/100, 1) * 25
	score += math.Min(float64(len(report.Peaks))/10, 1) * 15
	if score > 100 {
		return 100
	}
	return score
}

func ValidateContinuity(points []domain.TrajectoryPoint) bool {
	if len(points) < 1 {
		return false
	}
	for index := 1; index < len(points); index++ {
		if points[index].Time <= points[index-1].Time {
			return false
		}
	}
	return true
}

func maxAbsY(points []domain.TrajectoryPoint) float64 {
	maxValue := 0.0
	for _, point := range points {
		if value := math.Abs(point.Y); value > maxValue {
			maxValue = value
		}
	}
	return maxValue
}

func maxAbsZ(points []domain.TrajectoryPoint) float64 {
	maxValue := 0.0
	for _, point := range points {
		if value := math.Abs(point.Z); value > maxValue {
			maxValue = value
		}
	}
	return maxValue
}

func distance(first, second domain.TrajectoryPoint) float64 {
	dx := second.X - first.X
	dy := second.Y - first.Y
	dz := second.Z - first.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}
