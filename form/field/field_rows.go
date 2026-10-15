package field

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/xiriframework/xiri-go/component/core"
)

// RowsField is a list of rows with the same structure. Fields describe one row; the value is one map per row,
// keyed by sub-field ID. Sub-fields serve as templates only: their Parse/Validate run per cell, their Value is
// never written.
type RowsField struct {
	*BaseField
	Fields    []FormField
	MinRows   *int // applies from the first row on; an empty, optional list is valid
	MaxRows   *int
	AddText   string // translation key for the "add row" button
	Sortable  bool   // drag handle and up/down buttons (default true)
	Addable   bool   // add button and Enter in the last row (default true)
	Removable bool   // delete button (default true)
	Copy      bool   // duplicate button
	Unique    []string
	Widths    map[string]int // relative column width per sub-field ID, exported as "cols"
	Value     []map[string]any
}

// CellError is the failure of one cell. Row is 0-based and refers to the bound value (empty rows dropped).
type CellError struct {
	Row  int
	Sub  string // sub-field ID
	Name string // sub-field label (translation key)
	Err  error
}

// RowErrors collects the failures of a RowsField: Field for the row count, Cells per cell.
// FormGroup.AddFieldError fans them out to "<field>" and "<field>.<row>.<sub>".
type RowErrors struct {
	Field error
	Cells []CellError
}

func (e *RowErrors) Error() string {
	parts := make([]string, 0, len(e.Cells)+1)
	if e.Field != nil {
		parts = append(parts, e.Field.Error())
	}
	for _, c := range e.Cells {
		parts = append(parts, fmt.Sprintf("row %d %s: %v", c.Row+1, c.Sub, c.Err))
	}
	return strings.Join(parts, "; ")
}

func (e *RowErrors) add(row int, sub FormField, err error) {
	e.Cells = append(e.Cells, CellError{Row: row, Sub: sub.GetID(), Name: sub.GetName(), Err: err})
}

func (e *RowErrors) orNil() error {
	if e.Field == nil && len(e.Cells) == 0 {
		return nil
	}
	return e
}

// NewRowsField creates a field with repeatable rows. Supported sub-field types: text (not textarea/html),
// number, bool, select. Panics on duplicate sub-field IDs or unsupported types - both are programming errors
// that would silently break the row maps.
func NewRowsField(id, name string, required bool, fields []FormField) *RowsField {
	seen := make(map[string]bool, len(fields))
	for _, sub := range fields {
		switch sub.GetType() {
		case FieldTypeText, FieldTypeInt, FieldTypeBool, FieldTypeSelect:
		default:
			panic(fmt.Sprintf("rows field %s: sub-field %s has unsupported type %q", id, sub.GetID(), sub.GetType()))
		}
		if t, ok := sub.(*TextField); ok && (t.Subtype == "textarea" || t.Subtype == "html") {
			panic(fmt.Sprintf("rows field %s: sub-field %s: subtype %q is not supported in rows", id, sub.GetID(), t.Subtype))
		}
		if seen[sub.GetID()] {
			panic(fmt.Sprintf("rows field %s: duplicate sub-field id %q", id, sub.GetID()))
		}
		seen[sub.GetID()] = true
	}
	return &RowsField{
		BaseField: &BaseField{ID: id, Type: FieldTypeRows, Name: name, Required: required, Form: true},
		Fields:    fields,
		Sortable:  true,
		Addable:   true,
		Removable: true,
	}
}

// Sub returns the sub-field with the given ID, nil if there is none.
func (f *RowsField) Sub(id string) FormField {
	for _, sub := range f.Fields {
		if sub.GetID() == id {
			return sub
		}
	}
	return nil
}

func (f *RowsField) Parse(raw any) (any, error) {
	if raw == nil {
		if f.Default == nil {
			return nil, nil
		}
		raw = f.Default
	}
	list, err := rowsToRaw(raw)
	if err != nil {
		return nil, fmt.Errorf("rows field %s: %w", f.ID, err)
	}

	rows := make([]map[string]any, 0, len(list))
	errs := &RowErrors{}
	for _, item := range list {
		in, ok := item.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("rows field %s: row must be an object, got %T", f.ID, item)
		}
		row := make(map[string]any, len(f.Fields))
		failed := false
		for _, sub := range f.Fields {
			// Disabled or Form=false: the client must not set it, the default applies (like on form level).
			var v any
			if sub.GetForm() && !sub.IsDisabled() {
				v = in[sub.GetID()]
			}
			p, err := sub.Parse(v)
			if err != nil {
				errs.add(len(rows), sub, err)
				failed = true
				continue
			}
			row[sub.GetID()] = p
		}
		if !failed && f.rowIsEmpty(row) {
			continue
		}
		rows = append(rows, row)
	}
	if err := errs.orNil(); err != nil {
		return nil, err
	}
	return rows, nil
}

func (f *RowsField) Validate(value any) error {
	var rows []map[string]any
	if value != nil {
		var ok bool
		if rows, ok = value.([]map[string]any); !ok {
			return fmt.Errorf("invalid rows value type for %s", f.ID)
		}
	}

	errs := &RowErrors{}
	switch n := len(rows); {
	case n == 0 && f.Required:
		errs.Field = invalid("required", nil, "rows field %s needs at least one row", f.ID)
	case n > 0 && f.MinRows != nil && n < *f.MinRows:
		errs.Field = invalid("min_items", map[string]any{"min": *f.MinRows}, "rows field %s must have at least %d rows", f.ID, *f.MinRows)
	case f.MaxRows != nil && n > *f.MaxRows:
		errs.Field = invalid("max_items", map[string]any{"max": *f.MaxRows}, "rows field %s must have at most %d rows", f.ID, *f.MaxRows)
	}

	for i, row := range rows {
		for _, sub := range f.Fields {
			if err := sub.Validate(row[sub.GetID()]); err != nil {
				errs.add(i, sub, err)
			}
		}
	}
	f.checkUnique(rows, errs)
	return errs.orNil()
}

// checkUnique marks every row whose Unique columns repeat an earlier row. Rows whose Unique values are all
// blank do not count, so two new rows without key do not collide.
func (f *RowsField) checkUnique(rows []map[string]any, errs *RowErrors) {
	if len(f.Unique) == 0 {
		return
	}
	first := f.Sub(f.Unique[0])
	seen := make(map[string]bool, len(rows))
	for i, row := range rows {
		key, blank := uniqueKey(f.Unique, row)
		if blank {
			continue
		}
		if seen[key] {
			errs.add(i, first, invalid("duplicate", nil, "rows field %s: row %d repeats an earlier row", f.ID, i+1))
			continue
		}
		seen[key] = true
	}
}

func uniqueKey(ids []string, row map[string]any) (string, bool) {
	parts := make([]string, len(ids))
	blank := true
	for i, id := range ids {
		if v := row[id]; v != nil {
			parts[i] = strings.ToLower(strings.TrimSpace(fmt.Sprint(v)))
		}
		if parts[i] != "" {
			blank = false
		}
	}
	return strings.Join(parts, "\x00"), blank
}

// rowIsEmpty: every visible cell blank or on its default, every hidden cell blank.
// Counterpart: isEmptyRow in xiri-ng (rows.util.ts) - both must agree, or row indices of errors shift.
func (f *RowsField) rowIsEmpty(row map[string]any) bool {
	for _, sub := range f.Fields {
		v := row[sub.GetID()]
		if isBlank(v) {
			continue
		}
		if h, ok := sub.(interface{ IsHidden() bool }); ok && h.IsHidden() {
			return false
		}
		def, err := sub.Parse(nil)
		if err != nil || fmt.Sprint(v) != fmt.Sprint(def) {
			return false
		}
	}
	return true
}

func isBlank(v any) bool {
	switch x := v.(type) {
	case nil:
		return true
	case string:
		return x == ""
	case ModelListValue:
		return len(x) == 0
	}
	return false
}

// rowsToRaw brings []any, []map[string]any or a slice of structs into the request shape ([]any of maps).
// Structs go through JSON, so their json tags must match the sub-field IDs.
func rowsToRaw(v any) ([]any, error) {
	if list, ok := v.([]any); ok {
		return list, nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var list []any
	if err := json.Unmarshal(b, &list); err != nil {
		return nil, fmt.Errorf("value must be a list of rows: %w", err)
	}
	return list, nil
}

// BindValue parses, validates, and stores the value in the field.
func (f *RowsField) BindValue(raw any) error {
	parsed, err := f.Parse(raw)
	if err != nil {
		return fmt.Errorf("parsing field %s: %w", f.ID, err)
	}
	if err := f.Validate(parsed); err != nil {
		return fmt.Errorf("validating field %s: %w", f.ID, err)
	}
	f.Value, _ = parsed.([]map[string]any)
	return nil
}

// SetMinRows sets the minimum number of rows (applies from the first row on).
func (f *RowsField) SetMinRows(n int) *RowsField {
	f.MinRows = &n
	return f
}

// SetMaxRows sets the maximum number of rows.
func (f *RowsField) SetMaxRows(n int) *RowsField {
	f.MaxRows = &n
	return f
}

// SetUnique makes the combination of these sub-fields unique across rows (trimmed, case-insensitive; rows
// with all of them blank do not count). Panics on an unknown sub-field ID.
func (f *RowsField) SetUnique(ids ...string) *RowsField {
	for _, id := range ids {
		if f.Sub(id) == nil {
			panic(fmt.Sprintf("rows field %s: unique sub-field %q does not exist", f.ID, id))
		}
	}
	f.Unique = ids
	return f
}

// ExportForFrontend exports the field; value may be []map[string]any, []any or a slice of structs (json tags =
// sub-field IDs). Sub-fields export their default as "value": that is what a new row starts with.
func (f *RowsField) ExportForFrontend(ctx *core.UiContext, value any) map[string]any {
	if value == nil {
		value = f.GetDefault()
	}
	rows := []any{}
	if value != nil {
		if list, err := rowsToRaw(value); err == nil && list != nil {
			rows = list
		}
	}
	result := f.BaseField.GetBaseExport(ctx, rows)

	subs := make([]map[string]any, 0, len(f.Fields))
	for _, sub := range f.Fields {
		if !sub.GetForm() {
			continue
		}
		s := sub.ExportForFrontend(ctx, sub.GetDefault())
		if w := f.Widths[sub.GetID()]; w > 0 {
			s["cols"] = w
		}
		subs = append(subs, s)
	}
	result["fields"] = subs
	result["sortable"] = f.Sortable
	result["addable"] = f.Addable
	result["removable"] = f.Removable
	result["copy"] = f.Copy
	if len(f.Unique) > 0 {
		result["unique"] = f.Unique
	}
	if f.AddText != "" {
		result["addText"] = ctx.SafeTranslate(f.AddText)
	}
	if f.MinRows != nil {
		result["min"] = *f.MinRows
	}
	if f.MaxRows != nil {
		result["max"] = *f.MaxRows
	}
	return result
}

// SetAddText sets the translation key of the "add row" button.
func (f *RowsField) SetAddText(key string) *RowsField {
	f.AddText = key
	return f
}

// SetSortable shows or hides drag handle and up/down buttons (UI only: the app checks the order itself).
func (f *RowsField) SetSortable(v bool) *RowsField {
	f.Sortable = v
	return f
}

// SetAddable shows or hides the add button (UI only: only the app can check that no row was added).
func (f *RowsField) SetAddable(v bool) *RowsField {
	f.Addable = v
	return f
}

// SetRemovable shows or hides the delete button (UI only: only the app can check that no row was removed).
func (f *RowsField) SetRemovable(v bool) *RowsField {
	f.Removable = v
	return f
}

// SetCopy shows a duplicate button; the copy goes below the original, hidden sub-fields stay empty.
func (f *RowsField) SetCopy(v bool) *RowsField {
	f.Copy = v
	return f
}

// SetColumnWidth sets the relative width of a column (CSS fr; default: 1, bool narrow).
func (f *RowsField) SetColumnWidth(subID string, n int) *RowsField {
	if f.Widths == nil {
		f.Widths = map[string]int{}
	}
	f.Widths[subID] = n
	return f
}

// SetClass sets the CSS class for frontend styling
func (f *RowsField) SetClass(class string) *RowsField {
	f.BaseField.SetClass(class)
	return f
}

// SetHint sets the tooltip/help text for the field
func (f *RowsField) SetHint(hint string) *RowsField {
	f.BaseField.SetHint(hint)
	return f
}

// SetStep sets the step indicator for multi-step forms
func (f *RowsField) SetStep(step int) *RowsField {
	f.BaseField.SetStep(step)
	return f
}

// SetDisabled makes the whole field read-only (no add/delete/move)
func (f *RowsField) SetDisabled(disabled bool) *RowsField {
	f.BaseField.SetDisabled(disabled)
	return f
}

// SetForm sets whether to show in form
func (f *RowsField) SetForm(form bool) *RowsField {
	f.BaseField.SetForm(form)
	return f
}

// SetAccess stores role metadata. Metadata only: neither exported nor evaluated by xiri-go, no access control.
func (f *RowsField) SetAccess(access []string) *RowsField {
	f.BaseField.SetAccess(access)
	return f
}

// SetScenario stores scenario metadata. Metadata only: neither exported nor evaluated by xiri-go.
func (f *RowsField) SetScenario(scenario []string) *RowsField {
	f.BaseField.SetScenario(scenario)
	return f
}
