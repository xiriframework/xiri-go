package field

import (
	"fmt"
	"strconv"

	"github.com/xiriframework/xiri-go/uicontext"
)

// ModelOption represents a single option in a model dropdown
type ModelOption struct {
	ID    int32                  // Model ID
	Name  string                 // Display name (already translated if needed)
	Extra map[string]interface{} // Additional metadata
}

// ModelLoaderFunc is a function type for loading model options.
// The loader implementation is project-specific (e.g., database-backed).
type ModelLoaderFunc func(ctx *uicontext.UiContext, modelType string) ([]ModelOption, error)

// ModelField represents a single model/object selection form field
// This is used to select a single object from a list (dropdown/autocomplete)
type ModelField struct {
	*BaseField
	ModelType   string                 // Type of model (e.g., "device", "driver", "group")
	URL         string                 // API endpoint to fetch options (for frontend)
	List        []ModelOption          // Predefined list of options
	Filter      map[string]interface{} // Additional filter parameters
	Add         []ModelOption          // Additional options to add to the list
	Sub         []int32                // Subtract/remove specific IDs from list
	AllowSearch bool                   // If true, search/autocomplete is enabled
	Params      map[string]interface{} // Additional parameters for API call
	LoaderFunc  ModelLoaderFunc        // Function to load options from database
	Options     []ModelOption          // Loaded options (populated by LoadOptions)
	Value       int32                  // Parsed and validated value (type-safe access - always has value)
}

// LoadOptions implements FieldOptionsLoader interface
// Loads model options using the configured LoaderFunc
func (f *ModelField) LoadOptions(ctx *uicontext.UiContext) error {
	if f.LoaderFunc == nil {
		// No loader function = options are static or loaded externally
		return nil
	}

	options, err := f.LoaderFunc(ctx, f.ModelType)
	if err != nil {
		return fmt.Errorf("loading options for model field %s: %w", f.ID, err)
	}

	f.Options = options
	return nil
}

// SetLoaderFunc sets the function used to load options from database
func (f *ModelField) SetLoaderFunc(loader ModelLoaderFunc) {
	f.LoaderFunc = loader
}

func (f *ModelField) Validate(value interface{}) error {
	if value == nil {
		if f.Required {
			return fmt.Errorf("model field %s is required", f.ID)
		}
		return nil
	}

	// Accept int, int32, int64 for model ID
	switch value.(type) {
	case int, int32, int64:
		return nil
	default:
		return fmt.Errorf("invalid model value type for %s, expected int", f.ID)
	}
}

func (f *ModelField) Parse(raw interface{}) (interface{}, error) {
	if raw == nil {
		return f.GetDefault(), nil
	}

	// Parse to int32 (model ID)
	switch v := raw.(type) {
	case int:
		return int32(v), nil
	case int32:
		return v, nil
	case int64:
		return int32(v), nil
	case float64:
		return int32(v), nil
	case string:
		parsed, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid model ID: %s", v)
		}
		return int32(parsed), nil
	default:
		return nil, fmt.Errorf("cannot parse model ID from %T", raw)
	}
}

// BindValue parses, validates, and stores the value in the field
func (f *ModelField) BindValue(raw interface{}) error {
	parsed, err := f.Parse(raw)
	if err != nil {
		return fmt.Errorf("parsing field %s: %w", f.ID, err)
	}

	if err := f.Validate(parsed); err != nil {
		return fmt.Errorf("validating field %s: %w", f.ID, err)
	}

	if parsed != nil {
		if id, ok := parsed.(int32); ok {
			f.Value = id
		}
	} else {
		// Use default value if parsed is nil
		if defaultVal, ok := f.GetDefault().(int32); ok {
			f.Value = defaultVal
		} else {
			f.Value = 0
		}
	}

	return nil
}

// ============================================================================
// Builder Functions
// ============================================================================

// NewModelField creates a model (single selection) form field
func NewModelField(id, name string, required bool, modelType string, currentValue int32) *ModelField {
	return &ModelField{
		BaseField: &BaseField{
			ID:       id,
			Type:     FieldTypeModel,
			Name:     name,
			Required: required,
			Default:  currentValue,
			Form:     true,
		},
		ModelType:   modelType,
		AllowSearch: true,
	}
}

// ExportForFrontend exports the field for frontend rendering
func (f *ModelField) ExportForFrontend(ctx *uicontext.UiContext, value interface{}) map[string]interface{} {
	if value == nil {
		value = f.GetDefault()
	}
	result := f.BaseField.GetBaseExport(ctx, value)

	// Set type to object (single selection)
	result["type"] = "object"

	// Export options as list of {id, name} (loaded from LoadOptions or List field)
	options := f.Options
	if len(options) == 0 && len(f.List) > 0 {
		options = f.List
	}

	// Convert options to frontend format, filtering out Sub IDs
	exportedOptions := make([]map[string]interface{}, 0, len(options))
	for _, opt := range options {
		// Skip if in Sub (subtract) list
		skip := false
		for _, subID := range f.Sub {
			if opt.ID == subID {
				skip = true
				break
			}
		}
		if !skip {
			exportedOptions = append(exportedOptions, map[string]interface{}{
				"id":   opt.ID,
				"name": opt.Name,
			})
		}
	}
	result["list"] = exportedOptions

	// Add URL and params if specified
	if f.URL != "" {
		result["url"] = f.URL
	}
	if f.Params != nil {
		result["params"] = f.Params
	}

	// Add search flag
	result["search"] = f.AllowSearch

	return result
}

// ============================================================================
// Chainable Setter Methods
// ============================================================================

// SetClass sets the CSS class for frontend styling
func (f *ModelField) SetClass(class string) *ModelField {
	f.BaseField.SetClass(class)
	return f
}

// SetHint sets the tooltip/help text for the field
func (f *ModelField) SetHint(hint string) *ModelField {
	f.BaseField.SetHint(hint)
	return f
}

// SetStep sets the step indicator for multi-step forms
func (f *ModelField) SetStep(step int) *ModelField {
	f.BaseField.SetStep(step)
	return f
}

// SetDisabled sets whether the field is disabled
func (f *ModelField) SetDisabled(disabled bool) *ModelField {
	f.BaseField.SetDisabled(disabled)
	return f
}

// SetAccess sets the access control permissions
func (f *ModelField) SetAccess(access []string) *ModelField {
	f.BaseField.SetAccess(access)
	return f
}

// SetScenario sets which scenarios this field applies to
func (f *ModelField) SetScenario(scenario []string) *ModelField {
	f.BaseField.SetScenario(scenario)
	return f
}

// SetForm sets whether to show in form
func (f *ModelField) SetForm(form bool) *ModelField {
	f.BaseField.SetForm(form)
	return f
}
