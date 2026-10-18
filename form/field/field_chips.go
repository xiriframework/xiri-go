package field

import (
	"fmt"

	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
)

// ChipsField represents a tag/chip input form field.
//
// Each chip value is either an option ID (referencing an entry in List) or a
// free-text string. Chips chosen from List are transmitted and stored as their
// option ID: numeric IDs as int64 (IDs()), string IDs as string (StrIDs()). A
// string that matches no string option is free text (Texts()).
type ChipsField struct {
	*BaseField
	List           []SelectOption // Available options for autocomplete
	FreeText       bool           // Whether to allow free text input (not just from list)
	Value          []interface{}  // Parsed and validated values: int64 (option ID) or string (option ID or free text)
	SuggestionsURL string         // Endpoint for suggestions while typing: POST {search, <SearchWith...>} -> [{id, name}]
	SearchWith     []string       // IDs of other form fields whose current value is sent with the suggestion request
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

// stringOptionExists reports whether s is the value of a string option.
func (f *ChipsField) stringOptionExists(s string) bool {
	for _, opt := range f.List {
		if v, ok := opt.Value.(string); ok && v == s {
			return true
		}
	}
	return false
}

// Validate validates the chips field value
func (f *ChipsField) Validate(value interface{}) error {
	if value == nil {
		if f.Required {
			return invalid("required", nil, "chips field %s is required", f.ID)
		}
		return nil
	}

	arr, ok := value.([]interface{})
	if !ok {
		return fmt.Errorf("invalid chips value type for %s", f.ID)
	}

	if f.Required && len(arr) == 0 {
		return invalid("required", nil, "chips field %s is required", f.ID)
	}

	// Numeric chips must reference a known option ID; strings are string options or free text.
	for _, item := range arr {
		switch v := item.(type) {
		case int64:
			if !f.optionExists(v) {
				return invalid("not_allowed", nil, "chips field %s has no matching option for id %d", f.ID, v)
			}
		case string:
			if !f.FreeText && !f.stringOptionExists(v) {
				return invalid("not_allowed", nil, "chips field %s allows no free text: %q", f.ID, v)
			}
		}
	}

	return nil
}

// Parse parses the raw value into a slice of int64 (option IDs) and strings
// (free text). Numeric elements are normalized to int64.
func (f *ChipsField) Parse(raw interface{}) (interface{}, error) {
	if raw == nil {
		def := f.GetDefault()
		if def == nil {
			return nil, nil
		}
		// A directly set default ([]string instead of []interface{}) is normalized like a request value.
		return f.Parse(def)
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

// StrIDs returns the string option IDs among the chip values (free text and
// numeric IDs are excluded).
func (f *ChipsField) StrIDs() []string {
	return f.strings(true)
}

// Texts returns the free-text chip values (option-ID chips are excluded).
func (f *ChipsField) Texts() []string {
	return f.strings(false)
}

func (f *ChipsField) strings(options bool) []string {
	out := make([]string, 0, len(f.Value))
	for _, item := range f.Value {
		if s, ok := item.(string); ok && f.stringOptionExists(s) == options {
			out = append(out, s)
		}
	}
	return out
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
	if !f.FreeText {
		result["freeText"] = false
	}
	if f.SuggestionsURL != "" {
		result["url"] = f.SuggestionsURL
		if len(f.SearchWith) > 0 {
			result["searchWith"] = f.SearchWith
		}
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

// SetSuggestionsURL enables server-side suggestions while typing, like
// TextField.SetSuggestionsURL: the frontend posts {search: "<input>", <id>: <value> ...}
// for every enabled form field named in searchWith and expects [{id, name}]. The
// chosen suggestion's name is inserted as a free-text chip (unless it equals the
// name of a List option, then that option's ID), so this is only useful with
// FreeText=true. Use builder.BindSuggest in the handler. A nil url disables the
// feature; the reserved id "search" is dropped from searchWith.
func (f *ChipsField) SetSuggestionsURL(u *url.Url, searchWith ...string) *ChipsField {
	f.SuggestionsURL, f.SearchWith = suggestionsURL(u, searchWith)
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
