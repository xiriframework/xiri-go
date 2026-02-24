// Package section provides a collapsible section component with title, icon, buttons, and nested components.
package section

import (
	"github.com/xiriframework/xiri-go/component/button"
	"github.com/xiriframework/xiri-go/component/core"
)

// Section represents a page section with optional header, collapsibility, and nested components.
type Section struct {
	title       *string
	subtitle    *string
	icon        *string
	iconColor   *core.Color
	collapsible bool
	collapsed   bool
	buttons     *button.ButtonLine
	components  []core.Component
	display     *string
}

// New creates a new Section.
func New() *Section {
	return &Section{
		components: make([]core.Component, 0),
	}
}

// Title sets the section title.
func (s *Section) Title(t string) *Section {
	s.title = &t
	return s
}

// Subtitle sets the section subtitle.
func (s *Section) Subtitle(t string) *Section {
	s.subtitle = &t
	return s
}

// Icon sets the section icon and its color.
func (s *Section) Icon(icon string, color core.Color) *Section {
	s.icon = &icon
	s.iconColor = &color
	return s
}

// Collapsible makes the section collapsible with initial collapsed state.
func (s *Section) Collapsible(collapsed bool) *Section {
	s.collapsible = true
	s.collapsed = collapsed
	return s
}

// Buttons sets the action buttons for the section header.
func (s *Section) Buttons(b *button.ButtonLine) *Section {
	s.buttons = b
	return s
}

// Add adds a child component to the section.
func (s *Section) Add(component core.Component) *Section {
	s.components = append(s.components, component)
	return s
}

// WithDisplay sets the display/layout class.
func (s *Section) WithDisplay(display string) *Section {
	s.display = &display
	return s
}

// Print returns the JSON representation of the section component.
func (s *Section) Print(translator core.TranslateFunc) map[string]any {
	data := make(map[string]any)

	if s.title != nil {
		data["title"] = core.Translate(translator, *s.title)
	}
	if s.subtitle != nil {
		data["subtitle"] = core.Translate(translator, *s.subtitle)
	}
	if s.icon != nil {
		data["icon"] = *s.icon
	}
	if s.iconColor != nil {
		data["iconColor"] = string(*s.iconColor)
	}
	if s.collapsible {
		data["collapsible"] = true
		data["collapsed"] = s.collapsed
	}
	if s.buttons != nil {
		data["buttons"] = s.buttons.PrintData(translator)
	}

	if len(s.components) > 0 {
		componentData := make([]map[string]any, len(s.components))
		for i, comp := range s.components {
			componentData[i] = comp.Print(translator)
		}
		data["components"] = componentData
	}

	result := map[string]any{
		"type": "section",
		"data": data,
	}

	if s.display != nil {
		result["display"] = *s.display
	}

	return result
}
