package callout_test

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/button"
	"github.com/xiriframework/xiri-go/component/callout"
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
)

func TestCallout_Print_TypeAndRequiredFields(t *testing.T) {
	c := callout.New("info", "Heads up.")

	out := c.Print(nil)

	if out["type"] != "callout" {
		t.Errorf("type = %v, want %q", out["type"], "callout")
	}

	data, ok := out["data"].(map[string]any)
	if !ok {
		t.Fatalf("data is not a map: %T", out["data"])
	}
	if data["tone"] != "info" {
		t.Errorf("tone = %v, want %q", data["tone"], "info")
	}
	if data["text"] != "Heads up." {
		t.Errorf("text = %v, want %q", data["text"], "Heads up.")
	}
}

func TestCallout_Print_OmitsUnsetOptionals(t *testing.T) {
	data := callout.New("success", "Done.").Print(nil)["data"].(map[string]any)

	for _, key := range []string{"title", "icon", "actions", "dismissible", "compact"} {
		if _, has := data[key]; has {
			t.Errorf("expected no %q key, got %v", key, data[key])
		}
	}
}

func TestCallout_Print_SetsOptionals(t *testing.T) {
	data := callout.New("warning", "Careful.").
		Title("Warning").
		Icon("warning").
		Dismissible().
		Compact().
		Print(nil)["data"].(map[string]any)

	if data["title"] != "Warning" {
		t.Errorf("title = %v, want %q", data["title"], "Warning")
	}
	if data["icon"] != "warning" {
		t.Errorf("icon = %v, want %q", data["icon"], "warning")
	}
	if data["dismissible"] != true {
		t.Errorf("dismissible = %v, want true", data["dismissible"])
	}
	if data["compact"] != true {
		t.Errorf("compact = %v, want true", data["compact"])
	}
}

func TestCallout_Print_SerializesActions(t *testing.T) {
	btn := button.NewSimpleLinkButton("Reload", url.NewUrl("/reload"), core.ColorPrimary)

	data := callout.New("info", "Update available.").
		AddAction(btn).
		Print(nil)["data"].(map[string]any)

	actions, ok := data["actions"].([]map[string]any)
	if !ok {
		t.Fatalf("actions is not a []map[string]any: %T", data["actions"])
	}
	if len(actions) != 1 {
		t.Fatalf("len(actions) = %d, want 1", len(actions))
	}
	if actions[0]["text"] != "Reload" {
		t.Errorf("actions[0].text = %v, want %q", actions[0]["text"], "Reload")
	}
	if actions[0]["action"] != string(core.ButtonActionLink) {
		t.Errorf("actions[0].action = %v, want %q", actions[0]["action"], core.ButtonActionLink)
	}
}
