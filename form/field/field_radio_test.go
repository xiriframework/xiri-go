package field

import "testing"

func TestRadioField_ExportType(t *testing.T) {
	f := NewRadioField("anrede", "Anrede", true, selectTestOptions())

	out := f.ExportForFrontend(nil, nil)
	if out["type"] != "radio" {
		t.Errorf("type=%v want radio", out["type"])
	}
	list, ok := out["list"].([]map[string]interface{})
	if !ok || len(list) != 3 {
		t.Fatalf("list=%v want 3 options", out["list"])
	}
	if list[0]["name"] != "status.active" {
		t.Errorf("list[0].name=%v", list[0]["name"])
	}
	if _, ok := out["multiple"]; ok {
		t.Error("radio must not export multiple")
	}
}

func TestRadioField_BindSingle(t *testing.T) {
	f := NewRadioField("anrede", "Anrede", true, selectTestOptions())
	if err := f.BindValue(float64(2)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Value != 2 {
		t.Errorf("value=%v want 2", f.Value)
	}
}
