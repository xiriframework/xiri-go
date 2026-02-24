// Package button provides button components for the Angular frontend.
package button

import (
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
)

// Button represents a button component with action, styling, and behavior
type Button struct {
	action     core.ButtonAction
	text       string
	url        *url.Url
	color      core.Color
	buttonType core.ButtonType
	hint       string
	icon       string
	filename   *string
	disabled   bool
	tabIndex   *int
	isDefault  bool
	target     string
	options    map[string]any
}

// TableButton represents a type-safe button specifically for table actions.
// This is a wrapper around Button to provide type safety for table button operations.
type TableButton struct {
	button *Button
}

// Print renders the TableButton to JSON format by delegating to the underlying Button
func (tb *TableButton) Print(translator core.TranslateFunc) map[string]any {
	return tb.button.Print(translator)
}

// GetButton returns the underlying Button for internal use.
// This is primarily for compatibility with legacy code that works with *Button.
func (tb *TableButton) GetButton() *Button {
	return tb.button
}

// NewButton creates a new button with full configuration
func NewButton(
	action core.ButtonAction,
	text string,
	u *url.Url,
	color core.Color,
	buttonType core.ButtonType,
	hint string,
	icon string,
	disabled bool,
	tabIndex *int,
	isDefault bool,
	target string,
	options map[string]any,
) *Button {
	if options == nil {
		options = make(map[string]any)
	}
	if tabIndex == nil {
		defaultTabIndex := -1
		tabIndex = &defaultTabIndex
	}
	return &Button{
		action:     action,
		text:       text,
		url:        u,
		color:      color,
		buttonType: buttonType,
		hint:       hint,
		icon:       icon,
		filename:   nil,
		disabled:   disabled,
		tabIndex:   tabIndex,
		isDefault:  isDefault,
		target:     target,
		options:    options,
	}
}

// NewApiButton creates a button that triggers an API call
func NewApiButton(
	text string,
	u *url.Url,
	color core.Color,
	buttonType core.ButtonType,
	hint string,
	disabled bool,
	tabIndex *int,
	options map[string]any,
) *Button {
	return NewButton(
		core.ButtonActionApi,
		text,
		u,
		color,
		buttonType,
		hint,
		"",
		disabled,
		tabIndex,
		false,
		"_self",
		options,
	)
}

// NewBackButton creates a back navigation button
func NewBackButton(
	text string,
	color core.Color,
	buttonType core.ButtonType,
	hint string,
	disabled bool,
	tabIndex *int,
	isDefault bool,
	options map[string]any,
) *Button {
	return NewButton(
		core.ButtonActionBack,
		text,
		url.NewUrl(""),
		color,
		buttonType,
		hint,
		"",
		disabled,
		tabIndex,
		isDefault,
		"_self",
		options,
	)
}

// NewCloseButton creates a close button
func NewCloseButton(
	text string,
	color core.Color,
	buttonType core.ButtonType,
	hint string,
	disabled bool,
	tabIndex *int,
	isDefault bool,
	options map[string]any,
) *Button {
	return NewButton(
		core.ButtonActionClose,
		text,
		url.NewUrl(""),
		color,
		buttonType,
		hint,
		"",
		disabled,
		tabIndex,
		isDefault,
		"_self",
		options,
	)
}

// NewDialogButton creates a button that opens a dialog
func NewDialogButton(
	text string,
	u *url.Url,
	color core.Color,
	buttonType core.ButtonType,
	hint string,
	disabled bool,
	tabIndex *int,
	options map[string]any,
) *Button {
	return NewButton(
		core.ButtonActionDialog,
		text,
		u,
		color,
		buttonType,
		hint,
		"",
		disabled,
		tabIndex,
		false,
		"_self",
		options,
	)
}

// NewDownloadButton creates a download button
func NewDownloadButton(
	text string,
	u *url.Url,
	color core.Color,
	buttonType core.ButtonType,
	hint string,
	filename *string,
	disabled bool,
	tabIndex *int,
	options map[string]any,
) *Button {
	btn := NewButton(
		core.ButtonActionDownload,
		text,
		u,
		color,
		buttonType,
		hint,
		"",
		disabled,
		tabIndex,
		false,
		"_self",
		options,
	)
	btn.filename = filename
	return btn
}

// NewFormButton creates a form submit button
// Note: tabIndex defaults to nil (not -1)
func NewFormButton(
	text string,
	u *url.Url,
	color core.Color,
	buttonType core.ButtonType,
	hint string,
	disabled bool,
	tabIndex *int,
	isDefault bool,
	options map[string]any,
) *Button {
	// Don't force default tabIndex like other constructors - pass through as-is
	return NewButton(
		core.ButtonActionForm,
		text,
		u,
		color,
		buttonType,
		hint,
		"",
		disabled,
		tabIndex, // Pass through nil if provided
		isDefault,
		"_self",
		options,
	)
}

// NewLinkButton creates a navigation link button
func NewLinkButton(
	text string,
	u *url.Url,
	color core.Color,
	buttonType core.ButtonType,
	hint string,
	disabled bool,
	tabIndex *int,
	options map[string]any,
) *Button {
	return NewButton(
		core.ButtonActionLink,
		text,
		u,
		color,
		buttonType,
		hint,
		"",
		disabled,
		tabIndex,
		false,
		"_self",
		options,
	)
}

// NewHrefButton creates an external link button
func NewHrefButton(
	text string,
	u *url.Url,
	color core.Color,
	buttonType core.ButtonType,
	hint string,
	disabled bool,
	tabIndex *int,
	target string,
	options map[string]any,
) *Button {
	return NewButton(
		core.ButtonActionHref,
		text,
		u,
		color,
		buttonType,
		hint,
		"",
		disabled,
		tabIndex,
		false,
		target,
		options,
	)
}

// NewTableButton creates a table action button (icon-only) with type safety
func NewTableButton(
	action core.ButtonAction,
	icon string,
	u *url.Url,
	hint string,
	color core.Color,
	disabled bool,
	options map[string]any,
) *TableButton {
	button := NewButton(
		action,
		icon, // Icon is stored in text field for icon buttons
		u,
		color,
		core.ButtonTypeIcon,
		hint,
		"",
		disabled,
		nil,
		false,
		"_self",
		options,
	)
	return &TableButton{button: button}
}

// ============================================================================
// Simplified Constructors
// ============================================================================

// NewSimpleCloseButton creates a close button with default styling (Primary, Stroked)
func NewSimpleCloseButton(text string) *Button {
	return NewCloseButton(text, core.ColorPrimary, core.ButtonTypeStroked, "", false, nil, false, nil)
}

// NewSimpleBackButton creates a back button with default styling (Primary, Stroked)
func NewSimpleBackButton(text string) *Button {
	return NewBackButton(text, core.ColorPrimary, core.ButtonTypeStroked, "", false, nil, false, nil)
}

// NewSimpleFormButton creates a form submit button with default styling (Primary, Raised, default=true)
func NewSimpleFormButton(text string, u *url.Url) *Button {
	return NewFormButton(text, u, core.ColorPrimary, core.ButtonTypeRaised, "", false, nil, true, nil)
}

// NewSimpleLinkButton creates a link button with the given color (Stroked)
func NewSimpleLinkButton(text string, u *url.Url, color core.Color) *Button {
	return NewLinkButton(text, u, color, core.ButtonTypeStroked, "", false, nil, nil)
}

// NewSimpleDialogButton creates a dialog button with the given color (Stroked)
func NewSimpleDialogButton(text string, u *url.Url, color core.Color) *Button {
	return NewDialogButton(text, u, color, core.ButtonTypeStroked, "", false, nil, nil)
}

// NewSimpleApiButton creates an API button with the given color (Stroked)
func NewSimpleApiButton(text string, u *url.Url, color core.Color) *Button {
	return NewApiButton(text, u, color, core.ButtonTypeStroked, "", false, nil, nil)
}

// DefaultFormButtons creates the standard Back + Save button pair used in forms.
// backText and saveText are typically translated strings.
func DefaultFormButtons(backText, saveText string, saveUrl *url.Url) []*Button {
	defaultTabIndex := -1
	return []*Button{
		NewBackButton(backText, core.ColorPrimary, core.ButtonTypeStroked, "", false, &defaultTabIndex, false, nil),
		NewFormButton(saveText, saveUrl, core.ColorPrimary, core.ButtonTypeRaised, "", false, nil, true, nil),
	}
}

// Builder methods for optional button properties

// WithHint sets the tooltip/hint text (optional)
// Returns the Button for method chaining
func (b *Button) WithHint(hint string) *Button {
	b.hint = hint
	return b
}

// WithTabIndex sets the tab index for keyboard navigation (optional)
// Returns the Button for method chaining
func (b *Button) WithTabIndex(tabIndex int) *Button {
	b.tabIndex = &tabIndex
	return b
}

// WithDisabled sets the disabled state (optional)
// Returns the Button for method chaining
func (b *Button) WithDisabled(disabled bool) *Button {
	b.disabled = disabled
	return b
}

// WithDefault marks this button as the default/primary action (optional)
// Returns the Button for method chaining
func (b *Button) WithDefault(isDefault bool) *Button {
	b.isDefault = isDefault
	return b
}

// WithTarget sets the link target (optional, e.g., "_blank", "_self")
// Returns the Button for method chaining
func (b *Button) WithTarget(target string) *Button {
	b.target = target
	return b
}

// WithOptions sets additional custom options (optional)
// Returns the Button for method chaining
func (b *Button) WithOptions(options map[string]any) *Button {
	if options != nil {
		b.options = options
	}
	return b
}

// WithOption sets a single custom option (optional)
// Returns the Button for method chaining
func (b *Button) WithOption(key string, value any) *Button {
	b.options[key] = value
	return b
}

// Print returns the JSON representation of the button
func (b *Button) Print(translator core.TranslateFunc) map[string]any {
	data := map[string]any{
		"action":   string(b.action),
		"url":      b.url.PrintPrefix(),
		"type":     string(b.buttonType),
		"color":    string(b.color),
		"disabled": b.disabled,
		"default":  b.isDefault,
		"hint":     core.Translate(translator, b.hint),
		"target":   b.target,
		"tabIndex": b.tabIndex,
	}

	// Determine which field to use for text/icon based on button type
	if b.buttonType == core.ButtonTypeIcon || b.buttonType == core.ButtonTypeFab || b.buttonType == core.ButtonTypeMiniFab {
		// For icon types, use icon field
		if b.icon != "" {
			data["icon"] = b.icon
		} else {
			data["icon"] = b.text
		}
	} else if b.buttonType == core.ButtonTypeIconText {
		// For icon+text type, include both
		data["icon"] = b.icon
		data["text"] = core.Translate(translator, b.text)
	} else {
		// For other types, use text field
		data["text"] = core.Translate(translator, b.text)
	}

	// Add filename for download action
	if b.action == core.ButtonActionDownload && b.filename != nil {
		data["filename"] = *b.filename
	}

	// Merge options into data
	for key, value := range b.options {
		data[key] = value
	}

	return data
}
