package scatterchart

import (
	"reflect"
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
)

func ctx() *core.UiContext {
	return &core.UiContext{Translate: func(k string) string { return "T:" + k }}
}

func TestScatterChart(t *testing.T) {
	out := New("shops").Title("title").Axes("area", "revenue").
		Series("Filialen", core.ColorGreen, P(120, 4.2, "Wien"), P(80, 2.5, "")).
		Print(ctx())
	if out["type"] != "scatterchart" || out["id"] != "shops" {
		t.Fatalf("envelope=%v", out)
	}
	want := map[string]any{
		"title":  "T:title",
		"xLabel": "T:area",
		"yLabel": "T:revenue",
		"series": []map[string]any{{
			"name":  "Filialen",
			"color": "green",
			"points": []map[string]any{
				{"x": 120.0, "y": 4.2, "name": "Wien"},
				{"x": 80.0, "y": 2.5},
			},
		}},
	}
	if data := out["data"]; !reflect.DeepEqual(data, want) {
		t.Errorf("data=%#v", data)
	}
}
