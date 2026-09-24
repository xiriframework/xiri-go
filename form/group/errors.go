package group

import (
	"errors"
	"fmt"
	"sort"
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
	var ve *field.ValidationError
	if !errors.As(err, &ve) {
		return err.Error()
	}
	key := "validation." + ve.Code
	t := fg.ctx.SafeTranslate(key)
	if t == key || t == "" {
		return err.Error()
	}
	pairs := []string{"{field}", fg.label(fieldID)}
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
	if name := fg.GetTranslatedName(fieldID); name != "" {
		return name
	}
	return fieldID
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
