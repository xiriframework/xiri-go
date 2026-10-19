package form_test

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/form"
	"github.com/xiriframework/xiri-go/component/url"
)

// The default Back and Save buttons stay reachable by keyboard.
func TestNewForm_DefaultButtons_NoTabIndex(t *testing.T) {
	out := form.NewForm(nil, url.NewUrl("/x"), nil, nil, nil, nil).Print(nil)

	buttons := out["data"].(map[string]any)["buttons"].([]map[string]any)
	if len(buttons) != 2 {
		t.Fatalf("expected 2 default buttons, got %d", len(buttons))
	}
	for _, b := range buttons {
		if v, has := b["tabIndex"]; has {
			t.Errorf("%s button: expected no \"tabIndex\" key, got %v", b["action"], v)
		}
	}
}
