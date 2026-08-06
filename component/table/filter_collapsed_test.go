package table_test

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/table"
	"github.com/xiriframework/xiri-go/form/field"
	"github.com/xiriframework/xiri-go/form/group"
)

// buildFilteredTable returns a table that wraps itself in a Query because a filter is set.
func buildFilteredTable(t *testing.T, ctx *core.UiContext) *table.TableBuilder[DeviceTableRow] {
	t.Helper()

	fg, err := group.NewFormGroupWithContext([]field.FormField{
		field.NewTextField("search", "common.search", false, ""),
	}, ctx)
	if err != nil {
		t.Fatal(err)
	}

	b := table.NewBuilder[DeviceTableRow]()
	b.IdField("id", "device.id", func(r DeviceTableRow) int64 { return r.ID })
	b.TextField("name", "device.name", func(r DeviceTableRow) string { return r.Name })
	return b.SetFilter(fg)
}

// queryDataOf prints the table and returns the Query data section it is wrapped in.
func queryDataOf(t *testing.T, tbl *table.Table[DeviceTableRow], ctx *core.UiContext) map[string]any {
	t.Helper()

	out := tbl.Print(ctx)
	if out["type"] != "query" {
		t.Fatalf("type=%v, want query (table with filter must be wrapped)", out["type"])
	}
	data, ok := out["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data map, got %T", out["data"])
	}
	return data
}

func TestTableFilterCollapsed(t *testing.T) {
	ctx := exampleContext()
	tbl := buildFilteredTable(t, ctx).SetFilterCollapsed(true).Build()

	if got := queryDataOf(t, tbl, ctx)["collapsed"]; got != true {
		t.Errorf("collapsed=%v, want true", got)
	}
}

// Explicit false must still serialize — it renders an expanded panel, not "no panel".
func TestTableFilterCollapsedFalseIsSerialized(t *testing.T) {
	ctx := exampleContext()
	tbl := buildFilteredTable(t, ctx).SetFilterCollapsed(false).Build()

	data := queryDataOf(t, tbl, ctx)
	val, ok := data["collapsed"]
	if !ok {
		t.Fatalf("expected collapsed key when explicitly set")
	}
	if val != false {
		t.Errorf("collapsed=%v, want false", val)
	}
}

// Without the option the key stays absent, so the frontend renders no panel at all.
func TestTableFilterCollapsedOmittedByDefault(t *testing.T) {
	ctx := exampleContext()
	tbl := buildFilteredTable(t, ctx).Build()

	if _, ok := queryDataOf(t, tbl, ctx)["collapsed"]; ok {
		t.Errorf("did not expect collapsed key by default")
	}
}
