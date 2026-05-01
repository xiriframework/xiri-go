package tree

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
)

func ctxDe() *core.UiContext { return &core.UiContext{Translate: func(k string) string { return k }} }

func TestTree_Basic(t *testing.T) {
	root := NewNode("root").WithValue(100).AppendChild(
		NewNode("A").WithValue(60).AppendChild(
			NewNode("A1").WithValue(30),
			NewNode("A2").WithValue(30),
		),
		NewNode("B").WithValue(40).Collapse(),
	)

	tr := New("org").Title("Org chart").Root(root).Orient("TB")
	out := tr.Print(ctxDe())
	if out["type"] != "tree" {
		t.Errorf("type=%v", out["type"])
	}
	data := out["data"].(map[string]any)
	if data["orient"] != "TB" {
		t.Errorf("orient=%v", data["orient"])
	}
	r := data["root"].(map[string]any)
	if r["name"] != "root" || r["value"].(float64) != 100 {
		t.Errorf("root=%v", r)
	}
	children := r["children"].([]map[string]any)
	if len(children) != 2 || children[0]["name"] != "A" {
		t.Errorf("children=%v", children)
	}
	if children[1]["collapsed"] != true {
		t.Errorf("B not collapsed: %v", children[1])
	}
	a := children[0]
	aChildren := a["children"].([]map[string]any)
	if len(aChildren) != 2 {
		t.Errorf("A children=%d", len(aChildren))
	}
}
