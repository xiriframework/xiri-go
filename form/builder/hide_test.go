package builder

import (
	"testing"

	"github.com/xiriframework/xiri-go/form/field"
	"github.com/xiriframework/xiri-go/form/group"
)

// Hidden fields (SetHide(true)) must still be parsed from the request body.
func TestBindFromMap_HiddenFieldIsParsed(t *testing.T) {
	hidden := field.NewTextField("token", "TOKEN", true, "")
	hidden.SetHide(true)

	fg := group.NewFormGroup([]field.FormField{hidden})

	data := map[string]interface{}{
		"token": "abc-123",
	}

	if err := BindFromMap(data, fg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if hidden.Value == nil || *hidden.Value != "abc-123" {
		t.Errorf("expected hidden field value 'abc-123', got %v", hidden.Value)
	}
}

// Fields with SetForm(false) must NOT be parsed — they are backend-only.
func TestBindFromMap_FormFalseFieldIsIgnored(t *testing.T) {
	backend := field.NewTextField("internal", "INTERNAL", false, "default")
	backend.SetForm(false)

	fg := group.NewFormGroup([]field.FormField{backend})

	data := map[string]interface{}{
		"internal": "should-be-ignored",
	}

	if err := BindFromMap(data, fg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if backend.Value != nil && *backend.Value == "should-be-ignored" {
		t.Errorf("expected Form=false field to be ignored, but got %q", *backend.Value)
	}
}

// Hidden + Form=false combination: backend-only still wins, value ignored.
func TestBindFromMap_HideAndFormFalse(t *testing.T) {
	f := field.NewTextField("x", "X", false, "")
	f.SetHide(true)
	f.SetForm(false)

	fg := group.NewFormGroup([]field.FormField{f})

	data := map[string]interface{}{
		"x": "submitted",
	}

	if err := BindFromMap(data, fg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if f.Value != nil && *f.Value == "submitted" {
		t.Errorf("expected Form=false to override Hide and ignore submitted value")
	}
}
