package render

import (
	"fmt"
	"strings"

	"bridge-trajectory/domain"
)

func ASCII(projection domain.Projection, width, height int) string {
	if width < 10 {
		width = 10
	}
	if height < 4 {
		height = 4
	}
	grid := make([][]rune, height)
	for row := range grid {
		grid[row] = make([]rune, width)
		for col := range grid[row] {
			grid[row][col] = ' '
		}
	}
	minA, maxA, minB, maxB := projectionBounds(projection)
	for index, point := range projection.Points {
		col := scale(point.A, minA, maxA, width-1)
		row := height - 1 - scale(point.B, minB, maxB, height-1)
		if row < 0 {
			row = 0
		}
		if row >= height {
			row = height - 1
		}
		if col < 0 {
			col = 0
		}
		if col >= width {
			col = width - 1
		}
		if index%10 == 0 {
			grid[row][col] = 'O'
		} else {
			grid[row][col] = '.'
		}
	}
	lines := make([]string, 0, height)
	for _, row := range grid {
		lines = append(lines, strings.TrimRight(string(row), " "))
	}
	return strings.Join(lines, "\n")
}

func projectionBounds(projection domain.Projection) (float64, float64, float64, float64) {
	if len(projection.Points) == 0 {
		return 0, 1, 0, 1
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
	if minA == maxA {
		maxA++
	}
	if minB == maxB {
		maxB++
	}
	return minA, maxA, minB, maxB
}

func scale(value, low, high float64, size int) int {
	if high <= low {
		return 0
	}
	return int((value-low)/(high-low)*float64(size) + 0.5)
}

func Legend(mode domain.ViewMode) string {
	if mode == domain.ViewTop {
		return fmt.Sprintf("top view: X/Y plane")
	}
	return fmt.Sprintf("side view: X/Z plane")
}
