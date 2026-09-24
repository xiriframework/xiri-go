package builder

import (
	"errors"
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/form/field"
	"github.com/xiriframework/xiri-go/form/group"
	"github.com/xiriframework/xiri-go/response"
)

func TestBindFromMap_Basic(t *testing.T) {
	name := field.NewTextField("name", "NAME", true, "")
	active := field.NewBoolField("active", "ACTIVE", false, false)

	fg := group.NewFormGroup([]field.FormField{name, active})

	data := map[string]interface{}{
		"name":   "test-name",
		"active": true,
	}

	if err := BindFromMap(data, fg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if name.Value == nil || *name.Value != "test-name" {
		t.Errorf("expected 'test-name', got %v", name.Value)
	}
	if active.Value == nil || *active.Value != true {
		t.Errorf("expected true, got %v", active.Value)
	}
}

func TestBindFromMap_MissingOptionalField(t *testing.T) {
	name := field.NewTextField("name", "NAME", true, "")
	count := field.NewIntField("count", "COUNT", false, 99)

	fg := group.NewFormGroup([]field.FormField{name, count})

	data := map[string]interface{}{
		"name": "hello",
		// count is missing - should use default
	}

	if err := BindFromMap(data, fg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if name.Value == nil || *name.Value != "hello" {
		t.Errorf("expected 'hello', got %v", name.Value)
	}
	if count.Value == nil || *count.Value != 99 {
		t.Errorf("expected default 99, got %v", count.Value)
	}
}

func TestBindFromMap_ValidationError(t *testing.T) {
	name := field.NewTextFieldWithLength("name", "NAME", true, "", 5, 100)

	fg := group.NewFormGroup([]field.FormField{name})

	data := map[string]interface{}{
		"name": "ab", // too short
	}

	err := BindFromMap(data, fg)
	if err == nil {
		t.Fatal("expected validation error for too-short text")
	}
}

func TestNewFormBuilder_BuildAdd(t *testing.T) {
	name := field.NewTextField("name", "NAME", true, "default-name")
	active := field.NewBoolField("active", "ACTIVE", false, true)

	builder := NewFormBuilder(nil)
	builder.AddField(name).AddField(active)

	fg, defaults, err := builder.BuildAdd()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(fg.GetFields()) != 2 {
		t.Errorf("expected 2 fields, got %d", len(fg.GetFields()))
	}

	if defaults["name"] != "default-name" {
		t.Errorf("expected default 'default-name', got %v", defaults["name"])
	}
	if defaults["active"] != true {
		t.Errorf("expected default true, got %v", defaults["active"])
	}
}

func TestNewFormBuilder_BuildEdit(t *testing.T) {
	name := field.NewTextField("name", "NAME", true, "current-name")

	builder := NewFormBuilder(nil)
	builder.AddField(name)

	fg, values, err := builder.BuildEdit()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(fg.GetFields()) != 1 {
		t.Errorf("expected 1 field, got %d", len(fg.GetFields()))
	}

	if values["name"] != "current-name" {
		t.Errorf("expected 'current-name', got %v", values["name"])
	}
}

func TestNewFormBuilder_BuildAddForDisplay(t *testing.T) {
	name := field.NewTextField("name", "NAME", true, "test")

	builder := NewFormBuilder(nil)
	builder.AddField(name)

	exported, err := builder.BuildAddForDisplay()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(exported) != 1 {
		t.Fatalf("expected 1 exported field, got %d", len(exported))
	}
	if exported[0]["id"] != "name" {
		t.Errorf("expected field id 'name', got %v", exported[0]["id"])
	}
}

// TestBindFromMap_DisabledFieldUsesDefault verifies that disabled fields use
// their default value and ignore submitted data (M5).
func TestBindFromMap_DisabledFieldUsesDefault(t *testing.T) {
	name := field.NewTextField("name", "NAME", true, "secure-default")
	name.SetDisabled(true)
	active := field.NewBoolField("active", "ACTIVE", false, false)

	fg := group.NewFormGroup([]field.FormField{name, active})

	data := map[string]interface{}{
		"name":   "attacker-value",
		"active": true,
	}

	if err := BindFromMap(data, fg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Disabled field must use default, not submitted value
	if name.Value == nil || *name.Value != "secure-default" {
		t.Errorf("expected 'secure-default', got %v", name.Value)
	}
	// Non-disabled field must accept submitted value
	if active.Value == nil || *active.Value != true {
		t.Errorf("expected true, got %v", active.Value)
	}
}

// TestBindFromMap_DisabledIntField verifies that disabled IntField uses
// its default value and ignores submitted data (M5).
func TestBindFromMap_DisabledIntField(t *testing.T) {
	count := field.NewIntField("count", "COUNT", false, 42)
	count.SetDisabled(true)

	fg := group.NewFormGroup([]field.FormField{count})

	data := map[string]interface{}{
		"count": float64(999),
	}

	if err := BindFromMap(data, fg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if count.Value == nil || *count.Value != 42 {
		t.Errorf("expected 42, got %v", count.Value)
	}
}

// TestBindFromMap_DisabledModelField verifies that disabled ModelField uses
// its default value and ignores submitted data (M5).
func TestBindFromMap_DisabledModelField(t *testing.T) {
	device := field.NewModelField("device", "DEVICE", false, "Device", int32(10))
	device.SetDisabled(true)

	fg := group.NewFormGroup([]field.FormField{device})

	data := map[string]interface{}{
		"device": float64(999),
	}

	if err := BindFromMap(data, fg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if device.Value != 10 {
		t.Errorf("expected 10, got %v", device.Value)
	}
}

func TestNewFormBuilder_OnEditValueCheck(t *testing.T) {
	name := field.NewTextField("name", "NAME", true, "default")

	builder := NewFormBuilder(nil)
	builder.AddField(name)

	hookCalled := false
	builder.OnEditValueCheck = func(fg *group.FormGroup, values map[string]interface{}) error {
		hookCalled = true
		values["name"] = "modified-by-hook"
		return nil
	}

	_, values, err := builder.BuildEdit()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !hookCalled {
		t.Error("expected OnEditValueCheck hook to be called")
	}
	if values["name"] != "modified-by-hook" {
		t.Errorf("expected 'modified-by-hook', got %v", values["name"])
	}
}

func TestBindFromMap_CollectsAllFieldErrors(t *testing.T) {
	name := field.NewTextFieldWithLength("name", "NAME", true, "", 3, 10)
	email := field.NewTextFieldWithLength("email", "EMAIL", true, "", 5, 50)
	fg := group.NewFormGroup([]field.FormField{name, email})

	err := BindFromMap(map[string]interface{}{"name": "ab", "email": "x"}, fg)

	var fe group.FieldErrors
	if !errors.As(err, &fe) {
		t.Fatalf("expected group.FieldErrors, got %T: %v", err, err)
	}
	if _, ok := fe["name"]; !ok {
		t.Errorf("missing error for name: %v", fe)
	}
	if _, ok := fe["email"]; !ok {
		t.Errorf("missing error for email: %v", fe)
	}
}

func TestBindFromMap_DisabledTextFieldDefaultTrimmed(t *testing.T) {
	name := field.NewTextField("name", "NAME", false, " x ").SetDisabled(true)
	fg := group.NewFormGroup([]field.FormField{name})

	if err := BindFromMap(map[string]interface{}{}, fg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name.Value == nil || *name.Value != "x" {
		t.Errorf("expected \"x\", got %v", name.Value)
	}
}

func TestBindFromMap_PintRejectsNegative(t *testing.T) {
	n := field.NewIntField("n", "N", false, 0)
	n.Subtype = "pint"
	fg := group.NewFormGroup([]field.FormField{n})

	err := BindFromMap(map[string]interface{}{"n": float64(-1)}, fg)

	var fe group.FieldErrors
	if !errors.As(err, &fe) {
		t.Fatalf("expected group.FieldErrors, got %T: %v", err, err)
	}
	if _, ok := fe["n"]; !ok {
		t.Errorf("missing error for n: %v", fe)
	}
}

func TestBindFromMap_ModelFieldRejectsForeignID(t *testing.T) {
	g := field.NewModelField("group_id", "GROUP", true, "group", 0)
	g.SetLoaderFunc(func(*core.UiContext, string) ([]field.ModelOption, error) {
		return []field.ModelOption{{ID: 1, Name: "Mine"}}, nil
	})
	fg, err := group.NewFormGroupWithContext([]field.FormField{g}, &core.UiContext{})
	if err != nil {
		t.Fatal(err)
	}

	err = BindFromMap(map[string]interface{}{"group_id": float64(999)}, fg)

	var fe group.FieldErrors
	if !errors.As(err, &fe) || fe["group_id"] == "" {
		t.Fatalf("expected FieldErrors for group_id, got %T: %v", err, err)
	}
}

// BindReloadFromMap meldet keine Feldfehler: eine fremde ID wird verworfen, der Default bleibt.
func TestBindReloadFromMap_ModelFieldForeignIDKeepsDefault(t *testing.T) {
	g := field.NewModelField("group_id", "GROUP", false, "group", 1)
	g.SetLoaderFunc(func(*core.UiContext, string) ([]field.ModelOption, error) {
		return []field.ModelOption{{ID: 1, Name: "Mine"}, {ID: 2, Name: "Also mine"}}, nil
	})
	fg, err := group.NewFormGroupWithContext([]field.FormField{g}, &core.UiContext{})
	if err != nil {
		t.Fatal(err)
	}

	if err := BindReloadFromMap(map[string]interface{}{"group_id": float64(999)}, fg); err != nil {
		t.Fatalf("BindReloadFromMap: %v", err)
	}
	if g.Value != 1 {
		t.Errorf("Value = %d, want default 1", g.Value)
	}
}

// Dokumentiertes Muster: Loader + Server-Suche, abgesichert per AllowedFunc.
func TestBindFromMap_ModelFieldAllowedFunc(t *testing.T) {
	g := field.NewModelField("group_id", "GROUP", true, "group", 0)
	g.SetLoaderFunc(func(*core.UiContext, string) ([]field.ModelOption, error) {
		return []field.ModelOption{{ID: 1, Name: "Mine"}}, nil
	})
	g.URL = "/api/groups/search"
	g.SetAllowedFunc(func(ids []int32) bool { return ids[0] == 1 || ids[0] == 2 })
	fg, err := group.NewFormGroupWithContext([]field.FormField{g}, &core.UiContext{})
	if err != nil {
		t.Fatal(err)
	}

	if err := BindFromMap(map[string]interface{}{"group_id": float64(2)}, fg); err != nil {
		t.Fatalf("search hit 2 rejected: %v", err)
	}
	err = BindFromMap(map[string]interface{}{"group_id": float64(999)}, fg)
	var fe group.FieldErrors
	if !errors.As(err, &fe) || fe["group_id"] == "" {
		t.Fatalf("expected FieldErrors for group_id, got %T: %v", err, err)
	}
}

func TestBindReloadFromMap_ModelFieldAllowedFunc(t *testing.T) {
	g := field.NewModelField("group_id", "GROUP", false, "group", 1)
	g.URL = "/api/groups/search"
	g.SetAllowedFunc(func(ids []int32) bool { return false })
	fg := group.NewFormGroup([]field.FormField{g})

	if err := BindReloadFromMap(map[string]interface{}{"group_id": float64(999)}, fg); err != nil {
		t.Fatalf("BindReloadFromMap: %v", err)
	}
	if g.Value != 1 {
		t.Errorf("Value = %d, want default 1", g.Value)
	}
}

// Typisierte Defaults (TimeRangeValue, TimeLimitValue, *GeoformValue) sind kein Request-Wert:
// ein fehlendes oder gesperrtes Feld bindet sie über Parse(nil), nicht über Parse(default).
func typedDefaultFields() (*field.TimeRangeField, *field.TimeLimitField, *field.GeoformField) {
	geo := field.NewGeoformField("geo", "GEO", false)
	geo.Default = &field.GeoformValue{Type: 2, Path: map[string]string{"lat": "1", "lng": "2", "radius": "3"}}
	return field.NewTimeRangeFieldWithDefault("tr", "TR", false, 7), field.NewTimeLimitField("tl", "TL", false), geo
}

func TestBindFromMap_MissingTypedDefaults(t *testing.T) {
	tr, tl, geo := typedDefaultFields()
	fg := group.NewFormGroup([]field.FormField{tr, tl, geo})

	if err := BindFromMap(map[string]interface{}{}, fg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.Value != tr.Default {
		t.Errorf("tr.Value = %v, want default %v", tr.Value, tr.Default)
	}
}

func TestBindFromMap_DisabledTimeRangeKeepsDefault(t *testing.T) {
	tr := field.NewTimeRangeFieldWithDefault("tr", "TR", false, 7)
	tr.SetDisabled(true)
	fg := group.NewFormGroup([]field.FormField{tr})

	err := BindFromMap(map[string]interface{}{"tr": map[string]interface{}{"start": "2026-01-01", "end": "2026-01-02"}}, fg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.Value != tr.Default {
		t.Errorf("tr.Value = %v, want default %v", tr.Value, tr.Default)
	}
}

func TestBindReloadFromMap_MissingTypedDefaults(t *testing.T) {
	tr, tl, geo := typedDefaultFields()
	fg := group.NewFormGroup([]field.FormField{tr, tl, geo})

	if err := BindReloadFromMap(map[string]interface{}{}, fg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.Value != tr.Default {
		t.Errorf("tr.Value = %v, want default %v", tr.Value, tr.Default)
	}
}

func TestBindFromMap_MissingDefaultsStillBind(t *testing.T) {
	n := field.NewIntField("n", "N", false, 7)
	m := field.NewModelField("m", "M", false, "device", 0)
	m.Default = 42
	txt := field.NewTextField("t", "T", false, " x ")
	chips := field.NewChipsField("c", "C", false)
	chips.Default = []string{"a"}
	req := field.NewTextField("r", "R", true, "")
	req.Default = nil
	fg := group.NewFormGroup([]field.FormField{n, m, txt, chips, req})

	err := BindFromMap(map[string]interface{}{}, fg)

	var fe group.FieldErrors
	if !errors.As(err, &fe) || len(fe) != 1 || fe["r"] == "" {
		t.Fatalf("expected only a required error for r, got %T: %v", err, err)
	}
	if n.Value == nil || *n.Value != 7 || m.Value != 42 || txt.Value == nil || *txt.Value != "x" || len(chips.Value) != 1 {
		t.Errorf("defaults not bound: n=%v m=%v t=%v c=%v", n.Value, m.Value, txt.Value, chips.Value)
	}
}

func TestBindFromMap_TranslatesWrappedValidationError(t *testing.T) {
	ort := field.NewTextFieldWithLength("ort", "ORT", false, "", 0, 5)
	texts := map[string]string{"ORT": "Ort", "validation.max_length": "Höchstens {max} Zeichen"}
	ctx := &core.UiContext{Translate: func(k string) string {
		if t, ok := texts[k]; ok {
			return t
		}
		return k
	}}
	fg, err := group.NewFormGroupWithContext([]field.FormField{ort}, ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = BindFromMap(map[string]interface{}{"ort": "Wien-Floridsdorf"}, fg)

	var fe group.FieldErrors
	if !errors.As(err, &fe) || fe["ort"] != "Höchstens 5 Zeichen" {
		t.Fatalf("expected translated field error, got %T: %v", err, err)
	}
	if err.Error() != "Ort: Höchstens 5 Zeichen" {
		t.Fatalf("banner: %q", err.Error())
	}
	if resp := response.NewErrorResponseFromError(err); resp.Fields["ort"] != "Höchstens 5 Zeichen" {
		t.Fatalf("response fields: %v", resp.Fields)
	}
}
