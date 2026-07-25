package field

import (
	"fmt"
	"time"

	"github.com/xiriframework/xiri-go/component/core"
)

// TimeRangeField represents a date/time range form field
type TimeRangeField struct {
	*BaseField
	AllowSingleDay bool            // If true, end date can equal start date
	Subtype        string          // "daterange" or "time" (default: "time")
	Min            *int64          // Minimum date (Unix timestamp) - can be days offset or absolute timestamp
	Max            *int64          // Maximum date (Unix timestamp) - can be days offset or absolute timestamp
	Value          *TimeRangeValue // Parsed and validated value (type-safe access)
}

// TimeRangeValue represents a parsed time range
type TimeRangeValue struct {
	Start time.Time
	End   time.Time
}

func (f *TimeRangeField) Validate(value interface{}) error {
	if value == nil {
		if f.Required {
			return fmt.Errorf("timerange field %s is required", f.ID)
		}
		return nil
	}

	tr, ok := value.(*TimeRangeValue)
	if !ok {
		return fmt.Errorf("invalid timerange value type for %s", f.ID)
	}

	if tr.Start.IsZero() || tr.End.IsZero() {
		return fmt.Errorf("timerange %s cannot have zero dates", f.ID)
	}

	if f.AllowSingleDay {
		if tr.Start.After(tr.End) {
			return fmt.Errorf("timerange %s start cannot be after end", f.ID)
		}
	} else {
		if !tr.Start.Before(tr.End) {
			return fmt.Errorf("timerange %s start must be before end", f.ID)
		}
	}

	return nil
}

func (f *TimeRangeField) Parse(raw interface{}) (interface{}, error) {
	if raw == nil {
		return f.GetDefault(), nil
	}

	// Expect map with "start" and "end" keys (Unix timestamps or ISO strings)
	m, ok := raw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("timerange field expects map with start/end")
	}

	// Parse start date (support Unix timestamp or ISO string)
	start, err := parseDateTime(m["start"])
	if err != nil {
		return nil, fmt.Errorf("timerange field missing or invalid start date: %w", err)
	}

	// Parse end date (support Unix timestamp or ISO string)
	end, err := parseDateTime(m["end"])
	if err != nil {
		return nil, fmt.Errorf("timerange field missing or invalid end date: %w", err)
	}

	return &TimeRangeValue{Start: start, End: end}, nil
}

// BindValue parses, validates, and stores the value in the field
func (f *TimeRangeField) BindValue(raw interface{}) error {
	parsed, err := f.Parse(raw)
	if err != nil {
		return fmt.Errorf("parsing field %s: %w", f.ID, err)
	}

	if err := f.Validate(parsed); err != nil {
		return fmt.Errorf("validating field %s: %w", f.ID, err)
	}

	if parsed != nil {
		if tr, ok := parsed.(*TimeRangeValue); ok {
			f.Value = tr
		}
	} else {
		f.Value = nil
	}

	return nil
}

// ExportForFrontend exports the field for frontend rendering
func (f *TimeRangeField) ExportForFrontend(ctx *core.UiContext, value interface{}) map[string]interface{} {
	// Get base export fields
	result := f.BaseField.GetBaseExport(ctx, nil)

	// Determine type based on subtype
	if f.Subtype == "daterange" {
		result["type"] = "daterange"
	} else {
		result["type"] = "datetimerange"
	}
	result["subtype"] = f.Subtype

	now := time.Now()

	// Value takes precedence, then the field default, then now.
	var start, end time.Time
	if tr, ok := value.(*TimeRangeValue); ok && tr != nil {
		start = tr.Start
		end = tr.End
	}
	if start.IsZero() || end.IsZero() {
		if tr, ok := f.GetDefault().(*TimeRangeValue); ok && tr != nil {
			start = tr.Start
			end = tr.End
		}
	}
	if start.IsZero() || end.IsZero() {
		start = now
		end = now
	}

	// Export as Unix timestamps (seconds)
	result["value"] = map[string]int64{
		"start": start.Unix(),
		"end":   end.Unix(),
	}

	// Add min/max if specified — day offsets are anchored to local midnight
	// in the user's timezone (same as TimeField).
	if f.Min != nil {
		result["min"] = resolveDateBound(ctx, *f.Min, now)
	}
	if f.Max != nil {
		result["max"] = resolveDateBound(ctx, *f.Max, now)
	}

	return result
}

// ============================================================================
// Builder Functions
// ============================================================================

// NewTimeRangeField creates a new time range form field
func NewTimeRangeField(id, name string, required bool) *TimeRangeField {
	return &TimeRangeField{
		BaseField: &BaseField{
			ID:       id,
			Type:     FieldTypeTimeRange,
			Name:     name,
			Required: required,
			Form:     true,
		},
		AllowSingleDay: true,
	}
}

// NewTimeRangeFieldWithDefault creates a time range field with a default value
// defaultDays specifies how many days back from now the default range should start
func NewTimeRangeFieldWithDefault(id, name string, required bool, defaultDays int) *TimeRangeField {
	now := time.Now()
	start := now.AddDate(0, 0, -defaultDays)
	end := now

	return &TimeRangeField{
		BaseField: &BaseField{
			ID:       id,
			Type:     FieldTypeTimeRange,
			Name:     name,
			Required: required,
			Default: &TimeRangeValue{
				Start: start,
				End:   end,
			},
			Form: true,
		},
		AllowSingleDay: true,
	}
}

// ============================================================================
// Chainable Setter Methods
// ============================================================================

// SetClass sets the CSS class for frontend styling
func (f *TimeRangeField) SetClass(class string) *TimeRangeField {
	f.BaseField.SetClass(class)
	return f
}

// SetHint sets the tooltip/help text for the field
func (f *TimeRangeField) SetHint(hint string) *TimeRangeField {
	f.BaseField.SetHint(hint)
	return f
}

// SetStep sets the step indicator for multi-step forms
func (f *TimeRangeField) SetStep(step int) *TimeRangeField {
	f.BaseField.SetStep(step)
	return f
}

// SetDisabled sets whether the field is disabled
func (f *TimeRangeField) SetDisabled(disabled bool) *TimeRangeField {
	f.BaseField.SetDisabled(disabled)
	return f
}

// SetAccess sets the access control permissions
func (f *TimeRangeField) SetAccess(access []string) *TimeRangeField {
	f.BaseField.SetAccess(access)
	return f
}

// SetScenario sets which scenarios this field applies to
func (f *TimeRangeField) SetScenario(scenario []string) *TimeRangeField {
	f.BaseField.SetScenario(scenario)
	return f
}

// SetForm sets whether to show in form
func (f *TimeRangeField) SetForm(form bool) *TimeRangeField {
	f.BaseField.SetForm(form)
	return f
}
