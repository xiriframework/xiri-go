package field

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
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
	Trim       bool    // Parse trims surrounding whitespace; true via constructors, false in a struct literal; never for "password"
	Value      *string // Parsed and validated value (type-safe access)

	Suggestions    []string // Autocomplete suggestions; nil = none, empty = explicitly none (exports list: [])
	SuggestionsURL string   // Endpoint for suggestions while typing: POST {search, <SearchWith...>} -> [{id, name}]
	SearchWith     []string // IDs of other form fields whose current value is sent with the suggestion request
}

func (f *TextField) Validate(value interface{}) error {
	// Vor jedem Leerwert-Rücksprung, sonst fällt ein kaputtes Pattern an optionalen Feldern nie auf.
	var re *regexp.Regexp
	if f.Pattern != "" {
		var err error
		if re, err = compilePattern(f.Pattern); err != nil {
			return fmt.Errorf("text field %s: %w", f.ID, err)
		}
	}

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

	if f.Required && strings.TrimSpace(str) == "" {
		return fmt.Errorf("text field %s is required", f.ID)
	}

	if f.MinLength > 0 && len(str) < f.MinLength {
		return fmt.Errorf("text field %s must be at least %d characters", f.ID, f.MinLength)
	}

	if f.MaxLength > 0 && len(str) > f.MaxLength {
		return fmt.Errorf("text field %s must be at most %d characters", f.ID, f.MaxLength)
	}

	if re != nil && str != "" && !re.MatchString(str) {
		return fmt.Errorf("text field %s has invalid format", f.ID)
	}

	return nil
}

// goOnlySyntax lists RE2 constructs that make the browser's new RegExp throw or match something else.
// ponytail: substring blocklist, not a JS parser - incomplete, and an escaped backslash before e.g. A
// ("\\A") is a false positive. Replace with a real shared-subset parser if patterns get fancier.
var goOnlySyntax = regexp.MustCompile(`\(\?[^:]|\\[AzQpP]|\\x\{|\[\[:`)

// compilePattern compiles p the way Angular's Validators.pattern does: "^" is prepended unless p starts
// with it, "$" appended unless p ends with it - so "a|b" becomes "^a|b$", not "^(?:a|b)$".
func compilePattern(p string) (*regexp.Regexp, error) {
	if m := goOnlySyntax.FindString(p); m != "" {
		return nil, fmt.Errorf("pattern %q uses %q, which browsers do not support the same way", p, m)
	}
	if !strings.HasPrefix(p, "^") {
		p = "^" + p
	}
	if !strings.HasSuffix(p, "$") {
		p += "$"
	}
	return regexp.Compile(p)
}

func (f *TextField) Parse(raw interface{}) (interface{}, error) {
	if raw == nil {
		raw = f.GetDefault()
		if raw == nil {
			return nil, nil
		}
	}

	str, ok := raw.(string)
	if !ok {
		str = fmt.Sprintf("%v", raw)
	}

	if f.Trim && f.Subtype != "password" {
		str = strings.TrimSpace(str)
	}
	return str, nil
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
		Trim: true,
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
		Trim:      true,
	}
}

// ExportForFrontend exports the field for frontend rendering
func (f *TextField) ExportForFrontend(ctx *core.UiContext, value interface{}) map[string]interface{} {
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

	// An invalid pattern is left out so it cannot break the browser form; Validate rejects it anyway.
	if f.Pattern != "" {
		if _, err := compilePattern(f.Pattern); err == nil {
			result["pattern"] = f.Pattern
		}
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

	// Suggestions: nil means "not configured" (no key), an empty slice means "explicitly none" -
	// exported as list: [] so a reload patch can clear previous suggestions.
	if f.Suggestions != nil {
		list := make([]map[string]interface{}, len(f.Suggestions))
		for i, s := range f.Suggestions {
			list[i] = map[string]interface{}{"id": s, "name": s}
		}
		result["list"] = list
	}
	if f.SuggestionsURL != "" {
		result["url"] = f.SuggestionsURL
		if len(f.SearchWith) > 0 {
			result["searchWith"] = f.SearchWith
		}
	}

	return result
}

// SetSuggestions sets a fixed list of autocomplete suggestions; any other text stays valid.
// Calling it without arguments exports an empty list (clears suggestions via reload patch).
func (f *TextField) SetSuggestions(suggestions ...string) *TextField {
	if suggestions == nil {
		suggestions = []string{}
	}
	f.Suggestions = suggestions
	return f
}

// SetSuggestionsURL enables server-side suggestions while typing.
//
// The frontend posts {search: "<input>", <id>: <current value> ...} for every enabled form
// field named in searchWith and expects [{id, name}]; name is what gets inserted. Use
// builder.BindSuggest in the handler. A nil url disables the feature; the reserved id "search"
// is dropped from searchWith.
func (f *TextField) SetSuggestionsURL(u *url.Url, searchWith ...string) *TextField {
	if u == nil || u.PrintPrefix() == "" {
		f.SuggestionsURL, f.SearchWith = "", nil
		return f
	}
	f.SuggestionsURL = u.PrintPrefix()
	f.SearchWith = nil
	for _, id := range searchWith {
		if id != "search" {
			f.SearchWith = append(f.SearchWith, id)
		}
	}
	return f
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

// SetTrim sets whether Parse trims surrounding whitespace (constructors default to true).
// Subtype "password" is never trimmed.
func (f *TextField) SetTrim(trim bool) *TextField {
	f.Trim = trim
	return f
}

// SetPattern sets a regex the whole value must match, checked in the browser and in Validate.
// Anchoring follows Angular (see compilePattern). Use only syntax both Go (RE2) and JavaScript
// understand; SetPattern panics on invalid or Go-only patterns, since that is a programming error.
func (f *TextField) SetPattern(pattern string) *TextField {
	if pattern != "" {
		if _, err := compilePattern(pattern); err != nil {
			panic(err)
		}
	}
	f.Pattern = pattern
	return f
}
