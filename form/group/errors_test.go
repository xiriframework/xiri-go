package group

import "testing"

func TestFieldErrors_Error_SortedJoined(t *testing.T) {
	err := FieldErrors{"b": "zu kurz", "a": "fehlt"}
	if got := err.Error(); got != "a: fehlt; b: zu kurz" {
		t.Fatalf("unexpected message: %q", got)
	}
	if got := err.FieldErrors(); got["a"] != "fehlt" {
		t.Fatalf("FieldErrors() must expose the map, got %v", got)
	}
}
