package field

import (
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
