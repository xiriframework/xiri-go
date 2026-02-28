package table

// FieldMeta provides read-only metadata about a table field for external consumers
// (e.g., PDF/Excel/CSV renderers) without exposing the internal field[T] type.
type FieldMeta struct {
	ID    string
	Name  string      // Translation key
	Type  string      // "text", "number", "text2", "id", "html", "link", "icon", "buttons", "textn", "header", "input"
	Hidden bool
	Align *FieldAlign
	CSV   bool // Whether field is included in CSV export
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

// GetFieldMetas returns metadata for all fields. This is the public API for
// external consumers that need field information (renderers, exporters).
func (t *Table[T]) GetFieldMetas() []FieldMeta {
	metas := make([]FieldMeta, len(t.fields))
	for i, f := range t.fields {
		metas[i] = FieldMeta{
			ID:     f.id,
			Name:   f.name,
			Type:   string(f.fieldType),
			Hidden: f.hide,
			Align:  f.align,
			CSV:    f.csv,
		}
	}
	return metas
}
