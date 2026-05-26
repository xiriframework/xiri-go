package table

import (
	"log/slog"

	"github.com/xiriframework/xiri-go/form/group"
)

// TableBuilder provides a fluent API for building type-safe tables.
type TableBuilder[T any] struct {
	table *Table[T]
}

// NewBuilder creates a new table builder.
//
// Example:
//
//	builder := table.NewBuilder[DeviceRow]()
func NewBuilder[T any]() *TableBuilder[T] {
	// Set common defaults for table options
	defaultTrue := true
	defaultTextNoData := "KEINEDATEN" // Translation key, resolved lazily at Print time

	return &TableBuilder[T]{
		table: &Table[T]{
			tableCore: tableCore{
				fieldBases: make([]*fieldBase, 0),
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
			fields: make([]*field[T], 0),
		},
	}
}

// fieldInternal is the internal implementation for all typed field methods.
// This method is used by all type-safe field methods (IntField, TextField, etc.)
// to create fields with the correct configuration.
func (b *TableBuilder[T]) fieldInternal(id string, name string, fieldType fieldTypeHint, accessor func(T) any) *FieldBuilder {
	f := &field[T]{
		fieldBase: fieldBase{
			id:         id,
			name:       name,
			formatters: make(map[OutputType]OutputFormatter),
			csv:        true,          // Include in CSV by default
			footer:     FieldFooterNo, // No footer by default
		},
		accessor: accessor,
	}

	// Apply defaults based on field type (non-generic)
	applyFieldTypeDefaults(&f.fieldBase, fieldType)

	b.table.fields = append(b.table.fields, f)
	b.table.fieldBases = append(b.table.fieldBases, &f.fieldBase)

	return &FieldBuilder{base: &f.fieldBase, typedField: f}
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
	for _, f := range b.table.fieldBases {
		switch f.fieldType {
		case fieldTypeIcon:
			if len(f.icons) == 0 {
				slog.Warn("table.Build: icon field has no icon definitions, use IconFieldFromSet()", "fieldId", f.id)
			}
		case fieldTypeButtons:
			if len(f.buttons) == 0 {
				slog.Warn("table.Build: buttons field has no button definitions, use AddButton()", "fieldId", f.id)
			}
		}
	}
	b.validateTree()
}

// validateTree warns about common tree-mode configuration mistakes:
// referencing fields that don't exist, or a parentId field that is hidden (which would be
// dropped from the row data and break the hierarchy).
func (b *TableBuilder[T]) validateTree() {
	if b.table.options.Tree == nil {
		return
	}
	tree := b.table.options.Tree

	var idField, parentField *fieldBase
	for _, f := range b.table.fieldBases {
		if f.id == tree.IdField {
			idField = f
		}
		if f.id == tree.ParentIdField {
			parentField = f
		}
	}

	if tree.IdField == "" || idField == nil {
		slog.Warn("table.Build: tree mode IdField not found among table fields", "idField", tree.IdField)
	}
	if tree.ParentIdField == "" || parentField == nil {
		slog.Warn("table.Build: tree mode ParentIdField not found among table fields", "parentIdField", tree.ParentIdField)
	} else if parentField.hide {
		slog.Warn("table.Build: tree mode ParentIdField is hidden and will be dropped from row data; use an id-format field instead", "parentIdField", tree.ParentIdField)
	}
}
