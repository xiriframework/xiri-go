package expansion

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/button"
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/stat"
	"github.com/xiriframework/xiri-go/component/url"
	"github.com/xiriframework/xiri-go/response"
)

func expCtx() *core.UiContext {
	return &core.UiContext{Translate: func(k string) string { return k }}
}

// TestPanel_Buttons verifies that Panel.Buttons(...) emits the ButtonLine data
// ({class, buttons}) under "buttons", like Section and PageHeader do.
func TestPanel_Buttons(t *testing.T) {
	bl := button.NewButtonLine("right", nil).
		Add(button.NewSimpleCloseButton("A")).
		Add(button.NewSimpleCloseButton("B"))
	out := NewPanel("GPS").Buttons(bl).Print(expCtx())

	btns, ok := out["buttons"].(map[string]any)
	if !ok {
		t.Fatalf("expected buttons map, got %T", out["buttons"])
	}
	if btns["class"] != "right" {
		t.Errorf("class=%v want right", btns["class"])
	}
	list, ok := btns["buttons"].([]map[string]any)
	if !ok || len(list) != 2 {
		t.Fatalf("expected 2 buttons, got %T len %d", btns["buttons"], len(list))
	}
}

// TestPanel_NoButtons verifies the key is absent when no buttons were set.
func TestPanel_NoButtons(t *testing.T) {
	out := NewPanel("GPS").Print(expCtx())
	if _, has := out["buttons"]; has {
		t.Errorf("expected no 'buttons' key")
	}
}

// TestPanel_SetURL verifies the URL mode: Print emits the header fields plus "url" and an
// empty "data" array, and skips the static content — the frontend loads the panel itself.
func TestPanel_SetURL(t *testing.T) {
	p := NewPanel("Versicherung").
		Buttons(button.NewButtonLine("small", nil).Add(button.NewSimpleCloseButton("Edit"))).
		AddContent(stat.New("1", "x")).
		SetURL(url.NewUrlPrefix("/Vehicle/7/Panel/Insurance", "/api"))
	out := p.Print(expCtx())

	if out["url"] != "/api/Vehicle/7/Panel/Insurance" {
		t.Errorf("url=%v", out["url"])
	}
	if out["title"] != "Versicherung" {
		t.Errorf("title=%v", out["title"])
	}
	if _, has := out["buttons"]; !has {
		t.Errorf("expected buttons in URL mode")
	}
	if data, ok := out["data"].([]map[string]any); !ok || len(data) != 0 {
		t.Errorf("expected empty data in URL mode, got %v", out["data"])
	}
}

// TestPanel_PrintData_IgnoresURL: PrintData is the complete panel for the URL endpoint —
// with content, without "url".
func TestPanel_PrintData_IgnoresURL(t *testing.T) {
	p := NewPanel("Versicherung").AddContent(stat.New("1", "x")).SetURL(url.NewUrl("/x"))
	out := p.PrintData(expCtx())

	if _, has := out["url"]; has {
		t.Errorf("expected no url in PrintData, got %v", out["url"])
	}
	if data, ok := out["data"].([]map[string]any); !ok || len(data) != 1 {
		t.Errorf("expected one content component, got %v", out["data"])
	}
}

// TestPanel_DataResponseEnvelope: the URL endpoint answers with {"panel": <PrintData>} so the
// frontend can tell a complete panel from other responses.
func TestPanel_DataResponseEnvelope(t *testing.T) {
	p := NewPanel("Versicherung").AddContent(stat.New("1", "x"))
	res := p.DataResponse(expCtx())
	if res.Type != response.ResponseJSON {
		t.Fatalf("type=%v want JSON", res.Type)
	}
	body := res.Body.(map[string]any)
	if _, has := body["data"]; has {
		t.Errorf("expected no top-level 'data' key, got %v", body)
	}
	panel, ok := body["panel"].(map[string]any)
	if !ok {
		t.Fatalf("expected 'panel' map, got %T", body["panel"])
	}
	if panel["title"] != "Versicherung" {
		t.Errorf("panel=%v", panel)
	}
	if data, ok := panel["data"].([]map[string]any); !ok || len(data) != 1 {
		t.Errorf("expected content in panel data, got %v", panel["data"])
	}
}
