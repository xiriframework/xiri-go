package gantt

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
)

func ctxDe() *core.UiContext { return &core.UiContext{Translate: func(k string) string { return k }} }

func TestGantt_Basic(t *testing.T) {
	g := New("project").Title("Project plan").
		Rows("Design", "Build", "Test").
		Task(0, "Wireframes",     1700000000000, 1700604800000, core.ColorBlue).
		Task(1, "Implementation", 1700604800000, 1701820800000, core.ColorGreen).
		Task(2, "QA",             1701820800000, 1702425600000, core.ColorWarning).
		XRange(1700000000000, 1702425600000)

	out := g.Print(ctxDe())
	if out["type"] != "gantt" {
		t.Errorf("type=%v", out["type"])
	}
	data := out["data"].(map[string]any)
	rows := data["rows"].([]string)
	if len(rows) != 3 || rows[1] != "Build" {
		t.Errorf("rows=%v", rows)
	}
	tasks := data["tasks"].([]map[string]any)
	if len(tasks) != 3 || tasks[0]["name"] != "Wireframes" || tasks[2]["color"] != "warn" {
		t.Errorf("tasks=%v", tasks)
	}
	if data["rangeStart"].(int64) != 1700000000000 {
		t.Errorf("rangeStart=%v", data["rangeStart"])
	}
}
