package field

import (
	"github.com/xiriframework/xiri-go/uicontext"
)

import "fmt"

// JsonField represents a JSON object form field
type JsonField struct {
	*BaseField
	Schema interface{}            // Optional JSON schema for validation
	Value  map[string]interface{} // Parsed and validated value
}

func (f *JsonField) Validate(value interface{}) error {
	if value == nil {
		if f.Required {
			return fmt.Errorf("json field %s is required", f.ID)
		}
		return nil
	}

	// JSON validation - accept map or string
	switch value.(type) {
	case map[string]interface{}, []interface{}, string:
		return nil
	default:
		return fmt.Errorf("invalid json value type for %s", f.ID)
	}
}

func (f *JsonField) Parse(raw interface{}) (interface{}, error) {
	if raw == nil {
		return f.GetDefault(), nil
	}

	// Accept JSON string or already-parsed object
	return raw, nil
}

// BindValue parses, validates, and stores the value in the field
// This enables type-safe access via field.Value after form binding
func (f *JsonField) BindValue(raw interface{}) error {
	parsed, err := f.Parse(raw)
	if err != nil {
		return fmt.Errorf("parsing field %s: %w", f.ID, err)
	}

	if err := f.Validate(parsed); err != nil {
		return fmt.Errorf("validating field %s: %w", f.ID, err)
	}

	if parsed != nil {
		// Try to convert to map[string]interface{}
		if jsonMap, ok := parsed.(map[string]interface{}); ok {
			f.Value = jsonMap
		} else {
			// If it's a string, keep it as-is (frontend will parse)
			// For type safety, we store an empty map
			f.Value = make(map[string]interface{})
		}
	} else {
		f.Value = nil
	}

	return nil
}

// ============================================================================
// Builder Functions
// ============================================================================

// NewJsonField creates a JSON form field
func NewJsonField(id, name string, required bool, defaultValue map[string]interface{}) *JsonField {
	return &JsonField{
		BaseField: &BaseField{
			ID:       id,
			Type:     FieldTypeJson,
			Name:     name,
			Required: required,
			Default:  defaultValue,
			Form:     true,
		},
	}
}

// ExportForFrontend exports the field for frontend rendering
func (f *JsonField) ExportForFrontend(ctx *uicontext.UiContext, value interface{}) map[string]interface{} {
	if value == nil {
		value = f.GetDefault()
	}
	return f.BaseField.GetBaseExport(ctx, value)
}

// ============================================================================
// Chainable Setter Methods
// ============================================================================

// SetClass sets the CSS class for frontend styling
func (f *JsonField) SetClass(class string) *JsonField {
	f.BaseField.SetClass(class)
	return f
}

// SetHint sets the tooltip/help text for the field
func (f *JsonField) SetHint(hint string) *JsonField {
	f.BaseField.SetHint(hint)
	return f
}

// SetStep sets the step indicator for multi-step forms
func (f *JsonField) SetStep(step int) *JsonField {
	f.BaseField.SetStep(step)
	return f
}

// SetDisabled sets whether the field is disabled
func (f *JsonField) SetDisabled(disabled bool) *JsonField {
	f.BaseField.SetDisabled(disabled)
	return f
}

// SetAccess sets the access control permissions
func (f *JsonField) SetAccess(access []string) *JsonField {
	f.BaseField.SetAccess(access)
	return f
}

// SetScenario sets which scenarios this field applies to
func (f *JsonField) SetScenario(scenario []string) *JsonField {
	f.BaseField.SetScenario(scenario)
	return f
}

// SetForm sets whether to show in form
func (f *JsonField) SetForm(form bool) *JsonField {
	f.BaseField.SetForm(form)
	return f
}
