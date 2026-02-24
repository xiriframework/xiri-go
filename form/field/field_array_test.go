package field

import (
	"testing"
)

func TestArrayField_Validate_Required(t *testing.T) {
	f := NewArrayField("items", "ITEMS", true, "string", nil)
	err := f.Validate(nil)
	if err == nil {
		t.Fatal("expected error for nil value on required array field")
	}
}

func TestArrayField_Validate_Optional_Nil(t *testing.T) {
	f := NewArrayField("items", "ITEMS", false, "string", nil)
	err := f.Validate(nil)
	if err != nil {
		t.Fatalf("unexpected error for nil on optional array field: %v", err)
	}
}

func TestArrayField_Validate_EmptyArray(t *testing.T) {
	// AllowEmpty=false should reject empty array
	f := NewArrayField("items", "ITEMS", true, "string", nil)
	f.AllowEmpty = false
	err := f.Validate([]interface{}{})
	if err == nil {
		t.Fatal("expected error for empty array with AllowEmpty=false")
	}

	// AllowEmpty=true should accept empty array
	f.AllowEmpty = true
	err = f.Validate([]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error for empty array with AllowEmpty=true: %v", err)
	}
}

func TestArrayField_Validate_MinItems(t *testing.T) {
	f := NewArrayField("items", "ITEMS", true, "string", nil)
	f.AllowEmpty = true
	min := 3
	f.MinItems = &min

	err := f.Validate([]interface{}{"a", "b"})
	if err == nil {
		t.Fatal("expected error for array with fewer items than MinItems")
	}
}

func TestArrayField_Validate_MaxItems(t *testing.T) {
	f := NewArrayField("items", "ITEMS", true, "string", nil)
	f.AllowEmpty = true
	max := 2
	f.MaxItems = &max

	err := f.Validate([]interface{}{"a", "b", "c"})
	if err == nil {
		t.Fatal("expected error for array with more items than MaxItems")
	}
}

func TestArrayField_Validate_ValidArray(t *testing.T) {
	f := NewArrayField("items", "ITEMS", true, "string", nil)
	f.AllowEmpty = true
	min := 1
	max := 5
	f.MinItems = &min
	f.MaxItems = &max

	err := f.Validate([]interface{}{"a", "b", "c"})
	if err != nil {
		t.Fatalf("unexpected error for valid array: %v", err)
	}
}

func TestArrayField_Parse_Nil_WithDefault(t *testing.T) {
	def := []interface{}{"x", "y"}
	f := NewArrayField("items", "ITEMS", false, "string", def)

	result, err := f.Parse(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	arr, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected []interface{}, got %T", result)
	}
	if len(arr) != 2 || arr[0] != "x" || arr[1] != "y" {
		t.Errorf("expected default [x y], got %v", arr)
	}
}

func TestArrayField_Parse_Nil_NoDefault(t *testing.T) {
	f := NewArrayField("items", "ITEMS", false, "string", nil)
	// nil []interface{} stored in interface{} is non-nil, so GetDefault returns it.
	// Parse returns the default, which is []interface{}(nil).
	// If that is non-nil interface, it is returned; otherwise empty slice.
	result, err := f.Parse(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result from Parse(nil)")
	}
}

func TestArrayField_Parse_ValidArray(t *testing.T) {
	f := NewArrayField("items", "ITEMS", true, "string", nil)
	input := []interface{}{"a", "b", "c"}

	result, err := f.Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	arr, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected []interface{}, got %T", result)
	}
	if len(arr) != 3 {
		t.Errorf("expected 3 items, got %d", len(arr))
	}
}

func TestArrayField_Parse_InvalidType(t *testing.T) {
	f := NewArrayField("items", "ITEMS", true, "string", nil)

	_, err := f.Parse("not an array")
	if err == nil {
		t.Fatal("expected error for non-array input")
	}
}

func TestArrayField_BindValue(t *testing.T) {
	f := NewArrayField("items", "ITEMS", true, "string", nil)
	f.AllowEmpty = true

	input := []interface{}{"hello", "world"}
	err := f.BindValue(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if f.Value == nil {
		t.Fatal("expected non-nil Value after BindValue")
	}
	if len(f.Value) != 2 {
		t.Errorf("expected 2 items, got %d", len(f.Value))
	}
	if f.Value[0] != "hello" || f.Value[1] != "world" {
		t.Errorf("expected [hello world], got %v", f.Value)
	}
}
