// Package tachotime provides the TachoTime tachograph chart component for the Angular frontend.
package tachotime

import "github.com/xiriframework/xiri-go/component/core"

// TachoTime represents a tachograph time chart component.
type TachoTime struct {
	header  string
	display *string
	data    []TachoTimeDay
}

// TachoTimeDay represents one day of tachograph data.
type TachoTimeDay struct {
	date              string
	minDate           int64
	maxDate           int64
	activities        []TachoTimeData
	driveblocks       []TachoTimeDriveBlock
	drivedays         []TachoTimeDriveDay
	durationDriving   int32
	durationWorking   int32
	durationBreak     int32
	durationAvailable int32
	durationTotal     int32
}

// TachoTimeData represents a single activity entry.
type TachoTimeData struct {
	start    int64
	end      int64
	activity int32
	duration string
}

// TachoTimeDriveBlock represents a driving block with metadata.
type TachoTimeDriveBlock struct {
	start  int64
	end    int64
	length int32
	data   TachoTimeDriveBlockData
}

// TachoTimeDriveBlockData contains metadata for a drive block.
type TachoTimeDriveBlockData struct {
	driving  int32
	duration string
	start    string
	end      string
}

// TachoTimeDriveDay represents a driving day with metadata.
type TachoTimeDriveDay struct {
	start   int64
	end     int64
	unknown int32
	data    TachoTimeDriveDayData
}

// TachoTimeDriveDayData contains metadata for a drive day.
type TachoTimeDriveDayData struct {
	duration string
	start    string
	end      string
}

// ============================================================================
// Constructors
// ============================================================================

// NewTachoTime creates a new TachoTime component.
func NewTachoTime(header string, data []TachoTimeDay, display *string) *TachoTime {
	return &TachoTime{
		header:  header,
		display: display,
		data:    data,
	}
}

// NewTachoTimeDay creates a new day entry for a tachograph chart.
func NewTachoTimeDay(
	date string, minDate, maxDate int64,
	activities []TachoTimeData,
	driveblocks []TachoTimeDriveBlock,
	drivedays []TachoTimeDriveDay,
	durationDriving, durationWorking, durationBreak, durationAvailable, durationTotal int32,
) TachoTimeDay {
	return TachoTimeDay{
		date: date, minDate: minDate, maxDate: maxDate,
		activities: activities, driveblocks: driveblocks, drivedays: drivedays,
		durationDriving: durationDriving, durationWorking: durationWorking,
		durationBreak: durationBreak, durationAvailable: durationAvailable,
		durationTotal: durationTotal,
	}
}

// NewTachoTimeData creates a new activity entry.
func NewTachoTimeData(start, end int64, activity int32, duration string) TachoTimeData {
	return TachoTimeData{start: start, end: end, activity: activity, duration: duration}
}

// NewTachoTimeDriveBlock creates a new drive block entry.
func NewTachoTimeDriveBlock(start, end int64, length int32, data TachoTimeDriveBlockData) TachoTimeDriveBlock {
	return TachoTimeDriveBlock{start: start, end: end, length: length, data: data}
}

// NewTachoTimeDriveBlockData creates metadata for a drive block.
func NewTachoTimeDriveBlockData(driving int32, duration, start, end string) TachoTimeDriveBlockData {
	return TachoTimeDriveBlockData{driving: driving, duration: duration, start: start, end: end}
}

// NewTachoTimeDriveDay creates a new drive day entry.
func NewTachoTimeDriveDay(start, end int64, unknown int32, data TachoTimeDriveDayData) TachoTimeDriveDay {
	return TachoTimeDriveDay{start: start, end: end, unknown: unknown, data: data}
}

// NewTachoTimeDriveDayData creates metadata for a drive day.
func NewTachoTimeDriveDayData(duration, start, end string) TachoTimeDriveDayData {
	return TachoTimeDriveDayData{duration: duration, start: start, end: end}
}

// ============================================================================
// Print
// ============================================================================

// convertTachoTimeData converts TachoTimeData struct to array format [start, end, activity, duration]
func convertTachoTimeData(data TachoTimeData) []interface{} {
	return []interface{}{data.start, data.end, data.activity, data.duration}
}

// convertTachoTimeDriveBlock converts TachoTimeDriveBlock struct to array format [start, end, length, data]
func convertTachoTimeDriveBlock(block TachoTimeDriveBlock) []interface{} {
	return []interface{}{block.start, block.end, block.length, map[string]any{
		"driving":  block.data.driving,
		"duration": block.data.duration,
		"start":    block.data.start,
		"end":      block.data.end,
	}}
}

// convertTachoTimeDriveDay converts TachoTimeDriveDay struct to array format [start, end, unknown, data]
func convertTachoTimeDriveDay(day TachoTimeDriveDay) []interface{} {
	return []interface{}{day.start, day.end, day.unknown, map[string]any{
		"duration": day.data.duration,
		"start":    day.data.start,
		"end":      day.data.end,
	}}
}

// Print implements core.Component. Returns the JSON representation of the tachograph chart.
func (t *TachoTime) Print(translator core.TranslateFunc) map[string]any {
	// Convert each TachoTimeDay to use array format for activities, driveblocks, and drivedays
	convertedDays := make([]map[string]any, len(t.data))

	for i, day := range t.data {
		// Convert activities to arrays
		activities := make([]interface{}, len(day.activities))
		for j, act := range day.activities {
			activities[j] = convertTachoTimeData(act)
		}

		// Convert driveblocks to arrays
		driveblocks := make([]interface{}, len(day.driveblocks))
		for j, block := range day.driveblocks {
			driveblocks[j] = convertTachoTimeDriveBlock(block)
		}

		// Convert drivedays to arrays
		drivedays := make([]interface{}, len(day.drivedays))
		for j, dd := range day.drivedays {
			drivedays[j] = convertTachoTimeDriveDay(dd)
		}

		// Build the converted day object
		convertedDays[i] = map[string]any{
			"date":              day.date,
			"minDate":           day.minDate,
			"maxDate":           day.maxDate,
			"data":              activities,
			"driveblocks":       driveblocks,
			"drivedays":         drivedays,
			"durationDriving":   day.durationDriving,
			"durationWorking":   day.durationWorking,
			"durationBreak":     day.durationBreak,
			"durationAvailable": day.durationAvailable,
			"durationTotal":     day.durationTotal,
		}
	}

	return map[string]any{
		"type":    "tachotime",
		"display": t.display,
		"data": map[string]any{
			"header": t.header,
			"data":   convertedDays,
		},
	}
}
