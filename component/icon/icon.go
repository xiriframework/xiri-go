// Package icon provides icon components for the Angular frontend.
package icon

import "github.com/xiriframework/xiri-go/component/core"

// Icon represents an icon component with color and hint
type Icon struct {
	icon    string
	hint    string
	color   core.Color
	options map[string]any
	data    map[string]any
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

// WithOptions sets additional custom options merged at the top level of the
// rendered JSON. For custom payload consumed by the Angular frontend, prefer
// WithData.
//
// Deprecated: For custom payload, use WithData.
func (i *Icon) WithOptions(options map[string]any) *Icon {
	if options != nil {
		i.options = options
	}
	return i
}

// WithOption sets a single custom option merged at the top level of the
// rendered JSON. For custom payload consumed by the Angular frontend, prefer
// WithData.
//
// Deprecated: For custom payload, use WithData.
func (i *Icon) WithOption(key string, value any) *Icon {
	i.options[key] = value
	return i
}

// WithData sets the custom payload rendered under the explicit top-level
// "data" key. Calling with a nil or empty map omits the "data" key entirely.
// If both WithData and WithOption("data", …) are set, WithData wins.
func (i *Icon) WithData(payload map[string]any) *Icon {
	i.data = payload
	return i
}

// Print returns the JSON representation of the icon
func (i *Icon) Print(ctx *core.UiContext) map[string]any {
	data := map[string]any{
		"icon":  i.icon,
		"color": string(i.color),
		"hint":  core.Translate(ctx, i.hint),
	}

	// Merge options into data (legacy / structural top-level fields)
	for key, value := range i.options {
		data[key] = value
	}

	// Add custom payload under the explicit "data" key (only if non-empty);
	// overrides any options["data"] set via legacy WithOption("data", …).
	if len(i.data) > 0 {
		data["data"] = i.data
	}

	return data
}
