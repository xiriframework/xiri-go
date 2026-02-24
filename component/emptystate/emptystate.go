// Package emptystate provides an empty state display component for the Angular frontend.
package emptystate

import (
	"github.com/xiriframework/xiri-go/component/button"
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/response"
)

// EmptyState represents an empty state display with icon, title, description, and optional button.
type EmptyState struct {
	icon        string
	iconColor   core.Color
	title       string
	description *string
	button      *button.Button
	display     *string
}

// New creates a new empty state component with required parameters.
func New(icon string, iconColor core.Color, title string) *EmptyState {
	return &EmptyState{
		icon:      icon,
		iconColor: iconColor,
		title:     title,
	}
}

// WithDescription sets the description text (optional).
func (es *EmptyState) WithDescription(description string) *EmptyState {
	es.description = &description
	return es
}

// WithButton sets the call-to-action button (optional).
func (es *EmptyState) WithButton(btn *button.Button) *EmptyState {
	es.button = btn
	return es
}

// WithDisplay sets the display/layout class (optional).
func (es *EmptyState) WithDisplay(display string) *EmptyState {
	es.display = &display
	return es
}

// Print returns the JSON representation of the empty state component.
func (es *EmptyState) Print(translator core.TranslateFunc) map[string]any {
	return map[string]any{
		"type":    "empty-state",
		"display": es.display,
		"data":    es.printData(translator),
	}
}

// PrintData returns only the data portion (for use in table options).
func (es *EmptyState) PrintData(translator core.TranslateFunc) map[string]any {
	return es.printData(translator)
}

// DataResponse returns a DataResult wrapping the empty state data in {"data": ...} envelope.
func (es *EmptyState) DataResponse(translator core.TranslateFunc) response.DataResult {
	return response.NewJSONDataResult(es.PrintData(translator))
}

// printData builds the data map used by both Print and PrintData.
func (es *EmptyState) printData(translator core.TranslateFunc) map[string]any {
	data := map[string]any{
		"icon":      es.icon,
		"iconColor": string(es.iconColor),
		"title":     core.Translate(translator, es.title),
	}

	if es.description != nil {
		data["description"] = core.Translate(translator, *es.description)
	}

	if es.button != nil {
		data["button"] = es.button.Print(translator)
	}

	return data
}
