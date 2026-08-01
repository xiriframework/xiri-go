package table

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
)

// The frontend understands density: 'compact' | 'regular' | 'relaxed' and treats
// dense: true only as a legacy alias for 'compact'. Without SetDensity a Go
// consumer can never reach 'relaxed'.
func TestSetDensity_EmitsDensityOption(t *testing.T) {
	for _, d := range []Density{DensityCompact, DensityRegular, DensityRelaxed} {
		builder := NewBuilder[testOptionRow]()
		builder.SetDensity(d)

		opts := builder.Build().exportOptions(&core.UiContext{})

		if opts["density"] != string(d) {
			t.Errorf("opts[\"density\"] = %v, want %q", opts["density"], string(d))
		}
	}
}

func TestSetDensity_OmittedWhenUnset(t *testing.T) {
	opts := NewBuilder[testOptionRow]().Build().exportOptions(&core.UiContext{})

	if v, ok := opts["density"]; ok {
		t.Errorf("expected no \"density\" key by default, got %v", v)
	}
}

// SetDense stays functional for existing callers and keeps emitting "dense".
func TestSetDense_StillEmitsDense(t *testing.T) {
	builder := NewBuilder[testOptionRow]()
	builder.SetDense(true)

	opts := builder.Build().exportOptions(&core.UiContext{})

	if opts["dense"] != true {
		t.Errorf("opts[\"dense\"] = %v, want true", opts["dense"])
	}
	if _, ok := opts["density"]; ok {
		t.Error("SetDense must not synthesise a density value")
	}
}

// Both set: the frontend gives density precedence (mergeOptions only falls back
// to the dense alias when density is undefined), so both may travel together.
func TestSetDensityAndSetDense_BothEmitted(t *testing.T) {
	builder := NewBuilder[testOptionRow]()
	builder.SetDense(true).SetDensity(DensityRelaxed)

	opts := builder.Build().exportOptions(&core.UiContext{})

	if opts["dense"] != true {
		t.Errorf("opts[\"dense\"] = %v, want true", opts["dense"])
	}
	if opts["density"] != string(DensityRelaxed) {
		t.Errorf("opts[\"density\"] = %v, want %q", opts["density"], string(DensityRelaxed))
	}
}
