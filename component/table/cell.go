package table

import (
	"math"
	"time"

	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/formatter"
)

// cellValueFor returns v of the web cell object for a field type hint, and whether the hint delivers
// cell objects at all. v is the raw value in a form that both sorts and edits: ISO date / local ISO
// datetime (user's timezone) as string, seconds or numbers as numbers; for text2/textN the first
// value. nil means "empty" (zero time, negative duration, empty slice, unexpected value type).
// Number hints are not listed: they ship the legacy [display, value] pair.
//
// ponytail: numeric v travels as a JSON number, so integers beyond 2^53 tie in the browser while
// the display d keeps full precision. Upgrade path: zero-padded string v.
func cellValueFor(hint fieldTypeHint, value any, ctx *core.UiContext) (v any, ok bool) {
	switch hint {
	case date:
		return isoOrNil(toInt64(value), formatter.CellDateLayout, ctx), true
	case dateTime:
		return isoOrNil(toInt64(value), formatter.CellDateTimeLayout, ctx), true
	case timeLength:
		return nonNegative(toInt64(value)), true
	case text2Date, text2DateTime:
		if arr, ok := value.([2]time.Time); ok {
			return isoTimeOrNil(arr[0], layoutFor(hint), ctx), true
		}
		return nil, true
	case text2Int:
		if arr, ok := value.([2]int); ok {
			return int64(arr[0]), true
		}
		return nil, true
	case text2Float, text2Distance, text2Speed:
		if arr, ok := value.([2]float64); ok {
			return finiteOrNil(arr[0]), true
		}
		return nil, true
	case text2TimeLength:
		if arr, ok := value.([2]int64); ok {
			return nonNegative(arr[0]), true
		}
		return nil, true
	case dateN, dateTimeN:
		if arr, ok := value.([]time.Time); ok && len(arr) > 0 {
			return isoTimeOrNil(arr[0], layoutFor(hint), ctx), true
		}
		return nil, true
	case integerN:
		if arr, ok := value.([]int); ok && len(arr) > 0 {
			return int64(arr[0]), true
		}
		return nil, true
	case floatN, distanceN, speedN:
		if arr, ok := value.([]float64); ok && len(arr) > 0 {
			return finiteOrNil(arr[0]), true
		}
		return nil, true
	case timeLengthN:
		if arr, ok := value.([]int64); ok && len(arr) > 0 {
			return nonNegative(arr[0]), true
		}
		return nil, true
	}
	return nil, false
}

// cellObjectKind is the field JSON value of "cellObject": "" for plain cells, otherwise the kind of v
// ("string" for ISO dates/times, "number" for seconds and numbers). The client derives the edit value
// type from it, so an empty v or a custom inputType cannot change the type of a column.
func cellObjectKind(hint fieldTypeHint) string {
	if _, ok := cellValueFor(hint, nil, nil); !ok {
		return ""
	}
	switch hint {
	case date, dateTime, text2Date, text2DateTime, dateN, dateTimeN:
		return "string"
	}
	return "number"
}

// layoutFor picks the date or datetime layout for the date-like text2/textN hints.
func layoutFor(hint fieldTypeHint) string {
	if hint == text2DateTime || hint == dateTimeN {
		return formatter.CellDateTimeLayout
	}
	return formatter.CellDateLayout
}

// ctxLocation is the user's timezone as *time.Location, UTC when it cannot be loaded (ctx may be nil).
func ctxLocation(ctx *core.UiContext) *time.Location {
	loc, err := time.LoadLocation(ctx.SafeTimezone().GetIANA())
	if err != nil {
		return time.UTC
	}
	return loc
}

func isoOrNil(ts int64, layout string, ctx *core.UiContext) any {
	if ts == 0 {
		return nil
	}
	return time.Unix(ts, 0).In(ctxLocation(ctx)).Format(layout)
}

func isoTimeOrNil(t time.Time, layout string, ctx *core.UiContext) any {
	if t.IsZero() {
		return nil
	}
	return t.In(ctxLocation(ctx)).Format(layout)
}

func nonNegative(v int64) any {
	if v < 0 {
		return nil
	}
	return v
}

// finiteOrNil keeps NaN and ±Inf out of v: encoding/json refuses them and would fail the whole response.
func finiteOrNil(f float64) any {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return nil
	}
	return f
}
