// Package tabs provides tabbed container components for the Angular frontend.
package tabs

import "github.com/xiriframework/xiri-go/component/core"

// Tab represents an individual tab within a Tabs component
// Angular: XiriTabSettings interface
type Tab struct {
	label    string
	icon     *string
	disabled *bool
	lazy     *bool
	unload   *bool
	data     []core.Component
}

// NewTab creates a new Tab with the given label
func NewTab(label string) *Tab {
	return &Tab{
		label: label,
		data:  []core.Component{},
	}
}

// WithIcon sets the Material icon for the tab header
func (t *Tab) WithIcon(icon string) *Tab {
	t.icon = &icon
	return t
}

// WithDisabled sets whether the tab is disabled
func (t *Tab) WithDisabled(disabled bool) *Tab {
	t.disabled = &disabled
	return t
}

// WithLazy sets whether the tab content should be lazy-loaded
// This can override the global lazy setting on the Tabs container
func (t *Tab) WithLazy(lazy bool) *Tab {
	t.lazy = &lazy
	return t
}

// WithUnload sets whether the tab content should be destroyed when switching away
// When true, content is re-rendered every time the tab is selected
// This can override the global unload setting on the Tabs container
func (t *Tab) WithUnload(unload bool) *Tab {
	t.unload = &unload
	return t
}

// AddContent adds a component to be rendered in the tab content
func (t *Tab) AddContent(component core.Component) *Tab {
	t.data = append(t.data, component)
	return t
}

// Print serializes the Tab to a map for JSON output
func (t *Tab) Print(translator core.TranslateFunc) map[string]any {
	result := map[string]any{
		"label": core.Translate(translator, t.label),
	}

	if t.icon != nil {
		result["icon"] = *t.icon
	}
	if t.disabled != nil {
		result["disabled"] = *t.disabled
	}
	if t.lazy != nil {
		result["lazy"] = *t.lazy
	}
	if t.unload != nil {
		result["unload"] = *t.unload
	}

	// Convert content components to data array
	if len(t.data) > 0 {
		data := make([]map[string]any, len(t.data))
		for i, comp := range t.data {
			data[i] = comp.Print(translator)
		}
		result["data"] = data
	} else {
		result["data"] = []map[string]any{}
	}

	return result
}

// Tabs represents a tabbed container component
// Angular: XiriTabsSettings interface
type Tabs struct {
	tabs              []*Tab
	selectedIndex     *int
	dynamicHeight     *bool
	animationDuration *string
	lazy              *bool
	unload            *bool
	headerPosition    *core.TabHeaderPosition
	alignTabs         *core.TabAlignment
	stretchTabs       *bool
	display           *string
}

// NewTabs creates a new Tabs container
func NewTabs() *Tabs {
	return &Tabs{
		tabs: []*Tab{},
	}
}

// AddTab adds a tab to the container
func (t *Tabs) AddTab(tab *Tab) *Tabs {
	t.tabs = append(t.tabs, tab)
	return t
}

// WithSelectedIndex sets the initially selected tab index (0-based)
func (t *Tabs) WithSelectedIndex(index int) *Tabs {
	t.selectedIndex = &index
	return t
}

// WithDynamicHeight sets whether the container height adjusts to active tab
func (t *Tabs) WithDynamicHeight(dynamic bool) *Tabs {
	t.dynamicHeight = &dynamic
	return t
}

// WithAnimationDuration sets the CSS animation duration for tab transitions
func (t *Tabs) WithAnimationDuration(duration string) *Tabs {
	t.animationDuration = &duration
	return t
}

// WithLazy sets the global lazy-loading setting for all tabs
// Individual tabs can override this with Tab.WithLazy()
func (t *Tabs) WithLazy(lazy bool) *Tabs {
	t.lazy = &lazy
	return t
}

// WithUnload sets the global unload setting for all tabs
// When true, tab content is destroyed when switching away and re-rendered when selected
// Individual tabs can override this with Tab.WithUnload()
func (t *Tabs) WithUnload(unload bool) *Tabs {
	t.unload = &unload
	return t
}

// WithHeaderPosition sets where tab headers appear relative to content
func (t *Tabs) WithHeaderPosition(position core.TabHeaderPosition) *Tabs {
	t.headerPosition = &position
	return t
}

// WithAlignTabs sets the horizontal alignment of tab headers
func (t *Tabs) WithAlignTabs(align core.TabAlignment) *Tabs {
	t.alignTabs = &align
	return t
}

// WithStretchTabs sets whether tabs stretch to fill available width
func (t *Tabs) WithStretchTabs(stretch bool) *Tabs {
	t.stretchTabs = &stretch
	return t
}

// WithDisplay sets the CSS display class
func (t *Tabs) WithDisplay(display string) *Tabs {
	t.display = &display
	return t
}

// Print serializes the Tabs to a map for JSON output
func (t *Tabs) Print(translator core.TranslateFunc) map[string]any {
	data := map[string]any{}

	// Convert tabs
	tabsData := make([]map[string]any, len(t.tabs))
	for i, tab := range t.tabs {
		tabsData[i] = tab.Print(translator)
	}
	data["tabs"] = tabsData

	// Optional fields - only include if set
	if t.selectedIndex != nil {
		data["selectedIndex"] = *t.selectedIndex
	}
	if t.dynamicHeight != nil {
		data["dynamicHeight"] = *t.dynamicHeight
	}
	if t.animationDuration != nil {
		data["animationDuration"] = *t.animationDuration
	}
	if t.lazy != nil {
		data["lazy"] = *t.lazy
	}
	if t.unload != nil {
		data["unload"] = *t.unload
	}
	if t.headerPosition != nil {
		data["headerPosition"] = string(*t.headerPosition)
	}
	if t.alignTabs != nil {
		data["alignTabs"] = string(*t.alignTabs)
	}
	if t.stretchTabs != nil {
		data["stretchTabs"] = *t.stretchTabs
	}

	return map[string]any{
		"type":    "tabs",
		"display": t.display,
		"data":    data,
	}
}
