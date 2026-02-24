package group

import (
	"testing"

	"github.com/xiriframework/xiri-go/form/field"
)

func TestNewFormGroup(t *testing.T) {
	fields := []field.FormField{
		field.NewTextField("name", "NAME", true, ""),
		field.NewBoolField("active", "ACTIVE", false, true),
	}
	fg := NewFormGroup(fields)

	if len(fg.GetFields()) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(fg.GetFields()))
	}

	f, ok := fg.GetField("name")
	if !ok {
		t.Fatal("expected to find field 'name'")
	}
	if f.GetID() != "name" {
		t.Errorf("expected field ID 'name', got %q", f.GetID())
	}
}

func TestGetFieldIDs(t *testing.T) {
	fields := []field.FormField{
		field.NewTextField("a", "A", false, ""),
		field.NewIntField("b", "B", true, 42),
	}
	fg := NewFormGroup(fields)

	ids := fg.GetFieldIDs()
	if len(ids) != 2 || ids[0] != "a" || ids[1] != "b" {
		t.Errorf("unexpected field IDs: %v", ids)
	}
}

func TestGetDefaults(t *testing.T) {
	fields := []field.FormField{
		field.NewTextField("name", "NAME", true, "default-name"),
		field.NewIntField("count", "COUNT", false, 10),
	}
	fg := NewFormGroup(fields)

	defaults := fg.GetDefaults()
	if defaults["name"] != "default-name" {
		t.Errorf("expected default 'default-name', got %v", defaults["name"])
	}
	if defaults["count"] != int32(10) {
		t.Errorf("expected default 10, got %v", defaults["count"])
	}
}

func TestGetRequiredOptionalFields(t *testing.T) {
	fields := []field.FormField{
		field.NewTextField("name", "NAME", true, ""),
		field.NewBoolField("active", "ACTIVE", false, true),
		field.NewIntField("count", "COUNT", true, 0),
	}
	fg := NewFormGroup(fields)

	required := fg.GetRequiredFields()
	if len(required) != 2 {
		t.Errorf("expected 2 required fields, got %d", len(required))
	}

	optional := fg.GetOptionalFields()
	if len(optional) != 1 {
		t.Errorf("expected 1 optional field, got %d", len(optional))
	}
}

func TestGetFieldNotFound(t *testing.T) {
	fg := NewFormGroup(nil)
	_, ok := fg.GetField("nonexistent")
	if ok {
		t.Error("expected field not found")
	}
}

func TestParseAndValidate(t *testing.T) {
	fields := []field.FormField{
		field.NewTextField("name", "NAME", true, ""),
		field.NewIntField("count", "COUNT", false, 0),
	}
	fg := NewFormGroup(fields)

	// Valid data
	raw := map[string]interface{}{
		"name":  "test",
		"count": float64(42),
	}
	parsed, err := fg.ParseAndValidate(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if parsed["name"] != "test" {
		t.Errorf("expected 'test', got %v", parsed["name"])
	}
	if parsed["count"] != 42 {
		t.Errorf("expected 42, got %v", parsed["count"])
	}
}

func TestParseAndValidate_MissingRequired(t *testing.T) {
	fields := []field.FormField{
		field.NewTextField("name", "NAME", true, ""),
	}
	fg := NewFormGroup(fields)

	raw := map[string]interface{}{}
	_, err := fg.ParseAndValidate(raw)
	if err == nil {
		t.Fatal("expected error for missing required field")
	}
}

func TestParseAndValidate_OptionalUsesDefault(t *testing.T) {
	fields := []field.FormField{
		field.NewTextField("note", "NOTE", false, "default-note"),
	}
	fg := NewFormGroup(fields)

	raw := map[string]interface{}{}
	parsed, err := fg.ParseAndValidate(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if parsed["note"] != "default-note" {
		t.Errorf("expected default 'default-note', got %v", parsed["note"])
	}
}

func TestExportForFrontend(t *testing.T) {
	fields := []field.FormField{
		field.NewTextField("name", "NAME", true, "hello"),
	}
	fg := NewFormGroup(fields)

	exported := fg.ExportForFrontend()
	if len(exported) != 1 {
		t.Fatalf("expected 1 field export, got %d", len(exported))
	}
	if exported[0]["id"] != "name" {
		t.Errorf("expected field id 'name', got %v", exported[0]["id"])
	}
	if exported[0]["value"] != "hello" {
		t.Errorf("expected value 'hello', got %v", exported[0]["value"])
	}
}

func TestExportForFrontendWithValues(t *testing.T) {
	fields := []field.FormField{
		field.NewTextField("name", "NAME", true, "default"),
	}
	fg := NewFormGroup(fields)

	values := map[string]interface{}{"name": "override"}
	exported := fg.ExportForFrontendWithValues(values)
	if exported[0]["value"] != "override" {
		t.Errorf("expected value 'override', got %v", exported[0]["value"])
	}
}
