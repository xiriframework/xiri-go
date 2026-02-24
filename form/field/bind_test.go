package field

import (
	"testing"
)

func TestTextFieldBindValue(t *testing.T) {
	f := NewTextField("name", "NAME", true, "")
	if err := f.BindValue("hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Value == nil || *f.Value != "hello" {
		t.Errorf("expected 'hello', got %v", f.Value)
	}
}

func TestTextFieldBindValue_Nil(t *testing.T) {
	f := NewTextField("name", "NAME", false, "")
	if err := f.BindValue(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Default is "" so BindValue parses nil → default "" → *string{""}
	if f.Value == nil {
		t.Error("expected non-nil value (default empty string)")
	}
}

func TestTextFieldBindValue_RequiredNil(t *testing.T) {
	f := NewTextField("name", "NAME", true, "")
	err := f.BindValue(nil)
	// nil with default "" should bind to ""
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTextFieldBindValue_MinLength(t *testing.T) {
	f := NewTextFieldWithLength("name", "NAME", true, "", 3, 10)
	err := f.BindValue("ab")
	if err == nil {
		t.Fatal("expected min length error")
	}
}

func TestTextFieldBindValue_MaxLength(t *testing.T) {
	f := NewTextFieldWithLength("name", "NAME", true, "", 0, 5)
	err := f.BindValue("toolongstring")
	if err == nil {
		t.Fatal("expected max length error")
	}
}

func TestBoolFieldBindValue(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected bool
	}{
		{"bool true", true, true},
		{"bool false", false, false},
		{"string true", "true", true},
		{"string 1", "1", true},
		{"string yes", "yes", true},
		{"string false", "false", false},
		{"int 1", 1, true},
		{"int 0", 0, false},
		{"float 1", float64(1), true},
		{"float 0", float64(0), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewBoolField("active", "ACTIVE", false, false)
			if err := f.BindValue(tt.input); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if f.Value == nil || *f.Value != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, f.Value)
			}
		})
	}
}

func TestBoolFieldBindValue_Nil(t *testing.T) {
	f := NewBoolField("active", "ACTIVE", false, true)
	if err := f.BindValue(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should use default value (true)
	if f.Value == nil || *f.Value != true {
		t.Errorf("expected default true, got %v", f.Value)
	}
}

func TestIntFieldBindValue(t *testing.T) {
	f := NewIntField("count", "COUNT", true, 0)
	if err := f.BindValue(float64(42)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Value == nil || *f.Value != 42 {
		t.Errorf("expected 42, got %v", f.Value)
	}
}

func TestIntFieldBindValue_String(t *testing.T) {
	f := NewIntField("count", "COUNT", true, 0)
	if err := f.BindValue("123"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Value == nil || *f.Value != 123 {
		t.Errorf("expected 123, got %v", f.Value)
	}
}

func TestIntFieldBindValue_Bounds(t *testing.T) {
	f := NewIntFieldWithBounds("count", "COUNT", true, 0, 1, 100)

	// Below min
	err := f.BindValue(float64(0))
	if err == nil {
		t.Fatal("expected min bound error")
	}

	// Above max
	err = f.BindValue(float64(101))
	if err == nil {
		t.Fatal("expected max bound error")
	}

	// Within bounds
	err = f.BindValue(float64(50))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSelectFieldBindValue(t *testing.T) {
	opts := []SelectOption{
		{Value: int32(1), Label: "Option A"},
		{Value: int32(2), Label: "Option B"},
	}
	f := NewSelectField("choice", "CHOICE", true, opts)

	if err := f.BindValue(float64(2)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Value != 2 {
		t.Errorf("expected 2, got %v", f.Value)
	}
}

func TestSelectFieldBindValue_InvalidOption(t *testing.T) {
	opts := []SelectOption{
		{Value: int32(1), Label: "Option A"},
	}
	f := NewSelectField("choice", "CHOICE", true, opts)

	err := f.BindValue(float64(99))
	if err == nil {
		t.Fatal("expected error for invalid option")
	}
}

func TestModelFieldBindValue(t *testing.T) {
	f := NewModelField("device", "DEVICE", true, "Device", int32(0))

	if err := f.BindValue(float64(42)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Value != 42 {
		t.Errorf("expected 42, got %v", f.Value)
	}
}

func TestModelFieldBindValue_String(t *testing.T) {
	f := NewModelField("device", "DEVICE", true, "Device", int32(0))

	if err := f.BindValue("123"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Value != 123 {
		t.Errorf("expected 123, got %v", f.Value)
	}
}

func TestModelFieldBindValue_Nil(t *testing.T) {
	f := NewModelField("device", "DEVICE", false, "Device", int32(5))

	if err := f.BindValue(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Value != 5 {
		t.Errorf("expected default 5, got %v", f.Value)
	}
}

func TestFileFieldValidate(t *testing.T) {
	f := NewFileField("upload", "UPLOAD", true, 1024*1024)

	// Required field with nil value
	if err := f.Validate(nil); err == nil {
		t.Fatal("expected error for nil on required file field")
	}

	// Non-nil value should pass
	if err := f.Validate("somefile.pdf"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFileFieldValidate_Optional(t *testing.T) {
	f := NewFileField("upload", "UPLOAD", false, 1024*1024)

	if err := f.Validate(nil); err != nil {
		t.Fatalf("unexpected error for nil on optional file field: %v", err)
	}
}
