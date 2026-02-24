package group

import (
	"fmt"

	"github.com/xiriframework/xiri-go/formatter"
)

// FormatNumber formats a number according to the user's locale
// Locale 0 (German): uses comma as decimal separator (e.g., "1,234.56" -> "1.234,56")
// Locale 1+ (English, etc.): uses dot as decimal separator (e.g., "1,234.56")
func (fg *FormGroup) FormatNumber(value float64, decimals int) string {
	if fg.ctx == nil {
		// Default to English format
		return fmt.Sprintf("%%.%df", decimals)
	}

	return formatter.FormatNumberLocale(value, decimals, fg.ctx.Locale)
}

// FormatDistance formats a distance value according to user's distance unit preference
func (fg *FormGroup) FormatDistance(km float64) string {
	if fg.ctx == nil {
		return formatter.FormatNumberLocale(km, 0, 0) + " km"
	}

	return formatter.FormatDistanceLocaleWithDecimals(km, fg.ctx.Distance, fg.ctx.Locale, 2)
}

// FormatPressure formats a pressure value according to user's pressure unit preference
func (fg *FormGroup) FormatPressure(bar float64) string {
	if fg.ctx == nil {
		return formatter.FormatNumberLocale(bar, 2, 0) + " bar"
	}

	return formatter.FormatPressureLocale(bar, fg.ctx.Pressure, fg.ctx.Locale)
}

// FormatSpeed formats a speed value according to user's distance unit preference
func (fg *FormGroup) FormatSpeed(kmh float64) string {
	if fg.ctx == nil {
		return formatter.FormatNumberLocale(kmh, 0, 0) + " km/h"
	}

	return formatter.FormatSpeedLocale(kmh, fg.ctx.Distance, fg.ctx.Locale)
}
