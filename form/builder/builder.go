// Package formhelper provides form building and request binding utilities.
package builder

import (
	"github.com/xiriframework/xiri-go/form/field"
	"github.com/xiriframework/xiri-go/form/group"
	"github.com/xiriframework/xiri-go/uicontext"
)

// FormBuilder provides a fluent API for building forms
// that work seamlessly across Add/Edit/AddSave/EditSave workflows.
//
// KEY FEATURE: After BindAndValidate(), access values via field.Value (type-safe, no assertions!)
type FormBuilder struct {
	fields    []field.FormField
	ctx       *uicontext.UiContext
	translate func(string) string

	// Optional hook for edit edge cases (e.g., checking if current value is in options list)
	OnEditValueCheck func(fg *group.FormGroup, values map[string]interface{}) error
}

// NewFormBuilder creates a new FormBuilder.
//
// Parameters:
//   - ctx: User context containing locale, timezone, etc.
//   - translate: Translation function for field labels
//
// Returns a FormBuilder that can be used to define form fields via method chaining.
func NewFormBuilder(ctx *uicontext.UiContext, translate func(string) string) *FormBuilder {
	return &FormBuilder{
		fields:    make([]field.FormField, 0),
		ctx:       ctx,
		translate: translate,
	}
}

// BuildAdd builds a FormGroup for Add workflow.
// Returns the FormGroup and a map of default values.
//
// Default values are taken from field definitions (defaultValue parameters).
// Use this for displaying empty forms and for validation in AddSave.
func (fb *FormBuilder) BuildAdd() (*group.FormGroup, map[string]interface{}, error) {
	// Create FormGroup with context
	fg, err := group.NewFormGroupWithContext(fb.fields, fb.ctx)
	if err != nil {
		return nil, nil, err
	}

	// Collect default values from fields
	defaults := make(map[string]interface{})
	for _, f := range fb.fields {
		defaults[f.GetID()] = f.GetDefault()
	}

	return fg, defaults, nil
}

// BuildEdit builds a FormGroup for Edit workflow.
// Returns the FormGroup and a map of current values from field defaults.
//
// If OnEditValueCheck is set, it will be called to handle edge cases
// (e.g., current value not in accessible options list).
func (fb *FormBuilder) BuildEdit() (*group.FormGroup, map[string]interface{}, error) {
	// Create FormGroup with context
	fg, err := group.NewFormGroupWithContext(fb.fields, fb.ctx)
	if err != nil {
		return nil, nil, err
	}

	// Collect current values from field defaults
	values := make(map[string]interface{})
	for _, f := range fb.fields {
		values[f.GetID()] = f.GetDefault()
	}

	// Handle edge case via optional hook
	if fb.OnEditValueCheck != nil {
		if err := fb.OnEditValueCheck(fg, values); err != nil {
			return nil, nil, err
		}
	}

	return fg, values, nil
}

// AddField adds a field to the builder.
func (fb *FormBuilder) AddField(field field.FormField) *FormBuilder {
	fb.fields = append(fb.fields, field)
	return fb
}

// BuildAddForDisplay builds and exports field definitions for Add workflow.
// Returns frontend-ready JSON with default values.
func (fb *FormBuilder) BuildAddForDisplay() ([]map[string]interface{}, error) {
	fg, defaults, err := fb.BuildAdd()
	if err != nil {
		return nil, err
	}
	return fg.ExportForFrontendWithValues(defaults), nil
}

// BuildEditForDisplay builds and exports field definitions for Edit workflow.
// Returns frontend-ready JSON with current values from field defaults.
func (fb *FormBuilder) BuildEditForDisplay() ([]map[string]interface{}, error) {
	fg, values, err := fb.BuildEdit()
	if err != nil {
		return nil, err
	}
	return fg.ExportForFrontendWithValues(values), nil
}
