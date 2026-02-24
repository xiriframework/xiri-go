package builder

import (
	"testing"

	"github.com/xiriframework/xiri-go/form/field"
	"github.com/xiriframework/xiri-go/form/group"
)

func TestBindFromMap_Basic(t *testing.T) {
	name := field.NewTextField("name", "NAME", true, "")
	active := field.NewBoolField("active", "ACTIVE", false, false)

	fg := group.NewFormGroup([]field.FormField{name, active})

	data := map[string]interface{}{
		"name":   "test-name",
		"active": true,
	}

	if err := BindFromMap(data, fg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if name.Value == nil || *name.Value != "test-name" {
		t.Errorf("expected 'test-name', got %v", name.Value)
	}
	if active.Value == nil || *active.Value != true {
		t.Errorf("expected true, got %v", active.Value)
	}
}

func TestBindFromMap_MissingOptionalField(t *testing.T) {
	name := field.NewTextField("name", "NAME", true, "")
	count := field.NewIntField("count", "COUNT", false, 99)

	fg := group.NewFormGroup([]field.FormField{name, count})

	data := map[string]interface{}{
		"name": "hello",
		// count is missing - should use default
	}

	if err := BindFromMap(data, fg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if name.Value == nil || *name.Value != "hello" {
		t.Errorf("expected 'hello', got %v", name.Value)
	}
	if count.Value == nil || *count.Value != 99 {
		t.Errorf("expected default 99, got %v", count.Value)
	}
}

func TestBindFromMap_ValidationError(t *testing.T) {
	name := field.NewTextFieldWithLength("name", "NAME", true, "", 5, 100)

	fg := group.NewFormGroup([]field.FormField{name})

	data := map[string]interface{}{
		"name": "ab", // too short
	}

	err := BindFromMap(data, fg)
	if err == nil {
		t.Fatal("expected validation error for too-short text")
	}
}

func TestNewFormBuilder_BuildAdd(t *testing.T) {
	name := field.NewTextField("name", "NAME", true, "default-name")
	active := field.NewBoolField("active", "ACTIVE", false, true)

	builder := NewFormBuilder(nil, func(s string) string { return s })
	builder.AddField(name).AddField(active)

	fg, defaults, err := builder.BuildAdd()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(fg.GetFields()) != 2 {
		t.Errorf("expected 2 fields, got %d", len(fg.GetFields()))
	}

	if defaults["name"] != "default-name" {
		t.Errorf("expected default 'default-name', got %v", defaults["name"])
	}
	if defaults["active"] != true {
		t.Errorf("expected default true, got %v", defaults["active"])
	}
}

func TestNewFormBuilder_BuildEdit(t *testing.T) {
	name := field.NewTextField("name", "NAME", true, "current-name")

	builder := NewFormBuilder(nil, func(s string) string { return s })
	builder.AddField(name)

	fg, values, err := builder.BuildEdit()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(fg.GetFields()) != 1 {
		t.Errorf("expected 1 field, got %d", len(fg.GetFields()))
	}

	if values["name"] != "current-name" {
		t.Errorf("expected 'current-name', got %v", values["name"])
	}
}

func TestNewFormBuilder_BuildAddForDisplay(t *testing.T) {
	name := field.NewTextField("name", "NAME", true, "test")

	builder := NewFormBuilder(nil, func(s string) string { return s })
	builder.AddField(name)

	exported, err := builder.BuildAddForDisplay()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(exported) != 1 {
		t.Fatalf("expected 1 exported field, got %d", len(exported))
	}
	if exported[0]["id"] != "name" {
		t.Errorf("expected field id 'name', got %v", exported[0]["id"])
	}
}

func TestNewFormBuilder_OnEditValueCheck(t *testing.T) {
	name := field.NewTextField("name", "NAME", true, "default")

	builder := NewFormBuilder(nil, func(s string) string { return s })
	builder.AddField(name)

	hookCalled := false
	builder.OnEditValueCheck = func(fg *group.FormGroup, values map[string]interface{}) error {
		hookCalled = true
		values["name"] = "modified-by-hook"
		return nil
	}

	_, values, err := builder.BuildEdit()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !hookCalled {
		t.Error("expected OnEditValueCheck hook to be called")
	}
	if values["name"] != "modified-by-hook" {
		t.Errorf("expected 'modified-by-hook', got %v", values["name"])
	}
}
