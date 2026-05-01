package gaugechart

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
)

func ctxDe() *core.UiContext { return &core.UiContext{Translate: func(k string) string { return k }} }

func TestGauge_Basic(t *testing.T) {
	g := New("cpu").Title("CPU").Value(72).Range(0, 100).Color(core.ColorWarning).Label("%")

	out := g.Print(ctxDe())
	if out["type"] != "gaugechart" {
		t.Errorf("type=%v", out["type"])
	}
	data := out["data"].(map[string]any)
	if data["title"] != "CPU" {
		t.Errorf("title=%v", data["title"])
	}
	if data["value"].(float64) != 72 {
		t.Errorf("value=%v", data["value"])
	}
	if data["min"].(float64) != 0 || data["max"].(float64) != 100 {
		t.Errorf("min/max=%v/%v", data["min"], data["max"])
	}
	if data["color"] != "warn" {
		t.Errorf("color=%v want warn", data["color"])
	}
	if data["label"] != "%" {
		t.Errorf("label=%v", data["label"])
	}
}

func TestGauge_Compact(t *testing.T) {
	g := New("x").Value(50).Compact()
	data := g.Print(ctxDe())["data"].(map[string]any)
	if data["compact"] != true {
		t.Errorf("compact=%v", data["compact"])
	}
}
