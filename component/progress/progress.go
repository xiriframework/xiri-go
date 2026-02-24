// Package progress provides progress indicator components for the Angular frontend.
package progress

import (
	"fmt"
	"sort"

	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
	"github.com/xiriframework/xiri-go/response"
)

// MultiProgressLine represents a single progress line
type MultiProgressLine struct {
	name  string
	cnt   int
	value string
	color core.Color
	calc  bool
}

// MultiProgress represents a multi-line progress indicator component
type MultiProgress struct {
	header  string
	show    int
	sort    bool
	display *string
	data    []MultiProgressLine
	sum     int
	url     *url.Url
	reload  *bool
}

// NewMultiProgress creates a new multi-progress component
func NewMultiProgress(header string, show int, sortEnabled bool, display *string) *MultiProgress {
	return &MultiProgress{
		header:  header,
		show:    show,
		sort:    sortEnabled,
		display: display,
		data:    make([]MultiProgressLine, 0),
		sum:     0,
	}
}

// AddLine adds a progress line that contributes to the total
func (mp *MultiProgress) AddLine(name string, cnt int, color core.Color, value *string) *MultiProgress {
	val := value
	if val == nil {
		// Convert integer to string - this will be displayed as text
		valStr := fmt.Sprintf("%d", cnt)
		val = &valStr
	}

	mp.data = append(mp.data, MultiProgressLine{
		name:  name,
		cnt:   cnt,
		value: *val,
		color: color,
		calc:  true,
	})
	mp.sum += cnt
	return mp
}

// AddTotal adds a total/summary line that doesn't contribute to percentage calc
func (mp *MultiProgress) AddTotal(name string, sum *int, color core.Color, value *string) *MultiProgress {
	totalSum := mp.sum
	if sum != nil {
		totalSum = *sum
	}

	val := value
	if val == nil {
		// Convert integer to string - this will be displayed as text
		valStr := fmt.Sprintf("%d", totalSum)
		val = &valStr
	}

	mp.data = append(mp.data, MultiProgressLine{
		name:  name,
		cnt:   totalSum,
		value: *val,
		color: color,
		calc:  false,
	})
	return mp
}

// SetSum manually sets the sum for percentage calculations
func (mp *MultiProgress) SetSum(sum int) *MultiProgress {
	mp.sum = sum
	return mp
}

// WithDisplay sets the display/layout class (optional)
// Returns the MultiProgress for method chaining
func (mp *MultiProgress) WithDisplay(display string) *MultiProgress {
	mp.display = &display
	return mp
}

// SetURL sets the AJAX data URL. When set, data lines are cleared and the frontend loads data dynamically.
func (mp *MultiProgress) SetURL(url *url.Url) *MultiProgress {
	mp.url = url
	mp.data = nil
	return mp
}

// WithReload enables periodic reload of the progress data when using AJAX mode.
func (mp *MultiProgress) WithReload(reload bool) *MultiProgress {
	mp.reload = &reload
	return mp
}

// Print returns the JSON representation of the multi-progress
func (mp *MultiProgress) Print(translator core.TranslateFunc) map[string]any {
	if mp.url != nil {
		return map[string]any{
			"type":    "multiprogress",
			"display": mp.display,
			"data": map[string]any{
				"header": mp.header,
				"url":    mp.url.PrintPrefix(),
				"reload": mp.reload,
			},
		}
	}

	data := mp.printData(translator)
	if data == nil {
		return map[string]any{}
	}

	return map[string]any{
		"type":    "multiprogress",
		"display": mp.display,
		"data":    data,
	}
}

// PrintData returns only the data portion of the multi-progress (for use in data endpoints).
// Returns nil when there is no data.
func (mp *MultiProgress) PrintData(translator core.TranslateFunc) map[string]any {
	return mp.printData(translator)
}

// DataResponse returns a DataResult wrapping the progress data in {"data": ...} envelope.
func (mp *MultiProgress) DataResponse(translator core.TranslateFunc) response.DataResult {
	return response.NewJSONDataResult(mp.PrintData(translator))
}

// printData builds the data map used by both Print and PrintData.
// Returns nil when there is no data.
func (mp *MultiProgress) printData(translator core.TranslateFunc) map[string]any {
	if len(mp.data) == 0 {
		return nil
	}

	// Sort if enabled
	if mp.sort {
		sort.Slice(mp.data, func(i, j int) bool {
			if mp.data[i].cnt == mp.data[j].cnt {
				return mp.data[i].name > mp.data[j].name
			}
			return mp.data[i].cnt > mp.data[j].cnt
		})
	}

	// Build output list
	list := make([]map[string]any, 0, len(mp.data))
	for _, item := range mp.data {
		var cnt float64
		if mp.sum == 0 {
			cnt = 0
		} else {
			cnt = float64(item.cnt) / float64(mp.sum) * 100
		}

		var progress float64
		if item.calc {
			progress = cnt
		} else {
			progress = 100
		}

		list = append(list, map[string]any{
			"text":     item.name,
			"value":    item.value,
			"color":    string(item.color),
			"progress": progress,
		})
	}

	return map[string]any{
		"header": mp.header,
		"show":   mp.show,
		"data":   list,
	}
}
