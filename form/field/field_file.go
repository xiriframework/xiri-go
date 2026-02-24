package field

import (
	"github.com/xiriframework/xiri-go/uicontext"
)

import "fmt"

// FileField represents a file upload form field
type FileField struct {
	*BaseField
	MaxSize           int64    // Maximum file size in bytes
	AllowedTypes      []string // Allowed MIME types (e.g., ["image/jpeg", "image/png"])
	AllowedExtensions []string // Allowed file extensions (e.g., [".jpg", ".png"])
	Multiple          bool     // If true, multiple files can be uploaded
}

// Validate checks whether the file field value meets constraints.
// File content validation (size, MIME type) happens at upload time via the HTTP handler,
// not here. This method only enforces the Required constraint.
func (f *FileField) Validate(value interface{}) error {
	if value == nil {
		if f.Required {
			return fmt.Errorf("file field %s is required", f.ID)
		}
		return nil
	}
	return nil
}

// Parse returns the raw value as-is.
// File fields are handled differently from other form fields: the actual file data
// is processed by the HTTP multipart handler, not by the form field parser.
// The value here is typically a filename or file metadata reference.
func (f *FileField) Parse(raw interface{}) (interface{}, error) {
	if raw == nil {
		return f.GetDefault(), nil
	}
	return raw, nil
}

// ============================================================================
// Builder Functions
// ============================================================================

// NewFileField creates a file upload form field
func NewFileField(id, name string, required bool, maxSize int64) *FileField {
	return &FileField{
		BaseField: &BaseField{
			ID:       id,
			Type:     FieldTypeFile,
			Name:     name,
			Required: required,
			Default:  nil,
			Form:     true,
		},
		MaxSize:  maxSize,
		Multiple: false,
	}
}

// ExportForFrontend exports the field for frontend rendering
func (f *FileField) ExportForFrontend(ctx *uicontext.UiContext, value interface{}) map[string]interface{} {
	if value == nil {
		value = f.GetDefault()
	}
	result := f.BaseField.GetBaseExport(ctx, value)

	// Set type to file
	result["type"] = "file"

	// Add maxSize if specified
	if f.MaxSize > 0 {
		result["maxSize"] = f.MaxSize
	}

	// Add allowedTypes if specified
	if len(f.AllowedTypes) > 0 {
		result["allowedTypes"] = f.AllowedTypes
	}

	// Add allowedExtensions if specified
	if len(f.AllowedExtensions) > 0 {
		result["allowedExtensions"] = f.AllowedExtensions
	}

	// Add multiple flag
	result["multiple"] = f.Multiple

	return result
}

// ============================================================================
// Chainable Setter Methods
// ============================================================================

// SetClass sets the CSS class for frontend styling
func (f *FileField) SetClass(class string) *FileField {
	f.BaseField.SetClass(class)
	return f
}

// SetHint sets the tooltip/help text for the field
func (f *FileField) SetHint(hint string) *FileField {
	f.BaseField.SetHint(hint)
	return f
}

// SetStep sets the step indicator for multi-step forms
func (f *FileField) SetStep(step int) *FileField {
	f.BaseField.SetStep(step)
	return f
}

// SetDisabled sets whether the field is disabled
func (f *FileField) SetDisabled(disabled bool) *FileField {
	f.BaseField.SetDisabled(disabled)
	return f
}

// SetAccess sets the access control permissions
func (f *FileField) SetAccess(access []string) *FileField {
	f.BaseField.SetAccess(access)
	return f
}

// SetScenario sets which scenarios this field applies to
func (f *FileField) SetScenario(scenario []string) *FileField {
	f.BaseField.SetScenario(scenario)
	return f
}

// SetForm sets whether to show in form
func (f *FileField) SetForm(form bool) *FileField {
	f.BaseField.SetForm(form)
	return f
}
