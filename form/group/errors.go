package group

import (
	"sort"
	"strings"
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
