package field

import (
	"fmt"

	"github.com/xiriframework/xiri-go/component/core"
)

// SelectField represents a dropdown/select field with predefined options
type SelectField struct {
	*BaseField
	Options  []SelectOption
	Subtype  string         // Select subtype (e.g., "select", "radio", "checkbox")
	Search   *bool          // Enable search filter (nil = auto based on option count, true/false = force)
	Multiple bool           // If true, field allows selecting multiple options (uses Values instead of Value)
	Value    int32          // Parsed and validated value (single-select mode)
	Values   []int32        // Parsed and validated values (multi-select mode)
}

// SelectOption represents a single option in a select field
type SelectOption struct {
	Value interface{} // The actual value
	Label string      // Translation key for display
}

// optionExists checks if the given int32 id matches one of the field's options.
func (f *SelectField) optionExists(id int32) bool {
	for _, opt := range f.Options {
		switch v := opt.Value.(type) {
		case int32:
			if v == id {
				return true
			}
		case int:
			if int32(v) == id {
				return true
			}
		case int64:
			if int32(v) == id {
				return true
			}
		case float64:
			if int32(v) == id {
				return true
			}
		}
	}
	return false
}

func (f *SelectField) Validate(value interface{}) error {
	if f.Multiple {
		if value == nil {
			if f.Required {
				return fmt.Errorf("select field %s is required", f.ID)
			}
			return nil
		}
		list, ok := value.(ModelListValue)
		if !ok {
			return fmt.Errorf("invalid select multi value type for %s", f.ID)
		}
		if f.Required && len(list) == 0 {
			return fmt.Errorf("select field %s is required", f.ID)
		}
		for _, id := range list {
			if !f.optionExists(id) {
				return fmt.Errorf("select field %s has invalid value %d", f.ID, id)
			}
		}
		return nil
	}

	if value == nil {
		if f.Required {
			return fmt.Errorf("select field %s is required", f.ID)
		}
		return nil
	}

	// Check if value matches one of the options
	for _, opt := range f.Options {
		if opt.Value == value {
			return nil
		}
	}

	return fmt.Errorf("select field %s has invalid value", f.ID)
}

func (f *SelectField) Parse(raw interface{}) (interface{}, error) {
	if f.Multiple {
		parsed, err := parseModelListValue(raw, f.GetDefault())
		if err != nil {
			return nil, fmt.Errorf("select field %s: %w", f.ID, err)
		}
		for _, id := range parsed {
			if !f.optionExists(id) {
				return nil, fmt.Errorf("select field %s has no matching option for value %v", f.ID, id)
			}
		}
		return parsed, nil
	}

	if raw == nil {
		return f.GetDefault(), nil
	}

	// Try to match against options
	for _, opt := range f.Options {
		switch optVal := opt.Value.(type) {
		case string:
			if str, ok := raw.(string); ok && str == optVal {
				return optVal, nil
			}
		case int:
			if num, ok := raw.(float64); ok && int(num) == optVal {
				return optVal, nil
			}
			if num, ok := raw.(int); ok && num == optVal {
				return optVal, nil
			}
		case int32:
			// Handle int32 options (common for database IDs)
			if num, ok := raw.(float64); ok && int32(num) == optVal {
				return optVal, nil
			}
			if num, ok := raw.(int32); ok && num == optVal {
				return optVal, nil
			}
			if num, ok := raw.(int); ok && int32(num) == optVal {
				return optVal, nil
			}
		case int64:
			// Handle int64 options
			if num, ok := raw.(float64); ok && int64(num) == optVal {
				return optVal, nil
			}
			if num, ok := raw.(int64); ok && num == optVal {
				return optVal, nil
			}
			if num, ok := raw.(int); ok && int64(num) == optVal {
				return optVal, nil
			}
		}
	}

	return nil, fmt.Errorf("select field %s has no matching option for value %v", f.ID, raw)
}

// BindValue parses, validates, and stores the value in the field
func (f *SelectField) BindValue(raw interface{}) error {
	parsed, err := f.Parse(raw)
	if err != nil {
		return fmt.Errorf("parsing field %s: %w", f.ID, err)
	}

	if err := f.Validate(parsed); err != nil {
		return fmt.Errorf("validating field %s: %w", f.ID, err)
	}

	if f.Multiple {
		if list, ok := parsed.(ModelListValue); ok {
			f.Values = []int32(list)
		} else {
			f.Values = []int32{}
		}
		return nil
	}

	// Convert parsed value to int32
	switch v := parsed.(type) {
	case int32:
		f.Value = v
	case int:
		f.Value = int32(v)
	case float64:
		f.Value = int32(v)
	case nil:
		// Use default value if nil
		if defaultVal, ok := f.GetDefault().(int32); ok {
			f.Value = defaultVal
		} else {
			f.Value = 0
		}
	default:
		return fmt.Errorf("field %s: cannot convert %T to int32", f.ID, parsed)
	}

	return nil
}

// ============================================================================
// Builder Functions
// ============================================================================

// NewSelectField creates a select/dropdown form field
func NewSelectField(id, name string, required bool, options []SelectOption) *SelectField {
	var defaultValue interface{}
	if len(options) > 0 {
		defaultValue = options[0].Value
	}

	return &SelectField{
		BaseField: &BaseField{
			ID:       id,
			Type:     FieldTypeSelect,
			Name:     name,
			Required: required,
			Default:  defaultValue,
			Form:     true,
		},
		Options: options,
	}
}

// ExportForFrontend exports the field for frontend rendering
func (f *SelectField) ExportForFrontend(ctx *core.UiContext, value interface{}) map[string]interface{} {
	if value == nil {
		value = f.GetDefault()
	}

	// In multi-select mode the frontend FormControl expects an array.
	if f.Multiple {
		switch v := value.(type) {
		case nil:
			value = ModelListValue{}
		case ModelListValue:
			// ok
		case []int32:
			value = ModelListValue(v)
		case int32:
			if v == 0 {
				value = ModelListValue{}
			} else {
				value = ModelListValue{v}
			}
		}
	}

	result := f.BaseField.GetBaseExport(ctx, value)

	// Add subtype if specified
	if f.Subtype != "" {
		result["subtype"] = f.Subtype
	}

	if f.Multiple {
		result["multiple"] = true
	}

	// Export options as array of {id, name} maps
	// MUST be "list" not "options" to match frontend
	options := make([]map[string]interface{}, len(f.Options))
	for i, opt := range f.Options {
		options[i] = map[string]interface{}{
			"id":   opt.Value,
			"name": opt.Label, // Label should be translation key, FormGroup will translate
		}
	}
	result["list"] = options

	// Set search based on option count (disable if < 20 options unless explicitly set)
	if f.Search != nil {
		result["search"] = *f.Search
	} else {
		result["search"] = len(f.Options) >= 20
	}

	return result
}

// ============================================================================
// Chainable Setter Methods
// ============================================================================

// SetClass sets the CSS class for frontend styling
func (f *SelectField) SetClass(class string) *SelectField {
	f.BaseField.SetClass(class)
	return f
}

// SetHint sets the tooltip/help text for the field
func (f *SelectField) SetHint(hint string) *SelectField {
	f.BaseField.SetHint(hint)
	return f
}

// SetStep sets the step indicator for multi-step forms
func (f *SelectField) SetStep(step int) *SelectField {
	f.BaseField.SetStep(step)
	return f
}

// SetDisabled sets whether the field is disabled
func (f *SelectField) SetDisabled(disabled bool) *SelectField {
	f.BaseField.SetDisabled(disabled)
	return f
}

// SetAccess sets the access control permissions
func (f *SelectField) SetAccess(access []string) *SelectField {
	f.BaseField.SetAccess(access)
	return f
}

// SetScenario sets which scenarios this field applies to
func (f *SelectField) SetScenario(scenario []string) *SelectField {
	f.BaseField.SetScenario(scenario)
	return f
}

// SetForm sets whether to show in form
func (f *SelectField) SetForm(form bool) *SelectField {
	f.BaseField.SetForm(form)
	return f
}

// SetMultiple toggles multi-select mode. When enabled, the field exports
// multiple=true to the frontend (rendered as mat-select [multiple]) and
// stores parsed values in Values instead of Value. Default selection
// becomes an empty list.
func (f *SelectField) SetMultiple(multiple bool) *SelectField {
	f.Multiple = multiple
	if multiple {
		f.Default = ModelListValue{}
	}
	return f
}
