package field

import (
	"fmt"
	"strconv"

	"github.com/xiriframework/xiri-go/uicontext"
)

// IntField represents an integer field with optional min/max bounds
type IntField struct {
	*BaseField
	Min        *int
	Max        *int
	Pattern    string // Validation pattern (regex)
	Subtype    string // Number subtype: "int", "pint" (positive int), "float", "bigint", "real"
	TextPrefix string // Prefix text (e.g., "$")
	TextSuffix string // Suffix text (e.g., "km")
	IconPrefix string // Prefix icon name
	IconSuffix string // Suffix icon name
	Value      *int32 // Parsed and validated value (type-safe access)
}

func (f *IntField) Validate(value interface{}) error {
	if value == nil {
		if f.Required {
			return fmt.Errorf("int field %s is required", f.ID)
		}
		return nil
	}

	var num int
	switch v := value.(type) {
	case int:
		num = v
	case int32:
		num = int(v)
	case int64:
		num = int(v)
	default:
		return fmt.Errorf("invalid int value type for %s", f.ID)
	}

	if f.Min != nil && num < *f.Min {
		return fmt.Errorf("int field %s must be >= %d", f.ID, *f.Min)
	}

	if f.Max != nil && num > *f.Max {
		return fmt.Errorf("int field %s must be <= %d", f.ID, *f.Max)
	}

	return nil
}

func (f *IntField) Parse(raw interface{}) (interface{}, error) {
	if raw == nil {
		return f.GetDefault(), nil
	}

	switch v := raw.(type) {
	case int:
		return v, nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid int value: %s", v)
		}
		return parsed, nil
	default:
		return nil, fmt.Errorf("cannot parse int from %T", raw)
	}
}

// BindValue parses, validates, and stores the value in the field
func (f *IntField) BindValue(raw interface{}) error {
	parsed, err := f.Parse(raw)
	if err != nil {
		return fmt.Errorf("parsing field %s: %w", f.ID, err)
	}

	if err := f.Validate(parsed); err != nil {
		return fmt.Errorf("validating field %s: %w", f.ID, err)
	}

	if parsed != nil {
		if num, ok := parsed.(int); ok {
			val := int32(num)
			f.Value = &val
		}
	} else {
		f.Value = nil
	}

	return nil
}

// ============================================================================
// Builder Functions
// ============================================================================

// NewIntField creates an integer form field
func NewIntField(id, name string, required bool, currentValue int32) *IntField {
	return &IntField{
		BaseField: &BaseField{
			ID:       id,
			Type:     FieldTypeInt,
			Name:     name,
			Required: required,
			Default:  currentValue,
			Form:     true,
		},
	}
}

// NewIntFieldWithBounds creates an integer field with min/max bounds
func NewIntFieldWithBounds(id, name string, required bool, currentValue int32, min, max int) *IntField {
	return &IntField{
		BaseField: &BaseField{
			ID:       id,
			Type:     FieldTypeInt,
			Name:     name,
			Required: required,
			Default:  currentValue,
			Form:     true,
		},
		Min: &min,
		Max: &max,
	}
}

// NewNumberField creates a number form field (alias for IntField)
func NewNumberField(id, name string, required bool, defaultValue float64) *IntField {
	return &IntField{
		BaseField: &BaseField{
			ID:       id,
			Type:     FieldTypeInt,
			Name:     name,
			Required: required,
			Default:  int(defaultValue),
			Form:     true,
		},
	}
}

// ExportForFrontend exports the field for frontend rendering
func (f *IntField) ExportForFrontend(ctx *uicontext.UiContext, value interface{}) map[string]interface{} {
	if value == nil {
		value = f.GetDefault()
	}
	result := f.BaseField.GetBaseExport(ctx, value)

	// Set subtype to "number" for display
	result["subtype"] = "number"

	// Add min/max if specified
	if f.Min != nil {
		result["min"] = *f.Min
	} else if f.Subtype == "pint" {
		// Positive int has implicit min of 0
		result["min"] = 0
	}
	if f.Max != nil {
		result["max"] = *f.Max
	}

	// Add prefix/suffix text and icons
	if f.TextPrefix != "" {
		result["textPrefix"] = f.TextPrefix
	}
	if f.TextSuffix != "" {
		result["textSuffix"] = f.TextSuffix
	}
	if f.IconPrefix != "" {
		result["iconPrefix"] = f.IconPrefix
	}
	if f.IconSuffix != "" {
		result["iconSuffix"] = f.IconSuffix
	}

	// Add locale for number formatting
	if ctx != nil {
		result["locale"] = ctx.Locale.GetLocaleString()
	} else {
		result["locale"] = "de-DE" // Default to German locale
	}

	return result
}

// ============================================================================
// Chainable Setter Methods
// ============================================================================

// SetClass sets the CSS class for frontend styling
func (f *IntField) SetClass(class string) *IntField {
	f.BaseField.SetClass(class)
	return f
}

// SetHint sets the tooltip/help text for the field
func (f *IntField) SetHint(hint string) *IntField {
	f.BaseField.SetHint(hint)
	return f
}

// SetStep sets the step indicator for multi-step forms
func (f *IntField) SetStep(step int) *IntField {
	f.BaseField.SetStep(step)
	return f
}

// SetDisabled sets whether the field is disabled
func (f *IntField) SetDisabled(disabled bool) *IntField {
	f.BaseField.SetDisabled(disabled)
	return f
}

// SetAccess sets the access control permissions
func (f *IntField) SetAccess(access []string) *IntField {
	f.BaseField.SetAccess(access)
	return f
}

// SetScenario sets which scenarios this field applies to
func (f *IntField) SetScenario(scenario []string) *IntField {
	f.BaseField.SetScenario(scenario)
	return f
}

// SetForm sets whether to show in form
func (f *IntField) SetForm(form bool) *IntField {
	f.BaseField.SetForm(form)
	return f
}
