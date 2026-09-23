package field

import (
	"slices"
	"strings"
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
)

func loaderOf(ids ...int32) ModelLoaderFunc {
	return func(*core.UiContext, string) ([]ModelOption, error) {
		opts := make([]ModelOption, len(ids))
		for i, id := range ids {
			opts[i] = ModelOption{ID: id, Name: "x"}
		}
		return opts, nil
	}
}

func loadedModelField(def int32, ids ...int32) *ModelField {
	f := NewModelField("group_id", "GROUP", false, "group", def)
	f.SetLoaderFunc(loaderOf(ids...))
	if err := f.LoadOptions(&core.UiContext{}); err != nil {
		panic(err)
	}
	return f
}

func TestModelField_Validate_OnlyOfferedIDs(t *testing.T) {
	cases := []struct {
		name  string
		field func() *ModelField
		value interface{}
		ok    bool
	}{
		{"offered", func() *ModelField { return loadedModelField(0, 1, 2) }, int32(1), true},
		{"not offered", func() *ModelField { return loadedModelField(0, 1, 2) }, int32(999), false},
		{"zero means none", func() *ModelField { return loadedModelField(0, 1, 2) }, int32(0), true},
		{"unchanged default not offered", func() *ModelField { return loadedModelField(42, 1, 2) }, int32(42), true},
		{"other value not offered despite default", func() *ModelField { return loadedModelField(42, 1, 2) }, int32(43), false},
		{"sub removes offered", func() *ModelField {
			f := loadedModelField(0, 1, 2)
			f.Sub = []int32{2}
			return f
		}, int32(2), false},
		{"loader returned nothing", func() *ModelField { return loadedModelField(0) }, int32(1), false},
		{"loader never ran (no context)", func() *ModelField {
			f := NewModelField("g", "G", false, "group", 0)
			f.SetLoaderFunc(loaderOf(1))
			return f
		}, int32(1), false},
		{"static List", func() *ModelField {
			f := NewModelField("g", "G", false, "group", 0)
			f.List = []ModelOption{{ID: 5}}
			return f
		}, int32(6), false},
		{"loader + url: search hit outside list", func() *ModelField {
			f := loadedModelField(0, 1, 2)
			f.URL = "/api/groups/search"
			return f
		}, int32(999), true},
		{"url only: server cannot know", func() *ModelField {
			f := NewModelField("g", "G", false, "group", 0)
			f.URL = "/api/groups"
			return f
		}, int32(999), true},
		{"url: sub still enforced", func() *ModelField {
			f := loadedModelField(0, 1, 3)
			f.URL = "/api/groups"
			f.Sub = []int32{3}
			return f
		}, int32(3), false},
		{"no loader, no list, no url (reload-dependent)", func() *ModelField {
			return NewModelField("g", "G", false, "group", 0)
		}, int32(7), true},
		{"int64 value", func() *ModelField { return loadedModelField(0, 1) }, int64(999), false},
		{"unchanged default set directly as int", func() *ModelField {
			f := loadedModelField(0, 1)
			f.Default = 42
			return f
		}, 42, true},
		{"unchanged default set directly as int64", func() *ModelField {
			f := loadedModelField(0, 1)
			f.Default = int64(42)
			return f
		}, int32(42), true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.field().Validate(c.value)
			if (err == nil) != c.ok {
				t.Errorf("Validate(%v): err = %v, want ok=%v", c.value, err, c.ok)
			}
		})
	}
}

func TestModelListField_Validate_OnlyOfferedIDs(t *testing.T) {
	// 9 ist ein Alteintrag, den der User nicht (mehr) sieht.
	f := NewModelListField("devices", "DEVICES", false, "device", []int32{9})
	f.SetLoaderFunc(loaderOf(1, 2))
	if err := f.LoadOptions(&core.UiContext{}); err != nil {
		t.Fatal(err)
	}
	if err := f.Validate(ModelListValue{1, 2}); err != nil {
		t.Errorf("offered ids rejected: %v", err)
	}
	if err := f.Validate(ModelListValue{1, 9}); err != nil {
		t.Errorf("unchanged default id 9 rejected: %v", err)
	}
	if err := f.Validate(ModelListValue{1, 4242}); err == nil {
		t.Error("expected error for id 4242 not among options")
	}
}

func TestModelListField_Validate_Cases(t *testing.T) {
	loaded := func(ids ...int32) *ModelListField {
		f := NewModelListField("devices", "DEVICES", false, "device", nil)
		f.SetLoaderFunc(loaderOf(ids...))
		if err := f.LoadOptions(&core.UiContext{}); err != nil {
			panic(err)
		}
		return f
	}
	cases := []struct {
		name  string
		field func() *ModelListField
		value ModelListValue
		ok    bool
	}{
		{"url: search hit outside list", func() *ModelListField {
			f := loaded(1)
			f.URL = "/api/devices"
			return f
		}, ModelListValue{1, 999}, true},
		{"url: sub still enforced", func() *ModelListField {
			f := loaded(1, 3)
			f.URL = "/api/devices"
			f.Sub = []int32{3}
			return f
		}, ModelListValue{1, 3}, false},
		{"loader returned nothing", func() *ModelListField { return loaded() }, ModelListValue{1}, false},
		{"loader never ran (no context)", func() *ModelListField {
			f := NewModelListField("d", "D", false, "device", nil)
			f.SetLoaderFunc(loaderOf(1))
			return f
		}, ModelListValue{1}, false},
		{"static List", func() *ModelListField {
			f := NewModelListField("d", "D", false, "device", nil)
			f.List = []ModelOption{{ID: 5}}
			return f
		}, ModelListValue{5, 6}, false},
		{"no loader, no list, no url (reload-dependent)", func() *ModelListField {
			return NewModelListField("d", "D", false, "device", nil)
		}, ModelListValue{7}, true},
		{"zero is a regular id", func() *ModelListField { return loaded(1) }, ModelListValue{0}, false},
		{"empty list", func() *ModelListField { return loaded(1) }, ModelListValue{}, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.field().Validate(c.value)
			if (err == nil) != c.ok {
				t.Errorf("Validate(%v): err = %v, want ok=%v", c.value, err, c.ok)
			}
		})
	}
}

// Ein direkt gesetzter Default (int statt int32) muss als int32 gebunden werden, sonst bleibt
// Value still auf 0 und überschreibt beim Speichern den aktuellen Wert.
func TestModelField_BindValue_DirectIntDefault(t *testing.T) {
	f := loadedModelField(0, 1)
	f.Default = 42
	if err := f.BindValue(nil); err != nil {
		t.Fatalf("BindValue(nil): %v", err)
	}
	if f.Value != 42 {
		t.Errorf("Value = %d, want 42", f.Value)
	}
}

// Ein direkt als []int32 gesetzter Default darf weder paniken noch die Default-Ausnahme verlieren.
func TestModelListField_DirectSliceDefault(t *testing.T) {
	f := NewModelListField("devices", "DEVICES", false, "device", nil)
	f.Default = []int32{9}
	f.SetLoaderFunc(loaderOf(1))
	if err := f.LoadOptions(&core.UiContext{}); err != nil {
		t.Fatal(err)
	}
	if err := f.BindValue(nil); err != nil {
		t.Fatalf("BindValue(nil): %v", err)
	}
	if len(f.Value) != 1 || f.Value[0] != 9 {
		t.Errorf("Value = %v, want [9]", f.Value)
	}
	if err := f.Validate(ModelListValue{1, 9}); err != nil {
		t.Errorf("unchanged default id 9 rejected: %v", err)
	}
}

func TestModelListField_UnsupportedDefaultFails(t *testing.T) {
	f := NewModelListField("devices", "DEVICES", false, "device", nil)
	f.Default = []int{9}
	if err := f.BindValue(nil); err == nil {
		t.Fatalf("expected error for []int default, got Value %v", f.Value)
	}
}

// AllowedFunc ersetzt die Listenprüfung, sieht aber weder 0, den unveränderten Default noch Sub.
func TestModelField_Validate_AllowedFunc(t *testing.T) {
	var calls [][]int32
	field := func(def int32) *ModelField {
		f := loadedModelField(def, 1, 2)
		f.URL = "/api/groups/search"
		f.Sub = []int32{5}
		f.SetAllowedFunc(func(ids []int32) bool {
			calls = append(calls, ids)
			return ids[0] == 7
		})
		return f
	}
	cases := []struct {
		name  string
		def   int32
		value interface{}
		ok    bool
		calls int
	}{
		{"search hit outside list", 0, int32(7), true, 1},
		{"hook replaces list", 0, int32(1), false, 1},
		{"zero means none", 0, int32(0), true, 0},
		{"unchanged default", 42, int32(42), true, 0},
		{"sub", 0, int32(5), false, 0},
		{"int64 value", 0, int64(7), true, 1},
		{"int64 overflow", 0, int64(1 << 40), false, 0},
		{"string rejected", 0, "7", false, 0},
		{"float rejected", 0, float64(7), false, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			calls = nil
			err := field(c.def).Validate(c.value)
			if (err == nil) != c.ok {
				t.Errorf("Validate(%v): err = %v, want ok=%v", c.value, err, c.ok)
			}
			if len(calls) != c.calls {
				t.Errorf("hook calls = %v, want %d", calls, c.calls)
			}
		})
	}
}

// Sub ist eine harte Sperre, auch für den unveränderten Default.
func TestModelField_Validate_SubBeatsDefault(t *testing.T) {
	f := loadedModelField(3, 1, 3)
	f.Sub = []int32{3}
	if err := f.Validate(int32(3)); err == nil {
		t.Error("default in Sub accepted")
	}
}

func TestModelListField_Validate_AllowedFunc(t *testing.T) {
	var calls [][]int32
	f := NewModelListField("devices", "DEVICES", false, "device", []int32{3})
	f.URL = "/api/devices"
	f.Tree = true
	f.SetAllowedFunc(func(ids []int32) bool {
		calls = append(calls, ids)
		for _, id := range ids {
			if id != 7 && id != 8 {
				return false
			}
		}
		return true
	})

	calls = nil
	if err := f.Validate(ModelListValue{3, 7, 8}); err != nil {
		t.Errorf("[3,7,8] rejected: %v", err)
	}
	if len(calls) != 1 || !slices.Equal(calls[0], []int32{7, 8}) {
		t.Errorf("calls = %v, want one call with [7 8]", calls)
	}

	err := f.Validate(ModelListValue{3, 7, 9})
	if err == nil {
		t.Fatal("[3,7,9] accepted")
	}
	if strings.Contains(err.Error(), "7") || strings.Contains(err.Error(), "9") {
		t.Errorf("error names an id, but the hook does not say which: %v", err)
	}

	calls = nil
	if err := f.Validate(ModelListValue{3}); err != nil || len(calls) != 0 {
		t.Errorf("[3]: err = %v, calls = %v, want ok without call", err, calls)
	}

	calls = nil
	if err := f.Validate(ModelListValue{0}); err == nil || len(calls) != 1 {
		t.Errorf("[0]: err = %v, calls = %v, want rejection by hook", err, calls)
	}

	f.Sub = []int32{3}
	calls = nil
	if err := f.Validate(ModelListValue{3}); err == nil || len(calls) != 0 {
		t.Errorf("default in Sub: err = %v, calls = %v, want rejection without call", err, calls)
	}
}
