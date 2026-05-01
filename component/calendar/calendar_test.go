package calendar

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
)

func ctxDe() *core.UiContext { return &core.UiContext{Translate: func(k string) string { return k }} }

func TestCalendar_Year(t *testing.T) {
	c := New("activity").
		Title("Daily activity").
		Year("2025").
		Cell("2025-01-15", 3).
		Cell("2025-02-01", 8).
		MinMax(0, 10)

	out := c.Print(ctxDe())
	if out["type"] != "calendar" {
		t.Errorf("type=%v", out["type"])
	}
	data := out["data"].(map[string]any)
	if data["range"] != "2025" {
		t.Errorf("range=%v", data["range"])
	}
	cells := data["cells"].([]map[string]any)
	if len(cells) != 2 || cells[0]["date"] != "2025-01-15" || cells[1]["value"].(float64) != 8 {
		t.Errorf("cells=%v", cells)
	}
}

func TestCalendar_Range(t *testing.T) {
	c := New("x").Range("2025-01-01", "2025-06-30").Cell("2025-03-01", 1)
	data := c.Print(ctxDe())["data"].(map[string]any)
	rng := data["range"].([]string)
	if rng[0] != "2025-01-01" || rng[1] != "2025-06-30" {
		t.Errorf("range=%v", rng)
	}
}
