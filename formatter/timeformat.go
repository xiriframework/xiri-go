package formatter

import (
	"fmt"

	"github.com/xiriframework/xiri-go/uicontext"
)

// FormatTimeLengthHM converts seconds to "HH:MM" format.
// Example: 3665 seconds → "01:01"
// Uses the user's locale settings from UiContext.
func FormatTimeLengthHM(seconds int64, ctx *uicontext.UiContext) string {
	if seconds < 0 {
		return "00:00"
	}

	hours := seconds / 3600
	minutes := (seconds % 3600) / 60

	return fmt.Sprintf("%02d:%02d", hours, minutes)
}

// FormatTimeLengthMin converts seconds to minutes with "min" suffix.
// Example: 3665 seconds → "61 min"
// Note: Translation of "min" should be done by caller if needed.
func FormatTimeLengthMin(seconds int64, ctx *uicontext.UiContext) string {
	if seconds < 0 {
		return "0 min"
	}

	minutes := seconds / 60
	return fmt.Sprintf("%d min", minutes)
}

// FormatTimeLengthH converts seconds to hours with decimals and "h" suffix.
// Example: 3665 seconds → "1.0 h"
// Note: Translation of "h" should be done by caller if needed.
func FormatTimeLengthH(seconds int64, ctx *uicontext.UiContext) string {
	if seconds < 0 {
		return "0.0 h"
	}

	hours := float64(seconds) / 3600.0
	return fmt.Sprintf("%.1f h", hours)
}

// FormatTimeLengthHMS converts seconds to "HH:MM:SS" format.
// Example: 3665 seconds → "01:01:05"
// Uses the user's locale settings from UiContext.
func FormatTimeLengthHMS(seconds int64, ctx *uicontext.UiContext) string {
	if seconds < 0 {
		return "00:00:00"
	}

	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
}
