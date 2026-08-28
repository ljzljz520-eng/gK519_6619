package domain

import (
	"fmt"
	"time"
)

type ViewMode string

const (
	ViewTop  ViewMode = "top"
	ViewSide ViewMode = "side"
)

type BridgeRecord struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Length      float64 `json:"length"`
	DeckHeight  float64 `json:"deck_height"`
	CreatedUnix int64   `json:"created_unix"`
	UpdatedUnix int64   `json:"updated_unix"`
}

type WindScenarioRecord struct {
	ID          string  `json:"id"`
	BridgeID    string  `json:"bridge_id"`
	WindSpeed   float64 `json:"wind_speed"`
	Amplitude   float64 `json:"amplitude"`
	Step        float64 `json:"step"`
	Duration    float64 `json:"duration"`
	Description string  `json:"description"`
}

type TrajectoryPoint struct {
	Time float64 `json:"time"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	Z    float64 `json:"z"`
}

type TrajectoryRecord struct {
	ID          string            `json:"id"`
	BridgeID    string            `json:"bridge_id"`
	ScenarioID  string            `json:"scenario_id"`
	Points      []TrajectoryPoint `json:"points"`
	View        ViewMode          `json:"view"`
	Status      string            `json:"status"`
	CreatedUnix int64             `json:"created_unix"`
}

type ViewPreferenceRecord struct {
	UserID  string   `json:"user_id"`
	Mode    ViewMode `json:"mode"`
	Updated int64    `json:"updated"`
}

type AuditEvent struct {
	ID      string `json:"id"`
	Kind    string `json:"kind"`
	Subject string `json:"subject"`
	Detail  string `json:"detail"`
	AtUnix  int64  `json:"at_unix"`
}

type CalculationRequest struct {
	BridgeID   string
	ScenarioID string
	WindSpeed  float64
	Amplitude  float64
	Step       float64
	Duration   float64
	View       ViewMode
}

type TrajectoryFilter struct {
	BridgeID   string
	ScenarioID string
	Status     string
	View       ViewMode
	Limit      int
}

type Projection struct {
	Mode   ViewMode
	Points []ProjectedPoint
}

type ProjectedPoint struct {
	Label string  `json:"label"`
	A     float64 `json:"a"`
	B     float64 `json:"b"`
}

func (r BridgeRecord) Validate() error {
	if r.ID == "" {
		return fmt.Errorf("bridge id is required")
	}
	if r.Length <= 0 || r.DeckHeight <= 0 {
		return fmt.Errorf("bridge dimensions must be positive")
	}
	if r.Name == "" {
		return fmt.Errorf("bridge name is required")
	}
	return nil
}

func (r WindScenarioRecord) Validate() error {
	if r.ID == "" || r.BridgeID == "" {
		return fmt.Errorf("scenario identity is required")
	}
	if r.WindSpeed < 0 || r.Amplitude < 0 || r.Step <= 0 || r.Duration <= 0 {
		return fmt.Errorf("scenario values are out of range")
	}
	if r.Step > r.Duration {
		return fmt.Errorf("step cannot exceed duration")
	}
	return nil
}

func (r TrajectoryRecord) Validate() error {
	if r.ID == "" || r.BridgeID == "" || r.ScenarioID == "" {
		return fmt.Errorf("trajectory identity is required")
	}
	if len(r.Points) == 0 {
		return fmt.Errorf("trajectory has no points")
	}
	if r.View != ViewTop && r.View != ViewSide {
		return fmt.Errorf("unsupported view %q", r.View)
	}
	return nil
}

func NewBridge(id, name string, length, deckHeight float64, now time.Time) BridgeRecord {
	stamp := now.Unix()
	return BridgeRecord{ID: id, Name: name, Length: length, DeckHeight: deckHeight, CreatedUnix: stamp, UpdatedUnix: stamp}
}

func NewScenario(id, bridgeID string, wind, amplitude, step, duration float64, description string) WindScenarioRecord {
	return WindScenarioRecord{ID: id, BridgeID: bridgeID, WindSpeed: wind, Amplitude: amplitude, Step: step, Duration: duration, Description: description}
}

func CopyPoints(points []TrajectoryPoint) []TrajectoryPoint {
	if points == nil {
		return nil
	}
	copyOf := make([]TrajectoryPoint, len(points))
	copy(copyOf, points)
	return copyOf
}

func (r TrajectoryRecord) Clone() TrajectoryRecord {
	r.Points = CopyPoints(r.Points)
	return r
}

func DefaultView(mode ViewMode) ViewMode {
	if mode == ViewSide {
		return ViewSide
	}
	return ViewTop
}
