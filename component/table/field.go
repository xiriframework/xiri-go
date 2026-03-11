package table

import "github.com/xiriframework/xiri-go/component/core"

// fieldBase contains all type-independent field properties.
// Extracted from field[T] to avoid generic monomorphization of methods
// that don't depend on the row type T.
type fieldBase struct {
	// Identity
	id        string
	name      string // Translation key
	fieldType fieldType

	// Formatting - default + per-output overrides
	defaultFormatter OutputFormatter                // Used for all outputs by default
	formatters       map[OutputType]OutputFormatter // Overrides for specific output types

	// Field type hint tracking (for formatter recreation via WithXXX methods)
	fieldTypeHint fieldTypeHint // Semantic type (integer, float, distanceHint, etc.)
	decimals      int           // For Float, Distance, Pressure, Speed
	boolTrueText  string        // For Bool
	boolFalseText string        // For Bool

	// Display configuration (maps to TableFieldJSON properties)
	footer   FieldFooter
	hide     bool
	csv      bool // Include in CSV export
	align    *FieldAlign
	width    *string
	minWidth *string
	hint     *string
	display  *string

	// Display and behavior configuration
	search      bool    // Column is searchable
	sort        bool    // Column is sortable
	sticky      bool    // Column is sticky (fixed during scroll)
	header      *string // Custom header text override
	headerSpan  *int    // Header column span for grouped headers
	columnOrder int     // Column ordering/positioning

	// Input field configuration (for fieldTypeInput)
	inputType     *string // Input type (text, number, email, etc.)
	inputRequired *bool   // Input is required
	inputLang     *string // Input language/locale
	inputPaste    *bool   // Paste functionality enabled

	// Text decoration
	textPrefix *string // Text prefix (e.g., "$", "€")
	textSuffix *string // Text suffix (e.g., " km", " €")

	// Inline editing
	editable           bool              // Field supports inline editing
	editableOptions    []editableOption  // Predefined options for select-based inline editing
	editableOptionsUrl string            // URL to load options dynamically per row

	// Access control (potentially handled at higher level)
	access []string // Required permissions

	// Type-specific data for buttons/icon fields
	buttons   map[int]*buttonDef
	icons     map[string]*iconDef
	menuItems map[int][]*menuItemDef // Per-button menu item definitions (key = button index)
}

// field represents a table column with type-safe accessor and formatting.
// Each field extracts a value from the row struct using an accessor function,
// then formats it using an OutputFormatter that adapts to different output types.
type field[T any] struct {
	fieldBase // Embedded non-generic base (all T-independent properties)

	// Data extraction (T-dependent)
	accessor func(T) any // Extract value from row struct

	// T-dependent accessors for icon fields and menu buttons
	hintAccessor  func(T) string          // Optional: extracts per-row hint text for icon fields
	menuAccessors map[int]func(T) []string // Per-button menu data accessor (key = button index)
}

// buttonDef defines a button in a buttons-type field
type buttonDef struct {
	action  FieldButtonAction
	icon    string
	color   core.Color
	hint    string
	options map[string]any // Additional custom options
}

// iconDef defines an icon mapping in an icon-type field
type iconDef struct {
	icon    string
	color   core.Color
	hint    string
	options map[string]any // Additional custom options
}

// editableOption defines a selectable option for inline editing.
type editableOption struct {
	value string     // Option value (sent to backend)
	label string     // Display label (shown to user)
	color core.Color // Optional color for chips display
}

// menuItemDef defines a menu item in a menu-type button
type menuItemDef struct {
	action FieldButtonAction // "link", "href", or "dialog"
	icon   string
	color  core.Color
	text   string // Translation key
}

// getFormatter returns the formatter for a specific output type.
// Returns the per-output override if one exists, otherwise returns the default formatter.
func (f *fieldBase) getFormatter(output OutputType) OutputFormatter {
	if formatter, ok := f.formatters[output]; ok {
		return formatter
	}
	return f.defaultFormatter
}

// format formats a value using the appropriate formatter for the output type.
//
// CRITICAL: For number-type fields on web output, this wraps the result in [display, value] array
// to maintain exact JSON compatibility with xiri-ui frontend expectations.
func (f *fieldBase) format(value any, row Row, output OutputType, ctx *core.UiContext) any {
	formatter := f.getFormatter(output)
	formatted := formatter.Format(value, row, output, ctx)

	// CRITICAL: Number type fields on web output MUST return [display, value] array
	// This is required for sortable number columns in xiri-ui Angular frontend
	if output == OutputWeb && f.fieldType == fieldTypeNumber {
		return []any{formatted, value}
	}

	return formatted
}

// addButton adds a button definition to a buttons-type field.
func (f *fieldBase) addButton(key int, action FieldButtonAction, icon string, color core.Color, hint string) {
	if f.buttons == nil {
		f.buttons = make(map[int]*buttonDef)
	}
	f.buttons[key] = &buttonDef{
		action:  action,
		icon:    icon,
		color:   color,
		hint:    hint,
		options: make(map[string]any),
	}
}

// addIcon adds an icon mapping to an icon-type field.
func (f *fieldBase) addIcon(value string, icon string, color core.Color, hint string) {
	if f.icons == nil {
		f.icons = make(map[string]*iconDef)
	}
	f.icons[value] = &iconDef{
		icon:    icon,
		color:   color,
		hint:    hint,
		options: make(map[string]any),
	}
}

// toTableField converts fieldBase to tableFieldJSON for JSON serialization.
// This ensures all field properties are preserved in the conversion for the Angular frontend.
func (f *fieldBase) toTableField() *tableFieldJSON {
	tf := &tableFieldJSON{
		ID:          f.id,
		fieldType:   f.fieldType,
		name:        f.name,
		footer:      f.footer,
		hide:        f.hide,
		csv:         f.csv,
		columnOrder: f.columnOrder,

		// behavior
		search: f.search,
		sort:   f.sort,
		sticky: f.sticky,

		// Unexported fields - display
		width:    f.width,
		minWidth: f.minWidth,
		hint:     f.hint,
		display:  f.display,
		align:    f.align,

		// Unexported fields - header
		header:     f.header,
		headerSpan: f.headerSpan,

		// Unexported fields - input
		inputType:     f.inputType,
		inputRequired: f.inputRequired,
		inputLang:     f.inputLang,
		inputPaste:    f.inputPaste,

		// Unexported fields - text decoration
		textPrefix: f.textPrefix,
		textSuffix: f.textSuffix,

		// Unexported fields - inline editing
		editable:           f.editable,
		editableOptions:    f.editableOptions,
		editableOptionsUrl: f.editableOptionsUrl,

		// Unexported fields - access control
		access: f.access,

		// Initialize empty maps
		buttons: make(map[int]*fieldButtonJSON),
		icons:   make(map[string]*fieldIcon),
	}

	// Copy buttons
	if f.fieldType == fieldTypeButtons && len(f.buttons) > 0 {
		for key, btn := range f.buttons {
			icon := newFieldIcon(btn.icon, btn.color, btn.hint)
			for k, v := range btn.options {
				icon.withOption(k, v)
			}
			tf.buttons[key] = &fieldButtonJSON{
				action: btn.action,
				icon:   icon,
			}
		}
	}

	// Copy menu items into button JSON
	if f.fieldType == fieldTypeButtons && len(f.menuItems) > 0 {
		for key, items := range f.menuItems {
			if btn, ok := tf.buttons[key]; ok {
				btn.menuItems = items
			}
		}
	}

	// Copy icons
	if f.fieldType == fieldTypeIcon && len(f.icons) > 0 {
		for value, def := range f.icons {
			icon := newFieldIcon(def.icon, def.color, def.hint)
			for k, v := range def.options {
				icon.withOption(k, v)
			}
			tf.icons[value] = icon
		}
	}

	return tf
}
