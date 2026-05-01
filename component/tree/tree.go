// Package tree provides a hierarchical tree-chart component (echarts tree series).
package tree

import (
	"github.com/xiriframework/xiri-go/component/chart"
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
	"github.com/xiriframework/xiri-go/response"
)

// Node is a tree node with optional value and recursive children.
// Children may be added via the AppendChild builder helper.
type Node struct {
	Name      string
	Value     *float64
	Collapsed bool
	Children  []*Node
}

// NewNode creates a new tree node.
func NewNode(name string) *Node { return &Node{Name: name} }

// WithValue sets an optional numeric value displayed in the tooltip.
func (n *Node) WithValue(v float64) *Node { n.Value = &v; return n }

// Collapse marks the node initially collapsed (children hidden).
func (n *Node) Collapse() *Node { n.Collapsed = true; return n }

// AppendChild appends one or more children and returns the parent for chaining.
func (n *Node) AppendChild(children ...*Node) *Node {
	n.Children = append(n.Children, children...)
	return n
}

// Tree is a hierarchical tree chart.
type Tree struct {
	base   *chart.BaseChart
	root   *Node
	orient string // 'LR' (default), 'RL', 'TB', 'BT'
	layout string // 'orthogonal' (default) or 'radial'
}

// New creates a new Tree bound to the given id.
func New(id string) *Tree { return &Tree{base: chart.New(id)} }

// Title sets the chart title.
func (t *Tree) Title(s string) *Tree { t.base.SetTitle(s); return t }

// Root sets the root node of the tree.
func (t *Tree) Root(root *Node) *Tree { t.root = root; return t }

// Orient sets the layout orientation: LR | RL | TB | BT.
func (t *Tree) Orient(o string) *Tree { t.orient = o; return t }

// Layout sets the layout mode: 'orthogonal' (default) or 'radial'.
func (t *Tree) Layout(l string) *Tree { t.layout = l; return t }

func (t *Tree) Compact() *Tree              { t.base.SetCompact(); return t }
func (t *Tree) WithDisplay(d string) *Tree  { t.base.SetDisplay(d); return t }
func (t *Tree) SetURL(u *url.Url) *Tree     { t.base.SetURL(u); return t }
func (t *Tree) WithReload(r bool) *Tree     { t.base.SetReload(r); return t }

func (t *Tree) Print(ctx *core.UiContext) map[string]any {
	if t.base.HasURL() {
		return t.base.Envelope("tree", t.base.PrintAjaxData(), nil)
	}
	return t.base.Envelope("tree", t.printData(ctx), nil)
}

func (t *Tree) PrintData(ctx *core.UiContext) map[string]any { return t.printData(ctx) }

func (t *Tree) DataResponse(ctx *core.UiContext) response.DataResult {
	return response.NewJSONDataResult(t.PrintData(ctx))
}

func (t *Tree) printData(_ *core.UiContext) map[string]any {
	data := t.base.PrintBaseData(nil)

	if t.root != nil {
		data["root"] = nodeToMap(t.root)
	}
	if t.orient != "" {
		data["orient"] = t.orient
	}
	if t.layout != "" {
		data["layout"] = t.layout
	}
	return data
}

func nodeToMap(n *Node) map[string]any {
	m := map[string]any{"name": n.Name}
	if n.Value != nil {
		m["value"] = *n.Value
	}
	if n.Collapsed {
		m["collapsed"] = true
	}
	if len(n.Children) > 0 {
		children := make([]map[string]any, len(n.Children))
		for i, c := range n.Children {
			children[i] = nodeToMap(c)
		}
		m["children"] = children
	}
	return m
}
