package field

import (
	"github.com/xiriframework/xiri-go/uicontext"
)

// HtmlField represents HTML content display (non-interactive)
type HtmlField struct {
	*BaseField
	Content string // HTML content to display
}

func (f *HtmlField) Validate(value interface{}) error {
	// HTML fields are display-only, no validation needed
	return nil
}

func (f *HtmlField) Parse(raw interface{}) (interface{}, error) {
	// HTML fields are display-only, no parsing needed
	return nil, nil
}

// ============================================================================
// Builder Functions
// ============================================================================

// NewHtmlField creates an HTML (display-only) form field
func NewHtmlField(id, content string) *HtmlField {
	return &HtmlField{
		BaseField: &BaseField{
			ID:       id,
			Type:     FieldTypeHtml,
			Name:     "",
			Required: false,
			Default:  nil,
			Form:     true,
		},
		Content: content,
	}
}

// ExportForFrontend exports the field for frontend rendering
func (f *HtmlField) ExportForFrontend(ctx *uicontext.UiContext, value interface{}) map[string]interface{} {
	return f.BaseField.GetBaseExport(ctx, f.Content)
}

// ============================================================================
// Chainable Setter Methods
// ============================================================================

// SetClass sets the CSS class for frontend styling
func (f *HtmlField) SetClass(class string) *HtmlField {
	f.BaseField.SetClass(class)
	return f
}

// SetHint sets the tooltip/help text for the field
func (f *HtmlField) SetHint(hint string) *HtmlField {
	f.BaseField.SetHint(hint)
	return f
}

// SetStep sets the step indicator for multi-step forms
func (f *HtmlField) SetStep(step int) *HtmlField {
	f.BaseField.SetStep(step)
	return f
}

// SetDisabled sets whether the field is disabled
func (f *HtmlField) SetDisabled(disabled bool) *HtmlField {
	f.BaseField.SetDisabled(disabled)
	return f
}

// SetAccess sets the access control permissions
func (f *HtmlField) SetAccess(access []string) *HtmlField {
	f.BaseField.SetAccess(access)
	return f
}

// SetScenario sets which scenarios this field applies to
func (f *HtmlField) SetScenario(scenario []string) *HtmlField {
	f.BaseField.SetScenario(scenario)
	return f
}

// SetForm sets whether to show in form
func (f *HtmlField) SetForm(form bool) *HtmlField {
	f.BaseField.SetForm(form)
	return f
}
