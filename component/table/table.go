package table

import (
	"log/slog"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/xiriframework/xiri-go/component/button"
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/emptystate"
	"github.com/xiriframework/xiri-go/component/query"
	xurl "github.com/xiriframework/xiri-go/component/url"
	"github.com/xiriframework/xiri-go/form/group"
)

// tableCore contains all type-independent table state and configuration.
// Methods on tableCore are non-generic and exist only once in the binary,
// avoiding monomorphization across all Table[T] instantiations.
type tableCore struct {
	fieldBases      []*fieldBase
	fieldsCanChange bool

	// Table configuration
	url        *xurl.Url
	filter     *group.FormGroup
	filterData map[string]any // Raw filter data from request
	flags      []string       // UI-only filter fields (excluded from parsed data)
	hasFilter  *bool          // Explicit hasFilter override (nil = use t.filter != nil)
	options    TableOptions
	outputType OutputType       // Current output mode (Web, CSV, PDF, Excel)
	components []core.Component // Additional components (charts, stats, progress bars, etc.)
	poll       int              // Auto-refresh poll interval in ms (0 = disabled); emitted for Web output
}

// Table is a generic type-safe table that maintains type safety internally
// while producing output compatible with the existing component.Table JSON format.
type Table[T any] struct {
	tableCore
	fields             []*field[T]
	data               []T
	treeAddSubAccessor func(T) bool // tree mode: per-row "+ sub" visibility (nil = button on all rows)
}

// TableOptions contains all table configuration options.
// These map directly to component.Table options for JSON compatibility.
type TableOptions struct {
	Class        *string
	Title        *string
	TextNoData   *string
	EmptyState   *emptystate.EmptyState
	ItemsPerPage *int
	PageSizes    []int
	ButtonsTop   []*button.TableButton
	Reload       *bool
	Dense        *bool
	Density      *Density
	Pagination   *bool
	Search       *bool
	MinWidth     *string
	Query        *bool
	// FilterCollapsed controls the expansion panel around the filter set via SetFilter.
	// nil = no panel at all (filter always open); see SetFilterCollapsed.
	FilterCollapsed *bool
	Csv             *bool
	Excel           *bool
	SaveState       *bool
	SaveStateId     *string
	SaveInput       *string
	SaveInputUrl    *string
	Borders         *bool
	BordersHeader   *bool
	Select          *bool
	SelectButtons   []*button.TableButton
	// UX-007 Bulk actions (additive). BulkActions, when non-empty, makes the frontend show a
	// selection column and a sticky context bar with these actions. SelectAllResults offers
	// "select all results" (whole filter, not just the page); StickyBulkBar pins the bar.
	BulkActions      []*button.TableButton
	SelectAllResults *bool
	StickyBulkBar    *bool
	Display          *string
	Footer           *bool
	ServerSide       *bool       // Enable server-side pagination (data fetched page-by-page)
	ScrollHeight     *string     // Custom scroll height for the table container (e.g., "400px", "80vh")
	EditUrl          *string     // URL for inline edit save requests (POST { id, field, value })
	Tree             *TreeConfig // Opt-in tree mode (indent + expand/collapse per row); nil = flat table
}

// TreeConfig configures the opt-in tree mode of a table.
// When set, the frontend renders rows hierarchically (indent + expand/collapse) based
// on the IdField/ParentIdField values of each flat row. When nil, the table is flat and
// the JSON output is byte-for-byte identical to a table built without Tree().
type TreeConfig struct {
	IdField            string    // Row field holding the node ID (required)
	ParentIdField      string    // Row field holding the parent node ID (required); null/0 = root
	TreeColumn         string    // Column that renders the indentation; empty = first column
	CollapseAllDefault bool      // Start collapsed (default: false → tree starts fully expanded)
	HideCounts         bool      // Hide the child-count "(5)" badge when collapsed (default: counts shown)
	PersistStateKey    string    // localStorage key for expand-state persistence; empty = no persistence
	AddSubURL          *xurl.Url // When set, renders a "+ sub" button per row; frontend navigates here ({id} placeholder)
	AddSubField        string    // Per-row field whose truthy value gates the "+ sub" button; set internally by TreeAddSubWhen
}

// ============================================================================
// Non-generic methods on tableCore (exist only once in the binary)
// ============================================================================

// GetURL returns the table URL
func (tc *tableCore) GetURL() *xurl.Url {
	return tc.url
}

// GetFilter returns the filter FormGroup
func (tc *tableCore) GetFilter() *group.FormGroup {
	return tc.filter
}

// GetFilterData returns the raw filter data from the request.
func (tc *tableCore) GetFilterData() map[string]any {
	return tc.filterData
}

// GetOutputType returns the current output type (Web, CSV, PDF, Excel).
func (tc *tableCore) GetOutputType() OutputType {
	return tc.outputType
}

// GetOptions returns the table options
func (tc *tableCore) GetOptions() TableOptions {
	return tc.options
}

// SetOutputType sets the output type for the table.
func (tc *tableCore) SetOutputType(output OutputType) {
	tc.outputType = output
}

// SetPoll sets the auto-refresh poll interval (in milliseconds). While set (> 0),
// the table data response (via wc.Data(tbl) / DataResponse / ToTableDataResponse)
// includes a "poll" field, causing the frontend to reload the whole table after the
// given interval and show a table-wide auto-refresh indicator with countdown.
// Pass 0 to disable. Typically set per request while a background worker is still
// running for one of the rows.
func (tc *tableCore) SetPoll(intervalMs int) {
	tc.poll = intervalMs
}

// AddButtonTop adds a button to the table's top toolbar.
func (tc *tableCore) AddButtonTop(btn *button.TableButton) {
	tc.options.ButtonsTop = append(tc.options.ButtonsTop, btn)
}

// LoadFilterData parses filter data from request, detects CSV flag, and returns parsed filter values.
// This is the main method controllers should use for handling table data requests.
//
// Returns:
//   - Parsed filter values (empty map if no filter)
//   - Error if parsing/validation fails
//
// Side effects:
//   - Sets outputType to OutputCSV if _csv flag is true
//   - Stores raw filter data in filterData field
//
// Example:
//
//	tbl := buildDeviceTable(ctx, translator)
//	parsedFilters, err := tbl.LoadFilterData(c)
//	if err != nil {
//	    return wc.BadRequest(err.Error())
//	}
//	rows := fetchDevicesWithFilters(parsedFilters)
//	tbl.SetData(rows)
//	return wc.TableDataFromTable(tbl)
func (tc *tableCore) LoadFilterData(c echo.Context) (map[string]any, error) {
	// Parse request body (contains filter fields + _csv flag)
	var requestData map[string]interface{}
	if err := c.Bind(&requestData); err != nil {
		// A malformed body must not be silently treated as "no filter" (which could
		// trigger an unfiltered full-table export). Empty bodies do not reach here:
		// echo's binder returns nil for zero-length content.
		return nil, err
	}

	// Check for CSV flag and set output type (only if CSV export is enabled in options)
	if tc.options.Csv == nil || *tc.options.Csv {
		if csvVal, ok := requestData["_csv"]; ok {
			if csvBool, isBool := csvVal.(bool); isBool && csvBool {
				tc.outputType = OutputCSV
			} else if csvStr, isStr := csvVal.(string); isStr && csvStr == "true" {
				tc.outputType = OutputCSV
			}
		}
	}

	// Check for Excel flag and set output type (only if Excel export is enabled in options)
	if tc.options.Excel == nil || *tc.options.Excel {
		if excelVal, ok := requestData["_excel"]; ok {
			if excelBool, isBool := excelVal.(bool); isBool && excelBool {
				tc.outputType = OutputExcel
			} else if excelStr, isStr := excelVal.(string); isStr && excelStr == "true" {
				tc.outputType = OutputExcel
			}
		}
	}

	// Store filter data (exclude _csv, _excel and flags)
	tc.filterData = make(map[string]any)
	for k, v := range requestData {
		if k == "_csv" || k == "_excel" {
			continue
		}
		// Check if field is a flag
		isFlag := false
		for _, flag := range tc.flags {
			if k == flag {
				isFlag = true
				break
			}
		}
		if !isFlag {
			tc.filterData[k] = v
		}
	}

	// Parse filter values (if filter exists)
	if tc.filter != nil {
		parsedFilters, err := tc.filter.ParseAndValidate(tc.filterData)
		if err != nil {
			return nil, err
		}
		return parsedFilters, nil
	}

	// No filter - return raw data excluding pagination params
	// (pagination params are kept in filterData for LoadPaginationParams)
	result := make(map[string]any)
	for k, v := range tc.filterData {
		// Exclude server-side pagination params from returned filter data
		if k == "_page" || k == "_pageSize" || k == "_sort" || k == "_sortDir" || k == "_search" {
			continue
		}
		result[k] = v
	}
	return result, nil
}

// PaginationParams holds server-side pagination parameters from request.
// These are extracted from request body when server-side pagination is enabled.
type PaginationParams struct {
	Page     int    // 0-based page index (from _page)
	PageSize int    // Items per page (from _pageSize)
	Sort     string // Column ID to sort by (from _sort, optional)
	SortDir  string // "asc" or "desc" (from _sortDir, optional)
	Search   string // Search text (from _search, optional)
}

// LoadPaginationParams extracts server-side pagination parameters from request body.
// Call this AFTER LoadFilterData() to get pagination parameters.
// Returns default values (page=0, pageSize=50) if parameters are not present.
//
// Example:
//
//	filters, _ := tbl.LoadFilterData(c)
//	pagination := tbl.LoadPaginationParams()
//	devices, total := dbm.Device.FindWithPagination(filters, pagination.Page, pagination.PageSize)
func (tc *tableCore) LoadPaginationParams() PaginationParams {
	params := PaginationParams{
		Page:     0,
		PageSize: 50, // Default page size
		SortDir:  "asc",
	}

	// Use ItemsPerPage from options as default if set
	if tc.options.ItemsPerPage != nil {
		params.PageSize = *tc.options.ItemsPerPage
	}

	// Extract from stored filter data (set by LoadFilterData)
	if tc.filterData == nil {
		return params
	}

	// _page (0-based)
	if pageVal, ok := tc.filterData["_page"]; ok {
		switch v := pageVal.(type) {
		case float64:
			params.Page = int(v)
		case int:
			params.Page = v
		case int64:
			params.Page = int(v)
		}
	}

	// _pageSize
	if pageSizeVal, ok := tc.filterData["_pageSize"]; ok {
		switch v := pageSizeVal.(type) {
		case float64:
			params.PageSize = int(v)
		case int:
			params.PageSize = v
		case int64:
			params.PageSize = int(v)
		}
	}

	// _sort
	if sortVal, ok := tc.filterData["_sort"]; ok {
		if s, ok := sortVal.(string); ok {
			params.Sort = s
		}
	}

	// _sortDir
	if sortDirVal, ok := tc.filterData["_sortDir"]; ok {
		if s, ok := sortDirVal.(string); ok && (s == "asc" || s == "desc") {
			params.SortDir = s
		}
	}

	// _search
	if searchVal, ok := tc.filterData["_search"]; ok {
		if s, ok := searchVal.(string); ok {
			params.Search = s
		}
	}

	// Validate sort against defined field IDs to prevent SQL injection
	if params.Sort != "" {
		valid := false
		for _, f := range tc.fieldBases {
			if f.id == params.Sort {
				valid = true
				break
			}
		}
		if !valid {
			params.Sort = ""
		}
	}

	// Clamp page size to prevent memory exhaustion
	if params.PageSize < 1 {
		params.PageSize = 50
	}
	if params.PageSize > 1000 {
		params.PageSize = 1000
	}

	// Prevent negative page index
	if params.Page < 0 {
		params.Page = 0
	}

	// Limit search string length to prevent performance issues
	if len(params.Search) > 200 {
		params.Search = params.Search[:200]
	}

	return params
}

// exportFields converts fieldBase array to JSON array for component output.
// Each field is converted to tableFieldJSON format for JSON serialization.
// Hidden fields are excluded from the output.
func (tc *tableCore) exportFields(ctx *core.UiContext) []map[string]any {
	fields := make([]map[string]any, 0, len(tc.fieldBases))
	for _, f := range tc.fieldBases {
		if f.hide {
			continue
		}
		jsonField := f.toTableField()
		fields = append(fields, jsonField.print(ctx))
	}
	return fields
}

// exportFieldsForCSV converts fieldBase array to JSON array for CSV export only.
// This filters out fields where csv=false or fields that are hidden.
func (tc *tableCore) exportFieldsForCSV(ctx *core.UiContext) []map[string]any {
	csvFields := make([]map[string]any, 0, len(tc.fieldBases))
	for _, f := range tc.fieldBases {
		if f.hide {
			continue
		}
		if !f.csv {
			continue
		}
		jsonField := f.toTableField()
		csvFields = append(csvFields, jsonField.print(ctx))
	}
	return csvFields
}

// exportOptions converts TableOptions to JSON map for component output.
// Matches component.Table options format exactly.
func (tc *tableCore) exportOptions(ctx *core.UiContext) map[string]any {
	opts := tc.options
	options := make(map[string]any)

	// Add all options that are set
	if opts.Class != nil {
		options["class"] = *opts.Class
	}
	if opts.Title != nil {
		options["title"] = *opts.Title
	}
	if opts.TextNoData != nil {
		options["textNoData"] = core.Translate(ctx, *opts.TextNoData)
	}
	if opts.EmptyState != nil {
		options["emptyState"] = opts.EmptyState.PrintData(ctx)
	}
	if opts.ItemsPerPage != nil {
		options["itemsPerPage"] = *opts.ItemsPerPage
	}
	if opts.PageSizes != nil && len(opts.PageSizes) > 0 {
		options["pageSizes"] = opts.PageSizes
	}

	// ButtonsTop: Add CSV button if enabled, then export all buttons
	var topButtons []*button.TableButton
	if opts.ButtonsTop != nil {
		topButtons = append(topButtons, opts.ButtonsTop...)
	}

	// Auto-generate CSV download button if CSV option enabled and URL exists
	if opts.Csv != nil && *opts.Csv && tc.url != nil {
		csvBtn := button.NewTableButton(
			core.ButtonActionDownload,
			"csv",
			tc.url,
			"CSV",
			core.ColorAccent,
			false,
			nil,
		).WithData(map[string]any{"_csv": true})
		topButtons = append(topButtons, csvBtn)
	}

	// Auto-generate Excel download button if Excel option enabled and URL exists
	if opts.Excel != nil && *opts.Excel && tc.url != nil {
		excelBtn := button.NewTableButton(
			core.ButtonActionDownload,
			"explicit",
			tc.url,
			"Excel",
			core.ColorAccent,
			false,
			nil,
		).WithData(map[string]any{"_excel": true})
		topButtons = append(topButtons, excelBtn)
	}

	// Export ButtonsTop in same format as SelectButtons
	if len(topButtons) > 0 {
		buttons := make([]map[string]any, len(topButtons))
		for i, btn := range topButtons {
			buttons[i] = btn.Print(ctx)
		}
		options["buttons"] = map[string]any{"buttons": buttons}
	}

	if opts.Reload != nil {
		options["reload"] = *opts.Reload
	}
	if opts.Dense != nil {
		options["dense"] = *opts.Dense
	}
	if opts.Density != nil {
		options["density"] = string(*opts.Density)
	}
	if opts.Pagination != nil {
		options["pagination"] = *opts.Pagination
	}
	if opts.Search != nil {
		options["search"] = *opts.Search
	}
	if opts.MinWidth != nil {
		options["minWidth"] = *opts.MinWidth
	}
	if opts.Query != nil {
		options["query"] = *opts.Query
	}
	if opts.Csv != nil {
		options["csv"] = *opts.Csv
	}
	if opts.SaveState != nil && opts.SaveStateId != nil {
		options["saveState"] = *opts.SaveState
		options["saveStateId"] = *opts.SaveStateId
	}
	if opts.SaveInput != nil {
		options["saveInput"] = *opts.SaveInput
	}
	if opts.SaveInputUrl != nil {
		options["saveInputUrl"] = *opts.SaveInputUrl
	}
	if opts.Borders != nil {
		options["borders"] = *opts.Borders
	}
	if opts.BordersHeader != nil {
		options["bordersHeader"] = *opts.BordersHeader
	}
	if opts.Select != nil {
		options["select"] = *opts.Select
	}
	// SelectButtons: serialize each button component
	if opts.SelectButtons != nil && len(opts.SelectButtons) > 0 {
		buttons := make([]map[string]any, len(opts.SelectButtons))
		for i, btn := range opts.SelectButtons {
			buttons[i] = btn.Print(ctx)
		}
		options["selectButtons"] = buttons
	}
	// BulkActions (UX-007): serialize like SelectButtons; only when set.
	if len(opts.BulkActions) > 0 {
		buttons := make([]map[string]any, len(opts.BulkActions))
		for i, btn := range opts.BulkActions {
			buttons[i] = btn.Print(ctx)
		}
		options["bulkActions"] = buttons
	}
	if opts.SelectAllResults != nil {
		options["selectAllResults"] = *opts.SelectAllResults
	}
	if opts.StickyBulkBar != nil {
		options["stickyBulkBar"] = *opts.StickyBulkBar
	}
	if opts.Footer != nil {
		options["footer"] = *opts.Footer
	}
	if opts.ServerSide != nil {
		options["serverSide"] = *opts.ServerSide
	}
	if opts.ScrollHeight != nil {
		options["scrollHeight"] = *opts.ScrollHeight
	}
	if opts.EditUrl != nil {
		options["editUrl"] = *opts.EditUrl
	}

	// Tree mode: serialize the nested tree settings object (only when opted in).
	// Keys must match the frontend XiriTableTreeSettings interface exactly.
	if opts.Tree != nil {
		tree := map[string]any{
			"idField":              opts.Tree.IdField,
			"parentIdField":        opts.Tree.ParentIdField,
			"collapseAllByDefault": opts.Tree.CollapseAllDefault,
			"showCounts":           !opts.Tree.HideCounts,
		}
		if opts.Tree.TreeColumn != "" {
			tree["treeColumn"] = opts.Tree.TreeColumn
		}
		if opts.Tree.PersistStateKey != "" {
			tree["persistStateKey"] = opts.Tree.PersistStateKey
		}
		if opts.Tree.AddSubURL != nil {
			tree["addSubUrl"] = opts.Tree.AddSubURL.Print()
		}
		if opts.Tree.AddSubField != "" {
			tree["addSubField"] = opts.Tree.AddSubField
		}
		options["tree"] = tree
	}

	return options
}

// printComponent builds the component JSON for the Angular frontend.
// staticData should be nil for AJAX mode, or the pre-formatted data for static mode.
func (tc *tableCore) printComponent(ctx *core.UiContext, staticData []map[string]any) map[string]any {
	// Build base component structure
	result := map[string]any{
		"type": "table",
	}

	// Add display class if set
	if tc.options.Display != nil {
		result["display"] = *tc.options.Display
	}

	// Build data section
	dataSection := make(map[string]any)

	// Add filter flag
	if tc.hasFilter != nil {
		dataSection["hasFilter"] = *tc.hasFilter
	} else {
		dataSection["hasFilter"] = tc.filter != nil
	}

	// Add fields
	dataSection["fields"] = tc.exportFields(ctx)

	// Add options
	dataSection["options"] = tc.exportOptions(ctx)

	// Determine mode: AJAX (url set) vs Static (data set)
	if tc.url != nil {
		// AJAX mode: URL for dynamic loading
		dataSection["url"] = tc.url.PrintPrefix()
		dataSection["data"] = nil
		dataSection["components"] = nil
	} else {
		// Static mode: embedded data
		dataSection["url"] = nil
		dataSection["data"] = staticData
		dataSection["components"] = nil
	}

	result["data"] = dataSection

	// If filter exists, automatically wrap table in Query component
	if tc.filter != nil {
		// Extra-Daten aus Form=false-Feldern sammeln
		fields := tc.filter.GetFields()
		extraData := make(map[string]any)
		for _, f := range fields {
			if !f.GetForm() {
				extraData[f.GetID()] = f.GetDefault()
			}
		}

		// Sichtbare Felder exportieren (ExportForFrontend überspringt Form=false)
		filterForm := tc.filter.ExportForFrontend()

		saveStateId := tc.options.SaveStateId
		q := query.NewQuery(filterForm, saveStateId, tc.options.Display)

		if tc.options.FilterCollapsed != nil {
			q.Collapsed(*tc.options.FilterCollapsed)
		}

		if len(extraData) > 0 {
			q.SetExtraData(extraData)
		}

		q.AddArray(result)
		return q.Print(ctx)
	}

	return result
}

// hideFieldByID hides a single field by its ID in the fieldBases slice.
// Returns true if the field was found.
func (tc *tableCore) hideFieldByID(fieldID string) bool {
	for _, f := range tc.fieldBases {
		if f.id == fieldID {
			f.hide = true
			return true
		}
	}
	return false
}

// showFieldByID shows a single field by its ID in the fieldBases slice.
// Returns true if the field was found.
func (tc *tableCore) showFieldByID(fieldID string) bool {
	for _, f := range tc.fieldBases {
		if f.id == fieldID {
			f.hide = false
			return true
		}
	}
	return false
}

// hideFieldsByID hides multiple fields by their IDs in the fieldBases slice.
func (tc *tableCore) hideFieldsByID(fieldIDs []string) {
	if len(fieldIDs) == 0 {
		return
	}
	hideMap := make(map[string]bool, len(fieldIDs))
	for _, id := range fieldIDs {
		hideMap[id] = true
	}
	for _, f := range tc.fieldBases {
		if hideMap[f.id] {
			f.hide = true
		}
	}
}

// showFieldsByID shows multiple fields by their IDs in the fieldBases slice.
func (tc *tableCore) showFieldsByID(fieldIDs []string) {
	if len(fieldIDs) == 0 {
		return
	}
	showMap := make(map[string]bool, len(fieldIDs))
	for _, id := range fieldIDs {
		showMap[id] = true
	}
	for _, f := range tc.fieldBases {
		if showMap[f.id] {
			f.hide = false
		}
	}
}

// HasEditableField checks if a field with the given ID exists and is marked as editable.
// Use this to validate incoming inline-edit requests before processing them.
func (tc *tableCore) HasEditableField(fieldID string) bool {
	for _, f := range tc.fieldBases {
		if f.id == fieldID {
			return f.editable
		}
	}
	return false
}

// InlineEditRequest represents the payload sent by the Angular frontend for inline cell edits.
// The frontend sends POST { id, field, value } to the editUrl.
type InlineEditRequest struct {
	ID    int64  `json:"id"`
	Field string `json:"field"`
	Value any    `json:"value"`
}

// ParseInlineEdit parses the inline edit request from the Echo context and validates
// that the field exists and is editable in this table.
//
// Returns an error if:
//   - The request body cannot be parsed
//   - The field does not exist in the table
//   - The field is not marked as editable
//
// Example:
//
//	func (ctrl *Controller) InlineEditSave(c echo.Context) error {
//	    tbl := buildTable(ctx)
//	    req, err := tbl.ParseInlineEdit(c)
//	    if err != nil {
//	        return c.JSON(400, response.NewErrorResponse(err.Error()))
//	    }
//	    // Use req.ID, req.Field, req.Value for business logic
//	}
func (tc *tableCore) ParseInlineEdit(c echo.Context) (InlineEditRequest, error) {
	var req InlineEditRequest
	if err := c.Bind(&req); err != nil {
		return req, err
	}
	if !tc.HasEditableField(req.Field) {
		return req, echo.NewHTTPError(400, "field not editable: "+req.Field)
	}
	return req, nil
}

// GetFieldMetas returns metadata for all fields. This is the public API for
// external consumers that need field information (renderers, exporters).
func (tc *tableCore) GetFieldMetas() []FieldMeta {
	metas := make([]FieldMeta, len(tc.fieldBases))
	for i, f := range tc.fieldBases {
		metas[i] = FieldMeta{
			ID:     f.id,
			Name:   f.name,
			Type:   string(f.fieldType),
			Hidden: f.hide,
			Align:  f.align,
			CSV:    f.csv,

			Header:     f.header,
			HeaderSpan: f.headerSpan,
		}
	}
	return metas
}

// ============================================================================
// Generic methods on Table[T] (must use T)
// ============================================================================

// GetData returns formatted table data for a specific output type.
// This is where the magic happens: raw row structs are converted to formatted map[string]any
// with all formatters applied and locale/unit conversions done automatically.
func (t *Table[T]) GetData(ctx *core.UiContext, output OutputType) []map[string]any {
	// Build field accessor map for Row interface
	fieldMap := t.buildFieldMap()

	rows := make([]map[string]any, len(t.data))

	for i, rowData := range t.data {
		// Create Row wrapper for cross-field access
		row := newTypedRow(rowData, fieldMap)

		// Extract and format each field
		rowMap := make(map[string]any)

		for _, f := range t.fields {
			// Skip hidden fields
			if f.hide {
				continue
			}

			// Skip non-CSV fields for CSV/Excel output
			if (output == OutputCSV || output == OutputExcel) && !f.csv {
				continue
			}

			// Extract value using accessor
			value := f.accessor(rowData)

			// Format with access to entire row (for cross-field dependencies)
			formatted := f.format(value, row, output, ctx)

			// Special handling for Link fields in OutputWeb/OutputPDF
			// Link fields output TWO fields: fieldName (display text) and fieldNameLink (URL)
			if output == OutputWeb && f.fieldTypeHint == link {
				// Extract [2]string array from formatted value
				if linkArray, ok := formatted.([2]string); ok {
					rowMap[f.id] = linkArray[0]        // Display text
					rowMap[f.id+"Link"] = linkArray[1] // URL
				} else {
					// Fallback for invalid data
					rowMap[f.id] = ""
					rowMap[f.id+"Link"] = ""
				}
			} else {
				// Normal field: single value
				rowMap[f.id] = formatted
			}

			// Add per-row hint for icon fields with hintAccessor
			if output == OutputWeb && f.hintAccessor != nil {
				hintStr := f.hintAccessor(rowData)
				if hintStr != "" {
					rowMap[f.id+"Hint"] = hintStr
				}
			}

			// Inject menu data into button row data
			if output == OutputWeb && len(f.menuAccessors) > 0 {
				if buttonMap, ok := rowMap[f.id].(map[string]any); ok {
					for key, menuAccessor := range f.menuAccessors {
						keyStr := strconv.Itoa(key)
						if val, exists := buttonMap[keyStr]; exists && val == false {
							continue
						}
						menuData := menuAccessor(rowData)
						if menuData == nil {
							buttonMap[keyStr] = false
							continue
						}
						result := make([]any, len(menuData))
						for j, v := range menuData {
							if v == "" {
								result[j] = false
							} else {
								result[j] = v
							}
						}
						buttonMap[keyStr] = result
					}
				}
			}
		}

		// Tree mode: per-row "+ sub" visibility flag (read by the frontend via tree.addSubField).
		if output == OutputWeb && t.treeAddSubAccessor != nil {
			rowMap[treeAddSubKey] = t.treeAddSubAccessor(rowData)
		}

		rows[i] = rowMap
	}

	return rows
}

// treeAddSubKey is the reserved row-data key carrying the per-row "+ sub" visibility flag.
const treeAddSubKey = "_addSub"

// buildFieldMap creates accessor map for Row interface.
// This allows formatters to access any field value in the row for cross-field dependencies.
func (t *Table[T]) buildFieldMap() map[string]func(T) any {
	fieldMap := make(map[string]func(T) any)
	for _, f := range t.fields {
		fieldMap[f.id] = f.accessor
	}
	return fieldMap
}

// CalculateFooter computes footer aggregations for all fields with footer enabled.
// Returns a map of field_id -> aggregated_value (formatted).
func (t *Table[T]) CalculateFooter(ctx *core.UiContext, output OutputType) map[string]any {
	footer := make(map[string]any)
	fieldMap := t.buildFieldMap()

	for _, f := range t.fields {
		if f.footer == FieldFooterNo {
			continue
		}

		var aggregated any

		switch f.footer {
		case FieldFooterSum:
			aggregated = t.sumField(f)
		case FieldFooterCount:
			aggregated = t.countField(f)
		case FieldFooterStatic:
			// Static footer values would be set separately
			continue
		}

		// Format footer value
		// Use first row for Row context (for device-specific formatting)
		if len(t.data) > 0 {
			row := newTypedRow(t.data[0], fieldMap)
			formatted := f.format(aggregated, row, output, ctx)
			footer[f.id] = formatted
		} else {
			// For empty tables, use raw aggregated value without formatting
			// to avoid nil row access in formatters that reference other fields
			footer[f.id] = aggregated
		}
	}

	return footer
}

// sumField sums all values for a field
func (t *Table[T]) sumField(f *field[T]) float64 {
	sum := 0.0
	for _, rowData := range t.data {
		val := f.accessor(rowData)
		sum += toFloat64(val)
	}
	return sum
}

// countField counts non-empty values for a field
func (t *Table[T]) countField(f *field[T]) int {
	count := 0
	for _, rowData := range t.data {
		val := f.accessor(rowData)
		if val != nil && val != "" {
			count++
		}
	}
	return count
}

// toFloat64 converts any numeric value to float64
func toFloat64(value any) float64 {
	if value == nil {
		return 0.0
	}

	switch v := value.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	}

	return 0.0
}

// ============================================================================
// Table Mutation Methods (generic thin wrappers)
// ============================================================================

// SetData sets or updates the table data after building.
func (t *Table[T]) SetData(data []T) {
	t.data = data
	t.url = nil // Clear URL to force static mode
}

// SetURL sets or updates the AJAX data URL after building.
func (t *Table[T]) SetURL(url *xurl.Url) {
	t.url = url
	t.data = nil // Clear data to force AJAX mode
}

// AddComponent adds a component to display alongside the table.
func (t *Table[T]) AddComponent(comp core.Component) *Table[T] {
	t.components = append(t.components, comp)
	return t
}

// HideField hides a single field by its ID.
func (t *Table[T]) HideField(fieldID string) *Table[T] {
	if !t.tableCore.hideFieldByID(fieldID) {
		slog.Warn("HideField: field not found", "fieldID", fieldID)
	}
	return t
}

// ShowField shows a previously hidden field by its ID.
func (t *Table[T]) ShowField(fieldID string) *Table[T] {
	if !t.tableCore.showFieldByID(fieldID) {
		slog.Warn("ShowField: field not found", "fieldID", fieldID)
	}
	return t
}

// HideFields hides multiple fields by their IDs.
func (t *Table[T]) HideFields(fieldIDs ...string) *Table[T] {
	t.tableCore.hideFieldsByID(fieldIDs)
	return t
}

// ShowFields shows multiple previously hidden fields by their IDs.
func (t *Table[T]) ShowFields(fieldIDs ...string) *Table[T] {
	t.tableCore.showFieldsByID(fieldIDs)
	return t
}

// SetFilterData manually sets filter data.
func (t *Table[T]) SetFilterData(data map[string]any) *Table[T] {
	t.filterData = data
	return t
}

// SetFlags sets UI-only filter fields that should be excluded from parsed data.
func (t *Table[T]) SetFlags(flags ...string) *Table[T] {
	t.flags = flags
	return t
}

// ============================================================================
// Component Interface Implementation
// ============================================================================

// Print implements the core.Component interface, returning JSON for the Angular frontend.
// The generic wrapper delegates to non-generic printComponent, only calling GetData when needed.
func (t *Table[T]) Print(ctx *core.UiContext) map[string]any {
	// For static mode (no URL), get data first (this is the only T-dependent part)
	var staticData []map[string]any
	if t.url == nil {
		staticData = t.GetData(ctx, OutputWeb)
	}

	// Delegate to non-generic printComponent
	return t.tableCore.printComponent(ctx, staticData)
}

// ExportFields is the public version of exportFields that allows external packages
// (like dialog) to access field definitions for building dialog tables.
func (t *Table[T]) ExportFields(ctx *core.UiContext) []map[string]any {
	return t.tableCore.exportFields(ctx)
}
