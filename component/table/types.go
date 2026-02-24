package table

// FieldTypeHint provides semantic field type constants for the Field() method.
// These constants automatically configure the field with appropriate formatters and settings.
type FieldTypeHint string

// Semantic field type constants
// These are used as the 3rd parameter to Field() to automatically apply default
// formatters and field type settings.
const (
	// Integer creates an int32/int64 field with locale-aware number formatting.
	// - Sets field type to FieldTypeNumber
	// - Web/PDF output: locale-aware formatting (e.g., "1,234" or "1.234")
	// - CSV/Excel output: plain number string
	// - Returns [display, value] array on web output for sorting
	Integer FieldTypeHint = "integer"

	// Float creates a float64 field with locale-aware decimal formatting.
	// - Sets field type to FieldTypeNumber
	// - Default decimals: 2 (override with .WithDecimals(n))
	// - Web/PDF output: locale-aware formatting (e.g., "123,45" or "123.45")
	// - CSV/Excel output: formatted number string
	// - Returns [display, value] array on web output for sorting
	Float FieldTypeHint = "float"

	// Text creates a string field with simple text formatting.
	// - Sets field type to FieldTypeText
	// - All outputs: plain string or empty string for nil
	Text FieldTypeHint = "text"

	// Bool creates a boolean field with custom true/false text mapping.
	// - Sets field type to FieldTypeText
	// - Requires .WithBoolText(trueText, falseText) to set display values
	// - Default: shows "true"/"false" if WithBoolText not called
	Bool FieldTypeHint = "bool"

	// DateTime creates a timestamp field with date+time formatting.
	// - Sets field type to FieldTypeText
	// - Expects int64 Unix timestamp (seconds)
	// - Web/PDF output: user timezone and locale (e.g., "2021-12-20 12:26")
	// - CSV/Excel output: ISO format "2006-01-02 15:04:05" in user timezone
	DateTime FieldTypeHint = "datetime"

	// Date creates a timestamp field with date-only formatting (no time component).
	// - Sets field type to FieldTypeText
	// - Expects int64 Unix timestamp (seconds)
	// - Web/PDF output: user timezone and locale (e.g., "2021-12-20")
	// - CSV/Excel output: ISO format "2006-01-02"
	Date FieldTypeHint = "date"

	// Distance creates a float64 field with automatic unit conversion (km/mi/NM).
	// - Sets field type to FieldTypeNumber
	// - Expects value in kilometers
	// - Default decimals: 2 (override with .WithDecimals(n))
	// - Automatically converts to miles or nautical miles based on user/device preference
	// - Web/PDF output: formatted with unit (e.g., "123,45 km" or "76.72 mi")
	// - CSV/Excel output: converted numeric value only
	Distance FieldTypeHint = "distance"

	// Pressure creates a float64 field with automatic unit conversion (bar/psi).
	// - Sets field type to FieldTypeNumber
	// - Expects value in bar
	// - Default decimals: 2 (override with .WithDecimals(n))
	// - Automatically converts to psi based on user/device preference
	// - Web/PDF output: formatted with unit (e.g., "2,50 bar" or "36.26 psi")
	// - CSV/Excel output: converted numeric value only
	Pressure FieldTypeHint = "pressure"

	// Speed creates a float64 field with automatic unit conversion (km/h, mph, knots).
	// - Sets field type to FieldTypeNumber
	// - Expects value in km/h
	// - Default decimals: 1 (override with .WithDecimals(n))
	// - Automatically converts to mph or knots based on user/device preference
	// - Web/PDF output: formatted with unit (e.g., "100,0 km/h" or "62.1 mph")
	// - CSV/Excel output: converted numeric value only
	Speed FieldTypeHint = "speed"

	// Buttons creates a buttons-type field with action buttons.
	// - Sets field type to FieldTypeButtons
	// - Used for row actions (edit, delete, view, download, etc.)
	// - Configure buttons with .AddButton(key, action, icon, color, hint)
	// - Each button can have custom actions: link, dialog, api, etc.
	// - All outputs: renders as button array in table cell
	Buttons FieldTypeHint = "buttons"

	// Icon creates an icon-type field with value-to-icon mappings.
	// - Sets field type to FieldTypeIcon
	// - Maps field values to specific icons with colors
	// - Configure with IconFieldFromSet()
	// - Useful for status indicators, priorities, categories
	// - All outputs: renders as colored icon in table cell
	Icon FieldTypeHint = "icon"

	// Link creates a link-type field with clickable URLs.
	// - Sets field type to FieldTypeLink
	// - Renders values as clickable hyperlinks
	// - Supports internal navigation and external URLs
	// - All outputs: clickable link in web, plain URL in CSV/Excel
	Link FieldTypeHint = "link"

	// Html creates an html-type field with raw HTML content.
	// - Sets field type to FieldTypeHtml
	// - Renders raw HTML in table cells
	// - Web/PDF output: rendered HTML
	// - CSV/Excel output: plain text (HTML tags stripped)
	Html FieldTypeHint = "html"

	// Input creates an input-type field with inline editing.
	// - Sets field type to FieldTypeInput
	// - Allows direct cell editing in table
	// - Supports text, number, select input types
	// - Web output: editable input field
	// - CSV/Excel output: current value
	Input FieldTypeHint = "input"

	// Text2 creates a text2-type field (alternative text style).
	// - Sets field type to FieldTypeText2
	// - Alternative text rendering with different styling
	// - Used for secondary information or subtitles
	// - All outputs: plain text with alternative formatting
	Text2 FieldTypeHint = "text2"

	// Text2Int creates a text2-type field with integer formatting.
	// - Sets field type to FieldTypeText2
	// - Expects [2]int accessor
	// - Web/PDF output: locale-aware number formatting for both lines
	// - CSV/Excel output: formatted numbers combined with " - "
	Text2Int FieldTypeHint = "text2int"

	// Text2Float creates a text2-type field with float formatting.
	// - Sets field type to FieldTypeText2
	// - Expects [2]float64 accessor
	// - Default decimals: 2 (override with .WithDecimals(n))
	// - Web/PDF output: locale-aware decimal formatting for both lines
	// - CSV/Excel output: formatted numbers combined with " - "
	Text2Float FieldTypeHint = "text2float"

	// Text2DateTime creates a text2-type field with datetime formatting.
	// - Sets field type to FieldTypeText2
	// - Expects [2]time.Time accessor
	// - Web/PDF output: user timezone and locale for both lines
	// - CSV/Excel output: ISO datetime format combined with " - "
	Text2DateTime FieldTypeHint = "text2datetime"

	// Text2Date creates a text2-type field with date-only formatting.
	// - Sets field type to FieldTypeText2
	// - Expects [2]time.Time accessor
	// - Web/PDF output: date-only format in user timezone for both lines
	// - CSV/Excel output: ISO date format combined with " - "
	Text2Date FieldTypeHint = "text2date"

	// Text2Distance creates a text2-type field with distance formatting and unit conversion.
	// - Sets field type to FieldTypeText2
	// - Expects [2]float64 accessor (values in kilometers)
	// - Default decimals: 2 (override with .WithDecimals(n))
	// - Automatically converts to miles or nautical miles based on user/device preference
	// - Web/PDF output: formatted with unit for both lines
	// - CSV/Excel output: converted numeric values combined with " - "
	Text2Distance FieldTypeHint = "text2distance"

	// Text2Speed creates a text2-type field with speed formatting and unit conversion.
	// - Sets field type to FieldTypeText2
	// - Expects [2]float64 accessor (values in km/h)
	// - Default decimals: 1 (override with .WithDecimals(n))
	// - Automatically converts to mph or knots based on user/device preference
	// - Web/PDF output: formatted with unit for both lines
	// - CSV/Excel output: converted numeric values combined with " - "
	Text2Speed FieldTypeHint = "text2speed"

	// Text2Bool creates a text2-type field with boolean formatting.
	// - Sets field type to FieldTypeText2
	// - Expects [2]bool accessor
	// - Web/PDF output: "Yes"/"No" for both lines (locale-aware)
	// - CSV/Excel output: "Yes"/"No" combined with " - "
	Text2Bool FieldTypeHint = "text2bool"

	// TimeLength creates a time duration field with HH:MM formatting.
	// - Sets field type to FieldTypeText
	// - Expects int64 accessor (value in seconds)
	// - Web/PDF output: "HH:MM" format (e.g., "05:30"), or "Xd HH:MM" for >= 24 hours (e.g., "2d 05:30")
	// - CSV/Excel output: integer minutes (no decimals)
	TimeLength FieldTypeHint = "timelength"

	// Text2TimeLength creates a text2-type field with time duration formatting.
	// - Sets field type to FieldTypeText2
	// - Expects [2]int64 accessor (values in seconds)
	// - Web/PDF output: "HH:MM" or "Xd HH:MM" format for both lines
	// - CSV/Excel output: integer minutes combined with " - "
	Text2TimeLength FieldTypeHint = "text2timelength"

	// Header creates a header-type field (section divider).
	// - Sets field type to FieldTypeHeader
	// - Used to group related columns visually
	// - Not searchable or sortable
	// - Web output: renders as table section header
	// - CSV/Excel output: column label
	Header FieldTypeHint = "header"

	// Id creates an id-type field with special export format.
	// - Sets field type to FieldTypeID
	// - Uses number formatting with no decimals
	// - Special "id" type for frontend identification
	// - Expects int64 value
	// - All outputs: formatted as integer with special "id" type marker
	Id FieldTypeHint = "id"
)

// ============================================================================
// Field Configuration Enums
// ============================================================================

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

// FieldType represents the type/format of a table field for JSON export
type FieldType string

const (
	FieldTypeText    FieldType = "text"
	FieldTypeButtons FieldType = "buttons"
	FieldTypeIcon    FieldType = "icon"
	FieldTypeHtml    FieldType = "html"
	FieldTypeLink    FieldType = "link"
	FieldTypeInput   FieldType = "input"
	FieldTypeText2   FieldType = "text2"
	FieldTypeHeader  FieldType = "header"
	FieldTypeNumber  FieldType = "number"
	FieldTypeID      FieldType = "id" // Special ID field type
)

// FieldColor represents theme colors for field components
type FieldColor string

const (
	FieldColorPrimary  FieldColor = "primary"
	FieldColorAccent   FieldColor = "accent"
	FieldColorWarning  FieldColor = "warn"
	FieldColorTertiary FieldColor = "tertiary"
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
