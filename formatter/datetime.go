package formatter

import (
	"fmt"
	"time"

	"github.com/xiriframework/xiri-go/types/locale"
	"github.com/xiriframework/xiri-go/uicontext"
)

// ToUnixTimestamp converts Go time.Time to Unix timestamp in SECONDS
// CRITICAL: Must return seconds, NOT milliseconds
// xiri-ui expects Unix timestamps in seconds, not ISO strings
func ToUnixTimestamp(t time.Time) int64 {
	return t.Unix() // NOT t.UnixMilli() / 1000
}

// ToUnixTimestampBigInt converts Go time.Time to Unix timestamp in MILLISECONDS
func ToUnixTimestampBigInt(t time.Time) int64 {
	return t.UnixMilli()
}

// FromUnixTimestamp converts Unix timestamp (seconds) to Go time.Time
func FromUnixTimestamp(ts int64) time.Time {
	return time.Unix(ts, 0) // NOT time.UnixMilli(ts)
}

// FromUnixTimestampBigInt converts Unix timestamp in milliseconds (int64) to Go time.Time
func FromUnixTimestampBigInt(ts int64) time.Time {
	return time.UnixMilli(ts)
}

// FormatTimestampToTextRange formats a Unix timestamp to a relative text range
// e.g., "gerade eben", "vor 2 min", "vor 3 h", "vor 2 d", "02.01.2006 15:04"
// Parameters:
//   - timestamp: Unix timestamp in seconds
//   - includeTime: If true, shows time for dates older than 7 days
//   - timezone: IANA timezone string (e.g., "Europe/Vienna")
//   - translate: Optional translation function for locale-aware units (variadic for backward compatibility)
func FormatTimestampToTextRange(timestamp int64, includeTime bool, timezone string, translate ...func(string) string) string {
	if timestamp == 0 {
		return "-"
	}

	// Translation function - use provided or default to returning key as-is
	trans := func(key string) string { return key }
	if len(translate) > 0 && translate[0] != nil {
		trans = translate[0]
	}

	// Load timezone
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.UTC
	}

	// Convert timestamp to time
	t := FromUnixTimestamp(timestamp).In(loc)
	now := time.Now().In(loc)

	// Calculate difference
	diff := now.Sub(t)

	// Format based on difference (locale-aware)
	if diff < time.Minute {
		return trans("T.JETZT") // "jetzt" (DE), "now" (EN), etc.
	} else if diff < time.Hour {
		minutes := int(diff.Minutes())
		// Use T.VOR for past ("ago"), T.SEIT would be for duration ("since")
		return fmt.Sprintf("%s %d %s", trans("T.VOR"), minutes, trans("T.MIN"))
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		return fmt.Sprintf("%s %d %s", trans("T.VOR"), hours, trans("T.HOUR"))
	} else if diff < 7*24*time.Hour {
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%s %d d", trans("T.VOR"), days)
	}

	// For older timestamps, show formatted date/time
	if includeTime {
		return t.Format("2006-01-02 15:04")
	}
	return t.Format("2006-01-02")
}

// FormatTimestampDateTime formats a Unix timestamp to date and time without seconds
func FormatTimestampDateTime(timestamp int64, ctx *uicontext.UiContext) string {
	if timestamp == 0 {
		return "-"
	}

	return FormatDateTime(FromUnixTimestamp(timestamp), ctx)
}

// FormatTimestampDate formats a Unix timestamp to date
func FormatTimestampDate(timestamp int64, ctx *uicontext.UiContext) string {
	if timestamp == 0 {
		return "-"
	}

	return FormatDate(FromUnixTimestamp(timestamp), ctx)
}

// dateLayout returns the Go date format string for a locale
func dateLayout(loc locale.Locale) string {
	switch loc {
	case locale.De, locale.DeAT, locale.DeCH, locale.Sv, locale.Nb, locale.Da, locale.Fi:
		return "2006-01-02" // ISO
	case locale.EnUS:
		return "01/02/2006" // US MDY
	case locale.Ja, locale.ZhCN:
		return "2006/01/02" // Asian YMD
	default:
		return "02/01/2006" // European DMY
	}
}

// dateTimeLayout returns the Go datetime format string for a locale
func dateTimeLayout(loc locale.Locale) string {
	switch loc {
	case locale.De, locale.DeAT, locale.DeCH, locale.Sv, locale.Nb, locale.Da, locale.Fi:
		return "2006-01-02 15:04"
	case locale.EnUS:
		return "01/02/2006 03:04 PM"
	case locale.Ja, locale.ZhCN:
		return "2006/01/02 15:04"
	default:
		return "02/01/2006 15:04"
	}
}

// timeLayout returns the Go time-only format string for a locale
func timeLayout(loc locale.Locale) string {
	switch loc {
	case locale.EnUS, locale.EnGB:
		return "03:04 PM"
	default:
		return "15:04"
	}
}

// FormatDate formats a time.Time to date
func FormatDate(t time.Time, ctx *uicontext.UiContext) string {
	if t.IsZero() {
		return "-"
	}

	loc, err := time.LoadLocation(ctx.Timezone.GetIANA())
	if err != nil {
		loc = time.UTC
	}

	return t.In(loc).Format(dateLayout(ctx.Locale))
}

// FormatDateTime formats a time.Time to datetime
func FormatDateTime(t time.Time, ctx *uicontext.UiContext) string {
	if t.IsZero() {
		return "-"
	}

	loc, err := time.LoadLocation(ctx.Timezone.GetIANA())
	if err != nil {
		loc = time.UTC
	}

	return t.In(loc).Format(dateTimeLayout(ctx.Locale))
}

// FormatTime formats a time.Time to time only
func FormatTime(t time.Time, ctx *uicontext.UiContext) string {
	if t.IsZero() {
		return "-"
	}

	loc, err := time.LoadLocation(ctx.Timezone.GetIANA())
	if err != nil {
		loc = time.UTC
	}

	return t.In(loc).Format(timeLayout(ctx.Locale))
}

// FormatTimestampFullDate formats a Unix timestamp to full date (Y-m-d H:i format)
func FormatTimestampFullDate(timestamp int64, ctx *uicontext.UiContext) string {
	return FormatTimestampDateTime(timestamp, ctx)
}

// FormatMinutesAfterMidnight converts "minutes after midnight" to HH:MM time string
// This is used for tachograph data where times are stored as minutes since midnight (0-1439)
// Parameters:
//   - dayTimestamp: Unix timestamp for the day start (midnight UTC)
//   - minutesAfterMidnight: Minutes since midnight (0-1439)
//   - timezone: IANA timezone string (e.g., "Europe/Vienna")
//
// Example: dayTimestamp=1609459200 (2021-01-01 00:00 UTC), minutesAfterMidnight=930 → "15:30"
func FormatMinutesAfterMidnight(dayTimestamp int32, minutesAfterMidnight int16, timezone string) string {
	if minutesAfterMidnight < 0 || minutesAfterMidnight >= 1440 {
		return "-"
	}

	// Load timezone
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.UTC
	}

	// Convert day timestamp to time in specified timezone
	dayStart := time.Unix(int64(dayTimestamp), 0).In(loc)

	// Add minutes offset
	actualTime := dayStart.Add(time.Duration(minutesAfterMidnight) * time.Minute)

	// Format as HH:MM
	return actualTime.Format("15:04")
}
