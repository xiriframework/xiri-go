package table

// FieldBuilder provides a fluent API for configuring a single field
type FieldBuilder[T any] struct {
	field       *Field[T]
	lastMenuKey int // Track last menu button index for AddMenuItem chaining
}

// WithFormatter sets the default formatter (used for all output types unless overridden)
func (fb *FieldBuilder[T]) WithFormatter(formatter OutputFormatter) *FieldBuilder[T] {
	fb.field.defaultFormatter = formatter
	return fb
}

// WithWebFormatter sets formatter specifically for web output (overrides default)
func (fb *FieldBuilder[T]) WithWebFormatter(formatter OutputFormatter) *FieldBuilder[T] {
	fb.field.formatters[OutputWeb] = formatter
	return fb
}

// WithCSVFormatter sets formatter specifically for CSV output (overrides default)
func (fb *FieldBuilder[T]) WithCSVFormatter(formatter OutputFormatter) *FieldBuilder[T] {
	fb.field.formatters[OutputCSV] = formatter
	return fb
}

// WithPDFFormatter sets formatter specifically for PDF output (overrides default)
func (fb *FieldBuilder[T]) WithPDFFormatter(formatter OutputFormatter) *FieldBuilder[T] {
	fb.field.formatters[OutputPDF] = formatter
	return fb
}

// WithExcelFormatter sets formatter specifically for Excel output (overrides default)
func (fb *FieldBuilder[T]) WithExcelFormatter(formatter OutputFormatter) *FieldBuilder[T] {
	fb.field.formatters[OutputExcel] = formatter
	return fb
}

// WithFooter sets footer aggregation type
func (fb *FieldBuilder[T]) WithFooter(footer FieldFooter) *FieldBuilder[T] {
	fb.field.footer = footer
	return fb
}

// WithFooterSum enables sum aggregation in footer
func (fb *FieldBuilder[T]) WithFooterSum() *FieldBuilder[T] {
	fb.field.footer = FieldFooterSum
	return fb
}

// WithFooterCount enables count aggregation in footer
func (fb *FieldBuilder[T]) WithFooterCount() *FieldBuilder[T] {
	fb.field.footer = FieldFooterCount
	return fb
}

// Hide hides the field (not displayed in any output)
func (fb *FieldBuilder[T]) Hide() *FieldBuilder[T] {
	fb.field.hide = true
	return fb
}

// HideInCSV excludes field from CSV export
func (fb *FieldBuilder[T]) HideInCSV() *FieldBuilder[T] {
	fb.field.csv = false
	return fb
}

// ShowInCSV includes field in CSV export (default)
func (fb *FieldBuilder[T]) ShowInCSV() *FieldBuilder[T] {
	fb.field.csv = true
	return fb
}

// WithAlign sets text alignment
func (fb *FieldBuilder[T]) WithAlign(align FieldAlign) *FieldBuilder[T] {
	fb.field.align = &align
	return fb
}

// AlignLeft sets left alignment
func (fb *FieldBuilder[T]) AlignLeft() *FieldBuilder[T] {
	align := FieldAlignLeft
	fb.field.align = &align
	return fb
}

// AlignCenter sets center alignment
func (fb *FieldBuilder[T]) AlignCenter() *FieldBuilder[T] {
	align := FieldAlignCenter
	fb.field.align = &align
	return fb
}

// AlignRight sets right alignment
func (fb *FieldBuilder[T]) AlignRight() *FieldBuilder[T] {
	align := FieldAlignRight
	fb.field.align = &align
	return fb
}

// WithWidth sets column width
func (fb *FieldBuilder[T]) WithWidth(width string) *FieldBuilder[T] {
	fb.field.width = &width
	return fb
}

// WithMinWidth sets minimum column width
func (fb *FieldBuilder[T]) WithMinWidth(minWidth string) *FieldBuilder[T] {
	fb.field.minWidth = &minWidth
	return fb
}

// WithHint sets tooltip/hint text
func (fb *FieldBuilder[T]) WithHint(hint string) *FieldBuilder[T] {
	fb.field.hint = &hint
	return fb
}

// WithDisplay sets CSS display class
func (fb *FieldBuilder[T]) WithDisplay(display string) *FieldBuilder[T] {
	fb.field.display = &display
	return fb
}

// AddButton adds a button to a buttons-type field
func (fb *FieldBuilder[T]) AddButton(
	key int,
	action FieldButtonAction,
	icon string,
	color FieldColor,
	hint string,
) *FieldBuilder[T] {
	fb.field.AddButton(key, action, icon, color, hint)
	return fb
}

// WithDecimals sets the number of decimal places for numeric fields.
// This is used with Float, Distance, Pressure, Speed, and Text2 numeric field types.
//
// Example:
//
//	builder.Field("distance", "trip.distance", table.Distance, accessor).
//	    WithDecimals(3) // Override default 2 decimals
func (fb *FieldBuilder[T]) WithDecimals(decimals int) *FieldBuilder[T] {
	// Store decimals for later use
	fb.field.decimals = decimals

	// Recreate formatter based on field type hint
	switch fb.field.fieldTypeHint {
	case Float:
		fb.field.defaultFormatter = createFloatFormatter(decimals)
	case Distance:
		fb.field.defaultFormatter = createDistanceFormatter(decimals)
	case Pressure:
		fb.field.defaultFormatter = createPressureFormatter(decimals)
	case Speed:
		fb.field.defaultFormatter = createSpeedFormatter(decimals)
	case Text2Float:
		fb.field.defaultFormatter = createText2FloatFormatter(decimals)
	case Text2Distance:
		fb.field.defaultFormatter = createText2DistanceFormatter(decimals)
	case Text2Speed:
		fb.field.defaultFormatter = createText2SpeedFormatter(decimals)
	}
	return fb
}

// WithBoolText sets the true/false text for boolean fields.
//
// Example:
//
//	builder.Field("active", "device.active", table.Bool, accessor).
//	    WithBoolText("Yes", "No")
func (fb *FieldBuilder[T]) WithBoolText(trueText, falseText string) *FieldBuilder[T] {
	fb.field.boolTrueText = trueText
	fb.field.boolFalseText = falseText
	fb.field.defaultFormatter = createBoolFormatter(trueText, falseText)
	return fb
}

// WithSearch sets searchable flag
func (fb *FieldBuilder[T]) WithSearch(search bool) *FieldBuilder[T] {
	fb.field.search = search
	return fb
}

// WithSort sets sortable flag
func (fb *FieldBuilder[T]) WithSort(sort bool) *FieldBuilder[T] {
	fb.field.sort = sort
	return fb
}

// WithSticky makes column sticky (fixed during scroll)
func (fb *FieldBuilder[T]) WithSticky(sticky bool) *FieldBuilder[T] {
	fb.field.sticky = sticky
	return fb
}

// WithHeader sets custom header text
func (fb *FieldBuilder[T]) WithHeader(header string) *FieldBuilder[T] {
	fb.field.header = &header
	return fb
}

// WithHeaderSpan sets header column span
func (fb *FieldBuilder[T]) WithHeaderSpan(span int) *FieldBuilder[T] {
	fb.field.headerSpan = &span
	return fb
}

// WithColumnOrder sets column ordering
func (fb *FieldBuilder[T]) WithColumnOrder(order int) *FieldBuilder[T] {
	fb.field.columnOrder = order
	return fb
}

// WithInputType sets input field type
func (fb *FieldBuilder[T]) WithInputType(inputType string) *FieldBuilder[T] {
	fb.field.inputType = &inputType
	return fb
}

// WithInputRequired sets input required flag
func (fb *FieldBuilder[T]) WithInputRequired(required bool) *FieldBuilder[T] {
	fb.field.inputRequired = &required
	return fb
}

// WithInputLang sets input language
func (fb *FieldBuilder[T]) WithInputLang(lang string) *FieldBuilder[T] {
	fb.field.inputLang = &lang
	return fb
}

// WithInputPaste sets input paste enabled
func (fb *FieldBuilder[T]) WithInputPaste(paste bool) *FieldBuilder[T] {
	fb.field.inputPaste = &paste
	return fb
}

// WithTextPrefix sets text prefix
func (fb *FieldBuilder[T]) WithTextPrefix(prefix string) *FieldBuilder[T] {
	fb.field.textPrefix = &prefix
	return fb
}

// WithTextSuffix sets text suffix
func (fb *FieldBuilder[T]) WithTextSuffix(suffix string) *FieldBuilder[T] {
	fb.field.textSuffix = &suffix
	return fb
}

// WithRowHint sets a per-row hint accessor for icon fields.
// When set, the hint text is extracted from each row and sent as row[fieldId + "Hint"].
// In the Angular frontend, this overrides the static per-icon hint.
func (fb *FieldBuilder[T]) WithRowHint(accessor func(T) string) *FieldBuilder[T] {
	fb.field.hintAccessor = accessor
	return fb
}

// WithAccess sets required permissions
func (fb *FieldBuilder[T]) WithAccess(access []string) *FieldBuilder[T] {
	fb.field.access = access
	return fb
}

// AddMenu adds a menu trigger button. The accessor provides per-row menu item data.
// Each element in the returned []string corresponds to a menu item:
// - non-empty string: URL/data for the menu item
// - "": hide the menu item for this row
// Returning nil hides the entire menu button for this row.
func (fb *FieldBuilder[T]) AddMenu(key int, icon string, color FieldColor, hint string, accessor func(T) []string) *FieldBuilder[T] {
	fb.field.AddButton(key, FieldButtonActionMenu, icon, color, hint)

	if fb.field.menuAccessors == nil {
		fb.field.menuAccessors = make(map[int]func(T) []string)
	}
	fb.field.menuAccessors[key] = accessor

	if fb.field.menuItems == nil {
		fb.field.menuItems = make(map[int][]*MenuItemDef)
	}
	fb.field.menuItems[key] = nil

	fb.lastMenuKey = key
	return fb
}

// AddMenuItem adds a menu item definition to the last added menu button.
func (fb *FieldBuilder[T]) AddMenuItem(action FieldButtonAction, icon string, color FieldColor, text string) *FieldBuilder[T] {
	if fb.field.menuItems == nil {
		fb.field.menuItems = make(map[int][]*MenuItemDef)
	}
	fb.field.menuItems[fb.lastMenuKey] = append(fb.field.menuItems[fb.lastMenuKey], &MenuItemDef{
		Action: action,
		Icon:   icon,
		Color:  color,
		Text:   text,
	})
	return fb
}
