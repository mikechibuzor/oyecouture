package orders

import (
	"math"
	"testing"
)

func TestSanitizeMeasurementsKeepsOnlyKnownValues(t *testing.T) {
	t.Parallel()

	measurements := sanitizeMeasurements(&Measurements{
		Profile: "native",
		Unit:    "in",
		Values: map[string]float64{
			"chest":        42.5,
			"cap":          23,
			"unknown":      10,
			"waist":        -2,
			"top-length":   201,
			"round-sleeve": math.NaN(),
		},
		Notes: "  leave room through the chest  ",
	})

	if measurements == nil {
		t.Fatal("sanitizeMeasurements() returned nil")
	}
	if measurements.Profile != "native" || measurements.Unit != "in" {
		t.Fatalf("profile/unit = %q/%q, want native/in", measurements.Profile, measurements.Unit)
	}
	if len(measurements.Values) != 2 || measurements.Values["chest"] != 42.5 || measurements.Values["cap"] != 23 {
		t.Fatalf("values = %#v, want only chest and cap", measurements.Values)
	}
	if measurements.Notes != "leave room through the chest" {
		t.Fatalf("notes = %q, want trimmed notes", measurements.Notes)
	}
}

func TestSanitizeMeasurementsDropsUntrustedStructuredPayload(t *testing.T) {
	t.Parallel()

	measurements := sanitizeMeasurements(&Measurements{
		Profile: "unknown",
		Unit:    "cm",
		Values:  map[string]float64{"chest": 42},
	})

	if measurements != nil {
		t.Fatalf("sanitizeMeasurements() = %#v, want nil", measurements)
	}
}
