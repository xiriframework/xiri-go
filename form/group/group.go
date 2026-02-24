package group

import (
	"fmt"

	"github.com/xiriframework/xiri-go/form/field"
	"github.com/xiriframework/xiri-go/uicontext"
)

// FormGroup represents a collection of form fields
type FormGroup struct {
	fields []field.FormField
	index  map[string]field.FormField // For quick lookup by ID
	ctx    *uicontext.UiContext
}

// NewFormGroup creates a new form group without context
// For backwards compatibility - use NewFormGroupWithContext for new code
func NewFormGroup(fields []field.FormField) *FormGroup {
	index := make(map[string]field.FormField)
	for _, f := range fields {
		index[f.GetID()] = f
	}

	return &FormGroup{
		fields: fields,
		index:  index,
		ctx:    nil,
	}
}

// NewFormGroupWithContext creates a new form group with user context
func NewFormGroupWithContext(fields []field.FormField, ctx *uicontext.UiContext) (*FormGroup, error) {
	index := make(map[string]field.FormField)
	for _, f := range fields {
		index[f.GetID()] = f
	}

	fg := &FormGroup{
		fields: fields,
		index:  index,
		ctx:    ctx,
	}

	// Auto-load field options for Model/ModelList fields
	if err := fg.LoadFieldOptions(); err != nil {
		return nil, fmt.Errorf("loading field options: %w", err)
	}

	return fg, nil
}

// SetContext sets the user context
// Also triggers loading of field options for Model/ModelList fields
func (fg *FormGroup) SetContext(ctx *uicontext.UiContext) error {
	fg.ctx = ctx

	// Auto-load field options when context is set
	return fg.LoadFieldOptions()
}

// LoadFieldOptions loads dynamic options for fields that implement FieldOptionsLoader
// This is called automatically when setting context, but can also be called manually
func (fg *FormGroup) LoadFieldOptions() error {
	if fg.ctx == nil {
		// No context = skip loading
		return nil
	}

	for _, f := range fg.fields {
		if loader, ok := f.(field.FieldOptionsLoader); ok {
			if err := loader.LoadOptions(fg.ctx); err != nil {
				return fmt.Errorf("loading options for field %s: %w", f.GetID(), err)
			}
		}
	}

	return nil
}

// GetFields returns all fields in the group
func (fg *FormGroup) GetFields() []field.FormField {
	return fg.fields
}

// GetField returns a specific field by ID
func (fg *FormGroup) GetField(id string) (field.FormField, bool) {
	f, exists := fg.index[id]
	return f, exists
}

// GetDefaults returns a map of all default values
func (fg *FormGroup) GetDefaults() map[string]interface{} {
	defaults := make(map[string]interface{})

	for _, f := range fg.fields {
		if def := f.GetDefault(); def != nil {
			defaults[f.GetID()] = def
		}
	}

	return defaults
}

// GetRequiredFields returns all required fields
func (fg *FormGroup) GetRequiredFields() []field.FormField {
	var required []field.FormField
	for _, f := range fg.fields {
		if f.IsRequired() {
			required = append(required, f)
		}
	}
	return required
}

// GetOptionalFields returns all optional fields
func (fg *FormGroup) GetOptionalFields() []field.FormField {
	var optional []field.FormField
	for _, f := range fg.fields {
		if !f.IsRequired() {
			optional = append(optional, f)
		}
	}
	return optional
}

// GetTranslatedName returns the translated name for a field
// If translation is not available, returns the original name/key
func (fg *FormGroup) GetTranslatedName(fieldID string) string {
	f, exists := fg.index[fieldID]
	if !exists {
		return fieldID
	}

	name := f.GetName()
	if fg.ctx != nil && fg.ctx.Translate != nil {
		return fg.ctx.Translate(name)
	}

	return name
}

// GetContext returns the user context
func (fg *FormGroup) GetContext() *uicontext.UiContext {
	return fg.ctx
}

// HasContext returns true if the form group has a user context
func (fg *FormGroup) HasContext() bool {
	return fg.ctx != nil
}

// GetFieldIDs returns a slice of all field IDs in the form group.
// This is used by BindAndValidate to safely extract only declared fields from requests.
func (fg *FormGroup) GetFieldIDs() []string {
	ids := make([]string, len(fg.fields))
	for i, f := range fg.fields {
		ids[i] = f.GetID()
	}
	return ids
}
