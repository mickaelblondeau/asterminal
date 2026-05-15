package converter

import (
	"math"
	"testing"
)

type NormalizedRadiusTestCase struct {
	expected                 float64
	radius, distance, fovDeg float64
}

func TestNormalizedRadius(t *testing.T) {
	testCases := []NormalizedRadiusTestCase{
		{2.980902, 100, 1, 60},
		{5.961804, 100, 1, 30},
	}

	for _, testCase := range testCases {
		v := NormalizedRadius(testCase.radius, testCase.distance, testCase.fovDeg)

		if math.Abs(v-testCase.expected) > 1e-6 {
			t.Errorf("Expected %f, got: %f", testCase.expected, v)
		}
	}
}
