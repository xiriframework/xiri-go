package group

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/xiriframework/xiri-go/form/field"
)

// FieldErrors sammelt Validierungsfehler pro Feld-ID. Implementiert error, damit bestehende
// Aufrufer (`wc.BadRequest(err.Error())`) unverändert weiterlaufen; response.NewErrorResponseFromError
// liest die Map über FieldErrors() aus, ohne dieses Paket zu importieren.
type FieldErrors map[string]string

func (e FieldErrors) Error() string {
	ids := make([]string, 0, len(e))
	for id := range e {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	parts := make([]string, len(ids))
	for i, id := range ids {
		parts[i] = id + ": " + e[id]
	}
	return strings.Join(parts, "; ")
}

// FieldErrors liefert die Map; Zielmethode für errors.As über ein Interface.
func (e FieldErrors) FieldErrors() map[string]string { return e }

// FieldErrorMessage is the message for one failing field. A field.ValidationError (also wrapped,
// e.g. by BindValue) is translated via the key "validation.<Code>"; {field} and the Params keys
// are filled in. Without a translation (no TranslateFunc, key returned unchanged or "") the
// English err.Error() stays.
func (fg *FormGroup) FieldErrorMessage(fieldID string, err error) string {
	return fg.message(fg.label(fieldID), err)
}

// AddFieldError stores the message for err under the field's ID. A field.RowErrors fans out: its Field
// error under the field ID, each cell under "<field>.<row>.<sub>" with the sub-field's label as {field}.
func (fg *FormGroup) AddFieldError(errs FieldErrors, f field.FormField, err error) {
	var re *field.RowErrors
	if !errors.As(err, &re) {
		errs[f.GetID()] = fg.FieldErrorMessage(f.GetID(), err)
		return
	}
	if re.Field != nil {
		errs[f.GetID()] = fg.FieldErrorMessage(f.GetID(), re.Field)
	}
	for _, c := range re.Cells {
		// First error per cell wins: a duplicate must not hide the cell's own validation error.
		key := fmt.Sprintf("%s.%d.%s", f.GetID(), c.Row, c.Sub)
		if _, taken := errs[key]; !taken {
			errs[key] = fg.message(fg.ctx.SafeTranslate(c.Name), c.Err)
		}
	}
}

func (fg *FormGroup) message(label string, err error) string {
	var ve *field.ValidationError
	if !errors.As(err, &ve) {
		return err.Error()
	}
	key := "validation." + ve.Code
	t := fg.ctx.SafeTranslate(key)
	if t == key || t == "" {
		return err.Error()
	}
	pairs := []string{"{field}", label}
	for k, v := range ve.Params {
		pairs = append(pairs, "{"+k+"}", fmt.Sprint(v))
	}
	return strings.NewReplacer(pairs...).Replace(t)
}

// WrapFieldErrors names the fields by their translated label in Error() (the banner) once a
// TranslateFunc is set; without one errs is returned unchanged.
func (fg *FormGroup) WrapFieldErrors(errs FieldErrors) error {
	if fg.ctx == nil || fg.ctx.Translate == nil {
		return errs
	}
	labels := make(map[string]string, len(errs))
	for id := range errs {
		labels[id] = fg.label(id)
	}
	return labeledFieldErrors{errs: errs, labels: labels}
}

func (fg *FormGroup) label(fieldID string) string {
	if l, ok := fg.cellLabel(fieldID); ok {
		return l
	}
	if name := fg.GetTranslatedName(fieldID); name != "" {
		return name
	}
	return fieldID
}

// cellLabel names a RowsField cell key "<field>.<row>.<sub>" for the banner: "Positionen 3 · Bezeichnung".
func (fg *FormGroup) cellLabel(key string) (string, bool) {
	parts := strings.SplitN(key, ".", 3)
	if len(parts) != 3 {
		return "", false
	}
	rows, ok := fg.index[parts[0]].(*field.RowsField)
	if !ok {
		return "", false
	}
	row, err := strconv.Atoi(parts[1])
	sub := rows.Sub(parts[2])
	if err != nil || sub == nil {
		return "", false
	}
	return fmt.Sprintf("%s %d · %s", fg.label(parts[0]), row+1, fg.ctx.SafeTranslate(sub.GetName())), true
}

// labeledFieldErrors does not embed FieldErrors: the embedded field would shadow the FieldErrors() method.
type labeledFieldErrors struct {
	errs   FieldErrors
	labels map[string]string
}

func (e labeledFieldErrors) Error() string {
	ids := make([]string, 0, len(e.errs))
	for id := range e.errs {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	parts := make([]string, len(ids))
	for i, id := range ids {
		parts[i] = e.labels[id] + ": " + e.errs[id]
	}
	return strings.Join(parts, "; ")
}

func (e labeledFieldErrors) FieldErrors() map[string]string { return e.errs }

func (e labeledFieldErrors) Unwrap() error { return e.errs }
