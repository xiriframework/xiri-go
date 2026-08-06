package group

// ReloadDependent is implemented by fields whose content depends on other field values.
// Optional interface - fields opt in by embedding BaseField and calling SetReloadOn.
type ReloadDependent interface {
	GetReloadOn() []string
	GetReloadURL() string
}

// ExportPatch exports only the fields that declare a reload dependency, keyed by field ID.
//
// This is the response body for a reload request: the frontend merges each entry into the
// matching field definition. Value resolution matches ExportForFrontendWithValues(nil) -
// the frontend keeps its own control value and only prunes what the new option list no
// longer contains.
func (fg *FormGroup) ExportPatch() map[string]interface{} {
	result := make(map[string]interface{})

	for _, f := range fg.fields {
		if !f.GetForm() {
			continue
		}

		// Beide Angaben, wie beim normalen Export: ein halb deklariertes Feld exportiert dort
		// keine Reload-Keys, das Frontend würde einen Patch für es also gar nicht anwenden.
		dep, ok := f.(ReloadDependent)
		if !ok || len(dep.GetReloadOn()) == 0 || dep.GetReloadURL() == "" {
			continue
		}

		result[f.GetID()] = f.ExportForFrontend(fg.ctx, f.GetDefault())
	}

	return result
}

// ExportForFrontend exports field definitions in frontend-compatible format
// with translated names and locale-specific formatting
func (fg *FormGroup) ExportForFrontend() []map[string]interface{} {
	return fg.ExportForFrontendWithValues(nil)
}

// ExportForFrontendWithValues exports field definitions with specific values
// If values is nil, uses default values from fields
func (fg *FormGroup) ExportForFrontendWithValues(values map[string]interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(fg.fields))

	for _, f := range fg.fields {
		if !f.GetForm() {
			continue
		}

		// Get value for this field (from values map or default)
		var value interface{}
		if values != nil {
			if v, exists := values[f.GetID()]; exists {
				value = v
			}
		}
		if value == nil {
			value = f.GetDefault()
		}

		// Call field's ExportForFrontend method - this handles all field-specific logic
		def := f.ExportForFrontend(fg.ctx, value)

		result = append(result, def)
	}

	return result
}
