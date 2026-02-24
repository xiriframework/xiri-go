package field

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ============================================================================
// Helper Functions
// ============================================================================

// parseDateTime parses a date/time from various formats:
// - Unix timestamp (int, int64, float64)
// - ISO date string ("2006-01-02")
// - ISO datetime string ("2006-01-02T15:04:05Z")
func parseDateTime(raw interface{}) (time.Time, error) {
	if raw == nil {
		return time.Time{}, fmt.Errorf("date value is nil")
	}

	switch v := raw.(type) {
	case string:
		// Try ISO date format first (2006-01-02)
		if t, err := time.Parse("2006-01-02", v); err == nil {
			return t, nil
		}
		// Try ISO datetime format (2006-01-02T15:04:05Z)
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			return t, nil
		}
		return time.Time{}, fmt.Errorf("invalid date string format: %s", v)

	case int:
		// Unix timestamp (seconds)
		return time.Unix(int64(v), 0), nil

	case int32:
		// Unix timestamp (seconds)
		return time.Unix(int64(v), 0), nil

	case int64:
		// Unix timestamp (seconds)
		return time.Unix(v, 0), nil

	case float64:
		// Unix timestamp (seconds, from JSON number)
		return time.Unix(int64(v), 0), nil

	default:
		return time.Time{}, fmt.Errorf("unsupported date type: %T", v)
	}
}

// ModelListValue represents a list of selected model IDs
type ModelListValue []int32

// parseModelListValue parses ModelListValue from various formats
func parseModelListValue(raw interface{}, defaultValue interface{}) (ModelListValue, error) {
	if raw == nil {
		if defaultValue != nil {
			return defaultValue.(ModelListValue), nil
		}
		return ModelListValue{}, nil
	}

	// Handle different input formats
	switch v := raw.(type) {
	case []interface{}:
		// Array of numbers
		result := make(ModelListValue, 0, len(v))
		for _, item := range v {
			switch id := item.(type) {
			case float64:
				result = append(result, int32(id))
			case int:
				result = append(result, int32(id))
			case int32:
				result = append(result, id)
			case int64:
				result = append(result, int32(id))
			case string:
				// Try to parse string to int
				parsed, err := strconv.ParseInt(id, 10, 32)
				if err != nil {
					return nil, fmt.Errorf("invalid model ID: %s", id)
				}
				result = append(result, int32(parsed))
			default:
				return nil, fmt.Errorf("invalid model ID type: %T", item)
			}
		}
		return result, nil

	case string:
		// Comma-separated string
		if v == "" {
			return ModelListValue{}, nil
		}
		parts := strings.Split(v, ",")
		result := make(ModelListValue, 0, len(parts))
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			parsed, err := strconv.ParseInt(part, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid model ID: %s", part)
			}
			result = append(result, int32(parsed))
		}
		return result, nil

	case ModelListValue:
		// Already the correct type
		return v, nil

	case []int32:
		// Direct []int32 slice (underlying type of ModelListValue)
		return v, nil

	default:
		return nil, fmt.Errorf("unsupported modellist format: %T", raw)
	}
}
