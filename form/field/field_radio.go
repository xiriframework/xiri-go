package field

// NewRadioField creates a single-select radio-group form field for small
// option sets. It reuses SelectField's parsing, validation and export logic
// but sets the field type to "radio" so the frontend renders a
// <mat-radio-group> instead of a dropdown. Options are exported as "list"
// (id/name), identical to select. Multi-select is not supported for radio.
func NewRadioField(id, name string, required bool, options []SelectOption) *SelectField {
	f := NewSelectField(id, name, required, options)
	f.Type = FieldTypeRadio
	return f
}
