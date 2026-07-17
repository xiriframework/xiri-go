package status_test

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/status"
)

func data(s *status.Status) map[string]any {
	out := s.Print(nil)
	if out["type"] != "status" {
		panic("type != status")
	}
	return out["data"].(map[string]any)
}

func TestStatus_MinimalOnlyLabel(t *testing.T) {
	d := data(status.New("Online"))

	if d["label"] != "Online" {
		t.Errorf("label = %v, want %q", d["label"], "Online")
	}
	// Defaults (tone/variant) werden vom Frontend gesetzt -> hier nicht enthalten.
	for _, k := range []string{"tone", "variant", "icon", "hint", "ariaLabel"} {
		if _, has := d[k]; has {
			t.Errorf("expected no %q key, got %v", k, d[k])
		}
	}
}

func TestStatus_AllFields(t *testing.T) {
	s := status.New("Delayed").
		Tone(status.ToneWarning).
		Variant(status.VariantDot).
		Icon("schedule").
		Hint("Since 5 minutes").
		AriaLabel("Delivery delayed")

	d := data(s)

	cases := map[string]any{
		"label":     "Delayed",
		"tone":      "warning",
		"variant":   "dot",
		"icon":      "schedule",
		"hint":      "Since 5 minutes",
		"ariaLabel": "Delivery delayed",
	}
	for k, want := range cases {
		if d[k] != want {
			t.Errorf("data[%q] = %v, want %v", k, d[k], want)
		}
	}
}

func TestStatus_ToneAndVariantConstants(t *testing.T) {
	tones := []status.Tone{status.ToneNeutral, status.ToneInfo, status.ToneSuccess, status.ToneWarning, status.ToneError}
	wantTones := []string{"neutral", "info", "success", "warning", "error"}
	for i, tone := range tones {
		if d := data(status.New("x").Tone(tone)); d["tone"] != wantTones[i] {
			t.Errorf("tone = %v, want %q", d["tone"], wantTones[i])
		}
	}

	variants := []status.Variant{status.VariantBadge, status.VariantDot, status.VariantText}
	wantVariants := []string{"badge", "dot", "text"}
	for i, v := range variants {
		if d := data(status.New("x").Variant(v)); d["variant"] != wantVariants[i] {
			t.Errorf("variant = %v, want %q", d["variant"], wantVariants[i])
		}
	}
}
