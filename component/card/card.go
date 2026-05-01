// Package card provides card container components for the Angular frontend.
package card

import (
	"github.com/xiriframework/xiri-go/component/button"
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
	"github.com/xiriframework/xiri-go/response"
)

// Card represents a card component with header and content
type Card struct {
	cardType        core.CardType
	content         any // Can be map[string]any or struct (e.g., *CardListContent)
	components      []core.Component // Optional sub-components rendered in xcol-grid (multi-component card)
	header          string
	headerSub       *string
	headerIcon      *string
	headerIconColor *string
	translate       bool
	forceMinWidth   bool
	display         *string
	buttonsTop      []*button.Button
	buttonsBottom   []*button.Button
	url             *url.Url
	reload          *bool
	collapsible     *bool
	collapsed       *bool
	maxHeight       *string
	padding         *string
}

// NewCard creates a new card component with full control over all parameters.
// Prefer using the convenience constructors instead:
//   - NewCardLinks(header, content) — for cards displaying a list of clickable links
//   - NewCardList(header, content) — for cards displaying key-value data rows
//
// Use NewCard directly only when you need custom CardType or advanced options
// like translateHeader, forceMinWidth, or headerIcon.
func NewCard(
	cardType core.CardType,
	content any, // Can be map[string]any or struct (e.g., *CardListContent)
	header string,
	headerSub *string,
	headerIcon *string,
	headerIconColor *string,
	translateHeader bool,
	forceMinWidth bool,
	display *string,
) *Card {
	return &Card{
		cardType:        cardType,
		content:         content,
		components:      make([]core.Component, 0),
		header:          header,
		headerSub:       headerSub,
		headerIcon:      headerIcon,
		headerIconColor: headerIconColor,
		translate:       translateHeader,
		forceMinWidth:   forceMinWidth,
		display:         display,
		buttonsTop:      make([]*button.Button, 0),
		buttonsBottom:   make([]*button.Button, 0),
	}
}

// Add appends a sub-component to the card. When at least one component is added,
// the frontend renders the card as a multi-component card (header + xcol-grid of
// the sub-components) and ignores the `content` / table fields.
//
// Layout per sub-component is controlled via the component's WithDisplay("xcol-…").
//
// Example:
//
//	c := card.NewCard(core.CardTypeTable, nil, "Activity", nil, nil, nil, false, false, nil)
//	c.Add(barchart.New("activity").Mode(barchart.ModeSimple)./* … */ .WithDisplay("xcol-12")).
//	  Add(stat.New("18h", "Today").Compact().WithDisplay("xcol-6")).
//	  Add(stat.New("32h", "Last 7 days").Compact().WithDisplay("xcol-6"))
func (c *Card) Add(component core.Component) *Card {
	c.components = append(c.components, component)
	return c
}

// ButtonTop adds a button to the top of the card
func (c *Card) ButtonTop(btn *button.Button) *Card {
	c.buttonsTop = append(c.buttonsTop, btn)
	return c
}

// ButtonBottom adds a button to the bottom of the card
func (c *Card) ButtonBottom(btn *button.Button) *Card {
	c.buttonsBottom = append(c.buttonsBottom, btn)
	return c
}

// WithHeaderSub sets the subtitle text (optional)
func (c *Card) WithHeaderSub(headerSub string) *Card {
	c.headerSub = &headerSub
	return c
}

// WithHeaderIcon sets the header icon (optional)
func (c *Card) WithHeaderIcon(icon string) *Card {
	c.headerIcon = &icon
	return c
}

// WithHeaderIconColor sets the header icon color (optional)
func (c *Card) WithHeaderIconColor(color string) *Card {
	c.headerIconColor = &color
	return c
}

// WithTranslate sets whether to translate the header text (optional)
func (c *Card) WithTranslate(translateHeader bool) *Card {
	c.translate = translateHeader
	return c
}

// WithForceMinWidth sets whether to force minimum width (optional)
func (c *Card) WithForceMinWidth(force bool) *Card {
	c.forceMinWidth = force
	return c
}

// WithDisplay sets the display/layout class (optional)
func (c *Card) WithDisplay(display string) *Card {
	c.display = &display
	return c
}

// SetURL sets the AJAX data URL. When set, static content is cleared and the frontend loads data dynamically.
func (c *Card) SetURL(url *url.Url) *Card {
	c.url = url
	c.content = nil
	return c
}

// WithReload enables periodic reload of the card data when using AJAX mode.
func (c *Card) WithReload(reload bool) *Card {
	c.reload = &reload
	return c
}

// WithCollapsible enables the card to be collapsed/expanded by the user.
func (c *Card) WithCollapsible(collapsible bool) *Card {
	c.collapsible = &collapsible
	return c
}

// WithCollapsed sets whether the card starts in collapsed state.
func (c *Card) WithCollapsed(collapsed bool) *Card {
	c.collapsed = &collapsed
	return c
}

// WithMaxHeight sets the max-height of the card content area (e.g. "300px", "80vh").
// The frontend default is "50vh".
func (c *Card) WithMaxHeight(maxHeight string) *Card {
	c.maxHeight = &maxHeight
	return c
}

// WithPadding sets the inner padding of the card content area in multi-component
// mode (i.e. when sub-components were added via Add(...)). Tabellen-Cards bleiben
// unverändert (padding:0 für Tabellen-Fluss).
//
// Accepted values:
//   - Token: "xs" | "sm" | "md" | "lg" | "xl" — vom Frontend auf --xiri-spacing-*
//     gemappt.
//   - Freier CSS-Wert: "16px", "1rem", "var(--xiri-spacing-md)".
//
// Default ist "md". **Auf xs-Viewport (<576px) wird immer 8px verwendet**, der
// hier konfigurierte Wert greift erst ab sm.
func (c *Card) WithPadding(value string) *Card {
	c.padding = &value
	return c
}

// Print returns the JSON representation of the card
func (c *Card) Print(ctx *core.UiContext) map[string]any {
	var data map[string]any

	if c.url != nil {
		// AJAX mode: header + buttons + URL, skip static content
		data = c.printHeader(ctx)
		data["url"] = c.url.PrintPrefix()
		if c.reload != nil {
			data["reload"] = *c.reload
		}
	} else {
		data = c.printData(ctx)
	}

	return map[string]any{
		"type":    "card",
		"display": c.display,
		"data":    data,
	}
}

// PrintData returns only the data portion of the card (for use in data endpoints).
func (c *Card) PrintData(ctx *core.UiContext) map[string]any {
	return c.printData(ctx)
}

// DataResponse returns a DataResult wrapping the card data in {"data": ...} envelope.
func (c *Card) DataResponse(ctx *core.UiContext) response.DataResult {
	return response.NewJSONDataResult(c.PrintData(ctx))
}

// printHeader builds the header/buttons/type map shared by both URL and static paths.
func (c *Card) printHeader(ctx *core.UiContext) map[string]any {
	data := map[string]any{
		"header":          c.header,
		"headerSub":       c.headerSub,
		"headerIcon":      c.headerIcon,
		"headerIconColor": c.headerIconColor,
		"type":            string(c.cardType),
		"buttonsBottom":   nil,
		"buttonsTop":      nil,
	}

	// Translate header if needed
	if c.translate {
		data["header"] = core.Translate(ctx, c.header)
	}

	// Add top buttons if any
	if len(c.buttonsTop) > 0 {
		btnLine := button.NewButtonLine("right", nil)
		for _, btn := range c.buttonsTop {
			btnLine.Add(btn)
		}
		data["buttonsTop"] = btnLine.PrintData(ctx)
	}

	// Add bottom buttons if any
	if len(c.buttonsBottom) > 0 {
		btnLine := button.NewButtonLine("right", nil)
		for _, btn := range c.buttonsBottom {
			btnLine.Add(btn)
		}
		data["buttonsBottom"] = btnLine.PrintData(ctx)
	}

	if c.collapsible != nil {
		data["collapsible"] = *c.collapsible
	}
	if c.collapsed != nil {
		data["collapsed"] = *c.collapsed
	}
	if c.maxHeight != nil {
		data["maxHeight"] = *c.maxHeight
	}
	if c.padding != nil {
		data["padding"] = *c.padding
	}

	return data
}

// printData builds the full data map (header + content) used by both Print and PrintData.
func (c *Card) printData(ctx *core.UiContext) map[string]any {
	data := c.printHeader(ctx)

	// Multi-component mode: if Add(...) was used, serialize sub-components and
	// skip table-/content-rendering. Frontend renders these via xiri-dyncomponent.
	if len(c.components) > 0 {
		componentData := make([]map[string]any, len(c.components))
		for i, comp := range c.components {
			componentData[i] = comp.Print(ctx)
		}
		data["components"] = componentData
		return data
	}

	// Add content based on card type
	if c.cardType == core.CardTypeTable {
		// For table type, handle both map and struct content
		if contentStruct, ok := c.content.(*CardListContent); ok {
			data["fields"] = contentStruct.Fields
			data["data"] = contentStruct.Data
			data["dense"] = contentStruct.Dense
		} else if contentStruct, ok := c.content.(*CardLinkContent); ok {
			data["fields"] = contentStruct.Fields
			data["data"] = contentStruct.Data
			data["dense"] = contentStruct.Dense
		} else if contentMap, ok := c.content.(map[string]any); ok {
			if fields, ok := contentMap["fields"]; ok {
				data["fields"] = fields
			}
			if tableData, ok := contentMap["data"]; ok {
				data["data"] = tableData
			}
			if dense, ok := contentMap["dense"]; ok {
				data["dense"] = dense
			}
		}
		data["forceMinWidth"] = c.forceMinWidth
	} else {
		data["content"] = c.content
	}

	return data
}
