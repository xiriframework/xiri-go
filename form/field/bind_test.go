package field

import (
	"math"
	"strconv"
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

// TestIntFieldBindValue_Overflow verifies that float64 values exceeding int32
// range are rejected to prevent integer overflow (M4).
func TestIntFieldBindValue_Overflow(t *testing.T) {
	tests := []struct {
		name      string
		input     float64
		expectErr bool
	}{
		{"max int32", math.MaxInt32, false},
		{"min int32", math.MinInt32, false},
		{"above max int32", math.MaxInt32 + 1, true},
		{"below min int32", math.MinInt32 - 1, true},
		{"large positive", 1e18, true},
		{"large negative", -1e18, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewIntField("count", "COUNT", false, 0)
			err := f.BindValue(tt.input)
			if tt.expectErr && err == nil {
				t.Errorf("input=%v: expected error, got nil", tt.input)
			}
			if !tt.expectErr && err != nil {
				t.Errorf("input=%v: unexpected error: %v", tt.input, err)
			}
		})
	}
}

// TestIntFieldBindValue_Lossless verifies that every input type is converted to
// int32 without silent truncation or wrap-around (#5). The float64 range case is
// covered by TestIntFieldBindValue_Overflow; this covers fractions plus the
// string and int64 branches, which had no range check at all.
func TestIntFieldBindValue_Lossless(t *testing.T) {
	tests := []struct {
		name      string
		input     any
		expectErr bool
		want      int32
	}{
		{"fractional float", 1.9, true, 0},
		{"integer float", 42.0, false, 42},
		{"string out of int32 range", "3000000000", true, 0},
		{"string fractional", "1.9", true, 0},
		{"string integer", "42", false, 42},
		{"int64 out of int32 range", int64(1) << 40, true, 0},
		{"int out of int32 range", 1 << 40, true, 0},
		{"int64 in range", int64(-7), false, -7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewIntField("count", "COUNT", false, 0)
			err := f.BindValue(tt.input)
			if tt.expectErr {
				if err == nil {
					t.Fatalf("input=%v: expected error, got nil (value=%v)", tt.input, intPtrString(f.Value))
				}
				return
			}
			if err != nil {
				t.Fatalf("input=%v: unexpected error: %v", tt.input, err)
			}
			if f.Value == nil || *f.Value != tt.want {
				t.Errorf("input=%v: expected %d, got %v", tt.input, tt.want, f.Value)
			}
		})
	}
}

// TestIntFieldValidate_Int32Range verifies Validate rejects values that cannot be
// stored in the field's *int32, for callers that use Validate directly instead of
// going through Parse (#5).
func TestIntFieldValidate_Int32Range(t *testing.T) {
	f := NewIntField("count", "COUNT", false, 0)

	if err := f.Validate(int64(1) << 40); err == nil {
		t.Error("expected error for int64 out of int32 range, got nil")
	}
	if err := f.Validate(1 << 40); err == nil {
		t.Error("expected error for int out of int32 range, got nil")
	}
	if err := f.Validate(int64(42)); err != nil {
		t.Errorf("unexpected error for in-range value: %v", err)
	}
}

// TestModelFieldBindValue_IntOverflow covers the int/int64 branches of
// ModelField.Parse, which cast to int32 without a range check (#5). The float64
// branch is covered by TestModelFieldBindValue_Overflow.
func TestModelFieldBindValue_IntOverflow(t *testing.T) {
	tests := []struct {
		name      string
		input     any
		expectErr bool
	}{
		{"int out of int32 range", 1 << 40, true},
		{"int64 out of int32 range", int64(1) << 40, true},
		{"int in range", 42, false},
		{"int64 in range", int64(42), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewModelField("device", "DEVICE", false, "Device", int32(0))
			err := f.BindValue(tt.input)
			if tt.expectErr && err == nil {
				t.Errorf("input=%v: expected error, got nil (value=%v)", tt.input, f.Value)
			}
			if !tt.expectErr && err != nil {
				t.Errorf("input=%v: unexpected error: %v", tt.input, err)
			}
		})
	}
}

// TestModelListBindValue_Lossless verifies that model IDs in a list are not
// silently truncated or wrapped (#5) — a fractional ID would otherwise select a
// different entity.
func TestModelListBindValue_Lossless(t *testing.T) {
	tests := []struct {
		name      string
		input     any
		expectErr bool
		want      []int32
	}{
		{"fractional id", []any{1.9}, true, nil},
		{"float id out of int32 range", []any{float64(int64(1) << 40)}, true, nil},
		{"int id out of int32 range", []any{1 << 40}, true, nil},
		{"integer floats", []any{1.0, 2.0}, false, []int32{1, 2}},
		{"comma string", "1,2", false, []int32{1, 2}},
		{"comma string out of range", "3000000000", true, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewModelListField("devices", "DEVICES", false, "Device", nil)
			err := f.BindValue(tt.input)
			if tt.expectErr {
				if err == nil {
					t.Fatalf("input=%v: expected error, got nil (values=%v)", tt.input, f.Value)
				}
				return
			}
			if err != nil {
				t.Fatalf("input=%v: unexpected error: %v", tt.input, err)
			}
			if len(f.Value) != len(tt.want) {
				t.Fatalf("input=%v: expected %v, got %v", tt.input, tt.want, f.Value)
			}
			for i, want := range tt.want {
				if f.Value[i] != want {
					t.Errorf("input=%v: index %d expected %d, got %d", tt.input, i, want, f.Value[i])
				}
			}
		})
	}
}

// TestNewNumberField_RejectsFractionalDefault verifies the constructor reports a
// default that cannot be stored in the field's *int32 instead of truncating it (#5).
func TestNewNumberField_RejectsFractionalDefault(t *testing.T) {
	if _, err := NewNumberField("price", "PRICE", false, 1.9); err == nil {
		t.Error("expected error for fractional default 1.9, got nil")
	}
	if _, err := NewNumberField("price", "PRICE", false, math.MaxInt32+1); err == nil {
		t.Error("expected error for out-of-range default, got nil")
	}
	f, err := NewNumberField("price", "PRICE", false, 42.0)
	if err != nil {
		t.Fatalf("unexpected error for integer default: %v", err)
	}
	if f.GetDefault() != 42 {
		t.Errorf("expected default 42, got %v", f.GetDefault())
	}
}

// TestModelFieldBindValue_Overflow verifies that float64 values exceeding int32
// range or non-integer floats are rejected for model IDs (M4).
func TestModelFieldBindValue_Overflow(t *testing.T) {
	tests := []struct {
		name      string
		input     float64
		expectErr bool
	}{
		{"max int32", math.MaxInt32, false},
		{"min int32", math.MinInt32, false},
		{"above max int32", math.MaxInt32 + 1, true},
		{"below min int32", math.MinInt32 - 1, true},
		{"large positive", 1e18, true},
		{"non-integer", 42.5, true},
		{"integer float", 42.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewModelField("device", "DEVICE", false, "Device", int32(0))
			err := f.BindValue(tt.input)
			if tt.expectErr && err == nil {
				t.Errorf("input=%v: expected error, got nil", tt.input)
			}
			if !tt.expectErr && err != nil {
				t.Errorf("input=%v: unexpected error: %v", tt.input, err)
			}
		})
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

// intPtrString renders a *int32 for test failure messages.
func intPtrString(v *int32) string {
	if v == nil {
		return "<nil>"
	}
	return strconv.FormatInt(int64(*v), 10)
}
