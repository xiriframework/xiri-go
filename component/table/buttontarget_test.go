package table

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
)

// buttonJSON digs the serialized cell button at the given key out of the exported fields.
func buttonJSON(t *testing.T, tbl *Table[struct{}], fieldID string, key int) map[string]any {
	t.Helper()

	for _, f := range tbl.exportFields(&core.UiContext{}) {
		if f["id"] != fieldID {
			continue
		}
		buttons, ok := f["buttons"].([]any)
		if !ok {
			t.Fatalf("field %q has no buttons array, got %T", fieldID, f["buttons"])
		}
		if key >= len(buttons) {
			t.Fatalf("button key %d out of range, only %d buttons", key, len(buttons))
		}
		btn, ok := buttons[key].(map[string]any)
		if !ok {
			t.Fatalf("button %d is %T, want map[string]any", key, buttons[key])
		}
		return btn
	}
	t.Fatalf("field %q not found", fieldID)
	return nil
}

// A cell button for a generated PDF must be able to say "display, do not save".
// The frontend reads that off "target" on the button JSON.
func TestWithButtonTarget_EmitsTargetOnCellButton(t *testing.T) {
	builder := NewBuilder[struct{}]()
	builder.ButtonsField("actions", "actions", func(r struct{}) map[string]string {
		return map[string]string{}
	}).
		AddButton(0, FieldButtonActionDownload, "picture_as_pdf", core.ColorPrimary, "PDF").
		WithButtonTarget(0, "_blank")

	btn := buttonJSON(t, builder.Build(), "actions", 0)

	if btn["target"] != "_blank" {
		t.Errorf("btn[\"target\"] = %v, want \"_blank\"", btn["target"])
	}
}

func TestWithButtonTarget_NoTargetByDefault(t *testing.T) {
	builder := NewBuilder[struct{}]()
	builder.ButtonsField("actions", "actions", func(r struct{}) map[string]string {
		return map[string]string{}
	}).
		AddButton(0, FieldButtonActionDownload, "download", core.ColorPrimary, "CSV")

	btn := buttonJSON(t, builder.Build(), "actions", 0)

	if v, ok := btn["target"]; ok {
		t.Errorf("expected no \"target\" key by default, got %v", v)
	}
}

// addButton rejects out-of-range keys (see buttonkey_test.go), so the setter must
// not panic on a key that was never recorded.
func TestWithButtonTarget_UnknownKeyIgnored(t *testing.T) {
	builder := NewBuilder[struct{}]()
	builder.ButtonsField("actions", "actions", func(r struct{}) map[string]string {
		return map[string]string{}
	}).
		AddButton(-5, FieldButtonActionDownload, "download", core.ColorPrimary, "x").
		WithButtonTarget(-5, "_blank"). // rejected key
		WithButtonTarget(3, "_blank")   // never added

	_ = builder.Build().exportFields(&core.UiContext{}) // must not panic
}
