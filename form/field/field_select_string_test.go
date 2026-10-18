package field

import (
	"encoding/json"
	"reflect"
	"testing"
)

func strOptions() []SelectOption {
	return []SelectOption{
		{Value: "allow", Label: "Erlauben"},
		{Value: "deny", Label: "Verweigern"},
	}
}

func TestSelectField_Str_Single_Bind(t *testing.T) {
	f := NewSelectField("action", "AKTION", true, strOptions())
	if err := f.BindValue("deny"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.StrValue != "deny" || f.Value != 0 {
		t.Errorf("want StrValue=deny Value=0, got %q %d", f.StrValue, f.Value)
	}
}

func TestSelectField_Str_Single_DefaultThenRequest(t *testing.T) {
	f := NewSelectField("action", "AKTION", true, strOptions())
	if err := f.BindValue(nil); err != nil {
		t.Fatalf("default: %v", err)
	}
	if f.StrValue != "allow" {
		t.Errorf("default: want allow, got %q", f.StrValue)
	}
	if err := f.BindValue("deny"); err != nil {
		t.Fatalf("request: %v", err)
	}
	if f.StrValue != "deny" {
		t.Errorf("request: want deny, got %q", f.StrValue)
	}
}

func TestSelectField_Str_Single_Unknown(t *testing.T) {
	f := NewSelectField("action", "AKTION", true, strOptions())
	if err := f.BindValue("maybe"); err == nil {
		t.Error("expected not_allowed for unknown string")
	}
}

func TestSelectField_Str_Single_NoDefault(t *testing.T) {
	f := NewSelectField("action", "AKTION", false, strOptions())
	f.Default = nil
	f.StrValue = "stale"
	if err := f.BindValue(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.StrValue != "" || f.Value != 0 {
		t.Errorf("want zero values, got %q %d", f.StrValue, f.Value)
	}
}

func TestSelectField_Mixed_Single_ResetsOther(t *testing.T) {
	f := NewSelectField("m", "M", true, []SelectOption{{Value: "a", Label: "A"}, {Value: int32(2), Label: "B"}})
	if err := f.BindValue("a"); err != nil {
		t.Fatal(err)
	}
	if err := f.BindValue(float64(2)); err != nil {
		t.Fatal(err)
	}
	if f.Value != 2 || f.StrValue != "" {
		t.Errorf("want Value=2 StrValue=\"\", got %d %q", f.Value, f.StrValue)
	}
	if err := f.BindValue("a"); err != nil {
		t.Fatal(err)
	}
	if f.Value != 0 || f.StrValue != "a" {
		t.Errorf("want Value=0 StrValue=a, got %d %q", f.Value, f.StrValue)
	}
}

func TestSelectField_Str_Multi_Bind(t *testing.T) {
	f := NewSelectField("actions", "AKTIONEN", false, strOptions()).SetMultiple(true)
	if err := f.BindValue([]interface{}{"deny", "allow"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(f.StrValues, []string{"deny", "allow"}) || len(f.Values) != 0 {
		t.Errorf("want StrValues=[deny allow] Values=[], got %v %v", f.StrValues, f.Values)
	}
}

func TestSelectField_Str_Multi_SingleString(t *testing.T) {
	f := NewSelectField("actions", "AKTIONEN", false, strOptions()).SetMultiple(true)
	if err := f.BindValue("allow"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(f.StrValues, []string{"allow"}) {
		t.Errorf("want [allow], got %v", f.StrValues)
	}
	if err := f.BindValue("allow,deny"); err == nil {
		t.Errorf("expected not_allowed for comma string, got %v", f.StrValues)
	}
	if err := f.BindValue(""); err != nil || len(f.StrValues) != 0 {
		t.Errorf("empty string: want empty list, got %v %v", f.StrValues, err)
	}
}

func TestSelectField_Str_Multi_Invalid(t *testing.T) {
	f := NewSelectField("actions", "AKTIONEN", false, strOptions()).SetMultiple(true)
	for name, raw := range map[string]interface{}{
		"unknown": []interface{}{"allow", "maybe"},
		"number":  []interface{}{float64(1)},
	} {
		if err := f.BindValue(raw); err == nil {
			t.Errorf("%s: expected error, got %v", name, f.StrValues)
		}
	}
}

func TestSelectField_Mixed_Multi_NumberNotAllowed(t *testing.T) {
	f := NewSelectField("m", "M", false, []SelectOption{{Value: "a", Label: "A"}, {Value: int32(2), Label: "B"}}).
		SetMultiple(true)
	if err := f.BindValue([]interface{}{float64(2)}); err == nil {
		t.Errorf("expected not_allowed for number in string mode, got %v", f.Values)
	}
}

func TestSelectField_Str_Multi_Empty(t *testing.T) {
	opt := NewSelectField("actions", "AKTIONEN", false, strOptions()).SetMultiple(true)
	opt.StrValues = []string{"stale"}
	if err := opt.BindValue(nil); err != nil {
		t.Fatalf("not required: %v", err)
	}
	if opt.StrValues == nil || len(opt.StrValues) != 0 {
		t.Errorf("want empty non-nil StrValues, got %#v", opt.StrValues)
	}

	req := NewSelectField("actions", "AKTIONEN", true, strOptions()).SetMultiple(true)
	for name, raw := range map[string]interface{}{"nil": nil, "empty": []interface{}{}} {
		if err := req.BindValue(raw); err == nil {
			t.Errorf("%s: expected required error", name)
		}
	}
}

func TestSelectField_Str_Multi_StringSliceDefault(t *testing.T) {
	f := NewSelectField("actions", "AKTIONEN", false, strOptions()).SetMultiple(true)
	f.Default = []string{"deny"}
	if err := f.BindValue(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(f.StrValues, []string{"deny"}) {
		t.Errorf("want [deny], got %v", f.StrValues)
	}
}

func TestSelectField_Int_Multi_ResetsStrValues(t *testing.T) {
	f := NewSelectField("status", "Status", false, selectTestOptions()).SetMultiple(true)
	f.StrValues = []string{"stale"}
	if err := f.BindValue([]interface{}{float64(1)}); err != nil {
		t.Fatal(err)
	}
	if len(f.StrValues) != 0 {
		t.Errorf("want empty StrValues, got %v", f.StrValues)
	}
}

// Nachweis: der Export reicht []string unverändert durch.
func TestSelectField_Str_Multi_Export(t *testing.T) {
	f := NewSelectField("actions", "AKTIONEN", false, strOptions()).SetMultiple(true)
	for _, tc := range []struct {
		value interface{}
		want  string
	}{{nil, `[]`}, {[]string{"allow"}, `["allow"]`}} {
		b, err := json.Marshal(f.ExportForFrontend(nil, tc.value)["value"])
		if err != nil || string(b) != tc.want {
			t.Errorf("value %v: want %s, got %s (%v)", tc.value, tc.want, b, err)
		}
	}
}

func TestRowsField_StrMultiSelectCell(t *testing.T) {
	name := NewTextField("name", "NAME", false, "")
	tags := NewSelectField("tags", "TAGS", false, strOptions()).SetMultiple(true)
	f := NewRowsField("r", "R", false, []FormField{name, tags}).SetUnique("tags")

	if err := f.BindValue([]any{
		map[string]any{"name": "a", "tags": []any{"deny"}},
		map[string]any{"name": "b", "tags": []any{}},
		map[string]any{"name": "c", "tags": []any{}},
		map[string]any{"name": "", "tags": []any{}},
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.Value) != 3 {
		t.Fatalf("want 3 rows (empty one dropped), got %d: %v", len(f.Value), f.Value)
	}
	if !reflect.DeepEqual(f.Value[0]["tags"], []string{"deny"}) {
		t.Errorf("want []string{deny}, got %#v", f.Value[0]["tags"])
	}

	if err := f.BindValue([]any{map[string]any{"name": "a", "tags": []any{"maybe"}}}); err == nil {
		t.Error("expected not_allowed for unknown tag")
	}
}

// Eine versteckte Zelle zählt nur leer, wenn isBlank sie erkennt (kein Default-Vergleich).
func TestRowsField_HiddenEmptyStrMultiSelect_DropsRow(t *testing.T) {
	name := NewTextField("name", "NAME", false, "")
	tags := NewSelectField("tags", "TAGS", false, strOptions()).SetMultiple(true)
	tags.BaseField.SetHide(true)
	f := NewRowsField("r", "R", false, []FormField{name, tags})

	if err := f.BindValue([]any{map[string]any{"name": "", "tags": []any{}}}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.Value) != 0 {
		t.Errorf("want empty row dropped, got %v", f.Value)
	}
}
