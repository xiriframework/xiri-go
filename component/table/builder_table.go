package table

import (
	"log/slog"

	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/form/group"
	"github.com/xiriframework/xiri-go/uicontext"
)

// TableBuilder provides a fluent API for building type-safe tables.
type TableBuilder[T any] struct {
	table *Table[T]
}

// NewBuilder creates a new table builder with the given context and translator.
//
// Parameters:
//   - ctx: User context with locale/language/timezone/unit preferences
//   - translator: Translation function for field labels
//
// Example:
//
//	builder := table.NewBuilder[DeviceRow](ctx, translator)
func NewBuilder[T any](ctx *uicontext.UiContext, translator core.TranslateFunc) *TableBuilder[T] {
	// Set common defaults for table options
	defaultTrue := true
	defaultTextNoData := translator("KEINEDATEN")

	return &TableBuilder[T]{
		table: &Table[T]{
			ctx:             ctx,
			translator:      translator,
			fields:          make([]*Field[T], 0),
			fieldsCanChange: false,
			options: TableOptions{
				Pagination: &defaultTrue, // Default: enabled
				Search:     &defaultTrue, // Default: enabled
				Reload:     &defaultTrue, // Default: enabled
				Csv:        &defaultTrue, // Default: enabled
				Excel:      &defaultTrue, // Default: enabled
				TextNoData: &defaultTextNoData,
			},
			outputType: OutputWeb, // Default output mode
		},
	}
}

// fieldInternal is the internal implementation for all typed field methods.
// This method is used by all type-safe field methods (IntField, TextField, etc.)
// to create fields with the correct configuration.
func (b *TableBuilder[T]) fieldInternal(id string, name string, fieldType FieldTypeHint, accessor func(T) any) *FieldBuilder[T] {
	field := &Field[T]{
		id:         id,
		name:       name,
		accessor:   accessor,
		formatters: make(map[OutputType]OutputFormatter),
		csv:        true,          // Include in CSV by default
		footer:     FieldFooterNo, // No footer by default
	}

	// Apply defaults based on field type
	builder := &FieldBuilder[T]{field: field}
	builder = applyFieldTypeDefaults(builder, fieldType)

	b.table.fields = append(b.table.fields, field)

	return builder
}

// SetFilter sets the filter FormGroup.
func (b *TableBuilder[T]) SetFilter(fg *group.FormGroup) *TableBuilder[T] {
	b.table.filter = fg
	return b
}

// SetHasFilter explicitly sets the hasFilter flag.
// Use this when the table receives filter data from a parent Query component
// rather than having its own filter form.
func (b *TableBuilder[T]) SetHasFilter(hasFilter bool) *TableBuilder[T] {
	b.table.hasFilter = &hasFilter
	return b
}

// SetFlags sets UI-only filter fields that should be excluded from parsed data.
// Flags are typically used for frontend state that shouldn't be sent to the backend.
func (b *TableBuilder[T]) SetFlags(flags ...string) *TableBuilder[T] {
	b.table.flags = flags
	return b
}

// SetFieldsCanChange marks that fields may change between page load and data load.
func (b *TableBuilder[T]) SetFieldsCanChange() *TableBuilder[T] {
	b.table.fieldsCanChange = true
	return b
}

// Build returns the final Table[T].
// It validates field configurations and logs warnings for common mistakes.
func (b *TableBuilder[T]) Build() *Table[T] {
	b.validateFields()
	return b.table
}

// validateFields checks field configurations and logs warnings for likely mistakes.
func (b *TableBuilder[T]) validateFields() {
	for _, f := range b.table.fields {
		switch f.fieldType {
		case FieldTypeIcon:
			if len(f.icons) == 0 {
				slog.Warn("table.Build: icon field has no icon definitions, use IconFieldFromSet()", "fieldId", f.id)
			}
		case FieldTypeButtons:
			if len(f.buttons) == 0 {
				slog.Warn("table.Build: buttons field has no button definitions, use AddButton()", "fieldId", f.id)
			}
		}
	}
}
