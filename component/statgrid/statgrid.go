// Package statgrid provides a grid container for multiple stat/KPI cards.
package statgrid

import (
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/stat"
	"github.com/xiriframework/xiri-go/component/url"
	"github.com/xiriframework/xiri-go/response"
)

// StatGrid represents a grid layout for stat/KPI cards.
type StatGrid struct {
	stats   []*stat.Stat
	columns *int
	title   *string
	display *string
	url     *url.Url
	reload  *bool
}

// New creates a new StatGrid.
func New() *StatGrid {
	return &StatGrid{
		stats: make([]*stat.Stat, 0),
	}
}

// Add adds a stat card to the grid.
func (g *StatGrid) Add(s *stat.Stat) *StatGrid {
	g.stats = append(g.stats, s)
	return g
}

// Columns sets the number of grid columns (default: 4).
func (g *StatGrid) Columns(c int) *StatGrid {
	g.columns = &c
	return g
}

// Title sets the grid title.
func (g *StatGrid) Title(t string) *StatGrid {
	g.title = &t
	return g
}

// WithDisplay sets the display/layout class.
func (g *StatGrid) WithDisplay(display string) *StatGrid {
	g.display = &display
	return g
}

// SetURL sets the AJAX data URL. When set, stats are cleared and the frontend loads data dynamically.
func (g *StatGrid) SetURL(url *url.Url) *StatGrid {
	g.url = url
	g.stats = nil
	return g
}

// WithReload enables periodic reload of the stat grid data when using AJAX mode.
func (g *StatGrid) WithReload(reload bool) *StatGrid {
	g.reload = &reload
	return g
}

// Print returns the JSON representation of the stat grid component.
func (g *StatGrid) Print(translator core.TranslateFunc) map[string]any {
	var data map[string]any

	if g.url != nil {
		data = map[string]any{
			"url":   g.url.PrintPrefix(),
			"stats": nil,
		}
		if g.reload != nil {
			data["reload"] = *g.reload
		}
	} else {
		data = g.printData(translator)
	}

	if g.columns != nil {
		data["columns"] = *g.columns
	}
	if g.title != nil {
		data["title"] = core.Translate(translator, *g.title)
	}

	result := map[string]any{
		"type": "stat-grid",
		"data": data,
	}

	if g.display != nil {
		result["display"] = *g.display
	}

	return result
}

// PrintData returns only the data portion of the stat grid (for use in data endpoints).
func (g *StatGrid) PrintData(translator core.TranslateFunc) map[string]any {
	data := g.printData(translator)

	if g.columns != nil {
		data["columns"] = *g.columns
	}
	if g.title != nil {
		data["title"] = core.Translate(translator, *g.title)
	}

	return data
}

// DataResponse returns a DataResult wrapping the stat grid data in {"data": ...} envelope.
func (g *StatGrid) DataResponse(translator core.TranslateFunc) response.DataResult {
	return response.NewJSONDataResult(g.PrintData(translator))
}

// printData builds the data map used by both Print and PrintData.
func (g *StatGrid) printData(translator core.TranslateFunc) map[string]any {
	statsData := make([]map[string]any, len(g.stats))
	for i, s := range g.stats {
		statsData[i] = s.PrintData(translator)
	}
	return map[string]any{
		"stats": statsData,
	}
}
