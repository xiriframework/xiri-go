package formatter

import (
	"math"
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/types/locale"
)

// #18: FormatInteger must keep full int64 precision (no float64 detour).
func TestFormatIntegerPrecision(t *testing.T) {
	ctx := &core.UiContext{Locale: locale.EnUS}

	got := FormatInteger(9007199254740993, ctx) // 2^53 + 1, not exactly representable as float64
	want := "9,007,199,254,740,993"
	if got != want {
		t.Errorf("FormatInteger(2^53+1) = %q, want %q", got, want)
	}
}

// #18: MinInt64 must not overflow (guards against a negation-based implementation).
func TestFormatIntegerMinInt64(t *testing.T) {
	ctx := &core.UiContext{Locale: locale.EnUS}

	got := FormatInteger(math.MinInt64, ctx)
	want := "-9,223,372,036,854,775,808"
	if got != want {
		t.Errorf("FormatInteger(MinInt64) = %q, want %q", got, want)
	}
}
