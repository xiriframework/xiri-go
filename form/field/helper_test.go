package field

import (
	"math"
	"testing"
	"time"

	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/types/timezone"
)

// TestToInt32 covers the shared lossless int32 conversion (#5). Every field type
// that stores an int32 routes through it, so truncation or wrap-around here would
// silently select the wrong record.
func TestToInt32(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		want      int32
		expectErr bool
	}{
		{"int32", int32(-7), -7, false},
		{"int", 42, 42, false},
		{"int64", int64(42), 42, false},
		{"int64 max int32", int64(math.MaxInt32), math.MaxInt32, false},
		{"int64 above max int32", int64(math.MaxInt32) + 1, 0, true},
		{"int64 below min int32", int64(math.MinInt32) - 1, 0, true},
		{"int64 far out of range", int64(1) << 40, 0, true},
		{"integer float", 5.0, 5, false},
		{"fractional float", 1.9, 0, true},
		{"negative fractional float", -1.5, 0, true},
		{"float max int32", float64(math.MaxInt32), math.MaxInt32, false},
		{"float above max int32", float64(math.MaxInt32) + 1, 0, true},
		{"NaN", math.NaN(), 0, true},
		{"positive infinity", math.Inf(1), 0, true},
		{"negative infinity", math.Inf(-1), 0, true},
		{"string integer", "42", 42, false},
		{"string negative", "-42", -42, false},
		{"string fractional", "1.9", 0, true},
		{"string out of int32 range", "3000000000", 0, true},
		{"string empty", "", 0, true},
		{"string not a number", "abc", 0, true},
		{"unsupported type", true, 0, true},
		{"nil", nil, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toInt32(tt.input)
			if tt.expectErr {
				if err == nil {
					t.Fatalf("input=%v (%T): expected error, got %d", tt.input, tt.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("input=%v (%T): unexpected error: %v", tt.input, tt.input, err)
			}
			if got != tt.want {
				t.Errorf("input=%v (%T): expected %d, got %d", tt.input, tt.input, tt.want, got)
			}
		})
	}
}

// TestResolveDateBound_DSTSafe pins the day-offset semantics across a DST transition
// (#15). Europe/Vienna springs forward on 2026-03-29 (02:00 CET → 03:00 CEST), so a
// day is not 86400 seconds there — "now + n*86400" lands an hour off.
func TestResolveDateBound_DSTSafe(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Vienna")
	if err != nil {
		t.Fatalf("loading location: %v", err)
	}
	ctx := &core.UiContext{Timezone: timezone.EuropeVienna}
	now := time.Date(2026, 3, 28, 12, 0, 0, 0, loc)

	tests := []struct {
		name  string
		bound int64
		want  time.Time
	}{
		{"today", 0, time.Date(2026, 3, 28, 0, 0, 0, 0, loc)},
		{"day of transition", 1, time.Date(2026, 3, 29, 0, 0, 0, 0, loc)},
		{"across spring forward", 2, time.Date(2026, 3, 30, 0, 0, 0, 0, loc)},
		{"negative offset", -3, time.Date(2026, 3, 25, 0, 0, 0, 0, loc)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveDateBound(ctx, tt.bound, now)
			if got != tt.want.Unix() {
				t.Errorf("bound=%d: expected %s (%d), got %s (%d)",
					tt.bound, tt.want.Format(time.RFC3339), tt.want.Unix(),
					time.Unix(got, 0).In(loc).Format(time.RFC3339), got)
			}
		})
	}

	// Guard: the "across spring forward" case must actually discriminate against the
	// old formula, otherwise this test proves nothing about DST.
	midnight := time.Date(2026, 3, 28, 0, 0, 0, 0, loc)
	if naive := midnight.Unix() + 2*86400; naive == time.Date(2026, 3, 30, 0, 0, 0, 0, loc).Unix() {
		t.Error("test does not discriminate: naive 86400-arithmetic yields the same result")
	}
}

// TestResolveDateBound_AbsoluteAndNilContext covers the pass-through of absolute
// timestamps and the nil-context fallback (core.Component.Print allows ctx == nil).
func TestResolveDateBound_AbsoluteAndNilContext(t *testing.T) {
	now := time.Unix(1750000000, 0)

	for _, bound := range []int64{dayOffsetLimit, -dayOffsetLimit, 1750000000, -1750000000} {
		if got := resolveDateBound(nil, bound, now); got != bound {
			t.Errorf("bound=%d: expected pass-through, got %d", bound, got)
		}
	}

	// Day offset with nil ctx must not panic and must land on midnight of the
	// timezone SafeTimezone() falls back to.
	loc, err := time.LoadLocation((*core.UiContext)(nil).SafeTimezone().GetIANA())
	if err != nil {
		t.Fatalf("loading fallback location: %v", err)
	}
	got := time.Unix(resolveDateBound(nil, 0, now), 0).In(loc)
	if got.Hour() != 0 || got.Minute() != 0 || got.Second() != 0 {
		t.Errorf("expected midnight in fallback timezone, got %s", got.Format(time.RFC3339))
	}
}
