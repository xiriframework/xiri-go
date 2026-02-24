package timeline

import (
	"github.com/xiriframework/xiri-go/component/core"
)

// Item represents a single timeline entry.
type Item struct {
	title       string
	description string
	datetime    string
	icon        string
	iconColor   string
}

// Timeline represents a timeline component for displaying chronological events.
type Timeline struct {
	items   []Item
	display string
}

// New creates a new Timeline component.
func New() *Timeline {
	return &Timeline{
		items: make([]Item, 0),
	}
}

// Add adds a timeline item with title.
func (t *Timeline) Add(title string) *Item {
	item := Item{title: title}
	t.items = append(t.items, item)
	return &t.items[len(t.items)-1]
}

// WithDisplay sets the display/layout class.
func (t *Timeline) WithDisplay(display string) *Timeline {
	t.display = display
	return t
}

// Description sets the description on the item.
func (i *Item) Description(description string) *Item {
	i.description = description
	return i
}

// Datetime sets the datetime string on the item.
func (i *Item) Datetime(datetime string) *Item {
	i.datetime = datetime
	return i
}

// Icon sets the Material icon on the item.
func (i *Item) Icon(icon string) *Item {
	i.icon = icon
	return i
}

// IconColor sets the icon color on the item.
func (i *Item) IconColor(color string) *Item {
	i.iconColor = color
	return i
}

// Print returns the JSON representation of the timeline component.
func (t *Timeline) Print(translator core.TranslateFunc) map[string]any {
	items := make([]map[string]any, len(t.items))
	for idx, item := range t.items {
		entry := map[string]any{
			"title": core.Translate(translator, item.title),
		}
		if item.description != "" {
			entry["description"] = core.Translate(translator, item.description)
		}
		if item.datetime != "" {
			entry["datetime"] = item.datetime
		}
		if item.icon != "" {
			entry["icon"] = item.icon
		}
		if item.iconColor != "" {
			entry["iconColor"] = item.iconColor
		}
		items[idx] = entry
	}

	result := map[string]any{
		"type":  "timeline",
		"items": items,
	}

	if t.display != "" {
		result["display"] = t.display
	}

	return result
}
