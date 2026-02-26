package table

import (
	"fmt"
	"strconv"
	"time"

	"github.com/xiriframework/xiri-go/formatter"
	"github.com/xiriframework/xiri-go/types/distance"
	"github.com/xiriframework/xiri-go/types/pressure"
	"github.com/xiriframework/xiri-go/uicontext"
)

// ============================================================================
// Formatter creation functions
// These create the actual formatter functions that will be used at runtime
// ============================================================================

func createIdFormatter() OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *uicontext.UiContext) any {
		return toInt64(value)
	})
}

func createIntegerFormatter() OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *uicontext.UiContext) any {
		num := toInt64(value)
		switch output {
		case OutputWeb, OutputPDF:
			return formatter.FormatNumberLocale(float64(num), 0, ctx.Locale)
		case OutputCSV, OutputExcel:
			return fmt.Sprint(num)
		}
		return fmt.Sprint(num)
	})
}

func createFloatFormatter(decimals int) OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *uicontext.UiContext) any {
		num := toFloat64(value)
		switch output {
		case OutputWeb, OutputPDF:
			return formatter.FormatNumberLocale(num, decimals, ctx.Locale)
		case OutputCSV, OutputExcel:
			format := "%." + strconv.Itoa(decimals) + "f"
			return fmt.Sprintf(format, num)
		}
		format := "%." + strconv.Itoa(decimals) + "f"
		return fmt.Sprintf(format, num)
	})
}

func createTextFormatter() OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *uicontext.UiContext) any {
		if value == nil {
			return ""
		}
		return fmt.Sprint(value)
	})
}

// createPassthroughFormatter returns a formatter that preserves complex data structures (maps, arrays)
// for JSON serialization. For CSV output, converts to string representation.
func createPassthroughFormatter() OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *uicontext.UiContext) any {
		if value == nil {
			return nil
		}
		// For CSV output, convert to string
		if output == OutputCSV || output == OutputPDF || output == OutputExcel {
			return fmt.Sprint(value)
		}
		// For web output, return as-is to preserve structure for JSON
		return value
	})
}

// createText2Formatter returns a formatter for two-line text fields.
// Expects [2]string array: [0] = primary text, [1] = secondary text
func createText2Formatter() OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *uicontext.UiContext) any {
		if value == nil {
			return [2]string{"", ""}
		}
		arr, ok := value.([2]string)
		if !ok {
			return [2]string{"", ""}
		}
		if output == OutputCSV || output == OutputExcel {
			if arr[1] == "" {
				return arr[0]
			}
			return arr[0] + " - " + arr[1]
		}
		return arr
	})
}

// createText2IntFormatter returns a formatter for two-line integer fields.
// Expects [2]int array: [0] = primary value, [1] = secondary value
func createText2IntFormatter() OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *uicontext.UiContext) any {
		if value == nil {
			return [2]string{"", ""}
		}
		arr, ok := value.([2]int)
		if !ok {
			return [2]string{"", ""}
		}
		switch output {
		case OutputWeb, OutputPDF:
			return [2]string{
				formatter.FormatNumberLocale(float64(arr[0]), 0, ctx.Locale),
				formatter.FormatNumberLocale(float64(arr[1]), 0, ctx.Locale),
			}
		case OutputCSV, OutputExcel:
			primary := fmt.Sprint(arr[0])
			secondary := fmt.Sprint(arr[1])
			if secondary == "" || secondary == "0" {
				return primary
			}
			return primary + " - " + secondary
		}
		return [2]string{fmt.Sprint(arr[0]), fmt.Sprint(arr[1])}
	})
}

// createText2FloatFormatter returns a formatter for two-line float fields.
// Expects [2]float64 array: [0] = primary value, [1] = secondary value
func createText2FloatFormatter(decimals int) OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *uicontext.UiContext) any {
		if value == nil {
			return [2]string{"", ""}
		}
		arr, ok := value.([2]float64)
		if !ok {
			return [2]string{"", ""}
		}
		format := "%." + strconv.Itoa(decimals) + "f"
		switch output {
		case OutputWeb, OutputPDF:
			return [2]string{
				formatter.FormatNumberLocale(arr[0], decimals, ctx.Locale),
				formatter.FormatNumberLocale(arr[1], decimals, ctx.Locale),
			}
		case OutputCSV, OutputExcel:
			primary := fmt.Sprintf(format, arr[0])
			secondary := fmt.Sprintf(format, arr[1])
			if secondary == "" {
				return primary
			}
			return primary + " - " + secondary
		}
		return [2]string{
			fmt.Sprintf(format, arr[0]),
			fmt.Sprintf(format, arr[1]),
		}
	})
}

// createText2DateTimeFormatter returns a formatter for two-line datetime fields.
// Expects [2]time.Time array: [0] = primary time, [1] = secondary time
func createText2DateTimeFormatter() OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *uicontext.UiContext) any {
		if value == nil {
			return [2]string{"", ""}
		}
		arr, ok := value.([2]time.Time)
		if !ok {
			return [2]string{"", ""}
		}

		loc, err := time.LoadLocation(ctx.Timezone.GetIANA())
		if err != nil {
			loc = time.UTC
		}

		switch output {
		case OutputWeb, OutputPDF:
			primary := ""
			secondary := ""
			if !arr[0].IsZero() {
				primary = formatter.FormatTimestampDateTime(arr[0].Unix(), ctx)
			}
			if !arr[1].IsZero() {
				secondary = formatter.FormatTimestampDateTime(arr[1].Unix(), ctx)
			}
			return [2]string{primary, secondary}
		case OutputCSV, OutputExcel:
			primary := ""
			secondary := ""
			if !arr[0].IsZero() {
				primary = arr[0].In(loc).Format("2006-01-02 15:04:05")
			}
			if !arr[1].IsZero() {
				secondary = arr[1].In(loc).Format("2006-01-02 15:04:05")
			}
			if secondary == "" {
				return primary
			}
			return primary + " - " + secondary
		}
		primary := ""
		secondary := ""
		if !arr[0].IsZero() {
			primary = arr[0].In(loc).Format("2006-01-02 15:04:05")
		}
		if !arr[1].IsZero() {
			secondary = arr[1].In(loc).Format("2006-01-02 15:04:05")
		}
		return [2]string{primary, secondary}
	})
}

// createText2DateFormatter returns a formatter for two-line date fields.
// Expects [2]time.Time array: [0] = primary date, [1] = secondary date
func createText2DateFormatter() OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *uicontext.UiContext) any {
		if value == nil {
			return [2]string{"", ""}
		}
		arr, ok := value.([2]time.Time)
		if !ok {
			return [2]string{"", ""}
		}

		loc, err := time.LoadLocation(ctx.Timezone.GetIANA())
		if err != nil {
			loc = time.UTC
		}

		switch output {
		case OutputWeb, OutputPDF:
			primary := ""
			secondary := ""
			if !arr[0].IsZero() {
				primary = formatter.FormatTimestampDate(arr[0].Unix(), ctx)
			}
			if !arr[1].IsZero() {
				secondary = formatter.FormatTimestampDate(arr[1].Unix(), ctx)
			}
			return [2]string{primary, secondary}
		case OutputCSV, OutputExcel:
			primary := ""
			secondary := ""
			if !arr[0].IsZero() {
				primary = arr[0].In(loc).Format("2006-01-02")
			}
			if !arr[1].IsZero() {
				secondary = arr[1].In(loc).Format("2006-01-02")
			}
			if secondary == "" {
				return primary
			}
			return primary + " - " + secondary
		}
		primary := ""
		secondary := ""
		if !arr[0].IsZero() {
			primary = arr[0].In(loc).Format("2006-01-02")
		}
		if !arr[1].IsZero() {
			secondary = arr[1].In(loc).Format("2006-01-02")
		}
		return [2]string{primary, secondary}
	})
}

// createText2DistanceFormatter returns a formatter for two-line distance fields.
// Expects [2]float64 array (values in kilometers): [0] = primary, [1] = secondary
func createText2DistanceFormatter(decimals int) OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *uicontext.UiContext) any {
		if value == nil {
			return [2]string{"", ""}
		}
		arr, ok := value.([2]float64)
		if !ok {
			return [2]string{"", ""}
		}

		distanceUnit := ctx.Distance

		switch output {
		case OutputWeb, OutputPDF:
			return [2]string{
				formatter.FormatDistanceLocaleWithDecimals(arr[0], distanceUnit, ctx.Locale, decimals),
				formatter.FormatDistanceLocaleWithDecimals(arr[1], distanceUnit, ctx.Locale, decimals),
			}
		case OutputCSV, OutputExcel:
			format := "%." + strconv.Itoa(decimals) + "f"
			primary := fmt.Sprintf(format, convertDistanceValue(arr[0], distanceUnit))
			secondary := fmt.Sprintf(format, convertDistanceValue(arr[1], distanceUnit))
			if secondary == "" {
				return primary
			}
			return primary + " - " + secondary
		}
		format := "%." + strconv.Itoa(decimals) + "f"
		return [2]string{
			fmt.Sprintf(format, convertDistanceValue(arr[0], distanceUnit)),
			fmt.Sprintf(format, convertDistanceValue(arr[1], distanceUnit)),
		}
	})
}

// createText2SpeedFormatter returns a formatter for two-line speed fields.
// Expects [2]float64 array (values in km/h): [0] = primary, [1] = secondary
func createText2SpeedFormatter(decimals int) OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *uicontext.UiContext) any {
		if value == nil {
			return [2]string{"", ""}
		}
		arr, ok := value.([2]float64)
		if !ok {
			return [2]string{"", ""}
		}

		distanceUnit := ctx.Distance // Speed uses distance unit

		// Determine unit suffix
		unit := " km/h"
		switch distanceUnit {
		case distance.Miles:
			unit = " mph"
		case distance.Seemiles:
			unit = " kn"
		}

		switch output {
		case OutputWeb, OutputPDF:
			primary := formatter.FormatNumberLocale(convertSpeedValue(arr[0], distanceUnit), decimals, ctx.Locale) + unit
			secondary := formatter.FormatNumberLocale(convertSpeedValue(arr[1], distanceUnit), decimals, ctx.Locale) + unit
			return [2]string{primary, secondary}
		case OutputCSV, OutputExcel:
			format := "%." + strconv.Itoa(decimals) + "f"
			primary := fmt.Sprintf(format, convertSpeedValue(arr[0], distanceUnit))
			secondary := fmt.Sprintf(format, convertSpeedValue(arr[1], distanceUnit))
			if secondary == "" {
				return primary
			}
			return primary + " - " + secondary
		}
		format := "%." + strconv.Itoa(decimals) + "f"
		return [2]string{
			fmt.Sprintf(format, convertSpeedValue(arr[0], distanceUnit)),
			fmt.Sprintf(format, convertSpeedValue(arr[1], distanceUnit)),
		}
	})
}

// createText2BoolFormatter returns a formatter for two-line boolean fields.
// Expects [2]bool array: [0] = primary value, [1] = secondary value
func createText2BoolFormatter() OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *uicontext.UiContext) any {
		if value == nil {
			return [2]string{"", ""}
		}
		arr, ok := value.([2]bool)
		if !ok {
			return [2]string{"", ""}
		}

		primary := "No"
		if arr[0] {
			primary = "Yes"
		}
		secondary := "No"
		if arr[1] {
			secondary = "Yes"
		}

		switch output {
		case OutputWeb, OutputPDF:
			return [2]string{primary, secondary}
		case OutputCSV, OutputExcel:
			if secondary == "" {
				return primary
			}
			return primary + " - " + secondary
		}
		return [2]string{primary, secondary}
	})
}

// createTimeLengthFormatter returns a formatter for time duration fields.
// Expects int64 value in seconds.
// Web/PDF: "HH:MM" or "Xd HH:MM" format.
// CSV/Excel: integer minutes.
func createTimeLengthFormatter() OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *uicontext.UiContext) any {
		seconds := toInt64(value)

		switch output {
		case OutputWeb, OutputPDF:
			return formatTimeLength(seconds)
		case OutputCSV, OutputExcel:
			return fmt.Sprintf("%d", seconds/60)
		}
		return fmt.Sprintf("%d", seconds/60)
	})
}

// createText2TimeLengthFormatter returns a formatter for two-line time duration fields.
// Expects [2]int64 array (values in seconds): [0] = primary, [1] = secondary
func createText2TimeLengthFormatter() OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *uicontext.UiContext) any {
		if value == nil {
			return [2]string{"", ""}
		}
		arr, ok := value.([2]int64)
		if !ok {
			return [2]string{"", ""}
		}

		switch output {
		case OutputWeb, OutputPDF:
			return [2]string{
				formatTimeLength(arr[0]),
				formatTimeLength(arr[1]),
			}
		case OutputCSV, OutputExcel:
			primary := fmt.Sprintf("%d", arr[0]/60)
			secondary := fmt.Sprintf("%d", arr[1]/60)
			if secondary == "" || secondary == "0" {
				return primary
			}
			return primary + " - " + secondary
		}
		return [2]string{
			fmt.Sprintf("%d", arr[0]/60),
			fmt.Sprintf("%d", arr[1]/60),
		}
	})
}

// ============================================================================
// N-line (variable length) formatter functions
// These create formatters for fields with a variable number of lines (slices).
// ============================================================================

// createTextNFormatter returns a formatter for variable-line text fields.
// Expects []string slice.
func createTextNFormatter() OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *uicontext.UiContext) any {
		if value == nil {
			return []string{}
		}
		arr, ok := value.([]string)
		if !ok {
			return []string{}
		}
		if output == OutputCSV || output == OutputExcel {
			return arr
		}
		return arr
	})
}

// createIntegerNFormatter returns a formatter for variable-line integer fields.
// Expects []int slice.
func createIntegerNFormatter() OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *uicontext.UiContext) any {
		if value == nil {
			return []string{}
		}
		arr, ok := value.([]int)
		if !ok {
			return []string{}
		}
		switch output {
		case OutputWeb, OutputPDF:
			result := make([]string, len(arr))
			for i, v := range arr {
				result[i] = formatter.FormatNumberLocale(float64(v), 0, ctx.Locale)
			}
			return result
		case OutputCSV, OutputExcel:
			strs := make([]string, len(arr))
			for i, v := range arr {
				strs[i] = fmt.Sprint(v)
			}
			return strs
		}
		result := make([]string, len(arr))
		for i, v := range arr {
			result[i] = fmt.Sprint(v)
		}
		return result
	})
}

// createFloatNFormatter returns a formatter for variable-line float fields.
// Expects []float64 slice.
func createFloatNFormatter(decimals int) OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *uicontext.UiContext) any {
		if value == nil {
			return []string{}
		}
		arr, ok := value.([]float64)
		if !ok {
			return []string{}
		}
		format := "%." + strconv.Itoa(decimals) + "f"
		switch output {
		case OutputWeb, OutputPDF:
			result := make([]string, len(arr))
			for i, v := range arr {
				result[i] = formatter.FormatNumberLocale(v, decimals, ctx.Locale)
			}
			return result
		case OutputCSV, OutputExcel:
			strs := make([]string, len(arr))
			for i, v := range arr {
				strs[i] = fmt.Sprintf(format, v)
			}
			return strs
		}
		result := make([]string, len(arr))
		for i, v := range arr {
			result[i] = fmt.Sprintf(format, v)
		}
		return result
	})
}

// createDateTimeNFormatter returns a formatter for variable-line datetime fields.
// Expects []time.Time slice.
func createDateTimeNFormatter() OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *uicontext.UiContext) any {
		if value == nil {
			return []string{}
		}
		arr, ok := value.([]time.Time)
		if !ok {
			return []string{}
		}

		loc, err := time.LoadLocation(ctx.Timezone.GetIANA())
		if err != nil {
			loc = time.UTC
		}

		switch output {
		case OutputWeb, OutputPDF:
			result := make([]string, len(arr))
			for i, v := range arr {
				if !v.IsZero() {
					result[i] = formatter.FormatTimestampDateTime(v.Unix(), ctx)
				}
			}
			return result
		case OutputCSV, OutputExcel:
			strs := make([]string, len(arr))
			for i, v := range arr {
				if !v.IsZero() {
					strs[i] = v.In(loc).Format("2006-01-02 15:04:05")
				}
			}
			return strs
		}
		result := make([]string, len(arr))
		for i, v := range arr {
			if !v.IsZero() {
				result[i] = v.In(loc).Format("2006-01-02 15:04:05")
			}
		}
		return result
	})
}

// createDateNFormatter returns a formatter for variable-line date fields.
// Expects []time.Time slice.
func createDateNFormatter() OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *uicontext.UiContext) any {
		if value == nil {
			return []string{}
		}
		arr, ok := value.([]time.Time)
		if !ok {
			return []string{}
		}

		loc, err := time.LoadLocation(ctx.Timezone.GetIANA())
		if err != nil {
			loc = time.UTC
		}

		switch output {
		case OutputWeb, OutputPDF:
			result := make([]string, len(arr))
			for i, v := range arr {
				if !v.IsZero() {
					result[i] = formatter.FormatTimestampDate(v.Unix(), ctx)
				}
			}
			return result
		case OutputCSV, OutputExcel:
			strs := make([]string, len(arr))
			for i, v := range arr {
				if !v.IsZero() {
					strs[i] = v.In(loc).Format("2006-01-02")
				}
			}
			return strs
		}
		result := make([]string, len(arr))
		for i, v := range arr {
			if !v.IsZero() {
				result[i] = v.In(loc).Format("2006-01-02")
			}
		}
		return result
	})
}

// createDistanceNFormatter returns a formatter for variable-line distance fields.
// Expects []float64 slice (values in kilometers).
func createDistanceNFormatter(decimals int) OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *uicontext.UiContext) any {
		if value == nil {
			return []string{}
		}
		arr, ok := value.([]float64)
		if !ok {
			return []string{}
		}

		distanceUnit := ctx.Distance

		switch output {
		case OutputWeb, OutputPDF:
			result := make([]string, len(arr))
			for i, v := range arr {
				result[i] = formatter.FormatDistanceLocaleWithDecimals(v, distanceUnit, ctx.Locale, decimals)
			}
			return result
		case OutputCSV, OutputExcel:
			format := "%." + strconv.Itoa(decimals) + "f"
			strs := make([]string, len(arr))
			for i, v := range arr {
				strs[i] = fmt.Sprintf(format, convertDistanceValue(v, distanceUnit))
			}
			return strs
		}
		format := "%." + strconv.Itoa(decimals) + "f"
		result := make([]string, len(arr))
		for i, v := range arr {
			result[i] = fmt.Sprintf(format, convertDistanceValue(v, distanceUnit))
		}
		return result
	})
}

// createSpeedNFormatter returns a formatter for variable-line speed fields.
// Expects []float64 slice (values in km/h).
func createSpeedNFormatter(decimals int) OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *uicontext.UiContext) any {
		if value == nil {
			return []string{}
		}
		arr, ok := value.([]float64)
		if !ok {
			return []string{}
		}

		distanceUnit := ctx.Distance

		unit := " km/h"
		switch distanceUnit {
		case distance.Miles:
			unit = " mph"
		case distance.Seemiles:
			unit = " kn"
		}

		switch output {
		case OutputWeb, OutputPDF:
			result := make([]string, len(arr))
			for i, v := range arr {
				result[i] = formatter.FormatNumberLocale(convertSpeedValue(v, distanceUnit), decimals, ctx.Locale) + unit
			}
			return result
		case OutputCSV, OutputExcel:
			format := "%." + strconv.Itoa(decimals) + "f"
			strs := make([]string, len(arr))
			for i, v := range arr {
				strs[i] = fmt.Sprintf(format, convertSpeedValue(v, distanceUnit))
			}
			return strs
		}
		format := "%." + strconv.Itoa(decimals) + "f"
		result := make([]string, len(arr))
		for i, v := range arr {
			result[i] = fmt.Sprintf(format, convertSpeedValue(v, distanceUnit))
		}
		return result
	})
}

// createBoolNFormatter returns a formatter for variable-line boolean fields.
// Expects []bool slice.
func createBoolNFormatter() OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *uicontext.UiContext) any {
		if value == nil {
			return []string{}
		}
		arr, ok := value.([]bool)
		if !ok {
			return []string{}
		}

		result := make([]string, len(arr))
		for i, v := range arr {
			if v {
				result[i] = "Yes"
			} else {
				result[i] = "No"
			}
		}

		switch output {
		case OutputCSV, OutputExcel:
			return result
		}
		return result
	})
}

// createTimeLengthNFormatter returns a formatter for variable-line time duration fields.
// Expects []int64 slice (values in seconds).
func createTimeLengthNFormatter() OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *uicontext.UiContext) any {
		if value == nil {
			return []string{}
		}
		arr, ok := value.([]int64)
		if !ok {
			return []string{}
		}

		switch output {
		case OutputWeb, OutputPDF:
			result := make([]string, len(arr))
			for i, v := range arr {
				result[i] = formatTimeLength(v)
			}
			return result
		case OutputCSV, OutputExcel:
			strs := make([]string, len(arr))
			for i, v := range arr {
				strs[i] = fmt.Sprintf("%d", v/60)
			}
			return strs
		}
		result := make([]string, len(arr))
		for i, v := range arr {
			result[i] = fmt.Sprintf("%d", v/60)
		}
		return result
	})
}

// formatTimeLength formats seconds as "HH:MM" or "Xd HH:MM"
func formatTimeLength(seconds int64) string {
	if seconds < 0 {
		return ""
	}

	if seconds == 0 {
		return "00:00"
	}

	totalMinutes := seconds / 60
	hours := totalMinutes / 60
	minutes := totalMinutes % 60

	if hours >= 24 {
		days := hours / 24
		hours = hours % 24
		return fmt.Sprintf("%dd %02d:%02d", days, hours, minutes)
	}

	return fmt.Sprintf("%02d:%02d", hours, minutes)
}

// createLinkFormatter returns a formatter for link fields.
// Expects [2]string array: [0] = display text, [1] = URL.
// For OutputWeb: Returns [2]string array (GetData() will split into two fields).
// For OutputPDF/OutputCSV/OutputExcel: Returns just the display text.
func createLinkFormatter() OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *uicontext.UiContext) any {
		if value == nil {
			return [2]string{"", ""}
		}
		arr, ok := value.([2]string)
		if !ok {
			return [2]string{"", ""}
		}
		if output == OutputPDF || output == OutputCSV || output == OutputExcel {
			return arr[0]
		}
		return arr
	})
}

func createBoolFormatter(trueText, falseText string) OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *uicontext.UiContext) any {
		b, ok := value.(bool)
		if !ok {
			return ""
		}
		if b {
			return trueText
		}
		return falseText
	})
}

func createDateTimeFormatter() OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *uicontext.UiContext) any {
		timestamp := toInt64(value)
		if timestamp == 0 {
			return ""
		}
		switch output {
		case OutputWeb, OutputPDF:
			return formatter.FormatTimestampDateTime(timestamp, ctx)
		case OutputCSV, OutputExcel:
			loc, err := time.LoadLocation(ctx.Timezone.GetIANA())
			if err != nil {
				loc = time.UTC
			}
			t := time.Unix(timestamp, 0).In(loc)
			return t.Format("2006-01-02 15:04:05")
		}
		loc, err := time.LoadLocation(ctx.Timezone.GetIANA())
		if err != nil {
			loc = time.UTC
		}
		t := time.Unix(timestamp, 0).In(loc)
		return t.Format("2006-01-02 15:04:05")
	})
}

func createDateFormatter() OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *uicontext.UiContext) any {
		timestamp := toInt64(value)
		if timestamp == 0 {
			return ""
		}
		switch output {
		case OutputWeb, OutputPDF:
			return formatter.FormatTimestampDate(timestamp, ctx)
		case OutputCSV, OutputExcel:
			loc, err := time.LoadLocation(ctx.Timezone.GetIANA())
			if err != nil {
				loc = time.UTC
			}
			t := time.Unix(timestamp, 0).In(loc)
			return t.Format("2006-01-02")
		}
		loc, err := time.LoadLocation(ctx.Timezone.GetIANA())
		if err != nil {
			loc = time.UTC
		}
		t := time.Unix(timestamp, 0).In(loc)
		return t.Format("2006-01-02")
	})
}

func createDistanceFormatter(decimals int) OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *uicontext.UiContext) any {
		km := toFloat64(value)
		distanceUnit := ctx.Distance

		switch output {
		case OutputWeb, OutputPDF:
			return formatter.FormatDistanceLocaleWithDecimals(km, distanceUnit, ctx.Locale, decimals)
		case OutputCSV, OutputExcel:
			converted := convertDistanceValue(km, distanceUnit)
			format := "%." + strconv.Itoa(decimals) + "f"
			return fmt.Sprintf(format, converted)
		}
		converted := convertDistanceValue(km, distanceUnit)
		format := "%." + strconv.Itoa(decimals) + "f"
		return fmt.Sprintf(format, converted)
	})
}

func createPressureFormatter(decimals int) OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *uicontext.UiContext) any {
		bar := toFloat64(value)
		pressureUnit := ctx.Pressure

		switch output {
		case OutputWeb, OutputPDF:
			converted := convertPressureValue(bar, pressureUnit)
			unit := " bar"
			if pressureUnit == pressure.Psi {
				unit = " psi"
			}
			return formatter.FormatNumberLocale(converted, decimals, ctx.Locale) + unit
		case OutputCSV, OutputExcel:
			converted := convertPressureValue(bar, pressureUnit)
			format := "%." + strconv.Itoa(decimals) + "f"
			return fmt.Sprintf(format, converted)
		}
		converted := convertPressureValue(bar, pressureUnit)
		format := "%." + strconv.Itoa(decimals) + "f"
		return fmt.Sprintf(format, converted)
	})
}

func createSpeedFormatter(decimals int) OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *uicontext.UiContext) any {
		kmh := toFloat64(value)
		distanceUnit := ctx.Distance // Speed uses distance unit

		switch output {
		case OutputWeb, OutputPDF:
			converted := convertSpeedValue(kmh, distanceUnit)
			unit := " km/h"
			switch distanceUnit {
			case distance.Miles:
				unit = " mph"
			case distance.Seemiles:
				unit = " kn"
			}
			return formatter.FormatNumberLocale(converted, decimals, ctx.Locale) + unit
		case OutputCSV, OutputExcel:
			converted := convertSpeedValue(kmh, distanceUnit)
			format := "%." + strconv.Itoa(decimals) + "f"
			return fmt.Sprintf(format, converted)
		}
		converted := convertSpeedValue(kmh, distanceUnit)
		format := "%." + strconv.Itoa(decimals) + "f"
		return fmt.Sprintf(format, converted)
	})
}

// ============================================================================
// Helper Functions
// ============================================================================

// convertDistanceValue converts km to the target distance unit
func convertDistanceValue(km float64, unit distance.Distance) float64 {
	switch unit {
	case distance.Miles:
		return km * 0.621371
	case distance.Seemiles:
		return km * 0.539957
	default:
		return km
	}
}

// convertPressureValue converts bar to the target pressure unit
func convertPressureValue(bar float64, unit pressure.Pressure) float64 {
	switch unit {
	case pressure.Psi:
		return bar * 14.5038
	default:
		return bar
	}
}

// convertSpeedValue converts km/h to the target speed unit
func convertSpeedValue(kmh float64, unit distance.Distance) float64 {
	switch unit {
	case distance.Miles:
		return kmh * 0.621371
	case distance.Seemiles:
		return kmh * 0.539957
	default:
		return kmh
	}
}

// toInt64 converts any numeric value to int64
func toInt64(value any) int64 {
	if value == nil {
		return 0
	}
	switch v := value.(type) {
	case int64:
		return v
	case int32:
		return int64(v)
	case int:
		return int64(v)
	case float64:
		return int64(v)
	case float32:
		return int64(v)
	}
	return 0
}
