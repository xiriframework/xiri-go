package field

import (
	"fmt"
	"time"

	"github.com/xiriframework/xiri-go/uicontext"
)

// TimeField represents a single date/time form field
type TimeField struct {
	*BaseField
	Format      string // Date format string (e.g., "2006-01-02 15:04:05")
	MinDate     *time.Time
	MaxDate     *time.Time
	Min         *int64 // Minimum date (Unix timestamp or days offset)
	Max         *int64 // Maximum date (Unix timestamp or days offset)
	AllowPast   bool   // If false, only future dates are allowed
	AllowFuture bool   // If false, only past dates are allowed
	Subtype     string // Subtype: "date", "datetime", "time"
	Value       *int64 // Parsed and validated value (Unix timestamp)
}

func (f *TimeField) Validate(value interface{}) error {
	if value == nil {
		if f.Required {
			return fmt.Errorf("time field %s is required", f.ID)
		}
		return nil
	}

	// Accept Unix timestamp (int, int64) or time.Time
	var t time.Time
	switch v := value.(type) {
	case int:
		t = time.Unix(int64(v), 0)
	case int64:
		t = time.Unix(v, 0)
	case time.Time:
		t = v
	default:
		return fmt.Errorf("invalid time value type for %s", f.ID)
	}

	if t.IsZero() {
		return fmt.Errorf("time %s cannot be zero", f.ID)
	}

	now := time.Now()
	if !f.AllowPast && t.Before(now) {
		return fmt.Errorf("time %s cannot be in the past", f.ID)
	}

	if !f.AllowFuture && t.After(now) {
		return fmt.Errorf("time %s cannot be in the future", f.ID)
	}

	if f.MinDate != nil && t.Before(*f.MinDate) {
		return fmt.Errorf("time %s is before minimum date", f.ID)
	}

	if f.MaxDate != nil && t.After(*f.MaxDate) {
		return fmt.Errorf("time %s is after maximum date", f.ID)
	}

	return nil
}

func (f *TimeField) Parse(raw interface{}) (interface{}, error) {
	if raw == nil {
		return f.GetDefault(), nil
	}

	// Parse various time formats
	switch v := raw.(type) {
	case string:
		// Try ISO date format first
		if t, err := time.Parse("2006-01-02", v); err == nil {
			return t.Unix(), nil
		}
		// Try ISO datetime format
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			return t.Unix(), nil
		}
		// Try custom format if specified
		if f.Format != "" {
			if t, err := time.Parse(f.Format, v); err == nil {
				return t.Unix(), nil
			}
		}
		return nil, fmt.Errorf("invalid date string format: %s", v)

	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case float64:
		return int64(v), nil
	case time.Time:
		return v.Unix(), nil

	default:
		return nil, fmt.Errorf("unsupported time type: %T", v)
	}
}

// BindValue parses, validates, and stores the value in the field
// This enables type-safe access via field.Value after form binding
func (f *TimeField) BindValue(raw interface{}) error {
	parsed, err := f.Parse(raw)
	if err != nil {
		return fmt.Errorf("parsing field %s: %w", f.ID, err)
	}

	if err := f.Validate(parsed); err != nil {
		return fmt.Errorf("validating field %s: %w", f.ID, err)
	}

	if parsed != nil {
		if timestamp, ok := parsed.(int64); ok {
			f.Value = &timestamp
		}
	} else {
		f.Value = nil
	}

	return nil
}

// ============================================================================
// Builder Functions
// ============================================================================

// NewTimeField creates a single time/datetime form field
func NewTimeField(id, name string, required bool, defaultValue int64) *TimeField {
	return &TimeField{
		BaseField: &BaseField{
			ID:       id,
			Type:     FieldTypeTime,
			Name:     name,
			Required: required,
			Default:  defaultValue,
			Form:     true,
		},
		AllowPast:   true,
		AllowFuture: true,
	}
}

// ExportForFrontend exports the field for frontend rendering
func (f *TimeField) ExportForFrontend(ctx *uicontext.UiContext, value interface{}) map[string]interface{} {
	if value == nil {
		value = f.GetDefault()
	}
	result := f.BaseField.GetBaseExport(ctx, value)

	// Determine type based on subtype
	if f.Subtype == "date" {
		result["type"] = "date"
	} else if f.Subtype == "time" {
		result["type"] = "time"
	} else {
		result["type"] = "datetime"
	}
	result["subtype"] = f.Subtype

	// Add min/max with day offset handling (same as TimeRangeField)
	// Days offset is calculated from midnight in user's timezone
	if f.Min != nil {
		minVal := *f.Min
		// If value is between -10000 and 10000, treat as days offset from midnight today
		if minVal > -10000 && minVal < 10000 {
			// Get user's timezone
			loc, err := time.LoadLocation(ctx.Timezone.GetIANA())
			if err != nil {
				loc = time.UTC
			}
			// Calculate midnight today in user's timezone, then add days offset
			// This ensures DST transitions are handled correctly
			now := time.Now().In(loc)
			midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
			targetDate := midnight.AddDate(0, 0, int(minVal))
			result["min"] = targetDate.Unix()
		} else {
			result["min"] = minVal
		}
	}
	if f.Max != nil {
		maxVal := *f.Max
		// If value is between -10000 and 10000, treat as days offset from midnight today
		if maxVal > -10000 && maxVal < 10000 {
			// Get user's timezone
			loc, err := time.LoadLocation(ctx.Timezone.GetIANA())
			if err != nil {
				loc = time.UTC
			}
			// Calculate midnight today in user's timezone, then add days offset
			// This ensures DST transitions are handled correctly
			now := time.Now().In(loc)
			midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
			targetDate := midnight.AddDate(0, 0, int(maxVal))
			result["max"] = targetDate.Unix()
		} else {
			result["max"] = maxVal
		}
	}

	return result
}

// ============================================================================
// Chainable Setter Methods
// ============================================================================

// SetClass sets the CSS class for frontend styling
func (f *TimeField) SetClass(class string) *TimeField {
	f.BaseField.SetClass(class)
	return f
}

// SetHint sets the tooltip/help text for the field
func (f *TimeField) SetHint(hint string) *TimeField {
	f.BaseField.SetHint(hint)
	return f
}

// SetStep sets the step indicator for multi-step forms
func (f *TimeField) SetStep(step int) *TimeField {
	f.BaseField.SetStep(step)
	return f
}

// SetDisabled sets whether the field is disabled
func (f *TimeField) SetDisabled(disabled bool) *TimeField {
	f.BaseField.SetDisabled(disabled)
	return f
}

// SetAccess sets the access control permissions
func (f *TimeField) SetAccess(access []string) *TimeField {
	f.BaseField.SetAccess(access)
	return f
}

// SetScenario sets which scenarios this field applies to
func (f *TimeField) SetScenario(scenario []string) *TimeField {
	f.BaseField.SetScenario(scenario)
	return f
}

// SetForm sets whether to show in form
func (f *TimeField) SetForm(form bool) *TimeField {
	f.BaseField.SetForm(form)
	return f
}
