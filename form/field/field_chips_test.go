package field

import (
	"reflect"
	"testing"
)

func chipsTestOptions() []SelectOption {
	return []SelectOption{
		{Value: int64(1), Label: "tag.red"},
		{Value: int64(2), Label: "tag.green"},
		{Value: int64(3), Label: "tag.blue"},
	}
}

func TestChipsField_Parse_Mixed(t *testing.T) {
	f := NewChipsField("tags", "Tags", false).SetList(chipsTestOptions())
	parsed, err := f.Parse([]interface{}{float64(1), "custom", float64(3)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(parsed, []interface{}{int64(1), "custom", int64(3)}) {
		t.Errorf("expected [1 custom 3] (int64+string), got %#v", parsed)
	}
}

func TestChipsField_Parse_StringSlice(t *testing.T) {
	f := NewChipsField("tags", "Tags", false)
	parsed, err := f.Parse([]string{"a", "b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(parsed, []interface{}{"a", "b"}) {
		t.Errorf("expected [a b], got %#v", parsed)
	}
}

func TestChipsField_FreeTextOnly_NoList(t *testing.T) {
	f := NewChipsField("tags", "Tags", false)
	if err := f.BindValue([]interface{}{"alpha", "beta"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(f.Value, []interface{}{"alpha", "beta"}) {
		t.Errorf("expected [alpha beta], got %#v", f.Value)
	}
	if len(f.IDs()) != 0 {
		t.Errorf("expected no IDs, got %v", f.IDs())
	}
	if !reflect.DeepEqual(f.Texts(), []string{"alpha", "beta"}) {
		t.Errorf("expected texts [alpha beta], got %v", f.Texts())
	}
}

func TestChipsField_BindValue_IDsAndTexts(t *testing.T) {
	f := NewChipsField("tags", "Tags", false).SetList(chipsTestOptions())
	if err := f.BindValue([]interface{}{float64(2), "custom", float64(1)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(f.IDs(), []int64{2, 1}) {
		t.Errorf("expected IDs [2 1], got %v", f.IDs())
	}
	if !reflect.DeepEqual(f.Texts(), []string{"custom"}) {
		t.Errorf("expected texts [custom], got %v", f.Texts())
	}
}

func TestChipsField_Validate_UnknownID(t *testing.T) {
	f := NewChipsField("tags", "Tags", false).SetList(chipsTestOptions())
	if err := f.BindValue([]interface{}{float64(99)}); err == nil {
		t.Error("expected error for unknown option id 99")
	}
}

func TestChipsField_Validate_FreeTextAlwaysValid(t *testing.T) {
	f := NewChipsField("tags", "Tags", false).SetList(chipsTestOptions())
	if err := f.BindValue([]interface{}{"anything-goes"}); err != nil {
		t.Fatalf("free text should be valid, got: %v", err)
	}
}

func TestChipsField_Required_Empty(t *testing.T) {
	f := NewChipsField("tags", "Tags", true).SetList(chipsTestOptions())
	if err := f.BindValue([]interface{}{}); err == nil {
		t.Error("expected required error for empty chips")
	}
}

func TestChipsField_Export_MixedValue(t *testing.T) {
	ctx := newTestCtx()
	f := NewChipsField("tags", "Tags", false).SetList(chipsTestOptions())
	if err := f.BindValue([]interface{}{float64(1), "custom"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := f.ExportForFrontend(ctx, f.Value)

	if !reflect.DeepEqual(out["value"], []interface{}{int64(1), "custom"}) {
		t.Errorf("expected value [1 custom], got %#v", out["value"])
	}
	list, ok := out["list"].([]map[string]interface{})
	if !ok || len(list) != 3 {
		t.Fatalf("expected list with 3 options, got %#v", out["list"])
	}
	if list[0]["id"] != int64(1) {
		t.Errorf("expected first option id int64(1), got %#v", list[0]["id"])
	}
}

func TestChipsField_Export_DefaultEmpty(t *testing.T) {
	ctx := newTestCtx()
	f := NewChipsField("tags", "Tags", false)
	out := f.ExportForFrontend(ctx, nil)
	if !reflect.DeepEqual(out["value"], []interface{}{}) {
		t.Errorf("expected empty default value, got %#v", out["value"])
	}
}
