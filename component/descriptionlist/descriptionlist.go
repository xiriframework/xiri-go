// Package descriptionlist provides a key-value pair display component for detail pages.
package descriptionlist

import (
	"github.com/xiriframework/xiri-go/component/core"
)

// Item represents a single key-value entry in the description list.
type Item struct {
	label string
	value string
	icon  *string
	color *core.Color
	typ   *string // "text", "link", "html", "badge"
	list  *DescriptionList
}

// Icon sets the item icon.
func (i *Item) Icon(icon string) *Item {
	i.icon = &icon
	return i
}

// Color sets the item color.
func (i *Item) Color(c core.Color) *Item {
	i.color = &c
	return i
}

// Type sets the item display type ("text", "link", "html", "badge").
func (i *Item) Type(t string) *Item {
	i.typ = &t
	return i
}

// Done returns the parent DescriptionList for chaining.
func (i *Item) Done() *DescriptionList {
	return i.list
}

// DescriptionList represents a collection of key-value pairs for detail views.
type DescriptionList struct {
	items   []Item
	columns *int
	layout  *string
	display *string
}

// New creates a new DescriptionList.
func New() *DescriptionList {
	return &DescriptionList{
		items: make([]Item, 0),
	}
}

// Add adds a new item with label and value, returning the Item for further configuration.
func (d *DescriptionList) Add(label, value string) *Item {
	item := Item{label: label, value: value, list: d}
	d.items = append(d.items, item)
	return &d.items[len(d.items)-1]
}

// Columns sets the number of columns (1, 2, or 3).
func (d *DescriptionList) Columns(c int) *DescriptionList {
	d.columns = &c
	return d
}

// Layout sets the layout mode ("horizontal" or "stacked").
func (d *DescriptionList) Layout(l string) *DescriptionList {
	d.layout = &l
	return d
}

// WithDisplay sets the display/layout class.
func (d *DescriptionList) WithDisplay(display string) *DescriptionList {
	d.display = &display
	return d
}

// Print returns the JSON representation of the description list component.
func (d *DescriptionList) Print(translator core.TranslateFunc) map[string]any {
	itemsData := make([]map[string]any, len(d.items))
	for idx, item := range d.items {
		entry := map[string]any{
			"label": core.Translate(translator, item.label),
			"value": item.value,
		}
		if item.icon != nil {
			entry["icon"] = *item.icon
		}
		if item.color != nil {
			entry["color"] = string(*item.color)
		}
		if item.typ != nil {
			entry["type"] = *item.typ
		}
		itemsData[idx] = entry
	}

	data := map[string]any{
		"items": itemsData,
	}

	if d.columns != nil {
		data["columns"] = *d.columns
	}
	if d.layout != nil {
		data["layout"] = *d.layout
	}

	result := map[string]any{
		"type": "description-list",
		"data": data,
	}

	if d.display != nil {
		result["display"] = *d.display
	}

	return result
}
