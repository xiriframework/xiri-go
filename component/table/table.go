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
	"github.com/xiriframework/xiri-go/uicontext"
)

// Table is a generic type-safe table that maintains type safety internally
// while producing output compatible with the existing component.Table JSON format.
type Table[T any] struct {
	fields          []*Field[T]
	fieldsCanChange bool
	data            []T
	ctx             *uicontext.UiContext
	translator      core.TranslateFunc

	// Table configuration
	url        *xurl.Url
	filter     *group.FormGroup
	filterData map[string]any // Raw filter data from request
	flags      []string       // UI-only filter fields (excluded from parsed data)
	hasFilter  *bool          // Explicit hasFilter override (nil = use t.filter != nil)
	options    TableOptions
	outputType OutputType       // Current output mode (Web, CSV, PDF, Excel)
	components []core.Component // Additional components (charts, stats, progress bars, etc.)
}

// TableOptions contains all table configuration options.
// These map directly to component.Table options for JSON compatibility.
type TableOptions struct {
	Class         *string
	Title         *string
	TextNoData    *string
	EmptyState    *emptystate.EmptyState
	ItemsPerPage  *int
	PageSizes     []int
	ButtonsTop    []*button.TableButton
	Reload        *bool
	Dense         *bool
	Pagination    *bool
	Search        *bool
	MinWidth      *string
	Query         *bool
	Csv           *bool
	Excel         *bool
	SaveState     *bool
	SaveStateId   *string
	SaveInput     *string
	SaveInputUrl  *string
	Borders       *bool
	BordersHeader *bool
	Select        *bool
	SelectButtons []*button.TableButton
	Display       *string
	Footer        *bool
	ServerSide    *bool   // Enable server-side pagination (data fetched page-by-page)
	ScrollHeight  *string // Custom scroll height for the table container (e.g., "400px", "80vh")
}

// GetData returns formatted table data for a specific output type.
// This is where the magic happens: raw row structs are converted to formatted map[string]any
// with all formatters applied and locale/unit conversions done automatically.
func (t *Table[T]) GetData(output OutputType) []map[string]any {
	// Build field accessor map for Row interface
	fieldMap := t.buildFieldMap()

	rows := make([]map[string]any, len(t.data))

	for i, rowData := range t.data {
		// Create Row wrapper for cross-field access
		row := NewTypedRow(rowData, fieldMap)

		// Extract and format each field
		rowMap := make(map[string]any)

		for _, field := range t.fields {
			// Skip hidden fields
			if field.IsHidden() {
				continue
			}

			// Skip non-CSV fields for CSV/Excel output
			if (output == OutputCSV || output == OutputExcel) && !field.IsCsvEnabled() {
				continue
			}

			// Extract value using accessor
			value := field.GetAccessor()(rowData)

			// Format with access to entire row (for cross-field dependencies)
			formatted := field.Format(value, row, output, t.ctx)

			// Special handling for Link fields in OutputWeb/OutputPDF
			// Link fields output TWO fields: fieldName (display text) and fieldNameLink (URL)
			if output == OutputWeb && field.GetFieldTypeHint() == Link {
				// Extract [2]string array from formatted value
				if linkArray, ok := formatted.([2]string); ok {
					rowMap[field.GetID()] = linkArray[0]        // Display text
					rowMap[field.GetID()+"Link"] = linkArray[1] // URL
				} else {
					// Fallback for invalid data
					rowMap[field.GetID()] = ""
					rowMap[field.GetID()+"Link"] = ""
				}
			} else {
				// Normal field: single value
				rowMap[field.GetID()] = formatted
			}

			// Add per-row hint for icon fields with hintAccessor
			if output == OutputWeb && field.GetHintAccessor() != nil {
				hintStr := field.GetHintAccessor()(rowData)
				if hintStr != "" {
					rowMap[field.GetID()+"Hint"] = hintStr
				}
			}

			// Inject menu data into button row data
			if output == OutputWeb && len(field.GetMenuAccessors()) > 0 {
				if buttonMap, ok := rowMap[field.GetID()].(map[string]any); ok {
					for key, menuAccessor := range field.GetMenuAccessors() {
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

		rows[i] = rowMap
	}

	return rows
}

// buildFieldMap creates accessor map for Row interface.
// This allows formatters to access any field value in the row for cross-field dependencies.
func (t *Table[T]) buildFieldMap() map[string]func(T) any {
	fieldMap := make(map[string]func(T) any)
	for _, field := range t.fields {
		fieldMap[field.GetID()] = field.GetAccessor()
	}
	return fieldMap
}

// GetFields returns all fields
func (t *Table[T]) GetFields() []*Field[T] {
	return t.fields
}

// GetContext returns the UiContext
func (t *Table[T]) GetContext() *uicontext.UiContext {
	return t.ctx
}

// GetTranslator returns the translator function
func (t *Table[T]) GetTranslator() core.TranslateFunc {
	return t.translator
}

// GetURL returns the table URL
func (t *Table[T]) GetURL() *xurl.Url {
	return t.url
}

// GetFilter returns the filter FormGroup
func (t *Table[T]) GetFilter() *group.FormGroup {
	return t.filter
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
func (t *Table[T]) LoadFilterData(c echo.Context) (map[string]any, error) {
	// Parse request body (contains filter fields + _csv flag)
	var requestData map[string]interface{}
	if err := c.Bind(&requestData); err != nil {
		slog.Debug("LoadFilterData: failed to bind request body, using empty map", "error", err)
		requestData = make(map[string]interface{})
	}

	// Check for CSV flag and set output type
	if csvVal, ok := requestData["_csv"]; ok {
		if csvBool, isBool := csvVal.(bool); isBool && csvBool {
			t.outputType = OutputCSV
		} else if csvStr, isStr := csvVal.(string); isStr && csvStr == "true" {
			t.outputType = OutputCSV
		}
	}

	// Check for Excel flag and set output type
	if excelVal, ok := requestData["_excel"]; ok {
		if excelBool, isBool := excelVal.(bool); isBool && excelBool {
			t.outputType = OutputExcel
		} else if excelStr, isStr := excelVal.(string); isStr && excelStr == "true" {
			t.outputType = OutputExcel
		}
	}

	// Store filter data (exclude _csv, _excel and flags)
	t.filterData = make(map[string]any)
	for k, v := range requestData {
		if k == "_csv" || k == "_excel" {
			continue
		}
		// Check if field is a flag
		isFlag := false
		for _, flag := range t.flags {
			if k == flag {
				isFlag = true
				break
			}
		}
		if !isFlag {
			t.filterData[k] = v
		}
	}

	// Parse filter values (if filter exists)
	if t.filter != nil {
		parsedFilters, err := t.filter.ParseAndValidate(t.filterData)
		if err != nil {
			return nil, err
		}
		return parsedFilters, nil
	}

	// No filter - return raw data excluding pagination params
	// (pagination params are kept in filterData for LoadPaginationParams)
	result := make(map[string]any)
	for k, v := range t.filterData {
		// Exclude server-side pagination params from returned filter data
		if k == "_page" || k == "_pageSize" || k == "_sort" || k == "_sortDir" || k == "_search" {
			continue
		}
		result[k] = v
	}
	return result, nil
}

// SetFilterData manually sets filter data (for testing or special cases).
// For normal use, prefer LoadFilterData() which handles request parsing automatically.
func (t *Table[T]) SetFilterData(data map[string]any) *Table[T] {
	t.filterData = data
	return t
}

// GetFilterData returns the raw filter data from the request.
func (t *Table[T]) GetFilterData() map[string]any {
	return t.filterData
}

// AddButtonTop adds a button to the table's top toolbar.
// Useful for adding action buttons (PDF, Excel, etc.) after building the table.
func (t *Table[T]) AddButtonTop(btn *button.TableButton) {
	t.options.ButtonsTop = append(t.options.ButtonsTop, btn)
}

// SetFlags sets UI-only filter fields that should be excluded from parsed data.
// Flags are typically used for frontend state that shouldn't be sent to the backend.
func (t *Table[T]) SetFlags(flags ...string) *Table[T] {
	t.flags = flags
	return t
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
func (t *Table[T]) LoadPaginationParams() PaginationParams {
	params := PaginationParams{
		Page:     0,
		PageSize: 50, // Default page size
		SortDir:  "asc",
	}

	// Use ItemsPerPage from options as default if set
	if t.options.ItemsPerPage != nil {
		params.PageSize = *t.options.ItemsPerPage
	}

	// Extract from stored filter data (set by LoadFilterData)
	if t.filterData == nil {
		return params
	}

	// _page (0-based)
	if pageVal, ok := t.filterData["_page"]; ok {
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
	if pageSizeVal, ok := t.filterData["_pageSize"]; ok {
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
	if sortVal, ok := t.filterData["_sort"]; ok {
		if s, ok := sortVal.(string); ok {
			params.Sort = s
		}
	}

	// _sortDir
	if sortDirVal, ok := t.filterData["_sortDir"]; ok {
		if s, ok := sortDirVal.(string); ok && (s == "asc" || s == "desc") {
			params.SortDir = s
		}
	}

	// _search
	if searchVal, ok := t.filterData["_search"]; ok {
		if s, ok := searchVal.(string); ok {
			params.Search = s
		}
	}

	return params
}

// GetOptions returns the table options
func (t *Table[T]) GetOptions() TableOptions {
	return t.options
}

// CalculateFooter computes footer aggregations for all fields with footer enabled.
// Returns a map of field_id -> aggregated_value (formatted).
func (t *Table[T]) CalculateFooter(output OutputType) map[string]any {
	footer := make(map[string]any)
	fieldMap := t.buildFieldMap()

	for _, field := range t.fields {
		if field.GetFooter() == FieldFooterNo {
			continue
		}

		var aggregated any

		switch field.GetFooter() {
		case FieldFooterSum:
			aggregated = t.sumField(field)
		case FieldFooterCount:
			aggregated = t.countField(field)
		case FieldFooterStatic:
			// Static footer values would be set separately
			continue
		}

		// Format footer value
		// Use first row for Row context (for device-specific formatting)
		if len(t.data) > 0 {
			row := NewTypedRow(t.data[0], fieldMap)
			formatted := field.Format(aggregated, row, output, t.ctx)
			footer[field.GetID()] = formatted
		} else {
			// For empty tables, use raw aggregated value without formatting
			// to avoid nil row access in formatters that reference other fields
			footer[field.GetID()] = aggregated
		}
	}

	return footer
}

// sumField sums all values for a field
func (t *Table[T]) sumField(field *Field[T]) float64 {
	sum := 0.0
	accessor := field.GetAccessor()

	for _, rowData := range t.data {
		val := accessor(rowData)
		sum += toFloat64(val)
	}

	return sum
}

// countField counts non-empty values for a field
func (t *Table[T]) countField(field *Field[T]) int {
	count := 0
	accessor := field.GetAccessor()

	for _, rowData := range t.data {
		val := accessor(rowData)
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
// Table Mutation Methods
// ============================================================================

// SetData sets or updates the table data after building.
// This allows reusing the same table definition for both page and data endpoints.
// When data is set, the URL is cleared to force static mode.
//
// Example:
//
//	// Shared table definition
//	tbl := buildDeviceTable(ctx, translator)
//
//	// Data endpoint: set data and return formatted rows
//	tbl.SetData(deviceRows)
//	return tbl.GetData(OutputWeb)
func (t *Table[T]) SetData(data []T) {
	t.data = data
	t.url = nil // Clear URL to force static mode
}

// SetURL sets or updates the AJAX data URL after building.
// This allows reusing the same table definition for both page and data endpoints.
// When URL is set, the data is cleared to force AJAX mode.
//
// Example:
//
//	// Shared table definition
//	tbl := buildDeviceTable(ctx, translator)
//
//	// Page endpoint: set URL and return component definition
//	tbl.SetURL(component.NewUrl("/Portal/Device/TableData", ""))
//	return tbl.Print(translator)
func (t *Table[T]) SetURL(url *xurl.Url) {
	t.url = url
	t.data = nil // Clear data to force AJAX mode
}

// GetOutputType returns the current output type (Web, CSV, PDF, Excel).
// This determines how data is formatted when calling GetData() or ToTableDataResponse().
func (t *Table[T]) GetOutputType() OutputType {
	return t.outputType
}

// AddComponent adds a component (e.g., MultiProgress, Chart, Info) to display alongside the table.
// Components are included in Web output but excluded from CSV/PDF/Excel exports.
//
// This is useful for adding statistics, charts, or informational messages above/below the table.
// Components are rendered in the order they are added.
//
// Example usage:
//
//	// Add statistics progress bars
//	kundeProgress := component.NewMultiProgress("Customers", 5, true, nil)
//	kundeProgress.AddLine("Customer A", 10, component.ColorPrimary, nil)
//	table.AddComponent(kundeProgress)
//
//	// Add informational message
//	info := component.NewInfo("Note: Data refreshed hourly", core.ColorAccent)
//	table.AddComponent(info)
//
//	return wc.TableDataFromTable(table)
func (t *Table[T]) AddComponent(comp core.Component) *Table[T] {
	t.components = append(t.components, comp)
	return t
}

// HideField hides a single field by its ID. The field will be excluded from GetData() output.
// This can be called after Build() to conditionally hide columns based on filters or permissions.
//
// Example:
//
//	tbl := builder.Build()
//	if !user.IsAdmin() {
//	    tbl.HideField("internal_id")
//	}
//	tbl.SetData(rows)
func (t *Table[T]) HideField(fieldID string) *Table[T] {
	for _, field := range t.fields {
		if field.GetID() == fieldID {
			field.setHide(true)
			return t
		}
	}
	slog.Warn("HideField: field not found", "fieldID", fieldID)
	return t
}

// ShowField shows a previously hidden field by its ID.
// This can be called after Build() to conditionally show columns.
//
// Example:
//
//	tbl := builder.Build()
//	if user.HasPermission("view_costs") {
//	    tbl.ShowField("cost")
//	}
//	tbl.SetData(rows)
func (t *Table[T]) ShowField(fieldID string) *Table[T] {
	for _, field := range t.fields {
		if field.GetID() == fieldID {
			field.setHide(false)
			return t
		}
	}
	slog.Warn("ShowField: field not found", "fieldID", fieldID)
	return t
}

// HideFields hides multiple fields by their IDs. The fields will be excluded from GetData() output.
// This is more efficient than calling HideField multiple times.
//
// Example:
//
//	tbl := builder.Build()
//	if !user.IsAdmin() {
//	    tbl.HideFields("device_id", "internal_notes", "cost")
//	}
//	tbl.SetData(rows)
func (t *Table[T]) HideFields(fieldIDs ...string) *Table[T] {
	if len(fieldIDs) == 0 {
		return t
	}

	// Build map for O(1) lookup
	hideMap := make(map[string]bool, len(fieldIDs))
	for _, id := range fieldIDs {
		hideMap[id] = true
	}

	// Set hide flag on matching fields
	for _, field := range t.fields {
		if hideMap[field.GetID()] {
			field.setHide(true)
		}
	}
	return t
}

// ShowFields shows multiple previously hidden fields by their IDs.
// This is more efficient than calling ShowField multiple times.
//
// Example:
//
//	tbl := builder.Build()
//	if user.HasPermission("view_all") {
//	    tbl.ShowFields("device_id", "internal_notes", "cost")
//	}
//	tbl.SetData(rows)
func (t *Table[T]) ShowFields(fieldIDs ...string) *Table[T] {
	if len(fieldIDs) == 0 {
		return t
	}

	// Build map for O(1) lookup
	showMap := make(map[string]bool, len(fieldIDs))
	for _, id := range fieldIDs {
		showMap[id] = true
	}

	// Clear hide flag on matching fields
	for _, field := range t.fields {
		if showMap[field.GetID()] {
			field.setHide(false)
		}
	}
	return t
}

// SetOutputType sets the output type for the table.
// This affects data formatting in subsequent GetData() or ToTableDataResponse() calls.
//
// Example:
//
//	// Check for CSV request parameter
//	if wc.QueryParam("_csv") == "true" {
//	    tbl.SetOutputType(table.OutputCSV)
//	}
//	return wc.TableDataFromTable(tbl)
func (t *Table[T]) SetOutputType(output OutputType) {
	t.outputType = output
}

// ============================================================================
// Component Interface Implementation
// ============================================================================

// Print implements the core.Component interface, returning JSON for the Angular frontend.
//
// The output format depends on table configuration:
// - AJAX mode (url != nil): Component definition with URL for dynamic data loading
// - Static mode (url == nil, data != nil): Component definition with embedded data
func (t *Table[T]) Print(translator core.TranslateFunc) map[string]any {
	// Use provided translator or fall back to table's translator
	trans := translator
	if trans == nil {
		trans = t.translator
	}

	// Build base component structure
	result := map[string]any{
		"type": "table",
	}

	// Add display class if set
	if t.options.Display != nil {
		result["display"] = *t.options.Display
	}

	// Build data section
	dataSection := make(map[string]any)

	// Add filter flag
	if t.hasFilter != nil {
		dataSection["hasFilter"] = *t.hasFilter
	} else {
		dataSection["hasFilter"] = t.filter != nil
	}

	// Add fields
	dataSection["fields"] = t.exportFields(trans)

	// Add options
	dataSection["options"] = t.exportOptions(trans)

	// Determine mode: AJAX (url set) vs Static (data set)
	if t.url != nil {
		// AJAX mode: URL for dynamic loading
		dataSection["url"] = t.url.PrintPrefix()
		dataSection["data"] = nil
		dataSection["components"] = nil
	} else {
		// Static mode: embedded data
		dataSection["url"] = nil
		dataSection["data"] = t.GetData(OutputWeb)
		dataSection["components"] = nil
	}

	result["data"] = dataSection

	// If filter exists, automatically wrap table in Query component
	if t.filter != nil {
		// Export filter form fields
		filterForm := t.filter.ExportForFrontend()

		// Get all fields from filter group to check Form property
		fields := t.filter.GetFields()
		extraData := make(map[string]any)
		visibleFilterForm := make([]map[string]any, 0)

		// Separate visible fields from hidden (form=false) fields
		for i, fieldExport := range filterForm {
			if i < len(fields) {
				field := fields[i]
				if !field.GetForm() {
					// Hidden field - add to extra data with default value
					extraData[field.GetID()] = field.GetDefault()
				} else {
					// Visible field - keep in filter form
					visibleFilterForm = append(visibleFilterForm, fieldExport)
				}
			} else {
				// Fallback: if field count mismatch, keep in form
				visibleFilterForm = append(visibleFilterForm, fieldExport)
			}
		}

		// Create Query component with visible fields only
		saveStateId := t.options.SaveStateId
		query := query.NewQuery(visibleFilterForm, saveStateId, t.options.Display)

		// Set extra data if any hidden fields exist
		if len(extraData) > 0 {
			query.SetExtraData(extraData)
		}

		// Add table as nested component in query
		query.AddArray(result)

		return query.Print(trans)
	}

	return result
}

// ExportFields is the public version of exportFields that allows external packages
// (like dialog) to access field definitions for building dialog tables.
func (t *Table[T]) ExportFields() []map[string]any {
	return t.exportFields(t.translator)
}

// exportOptions converts TableOptions to JSON map for component output.
// Matches component.Table options format exactly.
func (t *Table[T]) exportOptions(translator core.TranslateFunc) map[string]any {
	opts := t.options
	options := make(map[string]any)

	// Add all options that are set
	if opts.Class != nil {
		options["class"] = *opts.Class
	}
	if opts.Title != nil {
		options["title"] = *opts.Title
	}
	if opts.TextNoData != nil {
		options["textNoData"] = *opts.TextNoData
	}
	if opts.EmptyState != nil {
		options["emptyState"] = opts.EmptyState.PrintData(translator)
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
	if opts.Csv != nil && *opts.Csv && t.url != nil {
		csvBtn := button.NewTableButton(
			core.ButtonActionDownload,
			"csv",
			t.url,
			"CSV",
			core.ColorAccent,
			false,
			map[string]any{"data": map[string]bool{"_csv": true}},
		)
		topButtons = append(topButtons, csvBtn)
	}

	// Auto-generate Excel download button if Excel option enabled and URL exists
	if opts.Excel != nil && *opts.Excel && t.url != nil {
		excelBtn := button.NewTableButton(
			core.ButtonActionDownload,
			"explicit",
			t.url,
			"Excel",
			core.ColorAccent,
			false,
			map[string]any{"data": map[string]bool{"_excel": true}},
		)
		topButtons = append(topButtons, excelBtn)
	}

	// Export ButtonsTop in same format as SelectButtons
	if len(topButtons) > 0 {
		buttons := make([]map[string]any, len(topButtons))
		for i, btn := range topButtons {
			buttons[i] = btn.Print(translator)
		}
		options["buttons"] = map[string]any{"buttons": buttons}
	}

	if opts.Reload != nil {
		options["reload"] = *opts.Reload
	}
	if opts.Dense != nil {
		options["dense"] = *opts.Dense
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
			buttons[i] = btn.Print(translator)
		}
		options["selectButtons"] = buttons
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

	return options
}
