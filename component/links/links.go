// Package links provides a link list card component for the Angular frontend.
package links

import (
	"github.com/xiriframework/xiri-go/component/button"
	"github.com/xiriframework/xiri-go/component/core"
)

// Links represents a card displaying a list of link buttons.
type Links struct {
	buttons         []*button.Button
	header          *string
	headerSub       *string
	headerIcon      *string
	headerIconColor *core.Color
	display         *string
}

// New creates a new Links component.
func New() *Links {
	return &Links{
		buttons: make([]*button.Button, 0),
	}
}

// Add adds a button to the link list.
func (l *Links) Add(b *button.Button) *Links {
	l.buttons = append(l.buttons, b)
	return l
}

// Header sets the card title.
func (l *Links) Header(h string) *Links {
	l.header = &h
	return l
}

// HeaderSub sets the card subtitle.
func (l *Links) HeaderSub(s string) *Links {
	l.headerSub = &s
	return l
}

// HeaderIcon sets the header icon and its color.
func (l *Links) HeaderIcon(icon string, color core.Color) *Links {
	l.headerIcon = &icon
	l.headerIconColor = &color
	return l
}

// WithDisplay sets the display/layout class.
func (l *Links) WithDisplay(d string) *Links {
	l.display = &d
	return l
}

// Print returns the JSON representation of the links component.
func (l *Links) Print(translator core.TranslateFunc) map[string]any {
	buttonData := make([]map[string]any, len(l.buttons))
	for i, btn := range l.buttons {
		buttonData[i] = btn.Print(translator)
	}

	data := map[string]any{
		"data": buttonData,
	}

	if l.header != nil {
		data["header"] = core.Translate(translator, *l.header)
	}
	if l.headerSub != nil {
		data["headerSub"] = core.Translate(translator, *l.headerSub)
	}
	if l.headerIcon != nil {
		data["headerIcon"] = *l.headerIcon
	}
	if l.headerIconColor != nil {
		data["headerIconColor"] = string(*l.headerIconColor)
	}

	result := map[string]any{
		"type": "links",
		"data": data,
	}

	if l.display != nil {
		result["display"] = *l.display
	}

	return result
}
