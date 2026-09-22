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
