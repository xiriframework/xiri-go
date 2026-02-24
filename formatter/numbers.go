package formatter

import (
	"github.com/xiriframework/xiri-go/uicontext"
)

// FormatDouble2 formats a float with exactly 2 decimal places using locale settings.
// Example (en): 123.456 → "123.46"
// Example (de): 123.456 → "123,46"
// Uses the user's locale settings from UiContext.
func FormatDouble2(value float64, ctx *uicontext.UiContext) string {
	return FormatNumberLocale(value, 2, ctx.Locale)
}

// FormatInteger formats an integer value with thousand separators based on locale.
// Example (en): 1234567 → "1,234,567"
// Example (de): 1234567 → "1.234.567"
// Uses the user's locale settings from UiContext.
func FormatInteger(value int64, ctx *uicontext.UiContext) string {
	return FormatNumberLocale(float64(value), 0, ctx.Locale)
}

// FormatBigNumber formats a large number with thousand separators and no decimals.
// Example (en): 1234567.89 → "1,234,568"
// Example (de): 1234567.89 → "1.234.568"
func FormatBigNumber(value float64, ctx *uicontext.UiContext) string {
	return FormatNumberLocale(value, 0, ctx.Locale)
}
