package calc

import (
	"testing"

	"bridge-trajectory/domain"
)

func TestSamplerProducesDeterministicSamples(t *testing.T) {
	sampler := NewSampler()
	model := NewWindModel(20, 1.5)
	first, err := sampler.Sample(300, 30, model, 0.5, 3)
	if err != nil {
		t.Fatal(err)
	}
	second, err := sampler.Sample(300, 30, model, 0.5, 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(first) != 7 || len(second) != len(first) {
		t.Fatalf("unexpected sample count %d", len(first))
	}
	for index := range first {
		if first[index] != second[index] {
			t.Fatalf("sample %d changed", index)
		}
	}
}

func TestSamplerEnvelopeAndResample(t *testing.T) {
	points := []domain.TrajectoryPoint{{Y: -2, Z: 4}, {Y: 3, Z: 8}}
	sampler := NewSampler()
	minY, maxY, minZ, maxZ := sampler.SampleEnvelope(points)
	if minY != -2 || maxY != 3 || minZ != 4 || maxZ != 8 {
		t.Fatalf("unexpected envelope")
	}
	resampled := sampler.Resample(points, 2)
	if resampled[1].Y != 6 {
		t.Fatalf("unexpected resample")
	}
}
