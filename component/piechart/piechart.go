// Package piechart provides a pie/donut chart component for the Angular
// frontend. The `xiri-ng` side renders via echarts (optional peerDependency).
package piechart

import (
	"github.com/xiriframework/xiri-go/component/chart"
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
	"github.com/xiriframework/xiri-go/response"
)

// Slice represents one pie slice.
type Slice struct {
	Name  string
	Value float64
	Color core.Color
}

// PieChart is a pie or donut chart. Set Nightingale() for a rose-style
// rendering where each slice's radius scales with its value.
type PieChart struct {
	base *chart.BaseChart

	slices          []Slice
	donut           bool
	nightingale     bool
	nightingaleType string // "radius" (default) or "area"
}

// New creates a new PieChart bound to the given id.
func New(id string) *PieChart {
	return &PieChart{base: chart.New(id)}
}

// Title sets the chart title (translation key).
func (pc *PieChart) Title(t string) *PieChart { pc.base.SetTitle(t); return pc }

// Slice appends a slice to the chart.
func (pc *PieChart) Slice(name string, value float64, color core.Color) *PieChart {
	pc.slices = append(pc.slices, Slice{Name: name, Value: value, Color: color})
	return pc
}

// Donut renders the chart as a donut (ring) instead of a full pie.
func (pc *PieChart) Donut() *PieChart { pc.donut = true; return pc }

// Nightingale renders the chart in rose-style: each slice keeps the same
// angular share but its radius scales with the value. Default 'radius' mode;
// pass NightingaleArea() for area-scaling instead.
func (pc *PieChart) Nightingale() *PieChart { pc.nightingale = true; return pc }

// NightingaleArea is like Nightingale but uses 'area' scaling (slice areas
// are proportional to values).
func (pc *PieChart) NightingaleArea() *PieChart {
	pc.nightingale = true
	pc.nightingaleType = "area"
	return pc
}

// Compact switches the chart into compact mode.
func (pc *PieChart) Compact() *PieChart { pc.base.SetCompact(); return pc }

// WithDisplay sets the layout class.
func (pc *PieChart) WithDisplay(d string) *PieChart { pc.base.SetDisplay(d); return pc }

// SetURL switches the chart to AJAX mode.
func (pc *PieChart) SetURL(u *url.Url) *PieChart { pc.base.SetURL(u); return pc }

// WithReload enables periodic reload in AJAX mode.
func (pc *PieChart) WithReload(r bool) *PieChart { pc.base.SetReload(r); return pc }

// Print returns the JSON envelope for the pie chart.
func (pc *PieChart) Print(ctx *core.UiContext) map[string]any {
	if pc.base.HasURL() {
		return pc.base.Envelope("piechart", pc.base.PrintAjaxData(), nil)
	}
	return pc.base.Envelope("piechart", pc.printData(ctx), nil)
}

// PrintData returns only the data map (for AJAX endpoints).
func (pc *PieChart) PrintData(ctx *core.UiContext) map[string]any {
	return pc.printData(ctx)
}

// DataResponse wraps PrintData in a {"data": ...} envelope for DataResult.
func (pc *PieChart) DataResponse(ctx *core.UiContext) response.DataResult {
	return response.NewJSONDataResult(pc.PrintData(ctx))
}

func (pc *PieChart) printData(ctx *core.UiContext) map[string]any {
	data := pc.base.PrintBaseData(ctx)

	if pc.donut {
		data["donut"] = true
	}
	if pc.nightingale {
		data["nightingale"] = true
		if pc.nightingaleType != "" {
			data["nightingaleType"] = pc.nightingaleType
		}
	}

	slices := make([]map[string]any, len(pc.slices))
	for i, s := range pc.slices {
		m := map[string]any{
			"name":  s.Name,
			"value": s.Value,
		}
		if s.Color != "" {
			m["color"] = string(s.Color)
		}
		slices[i] = m
	}
	data["slices"] = slices

	return data
}
