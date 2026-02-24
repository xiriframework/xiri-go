package table

import (
	"github.com/xiriframework/xiri-go/uicontext"
)

// Field represents a table column with type-safe accessor and formatting.
// Each field extracts a value from the row struct using an accessor function,
// then formats it using an OutputFormatter that adapts to different output types.
type Field[T any] struct {
	// Identity
	id        string
	name      string // Translation key
	fieldType FieldType

	// Data extraction
	accessor func(T) any // Extract value from row struct

	// Formatting - default + per-output overrides
	defaultFormatter OutputFormatter                // Used for all outputs by default
	formatters       map[OutputType]OutputFormatter // Overrides for specific output types

	// Field type hint tracking (for formatter recreation via WithXXX methods)
	fieldTypeHint FieldTypeHint // Semantic type (Integer, Float, Distance, etc.)
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

	// Input field configuration (for FieldTypeInput)
	inputType     *string // Input type (text, number, email, etc.)
	inputRequired *bool   // Input is required
	inputLang     *string // Input language/locale
	inputPaste    *bool   // Paste functionality enabled

	// Text decoration
	textPrefix *string // Text prefix (e.g., "$", "€")
	textSuffix *string // Text suffix (e.g., " km", " €")

	// Access control (potentially handled at higher level)
	access []string // Required permissions

	// Type-specific data for buttons/icon fields
	buttons      map[int]*ButtonDef
	icons        map[string]*IconDef
	hintAccessor func(T) string // Optional: extracts per-row hint text for icon fields

	// Menu button data
	menuAccessors map[int]func(T) []string // Per-button menu data accessor (key = button index)
	menuItems     map[int][]*MenuItemDef   // Per-button menu item definitions (key = button index)
}

// ButtonDef defines a button in a buttons-type field
type ButtonDef struct {
	Action  FieldButtonAction
	Icon    string
	Color   FieldColor
	Hint    string
	Options map[string]any // Additional custom options
}

// IconDef defines an icon mapping in an icon-type field
type IconDef struct {
	Icon    string
	Color   FieldColor
	Hint    string
	Options map[string]any // Additional custom options
}

// MenuItemDef defines a menu item in a menu-type button
type MenuItemDef struct {
	Action FieldButtonAction // "link", "href", or "dialog"
	Icon   string
	Color  FieldColor
	Text   string // Translation key
}

// GetFormatter returns the formatter for a specific output type.
// Returns the per-output override if one exists, otherwise returns the default formatter.
func (f *Field[T]) GetFormatter(output OutputType) OutputFormatter {
	if formatter, ok := f.formatters[output]; ok {
		return formatter
	}
	return f.defaultFormatter
}

// Format formats a value using the appropriate formatter for the output type.
//
// CRITICAL: For number-type fields on web output, this wraps the result in [display, value] array
// to maintain exact JSON compatibility with xiri-ui frontend expectations.
func (f *Field[T]) Format(value any, row Row, output OutputType, ctx *uicontext.UiContext) any {
	formatter := f.GetFormatter(output)
	formatted := formatter.Format(value, row, output, ctx)

	// CRITICAL: Number type fields on web output MUST return [display, value] array
	// This is required for sortable number columns in xiri-ui Angular frontend
	if output == OutputWeb && f.fieldType == FieldTypeNumber {
		return []any{formatted, value}
	}

	return formatted
}

// AddButton adds a button definition to a buttons-type field.
//
// Parameters:
//   - key: Button index (0-based) - determines position in buttons array
//   - action: Button action type (link, dialog, api, etc.)
//   - icon: Icon name
//   - color: Icon color
//   - hint: Tooltip text (translation key)
func (f *Field[T]) AddButton(key int, action FieldButtonAction, icon string, color FieldColor, hint string) {
	if f.buttons == nil {
		f.buttons = make(map[int]*ButtonDef)
	}
	f.buttons[key] = &ButtonDef{
		Action:  action,
		Icon:    icon,
		Color:   color,
		Hint:    hint,
		Options: make(map[string]any),
	}
}

// AddIcon adds an icon mapping to an icon-type field.
//
// Parameters:
//   - value: The value that maps to this icon (as string)
//   - icon: Icon name
//   - color: Icon color
//   - hint: Tooltip text (translation key)
func (f *Field[T]) addIcon(value string, icon string, color FieldColor, hint string) {
	if f.icons == nil {
		f.icons = make(map[string]*IconDef)
	}
	f.icons[value] = &IconDef{
		Icon:    icon,
		Color:   color,
		Hint:    hint,
		Options: make(map[string]any),
	}
}

// GetID returns the field ID
func (f *Field[T]) GetID() string {
	return f.id
}

// GetName returns the field name (translation key)
func (f *Field[T]) GetName() string {
	return f.name
}

// GetFieldType returns the field type
func (f *Field[T]) GetFieldType() FieldType {
	return f.fieldType
}

// GetAccessor returns the accessor function
func (f *Field[T]) GetAccessor() func(T) any {
	return f.accessor
}

// IsHidden returns whether the field is hidden
func (f *Field[T]) IsHidden() bool {
	return f.hide
}

// IsCsvEnabled returns whether the field is included in CSV export
func (f *Field[T]) IsCsvEnabled() bool {
	return f.csv
}

// GetFooter returns the footer aggregation type
func (f *Field[T]) GetFooter() FieldFooter {
	return f.footer
}

// GetAlign returns the text alignment
func (f *Field[T]) GetAlign() *FieldAlign {
	return f.align
}

// GetWidth returns the column width
func (f *Field[T]) GetWidth() *string {
	return f.width
}

// GetMinWidth returns the minimum column width
func (f *Field[T]) GetMinWidth() *string {
	return f.minWidth
}

// GetHint returns the hint/tooltip text
func (f *Field[T]) GetHint() *string {
	return f.hint
}

// GetDisplay returns the display CSS class
func (f *Field[T]) GetDisplay() *string {
	return f.display
}

// GetButtons returns the button definitions (for buttons-type fields)
func (f *Field[T]) GetButtons() map[int]*ButtonDef {
	return f.buttons
}

// GetIcons returns the icon mappings (for icon-type fields)
func (f *Field[T]) GetIcons() map[string]*IconDef {
	return f.icons
}

// GetSearch returns searchable flag
func (f *Field[T]) GetSearch() bool {
	return f.search
}

// GetSort returns sortable flag
func (f *Field[T]) GetSort() bool {
	return f.sort
}

// GetSticky returns sticky flag
func (f *Field[T]) GetSticky() bool {
	return f.sticky
}

// GetHeader returns custom header text
func (f *Field[T]) GetHeader() *string {
	return f.header
}

// GetHeaderSpan returns header column span
func (f *Field[T]) GetHeaderSpan() *int {
	return f.headerSpan
}

// GetColumnOrder returns column order
func (f *Field[T]) GetColumnOrder() int {
	return f.columnOrder
}

// GetInputType returns input type
func (f *Field[T]) GetInputType() *string {
	return f.inputType
}

// GetInputRequired returns input required flag
func (f *Field[T]) GetInputRequired() *bool {
	return f.inputRequired
}

// GetInputLang returns input language
func (f *Field[T]) GetInputLang() *string {
	return f.inputLang
}

// GetInputPaste returns input paste flag
func (f *Field[T]) GetInputPaste() *bool {
	return f.inputPaste
}

// GetTextPrefix returns text prefix
func (f *Field[T]) GetTextPrefix() *string {
	return f.textPrefix
}

// GetTextSuffix returns text suffix
func (f *Field[T]) GetTextSuffix() *string {
	return f.textSuffix
}

// GetHintAccessor returns the per-row hint accessor function
func (f *Field[T]) GetHintAccessor() func(T) string {
	return f.hintAccessor
}

// GetMenuAccessors returns the per-button menu data accessors
func (f *Field[T]) GetMenuAccessors() map[int]func(T) []string {
	return f.menuAccessors
}

// GetMenuItems returns the per-button menu item definitions
func (f *Field[T]) GetMenuItems() map[int][]*MenuItemDef {
	return f.menuItems
}

// GetAccess returns required permissions
func (f *Field[T]) GetAccess() []string {
	return f.access
}

// GetFieldTypeHint returns the field type hint (for special processing like Link fields)
func (f *Field[T]) GetFieldTypeHint() FieldTypeHint {
	return f.fieldTypeHint
}

// setHide sets whether the field is hidden (used internally by HideField/ShowField).
func (f *Field[T]) setHide(hide bool) {
	f.hide = hide
}

// toTableField converts Field[T] to tableFieldJSON for JSON serialization.
// This ensures all field properties are preserved in the conversion for the Angular frontend.
func (f *Field[T]) toTableField() *tableFieldJSON {
	// Directly construct tableFieldJSON from Field[T]
	field := &tableFieldJSON{
		// Exported fields
		ID:          f.id,
		FieldType:   f.fieldType,
		Name:        f.name,
		Footer:      f.footer,
		Hide:        f.hide,
		Csv:         f.csv,
		ColumnOrder: f.columnOrder,

		// Unexported fields - behavior
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

		// Unexported fields - access control
		access: f.access,

		// Initialize empty maps
		buttons: make(map[int]*fieldButtonJSON),
		icons:   make(map[string]*fieldIcon),
	}

	// Copy buttons
	if f.fieldType == FieldTypeButtons && len(f.buttons) > 0 {
		for key, btn := range f.buttons {
			icon := newFieldIcon(btn.Icon, btn.Color, btn.Hint)
			for k, v := range btn.Options {
				icon.WithOption(k, v)
			}
			field.buttons[key] = &fieldButtonJSON{
				Action: btn.Action,
				Icon:   icon,
			}
		}
	}

	// Copy menu items into button JSON
	if f.fieldType == FieldTypeButtons && len(f.menuItems) > 0 {
		for key, items := range f.menuItems {
			if btn, ok := field.buttons[key]; ok {
				btn.MenuItems = items
			}
		}
	}

	// Copy icons
	if f.fieldType == FieldTypeIcon && len(f.icons) > 0 {
		for value, iconDef := range f.icons {
			icon := newFieldIcon(iconDef.Icon, iconDef.Color, iconDef.Hint)
			for k, v := range iconDef.Options {
				icon.WithOption(k, v)
			}
			field.icons[value] = icon
		}
	}

	return field
}
