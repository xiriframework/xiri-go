// Package status provides a status indicator (badge/dot/text) component for the
// Angular frontend. Das Label ist immer sichtbar — Farbe und Icon sind rein redundant.
package status

import "github.com/xiriframework/xiri-go/component/core"

// Tone is the semantic color level of a status. Never the sole carrier of meaning.
type Tone string

const (
	ToneNeutral Tone = "neutral"
	ToneInfo    Tone = "info"
	ToneSuccess Tone = "success"
	ToneWarning Tone = "warning"
	ToneError   Tone = "error"
)

// Variant is the visual form of the status indicator.
type Variant string

const (
	VariantBadge Variant = "badge"
	VariantDot   Variant = "dot"
	VariantText  Variant = "text"
)

// Status represents a status indicator component.
type Status struct {
	label     string
	tone      *Tone
	variant   *Variant
	icon      *string
	hint      *string
	ariaLabel *string
}

// New creates a new Status with the always-visible label.
func New(label string) *Status {
	return &Status{label: label}
}

// Tone sets the semantic color level (default 'neutral' on the frontend).
func (s *Status) Tone(t Tone) *Status {
	s.tone = &t
	return s
}

// Variant sets the visual form (default 'badge' on the frontend).
func (s *Status) Variant(v Variant) *Status {
	s.variant = &v
	return s
}

// Icon sets an optional Material icon name, redundant to the label.
func (s *Status) Icon(icon string) *Status {
	s.icon = &icon
	return s
}

// Hint sets an optional tooltip/extra text.
func (s *Status) Hint(hint string) *Status {
	s.hint = &hint
	return s
}

// AriaLabel overrides the accessible name.
func (s *Status) AriaLabel(label string) *Status {
	s.ariaLabel = &label
	return s
}

// Print returns the JSON representation of the status component. Only set fields
// are written; defaults (tone/variant) are resolved by the frontend.
func (s *Status) Print(ctx *core.UiContext) map[string]any {
	data := map[string]any{
		"label": core.Translate(ctx, s.label),
	}

	if s.tone != nil {
		data["tone"] = string(*s.tone)
	}
	if s.variant != nil {
		data["variant"] = string(*s.variant)
	}
	if s.icon != nil {
		data["icon"] = *s.icon
	}
	if s.hint != nil {
		data["hint"] = core.Translate(ctx, *s.hint)
	}
	if s.ariaLabel != nil {
		data["ariaLabel"] = core.Translate(ctx, *s.ariaLabel)
	}

	return map[string]any{
		"type": "status",
		"data": data,
	}
}
