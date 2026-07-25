package table

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/types/locale"
)

// #18 (extended): the table integer formatters must not lose precision via a float64 detour.
func TestIntegerFormattersPrecision(t *testing.T) {
	ctx := &core.UiContext{Locale: locale.EnUS}
	const big = 9007199254740993 // 2^53 + 1, not exactly representable as float64
	const want = "9,007,199,254,740,993"

	if got := createIntegerFormatter().Format(int64(big), nil, OutputWeb, ctx); got != want {
		t.Errorf("createIntegerFormatter = %v, want %q", got, want)
	}

	got2 := createText2IntFormatter().Format([2]int{big, 0}, nil, OutputWeb, ctx)
	if arr, ok := got2.([2]string); !ok || arr[0] != want {
		t.Errorf("createText2IntFormatter = %v, want first %q", got2, want)
	}

	gotN := createIntegerNFormatter().Format([]int{big}, nil, OutputWeb, ctx)
	if arr, ok := gotN.([]string); !ok || len(arr) != 1 || arr[0] != want {
		t.Errorf("createIntegerNFormatter = %v, want [%q]", gotN, want)
	}
}
