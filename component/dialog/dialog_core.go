// Package dialog provides dialog/modal components for the Angular frontend.
package dialog

import (
	"github.com/xiriframework/xiri-go/component/button"
	"github.com/xiriframework/xiri-go/component/core"
)

// Dialog represents any component that can be rendered as a dialog
type Dialog interface {
	Print(ctx *core.UiContext) map[string]any
	WithExtra(extra map[string]any) Dialog
	WithOptions(options map[string]any) Dialog
	WithOption(key string, value any) Dialog
	WithData(payload map[string]any) Dialog
	SetButtons(buttons []*button.Button) Dialog
}

// dialogImpl represents a dialog/modal component
type dialogImpl struct {
	dialogType  core.DialogType
	header      string
	content     any
	buttons     []*button.Button
	extra       map[string]any
	options     map[string]any
	data        map[string]any
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
// Used for structural top-level fields (e.g. url, checkTime, size). For custom
// payload consumed by the Angular frontend, prefer WithData.
func (d *dialogImpl) WithOption(key string, value any) Dialog {
	d.options[key] = value
	return d
}

// WithData sets the custom payload rendered under the explicit top-level
// "data" key of the dialog JSON. Calling with a nil or empty map omits the
// "data" key entirely. If both WithData and WithOption("data", …) are set,
// WithData wins.
func (d *dialogImpl) WithData(payload map[string]any) Dialog {
	d.data = payload
	return d
}

// SetButtons replaces the dialog's buttons
func (d *dialogImpl) SetButtons(buttons []*button.Button) Dialog {
	d.buttons = buttons
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
func (d *dialogImpl) Print(ctx *core.UiContext) map[string]any {
	buttonData := make([]map[string]any, len(d.buttons))
	for i, btn := range d.buttons {
		buttonData[i] = btn.Print(ctx)
	}

	content := d.content
	if d.hookContent != nil {
		d.hookContent(content)
	}

	var processedContent any
	if contentPrinter, ok := content.(DialogContent); ok {
		processedContent = contentPrinter.Print(ctx)
	} else {
		processedContent = content
	}

	data := map[string]any{
		"header":  core.Translate(ctx, d.header),
		"type":    string(d.dialogType),
		"content": processedContent,
		"extra":   d.extra,
		"buttons": buttonData,
	}

	for key, value := range d.options {
		data[key] = value
	}

	// Add custom payload under the explicit "data" key (only if non-empty);
	// overrides any options["data"] set via legacy WithOption("data", …).
	if len(d.data) > 0 {
		data["data"] = d.data
	}

	return data
}
