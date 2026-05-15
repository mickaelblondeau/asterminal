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
	obliquity  = 23.43928
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

func heliocentricEcliptic(a, e, I, L, wbar, Omega float64) Vector {
	w := wbar - Omega
	M := degToRad(math.Mod(L-wbar, 360))
	E := solveKepler(M, e)

	xOrb := a * (math.Cos(E) - e)
	yOrb := a * math.Sqrt(1-e*e) * math.Sin(E)

	wR := degToRad(w)
	IR := degToRad(I)
	OR := degToRad(Omega)

	cosW, sinW := math.Cos(wR), math.Sin(wR)
	cosI, sinI := math.Cos(IR), math.Sin(IR)
	cosO, sinO := math.Cos(OR), math.Sin(OR)

	x := (cosO*cosW-sinO*sinW*cosI)*xOrb + (-cosO*sinW-sinO*cosW*cosI)*yOrb
	y := (sinO*cosW+cosO*sinW*cosI)*xOrb + (-sinO*sinW+cosO*cosW*cosI)*yOrb
	z := (sinW*sinI)*xOrb + (cosW*sinI)*yOrb

	return Vector{x, y, z}
}

func eclipticToEquatorial(v Vector) Vector {
	eps := degToRad(obliquity)
	cosE, sinE := math.Cos(eps), math.Sin(eps)

	return Vector{
		X: v.X,
		Y: cosE*v.Y - sinE*v.Z,
		Z: sinE*v.Y + cosE*v.Z,
	}
}

func equatorialToECEF(v Vector, unixSec int64) Vector {
	theta := gmst(unixSec)

	return Vector{
		X: v.X*math.Cos(theta) + v.Y*math.Sin(theta),
		Y: -v.X*math.Sin(theta) + v.Y*math.Cos(theta),
		Z: v.Z,
	}
}

func PlanetGeocentricECEF(el OrbitalElements, unixSec int64) (pos Vector, distMeters float64) {
	T := julianCenturies(unixSec)
	aP, eP, IP, LP, wbarP, OmegaP := elementAtEpoch(el, T)
	aE, eE, IE, LE, wbarE, OmegaE := elementAtEpoch(EarthElements, T)

	planetHelio := heliocentricEcliptic(aP, eP, IP, LP, wbarP, OmegaP)
	earthHelio := heliocentricEcliptic(aE, eE, IE, LE, wbarE, OmegaE)

	geoEcl := sub(planetHelio, earthHelio)
	geoEq := eclipticToEquatorial(geoEcl)
	ecef := equatorialToECEF(geoEq, unixSec)

	ecefMeters := scale(ecef, auToMeters)
	dist := length(geoEcl) * auToMeters

	return ecefMeters, dist
}

func SunGeocentricECEF(unixSec int64) (pos Vector, distMeters float64) {
	T := julianCenturies(unixSec)
	aE, eE, IE, LE, wbarE, OmegaE := elementAtEpoch(EarthElements, T)

	earthHelio := heliocentricEcliptic(aE, eE, IE, LE, wbarE, OmegaE)

	sunGeoEcl := scale(earthHelio, -1)
	sunGeoEq := eclipticToEquatorial(sunGeoEcl)
	ecef := equatorialToECEF(sunGeoEq, unixSec)

	ecefMeters := scale(ecef, auToMeters)
	dist := length(sunGeoEcl) * auToMeters

	return ecefMeters, dist
}
