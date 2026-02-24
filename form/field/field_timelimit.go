package field

import (
	"fmt"

	"github.com/xiriframework/xiri-go/uicontext"
)

// TimeLimitField represents a time limit form field with weekdays and time range
// This is a "multi-field" that maps to 5 database columns: time_check, time_weekdays, time_from, time_to, time_in
type TimeLimitField struct {
	*BaseField
}

// TimeLimitValue represents the value of a timelimit field
// Format: {check, wd, fromhour, frommin, tohour, tomin, in}
type TimeLimitValue struct {
	Check    bool    // Whether time limit is enabled
	Weekdays [7]bool // Weekdays (0=Sunday, 1=Monday, ..., 6=Saturday)
	FromHour string  // Start hour (00-23) as string for frontend
	FromMin  string  // Start minute (00-55) as string for frontend
	ToHour   string  // End hour (00-24) as string for frontend
	ToMin    string  // End minute (00-55) as string for frontend
	In       bool    // If true, applies INSIDE time range; if false, OUTSIDE time range
}

// NewTimeLimitField creates a timelimit form field
// Default: Check disabled, all weekdays false, 00:00-24:00, inside=true
func NewTimeLimitField(id, name string, required bool) *TimeLimitField {
	return &TimeLimitField{
		BaseField: &BaseField{
			ID:       id,
			Type:     FieldTypeTimelimit,
			Name:     name,
			Required: required,
			Default: TimeLimitValue{
				Check:    false,
				Weekdays: [7]bool{false, false, false, false, false, false, false},
				FromHour: "00",
				FromMin:  "00",
				ToHour:   "24",
				ToMin:    "00",
				In:       true,
			},
			Form: true,
		},
	}
}

func (f *TimeLimitField) Validate(value interface{}) error {
	if value == nil {
		if f.Required {
			return fmt.Errorf("timelimit field %s is required", f.ID)
		}
		return nil
	}

	tl, ok := value.(TimeLimitValue)
	if !ok {
		return fmt.Errorf("invalid timelimit value type for %s", f.ID)
	}

	// Validate hour/minute ranges if check is enabled
	if tl.Check {
		// Validate from hour/min
		fromHour := parseHourMin(tl.FromHour)
		fromMin := parseHourMin(tl.FromMin)
		toHour := parseHourMin(tl.ToHour)
		toMin := parseHourMin(tl.ToMin)

		if fromHour < 0 || fromHour > 23 {
			return fmt.Errorf("timelimit %s: fromhour must be 0-23", f.ID)
		}
		if fromMin < 0 || fromMin > 55 {
			return fmt.Errorf("timelimit %s: frommin must be 0-55", f.ID)
		}
		if toHour < 0 || toHour > 24 {
			return fmt.Errorf("timelimit %s: tohour must be 0-24", f.ID)
		}
		if toMin < 0 || toMin > 55 {
			return fmt.Errorf("timelimit %s: tomin must be 0-55", f.ID)
		}
	}

	return nil
}

func (f *TimeLimitField) Parse(raw interface{}) (interface{}, error) {
	if raw == nil {
		return f.GetDefault(), nil
	}

	// Expect map[string]interface{} from JSON
	data, ok := raw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("timelimit field %s: expected object, got %T", f.ID, raw)
	}

	// Parse check
	check := false
	if checkVal, ok := data["check"]; ok {
		if checkBool, ok := checkVal.(bool); ok {
			check = checkBool
		}
	}

	// Parse weekdays array [7]bool
	weekdays := [7]bool{false, false, false, false, false, false, false}
	if wdVal, ok := data["wd"]; ok {
		if wdArr, ok := wdVal.([]interface{}); ok {
			for i := 0; i < 7 && i < len(wdArr); i++ {
				if wdBool, ok := wdArr[i].(bool); ok {
					weekdays[i] = wdBool
				}
			}
		}
	}

	// Parse time components (as strings for frontend)
	fromHour := "00"
	fromMin := "00"
	toHour := "24"
	toMin := "00"
	inRange := true

	if val, ok := data["fromhour"]; ok {
		fromHour = fmt.Sprintf("%v", val)
	}
	if val, ok := data["frommin"]; ok {
		fromMin = fmt.Sprintf("%v", val)
	}
	if val, ok := data["tohour"]; ok {
		toHour = fmt.Sprintf("%v", val)
	}
	if val, ok := data["tomin"]; ok {
		toMin = fmt.Sprintf("%v", val)
	}
	if val, ok := data["in"]; ok {
		if inBool, ok := val.(bool); ok {
			inRange = inBool
		}
	}

	return TimeLimitValue{
		Check:    check,
		Weekdays: weekdays,
		FromHour: fromHour,
		FromMin:  fromMin,
		ToHour:   toHour,
		ToMin:    toMin,
		In:       inRange,
	}, nil
}

// ExportForFrontend exports the field for frontend rendering
func (f *TimeLimitField) ExportForFrontend(ctx *uicontext.UiContext, value interface{}) map[string]interface{} {
	if value == nil {
		value = f.GetDefault()
	}

	result := f.BaseField.GetBaseExport(ctx, value)

	// Type and subtype
	result["type"] = "timelimit"
	result["subtype"] = "timelimit"

	// Translation texts for frontend labels
	t := ctx.Translate
	result["texts"] = map[string]interface{}{
		"check":    t("TL.CHECK"),
		"weekdays": t("TL.WEEKDAYS"),
		"wd0":      t("TL.WD0"), // Sunday
		"wd1":      t("TL.WD1"), // Monday
		"wd2":      t("TL.WD2"), // Tuesday
		"wd3":      t("TL.WD3"), // Wednesday
		"wd4":      t("TL.WD4"), // Thursday
		"wd5":      t("TL.WD5"), // Friday
		"wd6":      t("TL.WD6"), // Saturday
		"from":     t("TL.FROM"),
		"to":       t("TL.TO"),
		"inout":    t("TL.INOUT"),
	}

	// Convert value to frontend format
	tl, ok := value.(TimeLimitValue)
	if !ok {
		// Use default if value is wrong type
		tl = f.GetDefault().(TimeLimitValue)
	}

	// Export value: {check, wd, fromhour, frommin, tohour, tomin, in}
	result["value"] = map[string]interface{}{
		"check":    tl.Check,
		"wd":       tl.Weekdays,
		"fromhour": tl.FromHour,
		"frommin":  tl.FromMin,
		"tohour":   tl.ToHour,
		"tomin":    tl.ToMin,
		"in":       tl.In,
	}

	return result
}

// ToMultiDb converts the timelimit value to database format (5 columns)
// Returns: {cols: [...], names: [...], data: {...}}
func (f *TimeLimitField) ToMultiDb(value interface{}) map[string]interface{} {
	tl, ok := value.(TimeLimitValue)
	if !ok {
		// Use default if value is wrong type
		tl = f.GetDefault().(TimeLimitValue)
	}

	// Convert weekdays [7]bool to PostgreSQL boolean array format "{t,f,t,f,t,f,f}"
	weekdaysDB := "{"
	for i, wd := range tl.Weekdays {
		if i > 0 {
			weekdaysDB += ","
		}
		if wd {
			weekdaysDB += "t"
		} else {
			weekdaysDB += "f"
		}
	}
	weekdaysDB += "}"

	// Format time as HH:MM
	timeFrom := fmt.Sprintf("%s:%s", tl.FromHour, tl.FromMin)
	timeTo := fmt.Sprintf("%s:%s", tl.ToHour, tl.ToMin)

	// Convert to DB format
	return map[string]interface{}{
		"cols":  []string{"time_check", "time_weekdays", "time_from", "time_to", "time_in"},
		"names": []string{":time_check", ":time_weekdays", ":time_from", ":time_to", ":time_in"},
		"data": map[string]interface{}{
			"time_check":    tl.Check,
			"time_weekdays": weekdaysDB,
			"time_from":     timeFrom,
			"time_to":       timeTo,
			"time_in":       tl.In,
		},
	}
}

// parseHourMin converts string to int for validation
func parseHourMin(s string) int {
	var val int
	_, _ = fmt.Sscanf(s, "%d", &val)
	return val
}

// ============================================================================
// Chainable Setter Methods
// ============================================================================

// SetClass sets the CSS class for frontend styling
func (f *TimeLimitField) SetClass(class string) *TimeLimitField {
	f.BaseField.SetClass(class)
	return f
}

// SetHint sets the tooltip/help text for the field
func (f *TimeLimitField) SetHint(hint string) *TimeLimitField {
	f.BaseField.SetHint(hint)
	return f
}

// SetStep sets the step indicator for multi-step forms
func (f *TimeLimitField) SetStep(step int) *TimeLimitField {
	f.BaseField.SetStep(step)
	return f
}

// SetDisabled sets whether the field is disabled
func (f *TimeLimitField) SetDisabled(disabled bool) *TimeLimitField {
	f.BaseField.SetDisabled(disabled)
	return f
}

// SetAccess sets the access control permissions
func (f *TimeLimitField) SetAccess(access []string) *TimeLimitField {
	f.BaseField.SetAccess(access)
	return f
}

// SetScenario sets which scenarios this field applies to
func (f *TimeLimitField) SetScenario(scenario []string) *TimeLimitField {
	f.BaseField.SetScenario(scenario)
	return f
}

// SetForm sets whether to show in form
func (f *TimeLimitField) SetForm(form bool) *TimeLimitField {
	f.BaseField.SetForm(form)
	return f
}
