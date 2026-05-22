package converter

import (
	"fmt"
	"math"
	"testing"
	"time"
)

type geocentricTestCase struct {
	time     string
	expected Vector
}

const (
	diff     = 0.02
	diffMoon = 0.65
)

func checkVectorError(value, expected Vector, diff float64, t string) string {
	errorDist := math.Sqrt(math.Pow(expected.X-value.X, 2) + math.Pow(expected.Y-value.Y, 2) + math.Pow(expected.Z-value.Z, 2))
	totalDist := math.Sqrt(math.Pow(expected.X, 2) + math.Pow(expected.Y, 2) + math.Pow(expected.Z, 2))
	diffPercent := (errorDist / totalDist) * 100

	if diffPercent > diff {
		return fmt.Sprintf("[%s] Vector error too high: %.2f%% (Allowed: %.2f%%)", t, diffPercent, diff)
	}

	return ""
}

func TestPlanetGeocentricEq(t *testing.T) {
	testCases := []geocentricTestCase{
		{"2026-01-01T00:00:00", Vector{7.702213291044337e+07, -3.522672220630372e+08, -5.588644593033329e+06}},
		{"2026-02-01T00:00:00", Vector{2.109957440539093e+08, -2.864863632116151e+08, -6.459789913710088e+06}},
		{"2026-03-01T00:00:00", Vector{2.982227481880773e+08, -1.835666621149857e+08, -6.681025161379695e+06}},
		{"2026-04-01T00:00:00", Vector{3.405415655783519e+08, -4.326704212551549e+07, -6.254030005801909e+06}},
		{"2026-05-01T00:00:00", Vector{3.228454551484687e+08, 9.265572461992887e+07, -5.204188303183079e+06}},
		{"2026-06-01T00:00:00", Vector{2.524771564743362e+08, 2.073699162303412e+08, -3.583063701831698e+06}},
	}

	for _, testCase := range testCases {
		time, _ := time.ParseInLocation("2006-01-02T15:04:05", testCase.time, time.UTC)
		expected := scale(testCase.expected, 1e+3)
		expected = eclipticToEquatorial(expected)

		_, value := planetGeocentricEq(MarsElements, time.Unix())

		if err := checkVectorError(value, expected, diff, testCase.time); err != "" {
			t.Error(err)
		}
	}
}

func TestSunGeocentricEq(t *testing.T) {
	testCases := []geocentricTestCase{
		{"2026-01-01T00:00:00", Vector{2.607213844816194e+07, -1.447746738210197e+08, 8.892861905746162e+03}},
		{"2026-02-01T00:00:00", Vector{9.817500099241482e+07, -1.099446604634764e+08, 6.298976310864091e+03}},
		{"2026-03-01T00:00:00", Vector{1.393090841151139e+08, -5.058478056218451e+07, 2.469433247771114e+03}},
		{"2026-04-01T00:00:00", Vector{1.467588947990943e+08, 2.829438070141354e+07, -2.658867232609540e+03}},
		{"2026-05-01T00:00:00", Vector{1.149707719187170e+08, 9.743243572516793e+07, -6.843554095692933e+03}},
		{"2026-06-01T00:00:00", Vector{5.144685903194299e+07, 1.426831420847642e+08, -9.270717556029558e+03}},
	}

	for _, testCase := range testCases {
		time, _ := time.ParseInLocation("2006-01-02T15:04:05", testCase.time, time.UTC)
		expected := scale(testCase.expected, 1e+3)
		expected = eclipticToEquatorial(expected)

		_, value := sunGeocentricEq(time.Unix())

		if err := checkVectorError(value, expected, diff, testCase.time); err != "" {
			t.Error(err)
		}
	}
}

func TestMoonGeocentricEq(t *testing.T) {
	testCases := []geocentricTestCase{
		{"2026-01-01T00:00:00", Vector{1.443257274919273e+05, 3.293958308956202e+05, 3.175297634320188e+04}},
		{"2026-02-01T00:00:00", Vector{-1.813509759236951e+05, 3.201060471523177e+05, 2.121226595263975e+04}},
		{"2026-03-01T00:00:00", Vector{-2.341845363633380e+05, 2.919083048613686e+05, 1.721607887346357e+04}},
		{"2026-04-01T00:00:00", Vector{-3.894680455007774e+05, 1.289648175008370e+04, -1.194896872279250e+04}},
		{"2026-05-01T00:00:00", Vector{-3.380133989560704e+05, -2.124702485574115e+05, -2.903667898111700e+04}},
		{"2026-06-01T00:00:00", Vector{-9.040360754971113e+04, -3.946802390341145e+05, -3.432130324998641e+04}},
	}

	for _, testCase := range testCases {
		time, _ := time.ParseInLocation("2006-01-02T15:04:05", testCase.time, time.UTC)
		expected := testCase.expected
		expected = eclipticToEquatorial(expected)

		_, value := moonGeocentricEq(time.Unix())

		if err := checkVectorError(value, expected, diffMoon, testCase.time); err != "" {
			t.Error(err)
		}
	}
}
