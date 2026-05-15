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

func drawCircle(objectMap map[string]string, cx int, cy int, r float64, color string) {
	if r < raidusForPoint {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(color))

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

		objectMap[fmt.Sprintf("%d:%d", cx, cy)] = style.Render(char)
		return
	}

	style := lipgloss.NewStyle().Background(lipgloss.Color(color))

	rx := int(r / aspect)
	ry := int(r)

	for y := -ry; y <= ry; y++ {
		for x := -rx; x <= rx; x++ {
			dx := float64(x) * aspect
			dy := float64(y)

			if dx*dx+dy*dy <= r*r {
				objectMap[fmt.Sprintf("%d:%d", cx+x, cy+y)] = style.Render(" ")
			}
		}
	}
}

func drawStar(objectMap map[string]string, x int, y int, color string) {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(color))

	objectMap[fmt.Sprintf("%d:%d", x, y)] = style.Render("✦")
}
