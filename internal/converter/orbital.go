package converter

import "math"

type OrbitalElements struct {
	a, e, I, L, wbar, Omega       float64
	aR, eR, IR, LR, wbarR, OmegaR float64
}

var (
	MercuryElements = OrbitalElements{
		0.38709927, 0.20563593, 7.00497902, 252.25032350, 77.45779628, 48.33076593,
		0.00000037, 0.00001906, -0.00594749, 149472.67411175, 0.16047689, -0.12534081,
	}
	VenusElements = OrbitalElements{
		0.72333566, 0.00677672, 3.39467605, 181.97909950, 131.60246718, 76.67984255,
		0.00000390, -0.00004107, -0.00078890, 58517.81538729, 0.00268329, -0.27769418,
	}
	EarthElements = OrbitalElements{
		1.00000261, 0.01671123, -0.00001531, 100.46457166, 102.93768193, 0.0,
		0.00000562, -0.00004392, -0.01294668, 35999.37244981, 0.32327364, 0.0,
	}
	MarsElements = OrbitalElements{
		1.52371034, 0.09339410, 1.84969142, -4.55343205, -23.94362959, 49.55953891,
		0.00001847, 0.00007882, -0.00813131, 19140.30268499, 0.44441088, -0.29257343,
	}
	JupiterElements = OrbitalElements{
		5.20288700, 0.04838624, 1.30439695, 34.39644051, 14.72847983, 100.47390909,
		-0.00011607, -0.00013253, -0.00183714, 3034.74612775, 0.21252668, 0.20469106,
	}
	SaturnElements = OrbitalElements{
		9.53667594, 0.05386179, 2.48599187, 49.95424423, 92.59887831, 113.66242448,
		-0.00125060, -0.00050991, 0.00193609, 1222.49362201, -0.41897216, -0.28867794,
	}
	UranusElements = OrbitalElements{
		19.18916464, 0.04725744, 0.77263783, 313.23810451, 170.95427630, 74.01692503,
		-0.00196176, -0.00004397, -0.00242939, 428.48202785, 0.40805281, 0.04240589,
	}
	NeptuneElements = OrbitalElements{
		30.06992276, 0.00859048, 1.77004347, -55.12002969, 44.96476227, 131.78422574,
		0.00026291, 0.00005105, 0.00035372, 218.45945325, -0.32241464, -0.00508664,
	}
	PlutoElements = OrbitalElements{
		39.48168677, 0.24880766, 17.14175, 238.92903833, 224.09702263, 110.30347002,
		-0.00076912, 0.00006465, 0.00000501, 145.20780515, -0.00968827, -0.01183482,
	}
)

const (
	auToMeters = 1.495978707e11
	obliquity  = 23.439291
)

func julianCenturies(unixSec int64) float64 {
	jd := float64(unixSec)/secondsPerDay + julianDateUnixEpoch

	return (jd - julianDateJ2000) / julianDaysPerCentury
}

func elementAtEpoch(el OrbitalElements, T float64) (a, e, I, L, wbar, Omega float64) {
	a = el.a + el.aR*T
	e = el.e + el.eR*T
	I = el.I + el.IR*T
	L = el.L + el.LR*T
	wbar = el.wbar + el.wbarR*T
	Omega = el.Omega + el.OmegaR*T

	return
}

func solveKepler(M, e float64) float64 {
	M = math.Mod(M, 2*math.Pi)

	if M > math.Pi {
		M -= 2 * math.Pi
	}
	if M < -math.Pi {
		M += 2 * math.Pi
	}

	E := M

	for range 50 {
		dE := (M - (E - e*math.Sin(E))) / (1 - e*math.Cos(E))
		E += dE

		if math.Abs(dE) < 1e-12 {
			break
		}
	}

	return E
}

func keplerPosition(a, e, M float64) (x, y float64) {
	E := solveKepler(M, e)

	x = a * (math.Cos(E) - e)
	y = a * math.Sqrt(1-e*e) * math.Sin(E)
	return
}

func orbitalToEcliptic(x, y float64, I, Omega, w float64) Vector {
	cO, sO := math.Cos(Omega), math.Sin(Omega)
	cI, sI := math.Cos(I), math.Sin(I)
	cW, sW := math.Cos(w), math.Sin(w)

	X := (cO*cW-sO*sW*cI)*x + (-cO*sW-sO*cW*cI)*y
	Y := (sO*cW+cO*sW*cI)*x + (-sO*sW+cO*cW*cI)*y
	Z := (sW*sI)*x + (cW*sI)*y

	return Vector{X, Y, Z}
}

func heliocentricEcliptic(a, e, I, L, wbar, Omega float64) Vector {
	M := degToRad(math.Mod(L-wbar, 360))
	xOrb, yOrb := keplerPosition(a, e, M)

	w := degToRad(wbar - Omega)
	Ir := degToRad(I)
	Or := degToRad(Omega)

	return orbitalToEcliptic(xOrb, yOrb, Ir, Or, w)
}

func eclipticToEquatorial(v Vector) Vector {
	eps := degToRad(obliquity)
	cE, sE := math.Cos(eps), math.Sin(eps)

	return Vector{
		X: v.X,
		Y: cE*v.Y - sE*v.Z,
		Z: sE*v.Y + cE*v.Z,
	}
}

func equatorialToECEF(v Vector, unixSec int64) Vector {
	theta := gmst(unixSec)
	c, s := math.Cos(theta), math.Sin(theta)

	return Vector{
		X: v.X*c + v.Y*s,
		Y: -v.X*s + v.Y*c,
		Z: v.Z,
	}
}

func planetGeocentricEq(el OrbitalElements, unixSec int64) (Vector, Vector) {
	T := julianCenturies(unixSec)
	aP, eP, IP, LP, wbarP, OmegaP := elementAtEpoch(el, T)
	aE, eE, IE, LE, wbarE, OmegaE := elementAtEpoch(EarthElements, T)

	planetHelio := heliocentricEcliptic(aP, eP, IP, LP, wbarP, OmegaP)
	earthHelio := heliocentricEcliptic(aE, eE, IE, LE, wbarE, OmegaE)

	geoEcl := sub(planetHelio, earthHelio)
	geoEq := eclipticToEquatorial(geoEcl)
	geoEq = scale(geoEq, auToMeters)

	return geoEcl, geoEq
}

func sunGeocentricEq(unixSec int64) (Vector, Vector) {
	T := julianCenturies(unixSec)
	aE, eE, IE, LE, wbarE, OmegaE := elementAtEpoch(EarthElements, T)

	earthHelio := heliocentricEcliptic(aE, eE, IE, LE, wbarE, OmegaE)

	sunGeoEcl := scale(earthHelio, -1)
	sunGeoEq := eclipticToEquatorial(sunGeoEcl)
	sunGeoEq = scale(sunGeoEq, auToMeters)

	return sunGeoEcl, sunGeoEq
}

func PlanetGeocentricECEF(el OrbitalElements, unixSec int64) (Vector, float64) {
	geoEcl, geoEq := planetGeocentricEq(el, unixSec)
	ecef := equatorialToECEF(geoEq, unixSec)

	ecefMeters := scale(ecef, auToMeters)
	dist := length(geoEcl) * auToMeters

	return ecefMeters, dist
}

func SunGeocentricECEF(unixSec int64) (Vector, float64) {
	sunGeoEcl, sunGeoEq := sunGeocentricEq(unixSec)
	ecef := equatorialToECEF(sunGeoEq, unixSec)

	ecefMeters := scale(ecef, auToMeters)
	dist := length(sunGeoEcl) * auToMeters

	return ecefMeters, dist
}

func MoonGeocentricECEF(unixSec int64) (Vector, float64) {
	distKm, moonGeoEq := moonGeocentricEq(unixSec)
	ecef := equatorialToECEF(moonGeoEq, unixSec)

	ecefMeters := scale(ecef, 1000)
	dist := distKm * 1000

	return ecefMeters, dist
}
