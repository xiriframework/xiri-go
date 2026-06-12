package heatmap

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
)

func ctxDe() *core.UiContext { return &core.UiContext{Translate: func(k string) string { return k }} }

func TestHeatmap_Basic(t *testing.T) {
	h := New("activity").
		Title("Hourly activity").
		XLabels("00", "06", "12", "18").
		YLabels("Mon", "Tue").
		Cell(0, 0, 1).Cell(2, 0, 5).Cell(2, 1, 3).
		Range(0, 10).
		ColorRange("#eee", "#7c3aed").
		ShowValues()

	out := h.Print(ctxDe())
	if out["type"] != "heatmap" {
		t.Errorf("type=%v", out["type"])
	}
	data := out["data"].(map[string]any)
	cells := data["cells"].([]map[string]any)
	if len(cells) != 3 || cells[0]["x"].(int) != 0 || cells[1]["value"].(float64) != 5 {
		t.Errorf("cells=%v", cells)
	}
	if data["min"].(float64) != 0 || data["max"].(float64) != 10 {
		t.Errorf("min/max=%v/%v", data["min"], data["max"])
	}
	if data["showValues"] != true {
		t.Errorf("showValues=%v", data["showValues"])
	}
}

func TestHeatmap_XLabelRotate(t *testing.T) {
	h := New("activity").
		XLabels("Morning", "Afternoon", "Evening").
		YLabels("Mon").
		Cell(0, 0, 1).
		XLabelRotate(90)

	data := h.Print(ctxDe())["data"].(map[string]any)
	if data["xLabelRotate"] != 90 {
		t.Errorf("xLabelRotate=%v want 90", data["xLabelRotate"])
	}

	noRotate := New("x").XLabels("a").YLabels("b").Cell(0, 0, 1)
	data2 := noRotate.Print(ctxDe())["data"].(map[string]any)
	if _, ok := data2["xLabelRotate"]; ok {
		t.Errorf("xLabelRotate key should be absent by default")
	}
}
