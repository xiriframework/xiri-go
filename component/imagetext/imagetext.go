// Package imagetext provides a combined image and text display component for the Angular frontend.
package imagetext

import "github.com/xiriframework/xiri-go/component/core"

// ImageText represents a component displaying an image alongside text.
type ImageText struct {
	url             string
	info            string
	header          *string
	headerSub       *string
	headerIcon      *string
	headerIconColor *core.Color
	display         *string
}

// New creates a new ImageText component with the image URL and info text.
func New(url, info string) *ImageText {
	return &ImageText{
		url:  url,
		info: info,
	}
}

// Header sets the card title.
func (i *ImageText) Header(h string) *ImageText {
	i.header = &h
	return i
}

// HeaderSub sets the card subtitle.
func (i *ImageText) HeaderSub(s string) *ImageText {
	i.headerSub = &s
	return i
}

// HeaderIcon sets the header icon and its color.
func (i *ImageText) HeaderIcon(icon string, color core.Color) *ImageText {
	i.headerIcon = &icon
	i.headerIconColor = &color
	return i
}

// WithDisplay sets the display/layout class.
func (i *ImageText) WithDisplay(d string) *ImageText {
	i.display = &d
	return i
}

// Print returns the JSON representation of the imagetext component.
func (i *ImageText) Print(translator core.TranslateFunc) map[string]any {
	data := map[string]any{
		"url":  i.url,
		"info": i.info,
	}

	if i.header != nil {
		data["header"] = core.Translate(translator, *i.header)
	}
	if i.headerSub != nil {
		data["headerSub"] = core.Translate(translator, *i.headerSub)
	}
	if i.headerIcon != nil {
		data["headerIcon"] = *i.headerIcon
	}
	if i.headerIconColor != nil {
		data["headerIconColor"] = string(*i.headerIconColor)
	}

	result := map[string]any{
		"type": "imagetext",
		"data": data,
	}

	if i.display != nil {
		result["display"] = *i.display
	}

	return result
}
