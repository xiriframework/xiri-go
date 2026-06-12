// Package heatmap provides a 2D matrix-heatmap component (echarts heatmap series).
package heatmap

import (
	"github.com/xiriframework/xiri-go/component/chart"
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
	"github.com/xiriframework/xiri-go/response"
)

// Cell is one filled (x,y) position with its numeric value.
type Cell struct {
	X     int     // index into XLabels
	Y     int     // index into YLabels
	Value float64
}

// Heatmap is a 2D matrix heatmap.
type Heatmap struct {
	base *chart.BaseChart

	xLabels    []string
	yLabels    []string
	cells      []Cell
	min        *float64
	max        *float64
	colorRange [2]string
	showValues bool
}

// New creates a new Heatmap bound to the given id.
func New(id string) *Heatmap { return &Heatmap{base: chart.New(id)} }

// Title sets the chart title.
func (h *Heatmap) Title(t string) *Heatmap { h.base.SetTitle(t); return h }

// XLabels sets the column labels (X axis).
func (h *Heatmap) XLabels(labels ...string) *Heatmap { h.xLabels = labels; return h }

// YLabels sets the row labels (Y axis).
func (h *Heatmap) YLabels(labels ...string) *Heatmap { h.yLabels = labels; return h }

// Cell adds one cell at (x,y) with the given value.
func (h *Heatmap) Cell(x, y int, value float64) *Heatmap {
	h.cells = append(h.cells, Cell{X: x, Y: y, Value: value})
	return h
}

// Range sets the color-scale min/max. If unset, derived from cell values.
func (h *Heatmap) Range(min, max float64) *Heatmap { h.min = &min; h.max = &max; return h }

// ColorRange sets the [low, high] CSS colors used by the visualMap.
func (h *Heatmap) ColorRange(low, high string) *Heatmap {
	h.colorRange = [2]string{low, high}
	return h
}

// ShowValues prints values inside cells (useful for small grids).
func (h *Heatmap) ShowValues() *Heatmap { h.showValues = true; return h }

// XLabelRotate rotates the X-axis (column) labels by deg degrees (-90..90).
// 90 = vertical. Helps when columns are narrow and labels are long.
func (h *Heatmap) XLabelRotate(deg int) *Heatmap { h.base.SetXLabelRotate(deg); return h }

// Compact, WithDisplay, SetURL, WithReload — Forward to BaseChart.
func (h *Heatmap) Compact() *Heatmap                  { h.base.SetCompact(); return h }
func (h *Heatmap) WithDisplay(d string) *Heatmap      { h.base.SetDisplay(d); return h }
func (h *Heatmap) SetURL(u *url.Url) *Heatmap         { h.base.SetURL(u); return h }
func (h *Heatmap) WithReload(r bool) *Heatmap         { h.base.SetReload(r); return h }

// Print, PrintData, DataResponse.
func (h *Heatmap) Print(ctx *core.UiContext) map[string]any {
	if h.base.HasURL() {
		return h.base.Envelope("heatmap", h.base.PrintAjaxData(), nil)
	}
	return h.base.Envelope("heatmap", h.printData(ctx), nil)
}

func (h *Heatmap) PrintData(ctx *core.UiContext) map[string]any { return h.printData(ctx) }

func (h *Heatmap) DataResponse(ctx *core.UiContext) response.DataResult {
	return response.NewJSONDataResult(h.PrintData(ctx))
}

func (h *Heatmap) printData(ctx *core.UiContext) map[string]any {
	data := h.base.PrintBaseData(ctx)

	data["xLabels"] = h.xLabels
	data["yLabels"] = h.yLabels

	cells := make([]map[string]any, len(h.cells))
	for i, c := range h.cells {
		cells[i] = map[string]any{"x": c.X, "y": c.Y, "value": c.Value}
	}
	data["cells"] = cells

	if h.min != nil {
		data["min"] = *h.min
	}
	if h.max != nil {
		data["max"] = *h.max
	}
	if h.colorRange[0] != "" {
		data["colorRange"] = []string{h.colorRange[0], h.colorRange[1]}
	}
	if h.showValues {
		data["showValues"] = true
	}
	return data
}
