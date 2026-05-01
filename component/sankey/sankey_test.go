package sankey

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
)

func ctxDe() *core.UiContext { return &core.UiContext{Translate: func(k string) string { return k }} }

func TestSankey_Basic(t *testing.T) {
	s := New("flow").Title("Flow").
		Node("A", core.ColorBlue).
		Node("B", core.ColorGreen).
		Node("C", "").
		Link("A", "B", 5).
		Link("B", "C", 3)

	out := s.Print(ctxDe())
	if out["type"] != "sankey" {
		t.Errorf("type=%v", out["type"])
	}
	data := out["data"].(map[string]any)
	nodes := data["nodes"].([]map[string]any)
	if len(nodes) != 3 || nodes[0]["color"] != "blue" {
		t.Errorf("nodes=%v", nodes)
	}
	if _, has := nodes[2]["color"]; has {
		t.Errorf("default node should not have color: %v", nodes[2])
	}
	links := data["links"].([]map[string]any)
	if len(links) != 2 || links[0]["source"] != "A" || links[1]["value"].(float64) != 3 {
		t.Errorf("links=%v", links)
	}
}

func TestSankey_Vertical(t *testing.T) {
	s := New("x").Vertical().Node("A", "").Link("A", "B", 1)
	data := s.Print(ctxDe())["data"].(map[string]any)
	if data["orient"] != "vertical" {
		t.Errorf("orient=%v", data["orient"])
	}
}
