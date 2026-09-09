package expansion

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/button"
	"github.com/xiriframework/xiri-go/component/core"
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
