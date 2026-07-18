// Package multistat provides a KPI card that holds multiple numbers in one
// surface — each with its own value, label, color and icon — under a shared
// header. Unlike statgrid (N separate stat cards), multi-stat renders all
// items inside a single card, side by side.
package multistat

import (
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/stat"
	"github.com/xiriframework/xiri-go/component/url"
	"github.com/xiriframework/xiri-go/response"
)

// MultiStat represents a single KPI card holding multiple stat items.
type MultiStat struct {
	items     []*stat.Stat
	title     *string
	icon      *string
	iconColor *string
	display   *string
	url       *url.Url
	reload    *bool
	vertical  bool
}

// New creates a new MultiStat.
func New() *MultiStat {
	return &MultiStat{
		items: make([]*stat.Stat, 0),
	}
}

// Add appends a stat item to the card. Reuses stat.Stat as the item type —
// its per-item fields (value, label, icon, color, prefix/suffix, trend) are
// exactly what an item needs.
func (m *MultiStat) Add(s *stat.Stat) *MultiStat {
	m.items = append(m.items, s)
	return m
}

// Title sets the card header title (translated).
func (m *MultiStat) Title(t string) *MultiStat {
	m.title = &t
	return m
}

// Icon sets the header icon.
func (m *MultiStat) Icon(icon string) *MultiStat {
	m.icon = &icon
	return m
}

// IconColor sets the header icon color.
func (m *MultiStat) IconColor(color string) *MultiStat {
	m.iconColor = &color
	return m
}

// WithDisplay sets the display/layout class.
func (m *MultiStat) WithDisplay(display string) *MultiStat {
	m.display = &display
	return m
}

// SetURL sets the AJAX data URL. When set, items are cleared and the frontend
// loads them dynamically; the header (title/icon) stays visible while loading.
func (m *MultiStat) SetURL(u *url.Url) *MultiStat {
	m.url = u
	m.items = nil
	return m
}

// WithReload enables a manual reload button when using AJAX mode.
func (m *MultiStat) WithReload(reload bool) *MultiStat {
	m.reload = &reload
	return m
}

// VerticalItems stacks each item's parts vertically (icon over value/label).
// Default is horizontal (icon beside value/label, like a single stat card).
func (m *MultiStat) VerticalItems() *MultiStat {
	m.vertical = true
	return m
}

// Print returns the JSON representation of the multi-stat component.
func (m *MultiStat) Print(ctx *core.UiContext) map[string]any {
	var data map[string]any

	if m.url != nil {
		data = map[string]any{
			"url":   m.url.PrintPrefix(),
			"items": nil,
		}
		if m.reload != nil {
			data["reload"] = *m.reload
		}
	} else {
		data = m.printData(ctx)
	}

	m.addHeader(ctx, data)

	result := map[string]any{
		"type": "multi-stat",
		"data": data,
	}
	if m.display != nil {
		result["display"] = *m.display
	}
	return result
}

// PrintData returns only the data portion (for use in data endpoints).
func (m *MultiStat) PrintData(ctx *core.UiContext) map[string]any {
	data := m.printData(ctx)
	m.addHeader(ctx, data)
	return data
}

// printData builds the items portion of the data map.
func (m *MultiStat) printData(ctx *core.UiContext) map[string]any {
	itemsData := make([]map[string]any, len(m.items))
	for i, s := range m.items {
		itemsData[i] = s.PrintData(ctx)
	}
	return map[string]any{
		"items": itemsData,
	}
}

// addHeader adds the optional title/icon/iconColor keys to the data map.
func (m *MultiStat) addHeader(ctx *core.UiContext, data map[string]any) {
	if m.title != nil {
		data["title"] = core.Translate(ctx, *m.title)
	}
	if m.icon != nil {
		data["icon"] = *m.icon
	}
	if m.iconColor != nil {
		data["iconColor"] = *m.iconColor
	}
	if m.vertical {
		data["verticalItems"] = true
	}
}

// DataResponse returns a DataResult wrapping the data in {"data": ...} envelope.
func (m *MultiStat) DataResponse(ctx *core.UiContext) response.DataResult {
	return response.NewJSONDataResult(m.PrintData(ctx))
}
