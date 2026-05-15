package converter

import (
	"math"
)

type Vector struct{ X, Y, Z float64 }

func degToRad(d float64) float64 { return d * math.Pi / 180 }

func dot(a, b Vector) float64 { return a.X*b.X + a.Y*b.Y + a.Z*b.Z }

func cross(a, b Vector) Vector {
	return Vector{
		a.Y*b.Z - a.Z*b.Y,
		a.Z*b.X - a.X*b.Z,
		a.X*b.Y - a.Y*b.X,
	}
}

func sub(a, b Vector) Vector { return Vector{a.X - b.X, a.Y - b.Y, a.Z - b.Z} }

func scale(v Vector, s float64) Vector { return Vector{v.X * s, v.Y * s, v.Z * s} }

func length(v Vector) float64 { return math.Sqrt(dot(v, v)) }

func normalize(v Vector) Vector {
	l := length(v)
	if l < 1e-12 {
		return Vector{}
	}
	return scale(v, 1/l)
}

const (
	earthA               = 6_378_137.0
	earthE2              = 6.6943799901414e-3
	secondsPerDay        = 86400.0
	julianDateUnixEpoch  = 2440587.5
	julianDateJ2000      = 2451545.0
	julianDaysPerCentury = 36525.0
	gmstJ2000Deg         = 280.46061837
	earthRotationRate    = 360.98564736629
	gmstPrecessionDeg    = 0.000387933
)

func gpsToECEF(lat, lon float64) Vector {
	latR := degToRad(lat)
	lonR := degToRad(lon)

	sinLat, cosLat := math.Sin(latR), math.Cos(latR)
	sinLon, cosLon := math.Sin(lonR), math.Cos(lonR)

	N := earthA / math.Sqrt(1-earthE2*sinLat*sinLat)

	return Vector{
		X: N * cosLat * cosLon,
		Y: N * cosLat * sinLon,
		Z: N * (1 - earthE2) * sinLat,
	}
}

func gmst(unixSec int64) float64 {
	julianDay := float64(unixSec)/secondsPerDay + julianDateUnixEpoch
	centuriesSinceJ2000 := (julianDay - julianDateJ2000) / julianDaysPerCentury
	deg := gmstJ2000Deg + earthRotationRate*(julianDay-julianDateJ2000) + gmstPrecessionDeg*centuriesSinceJ2000*centuriesSinceJ2000

	deg = math.Mod(deg, 360)
	if deg < 0 {
		deg += 360
	}

	return degToRad(deg)
}

func celestialToECEF(raDeg, decDeg float64, unixSec int64) Vector {
	ra := degToRad(raDeg)
	dec := degToRad(decDeg)

	v := Vector{
		X: math.Cos(dec) * math.Cos(ra),
		Y: math.Cos(dec) * math.Sin(ra),
		Z: math.Sin(dec),
	}

	return equatorialToECEF(v, unixSec)
}

func enuBasis(lat, lon float64) (east, north, up Vector) {
	latR := degToRad(lat)
	lonR := degToRad(lon)

	sinLat, cosLat := math.Sin(latR), math.Cos(latR)
	sinLon, cosLon := math.Sin(lonR), math.Cos(lonR)

	east = Vector{-sinLon, cosLon, 0}
	north = Vector{-sinLat * cosLon, -sinLat * sinLon, cosLat}
	up = Vector{cosLat * cosLon, cosLat * sinLon, sinLat}

	return
}

func cameraBasis(lat, lon, yawDeg, pitchDeg float64) (forward, right, camUp Vector) {
	if pitchDeg == 90 {
		pitchDeg -= 1e-4
	}

	east, north, localUp := enuBasis(lat, lon)

	yaw := degToRad(yawDeg)
	pitch := degToRad(pitchDeg)

	cosPitch, sinPitch := math.Cos(pitch), math.Sin(pitch)
	sinYaw, cosYaw := math.Sin(yaw), math.Cos(yaw)

	forward = normalize(Vector{
		X: cosPitch*(sinYaw*east.X+cosYaw*north.X) + sinPitch*localUp.X,
		Y: cosPitch*(sinYaw*east.Y+cosYaw*north.Y) + sinPitch*localUp.Y,
		Z: cosPitch*(sinYaw*east.Z+cosYaw*north.Z) + sinPitch*localUp.Z,
	})

	right = normalize(cross(forward, localUp))
	if length(cross(forward, localUp)) < 1e-6 {
		right = normalize(cross(forward, north))
	}

	camUp = normalize(cross(right, forward))

	return
}

func ObjectECEFVec(ra, dec, dist float64, t int64) Vector {
	return scale(celestialToECEF(ra, dec, t), dist)
}

func ConvertObjectToScreenSpace(lat, lon, yaw, pitch, fov float64, objECEF Vector, t int64) (float64, float64, bool) {
	camECEF := gpsToECEF(lat, lon)

	dir := normalize(sub(objECEF, camECEF))

	forward, right, camUp := cameraBasis(
		lat,
		lon,
		yaw,
		pitch,
	)

	x := dot(dir, right)
	y := dot(dir, camUp)
	z := dot(dir, forward)

	if z <= 1e-6 {
		return -1, -1, false
	}

	fov = degToRad(fov)
	s := 1.0 / math.Tan(fov/2)

	screenX := 0.5 + (x/z)*s*0.5
	screenY := 0.5 - (y/z)*s*0.5

	_, _, localUp := enuBasis(lat, lon)

	return screenX, screenY, dot(dir, localUp) >= 0
}

func ConvertPositionsToScreenSpace(lat, lon, yaw, pitch, fov, ra, dec, dist float64, t int64) (float64, float64, bool) {
	camECEF := gpsToECEF(lat, lon)
	objECEF := scale(celestialToECEF(ra, dec, t), dist)

	dir := normalize(sub(objECEF, camECEF))

	forward, right, camUp := cameraBasis(
		lat,
		lon,
		yaw,
		pitch,
	)

	x := dot(dir, right)
	y := dot(dir, camUp)
	z := dot(dir, forward)

	if z <= 1e-6 {
		return -1, -1, false
	}

	fov = degToRad(fov)
	s := 1.0 / math.Tan(fov/2)

	screenX := 0.5 + (x/z)*s*0.5
	screenY := 0.5 - (y/z)*s*0.5

	_, _, localUp := enuBasis(lat, lon)

	return screenX, screenY, dot(dir, localUp) >= 0
}

func NormalizedRadius(radius, distance, fovDeg float64) float64 {
	fovRad := degToRad(fovDeg)
	angularRadius := math.Atan(radius / distance)

	return angularRadius / (fovRad / 2.0)
}

func AzimuthElevation(lat, lon float64, objECEF Vector) (yawDeg, pitchDeg float64, visible bool) {
	camECEF := gpsToECEF(lat, lon)
	dir := normalize(sub(objECEF, camECEF))
	east, north, up := enuBasis(lat, lon)

	if dot(dir, up) < 0 {
		return 0, 0, false
	}

	pitchDeg = math.Asin(math.Max(-1, math.Min(1, dot(dir, up)))) * 180 / math.Pi
	yawDeg = math.Atan2(dot(dir, east), dot(dir, north)) * 180 / math.Pi

	if yawDeg < 0 {
		yawDeg += 360
	}

	return yawDeg, pitchDeg, true
}
