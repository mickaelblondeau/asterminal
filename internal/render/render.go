package render

import (
	"asterminal/internal/converter"
	"asterminal/internal/model"
	"fmt"
	"math"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	refreshRate  = time.Second / 60
	zoomSize     = 40
	dateSize     = 21
	trackingSize = 30
)

type tickMsg time.Time

var groundStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#007700"))

type Renderer struct {
	model                                                                                                            *model.Model
	infoStyle, controlStyle, trackingStyle, dateStyle, zoomStyle, mapStyle, trackedObjectStyle, pitchStyle, yawStyle lipgloss.Style
	trackedObject                                                                                                    *model.Object
	showControls                                                                                                     bool
}

func NewRenderer(model *model.Model) *Renderer {
	return &Renderer{model: model}
}

func tick() tea.Cmd {
	return tea.Tick(refreshRate, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m *Renderer) Init() tea.Cmd {
	return tick()
}

func (m *Renderer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.model.Update()

		if m.model.TrackedObject != nil && m.model.TrackedObject != m.trackedObject {
			m.trackedObjectStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(m.model.TrackedObject.Color))
			m.trackedObject = m.model.TrackedObject
		}

		return m, tick()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyRight:
			m.model.Control.X = 1
		case tea.KeyLeft:
			m.model.Control.X = -1
		case tea.KeyUp:
			m.model.Control.Y = 1
		case tea.KeyDown:
			m.model.Control.Y = -1
		case tea.KeyPgUp:
			m.model.TimeControl = 1
		case tea.KeyPgDown:
			m.model.TimeControl = -1
		case tea.KeyCtrlPgUp:
			m.model.TimeControl = 60
		case tea.KeyCtrlPgDown:
			m.model.TimeControl = -60
		case tea.KeyEsc:
			m.model.TrackStop()
			m.showControls = false
		}
		switch msg.String() {
		case "+":
			m.model.Control.Z = 1
		case "-":
			m.model.Control.Z = -1
		case "o":
			m.model.TrackControl = -1
		case "p":
			m.model.TrackControl = 1
		case "c":
			m.showControls = true
		}
	case tea.WindowSizeMsg:
		m.infoStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Height(3).Width(msg.Width - trackingSize - dateSize - zoomSize - 8).Padding(1)
		m.controlStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Height(msg.Height - 2).Width(msg.Width - 2).Padding(1)
		m.trackingStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Height(3).Width(trackingSize).Padding(1).Align(lipgloss.Center)
		m.dateStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Height(3).Width(dateSize).Padding(1)
		m.zoomStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Height(3).Width(zoomSize).Padding(1)
		m.mapStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Height(msg.Height - 12).Width(msg.Width - 9).Padding(1)
		m.pitchStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Height(msg.Height - 12).Width(5).Padding(1)
		m.yawStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Height(3).Width(msg.Width - 2).Padding(1)
	}

	return m, nil
}

func (m *Renderer) View() string {
	if m.showControls {
		return m.controlStyle.Render(renderControls(m.controlStyle.GetWidth() - 2))
	}

	top := lipgloss.JoinHorizontal(lipgloss.Left, m.infoStyle.Render(renderInfo(m.infoStyle.GetWidth()-2)), m.trackingStyle.Render(renderTracking(m.model.TrackedObject, m.trackedObjectStyle)), m.dateStyle.Render(renderDate(m.model.TimeOffset)), m.zoomStyle.Render(renderZoom(m.model.Camera.Zoom)))
	middle := lipgloss.JoinHorizontal(lipgloss.Left, m.mapStyle.Render(renderMap(m.mapStyle.GetWidth()-2, m.mapStyle.GetHeight()-2, m.model.Camera, m.model.Objects, m.model.TimeOffset)), m.pitchStyle.Render(renderElevation(m.pitchStyle.GetHeight(), m.model.Camera.Pitch)))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		top,
		middle,
		m.yawStyle.Render(renderRotation(m.yawStyle.GetWidth(), m.model.Camera.Yaw)),
	)
}

func renderMap(w, h int, camera model.Camera, objects []model.Object, timeOffset int64) string {
	var res strings.Builder
	objectMap := make(map[string]cell)
	t := time.Now().Unix() + timeOffset

	for i := range objects {
		object := &objects[i]
		ecef := object.ECEF

		if object.Kind == model.KindStar {
			ecef = converter.ObjectECEFVec(object.Position.RA, object.Position.Dec, object.Position.Dist, t)
		}

		posX, posY, visible := converter.ConvertObjectToScreenSpace(camera.Position.Lat, camera.Position.Lon, camera.Yaw.Val, camera.Pitch.Val, camera.Fov, w, h, ecef)

		if !visible {
			continue
		}

		coordX := int(math.Round(posX * float64(w)))
		coordY := int(math.Round(posY * float64(h)))

		dist := object.Position.Dist
		if object.Kind != model.KindStar {
			dist = object.Distance
		}

		r := converter.NormalizedRadius(object.Radius, dist, camera.Fov)
		screenR := r * float64(h) / 2.0
		intRadius := int(screenR)

		if coordX >= -intRadius*2 && coordX <= w+intRadius*2 && coordY >= -intRadius*2 && coordY <= h+intRadius*2 {
			if object.Radius == 1 {
				drawStar(objectMap, coordX, coordY, object.Color)
			} else {
				drawCircle(objectMap, coordX, coordY, screenR, w, h, object.Color)
			}
		}
	}

	groundY := int((0.5 + float64(camera.Pitch.Val)/float64(camera.Fov)) * float64(h))

	for y := range h {
		if y == groundY {
			res.WriteString(strings.Repeat("─", w))
		} else if y > groundY {
			res.WriteString(groundStyle.Render(strings.Repeat("░", w)))
		} else {
			for x := 0; x < w; {
				start := x
				val, ok := objectMap[fmt.Sprintf("%d:%d", x, y)]

				if ok {
					for x < w {
						x++
						val2, ok := objectMap[fmt.Sprintf("%d:%d", x, y)]

						if !ok || val2.bg != val.bg || val2.fg != val.fg {
							break
						}
					}

					res.WriteString(getCachedStyle(val.fg, val.bg).Render(strings.Repeat(val.char, x-start)))
				} else {
					res.WriteString("\u00A0")
					x++
				}
			}
		}
	}

	return res.String()
}

func renderControls(w int) string {
	title := "ASTerminal"
	subTitle := "Press <ESC> to go back"

	var res strings.Builder
	res.WriteString(title)

	spaces := w - len(title) - len(subTitle)

	for range spaces {
		res.WriteString(" ")
	}

	res.WriteString(subTitle + "\n\n")
	res.WriteString("Rotate: ↑ ↓ ← →\n")
	res.WriteString("Zoom +: +\n")
	res.WriteString("Zoom -: -\n")
	res.WriteString("Time +: ⇞\n")
	res.WriteString("Time ++: ctrl+⇞\n")
	res.WriteString("Time -: ⇟\n")
	res.WriteString("Time --: ctrl+⇟\n")
	res.WriteString("Track next object: p\n")
	res.WriteString("Track previous object: o\n")
	res.WriteString("Cancel tracking: esc\n")

	return res.String()
}

func renderInfo(w int) string {
	title := "ASTerminal"
	subTitle := "Press <C> for controls"

	var res strings.Builder
	res.WriteString(title)

	spaces := w - len(title) - len(subTitle)

	for range spaces {
		res.WriteString(" ")
	}

	return res.String() + subTitle
}

func renderTracking(obj *model.Object, style lipgloss.Style) string {
	if obj != nil {
		return fmt.Sprintf("Tracking: %s %s", style.Render("⬤"), obj.Name)
	}

	return "Tracking: Nothing"
}

func renderDate(date int64) string {
	return time.Unix(time.Now().Unix()+date, 0).Format("2006-01-02 15:04:05")
}

func renderZoom(zoom float64) string {
	var res strings.Builder
	res.WriteString("- ")
	scaledVal := int((zoom / 100) * (zoomSize - 7))

	for i := range zoomSize - 6 {
		if i == scaledVal {
			res.WriteString("●")
		} else {
			res.WriteString("┄")
		}
	}

	return res.String() + " +"
}

func getElevationDisplay(v float64) string {
	s := fmt.Sprintf("%.3f", v-math.Floor(v))
	parts := strings.Split(s, ".")
	return fmt.Sprintf("%03.0f\n%s", v, parts[1])
}

func renderElevation(h int, elevation model.ValueRange) string {
	var res strings.Builder
	res.WriteString(getElevationDisplay(elevation.Max))

	spaces := h/2 - 4

	for range spaces {
		res.WriteString("\n")
	}

	res.WriteString(" ↑\n\n" + getElevationDisplay(elevation.Val) + "\n\n ↓")

	for range spaces - 2 {
		res.WriteString("\n")
	}

	return res.String() + getElevationDisplay(elevation.Min)
}

func renderRotation(w int, rotation model.ValueRange) string {
	var res strings.Builder
	fmt.Fprintf(&res, "%03.2f", rotation.Min)

	spaces := (w-11)/2 - 8

	for range spaces {
		res.WriteString(" ")
	}

	fmt.Fprintf(&res, "←  %03.2f  →", rotation.Val)

	for range spaces {
		res.WriteString(" ")
	}

	return res.String() + fmt.Sprintf("%03.2f", rotation.Max)
}
