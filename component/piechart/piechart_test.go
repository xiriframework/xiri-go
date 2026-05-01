package piechart

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
)

func ctxDe() *core.UiContext { return &core.UiContext{Translate: func(k string) string { return k }} }

func TestPieChart_Basic(t *testing.T) {
	pc := New("traffic").
		Title("Traffic sources").
		Slice("Direct", 1234, core.ColorBlue).
		Slice("Search", 856, core.ColorGreen).
		Slice("Social", 423, core.ColorPurple)

	out := pc.Print(ctxDe())
	if out["type"] != "piechart" {
		t.Errorf("type=%v", out["type"])
	}
	data := out["data"].(map[string]any)
	if data["title"] != "Traffic sources" {
		t.Errorf("title=%v", data["title"])
	}
	if _, has := data["donut"]; has {
		t.Errorf("default chart should not have donut=true")
	}
	slices := data["slices"].([]map[string]any)
	if len(slices) != 3 {
		t.Fatalf("slices=%d", len(slices))
	}
	if slices[0]["name"] != "Direct" || slices[0]["color"] != "blue" {
		t.Errorf("slices[0]=%v", slices[0])
	}
	if slices[2]["value"].(float64) != 423 {
		t.Errorf("slices[2].value=%v", slices[2]["value"])
	}
}

func TestPieChart_DonutCompact(t *testing.T) {
	pc := New("x").Slice("A", 1, core.ColorBlue).Donut().Compact()
	data := pc.Print(ctxDe())["data"].(map[string]any)
	if data["donut"] != true {
		t.Errorf("donut=%v", data["donut"])
	}
	if data["compact"] != true {
		t.Errorf("compact=%v", data["compact"])
	}
}
