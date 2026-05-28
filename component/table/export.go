package table

import (
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/response"
)

// ToServerSideResponse creates an AJAX data response for server-side pagination.
// This is a convenience method that wraps ToTableDataResponse() and adds totalCount.
//
// The totalCount parameter is required - use the total count from your database query
// (the count BEFORE pagination is applied).
//
// Example:
//
//	// Query with pagination
//	devices, totalCount := dbm.Device.FindWithPagination(filters, page, pageSize)
//	tbl.SetData(devices)
//	return wc.TableData(tbl.ToServerSideResponse(ctx, totalCount))
func (t *Table[T]) ToServerSideResponse(ctx *core.UiContext, totalCount int) *TableDataResponse {
	resp := t.ToTableDataResponse(ctx)
	resp.WithTotalCount(totalCount)
	return resp
}

// DataResponse returns a DataResult by delegating to ToTableDataResponse().DataResponse().
func (t *Table[T]) DataResponse(ctx *core.UiContext) response.DataResult {
	return t.ToTableDataResponse(ctx).DataResponse(ctx)
}

// ToTableDataResponse creates an AJAX data response with exact JSON structure.
//
// This produces the TableDataResponse for frontend compatibility.
// The response includes:
// - data: Array of formatted row objects
// - footer: Optional footer aggregations
//
// This method uses the table's internal outputType (set via SetOutputType).
// For CSV export, set outputType to OutputCSV before calling this method.
//
// This method is specifically for AJAX endpoint handlers that return table data
// without the full component definition. Use Print() to get the full component JSON.
func (t *Table[T]) ToTableDataResponse(ctx *core.UiContext) *TableDataResponse {

	// Get formatted data with all formatters applied using internal outputType
	data := t.GetData(ctx, t.outputType)

	// Create response with outputType
	resp := NewTableDataResponse(data, t.outputType)

	// For CSV output, filter out fields with csv=false
	// Uses non-generic exportFieldsForCSV from tableCore via promotion
	if t.outputType == OutputCSV {
		fields := t.exportFieldsForCSV(ctx)
		resp.withFieldsForCSV(fields)
	} else if t.outputType == OutputExcel {
		fields := t.exportFieldsForCSV(ctx)
		resp.withFieldsForCSV(fields)
	} else if t.fieldsCanChange {
		fields := t.exportFields(ctx)
		resp.WithFields(fields)
	}

	// Add components (only for Web output, excluded from CSV/PDF/Excel)
	// Components like MultiProgress, Charts, Info messages are rendered alongside the table
	if t.outputType != OutputCSV {
		// Calculate and add footer if any fields have aggregations
		footer := t.CalculateFooter(ctx, t.outputType)
		if len(footer) > 0 {
			resp.WithFooter(footer)
		}

		for _, comp := range t.components {
			resp.AddComponent(comp)
		}

		// Propagate auto-refresh poll interval (Web output only)
		if t.poll > 0 {
			resp.WithPoll(t.poll)
		}
	}

	return resp
}
