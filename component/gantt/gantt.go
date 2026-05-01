// Package gantt provides a Gantt-chart component (rendered via echarts custom series).
package gantt

import (
	"github.com/xiriframework/xiri-go/component/chart"
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
	"github.com/xiriframework/xiri-go/response"
)

// Task is one bar on the Gantt chart.
//
// Row is the index into Gantt.rows. Start/End are unix milliseconds.
type Task struct {
	Row   int
	Name  string
	Start int64
	End   int64
	Color core.Color
}

// Gantt is a Gantt chart with one row per category and one bar per task.
type Gantt struct {
	base *chart.BaseChart

	rows       []string
	tasks      []Task
	rangeStart *int64
	rangeEnd   *int64
}

// New creates a new Gantt bound to the given id.
func New(id string) *Gantt { return &Gantt{base: chart.New(id)} }

// Title sets the chart title.
func (g *Gantt) Title(t string) *Gantt { g.base.SetTitle(t); return g }

// Rows sets the category rows (Y-axis), top-down.
func (g *Gantt) Rows(rows ...string) *Gantt { g.rows = rows; return g }

// Task appends one task bar.
func (g *Gantt) Task(row int, name string, start, end int64, color core.Color) *Gantt {
	g.tasks = append(g.tasks, Task{Row: row, Name: name, Start: start, End: end, Color: color})
	return g
}

// XRange sets the visible X-axis range (unix ms). If unset, echarts auto-fits.
func (g *Gantt) XRange(start, end int64) *Gantt { g.rangeStart = &start; g.rangeEnd = &end; return g }

func (g *Gantt) Compact() *Gantt              { g.base.SetCompact(); return g }
func (g *Gantt) WithDisplay(d string) *Gantt  { g.base.SetDisplay(d); return g }
func (g *Gantt) SetURL(u *url.Url) *Gantt     { g.base.SetURL(u); return g }
func (g *Gantt) WithReload(r bool) *Gantt     { g.base.SetReload(r); return g }

func (g *Gantt) Print(ctx *core.UiContext) map[string]any {
	if g.base.HasURL() {
		return g.base.Envelope("gantt", g.base.PrintAjaxData(), nil)
	}
	return g.base.Envelope("gantt", g.printData(ctx), nil)
}

func (g *Gantt) PrintData(ctx *core.UiContext) map[string]any { return g.printData(ctx) }

func (g *Gantt) DataResponse(ctx *core.UiContext) response.DataResult {
	return response.NewJSONDataResult(g.PrintData(ctx))
}

func (g *Gantt) printData(ctx *core.UiContext) map[string]any {
	data := g.base.PrintBaseData(ctx)

	data["rows"] = g.rows

	tasks := make([]map[string]any, len(g.tasks))
	for i, t := range g.tasks {
		m := map[string]any{
			"row":   t.Row,
			"name":  t.Name,
			"start": t.Start,
			"end":   t.End,
		}
		if t.Color != "" {
			m["color"] = string(t.Color)
		}
		tasks[i] = m
	}
	data["tasks"] = tasks

	if g.rangeStart != nil {
		data["rangeStart"] = *g.rangeStart
	}
	if g.rangeEnd != nil {
		data["rangeEnd"] = *g.rangeEnd
	}
	return data
}
