package domain

import (
	"fmt"
	"math"
)

func ValidateCalculation(req CalculationRequest) error {
	if req.BridgeID == "" || req.ScenarioID == "" {
		return fmt.Errorf("bridge and scenario are required")
	}
	if math.IsNaN(req.WindSpeed) || math.IsInf(req.WindSpeed, 0) {
		return fmt.Errorf("wind speed must be finite")
	}
	if math.IsNaN(req.Amplitude) || math.IsInf(req.Amplitude, 0) {
		return fmt.Errorf("amplitude must be finite")
	}
	if req.WindSpeed < 0 || req.Amplitude < 0 {
		return fmt.Errorf("wind speed and amplitude cannot be negative")
	}
	if req.Step <= 0 || req.Duration <= 0 || req.Step > req.Duration {
		return fmt.Errorf("invalid time range")
	}
	if req.View != ViewTop && req.View != ViewSide {
		return fmt.Errorf("view must be top or side")
	}
	return nil
}

func NormalizeLimit(limit int) int {
	if limit <= 0 {
		return 50
	}
	if limit > 500 {
		return 500
	}
	return limit
}

func IsTerminalStatus(status string) bool {
	switch status {
	case "ready", "archived", "failed":
		return true
	default:
		return false
	}
}
