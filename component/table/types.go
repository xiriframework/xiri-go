package table

// fieldTypeHint provides semantic field type constants for typed field methods.
// These constants automatically configure the field with appropriate formatters and settings.
type fieldTypeHint string

// Semantic field type constants.
// These are used internally by typed field methods to automatically apply
// default formatters and field type settings.
const (
	// integer creates an int32/int64 field with locale-aware number formatting.
	integer fieldTypeHint = "integer"

	// float creates a float64 field with locale-aware decimal formatting.
	float fieldTypeHint = "float"

	// text creates a string field with simple text formatting.
	text fieldTypeHint = "text"

	// boolean creates a boolean field with custom true/false text mapping.
	boolean fieldTypeHint = "bool"

	// dateTime creates a timestamp field with date+time formatting.
	dateTime fieldTypeHint = "datetime"

	// date creates a timestamp field with date-only formatting.
	date fieldTypeHint = "date"

	// distanceHint creates a float64 field with automatic unit conversion (km/mi/NM).
	distanceHint fieldTypeHint = "distance"

	// pressureHint creates a float64 field with automatic unit conversion (bar/psi).
	pressureHint fieldTypeHint = "pressure"

	// speedHint creates a float64 field with automatic unit conversion (km/h, mph, knots).
	speedHint fieldTypeHint = "speed"

	// buttons creates a buttons-type field with action buttons.
	buttons fieldTypeHint = "buttons"

	// icon creates an icon-type field with value-to-icon mappings.
	icon fieldTypeHint = "icon"

	// link creates a link-type field with clickable URLs.
	link fieldTypeHint = "link"

	// html creates an html-type field with raw HTML content.
	html fieldTypeHint = "html"

	// input creates an input-type field with inline editing.
	inputHint fieldTypeHint = "input"

	// text2 creates a text2-type field (alternative text style).
	text2 fieldTypeHint = "text2"

	// text2Int creates a text2-type field with integer formatting.
	text2Int fieldTypeHint = "text2int"

	// text2Float creates a text2-type field with float formatting.
	text2Float fieldTypeHint = "text2float"

	// text2DateTime creates a text2-type field with datetime formatting.
	text2DateTime fieldTypeHint = "text2datetime"

	// text2Date creates a text2-type field with date-only formatting.
	text2Date fieldTypeHint = "text2date"

	// text2Distance creates a text2-type field with distance formatting and unit conversion.
	text2Distance fieldTypeHint = "text2distance"

	// text2Speed creates a text2-type field with speed formatting and unit conversion.
	text2Speed fieldTypeHint = "text2speed"

	// text2Bool creates a text2-type field with boolean formatting.
	text2Bool fieldTypeHint = "text2bool"

	// timeLength creates a time duration field with HH:MM formatting.
	timeLength fieldTypeHint = "timelength"

	// text2TimeLength creates a text2-type field with time duration formatting.
	text2TimeLength fieldTypeHint = "text2timelength"

	// textN creates a textn-type field with variable number of lines.
	textN fieldTypeHint = "textn"

	// integerN creates a textn-type field with integer formatting for variable lines.
	integerN fieldTypeHint = "integern"

	// floatN creates a textn-type field with float formatting for variable lines.
	floatN fieldTypeHint = "floatn"

	// dateTimeN creates a textn-type field with datetime formatting for variable lines.
	dateTimeN fieldTypeHint = "datetimen"

	// dateN creates a textn-type field with date-only formatting for variable lines.
	dateN fieldTypeHint = "daten"

	// distanceN creates a textn-type field with distance formatting for variable lines.
	distanceN fieldTypeHint = "distancen"

	// speedN creates a textn-type field with speed formatting for variable lines.
	speedN fieldTypeHint = "speedn"

	// boolN creates a textn-type field with boolean formatting for variable lines.
	boolN fieldTypeHint = "booln"

	// timeLengthN creates a textn-type field with time duration formatting for variable lines.
	timeLengthN fieldTypeHint = "timelenthn"

	// header creates a header-type field (section divider).
	header fieldTypeHint = "header"

	// idHint creates an id-type field with special export format.
	idHint fieldTypeHint = "id"

	// chipsHint creates a chips-type field rendering a list of colored chips per row.
	chipsHint fieldTypeHint = "chips"
)

// ============================================================================
// Field Configuration Enums
// ============================================================================

// Density controls the row height of a table.
//
// The frontend accepts these three values; the older SetDense(true) is only a
// legacy alias for DensityCompact and cannot express DensityRelaxed.
type Density string

const (
	DensityCompact Density = "compact"
	DensityRegular Density = "regular"
	DensityRelaxed Density = "relaxed"
)

// FieldAlign represents field text alignment
type FieldAlign string

const (
	FieldAlignLeft   FieldAlign = "left"
	FieldAlignCenter FieldAlign = "center"
	FieldAlignRight  FieldAlign = "right"
)

// FieldFooter represents footer aggregation types
type FieldFooter string

const (
	FieldFooterNo     FieldFooter = "no"
	FieldFooterSum    FieldFooter = "sum"
	FieldFooterCount  FieldFooter = "count"
	FieldFooterStatic FieldFooter = "static"
)

// fieldType represents the type/format of a table field for JSON export
type fieldType string

const (
	fieldTypeText    fieldType = "text"
	fieldTypeButtons fieldType = "buttons"
	fieldTypeIcon    fieldType = "icon"
	fieldTypeHtml    fieldType = "html"
	fieldTypeLink    fieldType = "link"
	fieldTypeInput   fieldType = "input"
	fieldTypeText2   fieldType = "text2"
	fieldTypeTextN   fieldType = "textn"
	fieldTypeHeader  fieldType = "header"
	fieldTypeNumber  fieldType = "number"
	fieldTypeID      fieldType = "id"
	fieldTypeChips   fieldType = "chips"
)

// FieldButtonAction represents button action types
type FieldButtonAction string

const (
	FieldButtonActionLink     FieldButtonAction = "link"
	FieldButtonActionDialog   FieldButtonAction = "dialog"
	FieldButtonActionApi      FieldButtonAction = "api"
	FieldButtonActionDownload FieldButtonAction = "download"
	FieldButtonActionForm     FieldButtonAction = "form"
	FieldButtonActionBack     FieldButtonAction = "back"
	FieldButtonActionClose    FieldButtonAction = "close"
	FieldButtonActionSave     FieldButtonAction = "save"
	FieldButtonActionHref     FieldButtonAction = "href"
	FieldButtonActionGet      FieldButtonAction = "get"
	FieldButtonActionPost     FieldButtonAction = "post"
	FieldButtonActionPut      FieldButtonAction = "put"
	FieldButtonActionDelete   FieldButtonAction = "delete"
	FieldButtonActionMenu     FieldButtonAction = "menu"
)
