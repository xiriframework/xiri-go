// Package calendar provides a calendar-heatmap component (echarts calendar
// coordinate + heatmap series).
package calendar

import (
	"github.com/xiriframework/xiri-go/component/chart"
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
	"github.com/xiriframework/xiri-go/response"
)

// Cell is one calendar day with a value.
type Cell struct {
	Date  string  // 'YYYY-MM-DD'
	Value float64
}

// Calendar is a calendar-heatmap (GitHub-style activity grid).
type Calendar struct {
	base *chart.BaseChart

	rangeStart string // 'YYYY' or 'YYYY-MM-DD'
	rangeEnd   string // optional second date for a range; empty = single year
	cells      []Cell
	min        *float64
	max        *float64
	colorRange [2]string
	cellSize   *int
}

// New creates a new Calendar bound to the given id.
func New(id string) *Calendar { return &Calendar{base: chart.New(id)} }

// Title sets the chart title.
func (c *Calendar) Title(t string) *Calendar { c.base.SetTitle(t); return c }

// Year sets the calendar to display a single year (e.g. "2025").
func (c *Calendar) Year(year string) *Calendar { c.rangeStart = year; c.rangeEnd = ""; return c }

// Range sets a date range with start/end (each "YYYY-MM-DD").
func (c *Calendar) Range(start, end string) *Calendar { c.rangeStart = start; c.rangeEnd = end; return c }

// Cell adds one (date, value) pair.
func (c *Calendar) Cell(date string, value float64) *Calendar {
	c.cells = append(c.cells, Cell{Date: date, Value: value})
	return c
}

// MinMax sets the color-scale range. If unset, derived from cell values.
func (c *Calendar) MinMax(min, max float64) *Calendar { c.min = &min; c.max = &max; return c }

// ColorRange sets the [low, high] CSS colors used by the visualMap.
func (c *Calendar) ColorRange(low, high string) *Calendar {
	c.colorRange = [2]string{low, high}
	return c
}

// CellSize sets the cell pixel size.
func (c *Calendar) CellSize(px int) *Calendar { c.cellSize = &px; return c }

func (c *Calendar) Compact() *Calendar              { c.base.SetCompact(); return c }
func (c *Calendar) WithDisplay(d string) *Calendar  { c.base.SetDisplay(d); return c }
func (c *Calendar) SetURL(u *url.Url) *Calendar     { c.base.SetURL(u); return c }
func (c *Calendar) WithReload(r bool) *Calendar     { c.base.SetReload(r); return c }

func (c *Calendar) Print(ctx *core.UiContext) map[string]any {
	if c.base.HasURL() {
		return c.base.Envelope("calendar", c.base.PrintAjaxData(), nil)
	}
	return c.base.Envelope("calendar", c.printData(ctx), nil)
}

func (c *Calendar) PrintData(ctx *core.UiContext) map[string]any { return c.printData(ctx) }

func (c *Calendar) DataResponse(ctx *core.UiContext) response.DataResult {
	return response.NewJSONDataResult(c.PrintData(ctx))
}

func (c *Calendar) printData(ctx *core.UiContext) map[string]any {
	data := c.base.PrintBaseData(ctx)

	if c.rangeEnd != "" {
		data["range"] = []string{c.rangeStart, c.rangeEnd}
	} else if c.rangeStart != "" {
		data["range"] = c.rangeStart
	}

	cells := make([]map[string]any, len(c.cells))
	for i, cell := range c.cells {
		cells[i] = map[string]any{"date": cell.Date, "value": cell.Value}
	}
	data["cells"] = cells

	if c.min != nil {
		data["min"] = *c.min
	}
	if c.max != nil {
		data["max"] = *c.max
	}
	if c.colorRange[0] != "" {
		data["colorRange"] = []string{c.colorRange[0], c.colorRange[1]}
	}
	if c.cellSize != nil {
		data["cellSize"] = *c.cellSize
	}
	return data
}
