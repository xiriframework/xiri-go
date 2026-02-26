package table

import "time"

// ============================================================================
// Type-Safe Field Methods
// ============================================================================
//
// These methods provide compile-time type checking for field accessors.
// Each method enforces the correct return type for its field type, preventing
// runtime errors from type mismatches.

// IdField adds an ID field with special export type "id".
// Returns int64 value with special "id" type for frontend.
//
// Example:
//
//	builder.IdField("id", "device.id", func(r DeviceRow) int64 {
//	    return r.ID
//	})
func (b *TableBuilder[T]) IdField(id, name string, accessor func(T) int64) *FieldBuilder[T] {
	return b.fieldInternal(id, name, Id, func(row T) any {
		return accessor(row)
	})
}

// IntField adds an integer field.
// Returns formatted number with locale-aware thousands separator.
//
// Example:
//
//	builder.IntField("count", "device.count", func(r DeviceRow) int {
//	    return r.Count
//	})
func (b *TableBuilder[T]) IntField(id, name string, accessor func(T) int) *FieldBuilder[T] {
	return b.fieldInternal(id, name, Integer, func(row T) any {
		return accessor(row)
	})
}

// Int32Field adds an int32 field.
// Returns formatted number with locale-aware thousands separator.
//
// Example:
//
//	builder.Int32Field("value", "sensor.value", func(r SensorRow) int32 {
//	    return r.Value
//	})
func (b *TableBuilder[T]) Int32Field(id, name string, accessor func(T) int32) *FieldBuilder[T] {
	return b.fieldInternal(id, name, Integer, func(row T) any {
		return accessor(row)
	})
}

// Int64Field adds an int64 field.
// Returns formatted number with locale-aware thousands separator.
//
// Example:
//
//	builder.Int64Field("size", "file.size", func(r FileRow) int64 {
//	    return r.Size
//	})
func (b *TableBuilder[T]) Int64Field(id, name string, accessor func(T) int64) *FieldBuilder[T] {
	return b.fieldInternal(id, name, Integer, func(row T) any {
		return accessor(row)
	})
}

// FloatField adds a float64 field with decimal formatting.
// Default decimals: 2 (override with .WithDecimals(n))
//
// Example:
//
//	builder.FloatField("price", "product.price", func(r Product) float64 {
//	    return r.Price
//	}).WithDecimals(2)
func (b *TableBuilder[T]) FloatField(id, name string, accessor func(T) float64) *FieldBuilder[T] {
	return b.fieldInternal(id, name, Float, func(row T) any {
		return accessor(row)
	})
}

// TextField adds a string field.
//
// Example:
//
//	builder.TextField("name", "device.name", func(r DeviceRow) string {
//	    return r.Name
//	})
func (b *TableBuilder[T]) TextField(id, name string, accessor func(T) string) *FieldBuilder[T] {
	return b.fieldInternal(id, name, Text, func(row T) any {
		return accessor(row)
	})
}

// Deprecated: Use [TableBuilder.TextNField] instead.
//
// Text2Field adds a two-line text field.
// Accessor returns [2]string: [0] = primary text, [1] = secondary text
//
// Example:
//
//	builder.Text2Field("device", "device.info", func(r DeviceRow) [2]string {
//	    return [2]string{r.Name, r.Model}
//	})
func (b *TableBuilder[T]) Text2Field(id, name string, accessor func(T) [2]string) *FieldBuilder[T] {
	return b.fieldInternal(id, name, Text2, func(row T) any {
		return accessor(row)
	})
}

// Deprecated: Use [TableBuilder.IntNField] instead.
//
// Text2IntField adds a two-line integer field with locale-aware formatting.
// Accessor returns [2]int: [0] = primary value, [1] = secondary value
//
// Example:
//
//	builder.Text2IntField("stats", "device.stats", func(r DeviceRow) [2]int {
//	    return [2]int{r.TripsToday, r.TotalTrips}
//	})
func (b *TableBuilder[T]) Text2IntField(id, name string, accessor func(T) [2]int) *FieldBuilder[T] {
	return b.fieldInternal(id, name, Text2Int, func(row T) any {
		return accessor(row)
	})
}

// Deprecated: Use [TableBuilder.FloatNField] instead.
//
// Text2FloatField adds a two-line float field with locale-aware decimal formatting.
// Accessor returns [2]float64: [0] = primary value, [1] = secondary value
// Default decimals: 2 (override with .WithDecimals(n))
//
// Example:
//
//	builder.Text2FloatField("fuel", "device.fuel", func(r DeviceRow) [2]float64 {
//	    return [2]float64{r.FuelCurrent, r.FuelAverage}
//	}).WithDecimals(2)
func (b *TableBuilder[T]) Text2FloatField(id, name string, accessor func(T) [2]float64) *FieldBuilder[T] {
	return b.fieldInternal(id, name, Text2Float, func(row T) any {
		return accessor(row)
	})
}

// Deprecated: Use [TableBuilder.DateTimeNField] instead.
//
// Text2DateTimeField adds a two-line datetime field with timezone-aware formatting.
// Accessor returns [2]time.Time: [0] = primary time, [1] = secondary time
//
// Example:
//
//	builder.Text2DateTimeField("times", "device.times", func(r DeviceRow) [2]time.Time {
//	    return [2]time.Time{r.LastSeen, r.Created}
//	})
func (b *TableBuilder[T]) Text2DateTimeField(id, name string, accessor func(T) [2]time.Time) *FieldBuilder[T] {
	return b.fieldInternal(id, name, Text2DateTime, func(row T) any {
		return accessor(row)
	})
}

// Deprecated: Use [TableBuilder.DateNField] instead.
//
// Text2DateField adds a two-line date field with date-only formatting (no time).
// Accessor returns [2]time.Time: [0] = primary date, [1] = secondary date
//
// Example:
//
//	builder.Text2DateField("dates", "device.dates", func(r DeviceRow) [2]time.Time {
//	    return [2]time.Time{r.RegistrationDate, r.ExpiryDate}
//	})
func (b *TableBuilder[T]) Text2DateField(id, name string, accessor func(T) [2]time.Time) *FieldBuilder[T] {
	return b.fieldInternal(id, name, Text2Date, func(row T) any {
		return accessor(row)
	})
}

// Deprecated: Use [TableBuilder.DistanceNField] instead.
//
// Text2DistanceField adds a two-line distance field with automatic unit conversion.
// Accessor returns [2]float64 (values in kilometers): [0] = primary, [1] = secondary
// Default decimals: 2 (override with .WithDecimals(n))
//
// Example:
//
//	builder.Text2DistanceField("distances", "device.distances", func(r DeviceRow) [2]float64 {
//	    return [2]float64{r.TodayKm, r.TotalKm}
//	}).WithDecimals(1)
func (b *TableBuilder[T]) Text2DistanceField(id, name string, accessor func(T) [2]float64) *FieldBuilder[T] {
	return b.fieldInternal(id, name, Text2Distance, func(row T) any {
		return accessor(row)
	})
}

// Deprecated: Use [TableBuilder.SpeedNField] instead.
//
// Text2SpeedField adds a two-line speed field with automatic unit conversion.
// Accessor returns [2]float64 (values in km/h): [0] = primary, [1] = secondary
// Default decimals: 1 (override with .WithDecimals(n))
//
// Example:
//
//	builder.Text2SpeedField("speeds", "device.speeds", func(r DeviceRow) [2]float64 {
//	    return [2]float64{r.MaxSpeed, r.AvgSpeed}
//	}).WithDecimals(1)
func (b *TableBuilder[T]) Text2SpeedField(id, name string, accessor func(T) [2]float64) *FieldBuilder[T] {
	return b.fieldInternal(id, name, Text2Speed, func(row T) any {
		return accessor(row)
	})
}

// Deprecated: Use [TableBuilder.BoolNField] instead.
//
// Text2BoolField adds a two-line boolean field.
// Accessor returns [2]bool: [0] = primary value, [1] = secondary value
// Formatted as "Yes"/"No" for each line.
//
// Example:
//
//	builder.Text2BoolField("status", "device.status", func(r DeviceRow) [2]bool {
//	    return [2]bool{r.Active, r.Online}
//	})
func (b *TableBuilder[T]) Text2BoolField(id, name string, accessor func(T) [2]bool) *FieldBuilder[T] {
	return b.fieldInternal(id, name, Text2Bool, func(row T) any {
		return accessor(row)
	})
}

// TimeLengthField adds a time duration field with HH:MM formatting.
// Accessor returns int64 (value in seconds).
// Web/PDF output: "HH:MM" or "Xd HH:MM" for durations >= 24 hours
// CSV/Excel output: integer minutes
//
// Example:
//
//	builder.TimeLengthField("duration", "trip.duration", func(r TripRow) int64 {
//	    return r.DurationSeconds
//	})
func (b *TableBuilder[T]) TimeLengthField(id, name string, accessor func(T) int64) *FieldBuilder[T] {
	return b.fieldInternal(id, name, TimeLength, func(row T) any {
		return accessor(row)
	})
}

// Deprecated: Use [TableBuilder.TimeLengthNField] instead.
//
// Text2TimeLengthField adds a two-line time duration field.
// Accessor returns [2]int64 (values in seconds): [0] = primary, [1] = secondary
// Web/PDF output: "HH:MM" or "Xd HH:MM" format for both lines
// CSV/Excel output: integer minutes combined with " - "
//
// Example:
//
//	builder.Text2TimeLengthField("times", "trip.times", func(r TripRow) [2]int64 {
//	    return [2]int64{r.DurationSeconds, r.IdleSeconds}
//	})
func (b *TableBuilder[T]) Text2TimeLengthField(id, name string, accessor func(T) [2]int64) *FieldBuilder[T] {
	return b.fieldInternal(id, name, Text2TimeLength, func(row T) any {
		return accessor(row)
	})
}

// TextNField adds a variable-line text field.
// Accessor returns []string: each element is one line of text.
//
// Example:
//
//	builder.TextNField("info", "device.info", func(r DeviceRow) []string {
//	    return []string{r.Name, r.Model, r.Serial}
//	})
func (b *TableBuilder[T]) TextNField(id, name string, accessor func(T) []string) *FieldBuilder[T] {
	return b.fieldInternal(id, name, TextN, func(row T) any {
		return accessor(row)
	})
}

// IntNField adds a variable-line integer field with locale-aware formatting.
// Accessor returns []int: each element is one line value.
//
// Example:
//
//	builder.IntNField("stats", "device.stats", func(r DeviceRow) []int {
//	    return []int{r.TripsToday, r.TripsWeek, r.TotalTrips}
//	})
func (b *TableBuilder[T]) IntNField(id, name string, accessor func(T) []int) *FieldBuilder[T] {
	return b.fieldInternal(id, name, IntegerN, func(row T) any {
		return accessor(row)
	})
}

// FloatNField adds a variable-line float field with locale-aware decimal formatting.
// Accessor returns []float64: each element is one line value.
// Default decimals: 2 (override with .WithDecimals(n))
//
// Example:
//
//	builder.FloatNField("fuel", "device.fuel", func(r DeviceRow) []float64 {
//	    return []float64{r.FuelCurrent, r.FuelAverage, r.FuelMax}
//	}).WithDecimals(2)
func (b *TableBuilder[T]) FloatNField(id, name string, accessor func(T) []float64) *FieldBuilder[T] {
	return b.fieldInternal(id, name, FloatN, func(row T) any {
		return accessor(row)
	})
}

// DateTimeNField adds a variable-line datetime field with timezone-aware formatting.
// Accessor returns []time.Time: each element is one line timestamp.
//
// Example:
//
//	builder.DateTimeNField("times", "device.times", func(r DeviceRow) []time.Time {
//	    return []time.Time{r.LastSeen, r.Created, r.Updated}
//	})
func (b *TableBuilder[T]) DateTimeNField(id, name string, accessor func(T) []time.Time) *FieldBuilder[T] {
	return b.fieldInternal(id, name, DateTimeN, func(row T) any {
		return accessor(row)
	})
}

// DateNField adds a variable-line date field with date-only formatting (no time).
// Accessor returns []time.Time: each element is one line date.
//
// Example:
//
//	builder.DateNField("dates", "device.dates", func(r DeviceRow) []time.Time {
//	    return []time.Time{r.RegistrationDate, r.ExpiryDate, r.InspectionDate}
//	})
func (b *TableBuilder[T]) DateNField(id, name string, accessor func(T) []time.Time) *FieldBuilder[T] {
	return b.fieldInternal(id, name, DateN, func(row T) any {
		return accessor(row)
	})
}

// DistanceNField adds a variable-line distance field with automatic unit conversion.
// Accessor returns []float64 (values in kilometers): each element is one line.
// Default decimals: 2 (override with .WithDecimals(n))
//
// Example:
//
//	builder.DistanceNField("distances", "device.distances", func(r DeviceRow) []float64 {
//	    return []float64{r.TodayKm, r.WeekKm, r.TotalKm}
//	}).WithDecimals(1)
func (b *TableBuilder[T]) DistanceNField(id, name string, accessor func(T) []float64) *FieldBuilder[T] {
	return b.fieldInternal(id, name, DistanceN, func(row T) any {
		return accessor(row)
	})
}

// SpeedNField adds a variable-line speed field with automatic unit conversion.
// Accessor returns []float64 (values in km/h): each element is one line.
// Default decimals: 1 (override with .WithDecimals(n))
//
// Example:
//
//	builder.SpeedNField("speeds", "device.speeds", func(r DeviceRow) []float64 {
//	    return []float64{r.MaxSpeed, r.AvgSpeed, r.MinSpeed}
//	}).WithDecimals(1)
func (b *TableBuilder[T]) SpeedNField(id, name string, accessor func(T) []float64) *FieldBuilder[T] {
	return b.fieldInternal(id, name, SpeedN, func(row T) any {
		return accessor(row)
	})
}

// BoolNField adds a variable-line boolean field.
// Accessor returns []bool: each element is one line value.
// Formatted as "Yes"/"No" for each line.
//
// Example:
//
//	builder.BoolNField("flags", "device.flags", func(r DeviceRow) []bool {
//	    return []bool{r.Active, r.Online, r.Licensed}
//	})
func (b *TableBuilder[T]) BoolNField(id, name string, accessor func(T) []bool) *FieldBuilder[T] {
	return b.fieldInternal(id, name, BoolN, func(row T) any {
		return accessor(row)
	})
}

// TimeLengthNField adds a variable-line time duration field.
// Accessor returns []int64 (values in seconds): each element is one line.
// Web/PDF output: "HH:MM" or "Xd HH:MM" format for each line
// CSV/Excel output: integer minutes joined with " - "
//
// Example:
//
//	builder.TimeLengthNField("durations", "trip.durations", func(r TripRow) []int64 {
//	    return []int64{r.DriveTime, r.IdleTime, r.StopTime}
//	})
func (b *TableBuilder[T]) TimeLengthNField(id, name string, accessor func(T) []int64) *FieldBuilder[T] {
	return b.fieldInternal(id, name, TimeLengthN, func(row T) any {
		return accessor(row)
	})
}

// BoolField adds a boolean field.
// Use .WithBoolText("Yes", "No") to customize display values.
//
// Example:
//
//	builder.BoolField("active", "device.active", func(r DeviceRow) bool {
//	    return r.Active
//	}).WithBoolText("Active", "Inactive")
func (b *TableBuilder[T]) BoolField(id, name string, accessor func(T) bool) *FieldBuilder[T] {
	return b.fieldInternal(id, name, Bool, func(row T) any {
		return accessor(row)
	})
}

// DateTimeField adds a timestamp field with date+time formatting.
// Accessor returns time.Time (converted to Unix seconds internally).
//
// Example:
//
//	builder.DateTimeField("last_seen", "device.last_seen", func(r DeviceRow) time.Time {
//	    return r.LastSeen
//	})
func (b *TableBuilder[T]) DateTimeField(id, name string, accessor func(T) time.Time) *FieldBuilder[T] {
	return b.fieldInternal(id, name, DateTime, func(row T) any {
		return accessor(row).Unix() // Convert to Unix seconds
	})
}

// DateField adds a timestamp field with date-only formatting.
// Accessor returns time.Time (converted to Unix seconds internally).
//
// Example:
//
//	builder.DateField("created", "device.created", func(r DeviceRow) time.Time {
//	    return r.Created
//	})
func (b *TableBuilder[T]) DateField(id, name string, accessor func(T) time.Time) *FieldBuilder[T] {
	return b.fieldInternal(id, name, Date, func(row T) any {
		return accessor(row).Unix() // Convert to Unix seconds
	})
}

// DistanceField adds a distance field with automatic km/mi/NM conversion.
// Expects value in kilometers.
//
// Example:
//
//	builder.DistanceField("distance", "trip.distance", func(r TripRow) float64 {
//	    return r.DistanceKm
//	}).WithDecimals(2)
func (b *TableBuilder[T]) DistanceField(id, name string, accessor func(T) float64) *FieldBuilder[T] {
	return b.fieldInternal(id, name, Distance, func(row T) any {
		return accessor(row)
	})
}

// PressureField adds a pressure field with automatic bar/psi conversion.
// Expects value in bar.
//
// Example:
//
//	builder.PressureField("pressure", "sensor.pressure", func(r SensorRow) float64 {
//	    return r.PressureBar
//	}).WithDecimals(2)
func (b *TableBuilder[T]) PressureField(id, name string, accessor func(T) float64) *FieldBuilder[T] {
	return b.fieldInternal(id, name, Pressure, func(row T) any {
		return accessor(row)
	})
}

// SpeedField adds a speed field with automatic km/h, mph, knots conversion.
// Expects value in km/h.
//
// Example:
//
//	builder.SpeedField("speed", "trip.speed", func(r TripRow) float64 {
//	    return r.SpeedKmh
//	}).WithDecimals(1)
func (b *TableBuilder[T]) SpeedField(id, name string, accessor func(T) float64) *FieldBuilder[T] {
	return b.fieldInternal(id, name, Speed, func(row T) any {
		return accessor(row)
	})
}

// ButtonsField adds a buttons field for row actions.
// Accessor returns map[string]string where key is button index ("0", "1", etc.)
// and value is URL string or empty string to hide button.
//
// Example:
//
//	builder.ButtonsField("actions", "device.actions", func(r DeviceRow) map[string]string {
//	    return map[string]string{
//	        "0": fmt.Sprintf("/Portal/Device/Edit?id=%d", r.ID),
//	        "1": fmt.Sprintf("/Portal/Device/Delete?id=%d", r.ID),
//	    }
//	}).AddButton(0, FieldButtonActionLink, "edit", FieldColorPrimary, "common.edit")
func (b *TableBuilder[T]) ButtonsField(id, name string, accessor func(T) map[string]string) *FieldBuilder[T] {
	return b.fieldInternal(id, name, Buttons, func(row T) any {
		// Convert map[string]string to map[string]any
		// Empty string values become false to hide buttons
		buttons := accessor(row)
		result := make(map[string]any, len(buttons))
		for k, v := range buttons {
			if v == "" {
				result[k] = false // Empty string = hide button
			} else {
				result[k] = v
			}
		}
		return result
	})
}

// IconFieldFromSet adds an icon field using a predefined [IconSet].
// The accessor returns *[IconRef] instead of string, ensuring at compile time
// that only registered icon values can be used. This prevents typos and
// missing icon definitions.
//
// When the accessor returns nil (e.g., from [IconSet.Resolve] for unknown values),
// an empty string is written to the row data.
//
// Example:
//
//	statusIcons := table.NewIconSet()
//	iconOnline  := statusIcons.Add("online",  "check_circle", table.FieldColorAccent,  "Online")
//	iconOffline := statusIcons.Add("offline", "cancel",       table.FieldColorWarning, "Offline")
//
//	builder.IconFieldFromSet("status", "device.status",
//	    func(r DeviceRow) *table.IconRef {
//	        if r.IsOnline {
//	            return iconOnline
//	        }
//	        return iconOffline
//	    },
//	    statusIcons,
//	)
func (b *TableBuilder[T]) IconFieldFromSet(id, name string, accessor func(T) *IconRef, iconSet *IconSet) *FieldBuilder[T] {
	fb := b.fieldInternal(id, name, Icon, func(row T) any {
		ref := accessor(row)
		if ref == nil {
			return ""
		}
		return ref.value
	})

	// Copy icon definitions from IconSet into the field
	for _, value := range iconSet.order {
		def := iconSet.icons[value]
		fb.field.addIcon(value, def.Icon, def.Color, def.Hint)
		// Copy custom options
		for k, v := range def.Options {
			fb.field.icons[value].Options[k] = v
		}
	}

	return fb
}

// LinkField adds a clickable link field.
// Accessor returns [2]string: [0] = display text, [1] = URL
//
// Example:
//
//	builder.LinkField("url", "device.url", func(r DeviceRow) [2]string {
//	    return [2]string{r.Name, r.URL}
//	})
func (b *TableBuilder[T]) LinkField(id, name string, accessor func(T) [2]string) *FieldBuilder[T] {
	return b.fieldInternal(id, name, Link, func(row T) any {
		return accessor(row)
	})
}

// HtmlField adds an HTML content field.
// Accessor returns raw HTML string.
//
// Example:
//
//	builder.HtmlField("notes", "device.notes", func(r DeviceRow) string {
//	    return r.NotesHTML
//	})
func (b *TableBuilder[T]) HtmlField(id, name string, accessor func(T) string) *FieldBuilder[T] {
	return b.fieldInternal(id, name, Html, func(row T) any {
		return accessor(row)
	})
}

// InputField adds an inline editable input field.
// Accepts any type since input fields can be text, number, select, etc.
//
// Example:
//
//	builder.InputField("notes", "device.notes", func(r DeviceRow) any {
//	    return r.Notes
//	}).WithInputType("text")
func (b *TableBuilder[T]) InputField(id, name string, accessor func(T) any) *FieldBuilder[T] {
	return b.fieldInternal(id, name, Input, func(row T) any {
		return accessor(row)
	})
}

// HeaderField adds a section header field.
//
// Example:
//
//	builder.HeaderField("section", "device.section", func(r DeviceRow) string {
//	    return "Device Information"
//	})
func (b *TableBuilder[T]) HeaderField(id, name string, accessor func(T) string) *FieldBuilder[T] {
	return b.fieldInternal(id, name, Header, func(row T) any {
		return accessor(row)
	})
}
