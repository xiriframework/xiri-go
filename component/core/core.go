package core

// Component is the base interface for all Xiri UI components.
// All components must implement a Print method that returns their JSON representation
// for the Angular frontend.
//
// Components follow these design principles:
// 1. Type safety - avoid interface{} where possible
// 2. Builder pattern - provide Gen() and specialized constructors
// 3. Fluent interface - return *Type for method chaining where appropriate
// 4. Translation support - accept optional TranslateFunc for i18n
type Component interface {
	// Print returns the JSON representation of the component.
	// The translator parameter is optional and can be nil.
	Print(translator TranslateFunc) map[string]any
}

// TranslateFunc is a function type for translating text keys to localized strings.
// If nil, components should output text directly without translation.
type TranslateFunc func(key string) string

// WithNewRow marks a printed component result to start a new grid row.
func WithNewRow(result map[string]any) map[string]any {
	result["newRow"] = true
	return result
}

// Translate applies translation if translator is provided
func Translate(translator TranslateFunc, key string) string {
	if translator != nil {
		return translator(key)
	}
	return key
}
