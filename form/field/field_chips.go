package field

import (
	"fmt"

	"github.com/xiriframework/xiri-go/component/core"
)

// ChipsField represents a tag/chip input form field.
//
// Each chip value is either a numeric option ID (int64, referencing an entry in
// List) or a free-text string. Chips chosen from List are transmitted and stored
// as their option ID; free-text input stays a string. The JSON type of each array
// element carries the distinction: number = option ID, string = free text.
type ChipsField struct {
	*BaseField
	List     []SelectOption // Available options for autocomplete
	FreeText bool           // Whether to allow free text input (not just from list)
	Value    []interface{}  // Parsed and validated values: int64 (option ID) or string (free text)
}

// optionExists reports whether the given int64 id matches one of the field's
// option values (which may be stored as int, int32, int64 or float64).
func (f *ChipsField) optionExists(id int64) bool {
	for _, opt := range f.List {
		switch v := opt.Value.(type) {
		case int:
			if int64(v) == id {
				return true
			}
		case int32:
			if int64(v) == id {
				return true
			}
		case int64:
			if v == id {
				return true
			}
		case float64:
			if int64(v) == id {
				return true
			}
		}
	}
	return false
}

// Validate validates the chips field value
func (f *ChipsField) Validate(value interface{}) error {
	if value == nil {
		if f.Required {
			return fmt.Errorf("chips field %s is required", f.ID)
		}
		return nil
	}

	arr, ok := value.([]interface{})
	if !ok {
		return fmt.Errorf("invalid chips value type for %s", f.ID)
	}

	if f.Required && len(arr) == 0 {
		return fmt.Errorf("chips field %s is required", f.ID)
	}

	// Numeric chips must reference a known option ID; strings are free text.
	for _, item := range arr {
		if id, ok := item.(int64); ok {
			if !f.optionExists(id) {
				return fmt.Errorf("chips field %s has no matching option for id %d", f.ID, id)
			}
		}
	}

	return nil
}

// Parse parses the raw value into a slice of int64 (option IDs) and strings
// (free text). Numeric elements are normalized to int64.
func (f *ChipsField) Parse(raw interface{}) (interface{}, error) {
	if raw == nil {
		return f.GetDefault(), nil
	}

	var items []interface{}
	switch v := raw.(type) {
	case []interface{}:
		items = v
	case []string:
		items = make([]interface{}, len(v))
		for i, s := range v {
			items[i] = s
		}
	default:
		return nil, fmt.Errorf("cannot parse chips from %T", raw)
	}

	result := make([]interface{}, len(items))
	for i, item := range items {
		switch val := item.(type) {
		case string:
			result[i] = val
		case int:
			result[i] = int64(val)
		case int32:
			result[i] = int64(val)
		case int64:
			result[i] = val
		case float64:
			result[i] = int64(val)
		default:
			return nil, fmt.Errorf("invalid chip value at index %d: %T", i, item)
		}
	}
	return result, nil
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

	if arr, ok := parsed.([]interface{}); ok {
		f.Value = arr
	} else {
		f.Value = nil
	}

	return nil
}

// IDs returns the numeric option IDs among the chip values (free-text chips are
// excluded).
func (f *ChipsField) IDs() []int64 {
	ids := make([]int64, 0, len(f.Value))
	for _, item := range f.Value {
		if id, ok := item.(int64); ok {
			ids = append(ids, id)
		}
	}
	return ids
}

// Texts returns the free-text chip values (option-ID chips are excluded).
func (f *ChipsField) Texts() []string {
	texts := make([]string, 0, len(f.Value))
	for _, item := range f.Value {
		if s, ok := item.(string); ok {
			texts = append(texts, s)
		}
	}
	return texts
}

// ExportForFrontend exports the field for frontend rendering
func (f *ChipsField) ExportForFrontend(ctx *core.UiContext, value interface{}) map[string]interface{} {
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
			Default:  []interface{}{},
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
