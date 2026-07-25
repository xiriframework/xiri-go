package table

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
)

// #21: an out-of-range button key must not panic (negative) during serialization.
func TestButtonKeyNegativeNoPanic(t *testing.T) {
	builder := NewBuilder[struct{}]()
	builder.ButtonsField("actions", "actions", func(r struct{}) map[string]string {
		return map[string]string{}
	}).
		AddButton(-5, FieldButtonActionLink, "edit", core.ColorPrimary, "x").   // negative → panic without guard
		AddButton(1<<30, FieldButtonActionLink, "del", core.ColorWarning, "y"). // huge → ~8GB alloc without guard
		AddButton(0, FieldButtonActionLink, "ok", core.ColorPrimary, "z")

	tbl := builder.Build()
	_ = tbl.exportFields(&core.UiContext{}) // must not panic or over-allocate
}

// #21: AddMenu must not record an out-of-range key (rejected by addButton).
func TestAddMenuKeyOutOfRangeIgnored(t *testing.T) {
	builder := NewBuilder[struct{}]()
	fb := builder.ButtonsField("actions", "actions", func(r struct{}) map[string]string { return nil })

	AddMenu(fb, -1, "menu", core.ColorPrimary, "x", func(r struct{}) []string { return []string{"a"} })

	f := fb.typedField.(*field[struct{}])
	if _, ok := f.menuAccessors[-1]; ok {
		t.Error("AddMenu recorded out-of-range key in menuAccessors")
	}
	if _, ok := fb.base.menuItems[-1]; ok {
		t.Error("AddMenu recorded out-of-range key in menuItems")
	}
}
