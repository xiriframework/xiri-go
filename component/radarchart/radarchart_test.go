package radarchart

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
)

func ctx() *core.UiContext {
	return &core.UiContext{Translate: func(k string) string { return map[string]string{"t.title": "Vergleich"}[k] }}
}

func TestRadarChart(t *testing.T) {
	out := New("cmp").Title("t.title").
		Indicator("Tempo", 10).Indicator("Preis", 0).
		Series("A", core.ColorBlue, 8, 3).Series("B", "", 5, 6).
		Area().Compact().WithDisplay("xcol-md-6").
		Print(ctx())
	if out["type"] != "radarchart" || out["id"] != "cmp" || out["display"] != "xcol-md-6" {
		t.Fatalf("envelope=%v", out)
	}
	data := out["data"].(map[string]any)
	want := map[string]any{
		"title":   "Vergleich",
		"compact": true,
		"area":    true,
		"indicators": []map[string]any{
			{"name": "Tempo", "max": 10.0},
			{"name": "Preis"},
		},
		"series": []map[string]any{
			{"name": "A", "values": []float64{8, 3}, "color": "blue"},
			{"name": "B", "values": []float64{5, 6}},
		},
	}
	if !reflect.DeepEqual(data, want) {
		t.Errorf("data=%#v", data)
	}
}

// Angular erwartet values als Array, auch ohne Werte.
func TestRadarChartEmptySeries(t *testing.T) {
	b, _ := json.Marshal(New("x").Series("A", "").Print(ctx())["data"])
	if !strings.Contains(string(b), `"values":[]`) {
		t.Errorf("json=%s", b)
	}
}
