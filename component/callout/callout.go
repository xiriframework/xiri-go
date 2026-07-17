// Package callout provides a persistent, non-modal inline notice component for
// the Angular frontend. Unlike alert (modal dialog) or the transient snackbar,
// a callout stays in the page flow to carry context next to the content it
// refers to.
package callout

import (
	"github.com/xiriframework/xiri-go/component/button"
	"github.com/xiriframework/xiri-go/component/core"
)

// Callout represents a persistent, non-modal inline notice.
type Callout struct {
	tone        string
	text        string
	title       *string
	icon        *string
	actions     []*button.Button
	dismissible bool
	compact     bool
}

// New creates a new Callout with the given tone ("info", "success", "warning",
// "error") and text. Both are required.
func New(tone, text string) *Callout {
	return &Callout{
		tone: tone,
		text: text,
	}
}

// Title sets the optional bold heading shown above the text.
func (c *Callout) Title(title string) *Callout {
	c.title = &title
	return c
}

// Icon sets the optional leading icon.
func (c *Callout) Icon(icon string) *Callout {
	c.icon = &icon
	return c
}

// AddAction appends one or more action buttons. Buttons reuse the existing
// button component and its action semantics.
func (c *Callout) AddAction(btn ...*button.Button) *Callout {
	c.actions = append(c.actions, btn...)
	return c
}

// Dismissible lets the user close the callout locally (client-side only).
func (c *Callout) Dismissible() *Callout {
	c.dismissible = true
	return c
}

// Compact switches the callout into a denser variant.
func (c *Callout) Compact() *Callout {
	c.compact = true
	return c
}

// Print returns the JSON representation of the callout component.
func (c *Callout) Print(ctx *core.UiContext) map[string]any {
	data := map[string]any{
		"tone": c.tone,
		"text": core.Translate(ctx, c.text),
	}

	if c.title != nil {
		data["title"] = core.Translate(ctx, *c.title)
	}
	if c.icon != nil {
		data["icon"] = *c.icon
	}
	if len(c.actions) > 0 {
		actions := make([]map[string]any, len(c.actions))
		for i, btn := range c.actions {
			actions[i] = btn.Print(ctx)
		}
		data["actions"] = actions
	}
	if c.dismissible {
		data["dismissible"] = true
	}
	if c.compact {
		data["compact"] = true
	}

	return map[string]any{
		"type": "callout",
		"data": data,
	}
}
