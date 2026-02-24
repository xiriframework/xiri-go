package field

import (
	"strings"
	"testing"
	"time"
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
