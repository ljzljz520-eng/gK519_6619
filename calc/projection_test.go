package calc

import (
	"testing"

	"bridge-trajectory/domain"
)

func TestProjectionSwitchesBetweenViews(t *testing.T) {
	points := []domain.TrajectoryPoint{{Time: 0, X: 4, Y: 5, Z: 6}}
	top, err := Project(points, domain.ViewTop)
	if err != nil || top.Points[0].B != 5 {
		t.Fatalf("top projection failed: %v", err)
	}
	side, err := Project(points, domain.ViewSide)
	if err != nil || side.Points[0].B != 6 {
		t.Fatalf("side projection failed: %v", err)
	}
	if _, err := Project(points, "diagonal"); err == nil {
		t.Fatal("expected invalid view")
	}
}
