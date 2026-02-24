package field

import (
	"fmt"

	"github.com/xiriframework/xiri-go/uicontext"
)

// ArrayField represents an array/list form field
type ArrayField struct {
	*BaseField
	ItemType    string        // Type of items in the array (e.g., "string", "int")
	MinItems    *int          // Minimum number of items
	MaxItems    *int          // Maximum number of items
	AllowEmpty  bool          // If true, empty array is allowed
	UniqueItems bool          // If true, all items must be unique
	Value       []interface{} // Parsed and validated value
}

func (f *ArrayField) Validate(value interface{}) error {
	if value == nil {
		if f.Required && !f.AllowEmpty {
			return fmt.Errorf("array field %s is required", f.ID)
		}
		return nil
	}

	arr, ok := value.([]interface{})
	if !ok {
		return fmt.Errorf("invalid array value type for %s", f.ID)
	}

	if !f.AllowEmpty && len(arr) == 0 {
		return fmt.Errorf("array %s cannot be empty", f.ID)
	}

	if f.MinItems != nil && len(arr) < *f.MinItems {
		return fmt.Errorf("array %s must have at least %d items", f.ID, *f.MinItems)
	}

	if f.MaxItems != nil && len(arr) > *f.MaxItems {
		return fmt.Errorf("array %s must have at most %d items", f.ID, *f.MaxItems)
	}

	return nil
}

func (f *ArrayField) Parse(raw interface{}) (interface{}, error) {
	if raw == nil {
		if def := f.GetDefault(); def != nil {
			return def, nil
		}
		return []interface{}{}, nil
	}

	arr, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("array field expects array")
	}

	return arr, nil
}

// BindValue parses, validates, and stores the value in the field
// This enables type-safe access via field.Value after form binding
func (f *ArrayField) BindValue(raw interface{}) error {
	parsed, err := f.Parse(raw)
	if err != nil {
		return fmt.Errorf("parsing field %s: %w", f.ID, err)
	}

	if err := f.Validate(parsed); err != nil {
		return fmt.Errorf("validating field %s: %w", f.ID, err)
	}

	if parsed != nil {
		if arr, ok := parsed.([]interface{}); ok {
			f.Value = arr
		}
	} else {
		f.Value = nil
	}

	return nil
}

// ============================================================================
// Builder Functions
// ============================================================================

// NewArrayField creates an array form field
func NewArrayField(id, name string, required bool, itemType string, defaultValue []interface{}) *ArrayField {
	return &ArrayField{
		BaseField: &BaseField{
			ID:       id,
			Type:     FieldTypeArray,
			Name:     name,
			Required: required,
			Default:  defaultValue,
			Form:     true,
		},
		ItemType:   itemType,
		AllowEmpty: !required,
	}
}

// ExportForFrontend exports the field for frontend rendering
func (f *ArrayField) ExportForFrontend(ctx *uicontext.UiContext, value interface{}) map[string]interface{} {
	if value == nil {
		value = f.GetDefault()
	}
	result := f.BaseField.GetBaseExport(ctx, value)

	// Set type to array
	result["type"] = "array"

	// Add min/max items if specified
	if f.MinItems != nil {
		result["min"] = *f.MinItems
	}
	if f.MaxItems != nil {
		result["max"] = *f.MaxItems
	}

	return result
}

// ============================================================================
// Chainable Setter Methods
// ============================================================================

// SetClass sets the CSS class for frontend styling
func (f *ArrayField) SetClass(class string) *ArrayField {
	f.BaseField.SetClass(class)
	return f
}

// SetHint sets the tooltip/help text for the field
func (f *ArrayField) SetHint(hint string) *ArrayField {
	f.BaseField.SetHint(hint)
	return f
}

// SetStep sets the step indicator for multi-step forms
func (f *ArrayField) SetStep(step int) *ArrayField {
	f.BaseField.SetStep(step)
	return f
}

// SetDisabled sets whether the field is disabled
func (f *ArrayField) SetDisabled(disabled bool) *ArrayField {
	f.BaseField.SetDisabled(disabled)
	return f
}

// SetAccess sets the access control permissions
func (f *ArrayField) SetAccess(access []string) *ArrayField {
	f.BaseField.SetAccess(access)
	return f
}

// SetScenario sets which scenarios this field applies to
func (f *ArrayField) SetScenario(scenario []string) *ArrayField {
	f.BaseField.SetScenario(scenario)
	return f
}

// SetForm sets whether to show in form
func (f *ArrayField) SetForm(form bool) *ArrayField {
	f.BaseField.SetForm(form)
	return f
}
