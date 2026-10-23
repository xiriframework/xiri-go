// Package scatterchart provides a scatter chart for the Angular frontend: points on two value axes,
// grouped into series. The `xiri-ng` side renders via echarts (optional peerDependency).
package scatterchart

import (
	"github.com/xiriframework/xiri-go/component/chart"
	"github.com/xiriframework/xiri-go/component/core"
)

// Point is one dot. Name shows in the tooltip and may be "".
type Point struct {
	X, Y float64
	Name string
}

// P is a short constructor for a Point.
func P(x, y float64, name string) Point { return Point{X: x, Y: y, Name: name} }

type series struct {
	name   string
	color  core.Color
	points []Point
}

// ScatterChart shows how two values relate.
type ScatterChart struct {
	base           *chart.BaseChart
	series         []series
	xLabel, yLabel string
}

// New creates a new ScatterChart bound to the given id.
func New(id string) *ScatterChart {
	return &ScatterChart{base: chart.New(id)}
}

// Title sets the chart title (translation key).
func (sc *ScatterChart) Title(t string) *ScatterChart { sc.base.SetTitle(t); return sc }

// Axes sets the axis titles (translation keys); "" leaves an axis untitled.
func (sc *ScatterChart) Axes(x, y string) *ScatterChart { sc.xLabel, sc.yLabel = x, y; return sc }

// Series adds a group of points. color may be "".
func (sc *ScatterChart) Series(name string, color core.Color, points ...Point) *ScatterChart {
	sc.series = append(sc.series, series{name: name, color: color, points: points})
	return sc
}

// Compact switches the chart into compact mode.
func (sc *ScatterChart) Compact() *ScatterChart { sc.base.SetCompact(); return sc }

// WithDisplay sets the layout class.
func (sc *ScatterChart) WithDisplay(d string) *ScatterChart { sc.base.SetDisplay(d); return sc }

// Print returns the JSON envelope for the scatter chart.
func (sc *ScatterChart) Print(ctx *core.UiContext) map[string]any {
	data := sc.base.PrintBaseData(ctx)
	if sc.xLabel != "" {
		data["xLabel"] = core.Translate(ctx, sc.xLabel)
	}
	if sc.yLabel != "" {
		data["yLabel"] = core.Translate(ctx, sc.yLabel)
	}
	ser := make([]map[string]any, len(sc.series))
	for i, s := range sc.series {
		pts := make([]map[string]any, len(s.points))
		for j, p := range s.points {
			pts[j] = map[string]any{"x": p.X, "y": p.Y}
			if p.Name != "" {
				pts[j]["name"] = p.Name
			}
		}
		ser[i] = map[string]any{"name": s.name, "points": pts}
		if s.color != "" {
			ser[i]["color"] = string(s.color)
		}
	}
	data["series"] = ser
	return sc.base.Envelope("scatterchart", data, nil)
}
