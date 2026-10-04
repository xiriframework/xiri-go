package table_test

import (
	"reflect"
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/table"
)

type menuRow struct {
	ID     int
	Booked bool
	Hidden bool
}

func menuCell(t *testing.T, build func(fb *table.FieldBuilder), rows ...menuRow) []any {
	t.Helper()
	builder := table.NewBuilder[menuRow]()
	fb := builder.ButtonsField("actions", "Actions", func(r menuRow) map[string]string { return nil })
	build(fb)
	fb.AddMenuItem(table.FieldButtonActionLink, "edit", core.ColorPrimary, "Edit").
		AddMenuItem(table.FieldButtonActionDialog, "cancel", core.ColorWarning, "Cancel").
		AddMenuItem(table.FieldButtonActionLink, "delete", core.ColorWarning, "Delete")

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(exampleContext(), table.OutputWeb)
	cells := make([]any, len(data))
	for i, row := range data {
		cells[i] = row["actions"].(map[string]any)["0"]
	}
	return cells
}

func TestAddMenuEntriesSerializesDisabledItems(t *testing.T) {
	cells := menuCell(t, func(fb *table.FieldBuilder) {
		table.AddMenuEntries(fb, 0, "more_vert", core.ColorPrimary, "Actions", func(r menuRow) []table.MenuEntry {
			if r.Hidden {
				return nil
			}
			cancel := table.MenuEntry{URL: "/cancel"}
			if r.Booked {
				cancel = table.MenuEntry{URL: "/cancel", DisabledHint: "Schon verbucht"}
			}
			return []table.MenuEntry{{URL: "/edit"}, cancel, {}}
		})
	}, menuRow{ID: 1, Booked: true}, menuRow{ID: 2}, menuRow{ID: 3, Hidden: true})

	want := []any{
		[]any{"/edit", map[string]any{"disabled": "Schon verbucht"}, false},
		[]any{"/edit", "/cancel", false},
		false,
	}
	if !reflect.DeepEqual(cells, want) {
		t.Errorf("menu cells = %#v, want %#v", cells, want)
	}
}

func TestAddMenuStringsUnchanged(t *testing.T) {
	cells := menuCell(t, func(fb *table.FieldBuilder) {
		table.AddMenu(fb, 0, "more_vert", core.ColorPrimary, "Actions", func(r menuRow) []string {
			if r.Hidden {
				return nil
			}
			return []string{"/edit", "", "/delete"}
		})
	}, menuRow{ID: 1}, menuRow{ID: 2, Hidden: true})

	want := []any{[]any{"/edit", false, "/delete"}, false}
	if !reflect.DeepEqual(cells, want) {
		t.Errorf("menu cells = %#v, want %#v", cells, want)
	}
}
