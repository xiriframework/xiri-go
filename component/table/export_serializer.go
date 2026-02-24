// Package table - JSON serialization layer for Angular frontend
//
// This file provides types for converting Field[T] objects to JSON format
// compatible with the Angular frontend. This is an internal implementation
// detail of the table export system.
//
// Call chain:
//   table.Print() → exportFields() → field.toTableField() →
//   tableFieldJSON.Print() → map[string]any
//
// DO NOT use these types directly from outside the table package.
// Use the public Field[T] API instead.

package table

import "github.com/xiriframework/xiri-go/component/core"

// ============================================================================
// Icon Type
// ============================================================================

// fieldIcon represents icon properties for table field buttons and status indicators.
// Used for visual cues in the Angular frontend (status icons, action buttons, etc.).
type fieldIcon struct {
	Icon    string         // Icon name (e.g., "done", "error", "warning")
	Color   FieldColor     // Color enum (primary, accent, warn, tertiary)
	Hint    string         // Tooltip text shown on hover
	Options map[string]any // Additional custom properties for frontend
}

// newFieldIcon creates a new fieldIcon with required fields.
func newFieldIcon(icon string, color FieldColor, hint string) *fieldIcon {
	return &fieldIcon{
		Icon:    icon,
		Color:   color,
		Hint:    hint,
		Options: make(map[string]any),
	}
}

// WithOption adds a custom option to the icon (builder pattern).
func (i *fieldIcon) WithOption(key string, value any) *fieldIcon {
	i.Options[key] = value
	return i
}

// ToMap converts fieldIcon to map[string]any for JSON serialization.
func (i *fieldIcon) ToMap(translator core.TranslateFunc) map[string]any {
	data := map[string]any{
		"icon":  i.Icon,
		"color": string(i.Color),
		"hint":  translate(translator, i.Hint),
	}
	// Merge custom options
	for k, v := range i.Options {
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
	Action    FieldButtonAction // Button action type (edit, delete, view, etc.)
	Icon      *fieldIcon        // Icon configuration with color and hint
	MenuItems []*MenuItemDef    // Only for action "menu": menu item definitions
}

// ============================================================================
// TableFieldJSON - JSON Export Type
// ============================================================================

// tableFieldJSON represents a table column/field in JSON format for the Angular frontend.
// It contains all properties needed for rendering and interaction in the UI.
//
// This type serves as a serialization adapter between the internal Field[T] type
// and the JSON structure expected by the Angular frontend. It handles:
//   - Conditional field inclusion (only non-default values)
//   - Type-specific serialization (buttons, icons, input fields)
//   - Translation of labels and hints
//   - Proper JSON structure for various field types
type tableFieldJSON struct {
	// Exported fields (always included in JSON)
	ID          string
	FieldType   FieldType
	Name        string
	Footer      FieldFooter
	Hide        bool
	Csv         bool
	ColumnOrder int

	// Unexported fields - behavior
	search bool
	sort   bool
	sticky bool

	// Unexported fields - dimensions
	width    *string
	minWidth *string

	// Unexported fields - display
	hint    *string
	display *string
	align   *FieldAlign

	// Unexported fields - header
	header     *string
	headerSpan *int

	// Unexported fields - input configuration
	inputType     *string
	inputRequired *bool
	inputLang     *string
	inputPaste    *bool

	// Unexported fields - text decoration
	textPrefix *string
	textSuffix *string

	// Unexported fields - access control
	access []string

	// Unexported fields - type-specific data
	buttons map[int]*fieldButtonJSON
	icons   map[string]*fieldIcon
}

// AddIcon adds an icon mapping for specific values (used for type=icon fields).
func (tf *tableFieldJSON) AddIcon(value string, icon *fieldIcon) *tableFieldJSON {
	tf.icons[value] = icon
	return tf
}

// AddButton adds a button to the field (used for type=buttons fields).
func (tf *tableFieldJSON) AddButton(key int, action FieldButtonAction, icon *fieldIcon) *tableFieldJSON {
	tf.buttons[key] = &fieldButtonJSON{
		Action: action,
		Icon:   icon,
	}
	return tf
}

// SetCsv sets whether this field should be included in CSV exports.
func (tf *tableFieldJSON) SetCsv(csv bool) *tableFieldJSON {
	tf.Csv = csv
	return tf
}

// Print returns the JSON representation of the table field for the Angular frontend.
//
// This method implements complex serialization logic:
//   - Conditionally includes fields (only if non-default)
//   - Handles type-specific rendering (buttons → array, icons → map)
//   - Translates user-facing strings
//   - Preserves exact JSON structure for backward compatibility
func (tf *tableFieldJSON) Print(translator core.TranslateFunc) map[string]any {
	ret := map[string]any{
		"id":     tf.ID,
		"name":   translate(translator, tf.Name),
		"format": string(tf.FieldType),
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
	if tf.Hide {
		ret["hide"] = true
	}
	if tf.Footer != FieldFooterNo {
		ret["footer"] = string(tf.Footer)
	}

	// Text decoration
	if tf.textPrefix != nil {
		ret["textPrefix"] = *tf.textPrefix
	}
	if tf.textSuffix != nil {
		ret["textSuffix"] = *tf.textSuffix
	}

	// Type-specific serialization
	if tf.FieldType == FieldTypeButtons {
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
			iconData := button.Icon.ToMap(translator)
			iconData["action"] = string(button.Action)
			if button.Action == FieldButtonActionMenu && len(button.MenuItems) > 0 {
				menuItems := make([]map[string]any, len(button.MenuItems))
				for j, item := range button.MenuItems {
					menuItems[j] = map[string]any{
						"action": string(item.Action),
						"icon":   item.Icon,
						"color":  string(item.Color),
						"text":   translate(translator, item.Text),
					}
				}
				iconData["menuItems"] = menuItems
			}
			buttonsArray[key] = iconData
		}
		ret["buttons"] = buttonsArray

	} else if tf.FieldType == FieldTypeIcon {
		// Icons: convert map to JSON object
		iconsMap := make(map[string]any)
		for value, icon := range tf.icons {
			iconsMap[value] = icon.ToMap(translator)
		}
		ret["icons"] = iconsMap

	} else if tf.FieldType == FieldTypeInput {
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

// translate is a helper function for translating user-facing strings.
// Returns the original key if no translator is provided.
func translate(translator core.TranslateFunc, key string) string {
	if translator == nil {
		return key
	}
	return translator(key)
}
