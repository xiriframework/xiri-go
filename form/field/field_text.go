package field

import (
	"fmt"

	"github.com/xiriframework/xiri-go/uicontext"
)

// TextField represents a free-text field
type TextField struct {
	*BaseField
	MinLength  int
	MaxLength  int
	Pattern    string  // Validation pattern (regex)
	Subtype    string  // Text subtype: "text", "textarea", "html", "email", "url", "tel", "password"
	TextPrefix string  // Prefix text
	TextSuffix string  // Suffix text
	IconPrefix string  // Prefix icon name
	IconSuffix string  // Suffix icon name
	Trim       bool    // Whether to trim whitespace (default: true)
	Value      *string // Parsed and validated value (type-safe access)
}

func (f *TextField) Validate(value interface{}) error {
	if value == nil {
		if f.Required {
			return fmt.Errorf("text field %s is required", f.ID)
		}
		return nil
	}

	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid text value type for %s", f.ID)
	}

	if f.MinLength > 0 && len(str) < f.MinLength {
		return fmt.Errorf("text field %s must be at least %d characters", f.ID, f.MinLength)
	}

	if f.MaxLength > 0 && len(str) > f.MaxLength {
		return fmt.Errorf("text field %s must be at most %d characters", f.ID, f.MaxLength)
	}

	return nil
}

func (f *TextField) Parse(raw interface{}) (interface{}, error) {
	if raw == nil {
		return f.GetDefault(), nil
	}

	if str, ok := raw.(string); ok {
		return str, nil
	}

	return fmt.Sprintf("%v", raw), nil
}

// BindValue parses, validates, and stores the value in the field
func (f *TextField) BindValue(raw interface{}) error {
	parsed, err := f.Parse(raw)
	if err != nil {
		return fmt.Errorf("parsing field %s: %w", f.ID, err)
	}

	if err := f.Validate(parsed); err != nil {
		return fmt.Errorf("validating field %s: %w", f.ID, err)
	}

	if parsed != nil {
		if str, ok := parsed.(string); ok {
			f.Value = &str
		}
	} else {
		f.Value = nil
	}

	return nil
}

// ============================================================================
// Builder Functions
// ============================================================================

// NewTextField creates a text form field
func NewTextField(id, name string, required bool, currentValue string) *TextField {
	return &TextField{
		BaseField: &BaseField{
			ID:       id,
			Type:     FieldTypeText,
			Name:     name,
			Required: required,
			Default:  currentValue,
			Form:     true,
		},
	}
}

// NewTextFieldWithLength creates a text field with length constraints
func NewTextFieldWithLength(id, name string, required bool, currentValue string, minLen, maxLen int) *TextField {
	return &TextField{
		BaseField: &BaseField{
			ID:       id,
			Type:     FieldTypeText,
			Name:     name,
			Required: required,
			Default:  currentValue,
			Form:     true,
		},
		MinLength: minLen,
		MaxLength: maxLen,
	}
}

// ExportForFrontend exports the field for frontend rendering
func (f *TextField) ExportForFrontend(ctx *uicontext.UiContext, value interface{}) map[string]interface{} {
	if value == nil {
		value = f.GetDefault()
	}
	result := f.BaseField.GetBaseExport(ctx, value)

	// Determine type based on subtype
	fieldType := "text"
	if f.Subtype == "textarea" || f.Subtype == "html" {
		fieldType = "textarea"
	}
	result["type"] = fieldType
	result["subtype"] = f.Subtype

	// Add min/max length
	if f.MinLength > 0 {
		result["min"] = f.MinLength
	}
	if f.MaxLength > 0 {
		result["max"] = f.MaxLength
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

	return result
}

// ============================================================================
// Chainable Setter Methods
// ============================================================================

// SetClass sets the CSS class for frontend styling
func (f *TextField) SetClass(class string) *TextField {
	f.BaseField.SetClass(class)
	return f
}

// SetHint sets the tooltip/help text for the field
func (f *TextField) SetHint(hint string) *TextField {
	f.BaseField.SetHint(hint)
	return f
}

// SetStep sets the step indicator for multi-step forms
func (f *TextField) SetStep(step int) *TextField {
	f.BaseField.SetStep(step)
	return f
}

// SetDisabled sets whether the field is disabled
func (f *TextField) SetDisabled(disabled bool) *TextField {
	f.BaseField.SetDisabled(disabled)
	return f
}

// SetAccess sets the access control permissions
func (f *TextField) SetAccess(access []string) *TextField {
	f.BaseField.SetAccess(access)
	return f
}

// SetScenario sets which scenarios this field applies to
func (f *TextField) SetScenario(scenario []string) *TextField {
	f.BaseField.SetScenario(scenario)
	return f
}

// SetForm sets whether to show in form
func (f *TextField) SetForm(form bool) *TextField {
	f.BaseField.SetForm(form)
	return f
}
