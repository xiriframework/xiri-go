package field

import (
	"fmt"
	"slices"

	"github.com/xiriframework/xiri-go/component/core"
)

// ModelListField represents a multiple model/object selection form field
// This is used to select multiple objects from a list
type ModelListField struct {
	*BaseField
	ModelType   string                 // Type of model (e.g., "device", "driver", "group")
	URL         string                 // API endpoint to fetch options (for frontend)
	List        []ModelOption          // Predefined list of options
	Filter      map[string]interface{} // Additional filter parameters
	Add         []ModelOption          // Additional options to add to the list
	Sub         []int32                // Subtract/remove specific IDs from list
	MinItems    *int                   // Minimum number of selections
	MaxItems    *int                   // Maximum number of selections
	Params      map[string]interface{} // Additional parameters for API call
	LoaderFunc  ModelLoaderFunc        // Function to load options from database
	AllowedFunc ModelAllowedFunc       // Authorizes chosen IDs server-side (needed with URL)
	Options     []ModelOption          // Loaded options (populated by LoadOptions)
	AllowEmpty  bool                   // If true, empty selection is allowed
	SingleOnly  bool                   // If true, only one item can be selected
	Tree        bool                   // If true, export type as "treeselect" instead of "multiselect"
	Value       ModelListValue         // Parsed and validated value (type-safe access)
}

// LoadOptions implements FieldOptionsLoader interface
// Loads model options using the configured LoaderFunc
func (f *ModelListField) LoadOptions(ctx *core.UiContext) error {
	if f.LoaderFunc == nil {
		// No loader function = options are static or loaded externally
		return nil
	}

	options, err := f.LoaderFunc(ctx, f.ModelType)
	if err != nil {
		return fmt.Errorf("loading options for modellist field %s: %w", f.ID, err)
	}

	f.Options = options
	return nil
}

// SetLoaderFunc sets the function used to load options from database
func (f *ModelListField) SetLoaderFunc(loader ModelLoaderFunc) {
	f.LoaderFunc = loader
}

// SetAllowedFunc sets the server-side authorization for chosen IDs (see ModelAllowedFunc).
// It is called once per Validate with all IDs that are neither in Sub nor in the unchanged default.
func (f *ModelListField) SetAllowedFunc(allowed ModelAllowedFunc) {
	f.AllowedFunc = allowed
}

func (f *ModelListField) Validate(value interface{}) error {
	if value == nil {
		if f.Required && !f.AllowEmpty {
			return fmt.Errorf("modellist field %s is required", f.ID)
		}
		return nil
	}

	list, ok := value.(ModelListValue)
	if !ok {
		return fmt.Errorf("invalid modellist value type for %s", f.ID)
	}

	if !f.AllowEmpty && len(list) == 0 && f.Required {
		return fmt.Errorf("modellist %s cannot be empty", f.ID)
	}

	if f.MinItems != nil && len(list) < *f.MinItems {
		return fmt.Errorf("modellist %s must have at least %d items", f.ID, *f.MinItems)
	}

	if f.MaxItems != nil && len(list) > *f.MaxItems {
		return fmt.Errorf("modellist %s must have at most %d items", f.ID, *f.MaxItems)
	}

	if f.SingleOnly && len(list) > 1 {
		return fmt.Errorf("modellist %s can only have one item", f.ID)
	}

	def, _ := parseModelListValue(nil, f.GetDefault())
	var check []int32
	for _, id := range list {
		if slices.Contains(f.Sub, id) {
			return fmt.Errorf("modellist field %s: id %d is not an allowed option", f.ID, id)
		}
		if slices.Contains(def, id) {
			continue
		}
		if f.AllowedFunc == nil && !modelIDAllowed(f.URL, f.LoaderFunc != nil, f.Options, f.List, id) {
			return fmt.Errorf("modellist field %s: id %d is not an allowed option", f.ID, id)
		}
		check = append(check, id)
	}
	// The hook only says yes or no for the whole batch, so the message names no ID.
	if f.AllowedFunc != nil && len(check) > 0 && !f.AllowedFunc(check) {
		return fmt.Errorf("modellist field %s: contains an id that is not an allowed option", f.ID)
	}

	return nil
}

func (f *ModelListField) Parse(raw interface{}) (interface{}, error) {
	// Reuse the helper function from types.go
	return parseModelListValue(raw, f.GetDefault())
}

// BindValue parses, validates, and stores the value in the field
func (f *ModelListField) BindValue(raw interface{}) error {
	parsed, err := f.Parse(raw)
	if err != nil {
		return fmt.Errorf("parsing field %s: %w", f.ID, err)
	}

	if err := f.Validate(parsed); err != nil {
		return fmt.Errorf("validating field %s: %w", f.ID, err)
	}

	if parsed != nil {
		if list, ok := parsed.(ModelListValue); ok {
			f.Value = list
		}
	} else {
		f.Value = nil
	}

	return nil
}

// ============================================================================
// Builder Functions
// ============================================================================

// NewModelListField creates a modellist (multi-selection) form field
// Note: LoaderFunc must be set separately via SetLoaderFunc for project-specific loading.
// currentValue can be nil (defaults to empty slice)
func NewModelListField(id, name string, required bool, modelType string, currentValue []int32) *ModelListField {
	// Default to empty slice if nil
	if currentValue == nil {
		currentValue = []int32{}
	}

	return &ModelListField{
		BaseField: &BaseField{
			ID:       id,
			Type:     FieldTypeModelList,
			Name:     name,
			Required: required,
			Default:  ModelListValue(currentValue),
			Form:     true,
		},
		ModelType: modelType,
	}
}

// ExportForFrontend exports the field for frontend rendering
func (f *ModelListField) ExportForFrontend(ctx *core.UiContext, value interface{}) map[string]interface{} {
	if value == nil {
		value = f.GetDefault()
	}
	result := f.BaseField.GetBaseExport(ctx, value)

	// Set type based on tree mode
	if f.Tree {
		result["type"] = "treeselect"
	} else {
		result["type"] = "multiselect"
	}

	if f.Tree && f.URL == "" {
		fmt.Printf("WARNING: ModelListField %q has Tree=true but no URL set. Tree data requires a URL endpoint that returns nested children.\n", f.ID)
	}

	// Export options as list of {id, name} (loaded from LoadOptions or List field)
	options := f.Options
	if len(options) == 0 && len(f.List) > 0 {
		options = f.List
	}

	// Convert options to frontend format
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

	// Add min/max if specified
	if f.MinItems != nil {
		result["min"] = *f.MinItems
	}
	if f.MaxItems != nil {
		result["max"] = *f.MaxItems
	}

	return result
}

// ============================================================================
// Chainable Setter Methods
// ============================================================================

// SetClass sets the CSS class for frontend styling
func (f *ModelListField) SetClass(class string) *ModelListField {
	f.BaseField.SetClass(class)
	return f
}

// SetHint sets the tooltip/help text for the field
func (f *ModelListField) SetHint(hint string) *ModelListField {
	f.BaseField.SetHint(hint)
	return f
}

// SetStep sets the step indicator for multi-step forms
func (f *ModelListField) SetStep(step int) *ModelListField {
	f.BaseField.SetStep(step)
	return f
}

// SetDisabled sets whether the field is disabled
func (f *ModelListField) SetDisabled(disabled bool) *ModelListField {
	f.BaseField.SetDisabled(disabled)
	return f
}

// SetAccess sets the access control permissions
func (f *ModelListField) SetAccess(access []string) *ModelListField {
	f.BaseField.SetAccess(access)
	return f
}

// SetScenario sets which scenarios this field applies to
func (f *ModelListField) SetScenario(scenario []string) *ModelListField {
	f.BaseField.SetScenario(scenario)
	return f
}

// SetForm sets whether to show in form
func (f *ModelListField) SetForm(form bool) *ModelListField {
	f.BaseField.SetForm(form)
	return f
}

// SetTree sets whether this field renders as a tree select.
// When true, a URL must be set that returns data with nested children arrays.
func (f *ModelListField) SetTree(tree bool) *ModelListField {
	f.Tree = tree
	return f
}
