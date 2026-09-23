package field

import (
	"reflect"
	"testing"

	"github.com/xiriframework/xiri-go/component/url"
)

func TestTextField_Export_Suggestions(t *testing.T) {
	f := NewTextField("city", "Stadt", false, "").SetSuggestions("Wien", "Graz")
	got := f.ExportForFrontend(nil, nil)["list"]
	want := []map[string]interface{}{
		{"id": "Wien", "name": "Wien"},
		{"id": "Graz", "name": "Graz"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected list %v, got %#v", want, got)
	}
}

func TestTextField_Export_SuggestionsURL(t *testing.T) {
	u := url.NewUrlPrefix("/City/Search", "/api")
	f := NewTextField("city", "Stadt", false, "").SetSuggestionsURL(u)
	export := f.ExportForFrontend(nil, nil)
	if got := export["url"]; got != u.PrintPrefix() {
		t.Errorf("expected url %s, got %#v", u.PrintPrefix(), got)
	}
	if _, ok := export["searchWith"]; ok {
		t.Errorf("searchWith must be absent without context fields")
	}
}

func TestTextField_Export_SearchWith(t *testing.T) {
	u := url.NewUrlPrefix("/City/Search", "/api")
	f := NewTextField("city", "Stadt", false, "").SetSuggestionsURL(u, "dept", "search", "kind")
	got := f.ExportForFrontend(nil, nil)["searchWith"]
	// "search" ist der Suchtext selbst und wird als Kontextfeld ignoriert.
	if !reflect.DeepEqual(got, []string{"dept", "kind"}) {
		t.Errorf("expected searchWith [dept kind], got %#v", got)
	}
}

func TestTextField_Export_NilURL_Ignored(t *testing.T) {
	export := NewTextField("city", "Stadt", false, "").SetSuggestionsURL(nil, "dept").ExportForFrontend(nil, nil)
	for _, key := range []string{"url", "searchWith"} {
		if _, ok := export[key]; ok {
			t.Errorf("%s must be absent for nil url", key)
		}
	}
}

func TestTextField_Export_NoSuggestions_NoKeys(t *testing.T) {
	export := NewTextField("city", "Stadt", false, "").ExportForFrontend(nil, nil)
	for _, key := range []string{"list", "url", "searchWith"} {
		if _, ok := export[key]; ok {
			t.Errorf("%s must be absent without suggestions", key)
		}
	}
}

// Explizit leer muss als list: [] ankommen, sonst kann ein Reload-Patch Vorschläge nie abräumen.
func TestTextField_Export_EmptySuggestions_EmptyList(t *testing.T) {
	f := NewTextField("city", "Stadt", false, "").SetSuggestions()
	got, ok := f.ExportForFrontend(nil, nil)["list"].([]map[string]interface{})
	if !ok || len(got) != 0 {
		t.Errorf("expected empty list, got %#v", f.ExportForFrontend(nil, nil)["list"])
	}
}

func TestTextField_Suggestions_FreeTextStillBinds(t *testing.T) {
	f := NewTextField("city", "Stadt", false, "").SetSuggestions("Wien")
	if err := f.BindValue("Klagenfurt"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Value == nil || *f.Value != "Klagenfurt" {
		t.Errorf("expected Klagenfurt, got %v", f.Value)
	}
}

func TestTextField_Validate_RequiredRejectsBlank(t *testing.T) {
	f := NewTextField("name", "Name", true, "")
	for _, v := range []string{"", "   "} {
		if err := f.Validate(v); err == nil {
			t.Errorf("Validate(%q) on required field: want error, got nil", v)
		}
	}
	if err := NewTextField("name", "Name", false, "").Validate(""); err != nil {
		t.Errorf("Validate(\"\") on optional field: want nil, got %v", err)
	}
}
