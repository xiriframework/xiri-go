// Package info provides informational display components (InfoText, InfoPoint) for the Angular frontend.
package info

import "github.com/xiriframework/xiri-go/component/core"

// InfoText represents a simple informational text component
type InfoText struct {
	text    string
	display *string
}

// NewInfoText creates a new info text component
func NewInfoText(text string, display *string) *InfoText {
	return &InfoText{
		text:    text,
		display: display,
	}
}

// WithDisplay sets the display/layout class (optional)
// Returns the InfoText for method chaining
func (it *InfoText) WithDisplay(display string) *InfoText {
	it.display = &display
	return it
}

// Print returns the JSON representation of the info text
func (it *InfoText) Print(translator core.TranslateFunc) map[string]any {
	return map[string]any{
		"type":    "infotext",
		"display": it.display,
		"data": map[string]any{
			"text": it.text,
		},
	}
}

// InfoPoint represents an information display with icon and optional link
type InfoPoint struct {
	text      string
	icon      string
	iconColor string
	subtext   *string
	url       *string
	urlParams map[string]string
	iconSet   *string
	dense     *bool
	display   *string
}

// NewInfoPoint creates a new info point component with required parameters
// Optional parameters can be set using builder methods: WithSubtext(), WithUrl(), etc.
func NewInfoPoint(
	text string,
	icon string,
	iconColor string,
	subtext *string,
	url *string,
	urlParams map[string]string,
	iconSet *string,
	dense *bool,
	display *string,
) *InfoPoint {
	return &InfoPoint{
		text:      text,
		icon:      icon,
		iconColor: iconColor,
		subtext:   subtext,
		url:       url,
		urlParams: urlParams,
		iconSet:   iconSet,
		dense:     dense,
		display:   display,
	}
}

// WithSubtext sets the subtext/label (optional)
// Returns the InfoPoint for method chaining
func (ip *InfoPoint) WithSubtext(subtext string) *InfoPoint {
	ip.subtext = &subtext
	return ip
}

// WithUrl sets the link URL (optional)
// Returns the InfoPoint for method chaining
func (ip *InfoPoint) WithUrl(url string) *InfoPoint {
	ip.url = &url
	return ip
}

// WithUrlParams sets the URL query parameters (optional)
// Returns the InfoPoint for method chaining
func (ip *InfoPoint) WithUrlParams(params map[string]string) *InfoPoint {
	ip.urlParams = params
	return ip
}

// WithIconSet sets the icon set to use (optional, e.g., "material", "fontawesome")
// Returns the InfoPoint for method chaining
func (ip *InfoPoint) WithIconSet(iconSet string) *InfoPoint {
	ip.iconSet = &iconSet
	return ip
}

// WithDense sets whether to use dense/compact layout (optional)
// Returns the InfoPoint for method chaining
func (ip *InfoPoint) WithDense(dense bool) *InfoPoint {
	ip.dense = &dense
	return ip
}

// WithDisplay sets the display/layout class (optional)
// Returns the InfoPoint for method chaining
func (ip *InfoPoint) WithDisplay(display string) *InfoPoint {
	ip.display = &display
	return ip
}

// Print returns the JSON representation of the info point
func (ip *InfoPoint) Print(translator core.TranslateFunc) map[string]any {
	return map[string]any{
		"type":    "infopoint",
		"display": ip.display,
		"data": map[string]any{
			"info":      ip.text,
			"text":      ip.subtext,
			"icon":      ip.icon,
			"iconColor": ip.iconColor,
			"iconSet":   ip.iconSet,
			"url":       ip.url,
			"urlParams": ip.urlParams,
			"dense":     ip.dense,
		},
	}
}
