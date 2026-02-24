package table

import "github.com/xiriframework/xiri-go/uicontext"

// OutputFormatter formats a field value for a specific output type.
// Formatters receive the entire row for cross-field dependencies (e.g., distance formatter needs device_id).
//
// All formatters are automatically UiContext-aware and use locale-specific formatting
// (decimal separator, date format, timezone, unit conversion) based on the user's preferences.
type OutputFormatter interface {
	// Format formats the field value for the given output type.
	//
	// Parameters:
	//   value: The field's own value (from the accessor function)
	//   row: Access to all other field values in the row (for cross-field dependencies)
	//   output: Target output type (Web, CSV, PDF, Excel)
	//   ctx: User context with locale/language/timezone/unit preferences
	//
	// Returns:
	//   Formatted value - typically a string, but can be any (e.g., for web number fields)
	Format(value any, row Row, output OutputType, ctx *uicontext.UiContext) any
}

// FormatterFunc is a function adapter for the OutputFormatter interface.
// This allows using simple functions as formatters without defining a struct.
//
// Example:
//
//	upperCaseFormatter := FormatterFunc(func(value any, row Row, output OutputType, ctx *uicontext.UiContext) any {
//	    return strings.ToUpper(fmt.Sprint(value))
//	})
type FormatterFunc func(any, Row, OutputType, *uicontext.UiContext) any

// Format implements the OutputFormatter interface
func (f FormatterFunc) Format(value any, row Row, output OutputType, ctx *uicontext.UiContext) any {
	return f(value, row, output, ctx)
}
