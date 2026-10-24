// Package radarchart provides a radar (spider) chart for the Angular frontend: several series
// compared across the same indicators. The `xiri-ng` side renders via echarts (optional peerDependency).
package radarchart

import (
	"github.com/xiriframework/xiri-go/component/chart"
	"github.com/xiriframework/xiri-go/component/core"
)

type indicator struct {
	name string
	max  float64
}

type series struct {
	name   string
	color  core.Color
	values []float64
}

// RadarChart compares series across indicators (the axes of the net).
type RadarChart struct {
	base       *chart.BaseChart
	indicators []indicator
	series     []series
	area       bool
}

// New creates a new RadarChart bound to the given id.
func New(id string) *RadarChart {
	return &RadarChart{base: chart.New(id)}
}

// Title sets the chart title (translation key).
func (rc *RadarChart) Title(t string) *RadarChart { rc.base.SetTitle(t); return rc }

// Indicator adds an axis. max 0 lets the frontend scale the axis itself.
func (rc *RadarChart) Indicator(name string, max float64) *RadarChart {
	rc.indicators = append(rc.indicators, indicator{name: name, max: max})
	return rc
}

// Series adds a series with one value per indicator, in indicator order. color may be "".
func (rc *RadarChart) Series(name string, color core.Color, values ...float64) *RadarChart {
	rc.series = append(rc.series, series{name: name, color: color, values: append([]float64{}, values...)})
	return rc
}

// Area fills the area of each series.
func (rc *RadarChart) Area() *RadarChart { rc.area = true; return rc }

// Compact switches the chart into compact mode.
func (rc *RadarChart) Compact() *RadarChart { rc.base.SetCompact(); return rc }

// WithDisplay sets the layout class.
func (rc *RadarChart) WithDisplay(d string) *RadarChart { rc.base.SetDisplay(d); return rc }

// Print returns the JSON envelope for the radar chart.
func (rc *RadarChart) Print(ctx *core.UiContext) map[string]any {
	data := rc.base.PrintBaseData(ctx)
	if rc.area {
		data["area"] = true
	}
	inds := make([]map[string]any, len(rc.indicators))
	for i, in := range rc.indicators {
		inds[i] = map[string]any{"name": in.name}
		if in.max != 0 {
			inds[i]["max"] = in.max
		}
	}
	data["indicators"] = inds
	ser := make([]map[string]any, len(rc.series))
	for i, s := range rc.series {
		ser[i] = map[string]any{"name": s.name, "values": s.values}
		if s.color != "" {
			ser[i]["color"] = string(s.color)
		}
	}
	data["series"] = ser
	return rc.base.Envelope("radarchart", data, nil)
}
