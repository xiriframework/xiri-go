package markdown

import (
	"reflect"
	"testing"
)

func TestMarkdown(t *testing.T) {
	out := New("## Titel\n\n- Punkt").WithDisplay("xcol-md-6").Print(nil)
	want := map[string]any{
		"type":    "markdown",
		"display": "xcol-md-6",
		"data":    map[string]any{"text": "## Titel\n\n- Punkt"},
	}
	if !reflect.DeepEqual(out, want) {
		t.Errorf("out=%#v", out)
	}
	if _, has := New("x").Print(nil)["display"]; has {
		t.Errorf("display without WithDisplay")
	}
}
