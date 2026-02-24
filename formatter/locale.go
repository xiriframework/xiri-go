package formatter

import (
	"fmt"

	"github.com/xiriframework/xiri-go/types/distance"
	"github.com/xiriframework/xiri-go/types/locale"
	"github.com/xiriframework/xiri-go/types/pressure"
)

// usesCommaDecimal returns true if the locale uses comma as decimal separator
// and dot as thousands separator (continental European convention).
// Returns false for locales using dot as decimal separator (English, Japanese, Chinese, Arabic).
func usesCommaDecimal(loc locale.Locale) bool {
	switch loc {
	case locale.EnGB, locale.EnUS, locale.Ja, locale.ZhCN, locale.ArAE:
		return false
	default:
		return true
	}
}

// FormatNumberLocale formats a number according to the user's locale
// Comma-decimal locales: uses comma as decimal separator and dot as thousands separator (e.g., "5000" → "5.000")
// Dot-decimal locales: uses dot as decimal separator and comma as thousands separator (e.g., "5000" → "5,000")
func FormatNumberLocale(value float64, decimals int, loc locale.Locale) string {
	format := fmt.Sprintf("%%.%df", decimals)
	str := fmt.Sprintf(format, value)

	if usesCommaDecimal(loc) {
		return addThousandSeparatorsLocale(str, '.', ',')
	}

	return addThousandSeparatorsLocale(str, ',', '.')
}

// FormatDistanceLocaleWithDecimals formats distance with configurable decimal places
// Returns formatted string with appropriate unit (km, mi, or NM for nautical miles)
func FormatDistanceLocaleWithDecimals(km float64, distUnit distance.Distance, loc locale.Locale, decimals int) string {
	switch distUnit {
	case distance.Miles:
		miles := km * 0.621371
		return FormatNumberLocale(miles, decimals, loc) + " mi"
	case distance.Seemiles:
		nauticalMiles := km * 0.539957
		return FormatNumberLocale(nauticalMiles, decimals, loc) + " NM"
	default: // distance.Kilometer
		return FormatNumberLocale(km, decimals, loc) + " km"
	}
}

// ConvertDistanceToKm converts a distance value from user's unit to kilometers
// Input: value in user's preferred unit (km, mi, or NM)
// Output: value in kilometers
func ConvertDistanceToKm(value float64, distUnit distance.Distance) float64 {
	switch distUnit {
	case distance.Miles:
		return value / 0.621371 // miles to km
	case distance.Seemiles:
		return value / 0.539957 // nautical miles to km
	default: // distance.Kilometer
		return value // already in km
	}
}

// FormatPressureLocale formats pressure according to user's pressure unit preference
// Returns formatted string with appropriate unit (bar, psi, or kPa)
func FormatPressureLocale(bar float64, pressUnit pressure.Pressure, loc locale.Locale) string {
	switch pressUnit {
	case pressure.Psi:
		psi := bar * 14.5038
		return FormatNumberLocale(psi, 1, loc) + " psi"
	case pressure.Kpa:
		kpa := bar * 100
		return FormatNumberLocale(kpa, 1, loc) + " kPa"
	default:
		return FormatNumberLocale(bar, 1, loc) + " bar"
	}
}

// FormatSpeedLocale formats speed according to distance unit preference
// Returns formatted string with appropriate unit (km/h, mph, or knots)
func FormatSpeedLocale(kmh float64, distUnit distance.Distance, loc locale.Locale) string {
	switch distUnit {
	case distance.Miles:
		mph := kmh * 0.621371
		return FormatNumberLocale(mph, 1, loc) + " mph"
	case distance.Seemiles:
		knots := kmh * 0.539957
		return FormatNumberLocale(knots, 1, loc) + " kn"
	default: // distance.Kilometer
		return FormatNumberLocale(kmh, 1, loc) + " km/h"
	}
}

// addThousandSeparatorsLocale adds thousands separators and sets decimal separator according to locale
// thousandsSep: character to use for thousands separator (e.g., '.' for German, ',' for English)
// decimalSep: character to use for decimal separator (e.g., ',' for German, '.' for English)
func addThousandSeparatorsLocale(numStr string, thousandsSep rune, decimalSep rune) string {
	// Find decimal point in the formatted number
	decimalPos := -1
	for i, char := range numStr {
		if char == '.' {
			decimalPos = i
			break
		}
	}

	// Separate integer and decimal parts
	var intPart, decPart string
	if decimalPos == -1 {
		intPart = numStr
		decPart = ""
	} else {
		intPart = numStr[:decimalPos]
		decPart = numStr[decimalPos+1:] // Skip the '.' itself
	}

	// Handle negative sign
	negative := false
	if len(intPart) > 0 && intPart[0] == '-' {
		negative = true
		intPart = intPart[1:]
	}

	// Add thousands separators to integer part
	result := ""
	for i, digit := range reverseString(intPart) {
		if i > 0 && i%3 == 0 {
			result = string(thousandsSep) + result
		}
		result = string(digit) + result
	}

	if negative {
		result = "-" + result
	}

	// Add decimal part with correct separator
	if len(decPart) > 0 {
		result = result + string(decimalSep) + decPart
	}

	return result
}

// reverseString reverses a string
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
