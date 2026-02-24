package group

import (
	"fmt"
)

// ParseValues parses raw field values into typed values
func (fg *FormGroup) ParseValues(raw map[string]interface{}) (map[string]interface{}, error) {
	parsed := make(map[string]interface{})

	// Parse each field value
	for _, f := range fg.fields {
		rawValue, exists := raw[f.GetID()]

		// If value doesn't exist, use default
		if !exists {
			if f.IsRequired() {
				return nil, fmt.Errorf("required field %s is missing", f.GetID())
			}
			// Use default value
			if def := f.GetDefault(); def != nil {
				parsed[f.GetID()] = def
			}
			continue
		}

		// Parse the value
		value, err := f.Parse(rawValue)
		if err != nil {
			return nil, fmt.Errorf("error parsing field %s: %w", f.GetID(), err)
		}

		parsed[f.GetID()] = value
	}

	return parsed, nil
}

// ValidateValues validates parsed field values
func (fg *FormGroup) ValidateValues(values map[string]interface{}) error {
	// Validate each field
	for _, f := range fg.fields {
		value, exists := values[f.GetID()]

		// Check required fields
		if !exists && f.IsRequired() {
			return fmt.Errorf("required field %s is missing", f.GetID())
		}

		// Validate the value
		if err := f.Validate(value); err != nil {
			return err
		}
	}

	return nil
}

// ParseAndValidate is a convenience method that parses and validates in one call
func (fg *FormGroup) ParseAndValidate(raw map[string]interface{}) (map[string]interface{}, error) {
	// Parse values
	parsed, err := fg.ParseValues(raw)
	if err != nil {
		return nil, err
	}

	// Validate parsed values
	if err := fg.ValidateValues(parsed); err != nil {
		return nil, err
	}

	return parsed, nil
}
