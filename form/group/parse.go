package group

import (
	"fmt"
)

// ParseValues parses raw field values into typed values. Missing, non-required
// fields receive their default value.
func (fg *FormGroup) ParseValues(raw map[string]interface{}) (map[string]interface{}, error) {
	return fg.parseValues(raw, true)
}

// ParseValuesSparse is ParseValues without defaults: a key that is missing from raw
// stays missing in the result. Used for filters, where "not sent" means "no filter".
func (fg *FormGroup) ParseValuesSparse(raw map[string]interface{}) (map[string]interface{}, error) {
	return fg.parseValues(raw, false)
}

func (fg *FormGroup) parseValues(raw map[string]interface{}, useDefaults bool) (map[string]interface{}, error) {
	parsed := make(map[string]interface{})

	for _, f := range fg.fields {
		if !f.GetForm() {
			if def := f.GetDefault(); def != nil && useDefaults {
				parsed[f.GetID()] = def
			}
			continue
		}

		rawValue, exists := raw[f.GetID()]
		if !exists {
			if f.IsRequired() {
				return nil, fmt.Errorf("required field %s is missing", f.GetID())
			}
			if def := f.GetDefault(); def != nil && useDefaults {
				parsed[f.GetID()] = def
			}
			continue
		}

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
		if !exists && f.IsRequired() && f.GetForm() {
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

// ParseAndValidateSparse is ParseAndValidate without defaults (see ParseValuesSparse).
// Only keys present in the result are validated — a missing filter is no filter.
func (fg *FormGroup) ParseAndValidateSparse(raw map[string]interface{}) (map[string]interface{}, error) {
	parsed, err := fg.ParseValuesSparse(raw)
	if err != nil {
		return nil, err
	}
	for _, f := range fg.fields {
		if v, ok := parsed[f.GetID()]; ok {
			if err := f.Validate(v); err != nil {
				return nil, err
			}
		}
	}
	return parsed, nil
}
