package builder

import (
	"testing"

	"github.com/xiriframework/xiri-go/form/field"
	"github.com/xiriframework/xiri-go/form/group"
)

func reloadSelectOptions() []field.SelectOption {
	return []field.SelectOption{
		{Value: int32(1), Label: "one"},
		{Value: int32(2), Label: "two"},
	}
}

// A reload happens mid-edit, so a required field that is still empty is normal input -
// not an error. The counter-check pins that BindFromMap really does reject the same data.
//
// A multi-select is the honest case here: NewSelectField seeds a single select with the
// first option as its default, so a missing value never fails validation there.
func TestBindReloadFromMap_ToleratesMissingRequired(t *testing.T) {
	data := map[string]interface{}{}

	strict := field.NewSelectField("tags", "TAGS", true, reloadSelectOptions()).SetMultiple(true)
	if err := BindFromMap(data, group.NewFormGroup([]field.FormField{strict})); err == nil {
		t.Fatal("counter-check failed: BindFromMap accepted a missing required field")
	}

	lenient := field.NewSelectField("tags", "TAGS", true, reloadSelectOptions()).SetMultiple(true)
	if err := BindReloadFromMap(data, group.NewFormGroup([]field.FormField{lenient})); err != nil {
		t.Errorf("BindReloadFromMap rejected a missing required field: %v", err)
	}
}

// Binding is lenient, not skipped: a field whose value could not be bound must still hold
// its default, otherwise the business logic reads a zero value instead.
func TestBindReloadFromMap_KeepsDefaultOnError(t *testing.T) {
	name := field.NewTextFieldWithLength("name", "NAME", false, "fallback", 5, 100)

	err := BindReloadFromMap(map[string]interface{}{"name": "ab"}, group.NewFormGroup([]field.FormField{name}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if name.Value == nil || *name.Value != "fallback" {
		t.Errorf("value=%v want the default %q", name.Value, "fallback")
	}
}

func TestBindReloadFromMap_BindsValidValues(t *testing.T) {
	status := field.NewSelectField("status", "STATUS", false, reloadSelectOptions())
	note := field.NewTextField("note", "NOTE", false, "")

	err := BindReloadFromMap(map[string]interface{}{
		"status": float64(2),
		"note":   "hello",
	}, group.NewFormGroup([]field.FormField{status, note}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if status.Value != 2 {
		t.Errorf("status=%v want 2", status.Value)
	}
	if note.Value == nil || *note.Value != "hello" {
		t.Errorf("note=%v want hello", note.Value)
	}
}

// One unusable value must not stop the fields after it from binding - the trigger the
// server needs could be any of them.
func TestBindReloadFromMap_BadValueDoesNotStopLaterFields(t *testing.T) {
	broken := field.NewSelectField("broken", "BROKEN", false, reloadSelectOptions())
	status := field.NewSelectField("status", "STATUS", false, reloadSelectOptions())

	err := BindReloadFromMap(map[string]interface{}{
		"broken": float64(99), // not a valid option
		"status": float64(1),
	}, group.NewFormGroup([]field.FormField{broken, status}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if status.Value != 1 {
		t.Errorf("status=%v want 1 - a later field stopped binding", status.Value)
	}
}
