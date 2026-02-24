package formatter

import "fmt"

// FormatTimeLimitFromDB formats time limit from database fields
// Parses PostgreSQL boolean array string and time strings to create formatted output
//
// Parameters:
//   - weekdaysStr: PostgreSQL boolean array format "{t,f,t,f,t,f,f}" (Mon-Sun)
//   - timeFrom: Time string in "HH:MM:SS" format (can be nil)
//   - timeTo: Time string in "HH:MM:SS" format (can be nil)
//   - timeIn: If true, time limit applies inside range; if false, outside range (shows "Nicht: ")
//   - t: Translation function for weekday and "Nicht" keys
//
// Returns formatted string like:
//   - "Mo-Fr 08:00-17:00" (time limit applies Mon-Fri 8am-5pm)
//   - "Nicht: Sa,So 10:00-18:00" (time limit applies EXCEPT Sat-Sun 10am-6pm)
//   - "Di,Do 00:00-24:00" (time limit applies Tue and Thu all day)
func FormatTimeLimitFromDB(weekdaysStr string, timeFrom, timeTo *string, timeIn *bool, t func(string) string) string {
	// Parse weekdays from PostgreSQL array format "{t,f,t,f,t,f,f}"
	weekdays := parsePostgresArrayBool(weekdaysStr)

	// Build text with "out" prefix if time_in is false
	text := ""
	if timeIn != nil && !*timeIn {
		text = t("TL.OUTSHORT") + ": "
	}

	// Find which weekdays are selected
	var wdids []int
	for i, val := range weekdays {
		if val {
			wdids = append(wdids, i)
		}
	}

	// Compress consecutive weekdays into ranges
	var ranges []interface{} // Can be int (single day) or [2]int (range)
	var firstval *int
	var lastval *int

	for _, val := range wdids {
		if lastval != nil {
			if *lastval+1 == val {
				// Consecutive day
				if firstval == nil {
					firstval = new(int)
					*firstval = *lastval
				}
			} else {
				// Gap - save previous range or single day
				if firstval != nil {
					ranges = append(ranges, [2]int{*firstval, *lastval})
					firstval = nil
				} else {
					ranges = append(ranges, *lastval)
				}
			}
		}
		lastval = new(int)
		*lastval = val
	}

	// Add final range or single day
	if lastval != nil {
		if firstval != nil {
			ranges = append(ranges, [2]int{*firstval, *lastval})
		} else {
			ranges = append(ranges, *lastval)
		}
	}

	// Format weekdays using TL.WD0-TL.WD6 translation keys
	var wdtext []string
	for _, val := range ranges {
		switch v := val.(type) {
		case [2]int:
			// Range: Mo-Fr
			wdtext = append(wdtext, fmt.Sprintf("%s-%s",
				t(fmt.Sprintf("TL.WD%d", v[0])),
				t(fmt.Sprintf("TL.WD%d", v[1]))))
		case int:
			// Single day: Mo
			wdtext = append(wdtext, t(fmt.Sprintf("TL.WD%d", v)))
		}
	}

	// Join weekdays with comma
	if len(wdtext) > 0 {
		for i, wd := range wdtext {
			if i > 0 {
				text += ","
			}
			text += wd
		}
		text += " " // Add space before time range
	}

	// Add time range (HH:MM-HH:MM)
	if timeFrom != nil && timeTo != nil {
		fromStr := extractTimeString(*timeFrom)
		toStr := extractTimeString(*timeTo)

		// Only append if both times are valid
		if fromStr != "" && toStr != "" {
			text += fromStr + "-" + toStr
		}
	}

	return text
}

// extractTimeString extracts HH:MM from various time string formats
// Handles:
//   - "HH:MM:SS" -> "HH:MM"
//   - "HH:MM:SS.microseconds" -> "HH:MM"
//   - "0000-01-01T00:00:00Z" (ISO datetime from PostgreSQL TIME) -> "00:00"
//   - "0000-01-02T00:00:00Z" (PostgreSQL 24:00:00 becomes next day) -> "24:00"
func extractTimeString(timeStr string) string {
	if len(timeStr) < 5 {
		return ""
	}

	// Check if it's ISO datetime format (contains 'T')
	// Format: "0000-01-01T00:00:00Z" or "0000-01-02T00:00:00Z"
	if len(timeStr) >= 16 {
		tIndex := -1
		for i, ch := range timeStr {
			if ch == 'T' {
				tIndex = i
				break
			}
		}
		if tIndex >= 0 && tIndex+6 <= len(timeStr) {
			// Check if this is "0000-01-02T00:00:00Z" which represents 24:00:00
			// PostgreSQL stores 24:00:00 as midnight of next day
			if len(timeStr) >= 20 && timeStr[:10] == "0000-01-02" {
				// Extract time portion after 'T'
				timePart := timeStr[tIndex+1:]
				// If it's 00:00:xx on day 02, it's actually 24:00
				if len(timePart) >= 5 && timePart[:5] == "00:00" {
					return "24:00"
				}
			}

			// Extract time portion after 'T': "00:00:00" -> "00:00"
			timePart := timeStr[tIndex+1:]
			if len(timePart) >= 5 && timePart[2] == ':' {
				return timePart[:5]
			}
		}
	}

	// Standard time format: "HH:MM:SS" or "HH:MM:SS.microseconds"
	if len(timeStr) >= 5 && timeStr[2] == ':' {
		return timeStr[:5]
	}

	return ""
}

// parsePostgresArrayBool parses PostgreSQL boolean array string "{t,f,t,f,t,f,f}" to []bool
// PostgreSQL stores boolean arrays as text in format: {true,false,...} or {t,f,...}
func parsePostgresArrayBool(arr string) []bool {
	// Handle empty or invalid input
	if len(arr) < 2 {
		return []bool{false, false, false, false, false, false, false}
	}

	// Remove braces
	arr = arr[1 : len(arr)-1]

	// Split by comma and extract 't' or 'f' characters
	var parts []rune
	for _, c := range arr {
		if c == 't' || c == 'f' {
			parts = append(parts, c)
		}
	}

	// Convert to bool array
	result := make([]bool, len(parts))
	for i, c := range parts {
		result[i] = c == 't'
	}

	return result
}
