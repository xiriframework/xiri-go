// Package dialog provides dialog/modal components for the Angular frontend.
package dialog

import (
	"github.com/xiriframework/xiri-go/component/button"
	"github.com/xiriframework/xiri-go/component/core"
)

// DialogQuestionContent represents content for question-type dialogs (delete, warning)
type DialogQuestionContent struct {
	Icon     string `json:"icon"`
	Question string `json:"question"`
}

// Print returns the JSON representation of question content
func (d DialogQuestionContent) Print(translator core.TranslateFunc) map[string]any {
	return map[string]any{
		"icon":     d.Icon,
		"question": d.Question,
	}
}

// DialogWaitingContent represents content for waiting dialogs
type DialogWaitingContent struct {
	Icon string `json:"icon"`
	Text string `json:"text"`
}

// Print returns the JSON representation of waiting content
func (d DialogWaitingContent) Print(translator core.TranslateFunc) map[string]any {
	return map[string]any{
		"icon": d.Icon,
		"text": d.Text,
	}
}

// DialogContent is an interface for content that needs custom serialization
type DialogContent interface {
	Print(translator core.TranslateFunc) map[string]any
}

// Dialog represents any component that can be rendered as a dialog
type Dialog interface {
	Print(translator core.TranslateFunc) map[string]any
	WithExtra(extra map[string]any) Dialog
	WithOptions(options map[string]any) Dialog
	WithOption(key string, value any) Dialog
}

// dialogImpl represents a dialog/modal component
type dialogImpl struct {
	dialogType  core.DialogType
	header      string
	content     any
	buttons     []*button.Button
	extra       map[string]any
	options     map[string]any
	hookContent func(any)
}

// newDialog creates a new dialog component (package-private base constructor)
func newDialog(
	dialogType core.DialogType,
	header string,
	content any,
	buttons []*button.Button,
	extra map[string]any,
	options map[string]any,
) *dialogImpl {
	if options == nil {
		options = make(map[string]any)
	}
	return &dialogImpl{
		dialogType:  dialogType,
		header:      header,
		content:     content,
		buttons:     buttons,
		extra:       extra,
		options:     options,
		hookContent: nil,
	}
}

// NewDialog creates a new dialog component (public generic constructor)
//
// This is the most flexible constructor - use it when you need full control over
// dialog type, content, buttons, and options. For common dialog types, prefer the
// specialized constructors like NewDialogDelete, NewDialogForm, etc.
func NewDialog(
	dialogType core.DialogType,
	header string,
	content any,
	buttons []*button.Button,
	extra map[string]any,
	options map[string]any,
) Dialog {
	return newDialog(dialogType, header, content, buttons, extra, options)
}

// WithExtra sets the extra data map (optional)
//
// Extra data is passed to the frontend and can contain any additional
// information needed for the dialog (e.g., selected IDs for multi-operations).
func (d *dialogImpl) WithExtra(extra map[string]any) Dialog {
	d.extra = extra
	return d
}

// WithOptions sets the options map (optional)
//
// Options control dialog behavior (e.g., size, filter, URL).
// These are merged into the root of the dialog JSON output.
func (d *dialogImpl) WithOptions(options map[string]any) Dialog {
	if options != nil {
		d.options = options
	}
	return d
}

// WithOption sets a single option key-value pair
//
// Convenience method to set individual options without replacing the entire map.
func (d *dialogImpl) WithOption(key string, value any) Dialog {
	d.options[key] = value
	return d
}

// Print returns the JSON representation of the dialog
//
// The output structure matches Angular's XiriDialogSettings interface:
//   - header: Dialog title
//   - type: Dialog type (question, form, waiting, table)
//   - content: Dialog content (may be processed via DialogContent.Print())
//   - extra: Additional data passed to frontend
//   - buttons: Array of button configurations
//   - [options]: Any keys from options map are merged at root level
func (d *dialogImpl) Print(translator core.TranslateFunc) map[string]any {
	buttonData := make([]map[string]any, len(d.buttons))
	for i, btn := range d.buttons {
		buttonData[i] = btn.Print(translator)
	}

	content := d.content
	if d.hookContent != nil {
		d.hookContent(content)
	}

	var processedContent any
	if contentPrinter, ok := content.(DialogContent); ok {
		processedContent = contentPrinter.Print(translator)
	} else {
		processedContent = content
	}

	data := map[string]any{
		"header":  d.header,
		"type":    string(d.dialogType),
		"content": processedContent,
		"extra":   d.extra,
		"buttons": buttonData,
	}

	for key, value := range d.options {
		data[key] = value
	}

	return data
}
