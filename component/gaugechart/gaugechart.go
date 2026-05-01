// Package gaugechart provides a gauge chart component for the Angular
// frontend. The `xiri-ng` side renders via echarts (optional peerDependency).
package gaugechart

import (
	"github.com/xiriframework/xiri-go/component/chart"
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
	"github.com/xiriframework/xiri-go/response"
)

// GaugeChart shows a single value on a circular gauge.
type GaugeChart struct {
	base *chart.BaseChart

	value float64
	min   *float64
	max   *float64
	label string // optional unit label below the value (e.g. "%")
}

// New creates a new GaugeChart bound to the given id.
func New(id string) *GaugeChart {
	return &GaugeChart{base: chart.New(id)}
}

// Title sets the chart title (translation key).
func (g *GaugeChart) Title(t string) *GaugeChart { g.base.SetTitle(t); return g }

// Value sets the current value displayed by the gauge.
func (g *GaugeChart) Value(v float64) *GaugeChart { g.value = v; return g }

// Range sets the gauge's min/max range. Default is 0..100.
func (g *GaugeChart) Range(min, max float64) *GaugeChart { g.min = &min; g.max = &max; return g }

// Label sets the unit label below the value (e.g. "%", "MB/s").
func (g *GaugeChart) Label(l string) *GaugeChart { g.label = l; return g }

// Color sets the progress arc color.
func (g *GaugeChart) Color(c core.Color) *GaugeChart { g.base.SetColor(c); return g }

// Compact switches the chart into compact mode (no axis labels, no own card).
func (g *GaugeChart) Compact() *GaugeChart { g.base.SetCompact(); return g }

// WithDisplay sets the layout class.
func (g *GaugeChart) WithDisplay(d string) *GaugeChart { g.base.SetDisplay(d); return g }

// SetURL switches the chart to AJAX mode.
func (g *GaugeChart) SetURL(u *url.Url) *GaugeChart { g.base.SetURL(u); return g }

// WithReload enables periodic reload in AJAX mode.
func (g *GaugeChart) WithReload(r bool) *GaugeChart { g.base.SetReload(r); return g }

// Print returns the JSON envelope for the gauge.
func (g *GaugeChart) Print(ctx *core.UiContext) map[string]any {
	if g.base.HasURL() {
		return g.base.Envelope("gaugechart", g.base.PrintAjaxData(), nil)
	}
	return g.base.Envelope("gaugechart", g.printData(ctx), nil)
}

// PrintData returns only the data map (for AJAX endpoints).
func (g *GaugeChart) PrintData(ctx *core.UiContext) map[string]any {
	return g.printData(ctx)
}

// DataResponse wraps PrintData in a {"data": ...} envelope for DataResult.
func (g *GaugeChart) DataResponse(ctx *core.UiContext) response.DataResult {
	return response.NewJSONDataResult(g.PrintData(ctx))
}

func (g *GaugeChart) printData(ctx *core.UiContext) map[string]any {
	data := g.base.PrintBaseData(ctx)

	data["value"] = g.value
	if g.min != nil {
		data["min"] = *g.min
	}
	if g.max != nil {
		data["max"] = *g.max
	}
	if g.label != "" {
		data["label"] = g.label
	}
	return data
}
