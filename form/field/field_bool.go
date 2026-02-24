package field

import (
	"fmt"

	"github.com/xiriframework/xiri-go/uicontext"
)

// BoolField represents a boolean form field
type BoolField struct {
	*BaseField
	Value *bool // Parsed and validated value (type-safe access)
}

func (f *BoolField) Validate(value interface{}) error {
	if value == nil {
		if f.Required {
			return fmt.Errorf("bool field %s is required", f.ID)
		}
		return nil
	}

	_, ok := value.(bool)
	if !ok {
		return fmt.Errorf("invalid bool value type for %s", f.ID)
	}

	return nil
}

func (f *BoolField) Parse(raw interface{}) (interface{}, error) {
	if raw == nil {
		return f.GetDefault(), nil
	}

	switch v := raw.(type) {
	case bool:
		return v, nil
	case string:
		return v == "true" || v == "1" || v == "yes", nil
	case int:
		return v != 0, nil
	case float64:
		return v != 0, nil
	default:
		return false, fmt.Errorf("cannot parse bool from %T", raw)
	}
}

// BindValue parses, validates, and stores the value in the field
func (f *BoolField) BindValue(raw interface{}) error {
	parsed, err := f.Parse(raw)
	if err != nil {
		return fmt.Errorf("parsing field %s: %w", f.ID, err)
	}

	if err := f.Validate(parsed); err != nil {
		return fmt.Errorf("validating field %s: %w", f.ID, err)
	}

	if parsed != nil {
		if b, ok := parsed.(bool); ok {
			f.Value = &b
		}
	} else {
		f.Value = nil
	}

	return nil
}

// ============================================================================
// Builder Functions
// ============================================================================

// ExportForFrontend exports the field for frontend rendering
func (f *BoolField) ExportForFrontend(ctx *uicontext.UiContext, value interface{}) map[string]interface{} {
	// Use value if provided, otherwise use default
	if value == nil {
		value = f.GetDefault()
	}
	return f.BaseField.GetBaseExport(ctx, value)
}

// ============================================================================
// Builder Functions
// ============================================================================

// NewBoolField creates a boolean form field
func NewBoolField(id, name string, required bool, currentValue bool) *BoolField {
	return &BoolField{
		BaseField: &BaseField{
			ID:       id,
			Type:     FieldTypeBool,
			Name:     name,
			Required: required,
			Default:  currentValue,
			Form:     true,
		},
	}
}

// ============================================================================
// Chainable Setter Methods
// ============================================================================

// SetClass sets the CSS class for frontend styling
func (f *BoolField) SetClass(class string) *BoolField {
	f.BaseField.SetClass(class)
	return f
}

// SetHint sets the tooltip/help text for the field
func (f *BoolField) SetHint(hint string) *BoolField {
	f.BaseField.SetHint(hint)
	return f
}

// SetStep sets the step indicator for multi-step forms
func (f *BoolField) SetStep(step int) *BoolField {
	f.BaseField.SetStep(step)
	return f
}

// SetDisabled sets whether the field is disabled
func (f *BoolField) SetDisabled(disabled bool) *BoolField {
	f.BaseField.SetDisabled(disabled)
	return f
}

// SetAccess sets the access control permissions
func (f *BoolField) SetAccess(access []string) *BoolField {
	f.BaseField.SetAccess(access)
	return f
}

// SetScenario sets which scenarios this field applies to
func (f *BoolField) SetScenario(scenario []string) *BoolField {
	f.BaseField.SetScenario(scenario)
	return f
}

// SetForm sets whether to show in form
func (f *BoolField) SetForm(form bool) *BoolField {
	f.BaseField.SetForm(form)
	return f
}
