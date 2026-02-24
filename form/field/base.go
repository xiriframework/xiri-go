package field

import (
	"github.com/xiriframework/xiri-go/uicontext"
)

// FieldType represents the type of form field
type FieldType string

const (
	FieldTypeTimeRange  FieldType = "timerange"
	FieldTypeModelList  FieldType = "modellist"
	FieldTypeModel      FieldType = "model"
	FieldTypeBool       FieldType = "bool"
	FieldTypeInt        FieldType = "number"
	FieldTypeSelect     FieldType = "select"
	FieldTypeText       FieldType = "text"
	FieldTypeDeviceList FieldType = "devicelist"
	FieldTypeDriverList FieldType = "driverlist"
	FieldTypeArray      FieldType = "array"
	FieldTypeFile       FieldType = "file"
	FieldTypeHeader     FieldType = "header"
	FieldTypeHtml       FieldType = "html"
	FieldTypeInfo       FieldType = "info"
	FieldTypeJson       FieldType = "json"
	FieldTypeSerial     FieldType = "serial"
	FieldTypeTime       FieldType = "time"
	FieldTypeGeoform    FieldType = "geoform"   // Geometry field for geofencing/map drawing (16th field type)
	FieldTypeTimelimit  FieldType = "timelimit" // Time limit field with weekdays and time range (17th field type)
	FieldTypeChips      FieldType = "chips"     // Tag/chip input field (18th field type)
	FieldTypeDivider    FieldType = "divider"   // Visual divider/separator (19th field type)
)

// FormField is the base interface that all form fields must implement
type FormField interface {
	// GetID returns the field's unique identifier
	GetID() string

	// GetType returns the field type
	GetType() FieldType

	// GetName returns the translation key for the field label
	GetName() string

	// IsRequired returns whether this field is required
	IsRequired() bool

	// Validate validates a field value
	Validate(value interface{}) error

	// Parse parses a raw value into the expected type
	Parse(raw interface{}) (interface{}, error)

	// GetDefault returns the default value for this field
	GetDefault() interface{}

	// GetForm returns whether this field should be shown in the form
	GetForm() bool

	// ExportForFrontend exports the field definition for the frontend
	// The ctx parameter is used for translations and the value parameter
	// is the current field value (from form data)
	ExportForFrontend(ctx *uicontext.UiContext, value interface{}) map[string]interface{}
}

// ValueBinder is an interface for fields that support binding a raw value.
// All field types implement this implicitly via their BindValue method.
type ValueBinder interface {
	BindValue(raw interface{}) error
}

// FieldOptionsLoader is an interface for fields that load their options dynamically
// (e.g., Model, ModelList fields). The loader implementation is project-specific.
type FieldOptionsLoader interface {
	LoadOptions(ctx *uicontext.UiContext) error
}

// BaseField provides common functionality for all form fields
type BaseField struct {
	ID       string
	Type     FieldType
	Name     string // Translation key
	Required bool
	Default  interface{}

	// Display options
	Hint     string // Tooltip/help text for the field
	Class    string // CSS class for frontend styling (e.g., "xcol-md-6")
	Step     int    // Step indicator for multi-step forms (0 = no step)
	Disabled bool   // Whether the field is disabled

	// Conditional visibility
	ShowWhen []Condition // Conditions that must all be true for the field to be visible

	// Advanced options
	Access   []string // Access control permissions (nil = no restriction)
	Scenario []string // Which scenarios this field applies to (nil = all scenarios)
	DBName   string   // Database column name (empty = use ID)
	Form     bool     // Whether to show in form (false = hidden field)
}

func (f *BaseField) GetID() string {
	return f.ID
}

func (f *BaseField) GetType() FieldType {
	return f.Type
}

func (f *BaseField) GetName() string {
	return f.Name
}

func (f *BaseField) IsRequired() bool {
	return f.Required
}

func (f *BaseField) GetDefault() interface{} {
	return f.Default
}

func (f *BaseField) GetForm() bool {
	return f.Form
}

// GetBaseExport returns common export fields for all field types
func (f *BaseField) GetBaseExport(ctx *uicontext.UiContext, value interface{}) map[string]interface{} {
	result := map[string]interface{}{
		"id":       f.ID,
		"type":     string(f.Type),
		"required": f.Required,
		"step":     f.Step,
		"class":    f.Class,
		"value":    value,
		"form":     f.Form,
	}

	// Translate name if context available
	if ctx != nil && ctx.Translate != nil {
		result["name"] = ctx.Translate(f.Name)
		if f.Hint != "" {
			result["hint"] = ctx.Translate(f.Hint)
		} else {
			result["hint"] = nil
		}
	} else {
		result["name"] = f.Name
		if f.Hint != "" {
			result["hint"] = f.Hint
		} else {
			result["hint"] = nil
		}
	}

	// Add showWhen conditions if any
	if len(f.ShowWhen) > 0 {
		if len(f.ShowWhen) == 1 {
			result["showWhen"] = f.ShowWhen[0]
		} else {
			result["showWhen"] = f.ShowWhen
		}
	}

	return result
}

// ============================================================================
// Builder Pattern - Chainable Setter Methods
// ============================================================================

// SetClass sets the CSS class for frontend styling (e.g., "xcol-md-6")
func (f *BaseField) SetClass(class string) *BaseField {
	f.Class = class
	return f
}

// SetHint sets the tooltip/help text for the field
func (f *BaseField) SetHint(hint string) *BaseField {
	f.Hint = hint
	return f
}

// SetStep sets the step indicator for multi-step forms (0 = no step)
func (f *BaseField) SetStep(step int) *BaseField {
	f.Step = step
	return f
}

// SetDisabled sets whether the field is disabled
func (f *BaseField) SetDisabled(disabled bool) *BaseField {
	f.Disabled = disabled
	return f
}

// SetAccess sets the access control permissions (nil = no restriction)
func (f *BaseField) SetAccess(access []string) *BaseField {
	f.Access = access
	return f
}

// SetScenario sets which scenarios this field applies to (nil = all scenarios)
func (f *BaseField) SetScenario(scenario []string) *BaseField {
	f.Scenario = scenario
	return f
}

// SetForm sets whether to show in form (false = hidden field)
func (f *BaseField) SetForm(form bool) *BaseField {
	f.Form = form
	return f
}

// SetShowWhen adds a visibility condition for this field
func (f *BaseField) SetShowWhen(field string, operator ConditionOperator, value interface{}) *BaseField {
	f.ShowWhen = append(f.ShowWhen, NewCondition(field, operator, value))
	return f
}

// SetShowWhenNotEmpty adds a "notEmpty" visibility condition for this field
func (f *BaseField) SetShowWhenNotEmpty(field string) *BaseField {
	f.ShowWhen = append(f.ShowWhen, NewConditionNotEmpty(field))
	return f
}
