package field

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/xiriframework/xiri-go/component/core"
)

// FileField represents a file upload form field.
//
// xiri-ng reads the selected files as data URLs and submits them inside the normal form JSON
// as [{name, data}], data being "data:<mime>;base64,<payload>". There is no multipart upload.
type FileField struct {
	*BaseField
	MaxSize           int64    // Maximum size per file in bytes (0 = no limit); enforced by the browser and by Validate
	AllowedTypes      []string // Allowed MIME types (e.g., ["image/jpeg", "image/png"])
	AllowedExtensions []string // Allowed file extensions (e.g., [".jpg", ".png"])
	Multiple          bool     // If true, multiple files can be uploaded
}

// Validate checks the submitted [{name, data}] list: Required (nil or empty list fails) and
// MaxSize per file, measured on the decoded base64 payload. The MIME type is not checked:
// it comes from the client inside the data URL.
func (f *FileField) Validate(value interface{}) error {
	if value == nil {
		if f.Required {
			return invalid("required", nil, "file field %s is required", f.ID)
		}
		return nil
	}
	files, ok := value.([]interface{})
	if !ok {
		return fmt.Errorf("invalid file value type for %s", f.ID)
	}
	if len(files) == 0 && f.Required {
		return invalid("required", nil, "file field %s is required", f.ID)
	}
	for i, item := range files {
		file, ok := item.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid file entry at index %d for %s", i, f.ID)
		}
		name, _ := file["name"].(string)
		data, _ := file["data"].(string)
		_, payload, found := strings.Cut(data, ";base64,")
		if name == "" || !strings.HasPrefix(data, "data:") || !found {
			return fmt.Errorf("file %d of %s is no {name, data} entry with a base64 data URL", i, f.ID)
		}
		// DecodedLen overestimates by at most 2 padding bytes: reject oversized payloads
		// before allocating a decode buffer for them.
		if f.MaxSize > 0 && int64(base64.StdEncoding.DecodedLen(len(payload))) > f.MaxSize+2 {
			return invalid("max_size", map[string]any{"max": f.MaxSize}, "file %d of %s exceeds %d bytes", i, f.ID, f.MaxSize)
		}
		decoded, err := base64.StdEncoding.DecodeString(payload)
		if err != nil {
			return fmt.Errorf("file %d of %s has invalid base64 data: %w", i, f.ID, err)
		}
		if f.MaxSize > 0 && int64(len(decoded)) > f.MaxSize {
			return invalid("max_size", map[string]any{"max": f.MaxSize}, "file %d of %s exceeds %d bytes", i, f.ID, f.MaxSize)
		}
	}
	return nil
}

// Parse returns the raw value as-is: the JSON-decoded [{name, data}] list (see FileField).
// The app decodes the data URLs itself after Validate.
func (f *FileField) Parse(raw interface{}) (interface{}, error) {
	if raw == nil {
		return f.GetDefault(), nil
	}
	return raw, nil
}

// ============================================================================
// Builder Functions
// ============================================================================

// NewFileField creates a file upload form field. maxSize is the limit per file in bytes
// (0 = no limit); the browser rejects larger files and Validate checks it again on the server.
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
func (f *FileField) ExportForFrontend(ctx *core.UiContext, value interface{}) map[string]interface{} {
	if value == nil {
		value = f.GetDefault()
	}
	result := f.BaseField.GetBaseExport(ctx, value)

	// Set type to file
	result["type"] = "file"

	// xiri-ng reads "max" (bytes per file) and "accept" (the <input type=file> attribute)
	if f.MaxSize > 0 {
		result["max"] = f.MaxSize
	}
	if accept := append(append([]string{}, f.AllowedTypes...), f.AllowedExtensions...); len(accept) > 0 {
		result["accept"] = strings.Join(accept, ",")
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

// SetAccess stores role metadata. Metadata only: neither exported nor evaluated by xiri-go, no access control.
func (f *FileField) SetAccess(access []string) *FileField {
	f.BaseField.SetAccess(access)
	return f
}

// SetScenario stores scenario metadata. Metadata only: neither exported nor evaluated by xiri-go.
func (f *FileField) SetScenario(scenario []string) *FileField {
	f.BaseField.SetScenario(scenario)
	return f
}

// SetForm sets whether to show in form
func (f *FileField) SetForm(form bool) *FileField {
	f.BaseField.SetForm(form)
	return f
}
