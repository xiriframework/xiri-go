package dialog

import (
	"github.com/xiriframework/xiri-go/component/button"
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/table"
)

/*
Creates a table/list dialog with an Ok button.

RawTableContent Structure:
  - data (required): []map[string]any - Table row data
  - dense (optional): bool/int - Dense table spacing (0=normal, 1=dense)
  - forceMinWidth (optional): bool - Force minimum column widths

Angular Frontend Flow:
  1. Dialog receives response with type="table"
  2. Angular reads res.content (RawTableContent.Print() output)
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
	builder := table.NewBuilder[InfoRow](ctx, t)
	builder.TextField("label", "label", func(r InfoRow) string { return r.Label })
	builder.TextField("value", "value", func(r InfoRow) string { return r.Value })

	tbl := builder.Build()
	tbl.SetData([]InfoRow{
	    {Label: t("DEVICETYP"), Value: deviceTypeStr},
	    {Label: t("GERAET"), Value: deviceName},
	})

	// Convert to dialog content
	content := table.ToDialogContent(tbl, false)

	// Create dialog
	dialog := dialog.NewDialogTable(t("INFO"), content, nil, nil, nil, t)
*/

// RawTableContent represents the content structure for table dialogs
// Maps to Angular's XiriRawTableSettings interface
type RawTableContent struct {
	data          []map[string]any // Table row data (required)
	fields        []map[string]any // Table field/column definitions (required)
	dense         *int             // Dense spacing: 0=normal, 1=dense (optional)
	forceMinWidth *bool            // Force minimum column widths (optional)
}

// Print converts RawTableContent to map[string]any for JSON serialization
// This is called by Dialog.Print() when the content is a RawTableContent
func (r *RawTableContent) Print(translator core.TranslateFunc) map[string]any {

	result := map[string]any{
		"data": r.data,
	}
	if r.fields != nil && len(r.fields) > 0 {
		result["fields"] = r.fields
	}
	if r.dense != nil {
		result["dense"] = *r.dense
	}
	if r.forceMinWidth != nil {
		result["forceMinWidth"] = *r.forceMinWidth
	}
	return result
}

// NewDialogTable creates a table dialog
//
// The content parameter must be a *RawTableContent created via NewRawTableContent() or table.ToDialogContent().
// Use the builder methods to configure optional fields (dense, forceMinWidth).
func NewDialogTable[T any](
	header string,
	tbl *table.Table[T],
) Dialog {

	// Get formatted data from table (uses all formatters, locale, units, etc.)
	data := tbl.GetData(table.OutputWeb)

	content := &RawTableContent{
		data:   data,
		fields: tbl.ExportFields(),
		dense:  nil,
	}

	// Content should be a map containing XiriRawTableSettings structure
	// Angular will assign this to rawTable property and render <xiri-raw-table>
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
