package field

import (
	"reflect"
	"testing"
)

func selectTestOptions() []SelectOption {
	return []SelectOption{
		{Value: int32(1), Label: "status.active"},
		{Value: int32(2), Label: "status.locked"},
		{Value: int32(3), Label: "status.expired"},
	}
}

func TestSelectField_Single_Unchanged(t *testing.T) {
	f := NewSelectField("status", "Status", false, selectTestOptions())
	if err := f.BindValue(float64(2)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Value != 2 {
		t.Errorf("expected Value=2, got %d", f.Value)
	}
	if f.Multiple {
		t.Error("expected Multiple=false by default")
	}
}

func TestSelectField_Single_DefaultIsFirstOption(t *testing.T) {
	f := NewSelectField("status", "Status", false, selectTestOptions())
	def, ok := f.GetDefault().(int32)
	if !ok || def != 1 {
		t.Errorf("expected default int32(1), got %T %v", f.GetDefault(), f.GetDefault())
	}
}

func TestSelectField_Multi_BindSlice(t *testing.T) {
	f := NewSelectField("status", "Status", false, selectTestOptions()).SetMultiple(true)
	if err := f.BindValue([]interface{}{float64(1), float64(3)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(f.Values, []int32{1, 3}) {
		t.Errorf("expected [1 3], got %v", f.Values)
	}
}

func TestSelectField_Multi_BindNil_NotRequired(t *testing.T) {
	f := NewSelectField("status", "Status", false, selectTestOptions()).SetMultiple(true)
	if err := f.BindValue(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.Values) != 0 {
		t.Errorf("expected empty slice, got %v", f.Values)
	}
}

func TestSelectField_Multi_BindNil_Required(t *testing.T) {
	f := NewSelectField("status", "Status", true, selectTestOptions()).SetMultiple(true)
	if err := f.BindValue(nil); err == nil {
		t.Error("expected required error for empty multi-select")
	}
}

func TestSelectField_Multi_InvalidValue(t *testing.T) {
	f := NewSelectField("status", "Status", false, selectTestOptions()).SetMultiple(true)
	if err := f.BindValue([]interface{}{float64(99)}); err == nil {
		t.Error("expected error for unknown option id")
	}
}

func TestSelectField_Multi_ExportContainsMultiple(t *testing.T) {
	f := NewSelectField("status", "Status", false, selectTestOptions()).SetMultiple(true)
	_ = f.BindValue([]interface{}{float64(2)})
	out := f.ExportForFrontend(nil, f.Values)
	if m, ok := out["multiple"].(bool); !ok || !m {
		t.Errorf("expected multiple=true in export, got %v", out["multiple"])
	}
	// value must be array-like (ModelListValue / []int32)
	switch out["value"].(type) {
	case ModelListValue, []int32:
		// ok
	default:
		t.Errorf("expected array value, got %T", out["value"])
	}
}

func TestSelectField_Single_ExportHasNoMultiple(t *testing.T) {
	f := NewSelectField("status", "Status", false, selectTestOptions())
	out := f.ExportForFrontend(nil, int32(1))
	if _, present := out["multiple"]; present {
		t.Error("expected no multiple key in single-select export")
	}
}
