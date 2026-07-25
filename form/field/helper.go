package field

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/xiriframework/xiri-go/component/core"
)

// ============================================================================
// Helper Functions
// ============================================================================

// toInt32 converts numeric and string input to int32 without loss. Fractional
// numbers, NaN/Inf and values outside the int32 range are rejected instead of
// being truncated or wrapped — silently turning 1.9 into ID 1, or 3000000000
// into -1294967296, selects the wrong record.
func toInt32(raw interface{}) (int32, error) {
	switch v := raw.(type) {
	case int32:
		return v, nil
	case int:
		return toInt32(int64(v))
	case int64:
		if v > math.MaxInt32 || v < math.MinInt32 {
			return 0, fmt.Errorf("value %d out of int32 range", v)
		}
		return int32(v), nil
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return 0, fmt.Errorf("value %v is not a number", v)
		}
		if math.Trunc(v) != v {
			return 0, fmt.Errorf("value %v is not an integer", v)
		}
		if v > math.MaxInt32 || v < math.MinInt32 {
			return 0, fmt.Errorf("value %v out of int32 range", v)
		}
		return int32(v), nil
	case string:
		parsed, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return 0, fmt.Errorf("value %q is not an int32", v)
		}
		return int32(parsed), nil
	default:
		return 0, fmt.Errorf("unsupported int32 value type: %T", raw)
	}
}

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

// dayOffsetLimit is the threshold that separates a relative day offset from an
// absolute Unix timestamp in Min/Max bounds.
const dayOffsetLimit = 10000

// resolveDateBound converts a Min/Max bound into an absolute Unix timestamp.
// Values within ±dayOffsetLimit are day offsets, anchored to local midnight in the
// user's timezone; larger values are already absolute timestamps and pass through.
//
// Anchoring to local midnight and stepping with AddDate keeps the boundary correct
// across DST transitions, where a day is not 86400 seconds long.
func resolveDateBound(ctx *core.UiContext, bound int64, now time.Time) int64 {
	if bound <= -dayOffsetLimit || bound >= dayOffsetLimit {
		return bound
	}

	loc, err := time.LoadLocation(ctx.SafeTimezone().GetIANA())
	if err != nil {
		loc = time.UTC
	}

	local := now.In(loc)
	midnight := time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, loc)
	return midnight.AddDate(0, 0, int(bound)).Unix()
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
			id, err := toInt32(item)
			if err != nil {
				return nil, fmt.Errorf("invalid model ID %v: %w", item, err)
			}
			result = append(result, id)
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
			id, err := toInt32(part)
			if err != nil {
				return nil, fmt.Errorf("invalid model ID %q: %w", part, err)
			}
			result = append(result, id)
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
