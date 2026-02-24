// Package query provides the Query filter/search component that wraps other components with filter forms.
package query

import (
	"github.com/xiriframework/xiri-go/component/button"
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
	"github.com/xiriframework/xiri-go/form/group"
)

// Query represents a query/filter form component with dynamic content
type Query struct {
	filterForm  []map[string]any
	url         *url.Url
	buttonLine  *button.ButtonLine
	saveStateId *string
	extra       map[string]any
	display     *string
	data        []map[string]any
}

// NewQuery creates a new query component
func NewQuery(filterForm []map[string]any, saveStateId *string, display *string) *Query {
	return &Query{
		filterForm:  filterForm,
		url:         nil,
		buttonLine:  nil,
		saveStateId: saveStateId,
		extra:       nil,
		display:     display,
		data:        make([]map[string]any, 0),
	}
}

// NewQueryWithFields creates a query with pre-exported field definitions
func NewQueryWithFields(
	filterForm []map[string]any,
	u *url.Url,
	buttonLine *button.ButtonLine,
	saveStateId *string,
	extra map[string]any,
) *Query {
	return &Query{
		filterForm:  filterForm,
		url:         u,
		buttonLine:  buttonLine,
		saveStateId: saveStateId,
		extra:       extra,
		display:     nil,
		data:        make([]map[string]any, 0),
	}
}

// NewQueryWithFormGroup creates a query from a FormGroup
func NewQueryWithFormGroup(
	fg *group.FormGroup,
	fieldValues map[string]any,
	u *url.Url,
	buttonLine *button.ButtonLine,
	saveStateId *string,
	extra map[string]any,
) *Query {
	var filterForm []map[string]any
	if fieldValues != nil {
		filterForm = fg.ExportForFrontendWithValues(fieldValues)
	} else {
		filterForm = fg.ExportForFrontend()
	}

	return &Query{
		filterForm:  filterForm,
		url:         u,
		buttonLine:  buttonLine,
		saveStateId: saveStateId,
		extra:       extra,
		display:     nil,
		data:        make([]map[string]any, 0),
	}
}

// Add adds a dynamic component to the query
func (q *Query) Add(component core.Component, translator core.TranslateFunc) *Query {
	q.data = append(q.data, component.Print(translator))
	return q
}

// AddArray adds a raw component map to the query
func (q *Query) AddArray(component map[string]any) *Query {
	q.data = append(q.data, component)
	return q
}

// SetExtraData sets additional data for the query
func (q *Query) SetExtraData(extra map[string]any) *Query {
	q.extra = extra
	return q
}

// WithSaveStateId sets the save state ID (optional)
func (q *Query) WithSaveStateId(saveStateId string) *Query {
	q.saveStateId = &saveStateId
	return q
}

// WithDisplay sets the display/layout class (optional)
func (q *Query) WithDisplay(display string) *Query {
	q.display = &display
	return q
}

// Print returns the JSON representation of the query
func (q *Query) Print(translator core.TranslateFunc) map[string]any {
	data := map[string]any{
		"fields":      q.filterForm,
		"dyn":         q.data,
		"url":         nil,
		"buttonline":  nil,
		"saveState":   q.saveStateId != nil,
		"saveStateId": q.saveStateId,
		"extra":       q.extra,
	}

	if q.url != nil {
		data["url"] = q.url.PrintPrefix()
	}

	if q.buttonLine != nil {
		data["buttonline"] = q.buttonLine.Print(translator)
	}

	return map[string]any{
		"type":    "query",
		"display": q.display,
		"data":    data,
	}
}
