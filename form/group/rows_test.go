package group

import (
	"strings"
	"testing"

	"github.com/xiriframework/xiri-go/form/field"
)

func rowsGroup(t *testing.T) (*FormGroup, *field.RowsField) {
	t.Helper()
	label := field.NewTextField("label", "BEZEICHNUNG", true, "")
	items := field.NewRowsField("items", "POSITIONEN", true, []field.FormField{label}).SetMaxRows(2)
	name := field.NewTextField("name", "NAME", true, "")
	texts := map[string]string{
		"POSITIONEN":           "Positionen",
		"BEZEICHNUNG":          "Bezeichnung",
		"NAME":                 "Name",
		"validation.required":  "{field} ist ein Pflichtfeld",
		"validation.max_items": "Höchstens {max} Einträge",
	}
	return groupWithCtx(t, deCtx(texts), name, items), items
}

func TestParseAndValidate_RowsRequiredCell(t *testing.T) {
	label := field.NewTextField("label", "BEZEICHNUNG", true, "")
	cat := field.NewTextField("category", "KATEGORIE", false, "")
	items := field.NewRowsField("items", "POSITIONEN", true, []field.FormField{cat, label})
	fg := groupWithCtx(t, deCtx(map[string]string{
		"POSITIONEN": "Positionen", "BEZEICHNUNG": "Bezeichnung",
		"validation.required": "{field} ist ein Pflichtfeld",
	}), items)

	_, err := fg.ParseAndValidate(map[string]interface{}{"items": []any{
		map[string]any{"category": "Technik", "label": "Bremsen"},
		map[string]any{"category": "Technik", "label": ""},
	}})

	fe := fieldErrors(t, err)
	if got := fe["items.1.label"]; got != "Bezeichnung ist ein Pflichtfeld" {
		t.Fatalf("cell message: %q (all: %v)", got, fe)
	}
	if !strings.Contains(err.Error(), "Positionen 2 · Bezeichnung: Bezeichnung ist ein Pflichtfeld") {
		t.Fatalf("banner: %q", err.Error())
	}
}

func TestParseAndValidate_RowsFieldLevelError(t *testing.T) {
	fg, _ := rowsGroup(t)
	_, err := fg.ParseAndValidate(map[string]interface{}{
		"name":  "Abfahrt",
		"items": []any{map[string]any{"label": "A"}, map[string]any{"label": "B"}, map[string]any{"label": "C"}},
	})
	if got := fieldErrors(t, err)["items"]; got != "Höchstens 2 Einträge" {
		t.Fatalf("field message: %q", got)
	}
}

func TestParseValues_RowsParseErrorFansOut(t *testing.T) {
	kind := field.NewSelectField("kind", "ART", false, []field.SelectOption{{Value: int32(1), Label: "A"}})
	items := field.NewRowsField("items", "POSITIONEN", false, []field.FormField{kind})
	fg := groupWithCtx(t, deCtx(nil), items)

	_, err := fg.ParseValues(map[string]interface{}{"items": []any{map[string]any{"kind": float64(9)}}})
	if _, ok := fieldErrors(t, err)["items.0.kind"]; !ok {
		t.Fatalf("want items.0.kind, got %v", err)
	}
}

func TestLabel_DottedIDOfNormalFieldStaysPlain(t *testing.T) {
	odd := field.NewTextField("a.b.c", "ODD", true, "")
	fg := groupWithCtx(t, deCtx(map[string]string{"ODD": "Komisch"}), odd)
	_, err := fg.ParseAndValidate(map[string]interface{}{})
	if !strings.HasPrefix(err.Error(), "Komisch: ") {
		t.Fatalf("banner: %q", err.Error())
	}
}

// Trifft ein Duplikat eine Zelle, die schon selbst ungültig ist, bleibt der Zellfehler stehen.
func TestParseAndValidate_RowsFirstCellErrorWins(t *testing.T) {
	label := field.NewTextFieldWithLength("label", "BEZEICHNUNG", false, "", 0, 3)
	items := field.NewRowsField("items", "POSITIONEN", false, []field.FormField{label}).SetUnique("label")
	fg := groupWithCtx(t, deCtx(map[string]string{
		"validation.max_length": "Höchstens {max} Zeichen",
		"validation.duplicate":  "doppelt",
	}), items)

	_, err := fg.ParseAndValidate(map[string]interface{}{"items": []any{
		map[string]any{"label": "abcd"},
		map[string]any{"label": "abcd"},
	}})

	if got := fieldErrors(t, err)["items.1.label"]; got != "Höchstens 3 Zeichen" {
		t.Fatalf("cell message: %q", got)
	}
}
