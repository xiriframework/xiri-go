package bulletchart

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
)

func ctxDe() *core.UiContext { return &core.UiContext{Translate: func(k string) string { return k }} }

func TestBullet_Basic(t *testing.T) {
	b := New("sales").Title("Sales").Value(70).Target(90).Max(100).Color(core.ColorWarning).Label("k€")

	out := b.Print(ctxDe())
	if out["type"] != "bulletchart" {
		t.Errorf("type=%v", out["type"])
	}
	data := out["data"].(map[string]any)
	if data["title"] != "Sales" {
		t.Errorf("title=%v", data["title"])
	}
	if data["value"].(float64) != 70 {
		t.Errorf("value=%v", data["value"])
	}
	if data["target"].(float64) != 90 {
		t.Errorf("target=%v", data["target"])
	}
	if data["max"].(float64) != 100 {
		t.Errorf("max=%v", data["max"])
	}
	if data["color"] != "warn" {
		t.Errorf("color=%v want warn", data["color"])
	}
	if data["label"] != "k€" {
		t.Errorf("label=%v", data["label"])
	}
}

func TestBullet_Compact(t *testing.T) {
	b := New("x").Value(50).Compact()
	data := b.Print(ctxDe())["data"].(map[string]any)
	if data["compact"] != true {
		t.Errorf("compact=%v", data["compact"])
	}
}

func TestBullet_OptionalOmitted(t *testing.T) {
	// target/max/label are optional and must be absent when unset
	b := New("x").Value(42)
	data := b.Print(ctxDe())["data"].(map[string]any)
	if data["value"].(float64) != 42 {
		t.Errorf("value=%v", data["value"])
	}
	if _, ok := data["target"]; ok {
		t.Error("target should be omitted when unset")
	}
	if _, ok := data["max"]; ok {
		t.Error("max should be omitted when unset")
	}
	if _, ok := data["label"]; ok {
		t.Error("label should be omitted when empty")
	}
}
