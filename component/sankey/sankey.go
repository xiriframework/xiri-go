// Package sankey provides a sankey-diagram component (echarts sankey series).
package sankey

import (
	"github.com/xiriframework/xiri-go/component/chart"
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
	"github.com/xiriframework/xiri-go/response"
)

// Node is a sankey node identified by its name.
type Node struct {
	Name  string
	Color core.Color
}

// Link is a directed flow from source to target with a magnitude (value).
type Link struct {
	Source string
	Target string
	Value  float64
}

// Sankey is a sankey diagram.
type Sankey struct {
	base *chart.BaseChart

	nodes  []Node
	links  []Link
	orient string // 'horizontal' (default) | 'vertical'
}

// New creates a new Sankey bound to the given id.
func New(id string) *Sankey { return &Sankey{base: chart.New(id)} }

// Title sets the chart title.
func (s *Sankey) Title(t string) *Sankey { s.base.SetTitle(t); return s }

// Node adds one node by name (and optional color).
func (s *Sankey) Node(name string, color core.Color) *Sankey {
	s.nodes = append(s.nodes, Node{Name: name, Color: color})
	return s
}

// Link adds a flow between two named nodes.
func (s *Sankey) Link(source, target string, value float64) *Sankey {
	s.links = append(s.links, Link{Source: source, Target: target, Value: value})
	return s
}

// Vertical switches the layout to vertical orientation.
func (s *Sankey) Vertical() *Sankey { s.orient = "vertical"; return s }

func (s *Sankey) Compact() *Sankey              { s.base.SetCompact(); return s }
func (s *Sankey) WithDisplay(d string) *Sankey  { s.base.SetDisplay(d); return s }
func (s *Sankey) SetURL(u *url.Url) *Sankey     { s.base.SetURL(u); return s }
func (s *Sankey) WithReload(r bool) *Sankey     { s.base.SetReload(r); return s }

func (s *Sankey) Print(ctx *core.UiContext) map[string]any {
	if s.base.HasURL() {
		return s.base.Envelope("sankey", s.base.PrintAjaxData(), nil)
	}
	return s.base.Envelope("sankey", s.printData(ctx), nil)
}

func (s *Sankey) PrintData(ctx *core.UiContext) map[string]any { return s.printData(ctx) }

func (s *Sankey) DataResponse(ctx *core.UiContext) response.DataResult {
	return response.NewJSONDataResult(s.PrintData(ctx))
}

func (s *Sankey) printData(ctx *core.UiContext) map[string]any {
	data := s.base.PrintBaseData(ctx)

	nodes := make([]map[string]any, len(s.nodes))
	for i, n := range s.nodes {
		m := map[string]any{"name": n.Name}
		if n.Color != "" {
			m["color"] = string(n.Color)
		}
		nodes[i] = m
	}
	data["nodes"] = nodes

	links := make([]map[string]any, len(s.links))
	for i, l := range s.links {
		links[i] = map[string]any{"source": l.Source, "target": l.Target, "value": l.Value}
	}
	data["links"] = links

	if s.orient != "" {
		data["orient"] = s.orient
	}
	return data
}
