package barchart

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
)

func ctxDe() *core.UiContext {
	return &core.UiContext{
		Translate: func(k string) string {
			if k == "vehicle.strain" {
				return "Vehicle strain"
			}
			return k
		},
	}
}

func TestBarChart_Simple(t *testing.T) {
	bc := New("weekly").
		Mode(ModeSimple).
		Title("Weekly activities").
		YAxis(0, 12).
		Color(core.ColorPurple).
		Bar("M", 3).
		Bar("T", 8).
		Bar("W", 6)

	out := bc.Print(ctxDe())

	if out["type"] != "barchart" {
		t.Errorf("type=%v want barchart", out["type"])
	}
	if out["mode"] != "simple" {
		t.Errorf("mode=%v want simple", out["mode"])
	}
	data := out["data"].(map[string]any)
	if data["title"] != "Weekly activities" {
		t.Errorf("title=%v want translated", data["title"])
	}
	if data["yMin"].(float64) != 0 || data["yMax"].(float64) != 12 {
		t.Errorf("yMin/yMax=%v/%v", data["yMin"], data["yMax"])
	}
	if data["color"] != "purple" {
		t.Errorf("color=%v want purple", data["color"])
	}
	bars := data["bars"].([]map[string]any)
	if len(bars) != 3 {
		t.Fatalf("bars=%d want 3", len(bars))
	}
	if bars[1]["label"] != "T" || bars[1]["value"].(float64) != 8 {
		t.Errorf("bars[1]=%v", bars[1])
	}
}

func TestBarChart_Stacked(t *testing.T) {
	bc := New("strain").
		Mode(ModeStacked).
		Title("vehicle.strain").
		YAxis(0, 4).
		StackedBar("M", Seg(2, core.ColorGreen), Seg(1, core.ColorYellow), Seg(1, core.ColorRed)).
		StackedBar("T", Seg(3, core.ColorGreen), Seg(1, core.ColorYellow))

	out := bc.Print(ctxDe())

	if out["mode"] != "stacked" {
		t.Errorf("mode=%v want stacked", out["mode"])
	}
	data := out["data"].(map[string]any)
	if data["title"] != "Vehicle strain" {
		t.Errorf("title=%v want translated 'Vehicle strain'", data["title"])
	}
	bars := data["bars"].([]map[string]any)
	if len(bars) != 2 {
		t.Fatalf("bars=%d", len(bars))
	}
	segs := bars[0]["segments"].([]map[string]any)
	if len(segs) != 3 {
		t.Fatalf("segments=%d want 3", len(segs))
	}
	if segs[2]["color"] != "red" || segs[2]["value"].(float64) != 1 {
		t.Errorf("segs[2]=%v", segs[2])
	}
}

func TestBarChart_Heatmap(t *testing.T) {
	bc := New("engine").
		Mode(ModeHeatmap).
		Title("Engine system").
		Point(1700000000000, 1).
		Point(1700000060000, 0).
		Point(1700000120000, 3)

	out := bc.Print(ctxDe())

	if out["mode"] != "heatmap" {
		t.Errorf("mode=%v want heatmap", out["mode"])
	}
	data := out["data"].(map[string]any)
	pts := data["points"].([]map[string]any)
	if len(pts) != 3 {
		t.Fatalf("points=%d want 3", len(pts))
	}
	if pts[0]["time"].(int64) != 1700000000000 {
		t.Errorf("pts[0].time=%v", pts[0]["time"])
	}
	if pts[2]["value"].(float64) != 3 {
		t.Errorf("pts[2].value=%v", pts[2]["value"])
	}
}

func TestBarChart_AjaxMode(t *testing.T) {
	u := url.NewUrl("/api/strain")
	bc := New("strain").Mode(ModeStacked).SetURL(u).WithReload(true)

	out := bc.Print(ctxDe())

	data := out["data"].(map[string]any)
	if data["url"] == nil {
		t.Errorf("expected url in AJAX mode")
	}
	if data["reload"] != true {
		t.Errorf("expected reload=true, got %v", data["reload"])
	}
	if _, ok := data["bars"]; ok {
		t.Errorf("did not expect bars in AJAX mode")
	}
}

func TestBarChart_DefaultModeSimple(t *testing.T) {
	bc := New("x").Bar("a", 1)
	out := bc.Print(nil)
	if out["mode"] != "simple" {
		t.Errorf("default mode=%v want simple", out["mode"])
	}
}

func TestBarChart_NamedSimple(t *testing.T) {
	bc := New("weekly").Bar("M", 3).BarNamed("T", "Tuesday", 9)
	data := bc.Print(nil)["data"].(map[string]any)
	bars := data["bars"].([]map[string]any)
	if _, ok := bars[0]["name"]; ok {
		t.Errorf("Bar()-only entry should not have a name key")
	}
	if bars[1]["name"] != "Tuesday" {
		t.Errorf("BarNamed should set name, got %v", bars[1]["name"])
	}
}

func TestBarChart_NamedStacked(t *testing.T) {
	bc := New("strain").Mode(ModeStacked).
		StackedBarNamed("M", "Monday",
			SegNamed(2, "Low strain", core.ColorGreen),
			Seg(1, core.ColorYellow))
	data := bc.Print(nil)["data"].(map[string]any)
	bars := data["bars"].([]map[string]any)
	if bars[0]["name"] != "Monday" {
		t.Errorf("bar.name=%v want Monday", bars[0]["name"])
	}
	segs := bars[0]["segments"].([]map[string]any)
	if segs[0]["name"] != "Low strain" {
		t.Errorf("seg[0].name=%v want 'Low strain'", segs[0]["name"])
	}
	if _, ok := segs[1]["name"]; ok {
		t.Errorf("Seg()-only segment should not have a name key")
	}
}

func TestBarChart_NamedHeatmap(t *testing.T) {
	bc := New("engine").Mode(ModeHeatmap).
		Point(1, 0.1).
		PointNamed(2, "Repeat #1", 1.0)
	data := bc.Print(nil)["data"].(map[string]any)
	pts := data["points"].([]map[string]any)
	if _, ok := pts[0]["name"]; ok {
		t.Errorf("Point()-only entry should not have a name key")
	}
	if pts[1]["name"] != "Repeat #1" {
		t.Errorf("PointNamed should set name, got %v", pts[1]["name"])
	}
}

func TestBarChart_LinkSimple(t *testing.T) {
	bc := New("weekly").
		Bar("M", 3).Link(url.NewUrl("/day/mon")).
		Bar("T", 8)
	data := bc.Print(nil)["data"].(map[string]any)
	bars := data["bars"].([]map[string]any)
	if bars[0]["url"] != "/day/mon" {
		t.Errorf("bars[0].url=%v want /day/mon", bars[0]["url"])
	}
	if _, ok := bars[1]["url"]; ok {
		t.Errorf("Bar()-only entry should not have a url key")
	}
}

func TestBarChart_LinkUsesPrintWithoutPrefix(t *testing.T) {
	bc := New("weekly").
		Bar("M", 3).Link(url.NewUrlPrefix("/day/mon", "/api"))
	data := bc.Print(nil)["data"].(map[string]any)
	bars := data["bars"].([]map[string]any)
	if bars[0]["url"] != "/day/mon" {
		t.Errorf("bar url should use Print() (no prefix), got %v", bars[0]["url"])
	}
}

func TestBarChart_LinkStacked(t *testing.T) {
	bc := New("strain").Mode(ModeStacked).
		StackedBar("M", Seg(2, core.ColorGreen)).Link(url.NewUrl("/day/mon")).
		StackedBar("T", Seg(3, core.ColorGreen))
	data := bc.Print(nil)["data"].(map[string]any)
	bars := data["bars"].([]map[string]any)
	if bars[0]["url"] != "/day/mon" {
		t.Errorf("bars[0].url=%v want /day/mon", bars[0]["url"])
	}
	if _, ok := bars[1]["url"]; ok {
		t.Errorf("StackedBar()-only entry should not have a url key")
	}
}

func TestBarChart_LinkHeatmap(t *testing.T) {
	bc := New("engine").Mode(ModeHeatmap).
		Point(1, 0.1).Link(url.NewUrl("/event/1")).
		Point(2, 1.0)
	data := bc.Print(nil)["data"].(map[string]any)
	pts := data["points"].([]map[string]any)
	if pts[0]["url"] != "/event/1" {
		t.Errorf("pts[0].url=%v want /event/1", pts[0]["url"])
	}
	if _, ok := pts[1]["url"]; ok {
		t.Errorf("Point()-only entry should not have a url key")
	}
}
