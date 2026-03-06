// Package pageheader provides a page header component with title, subtitle, icon, and action buttons.
package pageheader

import (
	"github.com/xiriframework/xiri-go/component/button"
	"github.com/xiriframework/xiri-go/component/core"
)

// PageHeader represents a page header with title, optional subtitle, icon, and buttons.
type PageHeader struct {
	title     string
	subtitle  *string
	icon      *string
	iconColor *core.Color
	buttons   *button.ButtonLine
	display   *string
}

// New creates a new PageHeader with the given title.
func New(title string) *PageHeader {
	return &PageHeader{
		title: title,
	}
}

// Subtitle sets the subtitle text.
func (p *PageHeader) Subtitle(s string) *PageHeader {
	p.subtitle = &s
	return p
}

// Icon sets the icon and its color.
func (p *PageHeader) Icon(icon string, color core.Color) *PageHeader {
	p.icon = &icon
	p.iconColor = &color
	return p
}

// Buttons sets the action buttons.
func (p *PageHeader) Buttons(b *button.ButtonLine) *PageHeader {
	p.buttons = b
	return p
}

// WithDisplay sets the display/layout class.
func (p *PageHeader) WithDisplay(display string) *PageHeader {
	p.display = &display
	return p
}

// Print returns the JSON representation of the page header component.
func (p *PageHeader) Print(ctx *core.UiContext) map[string]any {
	data := map[string]any{
		"title": core.Translate(ctx, p.title),
	}

	if p.subtitle != nil {
		data["subtitle"] = core.Translate(ctx, *p.subtitle)
	}
	if p.icon != nil {
		data["icon"] = *p.icon
	}
	if p.iconColor != nil {
		data["iconColor"] = string(*p.iconColor)
	}
	if p.buttons != nil {
		data["buttons"] = p.buttons.PrintData(ctx)
	}

	result := map[string]any{
		"type": "page-header",
		"data": data,
	}

	if p.display != nil {
		result["display"] = *p.display
	}

	return result
}
