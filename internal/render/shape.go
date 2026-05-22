package render

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

const (
	aspect          = 0.5
	raidusForPoint  = 1
	radiusForMedium = 0.05
	radiusForSmall  = 0.001
	radiusForTiny   = 0.00001
)

type cell struct {
	char, fg, bg string
}

var styleCache = make(map[string]lipgloss.Style)

func getCachedStyle(fg, bg string) lipgloss.Style {
	key := fmt.Sprintf("%s:%s", fg, bg)

	if _, ok := styleCache[key]; !ok {
		style := lipgloss.NewStyle()

		if fg != "" {
			style = style.Foreground(lipgloss.Color(fg))
		}

		if bg != "" {
			style = style.Background(lipgloss.Color(bg))
		}

		styleCache[key] = style
	}

	return styleCache[key]
}

func drawCircle(objectMap map[string]cell, cx int, cy int, r float64, color string) {
	char := "⬤"

	if r < radiusForMedium {
		char = "●"
	}
	if r < radiusForSmall {
		char = "•"
	}
	if r < radiusForTiny {
		char = "·"
	}

	if r < raidusForPoint {
		key := fmt.Sprintf("%d:%d", cx, cy)
		currentCell := objectMap[key]
		objectMap[key] = cell{fg: color, char: char, bg: currentCell.bg}
		return
	}

	rx := int(r / aspect)
	ry := int(r)

	for y := -ry; y <= ry; y++ {
		for x := -rx; x <= rx; x++ {
			dx := float64(x) * aspect
			dy := float64(y)

			if dx*dx+dy*dy <= r*r {
				objectMap[fmt.Sprintf("%d:%d", cx+x, cy+y)] = cell{bg: color, char: " "}
			}
		}
	}
}

func drawStar(objectMap map[string]cell, x int, y int, color string) {
	objectMap[fmt.Sprintf("%d:%d", x, y)] = cell{fg: color, char: "✦"}
}
