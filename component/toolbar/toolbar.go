// Package toolbar provides a toolbar component with title, icon, search, and action buttons.
package toolbar

import (
	"github.com/xiriframework/xiri-go/component/button"
	"github.com/xiriframework/xiri-go/component/core"
)

// SearchConfig holds configuration for the toolbar search field.
type SearchConfig struct {
	Placeholder *string
}

// Toolbar represents a toolbar with optional title, icon, search, and buttons.
type Toolbar struct {
	title   *string
	icon    *string
	search  *SearchConfig
	buttons *button.ButtonLine
	display *string
}

// New creates a new Toolbar.
func New() *Toolbar {
	return &Toolbar{}
}

// Title sets the toolbar title.
func (t *Toolbar) Title(title string) *Toolbar {
	t.title = &title
	return t
}

// Icon sets the toolbar icon.
func (t *Toolbar) Icon(icon string) *Toolbar {
	t.icon = &icon
	return t
}

// Search enables the search field with an optional placeholder.
func (t *Toolbar) Search(placeholder *string) *Toolbar {
	t.search = &SearchConfig{Placeholder: placeholder}
	return t
}

// Buttons sets the action buttons.
func (t *Toolbar) Buttons(b *button.ButtonLine) *Toolbar {
	t.buttons = b
	return t
}

// WithDisplay sets the display/layout class.
func (t *Toolbar) WithDisplay(display string) *Toolbar {
	t.display = &display
	return t
}

// Print returns the JSON representation of the toolbar component.
func (t *Toolbar) Print(translator core.TranslateFunc) map[string]any {
	data := make(map[string]any)

	if t.title != nil {
		data["title"] = core.Translate(translator, *t.title)
	}
	if t.icon != nil {
		data["icon"] = *t.icon
	}
	if t.search != nil {
		if t.search.Placeholder != nil {
			data["search"] = map[string]any{
				"placeholder": core.Translate(translator, *t.search.Placeholder),
			}
		} else {
			data["search"] = true
		}
	}
	if t.buttons != nil {
		data["buttons"] = t.buttons.PrintData(translator)
	}

	result := map[string]any{
		"type": "toolbar",
		"data": data,
	}

	if t.display != nil {
		result["display"] = *t.display
	}

	return result
}
