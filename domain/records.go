package domain

import "fmt"

func (r TrajectoryRecord) Duration() float64 {
	if len(r.Points) == 0 {
		return 0
	}
	return r.Points[len(r.Points)-1].Time - r.Points[0].Time
}

func (r TrajectoryRecord) PointAt(index int) (TrajectoryPoint, error) {
	if index < 0 || index >= len(r.Points) {
		return TrajectoryPoint{}, fmt.Errorf("point index %d is out of range", index)
	}
	return r.Points[index], nil
}

func (r TrajectoryRecord) IsReady() bool { return r.Status == "ready" && len(r.Points) > 0 }

func (r TrajectoryRecord) WithView(mode ViewMode) (TrajectoryRecord, error) {
	if mode != ViewTop && mode != ViewSide {
		return TrajectoryRecord{}, fmt.Errorf("unsupported view %q", mode)
	}
	clone := r.Clone()
	clone.View = mode
	return clone, nil
}

func (r TrajectoryRecord) SampleTimes() []float64 {
	times := make([]float64, len(r.Points))
	for index, point := range r.Points {
		times[index] = point.Time
	}
	return times
}

func (r TrajectoryRecord) LastPoint() (TrajectoryPoint, bool) {
	if len(r.Points) == 0 {
		return TrajectoryPoint{}, false
	}
	return r.Points[len(r.Points)-1], true
}

func (r BridgeRecord) Rename(name string) (BridgeRecord, error) {
	if name == "" {
		return BridgeRecord{}, fmt.Errorf("bridge name is required")
	}
	r.Name = name
	return r, nil
}

func (r WindScenarioRecord) WithStep(step float64) (WindScenarioRecord, error) {
	updated := r
	updated.Step = step
	if err := updated.Validate(); err != nil {
		return WindScenarioRecord{}, err
	}
	return updated, nil
}

func (f TrajectoryFilter) Match(record TrajectoryRecord) bool {
	if f.BridgeID != "" && f.BridgeID != record.BridgeID {
		return false
	}
	if f.ScenarioID != "" && f.ScenarioID != record.ScenarioID {
		return false
	}
	if f.Status != "" && f.Status != record.Status {
		return false
	}
	if f.View != "" && f.View != record.View {
		return false
	}
	return true
}

func (m ViewMode) Valid() bool { return m == ViewTop || m == ViewSide }

func (e AuditEvent) Validate() error {
	if e.ID == "" {
		return fmt.Errorf("event id is required")
	}
	if e.Kind == "" {
		return fmt.Errorf("event kind is required")
	}
	if e.Subject == "" {
		return fmt.Errorf("event subject is required")
	}
	return nil
}

func (e AuditEvent) Label() string {
	if e.Detail == "" {
		return e.Kind + ":" + e.Subject
	}
	return e.Kind + ":" + e.Subject + " - " + e.Detail
}
