package table

import "testing"

// SetFlat marks a table as embedded (no elevation/background of its own), e.g. inside an
// expansion panel. Unset must stay absent from the JSON so existing output is unchanged.
func TestSetFlat(t *testing.T) {
	ctx := testContext()

	build := func(set func(*TableBuilder[testDeviceRow])) map[string]any {
		b := NewBuilder[testDeviceRow]()
		b.IdField("id", "device.id", func(r testDeviceRow) int64 { return r.ID })
		if set != nil {
			set(b)
		}
		tbl := b.Build()
		tbl.SetData([]testDeviceRow{{ID: 1, Name: "Device 1"}})
		return tbl.Print(ctx)["data"].(map[string]any)["options"].(map[string]any)
	}

	if got := build(func(b *TableBuilder[testDeviceRow]) { b.SetFlat(true) })["flat"]; got != true {
		t.Errorf("flat(true)=%v want true", got)
	}
	if got := build(func(b *TableBuilder[testDeviceRow]) { b.SetFlat(false) })["flat"]; got != false {
		t.Errorf("flat(false)=%v want false", got)
	}
	if _, has := build(nil)["flat"]; has {
		t.Error("expected no 'flat' key when unset")
	}
}
