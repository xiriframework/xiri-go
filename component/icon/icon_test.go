package icon_test

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/icon"
)

func TestIcon_WithData_AddsTopLevelDataKey(t *testing.T) {
	ic := icon.NewIcon("home", "", core.ColorPrimary, nil)
	ic.WithData(map[string]any{"badge": 3})

	out := ic.Print(nil)

	d, ok := out["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected out[\"data\"] to be map[string]any, got %T (%v)", out["data"], out["data"])
	}
	if d["badge"] != 3 {
		t.Errorf("d[\"badge\"] = %v, want 3", d["badge"])
	}
}

func TestIcon_NoData_OmitsDataKey(t *testing.T) {
	ic := icon.NewIcon("home", "", core.ColorPrimary, nil)

	out := ic.Print(nil)

	if _, has := out["data"]; has {
		t.Errorf("expected no \"data\" key, got %v", out["data"])
	}
}

func TestIcon_WithData_OverridesWithOptionData(t *testing.T) {
	ic := icon.NewIcon("home", "", core.ColorPrimary, nil)
	ic.WithOption("data", map[string]any{"old": true})
	ic.WithData(map[string]any{"new": true})

	out := ic.Print(nil)

	d, ok := out["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected out[\"data\"] to be map[string]any, got %T", out["data"])
	}
	if d["new"] != true {
		t.Errorf("expected new=true, got %v", d["new"])
	}
	if _, has := d["old"]; has {
		t.Errorf("expected WithData to win, but got old=%v", d["old"])
	}
}
