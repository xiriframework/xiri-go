// Package chart provides shared building blocks for echarts-based chart
// components (barchart, linechart, piechart, gaugechart, …). Each specialized
// chart package embeds a BaseChart for the cross-cutting concerns.
package chart

import (
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
)

// BaseChart holds the cross-cutting fields every chart component supports.
//
// Specialized chart structs (barchart.BarChart, linechart.LineChart, …) hold a
// `*BaseChart` field and forward the cross-cutting builder methods (Title,
// Color, Compact, AjaxURL, …) to it. The Print() method of the specialized
// chart calls Envelope() to build the JSON wrapper and adds its own
// type-specific data into the data map.
type BaseChart struct {
	id           string
	title        string
	color        *core.Color
	compact      bool
	display      *string
	url          *url.Url
	reload       *bool
	xLabelRotate *int
}

// New returns a new BaseChart bound to the given id.
func New(id string) *BaseChart {
	return &BaseChart{id: id}
}

// ID returns the chart's id.
func (b *BaseChart) ID() string { return b.id }

// SetTitle sets the chart title (translation key).
func (b *BaseChart) SetTitle(t string) { b.title = t }

// SetColor sets the default series color.
func (b *BaseChart) SetColor(c core.Color) { b.color = &c }

// SetCompact toggles the compact mode (no own mat-card surface, tighter grid).
func (b *BaseChart) SetCompact() { b.compact = true }

// SetDisplay sets the display/layout class.
func (b *BaseChart) SetDisplay(d string) { b.display = &d }

// SetURL switches the chart to AJAX mode. The frontend will load data lazily.
func (b *BaseChart) SetURL(u *url.Url) { b.url = u }

// SetReload toggles periodic reload in AJAX mode.
func (b *BaseChart) SetReload(r bool) { b.reload = &r }

// SetXLabelRotate rotates the category X-axis labels by deg degrees (-90..90).
// Useful for narrow bars/columns with long labels. 90 = vertical.
func (b *BaseChart) SetXLabelRotate(deg int) { b.xLabelRotate = &deg }

// HasURL reports whether the chart is in AJAX mode.
func (b *BaseChart) HasURL() bool { return b.url != nil }

// Display returns the configured display class (or nil).
func (b *BaseChart) Display() *string { return b.display }

// PrintAjaxData returns the data map for AJAX mode (only url + reload).
func (b *BaseChart) PrintAjaxData() map[string]any {
	data := map[string]any{"url": b.url.PrintPrefix()}
	if b.reload != nil {
		data["reload"] = *b.reload
	}
	return data
}

// PrintBaseData populates the cross-cutting data fields (title, color,
// compact). Specialized charts call this and then add their own fields.
func (b *BaseChart) PrintBaseData(ctx *core.UiContext) map[string]any {
	data := map[string]any{}
	if b.title != "" {
		data["title"] = core.Translate(ctx, b.title)
	}
	if b.color != nil {
		data["color"] = string(*b.color)
	}
	if b.compact {
		data["compact"] = true
	}
	if b.xLabelRotate != nil {
		data["xLabelRotate"] = *b.xLabelRotate
	}
	return data
}

// Envelope builds the standard `{type, id, display, data}` envelope.
// `extra` may carry type-specific top-level keys (e.g. "mode" for barchart).
func (b *BaseChart) Envelope(typ string, data map[string]any, extra map[string]any) map[string]any {
	out := map[string]any{
		"type": typ,
		"id":   b.id,
		"data": data,
	}
	if b.display != nil {
		out["display"] = *b.display
	}
	for k, v := range extra {
		out[k] = v
	}
	return out
}
