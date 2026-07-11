// Package bulletchart provides a bullet chart component for the Angular
// frontend — a compact alternative to a gauge for showing a value against a
// target. The `xiri-ng` side renders it via echarts (horizontal bar + target
// markline).
package bulletchart

import (
	"github.com/xiriframework/xiri-go/component/chart"
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
	"github.com/xiriframework/xiri-go/response"
)

// BulletChart shows a single value as a horizontal bar, optionally against a
// target marker.
type BulletChart struct {
	base *chart.BaseChart

	value  float64
	target *float64
	max    *float64
	label  string // optional label (e.g. a unit)
}

// New creates a new BulletChart bound to the given id.
func New(id string) *BulletChart {
	return &BulletChart{base: chart.New(id)}
}

// Title sets the chart title (translation key).
func (b *BulletChart) Title(t string) *BulletChart { b.base.SetTitle(t); return b }

// Value sets the current value displayed by the bar.
func (b *BulletChart) Value(v float64) *BulletChart { b.value = v; return b }

// Target sets the target marker (rendered as a markline).
func (b *BulletChart) Target(t float64) *BulletChart { b.target = &t; return b }

// Max sets the axis maximum. Default (frontend) is max(value, target) * 1.2.
func (b *BulletChart) Max(m float64) *BulletChart { b.max = &m; return b }

// Label sets an optional label (e.g. a unit).
func (b *BulletChart) Label(l string) *BulletChart { b.label = l; return b }

// Color sets the bar color.
func (b *BulletChart) Color(c core.Color) *BulletChart { b.base.SetColor(c); return b }

// Compact switches the chart into compact mode (thinner bar, no own card).
func (b *BulletChart) Compact() *BulletChart { b.base.SetCompact(); return b }

// WithDisplay sets the layout class.
func (b *BulletChart) WithDisplay(d string) *BulletChart { b.base.SetDisplay(d); return b }

// SetURL switches the chart to AJAX mode.
func (b *BulletChart) SetURL(u *url.Url) *BulletChart { b.base.SetURL(u); return b }

// WithReload enables periodic reload in AJAX mode.
func (b *BulletChart) WithReload(r bool) *BulletChart { b.base.SetReload(r); return b }

// Print returns the JSON envelope for the bullet chart.
func (b *BulletChart) Print(ctx *core.UiContext) map[string]any {
	if b.base.HasURL() {
		return b.base.Envelope("bulletchart", b.base.PrintAjaxData(), nil)
	}
	return b.base.Envelope("bulletchart", b.printData(ctx), nil)
}

// PrintData returns only the data map (for AJAX endpoints).
func (b *BulletChart) PrintData(ctx *core.UiContext) map[string]any {
	return b.printData(ctx)
}

// DataResponse wraps PrintData in a {"data": ...} envelope for DataResult.
func (b *BulletChart) DataResponse(ctx *core.UiContext) response.DataResult {
	return response.NewJSONDataResult(b.PrintData(ctx))
}

func (b *BulletChart) printData(ctx *core.UiContext) map[string]any {
	data := b.base.PrintBaseData(ctx)

	data["value"] = b.value
	if b.target != nil {
		data["target"] = *b.target
	}
	if b.max != nil {
		data["max"] = *b.max
	}
	if b.label != "" {
		data["label"] = b.label
	}
	return data
}
