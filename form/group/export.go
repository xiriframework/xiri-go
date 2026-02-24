package group

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
