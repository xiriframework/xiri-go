package progress

import (
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
	"github.com/xiriframework/xiri-go/response"
)

// Progress represents a single determinate progress bar showing "current of total"
// task progress — honest process tracking (e.g. backfills, backtest runs). For a
// share-of-sum breakdown use MultiProgress instead.
type Progress struct {
	current       int
	total         int
	label         string
	color         *core.Color
	indeterminate bool
	display       *string
	url           *url.Url
	reload        *bool
}

// NewProgress creates a single progress bar showing current of total with a label.
func NewProgress(current, total int, label string) *Progress {
	return &Progress{current: current, total: total, label: label}
}

// Color sets the bar color.
func (p *Progress) Color(color core.Color) *Progress {
	p.color = &color
	return p
}

// Indeterminate switches the bar into indeterminate mode (unknown progress).
func (p *Progress) Indeterminate() *Progress {
	p.indeterminate = true
	return p
}

// WithDisplay sets the display/layout class.
func (p *Progress) WithDisplay(display string) *Progress {
	p.display = &display
	return p
}

// SetURL sets the AJAX data URL. When set, the frontend loads progress data dynamically.
func (p *Progress) SetURL(u *url.Url) *Progress {
	p.url = u
	return p
}

// WithReload enables periodic reload of the progress data when using AJAX mode.
func (p *Progress) WithReload(reload bool) *Progress {
	p.reload = &reload
	return p
}

// Print returns the JSON representation of the progress bar.
func (p *Progress) Print(ctx *core.UiContext) map[string]any {
	var data map[string]any
	if p.url != nil {
		data = map[string]any{"url": p.url.PrintPrefix()}
		if p.reload != nil {
			data["reload"] = *p.reload
		}
	} else {
		data = p.printData(ctx)
	}

	result := map[string]any{
		"type": "progress",
		"data": data,
	}
	if p.display != nil {
		result["display"] = *p.display
	}
	return result
}

// PrintData returns only the data portion of the progress bar (for data endpoints).
func (p *Progress) PrintData(ctx *core.UiContext) map[string]any {
	return p.printData(ctx)
}

// DataResponse wraps the progress data in a {"data": ...} envelope.
func (p *Progress) DataResponse(ctx *core.UiContext) response.DataResult {
	return response.NewJSONDataResult(p.PrintData(ctx))
}

// printData builds the data map used by both Print and PrintData.
func (p *Progress) printData(ctx *core.UiContext) map[string]any {
	var value float64
	if p.total != 0 {
		value = float64(p.current) / float64(p.total) * 100
	}

	data := map[string]any{
		"label":   core.Translate(ctx, p.label),
		"current": p.current,
		"total":   p.total,
		"value":   value,
	}
	if p.color != nil {
		data["color"] = string(*p.color)
	}
	if p.indeterminate {
		data["indeterminate"] = true
	}
	return data
}
