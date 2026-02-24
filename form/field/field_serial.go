package field

import (
	"github.com/xiriframework/xiri-go/uicontext"
)

import "fmt"

// SerialField represents an auto-incrementing serial/ID form field
// Typically used for auto-generated IDs, read-only
type SerialField struct {
	*BaseField
	ReadOnly bool // Typically true for serial fields
}

func (f *SerialField) Validate(value interface{}) error {
	// Serial fields are typically auto-generated, validation not needed
	return nil
}

func (f *SerialField) Parse(raw interface{}) (interface{}, error) {
	if raw == nil {
		return f.GetDefault(), nil
	}

	// Parse to int64 (serial ID)
	switch v := raw.(type) {
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case float64:
		return int64(v), nil
	default:
		return nil, fmt.Errorf("cannot parse serial from %T", raw)
	}
}

// ============================================================================
// Builder Functions
// ============================================================================

// NewSerialField creates a serial/ID form field
func NewSerialField(id, name string) *SerialField {
	return &SerialField{
		BaseField: &BaseField{
			ID:       id,
			Type:     FieldTypeSerial,
			Name:     name,
			Required: false,
			Default:  nil,
			Form:     false,
		},
		ReadOnly: true,
	}
}

// ExportForFrontend exports the field for frontend rendering
func (f *SerialField) ExportForFrontend(ctx *uicontext.UiContext, value interface{}) map[string]interface{} {
	if value == nil {
		value = f.GetDefault()
	}
	return f.BaseField.GetBaseExport(ctx, value)
}

// ============================================================================
// Chainable Setter Methods
// ============================================================================

// SetClass sets the CSS class for frontend styling
func (f *SerialField) SetClass(class string) *SerialField {
	f.BaseField.SetClass(class)
	return f
}

// SetHint sets the tooltip/help text for the field
func (f *SerialField) SetHint(hint string) *SerialField {
	f.BaseField.SetHint(hint)
	return f
}

// SetStep sets the step indicator for multi-step forms
func (f *SerialField) SetStep(step int) *SerialField {
	f.BaseField.SetStep(step)
	return f
}

// SetDisabled sets whether the field is disabled
func (f *SerialField) SetDisabled(disabled bool) *SerialField {
	f.BaseField.SetDisabled(disabled)
	return f
}

// SetAccess sets the access control permissions
func (f *SerialField) SetAccess(access []string) *SerialField {
	f.BaseField.SetAccess(access)
	return f
}

// SetScenario sets which scenarios this field applies to
func (f *SerialField) SetScenario(scenario []string) *SerialField {
	f.BaseField.SetScenario(scenario)
	return f
}

// SetForm sets whether to show in form
func (f *SerialField) SetForm(form bool) *SerialField {
	f.BaseField.SetForm(form)
	return f
}
