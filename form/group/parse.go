package group

import "github.com/xiriframework/xiri-go/form/field"

// ParseValues parses raw field values into typed values. Missing, non-required
// fields receive their default value.
func (fg *FormGroup) ParseValues(raw map[string]interface{}) (map[string]interface{}, error) {
	return fg.parseValues(raw, true)
}

// ParseValuesSparse is ParseValues without defaults: a key that is missing from raw
// stays missing in the result. Used for filters, where "not sent" means "no filter".
// Exception: fields with Form=false or Disabled=true are never taken from the client,
// so their default is kept — that is how a server-set or locked filter transports its value.
func (fg *FormGroup) ParseValuesSparse(raw map[string]interface{}) (map[string]interface{}, error) {
	return fg.parseValues(raw, false)
}

func (fg *FormGroup) parseValues(raw map[string]interface{}, useDefaults bool) (map[string]interface{}, error) {
	parsed := make(map[string]interface{})
	errs := FieldErrors{}

	for _, f := range fg.fields {
		if !f.GetForm() || f.IsDisabled() {
			// Form=false: der Client sendet das Feld nie. Disabled: der Client darf es nicht
			// setzen (wie BindFromMap). In beiden Fällen ist der Default die einzige Quelle und
			// gilt auch im Sparse-Pfad — sonst verliert ein serverseitig gesetzter oder
			// gesperrter Filter (Einzelobjekt-Ansicht, eigener Mandant) seinen Wert.
			parseDefault(f, parsed, errs)
			continue
		}

		rawValue, exists := raw[f.GetID()]
		if !exists {
			if f.IsRequired() {
				errs[f.GetID()] = "required field " + f.GetID() + " is missing"
				continue
			}
			if useDefaults {
				parseDefault(f, parsed, errs)
			}
			continue
		}

		value, err := f.Parse(rawValue)
		if err != nil {
			errs[f.GetID()] = err.Error()
			continue
		}
		parsed[f.GetID()] = value
	}

	if len(errs) > 0 {
		return nil, errs
	}
	return parsed, nil
}

// parseDefault runs the default through Parse like the bind paths do, so directly set
// defaults (int instead of int32, []int32 instead of ModelListValue) come out normalized.
func parseDefault(f field.FormField, parsed map[string]interface{}, errs FieldErrors) {
	if f.GetDefault() == nil {
		return
	}
	v, err := f.Parse(nil)
	if err != nil {
		errs[f.GetID()] = err.Error()
		return
	}
	if v != nil {
		parsed[f.GetID()] = v
	}
}

// ValidateValues validates parsed field values and reports every failing field at once.
func (fg *FormGroup) ValidateValues(values map[string]interface{}) error {
	errs := FieldErrors{}
	for _, f := range fg.fields {
		value, exists := values[f.GetID()]
		if !exists && f.IsRequired() && f.GetForm() {
			errs[f.GetID()] = "required field " + f.GetID() + " is missing"
			continue
		}
		if err := f.Validate(value); err != nil {
			errs[f.GetID()] = err.Error()
		}
	}
	if len(errs) > 0 {
		return errs
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
	errs := FieldErrors{}
	for _, f := range fg.fields {
		if v, ok := parsed[f.GetID()]; ok {
			if err := f.Validate(v); err != nil {
				errs[f.GetID()] = err.Error()
			}
		}
	}
	if len(errs) > 0 {
		return nil, errs
	}
	return parsed, nil
}
