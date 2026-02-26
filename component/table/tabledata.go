package table

import (
	"bytes"
	"encoding/csv"
	"fmt"

	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/response"
	"github.com/xuri/excelize/v2"
)

// TableDataResponse represents table data returned by AJAX table endpoints.
// Provides type-safe API for building table responses with optional fields, footer, and components.
//
// This is the new location for TableDataResponse. The old ui/component/table_data.go is deprecated.
//
// JSON format: {"data": [...], "fields": [...], "footer": {...}, "components": [...]}
// CSV format: {"csv": "field1;field2\nval1;val2\n"}
// Excel format: {"excel": <binary bytes>}
type TableDataResponse struct {
	data          []map[string]any // Required - table rows
	fields        []map[string]any // Optional - field definitions if changed (for JSON output)
	fieldsForCSV  []map[string]any // Internal - field definitions for CSV/Excel headers (not exported to JSON)
	footer        map[string]any   // Optional - footer aggregations (sum, count)
	components    []core.Component // Optional - additional UI components
	outputType    OutputType       // Output type (Web, CSV, PDF, Excel)
	includeFields bool             // Whether to include fields in JSON output
	excelData     []byte           // Excel binary data (only populated for OutputExcel)
	totalCount    *int             // Optional - total record count for server-side pagination
}

// NewTableDataResponse creates a new TableDataResponse with the given row data and output type.
// This is the minimum required data for a table response.
//
// Usage:
//
//	td := table.NewTableDataResponse(rows, table.OutputWeb)
//	return wc.TableData(td)
func NewTableDataResponse(data []map[string]any, outputType OutputType) *TableDataResponse {
	// Ensure data is never nil - frontend expects array, not null
	if data == nil {
		data = make([]map[string]any, 0)
	}
	return &TableDataResponse{
		data:       data,
		fields:     nil,
		footer:     nil,
		components: make([]core.Component, 0),
		outputType: outputType,
	}
}

// WithFields sets the field definitions for the table.
// Use this when field definitions have changed and need to be sent to the frontend.
// This will include fields in the JSON output.
//
// Usage:
//
//	fields := []map[string]any{
//	    {"id": "name", "type": "text", "name": "Name"},
//	}
//	td.WithFields(fields)
func (td *TableDataResponse) WithFields(fields []map[string]any) *TableDataResponse {
	td.fields = fields
	td.includeFields = true
	return td
}

// withFieldsForCSV sets field definitions internally for CSV header generation only.
// These fields will NOT be included in JSON output.
// This is an internal method used by Table[T].ToTableDataResponse().
func (td *TableDataResponse) withFieldsForCSV(fields []map[string]any) *TableDataResponse {
	td.fieldsForCSV = fields
	return td
}

// WithFooter sets the footer aggregation data.
// Footer data typically contains sums or counts for numeric columns.
//
// Usage:
//
//	footer := map[string]any{
//	    "total": 1234.56,
//	    "count": 42,
//	}
//	td.WithFooter(footer)
func (td *TableDataResponse) WithFooter(footer map[string]any) *TableDataResponse {
	td.footer = footer
	return td
}

// WithTotalCount sets the total record count for server-side pagination.
// This is required when serverSide: true is enabled on the table.
// The totalCount is the total number of records matching the filter,
// not the number of records in the current page.
//
// Usage:
//
//	totalCount := dbm.Device.CountWithFilters(filters)
//	td.WithTotalCount(totalCount)
func (td *TableDataResponse) WithTotalCount(count int) *TableDataResponse {
	td.totalCount = &count
	return td
}

// AddComponent adds a UI component to be displayed alongside the table.
// Components are stored as Component objects and only rendered (Print()) when
// the final response is built. This allows translation to be applied correctly.
//
// Common use case: Adding MultiProgress components for statistics in admin tables.
//
// Usage:
//
//	td.AddComponent(component.NewMultiProgress("Statistics", 5, true, nil))
//	  .AddComponent(component.NewMultiProgress("More Stats", 5, true, nil))
func (td *TableDataResponse) AddComponent(comp core.Component) *TableDataResponse {
	if comp != nil {
		td.components = append(td.components, comp)
	}
	return td
}

// Print converts the TableDataResponse to the appropriate output format.
// Components are rendered using the provided translator function.
//
// Output format for OutputWeb (JSON):
//
//	{
//	  "data": [...],              // Always present
//	  "fields": [...],            // Only if WithFields() was called
//	  "footer": {...},            // Only if WithFooter() was called
//	  "components": [...]         // Only if AddComponent() was called
//	}
//
// Output format for OutputCSV:
//
//	{
//	  "csv": "field1;field2\nval1;val2\n"
//	}
//
// Output format for OutputExcel:
//
//	{
//	  "excel": <binary bytes>
//	}
func (td *TableDataResponse) Print(translator core.TranslateFunc) map[string]any {
	// Handle CSV output
	if td.outputType == OutputCSV {
		csvString := td.generateCSV(translator)
		return map[string]any{
			"csv": csvString,
		}
	}

	// Handle Excel output
	if td.outputType == OutputExcel {
		excelBytes, err := td.generateExcel(translator)
		if err != nil {
			// On error, return empty Excel data
			return map[string]any{
				"excel": []byte{},
			}
		}
		return map[string]any{
			"excel": excelBytes,
		}
	}

	// Handle regular JSON output (Web, PDF)
	response := map[string]any{
		"data": td.data,
	}

	// Add totalCount for server-side pagination
	if td.totalCount != nil {
		response["totalCount"] = *td.totalCount
	}

	// Add fields if set AND includeFields flag is true (field definitions changed)
	if td.includeFields && td.fields != nil && len(td.fields) > 0 {
		response["fields"] = td.fields
	}

	// Add footer if set
	if td.footer != nil && len(td.footer) > 0 {
		response["footer"] = td.footer
	}

	// Render components if any were added
	if len(td.components) > 0 {
		components := make([]map[string]any, 0, len(td.components))
		for _, comp := range td.components {
			printed := comp.Print(translator)
			if printed != nil && len(printed) > 0 {
				components = append(components, printed)
			}
		}
		// Only add components array if there are actually rendered components
		if len(components) > 0 {
			response["components"] = components
		}
	}

	return response
}

// DataResponse returns a DataResult with the appropriate response type (JSON, CSV, or Excel).
// Table responses are NOT wrapped in {"data": ...} — they have their own top-level structure.
func (td *TableDataResponse) DataResponse(translator core.TranslateFunc) response.DataResult {
	printed := td.Print(translator)

	if csv, ok := printed["csv"].(string); ok {
		return response.NewCSVDataResult(csv)
	}
	if excel, ok := printed["excel"].([]byte); ok {
		return response.NewExcelDataResult(excel)
	}
	return response.DataResult{Type: response.ResponseJSON, Body: printed}
}

// expandNFieldColumns expands []string values from N-field formatters into separate columns.
// For each field where any row contains a []string value, it determines the maximum slice length
// and creates that many individual columns. The first column keeps the original field ID/name,
// subsequent columns get "_2", "_3" suffixes (e.g., "distance", "distance_2", "distance_3").
func (td *TableDataResponse) expandNFieldColumns() {
	fieldsToUse := td.fieldsForCSV
	if fieldsToUse == nil {
		fieldsToUse = td.fields
	}
	if len(fieldsToUse) == 0 || len(td.data) == 0 {
		return
	}

	// Pass 1: Find maximum N for each field that has []string values
	maxN := make(map[string]int)
	for _, row := range td.data {
		for key, val := range row {
			if strs, ok := val.([]string); ok {
				if len(strs) > maxN[key] {
					maxN[key] = len(strs)
				}
			}
		}
	}

	if len(maxN) == 0 {
		return
	}

	// Pass 2: Expand field definitions
	expandedFields := make([]map[string]any, 0, len(fieldsToUse))
	for _, fieldDef := range fieldsToUse {
		fieldID, hasID := fieldDef["id"].(string)
		fieldName, _ := fieldDef["name"].(string)
		n, isNField := maxN[fieldID]
		if !hasID || !isNField || n <= 1 {
			expandedFields = append(expandedFields, fieldDef)
			continue
		}
		// Expand into N columns
		for i := range n {
			newDef := make(map[string]any, len(fieldDef))
			for k, v := range fieldDef {
				newDef[k] = v
			}
			if i == 0 {
				// First column keeps original ID and name
			} else {
				newDef["id"] = fmt.Sprintf("%s_%d", fieldID, i+1)
				newDef["name"] = fmt.Sprintf("%s %d", fieldName, i+1)
			}
			expandedFields = append(expandedFields, newDef)
		}
	}

	// Pass 3: Expand data rows
	for _, row := range td.data {
		for fieldID, n := range maxN {
			val, exists := row[fieldID]
			if !exists {
				continue
			}
			strs, ok := val.([]string)
			if !ok {
				continue
			}
			// Set first value
			if len(strs) > 0 {
				row[fieldID] = strs[0]
			} else {
				row[fieldID] = ""
			}
			// Set subsequent values
			for i := 1; i < n; i++ {
				expandedID := fmt.Sprintf("%s_%d", fieldID, i+1)
				if i < len(strs) {
					row[expandedID] = strs[i]
				} else {
					row[expandedID] = ""
				}
			}
		}
	}

	// Update field definitions
	if td.fieldsForCSV != nil {
		td.fieldsForCSV = expandedFields
	} else {
		td.fields = expandedFields
	}
}

// generateCSV creates a CSV string from the table data.
// Uses semicolon (;) delimiter for Excel compatibility.
// Only includes fields that are marked as CSV-enabled.
func (td *TableDataResponse) generateCSV(translator core.TranslateFunc) string {
	td.expandNFieldColumns()

	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)
	writer.Comma = ';' // Semicolon separator for Excel compatibility (European locale)

	// If no data, return empty CSV
	if len(td.data) == 0 {
		return ""
	}

	// Build field ID to name mapping from field definitions
	fieldIDToName := make(map[string]string)
	fieldIDsOrdered := make([]string, 0)

	// Use fieldsForCSV first (internal), then fall back to fields (public)
	fieldsToUse := td.fieldsForCSV
	if fieldsToUse == nil {
		fieldsToUse = td.fields
	}

	if fieldsToUse != nil && len(fieldsToUse) > 0 {
		// Use field definitions to get proper order and names
		for _, fieldDef := range fieldsToUse {
			fieldID, hasID := fieldDef["id"].(string)
			fieldName, hasName := fieldDef["name"].(string)
			if hasID && hasName {
				fieldIDToName[fieldID] = fieldName
				fieldIDsOrdered = append(fieldIDsOrdered, fieldID)
			}
		}
	}

	// Fallback: if no field definitions, use keys from first row
	if len(fieldIDsOrdered) == 0 {
		firstRow := td.data[0]
		for fieldID := range firstRow {
			fieldIDsOrdered = append(fieldIDsOrdered, fieldID)
			fieldIDToName[fieldID] = fieldID // Use ID as name
		}
	}

	// Write header row with field names (translated)
	headerRow := make([]string, len(fieldIDsOrdered))
	for i, fieldID := range fieldIDsOrdered {
		if name, ok := fieldIDToName[fieldID]; ok {
			headerRow[i] = name
		} else {
			headerRow[i] = fieldID
		}
	}
	if err := writer.Write(headerRow); err != nil {
		return fmt.Sprintf("Error writing CSV header: %v", err)
	}

	// Write data rows
	for _, rowData := range td.data {
		row := make([]string, len(fieldIDsOrdered))
		for i, fieldID := range fieldIDsOrdered {
			// Convert value to string
			value := rowData[fieldID]
			if value == nil {
				row[i] = ""
				continue
			}

			// Handle array values (e.g., [display, value] from formatters)
			if arr, ok := value.([]interface{}); ok && len(arr) > 0 {
				// For CSV, use the display value (first element)
				row[i] = fmt.Sprintf("%v", arr[0])
			} else {
				row[i] = fmt.Sprintf("%v", value)
			}
		}

		if err := writer.Write(row); err != nil {
			return fmt.Sprintf("Error writing CSV row: %v", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Sprintf("Error flushing CSV writer: %v", err)
	}

	return buf.String()
}

// generateExcel creates an Excel (.xlsx) file from the table data.
// Uses excelize library to create a properly formatted Excel workbook.
// Only includes fields that are marked as CSV-enabled (same logic as CSV).
func (td *TableDataResponse) generateExcel(translator core.TranslateFunc) ([]byte, error) {
	td.expandNFieldColumns()

	// If no data, return empty Excel file
	if len(td.data) == 0 {
		f := excelize.NewFile()
		defer f.Close()
		buf, err := f.WriteToBuffer()
		if err != nil {
			return nil, fmt.Errorf("error creating empty Excel file: %w", err)
		}
		return buf.Bytes(), nil
	}

	// Create new Excel file
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Sheet1"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, fmt.Errorf("error creating Excel sheet: %w", err)
	}
	f.SetActiveSheet(index)

	// Build field ID to name mapping from field definitions
	fieldIDToName := make(map[string]string)
	fieldIDsOrdered := make([]string, 0)

	// Use fieldsForCSV first (internal), then fall back to fields (public)
	fieldsToUse := td.fieldsForCSV
	if fieldsToUse == nil {
		fieldsToUse = td.fields
	}

	if fieldsToUse != nil && len(fieldsToUse) > 0 {
		// Use field definitions to get proper order and names
		for _, fieldDef := range fieldsToUse {
			fieldID, hasID := fieldDef["id"].(string)
			fieldName, hasName := fieldDef["name"].(string)
			if hasID && hasName {
				fieldIDToName[fieldID] = fieldName
				fieldIDsOrdered = append(fieldIDsOrdered, fieldID)
			}
		}
	}

	// Fallback: if no field definitions, use keys from first row
	if len(fieldIDsOrdered) == 0 {
		firstRow := td.data[0]
		for fieldID := range firstRow {
			fieldIDsOrdered = append(fieldIDsOrdered, fieldID)
			fieldIDToName[fieldID] = fieldID // Use ID as name
		}
	}

	// Write header row with field names
	for colIdx, fieldID := range fieldIDsOrdered {
		cellName, err := excelize.CoordinatesToCellName(colIdx+1, 1)
		if err != nil {
			return nil, fmt.Errorf("error getting cell name for header: %w", err)
		}
		name := fieldIDToName[fieldID]
		if name == "" {
			name = fieldID
		}
		if err := f.SetCellValue(sheetName, cellName, name); err != nil {
			return nil, fmt.Errorf("error writing header cell: %w", err)
		}
	}

	// Write data rows
	for rowIdx, rowData := range td.data {
		excelRow := rowIdx + 2 // Excel rows are 1-indexed, +1 for header
		for colIdx, fieldID := range fieldIDsOrdered {
			cellName, err := excelize.CoordinatesToCellName(colIdx+1, excelRow)
			if err != nil {
				return nil, fmt.Errorf("error getting cell name for data: %w", err)
			}

			// Get value
			value := rowData[fieldID]
			if value == nil {
				continue // Leave cell empty
			}

			// Handle array values (e.g., [display, value] from formatters)
			if arr, ok := value.([]interface{}); ok && len(arr) > 0 {
				// For Excel, use the display value (first element)
				value = arr[0]
			}

			// Write cell value
			if err := f.SetCellValue(sheetName, cellName, value); err != nil {
				return nil, fmt.Errorf("error writing data cell: %w", err)
			}
		}
	}

	// Auto-size columns based on content
	for colIdx, fieldID := range fieldIDsOrdered {
		maxWidth := float64(10) // Minimum width

		// Check header width
		headerName := fieldIDToName[fieldID]
		if headerName == "" {
			headerName = fieldID
		}
		headerWidth := float64(len(headerName)) * 1.2
		if headerWidth > maxWidth {
			maxWidth = headerWidth
		}

		// Check data values in this column
		for _, rowData := range td.data {
			value := rowData[fieldID]
			if value == nil {
				continue
			}

			// Handle array values
			if arr, ok := value.([]interface{}); ok && len(arr) > 0 {
				value = arr[0]
			}

			// Convert to string and measure
			valueStr := fmt.Sprintf("%v", value)
			valueWidth := float64(len(valueStr)) * 1.2
			if valueWidth > maxWidth {
				maxWidth = valueWidth
			}
		}

		// Cap at maximum width
		if maxWidth > 50 {
			maxWidth = 50
		}

		// Set column width (columns are 1-indexed)
		colName, err := excelize.ColumnNumberToName(colIdx + 1)
		if err != nil {
			continue // Skip on error, don't fail entire export
		}
		if err := f.SetColWidth(sheetName, colName, colName, maxWidth); err != nil {
			continue // Skip on error, don't fail entire export
		}
	}

	// Write to buffer
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("error writing Excel to buffer: %w", err)
	}

	return buf.Bytes(), nil
}
