package table

import (
	"fmt"
	"strconv"
)

// Row provides access to all field values in a table row.
// Used by formatters for cross-field dependencies (e.g., distance formatter needs device_id to determine units).
type Row interface {
	// Get returns the raw value by field ID
	Get(fieldID string) any

	// GetInt32 returns the value as int32 (with type conversion)
	GetInt32(fieldID string) int32

	// GetInt64 returns the value as int64 (with type conversion)
	GetInt64(fieldID string) int64

	// GetFloat64 returns the value as float64 (with type conversion)
	GetFloat64(fieldID string) float64

	// GetString returns the value as string (with type conversion)
	GetString(fieldID string) string

	// GetBool returns the value as bool (with type conversion)
	GetBool(fieldID string) bool
}

// TypedRow wraps a struct value and provides field access via accessor functions.
// This is the concrete implementation of the Row interface used by Table[T].
type TypedRow[T any] struct {
	data     T
	fieldMap map[string]func(T) any
}

// NewTypedRow creates a new TypedRow with the given data and field accessors
func NewTypedRow[T any](data T, fieldMap map[string]func(T) any) *TypedRow[T] {
	return &TypedRow[T]{
		data:     data,
		fieldMap: fieldMap,
	}
}

// Get returns the raw value by field ID
func (r *TypedRow[T]) Get(fieldID string) any {
	if accessor, ok := r.fieldMap[fieldID]; ok {
		return accessor(r.data)
	}
	return nil
}

// GetInt32 returns the value as int32 (with type conversion)
func (r *TypedRow[T]) GetInt32(fieldID string) int32 {
	v := r.Get(fieldID)
	if v == nil {
		return 0
	}

	switch val := v.(type) {
	case int32:
		return val
	case int:
		return int32(val)
	case int64:
		return int32(val)
	case float32:
		return int32(val)
	case float64:
		return int32(val)
	case string:
		if i, err := strconv.ParseInt(val, 10, 32); err == nil {
			return int32(i)
		}
	}

	return 0
}

// GetInt64 returns the value as int64 (with type conversion)
func (r *TypedRow[T]) GetInt64(fieldID string) int64 {
	v := r.Get(fieldID)
	if v == nil {
		return 0
	}

	switch val := v.(type) {
	case int64:
		return val
	case int:
		return int64(val)
	case int32:
		return int64(val)
	case float32:
		return int64(val)
	case float64:
		return int64(val)
	case string:
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			return i
		}
	}

	return 0
}

// GetFloat64 returns the value as float64 (with type conversion)
func (r *TypedRow[T]) GetFloat64(fieldID string) float64 {
	v := r.Get(fieldID)
	if v == nil {
		return 0.0
	}

	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int32:
		return float64(val)
	case int64:
		return float64(val)
	case string:
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
	}

	return 0.0
}

// GetString returns the value as string (with type conversion)
func (r *TypedRow[T]) GetString(fieldID string) string {
	v := r.Get(fieldID)
	if v == nil {
		return ""
	}

	return fmt.Sprint(v)
}

// GetBool returns the value as bool (with type conversion)
func (r *TypedRow[T]) GetBool(fieldID string) bool {
	v := r.Get(fieldID)
	if v == nil {
		return false
	}

	switch val := v.(type) {
	case bool:
		return val
	case int:
		return val != 0
	case int32:
		return val != 0
	case int64:
		return val != 0
	case string:
		return val == "true" || val == "1" || val == "yes"
	}

	return false
}
