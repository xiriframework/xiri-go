package field

import (
	"reflect"
	"testing"

	"github.com/xiriframework/xiri-go/component/url"
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

func chipsStringOptions() []SelectOption {
	return []SelectOption{
		{Value: "eth0", Label: "Ethernet 0"},
		{Value: "eth1", Label: "Ethernet 1"},
	}
}

// Eine String-Option kommt als String zurück und darf nicht als Freitext gelten.
func TestChipsField_StringOption_NotInTexts(t *testing.T) {
	f := NewChipsField("ifaces", "Interfaces", false).SetList(chipsStringOptions())
	if err := f.BindValue([]interface{}{"eth0", "custom"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(f.Texts(), []string{"custom"}) {
		t.Errorf("expected texts [custom], got %v", f.Texts())
	}
}

func TestChipsField_StrIDs_StringOptions(t *testing.T) {
	list := append(chipsStringOptions(), SelectOption{Value: int64(5), Label: "Fünf"})
	f := NewChipsField("ifaces", "Interfaces", false).SetList(list)
	if err := f.BindValue([]interface{}{"eth1", "custom", float64(5), "5", "eth0"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(f.StrIDs(), []string{"eth1", "eth0"}) {
		t.Errorf("expected StrIDs [eth1 eth0], got %v", f.StrIDs())
	}
	if !reflect.DeepEqual(f.IDs(), []int64{5}) {
		t.Errorf("expected IDs [5], got %v", f.IDs())
	}
	// "5" ist Freitext: die numerische Option 5 hat einen anderen Typ.
	if !reflect.DeepEqual(f.Texts(), []string{"custom", "5"}) {
		t.Errorf("expected texts [custom 5], got %v", f.Texts())
	}
}

func TestChipsField_FreeTextFalse_RejectsText(t *testing.T) {
	list := append(chipsStringOptions(), SelectOption{Value: int64(5), Label: "Fünf"})
	f := NewChipsField("ifaces", "Interfaces", false).SetList(list).SetFreeText(false)
	if err := f.BindValue([]interface{}{"eth0", float64(5)}); err != nil {
		t.Fatalf("options only must be valid, got: %v", err)
	}
	err := f.Validate([]interface{}{"eth0", "bogus"})
	if err == nil {
		t.Fatal("expected not_allowed for free text with FreeText=false")
	}
	if c := code(t, err); c != "not_allowed" {
		t.Errorf("expected code not_allowed, got %s", c)
	}
}

// Ohne Request-Wert greift der Default; auch der darf bei FreeText=false keinen Freitext enthalten.
func TestChipsField_FreeTextFalse_RejectsTextInDefault(t *testing.T) {
	f := NewChipsField("ifaces", "Interfaces", false).SetList(chipsStringOptions()).SetFreeText(false)
	f.Default = []string{"bogus"}
	if err := f.BindValue(nil); err == nil {
		t.Error("expected error for free-text default with FreeText=false")
	}
}

func TestChipsField_Export_FreeText(t *testing.T) {
	if _, ok := NewChipsField("t", "T", false).ExportForFrontend(nil, nil)["freeText"]; ok {
		t.Error("freeText must be absent by default (frontend default is true)")
	}
	out := NewChipsField("t", "T", false).SetFreeText(false).ExportForFrontend(nil, nil)
	if got, ok := out["freeText"]; !ok || got != false {
		t.Errorf("expected freeText false, got %#v (present=%v)", got, ok)
	}
}

func TestChipsField_Export_SuggestionsURL(t *testing.T) {
	u := url.NewUrlPrefix("/Iface/Search", "/api")
	out := NewChipsField("t", "T", false).SetSuggestionsURL(u, "host", "search").ExportForFrontend(nil, nil)
	if got := out["url"]; got != u.PrintPrefix() {
		t.Errorf("expected url %s, got %#v", u.PrintPrefix(), got)
	}
	// "search" ist der Suchtext selbst und wird als Kontextfeld ignoriert.
	if got := out["searchWith"]; !reflect.DeepEqual(got, []string{"host"}) {
		t.Errorf("expected searchWith [host], got %#v", got)
	}
}

func TestChipsField_Export_SuggestionsURL_NilDisables(t *testing.T) {
	u := url.NewUrlPrefix("/Iface/Search", "/api")
	f := NewChipsField("t", "T", false).SetSuggestionsURL(u, "host").SetSuggestionsURL(nil)
	out := f.ExportForFrontend(nil, nil)
	for _, key := range []string{"url", "searchWith"} {
		if _, ok := out[key]; ok {
			t.Errorf("%s must be absent after SetSuggestionsURL(nil)", key)
		}
	}
}
