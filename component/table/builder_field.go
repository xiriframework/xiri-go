package table

import "github.com/xiriframework/xiri-go/component/core"

// FieldBuilder provides a fluent API for configuring a single field.
// This is non-generic — all methods operate on *fieldBase.
type FieldBuilder struct {
	base        *fieldBase
	typedField  any // *field[T] — for T-dependent operations only
	lastMenuKey int // Track last menu button index for AddMenuItem chaining
}

// WithFormatter sets the default formatter (used for all output types unless overridden)
func (fb *FieldBuilder) WithFormatter(formatter OutputFormatter) *FieldBuilder {
	fb.base.defaultFormatter = formatter
	return fb
}

// WithWebFormatter sets formatter specifically for web output (overrides default)
func (fb *FieldBuilder) WithWebFormatter(formatter OutputFormatter) *FieldBuilder {
	fb.base.formatters[OutputWeb] = formatter
	return fb
}

// WithCSVFormatter sets formatter specifically for CSV output (overrides default)
func (fb *FieldBuilder) WithCSVFormatter(formatter OutputFormatter) *FieldBuilder {
	fb.base.formatters[OutputCSV] = formatter
	return fb
}

// WithPDFFormatter sets formatter specifically for PDF output (overrides default)
func (fb *FieldBuilder) WithPDFFormatter(formatter OutputFormatter) *FieldBuilder {
	fb.base.formatters[OutputPDF] = formatter
	return fb
}

// WithExcelFormatter sets formatter specifically for Excel output (overrides default)
func (fb *FieldBuilder) WithExcelFormatter(formatter OutputFormatter) *FieldBuilder {
	fb.base.formatters[OutputExcel] = formatter
	return fb
}

// WithFooter sets footer aggregation type
func (fb *FieldBuilder) WithFooter(footer FieldFooter) *FieldBuilder {
	fb.base.footer = footer
	return fb
}

// WithFooterSum enables sum aggregation in footer
func (fb *FieldBuilder) WithFooterSum() *FieldBuilder {
	fb.base.footer = FieldFooterSum
	return fb
}

// WithFooterCount enables count aggregation in footer
func (fb *FieldBuilder) WithFooterCount() *FieldBuilder {
	fb.base.footer = FieldFooterCount
	return fb
}

// Hide hides the field (not displayed in any output)
func (fb *FieldBuilder) Hide() *FieldBuilder {
	fb.base.hide = true
	return fb
}

// HideInCSV excludes field from CSV export
func (fb *FieldBuilder) HideInCSV() *FieldBuilder {
	fb.base.csv = false
	return fb
}

// ShowInCSV includes field in CSV export (default)
func (fb *FieldBuilder) ShowInCSV() *FieldBuilder {
	fb.base.csv = true
	return fb
}

// WithAlign sets text alignment
func (fb *FieldBuilder) WithAlign(align FieldAlign) *FieldBuilder {
	fb.base.align = &align
	return fb
}

// AlignLeft sets left alignment
func (fb *FieldBuilder) AlignLeft() *FieldBuilder {
	align := FieldAlignLeft
	fb.base.align = &align
	return fb
}

// AlignCenter sets center alignment
func (fb *FieldBuilder) AlignCenter() *FieldBuilder {
	align := FieldAlignCenter
	fb.base.align = &align
	return fb
}

// AlignRight sets right alignment
func (fb *FieldBuilder) AlignRight() *FieldBuilder {
	align := FieldAlignRight
	fb.base.align = &align
	return fb
}

// WithWidth sets column width
func (fb *FieldBuilder) WithWidth(width string) *FieldBuilder {
	fb.base.width = &width
	return fb
}

// WithMinWidth sets minimum column width
func (fb *FieldBuilder) WithMinWidth(minWidth string) *FieldBuilder {
	fb.base.minWidth = &minWidth
	return fb
}

// WithHint sets tooltip/hint text
func (fb *FieldBuilder) WithHint(hint string) *FieldBuilder {
	fb.base.hint = &hint
	return fb
}

// WithDisplay sets CSS display class
func (fb *FieldBuilder) WithDisplay(display string) *FieldBuilder {
	fb.base.display = &display
	return fb
}

// AddButton adds a button to a buttons-type field
func (fb *FieldBuilder) AddButton(
	key int,
	action FieldButtonAction,
	icon string,
	color core.Color,
	hint string,
) *FieldBuilder {
	fb.base.addButton(key, action, icon, color, hint)
	return fb
}

// WithDecimals sets the number of decimal places for numeric fields.
// This is used with Float, Distance, Pressure, Speed, and Text2 numeric field types.
//
// Example:
//
//	builder.Field("distance", "trip.distance", table.Distance, accessor).
//	    WithDecimals(3) // Override default 2 decimals
func (fb *FieldBuilder) WithDecimals(decimals int) *FieldBuilder {
	// Store decimals for later use
	fb.base.decimals = decimals

	// Recreate formatter based on field type hint
	switch fb.base.fieldTypeHint {
	case float:
		fb.base.defaultFormatter = createFloatFormatter(decimals)
	case distanceHint:
		fb.base.defaultFormatter = createDistanceFormatter(decimals)
	case pressureHint:
		fb.base.defaultFormatter = createPressureFormatter(decimals)
	case speedHint:
		fb.base.defaultFormatter = createSpeedFormatter(decimals)
	case text2Float:
		fb.base.defaultFormatter = createText2FloatFormatter(decimals)
	case text2Distance:
		fb.base.defaultFormatter = createText2DistanceFormatter(decimals)
	case text2Speed:
		fb.base.defaultFormatter = createText2SpeedFormatter(decimals)
	case floatN:
		fb.base.defaultFormatter = createFloatNFormatter(decimals)
	case distanceN:
		fb.base.defaultFormatter = createDistanceNFormatter(decimals)
	case speedN:
		fb.base.defaultFormatter = createSpeedNFormatter(decimals)
	}
	return fb
}

// WithBoolText sets the true/false text for boolean fields.
//
// Example:
//
//	builder.Field("active", "device.active", table.Bool, accessor).
//	    WithBoolText("Yes", "No")
func (fb *FieldBuilder) WithBoolText(trueText, falseText string) *FieldBuilder {
	fb.base.boolTrueText = trueText
	fb.base.boolFalseText = falseText
	fb.base.defaultFormatter = createBoolFormatter(trueText, falseText)
	return fb
}

// WithSearch sets searchable flag
func (fb *FieldBuilder) WithSearch(search bool) *FieldBuilder {
	fb.base.search = search
	return fb
}

// WithSort sets sortable flag
func (fb *FieldBuilder) WithSort(sort bool) *FieldBuilder {
	fb.base.sort = sort
	return fb
}

// WithSticky makes column sticky (fixed during scroll)
func (fb *FieldBuilder) WithSticky(sticky bool) *FieldBuilder {
	fb.base.sticky = sticky
	return fb
}

// WithHeader sets custom header text
func (fb *FieldBuilder) WithHeader(header string) *FieldBuilder {
	fb.base.header = &header
	return fb
}

// WithHeaderSpan sets header column span
func (fb *FieldBuilder) WithHeaderSpan(span int) *FieldBuilder {
	fb.base.headerSpan = &span
	return fb
}

// WithColumnOrder sets column ordering
func (fb *FieldBuilder) WithColumnOrder(order int) *FieldBuilder {
	fb.base.columnOrder = order
	return fb
}

// WithInputType sets input field type
func (fb *FieldBuilder) WithInputType(inputType string) *FieldBuilder {
	fb.base.inputType = &inputType
	return fb
}

// WithInputRequired sets input required flag
func (fb *FieldBuilder) WithInputRequired(required bool) *FieldBuilder {
	fb.base.inputRequired = &required
	return fb
}

// WithInputLang sets input language
func (fb *FieldBuilder) WithInputLang(lang string) *FieldBuilder {
	fb.base.inputLang = &lang
	return fb
}

// WithInputPaste sets input paste enabled
func (fb *FieldBuilder) WithInputPaste(paste bool) *FieldBuilder {
	fb.base.inputPaste = &paste
	return fb
}

// WithTextPrefix sets text prefix
func (fb *FieldBuilder) WithTextPrefix(prefix string) *FieldBuilder {
	fb.base.textPrefix = &prefix
	return fb
}

// WithTextSuffix sets text suffix
func (fb *FieldBuilder) WithTextSuffix(suffix string) *FieldBuilder {
	fb.base.textSuffix = &suffix
	return fb
}

// WithEditable marks the field as inline-editable.
// When true, the Angular frontend will allow clicking on the cell to edit the value.
// Requires editUrl to be set on the table options.
func (fb *FieldBuilder) WithEditable(editable bool) *FieldBuilder {
	fb.base.editable = editable
	return fb
}

// WithEditableOptionsUrl marks the field as inline-editable with options loaded from a URL.
// The Angular frontend will fetch GET {url}?id={rowId}&field={fieldId} when editing starts.
// Automatically sets editable=true.
func (fb *FieldBuilder) WithEditableOptionsUrl(url string) *FieldBuilder {
	fb.base.editable = true
	fb.base.editableOptionsUrl = url
	return fb
}

// WithEditableOptions marks the field as inline-editable with a select dropdown.
// The options define the allowed values shown in the dropdown.
// Automatically sets editable=true.
//
// Example:
//
//	builder.TextField("status", "Status", statusAccessor).
//	    WithEditableOptions(map[string]string{
//	        "Active":       "Active",
//	        "Discontinued": "Discontinued",
//	        "On Sale":      "On Sale",
//	    })
func (fb *FieldBuilder) WithEditableOptions(options map[string]string) *FieldBuilder {
	fb.base.editable = true
	fb.base.editableOptions = make([]editableOption, 0, len(options))
	for value, label := range options {
		fb.base.editableOptions = append(fb.base.editableOptions, editableOption{value: value, label: label})
	}
	return fb
}

// EditableChipOption defines a selectable chip option for inline editing with color.
type EditableChipOption struct {
	Value string
	Label string
	Color core.Color
}

// WithEditableChipOptions marks the field as inline-editable with a multi-select chip dropdown.
// Each option defines a chip value, display label, and optional color.
// Automatically sets editable=true.
//
// Example:
//
//	builder.ChipsField("tags", "Tags", tagsAccessor).
//	    WithEditableChipOptions([]table.EditableChipOption{
//	        {Value: "Frontend", Label: "Frontend", Color: core.ColorPrimary},
//	        {Value: "Backend", Label: "Backend", Color: core.ColorAccent},
//	    })
func (fb *FieldBuilder) WithEditableChipOptions(options []EditableChipOption) *FieldBuilder {
	fb.base.editable = true
	fb.base.editableOptions = make([]editableOption, len(options))
	for i, o := range options {
		fb.base.editableOptions[i] = editableOption{value: o.Value, label: o.Label, color: o.Color}
	}
	return fb
}

// WithEditableOptionsSearch adds a search box to the inline-edit select for
// client-side filtering of the option list. The options themselves still come
// from WithEditableOptions (static) or WithEditableOptionsUrl (loaded per row);
// this only enables the search input and local filtering by label.
// Automatically sets editable=true.
//
// For server-side search use WithEditableSearchOptionsUrl instead.
func (fb *FieldBuilder) WithEditableOptionsSearch(enable bool) *FieldBuilder {
	fb.base.editable = true
	fb.base.editableOptionsSearch = enable
	return fb
}

// WithEditableSearchOptionsUrl marks the field as inline-editable with a
// server-side searchable select. The Angular frontend shows a search box and,
// on each (debounced) keystroke, POSTs {id, field, search} to url and renders
// the returned []{value, label, color?} as options. Combine with
// WithEditableOptions to seed the dropdown with initial options shown before
// the user types.
// Automatically sets editable=true and enables the search box.
//
// Example:
//
//	builder.TextField("category", "Category", categoryAccessor).
//	    WithEditableSearchOptionsUrl("/api/categories/search")
func (fb *FieldBuilder) WithEditableSearchOptionsUrl(url string) *FieldBuilder {
	fb.base.editable = true
	fb.base.editableSearchUrl = url
	fb.base.editableOptionsSearch = true
	return fb
}

// WithAccess sets required permissions
func (fb *FieldBuilder) WithAccess(access []string) *FieldBuilder {
	fb.base.access = access
	return fb
}

// AddMenuItem adds a menu item definition to the last added menu button.
func (fb *FieldBuilder) AddMenuItem(action FieldButtonAction, icon string, color core.Color, text string) *FieldBuilder {
	if fb.base.menuItems == nil {
		fb.base.menuItems = make(map[int][]*menuItemDef)
	}
	fb.base.menuItems[fb.lastMenuKey] = append(fb.base.menuItems[fb.lastMenuKey], &menuItemDef{
		action: action,
		icon:   icon,
		color:  color,
		text:   text,
	})
	return fb
}

// ============================================================================
// T-dependent methods — generic standalone functions
// ============================================================================

// WithRowHint sets a per-row hint accessor for icon fields.
// When set, the hint text is extracted from each row and sent as row[fieldId + "Hint"].
// In the Angular frontend, this overrides the static per-icon hint.
func WithRowHint[T any](fb *FieldBuilder, accessor func(T) string) *FieldBuilder {
	fb.typedField.(*field[T]).hintAccessor = accessor
	return fb
}

// AddMenu adds a menu trigger button. The accessor provides per-row menu item data.
// Each element in the returned []string corresponds to a menu item:
// - non-empty string: URL/data for the menu item
// - "": hide the menu item for this row
// Returning nil hides the entire menu button for this row.
func AddMenu[T any](fb *FieldBuilder, key int, icon string, color core.Color, hint string, accessor func(T) []string) *FieldBuilder {
	if !fb.base.addButton(key, FieldButtonActionMenu, icon, color, hint) {
		return fb // out-of-range key rejected; do not record parallel menu state
	}

	f := fb.typedField.(*field[T])
	if f.menuAccessors == nil {
		f.menuAccessors = make(map[int]func(T) []string)
	}
	f.menuAccessors[key] = accessor

	if fb.base.menuItems == nil {
		fb.base.menuItems = make(map[int][]*menuItemDef)
	}
	fb.base.menuItems[key] = nil

	fb.lastMenuKey = key
	return fb
}
