package table

import (
	"github.com/xiriframework/xiri-go/component/button"
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/emptystate"
	xurl "github.com/xiriframework/xiri-go/component/url"
)

// ============================================================================
// Table Option Setters - Boolean Options
// ============================================================================

// SetReload enables/disables the reload button in the table
func (b *TableBuilder[T]) SetReload(reload bool) *TableBuilder[T] {
	b.table.options.Reload = &reload
	return b
}

// SetDense enables/disables dense spacing in the table.
//
// Legacy setter: the frontend treats dense: true as an alias for
// DensityCompact and cannot express DensityRelaxed this way. Prefer SetDensity.
func (b *TableBuilder[T]) SetDense(dense bool) *TableBuilder[T] {
	b.table.options.Dense = &dense
	return b
}

// SetDensity sets the row height of the table (compact, regular or relaxed).
//
// Takes precedence over SetDense in the frontend, so both may be set — the
// alias only applies when no explicit density arrived.
//
// Example:
//
//	builder.SetDensity(table.DensityRelaxed)
func (b *TableBuilder[T]) SetDensity(density Density) *TableBuilder[T] {
	b.table.options.Density = &density
	return b
}

// SetPagination enables/disables pagination in the table
func (b *TableBuilder[T]) SetPagination(pagination bool) *TableBuilder[T] {
	b.table.options.Pagination = &pagination
	return b
}

// SetSearch enables/disables search functionality in the table
func (b *TableBuilder[T]) SetSearch(search bool) *TableBuilder[T] {
	b.table.options.Search = &search
	return b
}

// SetQuery enables/disables query mode in the table
func (b *TableBuilder[T]) SetQuery(query bool) *TableBuilder[T] {
	b.table.options.Query = &query
	return b
}

// SetCsv enables/disables CSV export functionality
func (b *TableBuilder[T]) SetCsv(csv bool) *TableBuilder[T] {
	b.table.options.Csv = &csv
	return b
}

// SetExcel enables/disables Excel export functionality
func (b *TableBuilder[T]) SetExcel(excel bool) *TableBuilder[T] {
	b.table.options.Excel = &excel
	return b
}

// SetSaveState enables/disables state saving in local storage
func (b *TableBuilder[T]) SetSaveState(saveState bool) *TableBuilder[T] {
	b.table.options.SaveState = &saveState
	return b
}

// SetBorders enables/disables borders in the table
func (b *TableBuilder[T]) SetBorders(borders bool) *TableBuilder[T] {
	b.table.options.Borders = &borders
	return b
}

// SetBordersHeader enables/disables header borders in the table
func (b *TableBuilder[T]) SetBordersHeader(bordersHeader bool) *TableBuilder[T] {
	b.table.options.BordersHeader = &bordersHeader
	return b
}

// SetSelect enables/disables row selection in the table
func (b *TableBuilder[T]) SetSelect(selectEnabled bool) *TableBuilder[T] {
	b.table.options.Select = &selectEnabled
	return b
}

// SetFooter enables/disables footer row in the table
func (b *TableBuilder[T]) SetFooter(footer bool) *TableBuilder[T] {
	b.table.options.Footer = &footer
	return b
}

// SetScrollHeight sets a custom scroll height for the table container (e.g., "400px", "80vh").
// By default, the frontend uses 60vh. Set this to override with a specific height.
func (b *TableBuilder[T]) SetScrollHeight(scrollHeight string) *TableBuilder[T] {
	b.table.options.ScrollHeight = &scrollHeight
	return b
}

// SetEditUrl sets the URL for inline edit save requests.
// When set, fields marked as editable will allow inline editing in the Angular frontend.
// The frontend sends POST { id, field, value } to this URL when a cell is edited.
func (b *TableBuilder[T]) SetEditUrl(editUrl string) *TableBuilder[T] {
	b.table.options.EditUrl = &editUrl
	return b
}

// SetServerSide enables/disables server-side pagination.
// When enabled, pagination, sorting, and search are handled by the server.
// The frontend will send _page, _pageSize, _sort, _sortDir, _search parameters.
// The response must include totalCount for proper pagination.
func (b *TableBuilder[T]) SetServerSide(serverSide bool) *TableBuilder[T] {
	b.table.options.ServerSide = &serverSide
	return b
}

// Tree enables the opt-in tree mode for the table.
// Rows stay flat; the frontend constructs the hierarchy from the IdField/ParentIdField
// values and renders indentation + expand/collapse per node. Without this call the table
// renders as a flat list and produces identical JSON.
//
// Note: the field referenced by ParentIdField must be emitted in the row data, so it must
// NOT be hidden via .Hide() (hidden fields are dropped from row data). Add it as an id-format
// field instead — it is excluded from the visible columns but kept in the row data.
//
// Example:
//
//	builder.Tree(table.TreeConfig{
//	    IdField:         "id",
//	    ParentIdField:   "parentId",
//	    TreeColumn:      "name",
//	    PersistStateKey: "portal-groups",
//	    AddSubURL:       xurl.NewUrl("/Portal/Group/Add?parent={id}"),
//	})
func (b *TableBuilder[T]) Tree(opts TreeConfig) *TableBuilder[T] {
	b.table.options.Tree = &opts
	return b
}

// TreeAddSubWhen controls per-row visibility of the "+ sub" button: the button is only rendered
// for rows where fn returns true. Without this call the button shows on every row (when AddSubURL
// or an Angular addSubHandler is set). Call after Tree(...).
//
// Example:
//
//	builder.
//	    Tree(table.TreeConfig{ IdField: "id", ParentIdField: "parentId", AddSubURL: xurl.NewUrl("/Add?parent={id}") }).
//	    TreeAddSubWhen(func(r Region) bool { return r.AllowsChildren })
func (b *TableBuilder[T]) TreeAddSubWhen(fn func(T) bool) *TableBuilder[T] {
	b.table.treeAddSubAccessor = fn
	if b.table.options.Tree != nil {
		b.table.options.Tree.AddSubField = treeAddSubKey
	}
	return b
}

// ============================================================================
// Table Option Setters - String Options
// ============================================================================

// SetClass sets the CSS class for the table
func (b *TableBuilder[T]) SetClass(class string) *TableBuilder[T] {
	b.table.options.Class = &class
	return b
}

// SetTitle sets the table title
func (b *TableBuilder[T]) SetTitle(title string) *TableBuilder[T] {
	b.table.options.Title = &title
	return b
}

// SetTextNoData sets the text displayed when table has no data
func (b *TableBuilder[T]) SetTextNoData(textNoData string) *TableBuilder[T] {
	b.table.options.TextNoData = &textNoData
	return b
}

// SetEmptyState sets a rich empty state display for when the table has no data.
// This takes precedence over TextNoData in the frontend when both are set.
func (b *TableBuilder[T]) SetEmptyState(es *emptystate.EmptyState) *TableBuilder[T] {
	b.table.options.EmptyState = es
	return b
}

// SetMinWidth sets the minimum width CSS for the table
func (b *TableBuilder[T]) SetMinWidth(minWidth string) *TableBuilder[T] {
	b.table.options.MinWidth = &minWidth
	return b
}

// SetSaveStateId sets the state key for local storage
func (b *TableBuilder[T]) SetSaveStateId(saveStateId string) *TableBuilder[T] {
	b.table.options.SaveStateId = &saveStateId
	return b
}

// SetSaveInput sets the input saving configuration
func (b *TableBuilder[T]) SetSaveInput(saveInput string) *TableBuilder[T] {
	b.table.options.SaveInput = &saveInput
	return b
}

// SetSaveInputUrl sets the URL for saving input changes
func (b *TableBuilder[T]) SetSaveInputUrl(saveInputUrl string) *TableBuilder[T] {
	b.table.options.SaveInputUrl = &saveInputUrl
	return b
}

// SetDisplay sets the display CSS class for the table
func (b *TableBuilder[T]) SetDisplay(display string) *TableBuilder[T] {
	b.table.options.Display = &display
	return b
}

// ============================================================================
// Table Option Setters - Numeric Options
// ============================================================================

// SetItemsPerPage sets the number of items per page for pagination
func (b *TableBuilder[T]) SetItemsPerPage(itemsPerPage int) *TableBuilder[T] {
	b.table.options.ItemsPerPage = &itemsPerPage
	return b
}

// ============================================================================
// Table Option Setters - Slice Options
// ============================================================================

// SetPageSizes sets the available page size options for pagination
func (b *TableBuilder[T]) SetPageSizes(pageSizes []int) *TableBuilder[T] {
	b.table.options.PageSizes = pageSizes
	return b
}

// SetButtonsTop sets the buttons displayed at the top of the table
func (b *TableBuilder[T]) SetButtonsTop(buttons []*button.TableButton) *TableBuilder[T] {
	b.table.options.ButtonsTop = buttons
	return b
}

// ============================================================================
// Table Option Setters - Selection Options with Auto-Logic
// ============================================================================

// SetSelectButtons sets the selection action buttons and automatically enables row selection.
// When buttons slice is non-empty, Select is automatically set to true.
// When buttons slice is empty or nil, Select is automatically set to false.
//
// Example:
//
//	builder.SetSelectButtons([]*button.TableButton{
//	    button.NewTableButton(...),
//	    button.NewTableButton(...),
//	})
func (b *TableBuilder[T]) SetSelectButtons(buttons []*button.TableButton) *TableBuilder[T] {
	b.table.options.SelectButtons = buttons

	// Auto-set Select based on presence of buttons
	if len(buttons) > 0 {
		selectTrue := true
		b.table.options.Select = &selectTrue
	} else {
		selectFalse := false
		b.table.options.Select = &selectFalse
	}

	return b
}

// AddSelectButton appends a selection action button and automatically enables row selection.
// This is a convenience method for adding buttons one at a time.
//
// Example:
//
//	builder.
//	    AddSelectButton(button.NewTableButton(...)).
//	    AddSelectButton(button.NewTableButton(...))
func (b *TableBuilder[T]) AddSelectButton(button *button.TableButton) *TableBuilder[T] {
	b.table.options.SelectButtons = append(b.table.options.SelectButtons, button)

	// Auto-enable selection
	selectTrue := true
	b.table.options.Select = &selectTrue

	return b
}

// ClearSelectButtons removes all selection buttons and disables row selection
func (b *TableBuilder[T]) ClearSelectButtons() *TableBuilder[T] {
	b.table.options.SelectButtons = nil

	// Auto-disable selection
	selectFalse := false
	b.table.options.Select = &selectFalse

	return b
}

// ============================================================================
// Table Option Setters - Bulk Actions (UX-007)
// ============================================================================

// BulkActions sets the bulk-action buttons (UX-007). When set, the frontend shows a selection
// column and a sticky context bar with these actions once at least one row is selected. The
// triggered action receives the selected ids plus the selection mode (page vs. all results);
// in "all results" mode the active filter is sent too. Additive: unset = unchanged behaviour.
//
// Example:
//
//	builder.BulkActions(
//	    button.NewTableButton(core.ButtonActionApi, "archive", xurl.NewUrl("/Bulk/Archive"), "ARCHIVIEREN", core.ColorPrimary, false, nil),
//	    button.NewTableButton(core.ButtonActionApi, "delete", xurl.NewUrl("/Bulk/Delete"), "LOESCHEN", core.ColorWarning, false, nil),
//	)
func (b *TableBuilder[T]) BulkActions(buttons ...*button.TableButton) *TableBuilder[T] {
	b.table.options.BulkActions = buttons
	return b
}

// SelectAllResults enables the "select all results" affordance in the bulk bar, letting the
// user act on the whole filtered result set (all pages) instead of just the current page.
func (b *TableBuilder[T]) SelectAllResults() *TableBuilder[T] {
	enabled := true
	b.table.options.SelectAllResults = &enabled
	return b
}

// StickyBulkBar keeps the bulk context bar pinned to the top while the table is scrolled.
func (b *TableBuilder[T]) StickyBulkBar() *TableBuilder[T] {
	enabled := true
	b.table.options.StickyBulkBar = &enabled
	return b
}

// ============================================================================
// Table Option Setters - Standard Multi-Action Button Helpers
// ============================================================================

// AddMultiEditButton is a convenience method that adds a standard multi-edit button.
// Creates a button with dialog action, "edit" icon, primary color, and "BEARBEITEN" translation.
// This button is typically used to edit multiple selected rows at once.
//
// Example:
//
//	builder.AddMultiEditButton("/Portal/Device/MultiEdit")
func (b *TableBuilder[T]) AddMultiEditButton(url string) *TableBuilder[T] {
	editButton := button.NewTableButton(
		core.ButtonActionDialog,
		"edit",
		xurl.NewUrl(url),
		"BEARBEITEN", // Translation key, resolved lazily at Print time
		core.ColorPrimary,
		false,
		nil,
	)
	return b.AddSelectButton(editButton)
}

// AddMultiDeleteButton is a convenience method that adds a standard multi-delete button.
// Creates a button with dialog action, "delete" icon, warning color, and "LOESCHEN" translation.
// This button is typically used to delete multiple selected rows at once.
//
// Example:
//
//	builder.AddMultiDeleteButton("/Portal/Device/MultiDel")
func (b *TableBuilder[T]) AddMultiDeleteButton(url string) *TableBuilder[T] {
	deleteButton := button.NewTableButton(
		core.ButtonActionDialog,
		"delete",
		xurl.NewUrl(url),
		"LOESCHEN", // Translation key, resolved lazily at Print time
		core.ColorWarning,
		false,
		nil,
	)
	return b.AddSelectButton(deleteButton)
}

// AddMultiEditAndDeleteButtons is a convenience method that adds both edit and delete buttons.
// This is equivalent to calling AddMultiEditButton() and AddMultiDeleteButton() in sequence.
//
// Example:
//
//	builder.AddMultiEditAndDeleteButtons("/Portal/Device/MultiEdit", "/Portal/Device/MultiDel")
func (b *TableBuilder[T]) AddMultiEditAndDeleteButtons(editUrl, deleteUrl string) *TableBuilder[T] {
	return b.AddMultiEditButton(editUrl).AddMultiDeleteButton(deleteUrl)
}
