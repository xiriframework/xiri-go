// Package layout provides layout utility components (Spacer, Container, Header, Html) for the Angular frontend.
package layout

import "github.com/xiriframework/xiri-go/component/core"

// Spacer represents an empty spacer component for layout
type Spacer struct {
	display *string
}

// NewSpacer creates a new spacer component
// Spacers are empty divs used to add space between components
func NewSpacer(display *string) *Spacer {
	return &Spacer{
		display: display,
	}
}

// Print returns the JSON representation of the spacer
func (s *Spacer) Print(translator core.TranslateFunc) map[string]any {
	return map[string]any{
		"type":    "spacer",
		"display": s.display,
		"data":    nil,
	}
}

// Container represents a container component for grouping nested components
type Container struct {
	components []core.Component
	display    *string
}

// NewContainer creates a new container component
// Containers allow grouping multiple components together
func NewContainer(display *string) *Container {
	return &Container{
		components: make([]core.Component, 0),
		display:    display,
	}
}

// Add adds a component to the container.
func (c *Container) Add(component core.Component) *Container {
	c.components = append(c.components, component)
	return c
}

// Print returns the JSON representation of the container
func (c *Container) Print(translator core.TranslateFunc) map[string]any {
	// Print all nested components
	componentData := make([]map[string]any, 0, len(c.components))
	for _, comp := range c.components {
		componentData = append(componentData, comp.Print(translator))
	}

	return map[string]any{
		"type":    "container",
		"display": c.display,
		"data": map[string]any{
			"components": componentData,
		},
	}
}

// Header represents a header/title component
type Header struct {
	text    string
	color   core.Color
	size    *string
	display *string
}

// NewHeader creates a new header component with required parameters
// Optional parameters can be set using builder methods: WithSize(), WithDisplay()
func NewHeader(text string, color core.Color, size *string, display *string) *Header {
	return &Header{
		text:    text,
		color:   color,
		size:    size,
		display: display,
	}
}

// WithSize sets the header size (optional)
// Returns the Header for method chaining
func (h *Header) WithSize(size string) *Header {
	h.size = &size
	return h
}

// WithDisplay sets the display/layout class.
func (h *Header) WithDisplay(display string) *Header {
	h.display = &display
	return h
}

// Print returns the JSON representation of the header
func (h *Header) Print(translator core.TranslateFunc) map[string]any {
	return map[string]any{
		"type":    "header",
		"display": h.display,
		"data": map[string]any{
			"text":  h.text,
			"color": string(h.color),
			"size":  h.size,
		},
	}
}

// Divider represents a visual separator component with optional text and icon.
type Divider struct {
	text    *string
	icon    *string
	spacing *string
	display *string
}

// NewDivider creates a new divider component.
func NewDivider() *Divider {
	return &Divider{}
}

// Text sets the divider label text.
func (d *Divider) Text(t string) *Divider {
	d.text = &t
	return d
}

// Icon sets the divider icon.
func (d *Divider) Icon(i string) *Divider {
	d.icon = &i
	return d
}

// Spacing sets the divider spacing: "compact", "normal", or "large".
func (d *Divider) Spacing(s string) *Divider {
	d.spacing = &s
	return d
}

// WithDisplay sets the display/layout class.
func (d *Divider) WithDisplay(display string) *Divider {
	d.display = &display
	return d
}

// Print returns the JSON representation of the divider component.
func (d *Divider) Print(translator core.TranslateFunc) map[string]any {
	data := make(map[string]any)

	if d.text != nil {
		data["text"] = core.Translate(translator, *d.text)
	}
	if d.icon != nil {
		data["icon"] = *d.icon
	}
	if d.spacing != nil {
		data["spacing"] = *d.spacing
	}

	result := map[string]any{
		"type": "divider",
		"data": data,
	}

	if d.display != nil {
		result["display"] = *d.display
	}

	return result
}

// Html represents a raw HTML content component
type Html struct {
	text    string
	display *string
}

// NewHtml creates a new HTML component
func NewHtml(text string, display *string) *Html {
	return &Html{
		text:    text,
		display: display,
	}
}

// WithDisplay sets the display/layout class.
func (h *Html) WithDisplay(display string) *Html {
	h.display = &display
	return h
}

// Print returns the JSON representation of the HTML component
func (h *Html) Print(translator core.TranslateFunc) map[string]any {
	return map[string]any{
		"type":    "html",
		"display": h.display,
		"data": map[string]any{
			"html": h.text,
		},
	}
}
