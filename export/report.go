package export

import (
	"fmt"
	"sort"
	"strings"

	"bridge-trajectory/analysis"
	"bridge-trajectory/domain"
)

func Markdown(record domain.TrajectoryRecord, report analysis.Report) string {
	var builder strings.Builder
	builder.WriteString("# Bridge trajectory report\n\n")
	fmt.Fprintf(&builder, "Trajectory: `%s`\n\n", record.ID)
	fmt.Fprintf(&builder, "View: **%s**  \nPoints: **%d**  \nDuration: **%.2f**\n\n", record.View, report.PointCount, report.Duration)
	builder.WriteString("| Metric | Value |\n|---|---:|\n")
	fmt.Fprintf(&builder, "| Distance | %.4f |\n", report.Distance)
	fmt.Fprintf(&builder, "| Maximum speed | %.4f |\n", report.MaxSpeed)
	fmt.Fprintf(&builder, "| Maximum lateral | %.4f |\n", report.MaxLateral)
	fmt.Fprintf(&builder, "| Maximum vertical | %.4f |\n", report.MaxVertical)
	fmt.Fprintf(&builder, "| Risk score | %.2f (%s) |\n", report.RiskScore, analysis.RiskBand(report.RiskScore))
	if len(report.Warnings) > 0 {
		builder.WriteString("\nWarnings:\n")
		for _, warning := range report.Warnings {
			fmt.Fprintf(&builder, "- %s\n", warning)
		}
	}
	return builder.String()
}

func BatchMarkdown(reports []analysis.Report) string {
	sorted := analysis.Rank(reports)
	var builder strings.Builder
	builder.WriteString("# Trajectory comparison\n\n| Trajectory | Risk | Band |\n|---|---:|---|\n")
	for _, report := range sorted {
		fmt.Fprintf(&builder, "| %s | %.2f | %s |\n", report.TrajectoryID, report.RiskScore, analysis.RiskBand(report.RiskScore))
	}
	return builder.String()
}

func SummaryLines(report analysis.Report) []string {
	lines := []string{fmt.Sprintf("trajectory %s", report.TrajectoryID), fmt.Sprintf("risk %.2f", report.RiskScore), fmt.Sprintf("distance %.3f", report.Distance)}
	if len(report.Warnings) == 0 {
		lines = append(lines, "warnings none")
	} else {
		warnings := append([]string(nil), report.Warnings...)
		sort.Strings(warnings)
		lines = append(lines, "warnings "+strings.Join(warnings, "; "))
	}
	return lines
}
