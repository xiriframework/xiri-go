package field

import (
	"testing"
	"time"

	"github.com/xiriframework/xiri-go/types/distance"
	"github.com/xiriframework/xiri-go/types/language"
	"github.com/xiriframework/xiri-go/types/locale"
	"github.com/xiriframework/xiri-go/types/pressure"
	"github.com/xiriframework/xiri-go/types/timezone"
	"github.com/xiriframework/xiri-go/component/core"
)

func TestTimeFieldMinMaxMidnightCalculation(t *testing.T) {
	// Create a UiContext with Europe/Vienna timezone (UTC+1 in winter, UTC+2 in summer)
	ctx := &core.UiContext{
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
	ctx := &core.UiContext{
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

func TestTimeFieldParseYearMonth(t *testing.T) {
	field := NewYearMonthField("month", "test.month", false, 0)

	parsed, err := field.Parse("2026-04")
	if err != nil {
		t.Fatalf("Parse(\"2026-04\") returned error: %v", err)
	}

	ts, ok := parsed.(int64)
	if !ok {
		t.Fatalf("Parse returned %T, expected int64", parsed)
	}

	expected := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC).Unix()
	if ts != expected {
		t.Errorf("Parse(\"2026-04\") = %d, expected %d (%s)", ts, expected, time.Unix(expected, 0).UTC().Format(time.RFC3339))
	}
}

func TestYearMonthFieldExportType(t *testing.T) {
	ctx := &core.UiContext{
		Timezone: timezone.EuropeVienna,
		Lang:     language.Deutsch,
		Locale:   locale.De,
		Distance: distance.Kilometer,
		Pressure: pressure.Bar,
	}

	field := NewYearMonthField("month", "test.month", false, 0)
	result := field.ExportForFrontend(ctx, nil)

	if result["type"] != "yearmonth" {
		t.Errorf("Expected type \"yearmonth\", got %v", result["type"])
	}
	if result["subtype"] != "yearmonth" {
		t.Errorf("Expected subtype \"yearmonth\", got %v", result["subtype"])
	}
}

func TestTimeFieldParseYearMonthInvalidStrings(t *testing.T) {
	field := NewYearMonthField("month", "test.month", false, 0)

	cases := []string{
		"not-a-date",
		"2026/04",
		"2026-13",       // invalid month
		"2026-00",       // invalid month
		"26-04",         // 2-digit year not supported
		"April 2026",
		"",
	}

	for _, raw := range cases {
		_, err := field.Parse(raw)
		if err == nil {
			t.Errorf("Parse(%q) expected error, got nil", raw)
		}
	}
}

func TestTimeFieldParseYearMonthAcceptsOtherFormats(t *testing.T) {
	field := NewYearMonthField("month", "test.month", false, 0)

	// "2006-01-02" still works (full ISO date) — picks up the day component
	parsed, err := field.Parse("2026-04-15")
	if err != nil {
		t.Fatalf("Parse(\"2026-04-15\") returned unexpected error: %v", err)
	}
	ts, _ := parsed.(int64)
	expected := time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC).Unix()
	if ts != expected {
		t.Errorf("Parse(\"2026-04-15\") = %d, expected %d", ts, expected)
	}
}

func TestYearMonthFieldValidateRequired(t *testing.T) {
	required := NewYearMonthField("month", "test.month", true, 0)

	// nil with required → error
	if err := required.Validate(nil); err == nil {
		t.Error("Validate(nil) on required field expected error, got nil")
	}

	// valid timestamp → ok
	ts := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC).Unix()
	required.AllowFuture = true
	required.AllowPast = true
	if err := required.Validate(ts); err != nil {
		t.Errorf("Validate(valid timestamp) returned unexpected error: %v", err)
	}
}

func TestYearMonthFieldValidateInvalidType(t *testing.T) {
	field := NewYearMonthField("month", "test.month", false, 0)

	if err := field.Validate("string-not-allowed-here"); err == nil {
		t.Error("Validate(string) expected error, got nil")
	}
	if err := field.Validate(struct{}{}); err == nil {
		t.Error("Validate(struct) expected error, got nil")
	}
}

// Verifies that a "YYYY-MM" string parses to the same Unix timestamp regardless of
// the host machine's local timezone — Go's time.Parse uses UTC for the year-month layout.
func TestYearMonthFieldParseTimezoneStable(t *testing.T) {
	field := NewYearMonthField("month", "test.month", false, 0)

	// Switch goroutine-local time settings via TZ env / time.Local
	originalLocal := time.Local
	defer func() { time.Local = originalLocal }()

	timezones := []string{"Europe/Vienna", "America/New_York", "Asia/Tokyo", "UTC"}
	expected := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC).Unix()

	for _, tz := range timezones {
		loc, err := time.LoadLocation(tz)
		if err != nil {
			t.Fatalf("LoadLocation(%q) failed: %v", tz, err)
		}
		time.Local = loc

		parsed, err := field.Parse("2026-04")
		if err != nil {
			t.Fatalf("Parse(\"2026-04\") in %s returned error: %v", tz, err)
		}
		ts, _ := parsed.(int64)
		if ts != expected {
			t.Errorf("Parse(\"2026-04\") in %s = %d, expected %d (UTC-stable)", tz, ts, expected)
		}
	}
}

// Verifies ExportForFrontend produces the same type/subtype regardless of the user's locale/timezone.
func TestYearMonthFieldExportLocaleStable(t *testing.T) {
	cases := []struct {
		name string
		ctx  *core.UiContext
	}{
		{"Vienna/De", &core.UiContext{Timezone: timezone.EuropeVienna, Lang: language.Deutsch, Locale: locale.De, Distance: distance.Kilometer, Pressure: pressure.Bar}},
		{"London/EnGB", &core.UiContext{Timezone: timezone.EuropeLondon, Lang: language.Englisch, Locale: locale.EnGB, Distance: distance.Kilometer, Pressure: pressure.Bar}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			field := NewYearMonthField("month", "test.month", false, 0)
			result := field.ExportForFrontend(tc.ctx, nil)

			if result["type"] != "yearmonth" {
				t.Errorf("Expected type \"yearmonth\" in %s, got %v", tc.name, result["type"])
			}
			if result["subtype"] != "yearmonth" {
				t.Errorf("Expected subtype \"yearmonth\" in %s, got %v", tc.name, result["subtype"])
			}
		})
	}
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
			ctx := &core.UiContext{
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
