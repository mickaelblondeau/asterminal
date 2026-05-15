package model

import (
	"asterminal/internal/config"
	"asterminal/internal/converter"
	"math"
	"strings"
	"time"
)

type Vector struct {
	X, Y, Z int8
}

func (v *Vector) Reset() {
	v.X = 0
	v.Y = 0
	v.Z = 0
}

type ValueRange struct {
	Min, Max, Val float64
}

type CameraPosition struct {
	Lat, Lon float64
}

type ObjectPosition struct {
	RA, Dec, Dist float64
}

type OrbitalKind int

const (
	KindStar OrbitalKind = iota
	KindPlanet
	KindSun
	KindMoon
)

type Camera struct {
	Position CameraPosition
	Yaw      ValueRange
	Pitch    ValueRange
	Fov      float64
	Zoom     float64
}

type Object struct {
	Position  ObjectPosition
	Kind      OrbitalKind
	Elements  converter.OrbitalElements
	Radius    float64
	Name      string
	Color     string
	trackable bool
	ECEF      converter.Vector
	Distance  float64
}

type Model struct {
	Camera          Camera
	Objects         []Object
	Control         Vector
	TimeControl     int64
	TrackControl    int8
	TimeOffset      int64
	TrackedObject   *Object
	trackedObjectId int
}

const (
	startYaw       = 0
	startPitch     = 30
	baseYawSpeed   = 15
	basePitchSpeed = 15
	fovMin         = 0.04
	fovMax         = 120.0
	zoomSpeed      = 0.5
	zoomMultiplier = 3000
	timeOffset     = 1
)

func NewModel(config config.Config) *Model {
	m := &Model{
		Camera: Camera{CameraPosition{config.Lat, config.Lon}, ValueRange{Val: startYaw}, ValueRange{Val: startPitch}, fovMax, config.Zoom},
		Objects: []Object{
			// Major stars
			{Kind: KindStar, Position: ObjectPosition{37.95, 89.26, 1e20}, Name: "Polaris", Color: "#ffffff", Radius: 1},
			{Kind: KindStar, Position: ObjectPosition{101.29, -16.72, 1e20}, Name: "Sirius", Color: "#b0c8ff", Radius: 1},
			{Kind: KindStar, Position: ObjectPosition{279.23, 38.78, 1e20}, Name: "Vega", Color: "#d0d8ff", Radius: 1},
			{Kind: KindStar, Position: ObjectPosition{247.35, -26.43, 1e20}, Name: "Antares", Color: "#ff4020", Radius: 1},

			// Ursa Major
			{Kind: KindStar, Position: ObjectPosition{165.46, 56.38, 1e20}, Name: "Dubhe", Color: "#126b9e", Radius: 1},
			{Kind: KindStar, Position: ObjectPosition{165.93, 61.75, 1e20}, Name: "Merak", Color: "#126b9e", Radius: 1},
			{Kind: KindStar, Position: ObjectPosition{178.46, 53.69, 1e20}, Name: "Phecda", Color: "#126b9e", Radius: 1},
			{Kind: KindStar, Position: ObjectPosition{183.86, 57.03, 1e20}, Name: "Megrez", Color: "#126b9e", Radius: 1},
			{Kind: KindStar, Position: ObjectPosition{193.51, 55.96, 1e20}, Name: "Alioth", Color: "#126b9e", Radius: 1},
			{Kind: KindStar, Position: ObjectPosition{200.98, 54.93, 1e20}, Name: "Mizar", Color: "#126b9e", Radius: 1},
			{Kind: KindStar, Position: ObjectPosition{206.89, 49.31, 1e20}, Name: "Alkaid", Color: "#126b9e", Radius: 1},

			// Orion
			{Kind: KindStar, Position: ObjectPosition{88.79, 7.41, 1e20}, Name: "Betelgeuse", Color: "#ff6030", Radius: 1},
			{Kind: KindStar, Position: ObjectPosition{78.63, -8.20, 1e20}, Name: "Rigel", Color: "#ff6030", Radius: 1},
			{Kind: KindStar, Position: ObjectPosition{83.82, -5.91, 1e20}, Name: "Alnilam", Color: "#ff6030", Radius: 1},
			{Kind: KindStar, Position: ObjectPosition{81.28, -1.20, 1e20}, Name: "Mintaka", Color: "#ff6030", Radius: 1},
			{Kind: KindStar, Position: ObjectPosition{85.19, -1.94, 1e20}, Name: "Alnitak", Color: "#ff6030", Radius: 1},

			// Cassiopia
			{Kind: KindStar, Position: ObjectPosition{2.29, 59.15, 1e20}, Name: "Schedar", Color: "#5084ff", Radius: 1},
			{Kind: KindStar, Position: ObjectPosition{14.18, 60.72, 1e20}, Name: "Caph", Color: "#5084ff", Radius: 1},
			{Kind: KindStar, Position: ObjectPosition{9.24, 59.15, 1e20}, Name: "Gamma Cas", Color: "#5084ff", Radius: 1},
			{Kind: KindStar, Position: ObjectPosition{21.45, 60.24, 1e20}, Name: "Ruchbah", Color: "#5084ff", Radius: 1},
			{Kind: KindStar, Position: ObjectPosition{28.60, 63.67, 1e20}, Name: "Segin", Color: "#5084ff", Radius: 1},

			// Orbitals
			{Kind: KindPlanet, Elements: converter.PlutoElements, Name: "Pluto", Color: "#c9b08f", Radius: 1_188_300, trackable: true},
			{Kind: KindPlanet, Elements: converter.NeptuneElements, Name: "Neptune", Color: "#3f54ba", Radius: 24_622_000, trackable: true},
			{Kind: KindPlanet, Elements: converter.UranusElements, Name: "Uranus", Color: "#7de8e8", Radius: 25_362_000, trackable: true},
			{Kind: KindPlanet, Elements: converter.SaturnElements, Name: "Saturn", Color: "#e8d5a3", Radius: 58_232_000, trackable: true},
			{Kind: KindPlanet, Elements: converter.JupiterElements, Name: "Jupiter", Color: "#c88b3a", Radius: 69_911_000, trackable: true},
			{Kind: KindPlanet, Elements: converter.MarsElements, Name: "Mars", Color: "#ff4500", Radius: 3_389_500, trackable: true},
			{Kind: KindSun, Name: "Sun", Color: "#ffff00", Radius: 695_700_000, trackable: true},
			{Kind: KindPlanet, Elements: converter.VenusElements, Name: "Venus", Color: "#ffe5b4", Radius: 6_051_800, trackable: true},
			{Kind: KindPlanet, Elements: converter.MercuryElements, Name: "Mercury", Color: "#b5b8b1", Radius: 2_439_700, trackable: true},
			{Kind: KindMoon, Name: "Moon", Color: "#b8b8b8", Radius: 1_737_400, trackable: true},
		},
		trackedObjectId: -1,
		TimeOffset:      config.TimeOffset,
	}

	if config.Track != "" {
		m.TrackNamed(config.Track)
	}

	m.UpdatePositions()

	return m
}

func (m *Model) UpdatePositions() {
	t := time.Now().Unix() + m.TimeOffset

	for i := range m.Objects {
		obj := &m.Objects[i]
		switch obj.Kind {
		case KindPlanet:
			ecef, dist := converter.PlanetGeocentricECEF(obj.Elements, t)
			obj.ECEF = ecef
			obj.Distance = dist
		case KindSun:
			ecef, dist := converter.SunGeocentricECEF(t)
			obj.ECEF = ecef
			obj.Distance = dist
		case KindMoon:
			ecef, dist := converter.MoonGeocentricECEF(t)
			obj.ECEF = ecef
			obj.Distance = dist
		}
	}
}

func (m *Model) Update() {
	if m.Control.Z == 1 {
		m.Camera.Zoom += zoomSpeed

		if m.Camera.Zoom > 100 {
			m.Camera.Zoom = 100
		}
	}

	if m.Control.Z == -1 {
		m.Camera.Zoom -= zoomSpeed

		if m.Camera.Zoom < 0 {
			m.Camera.Zoom = 0
		}
	}

	z := m.Camera.Zoom / 100
	currentFov := fovMax * math.Pow(fovMin/fovMax, z)
	yawSpeed := baseYawSpeed * (currentFov / fovMax)
	pitchSpeed := basePitchSpeed * (currentFov / fovMax)

	m.UpdatePositions()

	if m.TrackedObject != nil {
		yaw, pitch, visible := m.GetTrackedObjectInfo()

		if visible {
			m.Camera.Yaw.Val = yaw
			m.Camera.Pitch.Val = pitch
		} else {
			m.TrackedObject = nil
			m.trackedObjectId = -1
		}
	} else {
		if m.Control.X == 1 {
			m.Camera.Yaw.Val += yawSpeed

			if m.Camera.Yaw.Val == 360 {
				m.Camera.Yaw.Val = 0
			}

			if m.Camera.Yaw.Val > 360 {
				m.Camera.Yaw.Val -= 360
			}
		}

		if m.Control.X == -1 {
			m.Camera.Yaw.Val -= yawSpeed

			if m.Camera.Yaw.Val < 0 {
				m.Camera.Yaw.Val += 360
			}
		}

		if m.Control.Y == 1 {
			m.Camera.Pitch.Val += pitchSpeed

			if m.Camera.Pitch.Val > 90 {
				m.Camera.Pitch.Val = 90
			}
		}

		if m.Control.Y == -1 {
			m.Camera.Pitch.Val -= pitchSpeed

			if m.Camera.Pitch.Val < 0 {
				m.Camera.Pitch.Val = 0
			}
		}
	}

	if m.TimeControl != 0 {
		m.TimeOffset += timeOffset * m.TimeControl
	}

	if m.TrackControl == 1 {
		m.TrackNext()
	}

	if m.TrackControl == -1 {
		m.TrackPrev()
	}

	m.Control.Reset()
	m.TimeControl = 0
	m.TrackControl = 0

	m.UpdateCamera()
}

func (m *Model) UpdateCamera() {
	t := m.Camera.Zoom / 100
	currentFov := fovMax * math.Pow(fovMin/fovMax, t)
	m.Camera.Fov = currentFov

	m.Camera.Yaw.Min = m.Camera.Yaw.Val - m.Camera.Fov/2
	m.Camera.Yaw.Max = m.Camera.Yaw.Val + m.Camera.Fov/2

	m.Camera.Pitch.Min = m.Camera.Pitch.Val - m.Camera.Fov/2
	m.Camera.Pitch.Max = m.Camera.Pitch.Val + m.Camera.Fov/2

	if m.Camera.Yaw.Min < 0 {
		m.Camera.Yaw.Min += 360
	}

	if m.Camera.Yaw.Max == 360 {
		m.Camera.Yaw.Max = 0
	}

	if m.Camera.Yaw.Max > 360 {
		m.Camera.Yaw.Max -= 360
	}
}

func (m *Model) GetTrackedObjectInfo() (float64, float64, bool) {
	return converter.AzimuthElevation(m.Camera.Position.Lat, m.Camera.Position.Lon, m.TrackedObject.ECEF)
}

func (m *Model) ResetCamera() {
	m.Camera.Pitch.Val = 0
	m.Camera.Yaw.Val = 0

	m.UpdateCamera()
}

func (m *Model) TrackStop() {
	if m.TrackedObject != nil {
		m.TrackedObject = nil
		m.trackedObjectId = -1
	}
}

func (m *Model) TrackNamed(name string) {
	for i := range m.Objects {
		if m.Objects[i].trackable && strings.ToLower(m.Objects[i].Name) == name {
			m.TrackedObject = &m.Objects[i]
			m.trackedObjectId = i
		}
	}
}

func (m *Model) TrackNext() {
	previousId := m.trackedObjectId

	for i := range m.Objects {
		if m.Objects[i].trackable && i > m.trackedObjectId {
			m.TrackedObject = &m.Objects[i]
			m.trackedObjectId = i

			if _, _, visible := m.GetTrackedObjectInfo(); visible {
				break
			}
		}
	}

	if m.trackedObjectId == previousId && m.trackedObjectId != -1 {
		m.trackedObjectId = -1
		m.TrackNext()
	}
}

func (m *Model) TrackPrev() {
	previousId := m.trackedObjectId

	for i := len(m.Objects) - 1; i >= 0; i-- {
		if m.Objects[i].trackable && i < previousId {
			m.TrackedObject = &m.Objects[i]
			m.trackedObjectId = i

			if _, _, visible := m.GetTrackedObjectInfo(); visible {
				break
			}
		}
	}

	if m.trackedObjectId == previousId && m.trackedObjectId != math.MaxInt {
		m.trackedObjectId = math.MaxInt
		m.TrackPrev()
	}
}
