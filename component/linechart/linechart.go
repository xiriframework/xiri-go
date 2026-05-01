// Package linechart provides a multi-line chart component for the Angular
// frontend. The `xiri-ng` side renders via echarts (optional peerDependency).
package linechart

import (
	"github.com/xiriframework/xiri-go/component/chart"
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
	"github.com/xiriframework/xiri-go/response"
)

// Line is a single data series in a line chart.
type Line struct {
	Name   string
	Values []float64
	Color  core.Color
	Dashed bool
	Area   bool
}

// LineChart is a multi-line chart with a shared category x-axis.
type LineChart struct {
	base *chart.BaseChart

	xLabels []string
	lines   []Line

	yMin   *float64
	yMax   *float64
	smooth bool
}

// New creates a new LineChart bound to the given id.
func New(id string) *LineChart {
	return &LineChart{base: chart.New(id)}
}

// Title sets the chart title (translation key).
func (lc *LineChart) Title(t string) *LineChart { lc.base.SetTitle(t); return lc }

// XLabels sets the category labels on the X-axis.
func (lc *LineChart) XLabels(labels ...string) *LineChart { lc.xLabels = labels; return lc }

// Line appends a data series. Use the chained Color/Dashed/Area methods for styling.
func (lc *LineChart) Line(name string, values []float64) *LineBuilder {
	lc.lines = append(lc.lines, Line{Name: name, Values: values})
	return &LineBuilder{lc: lc, idx: len(lc.lines) - 1}
}

// YAxis sets the Y-axis range.
func (lc *LineChart) YAxis(min, max float64) *LineChart { lc.yMin = &min; lc.yMax = &max; return lc }

// Smooth enables smoothed (curved) lines.
func (lc *LineChart) Smooth() *LineChart { lc.smooth = true; return lc }

// Color sets the default color (used as fallback for individual lines).
func (lc *LineChart) Color(c core.Color) *LineChart { lc.base.SetColor(c); return lc }

// Compact switches the chart into compact mode (no own mat-card, tight grid).
func (lc *LineChart) Compact() *LineChart { lc.base.SetCompact(); return lc }

// WithDisplay sets the layout class.
func (lc *LineChart) WithDisplay(d string) *LineChart { lc.base.SetDisplay(d); return lc }

// SetURL switches the chart to AJAX mode.
func (lc *LineChart) SetURL(u *url.Url) *LineChart { lc.base.SetURL(u); return lc }

// WithReload enables periodic reload in AJAX mode.
func (lc *LineChart) WithReload(r bool) *LineChart { lc.base.SetReload(r); return lc }

// LineBuilder allows chaining color/style methods on the most recently added line.
type LineBuilder struct {
	lc  *LineChart
	idx int
}

// Color sets the line color.
func (b *LineBuilder) Color(c core.Color) *LineBuilder {
	b.lc.lines[b.idx].Color = c
	return b
}

// Dashed renders the line as dashed instead of solid.
func (b *LineBuilder) Dashed() *LineBuilder {
	b.lc.lines[b.idx].Dashed = true
	return b
}

// Area fills the area below the line.
func (b *LineBuilder) Area() *LineBuilder {
	b.lc.lines[b.idx].Area = true
	return b
}

// Done returns the underlying LineChart for further chaining.
func (b *LineBuilder) Done() *LineChart { return b.lc }

// LineChart fluent shortcuts via LineBuilder (so you can keep chaining off Line()):
func (b *LineBuilder) Line(name string, values []float64) *LineBuilder { return b.lc.Line(name, values) }
func (b *LineBuilder) Title(t string) *LineBuilder                     { b.lc.Title(t); return b }
func (b *LineBuilder) XLabels(labels ...string) *LineBuilder           { b.lc.XLabels(labels...); return b }
func (b *LineBuilder) YAxis(min, max float64) *LineBuilder             { b.lc.YAxis(min, max); return b }
func (b *LineBuilder) Smooth() *LineBuilder                            { b.lc.Smooth(); return b }
func (b *LineBuilder) Compact() *LineBuilder                           { b.lc.Compact(); return b }
func (b *LineBuilder) WithDisplay(d string) *LineBuilder               { b.lc.WithDisplay(d); return b }

// Print returns the JSON envelope for the line chart.
func (lc *LineChart) Print(ctx *core.UiContext) map[string]any {
	if lc.base.HasURL() {
		return lc.base.Envelope("linechart", lc.base.PrintAjaxData(), nil)
	}
	return lc.base.Envelope("linechart", lc.printData(ctx), nil)
}

// PrintData returns only the data map (for AJAX endpoints).
func (lc *LineChart) PrintData(ctx *core.UiContext) map[string]any {
	return lc.printData(ctx)
}

// DataResponse wraps PrintData in a {"data": ...} envelope for DataResult.
func (lc *LineChart) DataResponse(ctx *core.UiContext) response.DataResult {
	return response.NewJSONDataResult(lc.PrintData(ctx))
}

func (lc *LineChart) printData(ctx *core.UiContext) map[string]any {
	data := lc.base.PrintBaseData(ctx)

	if lc.yMin != nil {
		data["yMin"] = *lc.yMin
	}
	if lc.yMax != nil {
		data["yMax"] = *lc.yMax
	}
	if lc.smooth {
		data["smooth"] = true
	}
	if len(lc.xLabels) > 0 {
		data["xLabels"] = lc.xLabels
	}

	lines := make([]map[string]any, len(lc.lines))
	for i, l := range lc.lines {
		m := map[string]any{
			"name":   l.Name,
			"values": l.Values,
		}
		if l.Color != "" {
			m["color"] = string(l.Color)
		}
		if l.Dashed {
			m["dashed"] = true
		}
		if l.Area {
			m["area"] = true
		}
		lines[i] = m
	}
	data["lines"] = lines

	return data
}
