package linechart

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
)

func ctxDe() *core.UiContext {
	return &core.UiContext{
		Translate: func(k string) string { return k },
	}
}

func TestLineChart_Basic(t *testing.T) {
	lc := New("revenue").
		Title("Monthly revenue").
		XLabels("Jan", "Feb", "Mar", "Apr").
		Line("Product A", []float64{100, 120, 150, 180}).Color(core.ColorBlue).Done().
		Line("Product B", []float64{80, 95, 110, 130}).Color(core.ColorGreen).Smooth().Done().
		YAxis(0, 200)

	out := lc.Print(ctxDe())
	if out["type"] != "linechart" {
		t.Errorf("type=%v want linechart", out["type"])
	}
	data := out["data"].(map[string]any)
	if data["title"] != "Monthly revenue" {
		t.Errorf("title=%v", data["title"])
	}
	if data["smooth"] != true {
		t.Errorf("smooth=%v want true", data["smooth"])
	}
	if data["yMin"].(float64) != 0 || data["yMax"].(float64) != 200 {
		t.Errorf("yMin/Max=%v/%v", data["yMin"], data["yMax"])
	}
	labels := data["xLabels"].([]string)
	if len(labels) != 4 || labels[0] != "Jan" {
		t.Errorf("xLabels=%v", labels)
	}
	lines := data["lines"].([]map[string]any)
	if len(lines) != 2 {
		t.Fatalf("lines=%d want 2", len(lines))
	}
	if lines[0]["name"] != "Product A" || lines[0]["color"] != "blue" {
		t.Errorf("lines[0]=%v", lines[0])
	}
	if lines[1]["color"] != "green" {
		t.Errorf("lines[1].color=%v", lines[1]["color"])
	}
}

func TestLineChart_Compact(t *testing.T) {
	lc := New("x").
		Line("A", []float64{1, 2, 3}).
		Done().
		Compact()

	out := lc.Print(ctxDe())
	data := out["data"].(map[string]any)
	if data["compact"] != true {
		t.Errorf("compact=%v want true", data["compact"])
	}
}

func TestLineChart_DashedArea(t *testing.T) {
	lc := New("x").
		Line("A", []float64{1, 2}).Dashed().Area().Done()

	out := lc.Print(ctxDe())
	lines := out["data"].(map[string]any)["lines"].([]map[string]any)
	if lines[0]["dashed"] != true || lines[0]["area"] != true {
		t.Errorf("dashed/area not set: %v", lines[0])
	}
}

func TestLineChart_XLabelRotate(t *testing.T) {
	lc := New("revenue").
		XLabels("January", "February", "March").
		Line("A", []float64{1, 2, 3}).XLabelRotate(45).Done()

	data := lc.Print(ctxDe())["data"].(map[string]any)
	if data["xLabelRotate"] != 45 {
		t.Errorf("xLabelRotate=%v want 45", data["xLabelRotate"])
	}

	noRotate := New("x").Line("A", []float64{1}).Done()
	data2 := noRotate.Print(ctxDe())["data"].(map[string]any)
	if _, ok := data2["xLabelRotate"]; ok {
		t.Errorf("xLabelRotate key should be absent by default")
	}
}
