package field

import (
	"github.com/xiriframework/xiri-go/uicontext"
)

// InfoField represents informational text display (non-interactive)
type InfoField struct {
	*BaseField
	Content string // Info content to display
}

func (f *InfoField) Validate(value interface{}) error {
	// Info fields are display-only, no validation needed
	return nil
}

func (f *InfoField) Parse(raw interface{}) (interface{}, error) {
	// Info fields are display-only, no parsing needed
	return nil, nil
}

// ============================================================================
// Builder Functions
// ============================================================================

// NewInfoField creates an info (display-only) form field
func NewInfoField(id, content string) *InfoField {
	return &InfoField{
		BaseField: &BaseField{
			ID:       id,
			Type:     FieldTypeInfo,
			Name:     "",
			Required: false,
			Default:  nil,
			Form:     true,
		},
		Content: content,
	}
}

// ExportForFrontend exports the field for frontend rendering
func (f *InfoField) ExportForFrontend(ctx *uicontext.UiContext, value interface{}) map[string]interface{} {
	// Info fields always use their Content as the value
	return f.BaseField.GetBaseExport(ctx, f.Content)
}

// ============================================================================
// Chainable Setter Methods
// ============================================================================

// SetClass sets the CSS class for frontend styling
func (f *InfoField) SetClass(class string) *InfoField {
	f.BaseField.SetClass(class)
	return f
}

// SetHint sets the tooltip/help text for the field
func (f *InfoField) SetHint(hint string) *InfoField {
	f.BaseField.SetHint(hint)
	return f
}

// SetStep sets the step indicator for multi-step forms
func (f *InfoField) SetStep(step int) *InfoField {
	f.BaseField.SetStep(step)
	return f
}

// SetDisabled sets whether the field is disabled
func (f *InfoField) SetDisabled(disabled bool) *InfoField {
	f.BaseField.SetDisabled(disabled)
	return f
}

// SetAccess sets the access control permissions
func (f *InfoField) SetAccess(access []string) *InfoField {
	f.BaseField.SetAccess(access)
	return f
}

// SetScenario sets which scenarios this field applies to
func (f *InfoField) SetScenario(scenario []string) *InfoField {
	f.BaseField.SetScenario(scenario)
	return f
}

// SetForm sets whether to show in form
func (f *InfoField) SetForm(form bool) *InfoField {
	f.BaseField.SetForm(form)
	return f
}
