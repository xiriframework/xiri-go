package field

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/url"
)

// Jeder Feldtyp, der den "+"-Button bekommt, exportiert addUrl über den Basis-Export.
func TestSetAddURL_Export(t *testing.T) {
	addURL := url.NewUrlPrefix("/Portal/Tag/AddDialog", "/api")
	cases := map[string]func() map[string]interface{}{
		"select": func() map[string]interface{} {
			f := NewSelectField("t", "T", false, selectTestOptions())
			f.BaseField.SetAddURL(addURL)
			return f.ExportForFrontend(nil, nil)
		},
		"modellist": func() map[string]interface{} {
			f := NewModelListField("t", "T", false, "Tag", nil)
			f.BaseField.SetAddURL(addURL)
			return f.ExportForFrontend(nil, nil)
		},
		"model": func() map[string]interface{} {
			f := NewModelField("t", "T", false, "Tag", 0)
			f.BaseField.SetAddURL(addURL)
			return f.ExportForFrontend(nil, nil)
		},
		"chips": func() map[string]interface{} {
			f := NewChipsField("t", "T", false)
			f.BaseField.SetAddURL(addURL)
			return f.ExportForFrontend(nil, nil)
		},
	}
	for name, export := range cases {
		t.Run(name, func(t *testing.T) {
			if got := export()["addUrl"]; got != "/api/Portal/Tag/AddDialog" {
				t.Errorf("addUrl=%v want /api/Portal/Tag/AddDialog", got)
			}
		})
	}
}

func TestSetAddURL_AbsentByDefault(t *testing.T) {
	out := NewSelectField("t", "T", false, selectTestOptions()).ExportForFrontend(nil, nil)
	if _, present := out["addUrl"]; present {
		t.Error("expected no addUrl key on a field without SetAddURL")
	}
}

func TestSetAddURL_IncompleteIsIgnored(t *testing.T) {
	cases := map[string]*url.Url{"nil url": nil, "empty url": url.NewUrl("")}
	for name, u := range cases {
		t.Run(name, func(t *testing.T) {
			f := NewSelectField("t", "T", false, selectTestOptions())
			f.BaseField.SetAddURL(u)
			if _, present := f.ExportForFrontend(nil, nil)["addUrl"]; present {
				t.Error("expected no addUrl key")
			}
		})
	}
}
