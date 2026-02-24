// Package form provides form UI components for the Angular frontend.
package form

import (
	"github.com/xiriframework/xiri-go/component/button"
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
)

// Form represents a form component with fields and action buttons
type Form struct {
	fields     []map[string]any
	url        *url.Url
	buttons    []*button.Button
	header     *string
	display    *string
	hookFields func([]map[string]any)
	translator core.TranslateFunc
}

// NewForm creates a new form component with default Back and Save buttons
// PHP equivalent: XiriForm::gen($xfields, $url, $obj, $header, $buttons, $data, $display)
func NewForm(
	fields []map[string]any,
	u *url.Url,
	header *string,
	buttons []*button.Button,
	display *string,
	translator core.TranslateFunc,
) *Form {
	// If no buttons provided, create default Back + Save buttons
	if buttons == nil {
		defaultTabIndex := -1
		buttons = []*button.Button{
			button.NewBackButton(
				core.Translate(translator, "Back"),
				core.ColorPrimary,
				core.ButtonTypeStroked,
				"",
				false,
				&defaultTabIndex,
				false,
				nil,
			),
			button.NewFormButton(
				core.Translate(translator, "Save"),
				u,
				core.ColorPrimary,
				core.ButtonTypeRaised,
				"",
				false,
				nil,
				true,
				nil,
			),
		}
	}

	return &Form{
		fields:     fields,
		url:        u,
		buttons:    buttons,
		header:     header,
		display:    display,
		hookFields: nil,
		translator: translator,
	}
}

// HookFields sets a hook function to modify fields before printing
// PHP equivalent: XiriForm->hookFields($hook)
func (f *Form) HookFields(hook func([]map[string]any)) *Form {
	f.hookFields = hook
	return f
}

// WithHeader sets the form header (optional)
// Returns the Form for method chaining
func (f *Form) WithHeader(header string) *Form {
	f.header = &header
	return f
}

// WithDisplay sets the display/layout class (optional)
// Returns the Form for method chaining
func (f *Form) WithDisplay(display string) *Form {
	f.display = &display
	return f
}

// Print returns the JSON representation of the form
// PHP equivalent: XiriForm->print()
func (f *Form) Print(translator core.TranslateFunc) map[string]any {
	// Use the form's translator if provided, otherwise use the parameter
	trans := translator
	if trans == nil && f.translator != nil {
		trans = f.translator
	}

	// Prepare buttons
	buttonData := make([]map[string]any, len(f.buttons))
	for i, btn := range f.buttons {
		buttonData[i] = btn.Print(trans)
	}

	// Apply hook if set
	fields := f.fields
	if f.hookFields != nil {
		f.hookFields(fields)
	}

	return map[string]any{
		"type":    "form",
		"display": f.display,
		"data": map[string]any{
			"header":  f.header,
			"url":     f.url.PrintPrefix(),
			"fields":  fields,
			"buttons": buttonData,
		},
	}
}
