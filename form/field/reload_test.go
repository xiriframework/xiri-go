package field

import (
	"reflect"
	"testing"

	"github.com/xiriframework/xiri-go/component/url"
)

func TestSetReloadOn_Export(t *testing.T) {
	f := NewSelectField("tags", "Tags", false, selectTestOptions())
	f.BaseField.SetReloadOn(url.NewUrlPrefix("/Portal/Thing/FormReload", "/api"), "status")

	out := f.ExportForFrontend(nil, nil)

	if got := out["reloadOn"]; !reflect.DeepEqual(got, []string{"status"}) {
		t.Errorf("reloadOn=%v want [status]", got)
	}
	if out["reloadUrl"] != "/api/Portal/Thing/FormReload" {
		t.Errorf("reloadUrl=%v want /api/Portal/Thing/FormReload", out["reloadUrl"])
	}
}

func TestSetReloadOn_MultipleTriggers(t *testing.T) {
	f := NewSelectField("tags", "Tags", false, selectTestOptions())
	f.BaseField.SetReloadOn(url.NewUrl("/reload"), "status", "kind")

	out := f.ExportForFrontend(nil, nil)

	if got := out["reloadOn"]; !reflect.DeepEqual(got, []string{"status", "kind"}) {
		t.Errorf("reloadOn=%v want [status kind]", got)
	}
}

func TestReloadKeysAbsentByDefault(t *testing.T) {
	f := NewSelectField("tags", "Tags", false, selectTestOptions())

	out := f.ExportForFrontend(nil, nil)

	if _, present := out["reloadOn"]; present {
		t.Error("expected no reloadOn key on a field without dependencies")
	}
	if _, present := out["reloadUrl"]; present {
		t.Error("expected no reloadUrl key on a field without dependencies")
	}
}

// A nil URL or an empty trigger list would produce a dependency the frontend cannot resolve —
// neither key must be exported in that case.
func TestSetReloadOn_IncompleteIsIgnored(t *testing.T) {
	cases := map[string]func(f *SelectField){
		"nil url":     func(f *SelectField) { f.BaseField.SetReloadOn(nil, "status") },
		"no triggers": func(f *SelectField) { f.BaseField.SetReloadOn(url.NewUrl("/reload")) },
		"empty url":   func(f *SelectField) { f.BaseField.SetReloadOn(url.NewUrl(""), "status") },
	}

	for name, setup := range cases {
		t.Run(name, func(t *testing.T) {
			f := NewSelectField("tags", "Tags", false, selectTestOptions())
			setup(f)

			out := f.ExportForFrontend(nil, nil)

			if _, present := out["reloadOn"]; present {
				t.Errorf("expected no reloadOn key, got %v", out["reloadOn"])
			}
			if _, present := out["reloadUrl"]; present {
				t.Errorf("expected no reloadUrl key, got %v", out["reloadUrl"])
			}
		})
	}
}

func TestGetReloadOn(t *testing.T) {
	f := NewSelectField("tags", "Tags", false, selectTestOptions())
	if got := f.GetReloadOn(); len(got) != 0 {
		t.Errorf("GetReloadOn=%v want empty", got)
	}

	f.BaseField.SetReloadOn(url.NewUrl("/reload"), "status")
	if got := f.GetReloadOn(); !reflect.DeepEqual(got, []string{"status"}) {
		t.Errorf("GetReloadOn=%v want [status]", got)
	}
}
