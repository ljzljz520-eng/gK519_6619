package calc

import (
	"fmt"
	"math"

	"bridge-trajectory/domain"
)

type Sampler struct{}

func NewSampler() Sampler { return Sampler{} }

func (Sampler) SampleInto(dst []domain.TrajectoryPoint, length, deckHeight float64, model WindModel, step, duration float64) ([]domain.TrajectoryPoint, error) {
	if length <= 0 || deckHeight <= 0 {
		return nil, fmt.Errorf("bridge geometry must be positive")
	}
	if step <= 0 || duration <= 0 || step > duration {
		return nil, fmt.Errorf("sampling interval is invalid")
	}
	count := int(math.Floor(duration/step)) + 1
	if count < 2 {
		count = 2
	}
	points := dst[:0]
	if cap(points) < count {
		points = make([]domain.TrajectoryPoint, 0, count)
	}
	for index := 0; index < count; index++ {
		t := float64(index) * step
		if t > duration {
			t = duration
		}
		points = append(points, domain.TrajectoryPoint{Time: t, X: length * t / duration, Y: model.Lateral(t), Z: deckHeight + model.Vertical(t) + model.Torsion(t)})
		if t >= duration {
			break
		}
	}
	return points, nil
}

func (Sampler) Sample(length, deckHeight float64, model WindModel, step, duration float64) ([]domain.TrajectoryPoint, error) {
	if length <= 0 || deckHeight <= 0 {
		return nil, fmt.Errorf("bridge geometry must be positive")
	}
	if step <= 0 || duration <= 0 || step > duration {
		return nil, fmt.Errorf("sampling interval is invalid")
	}
	count := int(math.Floor(duration/step)) + 1
	if count < 2 {
		count = 2
	}
	points := make([]domain.TrajectoryPoint, 0, count)
	for index := 0; index < count; index++ {
		t := float64(index) * step
		if t > duration {
			t = duration
		}
		points = append(points, domain.TrajectoryPoint{Time: t, X: length * t / duration, Y: model.Lateral(t), Z: deckHeight + model.Vertical(t) + model.Torsion(t)})
		if t >= duration {
			break
		}
	}
	return points, nil
}

func (Sampler) SampleEnvelope(points []domain.TrajectoryPoint) (minY, maxY, minZ, maxZ float64) {
	if len(points) == 0 {
		return 0, 0, 0, 0
	}
	minY, maxY = points[0].Y, points[0].Y
	minZ, maxZ = points[0].Z, points[0].Z
	for _, point := range points[1:] {
		if point.Y < minY {
			minY = point.Y
		}
		if point.Y > maxY {
			maxY = point.Y
		}
		if point.Z < minZ {
			minZ = point.Z
		}
		if point.Z > maxZ {
			maxZ = point.Z
		}
	}
	return minY, maxY, minZ, maxZ
}

func (Sampler) Resample(points []domain.TrajectoryPoint, scale float64) []domain.TrajectoryPoint {
	if scale <= 0 {
		scale = 1
	}
	result := make([]domain.TrajectoryPoint, len(points))
	for index, point := range points {
		point.Y *= scale
		point.Z = point.Z + (scale-1)*0.1
		result[index] = point
	}
	return result
}
