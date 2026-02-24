package card

import (
	"strconv"

	"github.com/xiriframework/xiri-go/component/core"
)

// CardContent represents the content structure for a card
type CardContent struct {
	Fields any  `json:"fields"`
	Data   any  `json:"data"` // Usually []map[string]interface{}
	Dense  bool `json:"dense"`
}

// SetDense sets the dense display mode for the card list
// Returns the CardContent for method chaining
func (c *CardContent) SetDense(dense bool) *CardContent {
	c.Dense = dense
	return c
}

// NewCardLinks creates a card with link content.
// Use NewCardLinkContent() to build the content parameter.
// Example:
//
//	lines := []CardLinkContentLine{{Text: "Dashboard", Link: "/dashboard", Icon: "home", Color: core.ColorPrimary}}
//	card := NewCardLinks("Navigation", NewCardLinkContent(lines))
func NewCardLinks(
	header string,
	content *CardLinkContent,
) *Card {
	return NewCard(
		core.CardTypeTable,
		content,
		header,
		nil,   // headerSub
		nil,   // headerIcon
		nil,   // headerIconColor
		false, // translate
		false, // forceMinWidth
		nil,   // display
	)
}

// CardLinkContent represents the content structure for a card link list
type CardLinkContent struct {
	CardContent
}

// CardLinkContentLine represents a single line in a card link list
type CardLinkContentLine struct {
	Text  string     `json:"text"`
	Link  string     `json:"textLink"`
	Icon  string     `json:"icon"`
	Color core.Color `json:"iconColor"`
}

// NewCardLinkContent creates a new CardLinkContent with default id/text fields
// Use SetDense() to enable dense mode
func NewCardLinkContent(data []CardLinkContentLine) *CardLinkContent {

	fields := []map[string]any{
		{"id": "text", "name": "Link", "format": "linkrow"},
	}

	return &CardLinkContent{
		CardContent{
			Fields: fields,
			Data:   data,
			Dense:  false,
		},
	}
}

// NewCardList creates a card with list content (key-value display).
// Use NewCardListContent() for simple name/content rows,
// or NewCardListContentFields() for custom field definitions.
// Example:
//
//	lines := []CardListContentLine{{Name: "Status", Content: "Active"}, {Name: "Version", Content: "2.1"}}
//	card := NewCardList("Details", NewCardListContent(lines))
func NewCardList(header string, content *CardListContent) *Card {
	return NewCard(
		core.CardTypeTable,
		content,
		header,
		nil,   // headerSub
		nil,   // headerIcon
		nil,   // headerIconColor
		false, // translate
		false, // forceMinWidth
		nil,   // display
	)
}

// CardListField represents a field definition in a card list
type CardListField struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Format   string  `json:"format"`
	Display  *string `json:"display,omitempty"`
	MinWidth *string `json:"minWidth,omitempty"`
}

// CardListContent represents the content structure for a card list
type CardListContent struct {
	CardContent
}

// CardListContentLine represents a single line in a card list
type CardListContentLine struct {
	Name    string
	Content string
}

// NewCardListContent creates a new CardListContent with default id/text fields
// Use SetDense() to enable dense mode
func NewCardListContent(data []CardListContentLine) *CardListContent {

	listData := make([]map[string]string, 0, len(data))
	for _, line := range data {
		listData = append(listData, map[string]string{
			"id":   line.Name,
			"text": line.Content,
		})
	}

	// Default fields for simple id/text display
	display1 := "info"
	display2 := "right"
	minWidth := "30px"
	fields := []CardListField{
		{
			ID:       "id",
			Name:     "Name",
			Format:   "text",
			Display:  &display1,
			MinWidth: &minWidth,
		},
		{
			ID:       "text",
			Name:     "Text",
			Format:   "text",
			Display:  &display2,
			MinWidth: &minWidth,
		},
	}
	return &CardListContent{
		CardContent{
			Fields: fields,
			Data:   listData,
			Dense:  false,
		},
	}
}

// CardListIconContentLine represents a single line in a card list with an icon on the right.
type CardListIconContentLine struct {
	Name      string     // Left column label
	Icon      string     // Right column icon name (Material Symbol)
	IconColor core.Color // Right column icon color
	IconHint  string     // Optional tooltip for the icon
}

// NewCardListIconContent creates a CardListContent that displays icons on the right
// instead of text. Each line shows a label on the left and an icon on the right.
// Use SetDense() to enable dense mode.
//
// Example:
//
//	lines := []CardListIconContentLine{{Name: "Status", Icon: "check_circle", IconColor: core.ColorSuccess, IconHint: "Active"}}
//	card := NewCardList("Details", NewCardListIconContent(lines))
func NewCardListIconContent(data []CardListIconContentLine) *CardListContent {

	icons := make(map[string]map[string]any, len(data))
	listData := make([]map[string]string, 0, len(data))

	for i, line := range data {
		key := strconv.Itoa(i)
		icons[key] = map[string]any{
			"icon":  line.Icon,
			"color": string(line.IconColor),
			"hint":  line.IconHint,
		}
		listData = append(listData, map[string]string{
			"id":   line.Name,
			"text": key,
		})
	}

	display := "info"
	minWidth := "30px"
	fields := []map[string]any{
		{"id": "id", "name": "Name", "format": "text", "display": display, "minWidth": minWidth},
		{"id": "text", "name": "Icon", "format": "icon", "display": "right", "icons": icons},
	}

	return &CardListContent{
		CardContent{
			Fields: fields,
			Data:   listData,
			Dense:  false,
		},
	}
}

// NewCardListContentFields creates a new CardListContent with NOT default fields
// if fields=nil, use NewCardListContent() instead
// Use SetDense() to enable dense mode
func NewCardListContentFields(fields []CardListField, data []map[string]interface{}) *CardListContent {

	return &CardListContent{
		CardContent{
			Fields: fields,
			Data:   data,
			Dense:  false,
		},
	}
}
