package table

import (
	"fmt"
	"strconv"
	"time"

	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/formatter"
	"github.com/xiriframework/xiri-go/types/distance"
	"github.com/xiriframework/xiri-go/types/pressure"
)

// ============================================================================
// Formatter creation functions
// These create the actual formatter functions that will be used at runtime
// ============================================================================

func createIdFormatter() OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *core.UiContext) any {
		return toInt64(value)
	})
}

func createIntegerFormatter() OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *core.UiContext) any {
		num := toInt64(value)
		if output == OutputWeb || output == OutputPDF {
			return formatter.FormatNumberLocale(float64(num), 0, ctx.SafeLocale())
		}
		return fmt.Sprint(num)
	})
}

func createFloatFormatter(decimals int) OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *core.UiContext) any {
		num := toFloat64(value)
		if output == OutputWeb || output == OutputPDF {
			return formatter.FormatNumberLocale(num, decimals, ctx.SafeLocale())
		}
		format := "%." + strconv.Itoa(decimals) + "f"
		return fmt.Sprintf(format, num)
	})
}

func createTextFormatter() OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *core.UiContext) any {
		if value == nil {
			return ""
		}
		return fmt.Sprint(value)
	})
}

// createChipsFormatter returns a formatter for chips-type fields.
// Expects []Chip. For web output emits []map[string]any{{"label", "color"}, ...}
// (matching the Angular xiri-table chips renderer). For CSV/Excel/PDF joins labels with ", ".
func createChipsFormatter() OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *core.UiContext) any {
		chips, ok := value.([]Chip)
		if !ok || chips == nil {
			if output == OutputCSV || output == OutputExcel || output == OutputPDF {
				return ""
			}
			return []map[string]any{}
		}
		if output == OutputCSV || output == OutputExcel || output == OutputPDF {
			labels := make([]string, len(chips))
			for i, c := range chips {
				labels[i] = c.Label
			}
			out := ""
			for i, l := range labels {
				if i > 0 {
					out += ", "
				}
				out += l
			}
			return out
		}
		out := make([]map[string]any, len(chips))
		for i, c := range chips {
			out[i] = map[string]any{
				"label": c.Label,
				"color": string(c.Color),
			}
		}
		return out
	})
}

// createPassthroughFormatter returns a formatter that preserves complex data structures (maps, arrays)
// for JSON serialization. For CSV output, converts to string representation.
func createPassthroughFormatter() OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *core.UiContext) any {
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
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *core.UiContext) any {
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
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *core.UiContext) any {
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
				formatter.FormatNumberLocale(float64(arr[0]), 0, ctx.SafeLocale()),
				formatter.FormatNumberLocale(float64(arr[1]), 0, ctx.SafeLocale()),
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
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *core.UiContext) any {
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
				formatter.FormatNumberLocale(arr[0], decimals, ctx.SafeLocale()),
				formatter.FormatNumberLocale(arr[1], decimals, ctx.SafeLocale()),
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
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *core.UiContext) any {
		if value == nil {
			return [2]string{"", ""}
		}
		arr, ok := value.([2]time.Time)
		if !ok {
			return [2]string{"", ""}
		}

		loc, err := time.LoadLocation(ctx.SafeTimezone().GetIANA())
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
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *core.UiContext) any {
		if value == nil {
			return [2]string{"", ""}
		}
		arr, ok := value.([2]time.Time)
		if !ok {
			return [2]string{"", ""}
		}

		loc, err := time.LoadLocation(ctx.SafeTimezone().GetIANA())
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
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *core.UiContext) any {
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
				formatter.FormatDistanceLocaleWithDecimals(arr[0], distanceUnit, ctx.SafeLocale(), decimals),
				formatter.FormatDistanceLocaleWithDecimals(arr[1], distanceUnit, ctx.SafeLocale(), decimals),
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
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *core.UiContext) any {
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
			primary := formatter.FormatNumberLocale(convertDistanceValue(arr[0], distanceUnit), decimals, ctx.SafeLocale()) + unit
			secondary := formatter.FormatNumberLocale(convertDistanceValue(arr[1], distanceUnit), decimals, ctx.SafeLocale()) + unit
			return [2]string{primary, secondary}
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

// createText2BoolFormatter returns a formatter for two-line boolean fields.
// Expects [2]bool array: [0] = primary value, [1] = secondary value
func createText2BoolFormatter() OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *core.UiContext) any {
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
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *core.UiContext) any {
		seconds := toInt64(value)

		if output == OutputWeb || output == OutputPDF {
			return formatTimeLength(seconds)
		}
		return fmt.Sprintf("%d", seconds/60)
	})
}

// createText2TimeLengthFormatter returns a formatter for two-line time duration fields.
// Expects [2]int64 array (values in seconds): [0] = primary, [1] = secondary
func createText2TimeLengthFormatter() OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *core.UiContext) any {
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
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *core.UiContext) any {
		if value == nil {
			return []string{}
		}
		arr, ok := value.([]string)
		if !ok {
			return []string{}
		}
		return arr
	})
}

// createIntegerNFormatter returns a formatter for variable-line integer fields.
// Expects []int slice.
func createIntegerNFormatter() OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *core.UiContext) any {
		if value == nil {
			return []string{}
		}
		arr, ok := value.([]int)
		if !ok {
			return []string{}
		}
		if output == OutputWeb || output == OutputPDF {
			result := make([]string, len(arr))
			for i, v := range arr {
				result[i] = formatter.FormatNumberLocale(float64(v), 0, ctx.SafeLocale())
			}
			return result
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
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *core.UiContext) any {
		if value == nil {
			return []string{}
		}
		arr, ok := value.([]float64)
		if !ok {
			return []string{}
		}
		format := "%." + strconv.Itoa(decimals) + "f"
		if output == OutputWeb || output == OutputPDF {
			result := make([]string, len(arr))
			for i, v := range arr {
				result[i] = formatter.FormatNumberLocale(v, decimals, ctx.SafeLocale())
			}
			return result
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
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *core.UiContext) any {
		if value == nil {
			return []string{}
		}
		arr, ok := value.([]time.Time)
		if !ok {
			return []string{}
		}

		loc, err := time.LoadLocation(ctx.SafeTimezone().GetIANA())
		if err != nil {
			loc = time.UTC
		}

		if output == OutputWeb || output == OutputPDF {
			result := make([]string, len(arr))
			for i, v := range arr {
				if !v.IsZero() {
					result[i] = formatter.FormatTimestampDateTime(v.Unix(), ctx)
				}
			}
			return result
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
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *core.UiContext) any {
		if value == nil {
			return []string{}
		}
		arr, ok := value.([]time.Time)
		if !ok {
			return []string{}
		}

		loc, err := time.LoadLocation(ctx.SafeTimezone().GetIANA())
		if err != nil {
			loc = time.UTC
		}

		if output == OutputWeb || output == OutputPDF {
			result := make([]string, len(arr))
			for i, v := range arr {
				if !v.IsZero() {
					result[i] = formatter.FormatTimestampDate(v.Unix(), ctx)
				}
			}
			return result
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
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *core.UiContext) any {
		if value == nil {
			return []string{}
		}
		arr, ok := value.([]float64)
		if !ok {
			return []string{}
		}

		distanceUnit := ctx.Distance

		if output == OutputWeb || output == OutputPDF {
			result := make([]string, len(arr))
			for i, v := range arr {
				result[i] = formatter.FormatDistanceLocaleWithDecimals(v, distanceUnit, ctx.SafeLocale(), decimals)
			}
			return result
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
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *core.UiContext) any {
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

		if output == OutputWeb || output == OutputPDF {
			result := make([]string, len(arr))
			for i, v := range arr {
				result[i] = formatter.FormatNumberLocale(convertDistanceValue(v, distanceUnit), decimals, ctx.SafeLocale()) + unit
			}
			return result
		}
		format := "%." + strconv.Itoa(decimals) + "f"
		result := make([]string, len(arr))
		for i, v := range arr {
			result[i] = fmt.Sprintf(format, convertDistanceValue(v, distanceUnit))
		}
		return result
	})
}

// createBoolNFormatter returns a formatter for variable-line boolean fields.
// Expects []bool slice.
func createBoolNFormatter() OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *core.UiContext) any {
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

		return result
	})
}

// createTimeLengthNFormatter returns a formatter for variable-line time duration fields.
// Expects []int64 slice (values in seconds).
func createTimeLengthNFormatter() OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *core.UiContext) any {
		if value == nil {
			return []string{}
		}
		arr, ok := value.([]int64)
		if !ok {
			return []string{}
		}

		if output == OutputWeb || output == OutputPDF {
			result := make([]string, len(arr))
			for i, v := range arr {
				result[i] = formatTimeLength(v)
			}
			return result
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
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *core.UiContext) any {
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
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *core.UiContext) any {
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
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *core.UiContext) any {
		timestamp := toInt64(value)
		if timestamp == 0 {
			return ""
		}
		if output == OutputWeb || output == OutputPDF {
			return formatter.FormatTimestampDateTime(timestamp, ctx)
		}
		loc, err := time.LoadLocation(ctx.SafeTimezone().GetIANA())
		if err != nil {
			loc = time.UTC
		}
		t := time.Unix(timestamp, 0).In(loc)
		return t.Format("2006-01-02 15:04:05")
	})
}

func createDateFormatter() OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *core.UiContext) any {
		timestamp := toInt64(value)
		if timestamp == 0 {
			return ""
		}
		if output == OutputWeb || output == OutputPDF {
			return formatter.FormatTimestampDate(timestamp, ctx)
		}
		loc, err := time.LoadLocation(ctx.SafeTimezone().GetIANA())
		if err != nil {
			loc = time.UTC
		}
		t := time.Unix(timestamp, 0).In(loc)
		return t.Format("2006-01-02")
	})
}

func createDistanceFormatter(decimals int) OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *core.UiContext) any {
		km := toFloat64(value)
		distanceUnit := ctx.Distance

		if output == OutputWeb || output == OutputPDF {
			return formatter.FormatDistanceLocaleWithDecimals(km, distanceUnit, ctx.SafeLocale(), decimals)
		}
		converted := convertDistanceValue(km, distanceUnit)
		format := "%." + strconv.Itoa(decimals) + "f"
		return fmt.Sprintf(format, converted)
	})
}

func createPressureFormatter(decimals int) OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *core.UiContext) any {
		bar := toFloat64(value)
		pressureUnit := ctx.Pressure

		if output == OutputWeb || output == OutputPDF {
			converted := convertPressureValue(bar, pressureUnit)
			unit := " bar"
			if pressureUnit == pressure.Psi {
				unit = " psi"
			}
			return formatter.FormatNumberLocale(converted, decimals, ctx.SafeLocale()) + unit
		}
		converted := convertPressureValue(bar, pressureUnit)
		format := "%." + strconv.Itoa(decimals) + "f"
		return fmt.Sprintf(format, converted)
	})
}

func createSpeedFormatter(decimals int) OutputFormatter {
	return FormatterFunc(func(value any, row Row, output OutputType, ctx *core.UiContext) any {
		kmh := toFloat64(value)
		distanceUnit := ctx.Distance // Speed uses distance unit

		if output == OutputWeb || output == OutputPDF {
			converted := convertDistanceValue(kmh, distanceUnit)
			unit := " km/h"
			switch distanceUnit {
			case distance.Miles:
				unit = " mph"
			case distance.Seemiles:
				unit = " kn"
			}
			return formatter.FormatNumberLocale(converted, decimals, ctx.SafeLocale()) + unit
		}
		converted := convertDistanceValue(kmh, distanceUnit)
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
