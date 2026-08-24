package render

import (
	"encoding/json"
	"fmt"

	"bridge-trajectory/domain"
)

func JSON(value any) (string, error) {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return "", fmt.Errorf("render json: %w", err)
	}
	return string(data), nil
}

func ProjectionJSON(projection domain.Projection) (string, error) {
	return JSON(projection)
}

func Summary(record domain.TrajectoryRecord, metrics string) string {
	return fmt.Sprintf("trajectory %s (%s) %d points: %s", record.ID, record.View, len(record.Points), metrics)
}
