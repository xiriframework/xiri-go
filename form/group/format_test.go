package group

import "testing"

// #17: FormatNumber must not return the format template when ctx is nil.
func TestFormatNumberNilContext(t *testing.T) {
	fg := NewFormGroup(nil)

	got := fg.FormatNumber(12.34, 2)
	if got != "12.34" {
		t.Errorf("FormatNumber(12.34, 2) with nil ctx = %q, want %q", got, "12.34")
	}
}
