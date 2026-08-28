package export

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"bridge-trajectory/domain"
)

func TrajectoryCSV(record domain.TrajectoryRecord) (string, error) {
	var builder strings.Builder
	if err := WriteTrajectoryCSV(&builder, record); err != nil {
		return "", err
	}
	return builder.String(), nil
}

func WriteTrajectoryCSV(writer io.Writer, record domain.TrajectoryRecord) error {
	if err := record.Validate(); err != nil {
		return err
	}
	encoder := csv.NewWriter(writer)
	if err := encoder.Write([]string{"trajectory_id", "time", "x", "y", "z", "view", "status"}); err != nil {
		return err
	}
	for _, point := range record.Points {
		row := []string{record.ID, number(point.Time), number(point.X), number(point.Y), number(point.Z), string(record.View), record.Status}
		if err := encoder.Write(row); err != nil {
			return err
		}
	}
	encoder.Flush()
	return encoder.Error()
}

func ProjectionCSV(projection domain.Projection) (string, error) {
	var builder strings.Builder
	encoder := csv.NewWriter(&builder)
	if err := encoder.Write([]string{"label", "a", "b", "mode"}); err != nil {
		return "", err
	}
	for _, point := range projection.Points {
		if err := encoder.Write([]string{point.Label, number(point.A), number(point.B), string(projection.Mode)}); err != nil {
			return "", err
		}
	}
	encoder.Flush()
	if err := encoder.Error(); err != nil {
		return "", err
	}
	return builder.String(), nil
}

func number(value float64) string { return strconv.FormatFloat(value, 'f', 6, 64) }

func ParseNumber(value string) (float64, error) {
	parsed, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	if err != nil {
		return 0, fmt.Errorf("parse number %q: %w", value, err)
	}
	return parsed, nil
}
