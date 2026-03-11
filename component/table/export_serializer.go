// Package table - JSON serialization layer for Angular frontend
//
// This file provides types for converting field[T] objects to JSON format
// compatible with the Angular frontend. This is an internal implementation
// detail of the table export system.
//
// Call chain:
//
//	table.Print() → exportFields() → field.toTableField() →
//	tableFieldJSON.print() → map[string]any

package table

import "github.com/xiriframework/xiri-go/component/core"

// ============================================================================
// Icon Type
// ============================================================================

// fieldIcon represents icon properties for table field buttons and status indicators.
// Used for visual cues in the Angular frontend (status icons, action buttons, etc.).
type fieldIcon struct {
	icon    string         // Icon name (e.g., "done", "error", "warning")
	color   core.Color     // Color enum (primary, accent, warn, tertiary)
	hint    string         // Tooltip text shown on hover
	options map[string]any // Additional custom properties for frontend
}

// newFieldIcon creates a new fieldIcon with required fields.
func newFieldIcon(icon string, color core.Color, hint string) *fieldIcon {
	return &fieldIcon{
		icon:    icon,
		color:   color,
		hint:    hint,
		options: make(map[string]any),
	}
}

// withOption adds a custom option to the icon (builder pattern).
func (i *fieldIcon) withOption(key string, value any) *fieldIcon {
	i.options[key] = value
	return i
}

// toMap converts fieldIcon to map[string]any for JSON serialization.
func (i *fieldIcon) toMap(ctx *core.UiContext) map[string]any {
	data := map[string]any{
		"icon":  i.icon,
		"color": string(i.color),
		"hint":  core.Translate(ctx, i.hint),
	}
	// Merge custom options
	for k, v := range i.options {
		data[k] = v
	}
	return data
}

// ============================================================================
// Button Type
// ============================================================================

// fieldButtonJSON represents a button in a buttons-type table field.
// Serializes to JSON for Angular frontend rendering.
type fieldButtonJSON struct {
	action    FieldButtonAction // Button action type (edit, delete, view, etc.)
	icon      *fieldIcon        // Icon configuration with color and hint
	menuItems []*menuItemDef    // Only for action "menu": menu item definitions
}

// ============================================================================
// tableFieldJSON - JSON Export Type
// ============================================================================

// tableFieldJSON represents a table column/field in JSON format for the Angular frontend.
// It contains all properties needed for rendering and interaction in the UI.
//
// This type serves as a serialization adapter between the internal field[T] type
// and the JSON structure expected by the Angular frontend. It handles:
//   - Conditional field inclusion (only non-default values)
//   - Type-specific serialization (buttons, icons, input fields)
//   - Translation of labels and hints
//   - Proper JSON structure for various field types
type tableFieldJSON struct {
	ID          string
	fieldType   fieldType
	name        string
	footer      FieldFooter
	hide        bool
	csv         bool
	columnOrder int

	// behavior
	search bool
	sort   bool
	sticky bool

	// dimensions
	width    *string
	minWidth *string

	// display
	hint    *string
	display *string
	align   *FieldAlign

	// header
	header     *string
	headerSpan *int

	// input configuration
	inputType     *string
	inputRequired *bool
	inputLang     *string
	inputPaste    *bool

	// text decoration
	textPrefix *string
	textSuffix *string

	// inline editing
	editable           bool
	editableOptions    []editableOption
	editableOptionsUrl string

	// access control
	access []string

	// type-specific data
	buttons map[int]*fieldButtonJSON
	icons   map[string]*fieldIcon
}

// print returns the JSON representation of the table field for the Angular frontend.
//
// This method implements complex serialization logic:
//   - Conditionally includes fields (only if non-default)
//   - Handles type-specific rendering (buttons → array, icons → map)
//   - Translates user-facing strings
//   - Preserves exact JSON structure for backward compatibility
func (tf *tableFieldJSON) print(ctx *core.UiContext) map[string]any {
	ret := map[string]any{
		"id":     tf.ID,
		"name":   core.Translate(ctx, tf.name),
		"format": string(tf.fieldType),
	}

	// Conditionally add dimension fields
	if tf.width != nil {
		ret["width"] = *tf.width
	}
	if tf.minWidth != nil {
		ret["minWidth"] = *tf.minWidth
	}

	// Conditionally add display fields
	if tf.hint != nil {
		ret["hint"] = *tf.hint
	}
	if tf.display != nil {
		ret["display"] = *tf.display
	}
	if tf.align != nil {
		ret["align"] = string(*tf.align)
	}

	// Conditionally add header fields
	if tf.header != nil {
		ret["header"] = *tf.header
	}
	if tf.headerSpan != nil {
		ret["headerSpan"] = *tf.headerSpan
	}

	// Boolean flags - only include if not default
	if !tf.search {
		ret["search"] = false
	}
	if !tf.sort {
		ret["sort"] = false
	}
	if tf.sticky {
		ret["sticky"] = true
	}
	if tf.hide {
		ret["hide"] = true
	}
	if tf.footer != FieldFooterNo {
		ret["footer"] = string(tf.footer)
	}
	if tf.editable {
		ret["editable"] = true
	}
	if len(tf.editableOptions) > 0 {
		opts := make([]map[string]any, len(tf.editableOptions))
		for i, o := range tf.editableOptions {
			opts[i] = map[string]any{"value": o.value, "label": core.Translate(ctx, o.label)}
			if o.color != "" {
				opts[i]["color"] = string(o.color)
			}
		}
		ret["editableOptions"] = opts
	}
	if tf.editableOptionsUrl != "" {
		ret["editableOptionsUrl"] = tf.editableOptionsUrl
	}

	// Text decoration
	if tf.textPrefix != nil {
		ret["textPrefix"] = *tf.textPrefix
	}
	if tf.textSuffix != nil {
		ret["textSuffix"] = *tf.textSuffix
	}

	// Type-specific serialization
	if tf.fieldType == fieldTypeButtons {
		// Buttons: convert map to array indexed by button key
		// Find max key to size array properly
		maxKey := -1
		for key := range tf.buttons {
			if key > maxKey {
				maxKey = key
			}
		}

		// Create array with proper size
		buttonsArray := make([]any, maxKey+1)
		for key, button := range tf.buttons {
			iconData := button.icon.toMap(ctx)
			iconData["action"] = string(button.action)
			if button.action == FieldButtonActionMenu && len(button.menuItems) > 0 {
				menuItems := make([]map[string]any, len(button.menuItems))
				for j, item := range button.menuItems {
					menuItems[j] = map[string]any{
						"action": string(item.action),
						"icon":   item.icon,
						"color":  string(item.color),
						"text":   core.Translate(ctx, item.text),
					}
				}
				iconData["menuItems"] = menuItems
			}
			buttonsArray[key] = iconData
		}
		ret["buttons"] = buttonsArray

	} else if tf.fieldType == fieldTypeIcon {
		// Icons: convert map to JSON object
		iconsMap := make(map[string]any)
		for value, icon := range tf.icons {
			iconsMap[value] = icon.toMap(ctx)
		}
		ret["icons"] = iconsMap

	} else if tf.fieldType == fieldTypeInput {
		// Input fields: add input-specific configuration
		ret["inputType"] = tf.inputType
		ret["inputRequired"] = tf.inputRequired
		ret["inputLang"] = tf.inputLang
		ret["inputPaste"] = tf.inputPaste
		ret["search"] = false
		ret["sort"] = false
	}

	return ret
}

// ============================================================================
// Helper Functions
// ============================================================================

