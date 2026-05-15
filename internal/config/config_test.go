package config

import (
	"fmt"
	"testing"
)

func TestGetArgValue(t *testing.T) {
	invalidValues := []string{"", "test", "test=", "test=test", "--test", "--test=", "--test=test"}
	for _, invalidValue := range invalidValues {
		value, found := getArgValue("test", invalidValue)

		if found {
			t.Errorf("Expected no result from %s, got %f", invalidValue, value)
		}
	}

	validValues := []float64{0, 1, 12, 12.3, 12.34}
	for _, validValue := range validValues {
		value, found := getArgValue("test", fmt.Sprintf("--test=%f", validValue))

		if !found {
			t.Errorf("Expected result from --test=%f, got none", validValue)
		}

		if value != validValue {
			t.Errorf("Expected %f, got %f", validValue, value)
		}
	}
}

func TestGetArgTimestampValue(t *testing.T) {
	invalidValues := []string{"", "test", "test=", "test=test", "--test", "--test=", "--test=test"}
	for _, invalidValue := range invalidValues {
		value, found := getArgTimestampValue("test", invalidValue)

		if found {
			t.Errorf("Expected no result from %s, got %d", invalidValue, value)
		}
	}

	const expectedValue = 1767225600
	validValues := []string{"2026-01-01", "2026-01-01T00:00:00"}
	for _, validValue := range validValues {
		value, found := getArgTimestampValue("test", fmt.Sprintf("--test=%s", validValue))

		if !found {
			t.Errorf("Expected result from --test=%s, got none", validValue)
		}

		if value != expectedValue {
			t.Errorf("Expected %d, got %d", expectedValue, value)
		}
	}
}
