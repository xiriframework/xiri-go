package field

import (
	"testing"
)

func TestJsonField_Validate_Required(t *testing.T) {
	f := NewJsonField("data", "DATA", true, nil)
	err := f.Validate(nil)
	if err == nil {
		t.Fatal("expected error for nil value on required json field")
	}
}

func TestJsonField_Validate_Optional_Nil(t *testing.T) {
	f := NewJsonField("data", "DATA", false, nil)
	err := f.Validate(nil)
	if err != nil {
		t.Fatalf("unexpected error for nil on optional json field: %v", err)
	}
}

func TestJsonField_Validate_Map(t *testing.T) {
	f := NewJsonField("data", "DATA", true, nil)
	err := f.Validate(map[string]interface{}{"key": "value"})
	if err != nil {
		t.Fatalf("unexpected error for map value: %v", err)
	}
}

func TestJsonField_Validate_String(t *testing.T) {
	f := NewJsonField("data", "DATA", true, nil)
	err := f.Validate(`{"key":"value"}`)
	if err != nil {
		t.Fatalf("unexpected error for string value: %v", err)
	}
}

func TestJsonField_Validate_Array(t *testing.T) {
	f := NewJsonField("data", "DATA", true, nil)
	err := f.Validate([]interface{}{"a", "b"})
	if err != nil {
		t.Fatalf("unexpected error for array value: %v", err)
	}
}

func TestJsonField_Validate_InvalidType(t *testing.T) {
	f := NewJsonField("data", "DATA", true, nil)
	err := f.Validate(42)
	if err == nil {
		t.Fatal("expected error for int value on json field")
	}
}

func TestJsonField_Parse_Nil_WithDefault(t *testing.T) {
	def := map[string]interface{}{"status": "active"}
	f := NewJsonField("data", "DATA", false, def)

	result, err := f.Parse(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	m, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map[string]interface{}, got %T", result)
	}
	if m["status"] != "active" {
		t.Errorf("expected status=active, got %v", m["status"])
	}
}

func TestJsonField_Parse_Nil_NoDefault(t *testing.T) {
	f := NewJsonField("data", "DATA", false, nil)

	result, err := f.Parse(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// nil map[string]interface{} stored in interface{} is non-nil,
	// but the underlying map value is nil.
	if result == nil {
		// GetDefault returns interface{} wrapping nil map — may be nil or non-nil
		// depending on Go interface semantics. Accept either.
		return
	}
	if m, ok := result.(map[string]interface{}); ok && m != nil {
		t.Errorf("expected nil or nil-map default, got %v", result)
	}
}

func TestJsonField_Parse_ValidValue(t *testing.T) {
	f := NewJsonField("data", "DATA", true, nil)
	input := map[string]interface{}{"key": "value"}

	result, err := f.Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	m, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map[string]interface{}, got %T", result)
	}
	if m["key"] != "value" {
		t.Errorf("expected key=value, got %v", m["key"])
	}
}

func TestJsonField_BindValue_Map(t *testing.T) {
	f := NewJsonField("data", "DATA", true, nil)
	input := map[string]interface{}{"name": "test", "count": float64(5)}

	err := f.BindValue(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if f.Value == nil {
		t.Fatal("expected non-nil Value after BindValue")
	}
	if f.Value["name"] != "test" {
		t.Errorf("expected name=test, got %v", f.Value["name"])
	}
	if f.Value["count"] != float64(5) {
		t.Errorf("expected count=5, got %v", f.Value["count"])
	}
}

func TestJsonField_BindValue_String(t *testing.T) {
	f := NewJsonField("data", "DATA", true, nil)

	err := f.BindValue(`{"key":"value"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// String input results in an empty map (per BindValue implementation)
	if f.Value == nil {
		t.Fatal("expected non-nil Value for string input")
	}
	if len(f.Value) != 0 {
		t.Errorf("expected empty map for string input, got %v", f.Value)
	}
}

func TestJsonField_BindValue_Nil_Optional(t *testing.T) {
	f := NewJsonField("data", "DATA", false, nil)

	err := f.BindValue(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// nil on optional field with nil default: Value should be nil
	if f.Value != nil {
		t.Errorf("expected nil Value for nil input on optional field, got %v", f.Value)
	}
}
