package field

import (
	"github.com/xiriframework/xiri-go/uicontext"
)

// HeaderField represents a header/section label (non-interactive)
type HeaderField struct {
	*BaseField
	Content     string // Info content to display
	Collapsible bool   // Whether the section can be collapsed
	Collapsed   bool   // Initial collapsed state
}

func (f *HeaderField) Validate(value interface{}) error {
	// Headers are display-only, no validation needed
	return nil
}

func (f *HeaderField) Parse(raw interface{}) (interface{}, error) {
	// Headers are display-only, no parsing needed
	return nil, nil
}

// ============================================================================
// Builder Functions
// ============================================================================

// NewHeaderField creates a header (display-only) form field
func NewHeaderField(id, content string) *HeaderField {
	return &HeaderField{
		BaseField: &BaseField{
			ID:       id,
			Type:     FieldTypeHeader,
			Name:     "",
			Required: false,
			Default:  nil,
			Form:     true,
		},
		Content: content,
	}
}

// ExportForFrontend exports the field for frontend rendering
func (f *HeaderField) ExportForFrontend(ctx *uicontext.UiContext, value interface{}) map[string]interface{} {
	result := f.BaseField.GetBaseExport(ctx, f.Content)
	if f.Collapsible {
		result["collapsible"] = true
		result["collapsed"] = f.Collapsed
	}
	return result
}

// ============================================================================
// Chainable Setter Methods
// ============================================================================

// SetClass sets the CSS class for frontend styling
func (f *HeaderField) SetClass(class string) *HeaderField {
	f.BaseField.SetClass(class)
	return f
}

// SetHint sets the tooltip/help text for the field
func (f *HeaderField) SetHint(hint string) *HeaderField {
	f.BaseField.SetHint(hint)
	return f
}

// SetStep sets the step indicator for multi-step forms
func (f *HeaderField) SetStep(step int) *HeaderField {
	f.BaseField.SetStep(step)
	return f
}

// SetDisabled sets whether the field is disabled
func (f *HeaderField) SetDisabled(disabled bool) *HeaderField {
	f.BaseField.SetDisabled(disabled)
	return f
}

// SetAccess sets the access control permissions
func (f *HeaderField) SetAccess(access []string) *HeaderField {
	f.BaseField.SetAccess(access)
	return f
}

// SetScenario sets which scenarios this field applies to
func (f *HeaderField) SetScenario(scenario []string) *HeaderField {
	f.BaseField.SetScenario(scenario)
	return f
}

// SetForm sets whether to show in form
func (f *HeaderField) SetForm(form bool) *HeaderField {
	f.BaseField.SetForm(form)
	return f
}

// SetCollapsible makes this header a collapsible section header
func (f *HeaderField) SetCollapsible(collapsible bool) *HeaderField {
	f.Collapsible = collapsible
	return f
}

// SetCollapsed sets the initial collapsed state (only effective with SetCollapsible(true))
func (f *HeaderField) SetCollapsed(collapsed bool) *HeaderField {
	f.Collapsed = collapsed
	return f
}
