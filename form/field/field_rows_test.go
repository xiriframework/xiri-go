package field

import (
	"errors"
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
)

// checklistRows ist das Abnahme-Beispiel aus der Anforderung von tracking.
func checklistRows() *RowsField {
	key := NewTextField("key", "KEY", false, "")
	key.BaseField.SetHide(true)
	category := NewTextField("category", "KATEGORIE", false, "")
	label := NewTextField("label", "BEZEICHNUNG", true, "")
	required := NewBoolField("required", "PFLICHT", false, true)
	photo := NewBoolField("photo", "FOTO", false, false)
	return NewRowsField("items", "POSITIONEN", true, []FormField{key, category, label, required, photo})
}

func rowErrors(t *testing.T, err error) *RowErrors {
	t.Helper()
	var re *RowErrors
	if !errors.As(err, &re) {
		t.Fatalf("expected *RowErrors, got %T: %v", err, err)
	}
	return re
}

func code(t *testing.T, err error) string {
	t.Helper()
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected *ValidationError, got %T: %v", err, err)
	}
	return ve.Code
}

func TestRowsField_ParseTrimsDropsEmptyRowsAndUnknownKeys(t *testing.T) {
	f := checklistRows()
	got, err := f.Parse([]any{
		map[string]any{"key": "brakes", "label": "  Bremsen ", "required": true, "photo": true, "evil": "x"},
		// leer: sichtbare Zellen leer bzw. auf Default (required=true, photo=false), key leer
		map[string]any{"key": "", "category": "", "label": "", "required": true, "photo": false},
	})
	if err != nil {
		t.Fatal(err)
	}
	rows := got.([]map[string]any)
	if len(rows) != 1 {
		t.Fatalf("want 1 row, got %d: %v", len(rows), rows)
	}
	if rows[0]["label"] != "Bremsen" {
		t.Errorf("label not trimmed: %q", rows[0]["label"])
	}
	if _, ok := rows[0]["evil"]; ok {
		t.Error("unknown key must be dropped")
	}
}

func TestRowsField_RowWithHiddenValueIsNeverEmpty(t *testing.T) {
	f := checklistRows()
	got, _ := f.Parse([]any{map[string]any{"key": "brakes", "label": "", "required": true}})
	if rows := got.([]map[string]any); len(rows) != 1 {
		t.Fatalf("row with key must stay, got %v", rows)
	}
}

func TestRowsField_ParseAcceptsJSONNumbers(t *testing.T) {
	qty := NewIntField("qty", "MENGE", false, 0)
	f := NewRowsField("items", "POSITIONEN", false, []FormField{qty})
	got, err := f.Parse([]any{map[string]any{"qty": float64(3)}})
	if err != nil {
		t.Fatal(err)
	}
	if v := got.([]map[string]any)[0]["qty"]; v != 3 {
		t.Fatalf("qty = %#v", v)
	}
}

func TestRowsField_DisabledSubFieldIgnoresClientValue(t *testing.T) {
	f := checklistRows()
	f.Sub("required").(*BoolField).BaseField.SetDisabled(true)
	got, _ := f.Parse([]any{map[string]any{"label": "Bremsen", "required": false}})
	if v := got.([]map[string]any)[0]["required"]; v != true {
		t.Fatalf("disabled sub-field must keep its default, got %v", v)
	}
}

func TestRowsField_ParseRejectsNonList(t *testing.T) {
	f := checklistRows()
	if _, err := f.Parse("x"); err == nil {
		t.Error("string must be rejected")
	}
	if _, err := f.Parse([]any{"x"}); err == nil {
		t.Error("non-object row must be rejected")
	}
}

func TestRowsField_ParseCellErrorCarriesRowAndSub(t *testing.T) {
	kind := NewSelectField("kind", "ART", false, []SelectOption{{Value: int32(1), Label: "A"}})
	label := NewTextField("label", "BEZEICHNUNG", false, "")
	f := NewRowsField("items", "POSITIONEN", false, []FormField{label, kind})
	_, err := f.Parse([]any{
		map[string]any{"label": "a", "kind": float64(1)},
		map[string]any{"label": "b", "kind": float64(9)},
	})
	re := rowErrors(t, err)
	if len(re.Cells) != 1 || re.Cells[0].Row != 1 || re.Cells[0].Sub != "kind" || re.Cells[0].Name != "ART" {
		t.Fatalf("cells: %+v", re.Cells)
	}
	if c := code(t, re.Cells[0].Err); c != "not_allowed" {
		t.Fatalf("code %q", c)
	}
}

func TestRowsField_ValidateCellRequired(t *testing.T) {
	f := checklistRows()
	v, err := f.Parse([]any{
		map[string]any{"label": "A"},
		map[string]any{"label": "B"},
		map[string]any{"label": "", "category": "Technik"},
	})
	if err != nil {
		t.Fatal(err)
	}
	re := rowErrors(t, f.Validate(v))
	if re.Field != nil || len(re.Cells) != 1 {
		t.Fatalf("errors: %+v", re)
	}
	if c := re.Cells[0]; c.Row != 2 || c.Sub != "label" || code(t, c.Err) != "required" {
		t.Fatalf("cell: %+v", c)
	}
}

func TestRowsField_ValidateRowCount(t *testing.T) {
	one := []map[string]any{{"label": "A", "required": true, "photo": false}}
	two := []map[string]any{one[0], {"label": "B", "required": true, "photo": false}}

	f := checklistRows()
	if c := code(t, rowErrors(t, f.Validate(nil)).Field); c != "required" {
		t.Errorf("nil + required: %q", c)
	}
	if c := code(t, rowErrors(t, f.Validate([]map[string]any{})).Field); c != "required" {
		t.Errorf("empty + required: %q", c)
	}

	f.SetMaxRows(1)
	if c := code(t, rowErrors(t, f.Validate(two)).Field); c != "max_items" {
		t.Errorf("max: %q", c)
	}

	f = checklistRows().SetMinRows(2)
	if c := code(t, rowErrors(t, f.Validate(one)).Field); c != "min_items" {
		t.Errorf("min: %q", c)
	}
	f.Required = false
	if err := f.Validate([]map[string]any{}); err != nil {
		t.Errorf("min applies only from the first row on: %v", err)
	}
}

func TestRowsField_ValidateUnique(t *testing.T) {
	f := checklistRows().SetUnique("label")
	v, _ := f.Parse([]any{
		map[string]any{"label": "Bremsen"},
		map[string]any{"label": "Licht"},
		map[string]any{"label": " bremsen"},
	})
	re := rowErrors(t, f.Validate(v))
	if len(re.Cells) != 1 || re.Cells[0].Row != 2 || re.Cells[0].Sub != "label" || code(t, re.Cells[0].Err) != "duplicate" {
		t.Fatalf("cells: %+v", re.Cells)
	}

	// Leere Werte zählen nicht: zwei neue Zeilen ohne key kollidieren nicht.
	f = checklistRows().SetUnique("key")
	v, _ = f.Parse([]any{map[string]any{"label": "A"}, map[string]any{"label": "B"}})
	if err := f.Validate(v); err != nil {
		t.Fatalf("blank keys must not collide: %v", err)
	}
}

func TestRowsField_BindValueStoresRows(t *testing.T) {
	f := checklistRows()
	if err := f.BindValue([]any{map[string]any{"key": "k1", "label": "A"}}); err != nil {
		t.Fatal(err)
	}
	if len(f.Value) != 1 || f.Value[0]["key"] != "k1" || f.Value[0]["required"] != true {
		t.Fatalf("Value: %v", f.Value)
	}
}

func TestNewRowsField_PanicsOnInvalidDefinition(t *testing.T) {
	cases := map[string]func(){
		"duplicate id": func() {
			NewRowsField("r", "R", false, []FormField{NewTextField("a", "A", false, ""), NewTextField("a", "A", false, "")})
		},
		"unsupported type": func() {
			NewRowsField("r", "R", false, []FormField{NewChipsField("c", "C", false)})
		},
		"textarea": func() {
			ta := NewTextField("t", "T", false, "")
			ta.Subtype = "textarea"
			NewRowsField("r", "R", false, []FormField{ta})
		},
		"nested rows": func() {
			NewRowsField("r", "R", false, []FormField{NewRowsField("n", "N", false, nil)})
		},
		"unknown unique id": func() {
			NewRowsField("r", "R", false, []FormField{NewTextField("a", "A", false, "")}).SetUnique("x")
		},
	}
	for name, fn := range cases {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("expected panic")
				}
			}()
			fn()
		})
	}
}

func TestRowsField_NonFormSubFieldIgnoresClientValue(t *testing.T) {
	cat := NewTextField("category", "KATEGORIE", false, "fix")
	cat.BaseField.SetForm(false)
	label := NewTextField("label", "BEZEICHNUNG", false, "")
	f := NewRowsField("items", "POSITIONEN", false, []FormField{cat, label})
	got, err := f.Parse([]any{map[string]any{"label": "a", "category": "hacked"}})
	if err != nil {
		t.Fatal(err)
	}
	if v := got.([]map[string]any)[0]["category"]; v != "fix" {
		t.Fatalf("Form=false sub-field must keep its default, got %v", v)
	}
}

type checklistItem struct {
	Key      string `json:"key"`
	Label    string `json:"label"`
	Required bool   `json:"required"`
	Qty      int32  `json:"qty"`
}

func TestRowsField_Export(t *testing.T) {
	ctx := &core.UiContext{Translate: func(k string) string {
		return map[string]string{"POSITIONEN": "Positionen", "ADD": "Position hinzufügen", "BEZEICHNUNG": "Bezeichnung"}[k]
	}}
	f := checklistRows().SetMinRows(1).SetMaxRows(50).SetAddText("ADD").SetSortable(false).SetCopy(true).
		SetUnique("label").SetColumnWidth("label", 3)

	out := f.ExportForFrontend(ctx, []checklistItem{{Key: "k1", Label: "Bremsen", Required: true, Qty: 2}})

	if out["type"] != "rows" || out["name"] != "Positionen" || out["addText"] != "Position hinzufügen" {
		t.Fatalf("head: %v", out)
	}
	if out["min"] != 1 || out["max"] != 50 || out["sortable"] != false || out["addable"] != true ||
		out["removable"] != true || out["copy"] != true {
		t.Fatalf("switches: %v", out)
	}
	if u := out["unique"].([]string); len(u) != 1 || u[0] != "label" {
		t.Fatalf("unique: %v", out["unique"])
	}
	subs := out["fields"].([]map[string]any)
	if len(subs) != 5 || subs[2]["id"] != "label" || subs[2]["cols"] != 3 || subs[2]["name"] != "Bezeichnung" {
		t.Fatalf("fields: %v", subs)
	}
	if subs[0]["hide"] != true || subs[3]["value"] != true {
		t.Fatalf("hide/default: %v / %v", subs[0], subs[3])
	}
	rows := out["value"].([]any)
	if r := rows[0].(map[string]any); r["key"] != "k1" || r["label"] != "Bremsen" || r["qty"] != nil {
		t.Fatalf("value: %v", rows)
	}
}

func TestRowsField_ExportWithoutValueIsEmptyList(t *testing.T) {
	out := checklistRows().ExportForFrontend(nil, nil)
	if rows, ok := out["value"].([]any); !ok || len(rows) != 0 {
		t.Fatalf("value: %#v", out["value"])
	}
}

func TestRowsField_ExportSkipsNonFormSubFields(t *testing.T) {
	secret := NewTextField("secret", "SECRET", false, "x")
	secret.BaseField.SetForm(false)
	f := NewRowsField("items", "POSITIONEN", false, []FormField{NewTextField("label", "L", false, ""), secret})
	subs := f.ExportForFrontend(&core.UiContext{}, nil)["fields"].([]map[string]any)
	if len(subs) != 1 || subs[0]["id"] != "label" {
		t.Fatalf("expected only label, got %v", subs)
	}
}

func TestRowsField_ExportRowsCarryOnlySubFieldKeys(t *testing.T) {
	secret := NewTextField("secret", "SECRET", false, "x")
	secret.BaseField.SetForm(false)
	f := NewRowsField("items", "POSITIONEN", false, []FormField{NewTextField("label", "L", false, ""), secret})
	out := f.ExportForFrontend(&core.UiContext{}, []map[string]any{
		{"label": "a", "secret": "TOPSECRET", "internal_hash": "HASH"},
	})
	row := out["value"].([]any)[0].(map[string]any)
	if len(row) != 1 || row["label"] != "a" {
		t.Fatalf("row must carry only exported sub-field keys, got %v", row)
	}
}

func TestRowsField_DisabledWithDefaultKeepsRows(t *testing.T) {
	f := NewRowsField("items", "POSITIONEN", true, []FormField{NewTextField("label", "L", true, "")}).
		SetDisabled(true).
		SetDefault([]map[string]any{{"label": "a"}, {"label": "b"}})
	if err := f.BindValue(nil); err != nil {
		t.Fatal(err)
	}
	if len(f.Value) != 2 || f.Value[1]["label"] != "b" {
		t.Fatalf("value: %v", f.Value)
	}
}

func TestRowsField_NumberCellsAreAlwaysInt(t *testing.T) {
	f := NewRowsField("items", "POSITIONEN", false, []FormField{NewTextField("label", "L", false, ""), NewIntField("qty", "Q", false, 0)})
	got, err := f.Parse([]any{
		map[string]any{"label": "a", "qty": nil},
		map[string]any{"label": "b", "qty": float64(3)},
	})
	if err != nil {
		t.Fatal(err)
	}
	for i, row := range got.([]map[string]any) {
		if _, ok := row["qty"].(int); !ok {
			t.Errorf("row %d qty is %T, want int", i, row["qty"])
		}
	}
}

func TestRowsField_CellErrorRowIndexSkipsDroppedEmptyRows(t *testing.T) {
	kind := NewSelectField("kind", "ART", false, []SelectOption{{Value: int32(1), Label: "A"}})
	label := NewTextField("label", "L", false, "")
	f := NewRowsField("items", "POSITIONEN", false, []FormField{label, kind})
	_, err := f.Parse([]any{
		map[string]any{"label": "a"},
		map[string]any{"label": ""},
		map[string]any{"label": "c", "kind": float64(9)},
	})
	re := rowErrors(t, err)
	if len(re.Cells) != 1 || re.Cells[0].Row != 1 || re.Cells[0].Sub != "kind" {
		t.Fatalf("cells: %+v", re.Cells)
	}
}

func TestRowsField_UniqueErrorOnFirstVisibleSubField(t *testing.T) {
	f := checklistRows().SetUnique("key", "label")
	v, _ := f.Parse([]any{map[string]any{"key": "k", "label": "a"}, map[string]any{"key": "k", "label": "a"}})
	re := rowErrors(t, f.Validate(v))
	if len(re.Cells) != 1 || re.Cells[0].Sub != "label" {
		t.Fatalf("cells: %+v", re.Cells)
	}
}
