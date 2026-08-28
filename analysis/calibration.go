package analysis

import (
	"fmt"
	"math"

	"bridge-trajectory/domain"
)

type Calibration struct {
	LateralBias float64 `json:"lateral_bias"`
	HeightBias  float64 `json:"height_bias"`
	Scale       float64 `json:"scale"`
}

func NewCalibration(lateralBias, heightBias, scale float64) (Calibration, error) {
	if math.IsNaN(lateralBias) || math.IsNaN(heightBias) || math.IsNaN(scale) {
		return Calibration{}, fmt.Errorf("calibration values must be finite")
	}
	if scale <= 0 {
		return Calibration{}, fmt.Errorf("calibration scale must be positive")
	}
	return Calibration{LateralBias: lateralBias, HeightBias: heightBias, Scale: scale}, nil
}

func (c Calibration) Apply(points []domain.TrajectoryPoint) []domain.TrajectoryPoint {
	result := make([]domain.TrajectoryPoint, len(points))
	for index, point := range points {
		point.Y = (point.Y - c.LateralBias) * c.Scale
		point.Z = (point.Z - c.HeightBias) * c.Scale
		result[index] = point
	}
	return result
}

func (c Calibration) Reverse(points []domain.TrajectoryPoint) []domain.TrajectoryPoint {
	if c.Scale == 0 {
		return append([]domain.TrajectoryPoint(nil), points...)
	}
	result := make([]domain.TrajectoryPoint, len(points))
	for index, point := range points {
		point.Y = point.Y/c.Scale + c.LateralBias
		point.Z = point.Z/c.Scale + c.HeightBias
		result[index] = point
	}
	return result
}

func (c Calibration) Equal(other Calibration) bool {
	return math.Abs(c.LateralBias-other.LateralBias) < 1e-9 && math.Abs(c.HeightBias-other.HeightBias) < 1e-9 && math.Abs(c.Scale-other.Scale) < 1e-9
}

func FitCalibration(reference, measured []domain.TrajectoryPoint) Calibration {
	count := len(reference)
	if len(measured) < count {
		count = len(measured)
	}
	if count == 0 {
		return Calibration{Scale: 1}
	}
	lateralBias, heightBias := 0.0, 0.0
	for index := 0; index < count; index++ {
		lateralBias += measured[index].Y - reference[index].Y
		heightBias += measured[index].Z - reference[index].Z
	}
	return Calibration{LateralBias: lateralBias / float64(count), HeightBias: heightBias / float64(count), Scale: 1}
}
