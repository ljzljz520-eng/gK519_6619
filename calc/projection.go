package calc

import (
	"fmt"

	"bridge-trajectory/domain"
)

func Project(points []domain.TrajectoryPoint, mode domain.ViewMode) (domain.Projection, error) {
	if mode != domain.ViewTop && mode != domain.ViewSide {
		return domain.Projection{}, fmt.Errorf("unsupported projection %q", mode)
	}
	projected := make([]domain.ProjectedPoint, 0, len(points))
	for index, point := range points {
		if mode == domain.ViewTop {
			projected = append(projected, domain.ProjectedPoint{Label: fmt.Sprintf("p%d", index), A: point.X, B: point.Y})
		} else {
			projected = append(projected, domain.ProjectedPoint{Label: fmt.Sprintf("p%d", index), A: point.X, B: point.Z})
		}
	}
	return domain.Projection{Mode: mode, Points: projected}, nil
}

func Bounds(projection domain.Projection) (float64, float64, float64, float64) {
	if len(projection.Points) == 0 {
		return 0, 0, 0, 0
	}
	minA, maxA := projection.Points[0].A, projection.Points[0].A
	minB, maxB := projection.Points[0].B, projection.Points[0].B
	for _, point := range projection.Points[1:] {
		if point.A < minA {
			minA = point.A
		}
		if point.A > maxA {
			maxA = point.A
		}
		if point.B < minB {
			minB = point.B
		}
		if point.B > maxB {
			maxB = point.B
		}
	}
	return minA, maxA, minB, maxB
}
