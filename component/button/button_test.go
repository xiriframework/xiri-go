package button_test

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/button"
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
)

func TestButton_WithData_AddsTopLevelDataKey(t *testing.T) {
	btn := button.NewSimpleApiButton("Go", url.NewUrl("/x"), core.ColorPrimary)
	btn.WithData(map[string]any{"_csv": true})

	out := btn.Print(nil)

	d, ok := out["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected out[\"data\"] to be map[string]any, got %T (%v)", out["data"], out["data"])
	}
	if d["_csv"] != true {
		t.Errorf("d[\"_csv\"] = %v, want true", d["_csv"])
	}
}

func TestButton_NoData_OmitsDataKey(t *testing.T) {
	btn := button.NewSimpleApiButton("Go", url.NewUrl("/x"), core.ColorPrimary)

	out := btn.Print(nil)

	if _, has := out["data"]; has {
		t.Errorf("expected no \"data\" key, got %v", out["data"])
	}
}

func TestButton_WithDataNil_OmitsDataKey(t *testing.T) {
	btn := button.NewSimpleApiButton("Go", url.NewUrl("/x"), core.ColorPrimary)
	btn.WithData(map[string]any{"_csv": true})
	btn.WithData(nil)

	out := btn.Print(nil)

	if _, has := out["data"]; has {
		t.Errorf("expected no \"data\" key after WithData(nil), got %v", out["data"])
	}
}

func TestButton_WithDataEmpty_OmitsDataKey(t *testing.T) {
	btn := button.NewSimpleApiButton("Go", url.NewUrl("/x"), core.ColorPrimary)
	btn.WithData(map[string]any{})

	out := btn.Print(nil)

	if _, has := out["data"]; has {
		t.Errorf("expected no \"data\" key for empty payload, got %v", out["data"])
	}
}

func TestButton_WithData_OverridesWithOptionData(t *testing.T) {
	btn := button.NewSimpleApiButton("Go", url.NewUrl("/x"), core.ColorPrimary)
	btn.WithOption("data", map[string]any{"old": true})
	btn.WithData(map[string]any{"new": true})

	out := btn.Print(nil)

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

func TestButton_LegacyOptionsDataStillWorks(t *testing.T) {
	// Backwards-compat: WithOption("data", payload) without WithData must still
	// surface payload at the top-level "data" key (legacy CSV-button workaround).
	btn := button.NewSimpleApiButton("Go", url.NewUrl("/x"), core.ColorPrimary)
	btn.WithOption("data", map[string]any{"_csv": true})

	out := btn.Print(nil)

	d, ok := out["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected out[\"data\"] to be map[string]any, got %T (%v)", out["data"], out["data"])
	}
	if d["_csv"] != true {
		t.Errorf("legacy options[\"data\"] payload broken: _csv=%v", d["_csv"])
	}
}

func TestTableButton_WithData_PassesThroughToButton(t *testing.T) {
	tb := button.NewTableButton(
		core.ButtonActionDownload,
		"csv",
		url.NewUrl("/x"),
		"CSV",
		core.ColorAccent,
		false,
		nil,
	).WithData(map[string]any{"_csv": true})

	out := tb.Print(nil)

	d, ok := out["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected out[\"data\"] to be map[string]any, got %T (%v)", out["data"], out["data"])
	}
	if d["_csv"] != true {
		t.Errorf("d[\"_csv\"] = %v, want true", d["_csv"])
	}
}
