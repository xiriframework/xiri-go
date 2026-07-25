package field

import (
	"strings"
	"testing"
	"time"

	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/types/timezone"
)

func TestTimeRangeField_Validate_Required(t *testing.T) {
	f := NewTimeRangeField("range", "RANGE", true)
	err := f.Validate(nil)
	if err == nil {
		t.Fatal("expected error for nil on required timerange field")
	}
	if !strings.Contains(err.Error(), "required") {
		t.Errorf("expected 'required' in error, got: %v", err)
	}
}

func TestTimeRangeField_Validate_Optional_Nil(t *testing.T) {
	f := NewTimeRangeField("range", "RANGE", false)
	if err := f.Validate(nil); err != nil {
		t.Fatalf("unexpected error for nil on optional field: %v", err)
	}
}

func TestTimeRangeField_Validate_InvalidType(t *testing.T) {
	f := NewTimeRangeField("range", "RANGE", true)
	err := f.Validate("not a timerange value")
	if err == nil {
		t.Fatal("expected error for non-TimeRangeValue input")
	}
	if !strings.Contains(err.Error(), "invalid timerange value type") {
		t.Errorf("expected 'invalid timerange value type' in error, got: %v", err)
	}
}

func TestTimeRangeField_Validate_ZeroDates(t *testing.T) {
	f := NewTimeRangeField("range", "RANGE", true)
	tr := &TimeRangeValue{
		Start: time.Time{}, // zero
		End:   time.Now(),
	}
	err := f.Validate(tr)
	if err == nil {
		t.Fatal("expected error for zero start date")
	}
	if !strings.Contains(err.Error(), "zero dates") {
		t.Errorf("expected 'zero dates' in error, got: %v", err)
	}
}

func TestTimeRangeField_Validate_StartAfterEnd(t *testing.T) {
	f := NewTimeRangeField("range", "RANGE", true)
	now := time.Now()
	tr := &TimeRangeValue{
		Start: now.Add(2 * time.Hour),
		End:   now,
	}
	err := f.Validate(tr)
	if err == nil {
		t.Fatal("expected error when start is after end")
	}
}

func TestTimeRangeField_Validate_Valid(t *testing.T) {
	f := NewTimeRangeField("range", "RANGE", true)
	now := time.Now()
	tr := &TimeRangeValue{
		Start: now,
		End:   now.Add(24 * time.Hour),
	}
	if err := f.Validate(tr); err != nil {
		t.Fatalf("unexpected error for valid range: %v", err)
	}
}

func TestTimeRangeField_Validate_AllowSingleDay(t *testing.T) {
	f := NewTimeRangeField("range", "RANGE", true)
	f.AllowSingleDay = true
	now := time.Now()
	tr := &TimeRangeValue{
		Start: now,
		End:   now, // same as start
	}
	if err := f.Validate(tr); err != nil {
		t.Fatalf("unexpected error when AllowSingleDay=true and start==end: %v", err)
	}
}

func TestTimeRangeField_Validate_DisallowSingleDay(t *testing.T) {
	f := NewTimeRangeField("range", "RANGE", true)
	f.AllowSingleDay = false
	now := time.Now()
	tr := &TimeRangeValue{
		Start: now,
		End:   now, // same as start
	}
	err := f.Validate(tr)
	if err == nil {
		t.Fatal("expected error when AllowSingleDay=false and start==end")
	}
	if !strings.Contains(err.Error(), "start must be before end") {
		t.Errorf("expected 'start must be before end' in error, got: %v", err)
	}
}

func TestTimeRangeField_Parse_Nil_WithDefault(t *testing.T) {
	f := NewTimeRangeFieldWithDefault("range", "RANGE", false, 7)
	result, err := f.Parse(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil default value")
	}
	tr, ok := result.(*TimeRangeValue)
	if !ok {
		t.Fatalf("expected *TimeRangeValue, got %T", result)
	}
	// Default should be 7 days back from now
	if tr.Start.IsZero() || tr.End.IsZero() {
		t.Error("expected non-zero start and end in default value")
	}
	diff := tr.End.Sub(tr.Start)
	// Allow some tolerance (should be roughly 7 days)
	if diff < 6*24*time.Hour || diff > 8*24*time.Hour {
		t.Errorf("expected ~7 day range, got %v", diff)
	}
}

func TestTimeRangeField_Parse_UnixTimestamps(t *testing.T) {
	f := NewTimeRangeField("range", "RANGE", true)
	input := map[string]interface{}{
		"start": float64(1640000000),
		"end":   float64(1640100000),
	}
	result, err := f.Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tr, ok := result.(*TimeRangeValue)
	if !ok {
		t.Fatalf("expected *TimeRangeValue, got %T", result)
	}
	expectedStart := time.Unix(1640000000, 0)
	expectedEnd := time.Unix(1640100000, 0)
	if !tr.Start.Equal(expectedStart) {
		t.Errorf("expected start %v, got %v", expectedStart, tr.Start)
	}
	if !tr.End.Equal(expectedEnd) {
		t.Errorf("expected end %v, got %v", expectedEnd, tr.End)
	}
}

func TestTimeRangeField_Parse_InvalidInput(t *testing.T) {
	f := NewTimeRangeField("range", "RANGE", true)
	_, err := f.Parse("not a map")
	if err == nil {
		t.Fatal("expected error for non-map input")
	}
	if !strings.Contains(err.Error(), "expects map") {
		t.Errorf("expected 'expects map' in error, got: %v", err)
	}
}

func TestTimeRangeField_BindValue(t *testing.T) {
	f := NewTimeRangeField("range", "RANGE", true)
	input := map[string]interface{}{
		"start": float64(1640000000),
		"end":   float64(1640100000),
	}
	err := f.BindValue(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Value == nil {
		t.Fatal("expected non-nil Value after BindValue")
	}
	expectedStart := time.Unix(1640000000, 0)
	expectedEnd := time.Unix(1640100000, 0)
	if !f.Value.Start.Equal(expectedStart) {
		t.Errorf("expected start %v, got %v", expectedStart, f.Value.Start)
	}
	if !f.Value.End.Equal(expectedEnd) {
		t.Errorf("expected end %v, got %v", expectedEnd, f.Value.End)
	}
}

// TestTimeRangeField_Export_DayOffsetAnchoredToLocalMidnight verifies that a relative
// Min/Max day offset is anchored to local midnight in the user's timezone instead of
// "now + n*86400" (#15). Anchoring to the current wall clock makes the boundary drift
// with the time of day and breaks across DST transitions.
func TestTimeRangeField_Export_DayOffsetAnchoredToLocalMidnight(t *testing.T) {
	ctx := &core.UiContext{Timezone: timezone.EuropeVienna}
	loc, err := time.LoadLocation("Europe/Vienna")
	if err != nil {
		t.Fatalf("loading location: %v", err)
	}

	minOffset, maxOffset := int64(-7), int64(1)
	f := NewTimeRangeField("range", "RANGE", false)
	f.Min = &minOffset
	f.Max = &maxOffset

	out := f.ExportForFrontend(ctx, nil)

	for key, offset := range map[string]int64{"min": minOffset, "max": maxOffset} {
		ts, ok := out[key].(int64)
		if !ok {
			t.Fatalf("%s: expected int64, got %T (%v)", key, out[key], out[key])
		}
		got := time.Unix(ts, 0).In(loc)
		if got.Hour() != 0 || got.Minute() != 0 || got.Second() != 0 {
			t.Errorf("%s: expected local midnight in Europe/Vienna, got %s", key, got.Format(time.RFC3339))
		}
		now := time.Now().In(loc)
		want := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).AddDate(0, 0, int(offset))
		if !got.Equal(want) {
			t.Errorf("%s: expected %s, got %s", key, want.Format(time.RFC3339), got.Format(time.RFC3339))
		}
	}
}

// TestTimeRangeField_Export_AbsoluteBoundsUnchanged verifies absolute timestamps
// (outside the ±10000 day-offset window) are passed through untouched.
func TestTimeRangeField_Export_AbsoluteBoundsUnchanged(t *testing.T) {
	absolute := int64(1750000000)
	f := NewTimeRangeField("range", "RANGE", false)
	f.Min = &absolute

	out := f.ExportForFrontend(&core.UiContext{Timezone: timezone.EuropeVienna}, nil)
	if out["min"] != absolute {
		t.Errorf("expected absolute min %d unchanged, got %v", absolute, out["min"])
	}
}

// TestTimeRangeField_Export_NilValueUsesDefault verifies the field default is used
// when no value is bound, instead of silently falling back to now (#15).
func TestTimeRangeField_Export_NilValueUsesDefault(t *testing.T) {
	f := NewTimeRangeFieldWithDefault("range", "RANGE", false, 7)
	def, ok := f.GetDefault().(*TimeRangeValue)
	if !ok {
		t.Fatalf("expected *TimeRangeValue default, got %T", f.GetDefault())
	}

	out := f.ExportForFrontend(&core.UiContext{Timezone: timezone.EuropeVienna}, nil)
	value, ok := out["value"].(map[string]int64)
	if !ok {
		t.Fatalf("expected map[string]int64 value, got %T", out["value"])
	}

	if value["start"] != def.Start.Unix() {
		t.Errorf("expected start from default (%d), got %d", def.Start.Unix(), value["start"])
	}
	if value["end"] != def.End.Unix() {
		t.Errorf("expected end from default (%d), got %d", def.End.Unix(), value["end"])
	}
}
