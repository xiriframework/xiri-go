// Package icon provides icon components for the Angular frontend.
package icon

import "github.com/xiriframework/xiri-go/component/core"

// Icon represents an icon component with color and hint
type Icon struct {
	icon    string
	hint    string
	color   core.Color
	options map[string]any
}

// NewIcon creates a new icon component with required parameters
// Optional parameters can be set using builder methods: WithHint(), WithOptions()
func NewIcon(icon string, hint string, color core.Color, options map[string]any) *Icon {
	if options == nil {
		options = make(map[string]any)
	}
	return &Icon{
		icon:    icon,
		hint:    hint,
		color:   color,
		options: options,
	}
}

// WithHint sets the tooltip/hint text (optional)
// Returns the Icon for method chaining
func (i *Icon) WithHint(hint string) *Icon {
	i.hint = hint
	return i
}

// WithOptions sets additional custom options (optional)
// Returns the Icon for method chaining
func (i *Icon) WithOptions(options map[string]any) *Icon {
	if options != nil {
		i.options = options
	}
	return i
}

// WithOption sets a single custom option (optional)
// Returns the Icon for method chaining
func (i *Icon) WithOption(key string, value any) *Icon {
	i.options[key] = value
	return i
}

// Print returns the JSON representation of the icon
func (i *Icon) Print(translator core.TranslateFunc) map[string]any {
	data := map[string]any{
		"icon":  i.icon,
		"color": string(i.color),
		"hint":  core.Translate(translator, i.hint),
	}

	// Merge options into data
	for key, value := range i.options {
		data[key] = value
	}

	return data
}
