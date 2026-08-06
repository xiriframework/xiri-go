package group

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/url"
	"github.com/xiriframework/xiri-go/form/field"
)

func reloadTestOptions() []field.SelectOption {
	return []field.SelectOption{
		{Value: int32(1), Label: "one"},
		{Value: int32(2), Label: "two"},
	}
}

func dependentSelect(id string, triggers ...string) *field.SelectField {
	f := field.NewSelectField(id, id, false, reloadTestOptions())
	f.BaseField.SetReloadOn(url.NewUrl("/reload"), triggers...)
	return f
}

func TestExportPatch_OnlyDependentFields(t *testing.T) {
	status := field.NewSelectField("status", "Status", false, reloadTestOptions())
	tags := dependentSelect("tags", "status")
	note := field.NewTextField("note", "Note", false, "")

	fg := NewFormGroup([]field.FormField{status, tags, note})

	patch := fg.ExportPatch()

	if len(patch) != 1 {
		t.Fatalf("patch has %d entries, want 1: %v", len(patch), patch)
	}
	entry, ok := patch["tags"].(map[string]interface{})
	if !ok {
		t.Fatalf("patch[tags]=%T want map", patch["tags"])
	}
	if entry["id"] != "tags" {
		t.Errorf("entry id=%v want tags", entry["id"])
	}
	if entry["list"] == nil {
		t.Error("expected the dependent field's option list in the patch")
	}
}

func TestExportPatch_EmptyWithoutDependencies(t *testing.T) {
	fg := NewFormGroup([]field.FormField{
		field.NewTextField("note", "Note", false, ""),
	})

	if patch := fg.ExportPatch(); len(patch) != 0 {
		t.Errorf("patch=%v want empty", patch)
	}
}

// Form=false fields are not rendered, so patching them would target a control that does
// not exist - consistent with ExportForFrontend skipping them.
func TestExportPatch_SkipsNonFormFields(t *testing.T) {
	hidden := dependentSelect("hidden", "status")
	hidden.BaseField.SetForm(false)

	fg := NewFormGroup([]field.FormField{
		field.NewSelectField("status", "Status", false, reloadTestOptions()),
		hidden,
	})

	if patch := fg.ExportPatch(); len(patch) != 0 {
		t.Errorf("patch=%v want empty", patch)
	}
}

func TestExportPatch_MultipleDependentFields(t *testing.T) {
	fg := NewFormGroup([]field.FormField{
		field.NewSelectField("status", "Status", false, reloadTestOptions()),
		dependentSelect("tags", "status"),
		dependentSelect("groups", "status"),
	})

	patch := fg.ExportPatch()

	if len(patch) != 2 {
		t.Fatalf("patch has %d entries, want 2: %v", len(patch), patch)
	}
	for _, id := range []string{"tags", "groups"} {
		if _, ok := patch[id]; !ok {
			t.Errorf("missing %q in patch", id)
		}
	}
}

// Ein Feld, dem die URL fehlt, exportiert im normalen Rendern keine Reload-Keys - dann darf es
// auch nicht im Patch landen. Erreichbar nur, wenn jemand ReloadOn direkt setzt statt SetReloadOn.
func TestExportPatch_SkipsFieldWithoutReloadURL(t *testing.T) {
	half := field.NewSelectField("tags", "Tags", false, reloadTestOptions())
	half.BaseField.ReloadOn = []string{"status"}

	fg := NewFormGroup([]field.FormField{
		field.NewSelectField("status", "Status", false, reloadTestOptions()),
		half,
	})

	if patch := fg.ExportPatch(); len(patch) != 0 {
		t.Errorf("patch=%v want empty", patch)
	}
}
