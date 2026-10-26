// Package markdown provides a Markdown text block for the Angular frontend. xiri-ng renders it with
// `marked`; raw HTML in the text is shown as text (only a bare <br> stays a line break) and the result goes
// through Angular's sanitizer, so the text may come from users or an LLM. For trusted raw HTML use
// info.NewInfoText(...).WithHtml().
package markdown

import "github.com/xiriframework/xiri-go/component/core"

// Markdown is a block of Markdown text (GitHub flavour: tables, lists, code).
type Markdown struct {
	text    string
	display *string
}

// New creates a Markdown block. The text is not translated.
func New(text string) *Markdown { return &Markdown{text: text} }

// WithDisplay sets the layout class.
func (m *Markdown) WithDisplay(d string) *Markdown { m.display = &d; return m }

// Print returns the JSON envelope for the Markdown block.
func (m *Markdown) Print(_ *core.UiContext) map[string]any {
	out := map[string]any{"type": "markdown", "data": map[string]any{"text": m.text}}
	if m.display != nil {
		out["display"] = *m.display
	}
	return out
}
