package render

import (
	"strings"
	"testing"

	"bridge-trajectory/domain"
)

func TestASCIIAndJSONRender(t *testing.T) {
	projection := domain.Projection{Mode: domain.ViewTop, Points: []domain.ProjectedPoint{{Label: "p0", A: 0, B: 0}, {Label: "p1", A: 1, B: 1}}}
	text := ASCII(projection, 20, 6)
	if !strings.Contains(text, "O") {
		t.Fatal("ascii output has no marker")
	}
	encoded, err := ProjectionJSON(projection)
	if err != nil || !strings.Contains(encoded, "p0") {
		t.Fatalf("json output failed: %v", err)
	}
}
