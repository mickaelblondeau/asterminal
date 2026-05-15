package converter

import (
	"math"
)

const (
	moonLpC0           = 218.3164477
	moonLpC1           = 481267.88123421
	moonLpC2           = -0.0015786
	moonLpC3           = 1.0 / 538841
	moonLpC4           = -1.0 / 65194000
	moonDC0            = 297.8501921
	moonDC1            = 445267.1114034
	moonDC2            = -0.0018819
	moonDC3            = 1.0 / 545868
	moonDC4            = -1.0 / 113065000
	moonMC0            = 357.5291092
	moonMC1            = 35999.0502909
	moonMC2            = -0.0001536
	moonMC3            = 1.0 / 24490000
	moonMpC0           = 134.9633964
	moonMpC1           = 477198.8675055
	moonMpC2           = 0.0087414
	moonMpC3           = 1.0 / 69699
	moonMpC4           = -1.0 / 14712000
	moonFC0            = 93.2720950
	moonFC1            = 483202.0175233
	moonFC2            = -0.0036539
	moonFC3            = -1.0 / 3526000
	moonFC4            = 1.0 / 863310000
	moonEccentricityC1 = 0.002516
	moonEccentricityC2 = 0.0000074
	moonA1C0           = 119.75
	moonA1C1           = 131.849
	moonA2C0           = 53.09
	moonA2C1           = 479264.290
	moonA3C0           = 313.45
	moonA3C1           = 481266.484
	moonSumLCorrA1     = 3958
	moonSumLCorrLpF    = 1962
	moonSumLCorrA2     = 318
	moonSumBCorrLp     = -2235
	moonSumBCorrA3     = 382
	moonSumBCorrA1F    = 175
	moonSumBCorrLpMp1  = 127
	moonSumBCorrLpMp2  = -115
	moonMeanDistanceKm = 385000.56
	moonSumRScale      = 1.0 / 1000
	moonSumScale       = 1.0 / 1e6
)

func norm360(x float64) float64 {
	x = math.Mod(x, 360)

	if x < 0 {
		x += 360
	}

	return x
}

func MoonGeocentricECEF(unixSec int64) (pos Vector, distMeters float64) {
	T := julianCenturies(unixSec)
	toR := math.Pi / 180

	LpDeg := norm360(moonLpC0 + moonLpC1*T + moonLpC2*T*T + moonLpC3*T*T*T + moonLpC4*T*T*T*T)
	DDeg := norm360(moonDC0 + moonDC1*T + moonDC2*T*T + moonDC3*T*T*T + moonDC4*T*T*T*T)
	MDeg := norm360(moonMC0 + moonMC1*T + moonMC2*T*T + moonMC3*T*T*T)
	MpDeg := norm360(moonMpC0 + moonMpC1*T + moonMpC2*T*T + moonMpC3*T*T*T + moonMpC4*T*T*T*T)
	FDeg := norm360(moonFC0 + moonFC1*T + moonFC2*T*T + moonFC3*T*T*T + moonFC4*T*T*T*T)

	Lp := LpDeg * toR
	D := DDeg * toR
	M := MDeg * toR
	Mp := MpDeg * toR
	F := FDeg * toR

	type term struct{ D, M, Mp, F, Sl, Sr float64 }
	terms := []term{
		{0, 0, 1, 0, 6288774, -20905355},
		{2, 0, -1, 0, 1274027, -3699111},
		{2, 0, 0, 0, 658314, -2955968},
		{0, 0, 2, 0, 213618, -569925},
		{0, 1, 0, 0, -185116, 48888},
		{0, 0, 0, 2, -114332, -3149},
		{2, 0, -2, 0, 58793, 246158},
		{2, -1, -1, 0, 57066, -152138},
		{2, 0, 1, 0, 53322, -170733},
		{2, -1, 0, 0, 45758, -204586},
		{0, 1, -1, 0, -40923, -129620},
		{1, 0, 0, 0, -34720, 108743},
		{0, 1, 1, 0, -30383, 104755},
		{2, 0, 0, -2, 15327, 10321},
		{0, 0, 1, 2, -12528, 0},
		{0, 0, 1, -2, 10980, 79661},
		{4, 0, -1, 0, 10675, -34782},
		{0, 0, 3, 0, 10034, -23210},
		{4, 0, -2, 0, 8548, -21636},
		{2, 1, -1, 0, -7888, 24208},
		{2, 1, 0, 0, -6766, 30824},
		{1, 0, -1, 0, -5163, -8379},
		{1, 1, 0, 0, 4987, -16675},
		{2, -1, 1, 0, 4036, -12831},
		{2, 0, 2, 0, 3994, -10445},
		{4, 0, 0, 0, 3861, -11650},
		{2, 0, -3, 0, 3665, 14403},
		{0, 1, -2, 0, -2689, -7003},
		{2, 0, -1, 2, -2602, 0},
		{2, -1, -2, 0, 2390, 10056},
		{1, 0, 1, 0, -2348, 6322},
		{2, -2, 0, 0, 2236, -9884},
		{0, 1, 2, 0, -2120, 5751},
		{0, 2, 0, 0, -2069, 0},
		{2, -2, -1, 0, 2048, -4950},
		{2, 0, 1, -2, -1773, 4130},
		{2, 0, 0, 2, -1595, 0},
		{4, -1, -1, 0, 1215, -3958},
		{0, 0, 2, 2, -1110, 0},
		{3, 0, -1, 0, -892, 3258},
		{2, 1, 1, 0, -810, 2616},
		{4, -1, -2, 0, 759, -1897},
		{0, 2, -1, 0, -713, -2117},
		{2, 2, -1, 0, -700, 2354},
		{2, 1, -2, 0, 691, 0},
		{2, -1, 0, -2, 596, 0},
		{4, 0, 1, 0, 549, -1423},
		{0, 0, 4, 0, 537, -1117},
		{4, -1, 0, 0, 520, -1571},
		{1, 0, -2, 0, -487, -1739},
		{2, 1, 0, -2, -399, 0},
		{0, 0, 2, -2, -381, -4421},
		{1, 1, 1, 0, 351, 0},
		{3, 0, -2, 0, -340, 0},
		{4, 0, -3, 0, 330, 0},
		{2, -1, 2, 0, 327, 0},
		{0, 2, 1, 0, -323, 1165},
		{1, 1, -1, 0, 299, 0},
		{2, 0, 3, 0, 294, 0},
	}

	type bterm struct{ D, M, Mp, F, Sb float64 }
	bterms := []bterm{
		{0, 0, 0, 1, 5128122},
		{0, 0, 1, 1, 280602},
		{0, 0, 1, -1, 277693},
		{2, 0, 0, -1, 173237},
		{2, 0, -1, 1, 55413},
		{2, 0, -1, -1, 46271},
		{2, 0, 0, 1, 32573},
		{0, 0, 2, 1, 17198},
		{2, 0, 1, -1, 9266},
		{0, 0, 2, -1, 8822},
		{2, -1, 0, -1, 8216},
		{2, 0, -2, -1, 4324},
		{2, 0, 1, 1, 4200},
		{2, 1, 0, -1, -3359},
		{2, -1, -1, 1, 2463},
		{2, -1, 0, 1, 2211},
		{2, -1, -1, -1, 2065},
		{0, 1, -1, -1, -1870},
		{4, 0, -1, -1, 1828},
		{0, 1, 0, 1, -1794},
		{0, 0, 0, 3, -1749},
		{0, 1, -1, 1, -1565},
		{1, 0, 0, 1, -1491},
		{0, 1, 1, 1, -1475},
		{0, 1, 1, -1, -1410},
		{0, 1, 0, -1, -1344},
		{1, 0, 0, -1, -1335},
		{0, 0, 3, 1, 1107},
		{4, 0, 0, -1, 1021},
		{4, 0, -1, 1, 833},
	}

	E := 1 - moonEccentricityC1*T - moonEccentricityC2*T*T

	var sumL, sumR, sumB float64

	for _, t := range terms {
		arg := t.D*D + t.M*M + t.Mp*Mp + t.F*F
		eCorr := 1.0
		absM := math.Round(math.Abs(t.M))

		switch absM {
		case 1:
			eCorr = E
		case 2:
			eCorr = E * E
		}

		sumL += eCorr * t.Sl * math.Sin(arg)
		sumR += eCorr * t.Sr * math.Cos(arg)
	}

	for _, t := range bterms {
		arg := t.D*D + t.M*M + t.Mp*Mp + t.F*F
		eCorr := 1.0
		if math.Abs(t.M) == 1 {
			eCorr = E
		} else if math.Abs(t.M) == 2 {
			eCorr = E * E
		}
		sumB += eCorr * t.Sb * math.Sin(arg)
	}

	A1 := moonA1C0 + moonA1C1*T
	A2 := moonA2C0 + moonA2C1*T
	A3 := moonA3C0 + moonA3C1*T

	sumL += moonSumLCorrA1*math.Sin(degToRad(A1)) +
		moonSumLCorrLpF*math.Sin(Lp-F) +
		moonSumLCorrA2*math.Sin(degToRad(A2))

	sumB += moonSumBCorrLp*math.Sin(Lp) +
		moonSumBCorrA3*math.Sin(degToRad(A3)) +
		moonSumBCorrA1F*math.Sin(degToRad(A1)-F) +
		moonSumBCorrA1F*math.Sin(degToRad(A1)+F) +
		moonSumBCorrLpMp1*math.Sin(Lp-Mp) +
		moonSumBCorrLpMp2*math.Sin(Lp+Mp)

	lambda := norm360(LpDeg + sumL/1e6)
	lambda *= toR

	beta := sumB * moonSumScale * toR
	distKm := moonMeanDistanceKm + sumR*moonSumRScale

	cosB := math.Cos(beta)
	ecliptic := Vector{
		X: distKm * cosB * math.Cos(lambda),
		Y: distKm * cosB * math.Sin(lambda),
		Z: distKm * math.Sin(beta),
	}

	eq := eclipticToEquatorial(ecliptic)
	ecef := equatorialToECEF(eq, unixSec)
	ecefMeters := scale(ecef, 1000)
	dist := distKm * 1000

	return ecefMeters, dist
}
