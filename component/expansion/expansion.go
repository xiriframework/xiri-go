// Package expansion provides accordion/expansion panel components for the Angular frontend.
package expansion

import "github.com/xiriframework/xiri-go/component/core"

// Panel represents an individual panel within an Expansion component
// Angular: XiriExpansionPanelSettings interface
type Panel struct {
	title       string
	description *string
	icon        *string
	disabled    *bool
	expanded    *bool
	lazy        *bool
	unload      *bool
	data        []core.Component
}

// NewPanel creates a new Panel with the given title
func NewPanel(title string) *Panel {
	return &Panel{
		title: title,
		data:  []core.Component{},
	}
}

// WithDescription sets the description text shown next to the title
func (p *Panel) WithDescription(description string) *Panel {
	p.description = &description
	return p
}

// WithIcon sets the Material icon for the panel header
func (p *Panel) WithIcon(icon string) *Panel {
	p.icon = &icon
	return p
}

// WithDisabled sets whether the panel is disabled
func (p *Panel) WithDisabled(disabled bool) *Panel {
	p.disabled = &disabled
	return p
}

// WithExpanded sets whether the panel is initially expanded
func (p *Panel) WithExpanded(expanded bool) *Panel {
	p.expanded = &expanded
	return p
}

// WithLazy sets whether the panel content should be lazy-loaded
// This can override the global lazy setting on the Expansion container
func (p *Panel) WithLazy(lazy bool) *Panel {
	p.lazy = &lazy
	return p
}

// WithUnload sets whether the panel content should be destroyed when closed
// When true, content is re-rendered every time the panel is opened
// This can override the global unload setting on the Expansion container
func (p *Panel) WithUnload(unload bool) *Panel {
	p.unload = &unload
	return p
}

// AddContent adds a component to be rendered in the panel content
func (p *Panel) AddContent(component core.Component) *Panel {
	p.data = append(p.data, component)
	return p
}

// Print serializes the Panel to a map for JSON output
func (p *Panel) Print(translator core.TranslateFunc) map[string]any {
	result := map[string]any{
		"title": core.Translate(translator, p.title),
	}

	if p.description != nil {
		result["description"] = core.Translate(translator, *p.description)
	}
	if p.icon != nil {
		result["icon"] = *p.icon
	}
	if p.disabled != nil {
		result["disabled"] = *p.disabled
	}
	if p.expanded != nil {
		result["expanded"] = *p.expanded
	}
	if p.lazy != nil {
		result["lazy"] = *p.lazy
	}
	if p.unload != nil {
		result["unload"] = *p.unload
	}

	// Convert content components to data array
	if len(p.data) > 0 {
		data := make([]map[string]any, len(p.data))
		for i, comp := range p.data {
			data[i] = comp.Print(translator)
		}
		result["data"] = data
	} else {
		result["data"] = []map[string]any{}
	}

	return result
}

// Expansion represents an accordion/expansion panel container component
// Angular: XiriExpansionSettings interface
type Expansion struct {
	panels         []*Panel
	multi          *bool
	displayMode    *core.ExpansionDisplayMode
	togglePosition *core.ExpansionTogglePosition
	hideToggle     *bool
	lazy           *bool
	unload         *bool
	display        *string
}

// NewExpansion creates a new Expansion container
func NewExpansion() *Expansion {
	return &Expansion{
		panels: []*Panel{},
	}
}

// AddPanel adds a panel to the container
func (e *Expansion) AddPanel(panel *Panel) *Expansion {
	e.panels = append(e.panels, panel)
	return e
}

// WithMulti sets whether multiple panels can be expanded at the same time
func (e *Expansion) WithMulti(multi bool) *Expansion {
	e.multi = &multi
	return e
}

// WithDisplayMode sets the display mode of the accordion
func (e *Expansion) WithDisplayMode(mode core.ExpansionDisplayMode) *Expansion {
	e.displayMode = &mode
	return e
}

// WithTogglePosition sets where the expand/collapse toggle icon is positioned
func (e *Expansion) WithTogglePosition(position core.ExpansionTogglePosition) *Expansion {
	e.togglePosition = &position
	return e
}

// WithHideToggle sets whether the expand/collapse toggle icon is hidden
func (e *Expansion) WithHideToggle(hide bool) *Expansion {
	e.hideToggle = &hide
	return e
}

// WithLazy sets the global lazy-loading setting for all panels
// Individual panels can override this with Panel.WithLazy()
func (e *Expansion) WithLazy(lazy bool) *Expansion {
	e.lazy = &lazy
	return e
}

// WithUnload sets the global unload setting for all panels
// When true, panel content is destroyed when closed and re-rendered when opened
// Individual panels can override this with Panel.WithUnload()
func (e *Expansion) WithUnload(unload bool) *Expansion {
	e.unload = &unload
	return e
}

// WithDisplay sets the CSS display class
func (e *Expansion) WithDisplay(display string) *Expansion {
	e.display = &display
	return e
}

// Print serializes the Expansion to a map for JSON output
func (e *Expansion) Print(translator core.TranslateFunc) map[string]any {
	data := map[string]any{}

	// Convert panels
	panelsData := make([]map[string]any, len(e.panels))
	for i, panel := range e.panels {
		panelsData[i] = panel.Print(translator)
	}
	data["panels"] = panelsData

	// Optional fields - only include if set
	if e.multi != nil {
		data["multi"] = *e.multi
	}
	if e.displayMode != nil {
		data["displayMode"] = string(*e.displayMode)
	}
	if e.togglePosition != nil {
		data["togglePosition"] = string(*e.togglePosition)
	}
	if e.hideToggle != nil {
		data["hideToggle"] = *e.hideToggle
	}
	if e.lazy != nil {
		data["lazy"] = *e.lazy
	}
	if e.unload != nil {
		data["unload"] = *e.unload
	}

	return map[string]any{
		"type":    "expansion",
		"display": e.display,
		"data":    data,
	}
}
