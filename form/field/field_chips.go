package field

import (
	"fmt"

	"github.com/xiriframework/xiri-go/uicontext"
)

// ChipsField represents a tag/chip input form field
type ChipsField struct {
	*BaseField
	List     []SelectOption // Available options for autocomplete
	FreeText bool           // Whether to allow free text input (not just from list)
	Value    []string       // Parsed and validated values
}

// Validate validates the chips field value
func (f *ChipsField) Validate(value interface{}) error {
	if value == nil {
		if f.Required {
			return fmt.Errorf("chips field %s is required", f.ID)
		}
		return nil
	}

	arr, ok := value.([]string)
	if !ok {
		return fmt.Errorf("invalid chips value type for %s", f.ID)
	}

	if f.Required && len(arr) == 0 {
		return fmt.Errorf("chips field %s is required", f.ID)
	}

	return nil
}

// Parse parses raw value into string slice
func (f *ChipsField) Parse(raw interface{}) (interface{}, error) {
	if raw == nil {
		return f.GetDefault(), nil
	}

	switch v := raw.(type) {
	case []string:
		return v, nil
	case []interface{}:
		result := make([]string, len(v))
		for i, item := range v {
			s, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("invalid chip value at index %d", i)
			}
			result[i] = s
		}
		return result, nil
	default:
		return nil, fmt.Errorf("cannot parse chips from %T", raw)
	}
}

// BindValue parses, validates, and stores the value
func (f *ChipsField) BindValue(raw interface{}) error {
	parsed, err := f.Parse(raw)
	if err != nil {
		return fmt.Errorf("parsing field %s: %w", f.ID, err)
	}

	if err := f.Validate(parsed); err != nil {
		return fmt.Errorf("validating field %s: %w", f.ID, err)
	}

	if parsed != nil {
		if arr, ok := parsed.([]string); ok {
			f.Value = arr
		}
	} else {
		f.Value = nil
	}

	return nil
}

// ExportForFrontend exports the field for frontend rendering
func (f *ChipsField) ExportForFrontend(ctx *uicontext.UiContext, value interface{}) map[string]interface{} {
	if value == nil {
		value = f.GetDefault()
	}
	result := f.BaseField.GetBaseExport(ctx, value)

	if len(f.List) > 0 {
		list := make([]map[string]interface{}, len(f.List))
		for i, opt := range f.List {
			list[i] = map[string]interface{}{
				"id":   opt.Value,
				"name": opt.Label,
			}
		}
		result["list"] = list
	}

	return result
}

// NewChipsField creates a new chips/tag input field
func NewChipsField(id, name string, required bool) *ChipsField {
	return &ChipsField{
		BaseField: &BaseField{
			ID:       id,
			Type:     FieldTypeChips,
			Name:     name,
			Required: required,
			Default:  []string{},
			Form:     true,
		},
		FreeText: true,
	}
}

// SetList sets the available options for autocomplete
func (f *ChipsField) SetList(list []SelectOption) *ChipsField {
	f.List = list
	return f
}

// SetFreeText sets whether free text input is allowed
func (f *ChipsField) SetFreeText(freeText bool) *ChipsField {
	f.FreeText = freeText
	return f
}

// SetClass sets the CSS class
func (f *ChipsField) SetClass(class string) *ChipsField {
	f.BaseField.SetClass(class)
	return f
}

// SetHint sets the hint text
func (f *ChipsField) SetHint(hint string) *ChipsField {
	f.BaseField.SetHint(hint)
	return f
}
