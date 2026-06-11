// Package barchart provides a bar-chart component for the Angular frontend.
//
// Three modes are supported:
//
//   - ModeSimple: one value per category. Use Bar(label, value).
//     Mockup example: "Weekly activities" Mon-Sun.
//
//   - ModeStacked: each category is split into colored segments stacked vertically.
//     Use StackedBar(label, segments...). Mockup example: "Vehicle strain"
//     where every weekday has green/yellow/red proportions.
//
//   - ModeHeatmap: a long sparse activity strip rendered as many thin bars over
//     a time axis. Use Point(timeMs, value). Mockup example: "Engine system"
//     activity track.
//
// The frontend renders all three via ECharts, but the JSON shape is library-agnostic.
package barchart

import (
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
	"github.com/xiriframework/xiri-go/response"
)

// Mode controls how the bar chart renders.
type Mode string

const (
	ModeSimple  Mode = "simple"
	ModeStacked Mode = "stacked"
	ModeHeatmap Mode = "heatmap"
)

// Segment is one colored slice of a stacked bar.
//
// Name is shown in tooltips. Color is the segment fill.
type Segment struct {
	Value float64
	Color core.Color
	Name  string
}

// Seg is a short constructor for Segment without a tooltip name.
func Seg(value float64, color core.Color) Segment {
	return Segment{Value: value, Color: color}
}

// SegNamed is a constructor for Segment with a tooltip name.
func SegNamed(value float64, name string, color core.Color) Segment {
	return Segment{Value: value, Color: color, Name: name}
}

// simpleBar holds one labeled value in ModeSimple.
type simpleBar struct {
	label string
	name  string
	value float64
	url   *url.Url
}

// stackedBar holds one labeled bar with stacked segments in ModeStacked.
type stackedBar struct {
	label    string
	name     string
	segments []Segment
	url      *url.Url
}

// heatPoint is one (time, value) pair in ModeHeatmap. Time is unix milliseconds.
type heatPoint struct {
	timeMs int64
	value  float64
	name   string
	url    *url.Url
}

// BarChart is a bar-chart component.
type BarChart struct {
	id    string
	mode  Mode
	title string

	simple  []simpleBar
	stacked []stackedBar
	heat    []heatPoint

	yMin    *float64
	yMax    *float64
	color   *core.Color // override single-color (ModeSimple/ModeHeatmap)
	display *string
	url     *url.Url
	reload  *bool
	compact bool
}

// New creates a new BarChart with the given id. Default mode is ModeSimple.
func New(id string) *BarChart {
	return &BarChart{id: id, mode: ModeSimple}
}

// Mode sets the rendering mode (simple/stacked/heatmap).
func (b *BarChart) Mode(m Mode) *BarChart {
	b.mode = m
	return b
}

// Title sets the chart title (translation key).
func (b *BarChart) Title(t string) *BarChart {
	b.title = t
	return b
}

// Bar appends one labeled value (use with ModeSimple).
//
// label is the short axis label. For longer tooltip names, use BarNamed.
func (b *BarChart) Bar(label string, value float64) *BarChart {
	b.simple = append(b.simple, simpleBar{label: label, value: value})
	return b
}

// BarNamed appends a labeled value with a separate tooltip name (use with ModeSimple).
// label is shown on the axis, name in tooltips.
func (b *BarChart) BarNamed(label, name string, value float64) *BarChart {
	b.simple = append(b.simple, simpleBar{label: label, name: name, value: value})
	return b
}

// StackedBar appends one labeled bar with colored segments (use with ModeStacked).
func (b *BarChart) StackedBar(label string, segments ...Segment) *BarChart {
	b.stacked = append(b.stacked, stackedBar{label: label, segments: segments})
	return b
}

// StackedBarNamed appends a stacked bar with a separate tooltip name (use with ModeStacked).
func (b *BarChart) StackedBarNamed(label, name string, segments ...Segment) *BarChart {
	b.stacked = append(b.stacked, stackedBar{label: label, name: name, segments: segments})
	return b
}

// Point appends one (time, value) pair (use with ModeHeatmap). timeMs is unix milliseconds.
func (b *BarChart) Point(timeMs int64, value float64) *BarChart {
	b.heat = append(b.heat, heatPoint{timeMs: timeMs, value: value})
	return b
}

// PointNamed is like Point but attaches a tooltip name to the point.
func (b *BarChart) PointNamed(timeMs int64, name string, value float64) *BarChart {
	b.heat = append(b.heat, heatPoint{timeMs: timeMs, value: value, name: name})
	return b
}

// Link attaches a navigation target to the most recently added bar/point.
// The frontend makes that bar clickable and navigates to the url on click.
// Call it directly after the Bar/StackedBar/Point it belongs to.
func (b *BarChart) Link(u *url.Url) *BarChart {
	switch b.mode {
	case ModeStacked:
		if n := len(b.stacked); n > 0 {
			b.stacked[n-1].url = u
		}
	case ModeHeatmap:
		if n := len(b.heat); n > 0 {
			b.heat[n-1].url = u
		}
	default:
		if n := len(b.simple); n > 0 {
			b.simple[n-1].url = u
		}
	}
	return b
}

// YAxis sets the Y-axis range.
func (b *BarChart) YAxis(min, max float64) *BarChart {
	b.yMin = &min
	b.yMax = &max
	return b
}

// Color sets a single color for ModeSimple / ModeHeatmap (ignored by ModeStacked).
func (b *BarChart) Color(c core.Color) *BarChart {
	b.color = &c
	return b
}

// WithDisplay sets the display/layout class.
func (b *BarChart) WithDisplay(display string) *BarChart {
	b.display = &display
	return b
}

// SetURL switches the chart to AJAX mode — the frontend loads data from this url.
func (b *BarChart) SetURL(u *url.Url) *BarChart {
	b.url = u
	return b
}

// WithReload enables periodic reload of the data in AJAX mode.
func (b *BarChart) WithReload(reload bool) *BarChart {
	b.reload = &reload
	return b
}

// Compact switches the BarChart into compact mode — no own mat-card surface.
// Use when nesting a barchart inside another card to avoid card-in-card visuals.
func (b *BarChart) Compact() *BarChart {
	b.compact = true
	return b
}

// Print returns the JSON representation of the bar chart.
func (b *BarChart) Print(ctx *core.UiContext) map[string]any {
	var data map[string]any
	if b.url != nil {
		data = map[string]any{"url": b.url.PrintPrefix()}
		if b.reload != nil {
			data["reload"] = *b.reload
		}
	} else {
		data = b.printData(ctx)
	}

	result := map[string]any{
		"type": "barchart",
		"id":   b.id,
		"mode": string(b.mode),
		"data": data,
	}
	if b.display != nil {
		result["display"] = *b.display
	}
	return result
}

// PrintData returns only the data portion (for AJAX endpoints).
func (b *BarChart) PrintData(ctx *core.UiContext) map[string]any {
	return b.printData(ctx)
}

// DataResponse wraps PrintData in a DataResult envelope.
func (b *BarChart) DataResponse(ctx *core.UiContext) response.DataResult {
	return response.NewJSONDataResult(b.PrintData(ctx))
}

func (b *BarChart) printData(ctx *core.UiContext) map[string]any {
	data := map[string]any{}

	if b.title != "" {
		data["title"] = core.Translate(ctx, b.title)
	}
	if b.yMin != nil {
		data["yMin"] = *b.yMin
	}
	if b.yMax != nil {
		data["yMax"] = *b.yMax
	}
	if b.color != nil {
		data["color"] = string(*b.color)
	}
	if b.compact {
		data["compact"] = true
	}

	switch b.mode {
	case ModeStacked:
		data["bars"] = b.stackedBarsJSON()
	case ModeHeatmap:
		data["points"] = b.heatPointsJSON()
	default:
		data["bars"] = b.simpleBarsJSON()
	}

	return data
}

func (b *BarChart) simpleBarsJSON() []map[string]any {
	out := make([]map[string]any, len(b.simple))
	for i, bar := range b.simple {
		m := map[string]any{
			"label": bar.label,
			"value": bar.value,
		}
		if bar.name != "" {
			m["name"] = bar.name
		}
		if bar.url != nil {
			m["url"] = bar.url.Print()
		}
		out[i] = m
	}
	return out
}

func (b *BarChart) stackedBarsJSON() []map[string]any {
	out := make([]map[string]any, len(b.stacked))
	for i, bar := range b.stacked {
		segs := make([]map[string]any, len(bar.segments))
		for j, s := range bar.segments {
			seg := map[string]any{
				"value": s.Value,
				"color": string(s.Color),
			}
			if s.Name != "" {
				seg["name"] = s.Name
			}
			segs[j] = seg
		}
		m := map[string]any{
			"label":    bar.label,
			"segments": segs,
		}
		if bar.name != "" {
			m["name"] = bar.name
		}
		if bar.url != nil {
			m["url"] = bar.url.Print()
		}
		out[i] = m
	}
	return out
}

func (b *BarChart) heatPointsJSON() []map[string]any {
	out := make([]map[string]any, len(b.heat))
	for i, p := range b.heat {
		m := map[string]any{
			"time":  p.timeMs,
			"value": p.value,
		}
		if p.name != "" {
			m["name"] = p.name
		}
		if p.url != nil {
			m["url"] = p.url.Print()
		}
		out[i] = m
	}
	return out
}
