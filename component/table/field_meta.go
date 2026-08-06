package table

// FieldMeta provides read-only metadata about a table field for external consumers
// (e.g., PDF/Excel/CSV renderers) without exposing the internal field[T] type.
type FieldMeta struct {
	ID     string
	Name   string      // Translation key
	Type   string      // "text", "number", "text2", "id", "html", "link", "icon", "buttons", "textn", "header", "input"
	Hidden bool
	Align  *FieldAlign
	CSV    bool // Whether field is included in CSV export

	Header     *string // Custom header text override (grouped headers)
	HeaderSpan *int    // Header column span for grouped headers
}

// FieldMeta type constants for use by renderers.
const (
	FieldMetaTypeText    = "text"
	FieldMetaTypeText2   = "text2"
	FieldMetaTypeNumber  = "number"
	FieldMetaTypeID      = "id"
	FieldMetaTypeHtml    = "html"
	FieldMetaTypeLink    = "link"
	FieldMetaTypeIcon    = "icon"
	FieldMetaTypeButtons = "buttons"
	FieldMetaTypeTextN   = "textn"
	FieldMetaTypeHeader  = "header"
	FieldMetaTypeInput   = "input"
)

// GetFieldMetas is now defined on tableCore in table.go
