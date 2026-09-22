package core

import "slices"

// componentTypes ist der Vertrag über alle "type"-Werte, die xiri-go im Komponenten-JSON ausgibt.
// xiri-ng führt dieselbe Liste als XIRI_COMPONENT_TYPES (plus reine Frontend-Typen wie cardlink);
// component/core/types_test.go prüft beide Richtungen. Sortiert halten.
var componentTypes = []string{
	"barchart", "bulletchart", "buttonline", "calendar", "callout", "card", "container",
	"description-list", "divider", "empty-state", "expansion", "form", "gantt", "gaugechart",
	"header", "heatmap", "html", "imagetext", "infopoint", "infotext", "linechart", "links", "list",
	"multi-stat", "multiprogress", "page-header", "piechart", "progress", "query", "sankey", "section",
	"spacer", "stat", "stat-grid", "status", "stepper", "table", "tabs", "tachotime",
	"timeline", "toolbar", "tree",
}

// ComponentTypes liefert eine Kopie der sortierten Liste aller "type"-Werte, die xiri-go ausgibt.
func ComponentTypes() []string {
	return slices.Clone(componentTypes)
}
