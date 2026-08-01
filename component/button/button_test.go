package button_test

import (
	"bytes"
	"log/slog"
	"strings"
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

func TestButton_WithAutoLoad_AddsAutoLoadKey(t *testing.T) {
	btn := button.NewSimpleApiButton("Go", url.NewUrl("/x"), core.ColorPrimary)
	btn.WithAutoLoad(true)

	out := btn.Print(nil)

	if out["autoLoad"] != true {
		t.Errorf("out[\"autoLoad\"] = %v, want true", out["autoLoad"])
	}
}

func TestButton_NoAutoLoad_OmitsAutoLoadKey(t *testing.T) {
	btn := button.NewSimpleApiButton("Go", url.NewUrl("/x"), core.ColorPrimary)

	out := btn.Print(nil)

	if _, ok := out["autoLoad"]; ok {
		t.Errorf("expected no \"autoLoad\" key by default, got %v", out["autoLoad"])
	}
}

func TestTableButton_WithAutoLoad_PassesThroughToButton(t *testing.T) {
	tb := button.NewTableButton(
		core.ButtonActionApi,
		"refresh",
		url.NewUrl("/x"),
		"Reload",
		core.ColorPrimary,
		false,
		nil,
	).WithAutoLoad(true)

	out := tb.Print(nil)

	if out["autoLoad"] != true {
		t.Errorf("out[\"autoLoad\"] = %v, want true", out["autoLoad"])
	}
}

// "_blank" on a download button is what makes the frontend display the file in a
// new tab instead of saving it — it has to survive the TableButton wrapper too.
func TestTableButton_WithTarget_PassesThroughToButton(t *testing.T) {
	tb := button.NewTableButton(
		core.ButtonActionDownload,
		"pdf",
		url.NewUrl("/x"),
		"PDF",
		core.ColorAccent,
		false,
		nil,
	).WithTarget("_blank")

	out := tb.Print(nil)

	if out["target"] != "_blank" {
		t.Errorf("out[\"target\"] = %v, want \"_blank\"", out["target"])
	}
}

func TestDownloadButton_DefaultTarget_IsSelf(t *testing.T) {
	filename := "report.csv"
	btn := button.NewDownloadButton("CSV", url.NewUrl("/x"), core.ColorPrimary, core.ButtonTypeRaised, "", &filename, false, nil, nil)

	out := btn.Print(nil)

	if out["target"] != "_self" {
		t.Errorf("out[\"target\"] = %v, want \"_self\"", out["target"])
	}
}

// Icon-only buttons carry no visible text, so the frontend can only build an
// accessible name from "hint" — Print() does not even emit "text" for those
// types. A missing hint therefore means a button screen readers cannot name.
func TestButton_IconTypeWithoutHint_Warns(t *testing.T) {
	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	defer slog.SetDefault(prev)

	btn := button.NewButton(core.ButtonActionApi, "", url.NewUrl("/x"), core.ColorPrimary,
		core.ButtonTypeIcon, "", "edit", false, nil, false, "_self", nil)
	btn.Print(nil)

	if !strings.Contains(buf.String(), "hint") {
		t.Errorf("expected a warning about the missing hint, got %q", buf.String())
	}
}

// The hint may arrive after the constructor — warning before the button is
// fully configured would flag correct fluent usage.
func TestButton_IconTypeHintViaBuilder_DoesNotWarn(t *testing.T) {
	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	defer slog.SetDefault(prev)

	btn := button.NewButton(core.ButtonActionApi, "", url.NewUrl("/x"), core.ColorPrimary,
		core.ButtonTypeIcon, "", "edit", false, nil, false, "_self", nil).
		WithHint("Bearbeiten")
	btn.Print(nil)

	if buf.Len() != 0 {
		t.Errorf("expected no warning, got %q", buf.String())
	}
}

// Print() may run more than once per button; the warning must not pile up.
func TestButton_IconTypeWithoutHint_WarnsOnlyOnce(t *testing.T) {
	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	defer slog.SetDefault(prev)

	btn := button.NewButton(core.ButtonActionApi, "", url.NewUrl("/x"), core.ColorPrimary,
		core.ButtonTypeIcon, "", "edit", false, nil, false, "_self", nil)
	btn.Print(nil)
	btn.Print(nil)
	btn.Print(nil)

	if n := strings.Count(buf.String(), "hint"); n != 1 {
		t.Errorf("expected exactly one warning, got %d in %q", n, buf.String())
	}
}

func TestButton_IconTypeWithHint_DoesNotWarn(t *testing.T) {
	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	defer slog.SetDefault(prev)

	btn := button.NewButton(core.ButtonActionApi, "", url.NewUrl("/x"), core.ColorPrimary,
		core.ButtonTypeIcon, "Bearbeiten", "edit", false, nil, false, "_self", nil)
	btn.Print(nil)

	if buf.Len() != 0 {
		t.Errorf("expected no warning, got %q", buf.String())
	}
}

// A text-bearing button needs no hint — its label is visible.
func TestButton_TextTypeWithoutHint_DoesNotWarn(t *testing.T) {
	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	defer slog.SetDefault(prev)

	button.NewSimpleApiButton("Speichern", url.NewUrl("/x"), core.ColorPrimary).Print(nil)

	if buf.Len() != 0 {
		t.Errorf("expected no warning, got %q", buf.String())
	}
}
