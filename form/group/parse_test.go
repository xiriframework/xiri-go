package group

import (
	"errors"
	"testing"

	"github.com/xiriframework/xiri-go/form/field"
)

func TestValidateValues_CollectsAllFieldErrors(t *testing.T) {
	name := field.NewTextFieldWithLength("name", "NAME", true, "", 3, 10)
	email := field.NewTextFieldWithLength("email", "EMAIL", true, "", 5, 50)
	fg := NewFormGroup([]field.FormField{name, email})

	err := fg.ValidateValues(map[string]interface{}{"name": "ab", "email": "x"})

	var fe FieldErrors
	if !errors.As(err, &fe) {
		t.Fatalf("expected FieldErrors, got %T: %v", err, err)
	}
	if len(fe) != 2 || fe["name"] == "" || fe["email"] == "" {
		t.Fatalf("expected errors for name and email, got %v", fe)
	}
}

func TestParseValues_MissingRequired_IsFieldError(t *testing.T) {
	name := field.NewTextField("name", "NAME", true, "")
	fg := NewFormGroup([]field.FormField{name})

	_, err := fg.ParseValues(map[string]interface{}{})

	var fe FieldErrors
	if !errors.As(err, &fe) || fe["name"] == "" {
		t.Fatalf("expected FieldErrors with key name, got %T: %v", err, err)
	}
}

func TestParseAndValidateSparse_CollectsAllFieldErrors(t *testing.T) {
	a := field.NewTextFieldWithLength("a", "A", false, "", 3, 10)
	b := field.NewTextFieldWithLength("b", "B", false, "", 3, 10)
	fg := NewFormGroup([]field.FormField{a, b})

	_, err := fg.ParseAndValidateSparse(map[string]interface{}{"a": "x", "b": "y"})

	var fe FieldErrors
	if !errors.As(err, &fe) || len(fe) != 2 {
		t.Fatalf("expected 2 field errors, got %T: %v", err, err)
	}
}

func TestParseValues_DisabledFieldIgnoresClientValue(t *testing.T) {
	owner := field.NewIntField("owner", "OWNER", true, 7)
	owner.SetDisabled(true)
	fg := NewFormGroup([]field.FormField{owner})

	for name, parse := range map[string]func(map[string]interface{}) (map[string]interface{}, error){
		"ParseAndValidate":       fg.ParseAndValidate,
		"ParseAndValidateSparse": fg.ParseAndValidateSparse,
	} {
		got, err := parse(map[string]interface{}{"owner": float64(999)})
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", name, err)
		}
		if got["owner"] != int32(7) {
			t.Errorf("%s: owner = %#v, want default int32(7)", name, got["owner"])
		}
	}
}

func TestParseValuesSparse_DisabledFieldKeepsDefaultWhenMissing(t *testing.T) {
	owner := field.NewIntField("owner", "OWNER", true, 7)
	owner.SetDisabled(true)
	fg := NewFormGroup([]field.FormField{owner})

	got, err := fg.ParseAndValidateSparse(map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["owner"] != int32(7) {
		t.Errorf("owner = %#v, want default int32(7)", got["owner"])
	}
}

// Direkt gesetzte Defaults laufen durch Parse wie in den Bind-Pfaden.
func TestParseValues_ModelIntDefaultIsNormalized(t *testing.T) {
	m := field.NewModelField("m", "M", false, "device", 0)
	m.Default = 42
	fg := NewFormGroup([]field.FormField{m})

	got, err := fg.ParseValues(map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["m"] != int32(42) {
		t.Errorf("m = %#v, want int32(42)", got["m"])
	}
}

func TestParseValues_ModelListInt32SliceDefault(t *testing.T) {
	ml := field.NewModelListField("ml", "ML", false, "device", nil)
	ml.Default = []int32{9}
	fg := NewFormGroup([]field.FormField{ml})

	got, err := fg.ParseAndValidate(map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := got["ml"].(field.ModelListValue); !ok || len(v) != 1 || v[0] != 9 {
		t.Errorf("ml = %#v, want ModelListValue{9}", got["ml"])
	}
}

func TestParseValuesSparse_DisabledModelIntDefault(t *testing.T) {
	m := field.NewModelField("m", "M", false, "device", 0)
	m.Default = 42
	m.SetDisabled(true)
	fg := NewFormGroup([]field.FormField{m})

	got, err := fg.ParseValuesSparse(map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["m"] != int32(42) {
		t.Errorf("m = %#v, want int32(42)", got["m"])
	}
}

func TestParseValues_InvalidDefaultIsFieldError(t *testing.T) {
	ml := field.NewModelListField("ml", "ML", false, "device", nil)
	ml.Default = []int{9}
	fg := NewFormGroup([]field.FormField{ml})

	_, err := fg.ParseValues(map[string]interface{}{})

	var fe FieldErrors
	if !errors.As(err, &fe) || fe["ml"] == "" {
		t.Fatalf("expected FieldErrors with key ml, got %T: %v", err, err)
	}
}

func TestParseValues_TextDefaultIsTrimmed(t *testing.T) {
	fg := NewFormGroup([]field.FormField{field.NewTextField("t", "T", false, " x ")})

	got, err := fg.ParseValues(map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["t"] != "x" {
		t.Errorf("t = %#v, want trimmed \"x\"", got["t"])
	}
}

func TestParseValues_MissingArrayWithoutDefaultStaysMissing(t *testing.T) {
	// Default explizit nil: der Konstruktor speichert ein typisiertes nil-Slice (≠ nil).
	a := field.NewArrayField("a", "A", false, "string", nil)
	a.Default = nil
	fg := NewFormGroup([]field.FormField{a})

	got, err := fg.ParseValues(map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := got["a"]; ok {
		t.Errorf("a = %#v, want missing", got["a"])
	}
}

func TestParseValuesSparse_MissingFilterStaysMissing(t *testing.T) {
	m := field.NewModelField("m", "M", false, "device", 0)
	m.Default = 42
	fg := NewFormGroup([]field.FormField{m})

	got, err := fg.ParseValuesSparse(map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := got["m"]; ok {
		t.Errorf("m = %#v, want missing", got["m"])
	}
}
