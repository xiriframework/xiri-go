package dialog

import (
	"github.com/xiriframework/xiri-go/component/button"
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/table"
)

/*
Creates a table/list dialog with an Ok button.

Angular Frontend Flow:
  1. Dialog receives response with type="table"
  2. Angular reads res.content (tableDialogContent.Print() output)
  3. Assigns content to rawTable property
  4. Renders: <xiri-raw-table [settings]="rawTable"></xiri-raw-table>

Additional dialog options via WithOption():
  - url (string): URL for data loading or form submission
  - size (string): Dialog size ("sm", "md", "lg", "xl", "full")
  - filter (any): Filter data for the table

Example usage:

	// Define row structure
	type InfoRow struct {
	    Label string
	    Value string
	}

	// Build table
	builder := table.NewBuilder[InfoRow]()
	builder.TextField("label", "label", func(r InfoRow) string { return r.Label })
	builder.TextField("value", "value", func(r InfoRow) string { return r.Value })

	tbl := builder.Build()
	tbl.SetData([]InfoRow{
	    {Label: t("DEVICETYP"), Value: deviceTypeStr},
	    {Label: t("GERAET"), Value: deviceName},
	})

	// Create dialog
	dialog := dialog.NewDialogTable(t("INFO"), tbl)
*/

// tableDialogContent wraps a Table[T] and implements DialogContent.
// Data and fields are resolved lazily at Print time with the provided ctx.
type tableDialogContent[T any] struct {
	tbl *table.Table[T]
}

// Print resolves table data and fields lazily using the provided ctx.
func (c *tableDialogContent[T]) Print(ctx *core.UiContext) map[string]any {
	data := c.tbl.GetData(ctx, table.OutputWeb)
	fields := c.tbl.ExportFields(ctx)
	return map[string]any{
		"data":   data,
		"fields": fields,
	}
}

// NewDialogTable creates a table dialog.
// The table data and fields are resolved lazily at Print time,
// so no ctx is needed at construction time.
func NewDialogTable[T any](
	header string,
	tbl *table.Table[T],
) Dialog {

	content := &tableDialogContent[T]{tbl: tbl}

	return NewDialog(
		core.DialogTypeTable,
		header,
		content,
		[]*button.Button{
			button.NewCloseButton("Ok", core.ColorPrimary, core.ButtonTypeRaised, "", false, nil, false, nil),
		},
		nil,
		nil,
	)
}
