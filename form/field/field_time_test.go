package field

import (
	"testing"
	"time"

	"github.com/xiriframework/xiri-go/types/distance"
	"github.com/xiriframework/xiri-go/types/language"
	"github.com/xiriframework/xiri-go/types/locale"
	"github.com/xiriframework/xiri-go/types/pressure"
	"github.com/xiriframework/xiri-go/types/timezone"
	"github.com/xiriframework/xiri-go/uicontext"
)

func TestTimeFieldMinMaxMidnightCalculation(t *testing.T) {
	// Create a UiContext with Europe/Vienna timezone (UTC+1 in winter, UTC+2 in summer)
	ctx := &uicontext.UiContext{
		Timezone: timezone.EuropeVienna,
		Lang:     language.Deutsch,
		Locale:   locale.De,
		Distance: distance.Kilometer,
		Pressure: pressure.Bar,
	}

	// Create a TimeField with min=-7 (7 days ago) and max=0 (today)
	field := NewTimeField("testtime", "test.time", true, 0)
	minVal := int64(-7)
	maxVal := int64(0)
	field.Min = &minVal
	field.Max = &maxVal

	// Export for frontend
	result := field.ExportForFrontend(ctx, nil)

	// Verify min and max are set
	if result["min"] == nil {
		t.Fatal("Expected min to be set")
	}
	if result["max"] == nil {
		t.Fatal("Expected max to be set")
	}

	minTimestamp := result["min"].(int64)
	maxTimestamp := result["max"].(int64)

	// Load Vienna timezone
	loc, err := time.LoadLocation("Europe/Vienna")
	if err != nil {
		t.Fatalf("Failed to load timezone: %v", err)
	}

	// Convert timestamps to time in Vienna timezone
	minTime := time.Unix(minTimestamp, 0).In(loc)
	maxTime := time.Unix(maxTimestamp, 0).In(loc)

	// Verify both are at midnight (00:00:00)
	if minTime.Hour() != 0 || minTime.Minute() != 0 || minTime.Second() != 0 {
		t.Errorf("Expected min to be at midnight, got %s", minTime.Format("15:04:05"))
	}
	if maxTime.Hour() != 0 || maxTime.Minute() != 0 || maxTime.Second() != 0 {
		t.Errorf("Expected max to be at midnight, got %s", maxTime.Format("15:04:05"))
	}

	// Verify min is approximately 7 days before max
	// Note: Due to DST transitions, the difference may not be exactly 7*24 hours
	// For example, when falling back from DST, there's an extra hour
	expectedDiffMin := int64(7 * 24 * 60 * 60)   // 7 days in seconds (no DST)
	expectedDiffMax := int64(7*24*60*60 + 60*60) // 7 days + 1 hour (DST fall back)
	actualDiff := maxTimestamp - minTimestamp
	if actualDiff < expectedDiffMin || actualDiff > expectedDiffMax {
		t.Errorf("Expected min to be 7 days before max (allowing for DST), got diff of %d seconds (expected %d-%d)", actualDiff, expectedDiffMin, expectedDiffMax)
	}

	t.Logf("Min time (7 days ago at midnight): %s", minTime.Format("2006-01-02 15:04:05 MST"))
	t.Logf("Max time (today at midnight): %s", maxTime.Format("2006-01-02 15:04:05 MST"))
}

func TestTimeFieldMinMaxAbsoluteTimestamp(t *testing.T) {
	// Create a UiContext
	ctx := &uicontext.UiContext{
		Timezone: timezone.EuropeVienna,
		Lang:     language.Deutsch,
		Locale:   locale.De,
		Distance: distance.Kilometer,
		Pressure: pressure.Bar,
	}

	// Create a TimeField with absolute timestamps (outside -10000 to 10000 range)
	field := NewTimeField("testtime", "test.time", true, 0)
	minVal := int64(1704067200) // 2024-01-01 00:00:00 UTC
	maxVal := int64(1735689600) // 2025-01-01 00:00:00 UTC
	field.Min = &minVal
	field.Max = &maxVal

	// Export for frontend
	result := field.ExportForFrontend(ctx, nil)

	// Verify min and max are unchanged (absolute timestamps)
	if result["min"].(int64) != minVal {
		t.Errorf("Expected min to be unchanged, got %d instead of %d", result["min"].(int64), minVal)
	}
	if result["max"].(int64) != maxVal {
		t.Errorf("Expected max to be unchanged, got %d instead of %d", result["max"].(int64), maxVal)
	}

	t.Logf("Min timestamp (absolute): %d", result["min"].(int64))
	t.Logf("Max timestamp (absolute): %d", result["max"].(int64))
}

func TestTimeFieldMinMaxDifferentTimezones(t *testing.T) {
	testCases := []struct {
		name     string
		timezone timezone.Timezone
		ianaName string
	}{
		{"Vienna", timezone.EuropeVienna, "Europe/Vienna"},
		{"Berlin", timezone.EuropeBerlin, "Europe/Berlin"},
		{"London", timezone.EuropeLondon, "Europe/London"},
		{"Madrid", timezone.EuropeMadrid, "Europe/Madrid"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a UiContext with specific timezone
			ctx := &uicontext.UiContext{
				Timezone: tc.timezone,
				Lang:     language.Deutsch,
				Locale:   locale.De,
				Distance: distance.Kilometer,
				Pressure: pressure.Bar,
			}

			// Create a TimeField with max=0 (today at midnight)
			field := NewTimeField("testtime", "test.time", true, 0)
			maxVal := int64(0)
			field.Max = &maxVal

			// Export for frontend
			result := field.ExportForFrontend(ctx, nil)

			maxTimestamp := result["max"].(int64)

			// Load the timezone
			loc, err := time.LoadLocation(tc.ianaName)
			if err != nil {
				t.Fatalf("Failed to load timezone %s: %v", tc.ianaName, err)
			}

			// Convert to time in that timezone
			maxTime := time.Unix(maxTimestamp, 0).In(loc)

			// Verify it's at midnight in that timezone
			if maxTime.Hour() != 0 || maxTime.Minute() != 0 || maxTime.Second() != 0 {
				t.Errorf("Expected midnight in %s timezone, got %s", tc.ianaName, maxTime.Format("15:04:05 MST"))
			}

			t.Logf("Midnight in %s: %s", tc.ianaName, maxTime.Format("2006-01-02 15:04:05 MST"))
		})
	}
}
