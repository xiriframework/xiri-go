package dialog

import (
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
)

// ExtractMultiSelectRequest parses HTTP requests for multi-select dialog workflows.
//
// Dialog Workflow:
//  1. User selects multiple table rows → POST {"data": [id1, id2, ...]}
//  2. Server shows dialog with NewDialogFormMultiEdit/Delete → extra["data"] + extra["done"]=true
//  3. User submits dialog → POST {"data": [id1, id2, ...], "done": true, "field": value, ...}
//
// Returns:
//   - selectedIDs: []int64 - The selected row IDs
//   - formData: map[string]interface{} - The complete request data (for form binding)
//   - isDialogOpen: bool - true if step 1 (showing dialog), false if step 3 (form submission)
//   - error: any parsing or validation error
func ExtractMultiSelectRequest(c echo.Context) ([]int64, map[string]interface{}, bool, error) {
	var requestData map[string]interface{}
	if err := c.Bind(&requestData); err != nil {
		return nil, nil, false, fmt.Errorf("invalid request data: %w", err)
	}

	ids, err := extractIDsFromData(requestData)
	if err != nil {
		return nil, nil, false, err
	}

	_, hasDone := requestData["done"]
	isDialogOpen := !hasDone

	return ids, requestData, isDialogOpen, nil
}

// extractIDsFromData is an internal helper that extracts IDs from the "data" field.
func extractIDsFromData(requestData map[string]interface{}) ([]int64, error) {
	data, ok := requestData["data"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("no data field found in request")
	}

	ids := make([]int64, 0, len(data))
	for _, id := range data {
		if idStr, ok := id.(string); ok {
			parsedID, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				continue
			}
			ids = append(ids, parsedID)
			continue
		}

		if idFloat, ok := id.(float64); ok {
			ids = append(ids, int64(idFloat))
			continue
		}
	}

	if len(ids) == 0 {
		return nil, fmt.Errorf("no items selected")
	}

	return ids, nil
}
