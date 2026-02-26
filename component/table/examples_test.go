package table_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/xiriframework/xiri-go/component/table"
	xurl "github.com/xiriframework/xiri-go/component/url"
	"github.com/xiriframework/xiri-go/form/field"
	"github.com/xiriframework/xiri-go/form/group"
	"github.com/xiriframework/xiri-go/types/distance"
	"github.com/xiriframework/xiri-go/types/language"
	"github.com/xiriframework/xiri-go/types/locale"
	"github.com/xiriframework/xiri-go/types/pressure"
	"github.com/xiriframework/xiri-go/types/timezone"
	"github.com/xiriframework/xiri-go/uicontext"
)

// ============================================================================
// TEST DATA & HELPERS
// ============================================================================

// DeviceTableRow is an example row struct for a device table
type DeviceTableRow struct {
	ID          int64
	Name        string
	NameLine2   string // For Text2 field
	IMEI        string
	DistanceKm  float64
	SpeedKmh    float64
	PressureBar float64
	LastSeen    time.Time
	CreatedDate time.Time
	Active      bool
	Status      string
	DeviceID    int32 // For cross-field access
	GroupName   string
	URL         string
	LinkText    string // For Link field display text
	Notes       string
	// Fields for Text2 variant examples
	TripsToday  int
	TotalTrips  int
	FuelCurrent float64
	FuelAverage float64
	TodayKm     float64
	TotalKm     float64
	MaxSpeed    float64
	AvgSpeed    float64
	Online      bool
}

// TripTableRow is an example row struct for a trip table
type TripTableRow struct {
	ID         int32
	DeviceName string
	StartTime  int64
	EndTime    int64
	DistanceKm float64
	Duration   int32 // seconds
	MaxSpeed   float64
	AvgSpeed   float64
	DeviceID   int32
	// Fields for TimeLength examples
	IdleTime       int64 // seconds
	DurationShort  int64 // seconds, < 1 hour
	DurationMedium int64 // seconds, < 24 hours
	DurationLong   int64 // seconds, >= 24 hours
}

// Example translator function
func exampleTranslator(key string) string {
	translations := map[string]string{
		"device.id":        "ID",
		"device.name":      "Device Name",
		"device.imei":      "IMEI",
		"device.distance":  "Distance",
		"device.speed":     "Speed",
		"device.pressure":  "Pressure",
		"device.last_seen": "Last Seen",
		"device.created":   "Created",
		"device.active":    "Active",
		"device.status":    "Status",
		"device.group":     "Group",
		"device.url":       "URL",
		"device.notes":     "Notes",
		"device.actions":   "Actions",
		"trip.id":          "Trip ID",
		"trip.device":      "Device",
		"trip.start":       "Start Time",
		"trip.end":         "End Time",
		"trip.distance":    "Distance",
		"trip.duration":    "Duration",
		"trip.max_speed":   "Max Speed",
		"trip.avg_speed":   "Avg Speed",
		"common.yes":       "Yes",
		"common.no":        "No",
		"common.online":    "Online",
		"common.offline":   "Offline",
		"common.edit":      "Edit",
		"common.delete":    "Delete",
		"common.view":      "View",
		"common.total":     "Total",
		"common.active":    "Active",
		"common.inactive":  "Inactive",
		"section.basic":    "Basic Information",
		"section.metrics":  "Metrics",
		"section.location": "Location Data",
	}
	if val, ok := translations[key]; ok {
		return val
	}
	return key
}

// Example UiContext (German user, km preference)
func exampleContext() *uicontext.UiContext {
	return &uicontext.UiContext{
		Timezone: timezone.EuropeVienna,
		Lang:     language.Deutsch,
		Locale:   locale.De, // German locale (comma decimal separator)
		Distance: distance.Kilometer,
		Pressure: pressure.Bar,
	}
}

// Example UiContext (English user, miles preference)
func exampleContextEnglish() *uicontext.UiContext {
	return &uicontext.UiContext{
		Timezone: timezone.EuropeLondon, // UK timezone
		Lang:     language.Englisch,     // English
		Locale:   locale.EnGB,           // English locale (dot decimal separator)
		Distance: distance.Miles,
		Pressure: pressure.Psi,
	}
}

// Sample device data generator
func generateDeviceData() []DeviceTableRow {
	return []DeviceTableRow{
		{
			ID: 1, Name: "Device 1", NameLine2: "Fleet A - Primary",
			IMEI:       "123456789012345",
			DistanceKm: 196.5, SpeedKmh: 80.0, PressureBar: 2.5,
			LastSeen: time.Unix(1640000000, 0), CreatedDate: time.Unix(1630000000, 0),
			Active: true, Status: "online", DeviceID: 1,
			GroupName: "Fleet A", URL: "https://example.com/device/1",
			LinkText: "View Device 1", Notes: "Primary vehicle",
			TripsToday: 5, TotalTrips: 1234,
			FuelCurrent: 45.5, FuelAverage: 52.3,
			TodayKm: 123.4, TotalKm: 98765.4,
			MaxSpeed: 120.5, AvgSpeed: 65.8,
			Online: true,
		},
		{
			ID: 2, Name: "Device 2", NameLine2: "Fleet B - Secondary",
			IMEI:       "234567890123456",
			DistanceKm: 342.8, SpeedKmh: 65.5, PressureBar: 2.2,
			LastSeen: time.Unix(1640010000, 0), CreatedDate: time.Unix(1631000000, 0),
			Active: false, Status: "offline", DeviceID: 2,
			GroupName: "Fleet B", URL: "https://example.com/device/2",
			LinkText: "View Device 2", Notes: "Secondary vehicle",
			TripsToday: 3, TotalTrips: 789,
			FuelCurrent: 38.2, FuelAverage: 48.7,
			TodayKm: 87.6, TotalKm: 54321.0,
			MaxSpeed: 110.0, AvgSpeed: 58.3,
			Online: false,
		},
		{
			ID: 3, Name: "Device 3", NameLine2: "Fleet A - Backup",
			IMEI:       "345678901234567",
			DistanceKm: 89.2, SpeedKmh: 45.0, PressureBar: 2.8,
			LastSeen: time.Unix(1640020000, 0), CreatedDate: time.Unix(1632000000, 0),
			Active: true, Status: "online", DeviceID: 3,
			GroupName: "Fleet A", URL: "https://example.com/device/3",
			LinkText: "View Device 3", Notes: "",
			TripsToday: 8, TotalTrips: 2456,
			FuelCurrent: 51.8, FuelAverage: 55.1,
			TodayKm: 234.5, TotalKm: 123456.7,
			MaxSpeed: 95.3, AvgSpeed: 72.1,
			Online: true,
		},
	}
}

// Sample trip data generator
func generateTripData() []TripTableRow {
	return []TripTableRow{
		{
			ID: 1, DeviceName: "Device 1",
			StartTime: 1640000000, EndTime: 1640003600,
			DistanceKm: 45.2, Duration: 3600,
			MaxSpeed: 90.0, AvgSpeed: 45.2,
			DeviceID:       1,
			IdleTime:       1200,   // 20 minutes
			DurationShort:  2700,   // 45 minutes
			DurationMedium: 19800,  // 5 hours 30 minutes
			DurationLong:   183900, // 2 days 3 hours 5 minutes
		},
		{
			ID: 2, DeviceName: "Device 2",
			StartTime: 1640010000, EndTime: 1640017200,
			DistanceKm: 82.5, Duration: 7200,
			MaxSpeed: 85.0, AvgSpeed: 41.2,
			DeviceID:       2,
			IdleTime:       1800,   // 30 minutes
			DurationShort:  1800,   // 30 minutes
			DurationMedium: 25200,  // 7 hours
			DurationLong:   259200, // 3 days
		},
		{
			ID: 3, DeviceName: "Device 1",
			StartTime: 1640020000, EndTime: 1640025600,
			DistanceKm: 65.8, Duration: 5600,
			MaxSpeed: 95.0, AvgSpeed: 42.2,
			DeviceID:       1,
			IdleTime:       900,    // 15 minutes
			DurationShort:  3300,   // 55 minutes
			DurationMedium: 14400,  // 4 hours
			DurationLong:   432000, // 5 days
		},
	}
}

// Helper function for pointer values
func ptr[T any](v T) *T {
	return &v
}

// ============================================================================
// SECTION 1: BASIC FIELD TYPE EXAMPLES
// ============================================================================

// TestExample01_IntegerField demonstrates integer number field
func TestExample01_IntegerField(t *testing.T) {
	ctx := exampleContext()
	rows := generateDeviceData()

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	// ID field - returns int64 value
	builder.IdField("id", "device.id", func(r DeviceTableRow) int64 {
		return r.ID
	})

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(table.OutputWeb)

	// Verify format: ["1", 1]
	if idValue, ok := data[0]["id"].([]any); ok {
		fmt.Printf("Integer field output: %v\n", idValue)
		// Output: ["1", 1]
	}

	// ID field with alignment and width
	builder2 := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)
	builder2.IdField("id", "device.id", func(r DeviceTableRow) int64 {
		return r.ID
	}).AlignRight().WithWidth("80px")

	tbl2 := builder2.Build()
	tbl.SetData(rows)
	_ = tbl2.Print(exampleTranslator)
}

// TestExample02_FloatField demonstrates float number field
func TestExample02_FloatField(t *testing.T) {
	ctx := exampleContext()
	rows := generateDeviceData()

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	// Float field with 2 decimals
	builder.FloatField("distance", "device.distance", func(r DeviceTableRow) float64 {
		return r.DistanceKm
	}).WithDecimals(2)

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(table.OutputWeb)

	// German locale uses comma: ["196,50", 196.5]
	if distValue, ok := data[0]["distance"].([]any); ok {
		fmt.Printf("Float field output (German locale): %v\n", distValue)
		// Output: ["196,50", 196.5]
	}

	// Float field with 3 decimals and footer sum
	builder2 := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)
	builder2.FloatField("distance", "device.distance", func(r DeviceTableRow) float64 {
		return r.DistanceKm
	}).WithDecimals(3).WithFooterSum().AlignRight()

	tbl2 := builder2.Build()
	tbl.SetData(rows)
	footer := tbl2.CalculateFooter(table.OutputWeb)

	if footerDist, ok := footer["distance"].([]any); ok {
		fmt.Printf("Footer sum: %v\n", footerDist)
		// Output: ["628,500", 628.5]
	}
}

// TestExample03_TextField demonstrates text field
func TestExample03_TextField(t *testing.T) {
	ctx := exampleContext()
	rows := generateDeviceData()

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	// Text field - simple string output
	builder.TextField("name", "device.name", func(r DeviceTableRow) string {
		return r.Name
	})

	// Text field with prefix/suffix
	builder.TextField("imei", "device.imei", func(r DeviceTableRow) string {
		return r.IMEI
	}).WithTextPrefix("IMEI:")

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(table.OutputWeb)

	// Text field returns string directly
	if nameValue, ok := data[0]["name"].(string); ok {
		fmt.Printf("Text field output: %s\n", nameValue)
		// Output: Device 1
	}
}

// TestExample04_BoolField demonstrates boolean field
func TestExample04_BoolField(t *testing.T) {
	ctx := exampleContext()
	rows := generateDeviceData()

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	// Bool field with Yes/No text
	builder.BoolField("active", "device.active", func(r DeviceTableRow) bool {
		return r.Active
	}).WithBoolText("Yes", "No")

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(table.OutputWeb)

	// Bool field returns formatted text
	if activeValue, ok := data[0]["active"].(string); ok {
		fmt.Printf("Bool field output: %s\n", activeValue)
		// Output: Yes
	}

	// Bool field with footer count
	builder2 := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)
	builder2.BoolField("active", "device.active", func(r DeviceTableRow) bool {
		return r.Active
	}).WithBoolText("Active", "Inactive").WithFooterCount()

	tbl2 := builder2.Build()
	tbl.SetData(rows)
	footer := tbl2.CalculateFooter(table.OutputWeb)

	if footerCount, ok := footer["active"].(int); ok {
		fmt.Printf("Active count: %d\n", footerCount)
		// Output: 2 (number of active devices)
	}
}

// TestExample05_DateTimeField demonstrates datetime field with timezone
func TestExample05_DateTimeField(t *testing.T) {
	ctx := exampleContext()
	rows := generateDeviceData()

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	// DateTime field - respects user timezone
	builder.DateTimeField("last_seen", "device.last_seen", func(r DeviceTableRow) time.Time {
		return r.LastSeen
	})

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(table.OutputWeb)

	// DateTime formatted with user timezone
	if lastSeenValue, ok := data[0]["last_seen"].(string); ok {
		fmt.Printf("DateTime field output: %s\n", lastSeenValue)
		// Output: 2021-12-20 12:26 (German format with Vienna timezone)
	}

	// Different format for CSV
	csvData := tbl.GetData(table.OutputCSV)
	if csvLastSeen, ok := csvData[0]["last_seen"].(string); ok {
		fmt.Printf("DateTime CSV output: %s\n", csvLastSeen)
		// Output: 2021-12-20 12:26:40 (ISO format)
	}
}

// TestExample06_DateField demonstrates date-only field
func TestExample06_DateField(t *testing.T) {
	ctx := exampleContext()
	rows := generateDeviceData()

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	// Date field - date only, no time
	builder.DateField("created", "device.created", func(r DeviceTableRow) time.Time {
		return r.CreatedDate
	})

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(table.OutputWeb)

	// Date formatted without time
	if createdValue, ok := data[0]["created"].(string); ok {
		fmt.Printf("Date field output: %s\n", createdValue)
		// Output: 2021-08-26
	}
}

// TestExample07_DistanceField demonstrates distance with unit conversion
func TestExample07_DistanceField(t *testing.T) {
	ctx := exampleContext()
	rows := generateDeviceData()

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	// Distance field with device-specific units
	// German user default: km, but each device can override
	builder.DistanceField("distance", "device.distance", func(r DeviceTableRow) float64 {
		return r.DistanceKm
	}).WithDecimals(2)

	// Hidden device_id field for cross-field access
	builder.Int32Field("device_id", "device.id", func(r DeviceTableRow) int32 {
		return r.DeviceID
	}).Hide()

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(table.OutputWeb)

	// Distance formatted with unit: ["196,50 km", 196.5]
	if distValue, ok := data[0]["distance"].([]any); ok {
		fmt.Printf("Distance field output: %v\n", distValue)
		// Output: ["196,50 km", 196.5]
	}

	// Test with English user (miles)
	ctxEnglish := exampleContextEnglish()
	builder2 := table.NewBuilder[DeviceTableRow](ctxEnglish, exampleTranslator)
	builder2.DistanceField("distance", "device.distance", func(r DeviceTableRow) float64 {
		return r.DistanceKm
	}).WithDecimals(2)
	builder2.Int32Field("device_id", "device.id", func(r DeviceTableRow) int32 {
		return r.DeviceID
	}).Hide()

	tbl2 := builder2.Build()
	tbl2.SetData(rows)
	data2 := tbl2.GetData(table.OutputWeb)

	// Distance converted to miles: ["122.09 mi", 122.09]
	if distValue2, ok := data2[0]["distance"].([]any); ok {
		fmt.Printf("Distance field output (miles): %v\n", distValue2)
		// Output: ["122.09 mi", 122.09]
	}
}

// TestExample08_SpeedField demonstrates speed with unit conversion
func TestExample08_SpeedField(t *testing.T) {
	ctx := exampleContext()
	rows := generateDeviceData()

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	// Speed field with device-specific units
	builder.SpeedField("speed", "device.speed", func(r DeviceTableRow) float64 {
		return r.SpeedKmh
	}).WithDecimals(1)

	builder.Int32Field("device_id", "device.id", func(r DeviceTableRow) int32 {
		return r.DeviceID
	}).Hide()

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(table.OutputWeb)

	// Speed formatted with unit: ["80,0 km/h", 80.0]
	if speedValue, ok := data[0]["speed"].([]any); ok {
		fmt.Printf("Speed field output: %v\n", speedValue)
		// Output: ["80,0 km/h", 80.0]
	}

	// Test with English user (mph)
	ctxEnglish := exampleContextEnglish()
	builder2 := table.NewBuilder[DeviceTableRow](ctxEnglish, exampleTranslator)
	builder2.SpeedField("speed", "device.speed", func(r DeviceTableRow) float64 {
		return r.SpeedKmh
	}).WithDecimals(1)
	builder2.Int32Field("device_id", "device.id", func(r DeviceTableRow) int32 {
		return r.DeviceID
	}).Hide()

	tbl2 := builder2.Build()
	tbl2.SetData(rows)
	data2 := tbl2.GetData(table.OutputWeb)

	// Speed converted to mph: ["49.7 mph", 49.7]
	if speedValue2, ok := data2[0]["speed"].([]any); ok {
		fmt.Printf("Speed field output (mph): %v\n", speedValue2)
		// Output: ["49.7 mph", 49.7]
	}
}

// TestExample09_PressureField demonstrates pressure with unit conversion
func TestExample09_PressureField(t *testing.T) {
	ctx := exampleContext()
	rows := generateDeviceData()

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	// Pressure field with device-specific units
	builder.PressureField("pressure", "device.pressure", func(r DeviceTableRow) float64 {
		return r.PressureBar
	}).WithDecimals(2)

	builder.Int32Field("device_id", "device.id", func(r DeviceTableRow) int32 {
		return r.DeviceID
	}).Hide()

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(table.OutputWeb)

	// Pressure formatted with unit: ["2,50 bar", 2.5]
	if pressureValue, ok := data[0]["pressure"].([]any); ok {
		fmt.Printf("Pressure field output: %v\n", pressureValue)
		// Output: ["2,50 bar", 2.5]
	}

	// Test with English user (psi)
	ctxEnglish := exampleContextEnglish()
	builder2 := table.NewBuilder[DeviceTableRow](ctxEnglish, exampleTranslator)
	builder2.PressureField("pressure", "device.pressure", func(r DeviceTableRow) float64 {
		return r.PressureBar
	}).WithDecimals(2)
	builder2.Int32Field("device_id", "device.id", func(r DeviceTableRow) int32 {
		return r.DeviceID
	}).Hide()

	tbl2 := builder2.Build()
	tbl2.SetData(rows)
	data2 := tbl2.GetData(table.OutputWeb)

	// Pressure converted to psi: ["36.26 psi", 36.26]
	if pressureValue2, ok := data2[0]["pressure"].([]any); ok {
		fmt.Printf("Pressure field output (psi): %v\n", pressureValue2)
		// Output: ["36.26 psi", 36.26]
	}
}

// TestExample10_IconField demonstrates icon field with mappings
func TestExample10_IconField(t *testing.T) {
	ctx := exampleContext()
	rows := generateDeviceData()

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	builder.TextField("name", "device.name", func(r DeviceTableRow) string {
		return r.Name
	})

	// Icon field with value -> icon/color/hint mappings
	statusIcons := table.NewIconSet()
	statusIcons.Add("online", "check_circle", table.FieldColorAccent, "Device is online")
	statusIcons.Add("offline", "cancel", table.FieldColorWarning, "Device is offline")

	builder.IconFieldFromSet("status", "device.status",
		func(r DeviceTableRow) *table.IconRef {
			return statusIcons.Resolve(r.Status)
		},
		statusIcons,
	)

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(table.OutputWeb)

	// Output icon field data
	fmt.Printf("Icon field for Device 1 (online): %v\n", data[0]["status"])
	fmt.Printf("Icon field for Device 2 (offline): %v\n", data[1]["status"])

	// Icon field output includes icon, color, and hint based on value
	// Frontend renders as colored icon with tooltip
	// Example output: map[icon:check_circle color:accent hint:Device is online]
}

// TestExample11_ButtonsField demonstrates buttons field with row actions
func TestExample11_ButtonsField(t *testing.T) {
	ctx := exampleContext()
	rows := generateDeviceData()

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	builder.IdField("id", "device.id", func(r DeviceTableRow) int64 {
		return r.ID
	})

	builder.TextField("name", "device.name", func(r DeviceTableRow) string {
		return r.Name
	})

	// Buttons field with edit/delete actions
	builder.ButtonsField("actions", "device.actions", func(r DeviceTableRow) map[string]string {
		// Return button data as map: button index -> URL (empty string to hide)
		return map[string]string{
			"0": fmt.Sprintf("/Portal/Device/Edit?id=%d", r.ID),   // Edit button (index 0)
			"1": fmt.Sprintf("/Portal/Device/Delete?id=%d", r.ID), // Delete button (index 1)
		}
	}).
		AddButton(0, table.FieldButtonActionLink, "edit", table.FieldColorPrimary, "Edit device").
		AddButton(1, table.FieldButtonActionDialog, "delete", table.FieldColorWarning, "Delete device")

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(table.OutputWeb)

	// Output buttons field data
	fmt.Printf("Buttons field for Device 1: %v\n", data[0]["actions"])
	// Expected output: map with button indices as keys, URLs or false as values
	// map[0:/Portal/Device/Edit?id=1 1:/Portal/Device/Delete?id=1]

	// Buttons field renders as action buttons in each row
	// Button 0 (edit) opens link, button 1 (delete) opens dialog
	// Set button value to false to hide it conditionally
}

// TestExample12_LinkField demonstrates link field with clickable URLs
func TestExample12_LinkField(t *testing.T) {
	ctx := exampleContext()
	rows := generateDeviceData()

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	builder.TextField("name", "device.name", func(r DeviceTableRow) string {
		return r.Name
	})

	// Link field - renders as clickable link
	builder.LinkField("url", "device.url", func(r DeviceTableRow) [2]string {
		return [2]string{r.LinkText, r.URL}
	})

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(table.OutputWeb)

	// Link field outputs TWO separate fields: "url" (display text) and "urlLink" (URL)
	if urlText, ok := data[0]["url"].(string); ok {
		fmt.Printf("Link field display text: %s\n", urlText)
		// Output: View Device 1
	}
	if urlLink, ok := data[0]["urlLink"].(string); ok {
		fmt.Printf("Link field URL: %s\n", urlLink)
		// Output: https://example.com/device/1
	}

	// Web output: Two separate fields for clickable link
	// CSV/PDF/Excel output: Display text only (single field)
	csvData := tbl.GetData(table.OutputCSV)
	if csvText, ok := csvData[0]["url"].(string); ok {
		fmt.Printf("Link field CSV output: %s\n", csvText)
		// Output: View Device 1
	}
}

// TestExample13_HtmlField demonstrates html field with raw HTML content
func TestExample13_HtmlField(t *testing.T) {
	ctx := exampleContext()
	rows := generateDeviceData()

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	builder.TextField("name", "device.name", func(r DeviceTableRow) string {
		return r.Name
	})

	// Html field - renders raw HTML
	builder.HtmlField("notes", "device.notes", func(r DeviceTableRow) string {
		// Return HTML content
		if r.Notes != "" {
			return fmt.Sprintf("<strong>%s</strong>", r.Notes)
		}
		return "<em>No notes</em>"
	})

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(table.OutputWeb)

	// Html field returns HTML string
	if htmlValue, ok := data[0]["notes"].(string); ok {
		fmt.Printf("Html field output: %s\n", htmlValue)
		// Output: <strong>Primary vehicle</strong>
	}

	// Web output: rendered HTML
	// CSV output: plain text (HTML tags stripped in real implementation)
	csvData := tbl.GetData(table.OutputCSV)
	if csvNotes, ok := csvData[0]["notes"].(string); ok {
		fmt.Printf("Html field CSV output: %s\n", csvNotes)
	}
}

// TestExample14_Text2Field demonstrates text2 field (alternative text styling)
func TestExample14_Text2Field(t *testing.T) {
	ctx := exampleContext()
	rows := generateDeviceData()

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	builder.TextField("name", "device.name", func(r DeviceTableRow) string {
		return r.Name
	})

	// Text2 field - two-line text (primary + secondary)
	builder.Text2Field("device_info", "device.info", func(r DeviceTableRow) [2]string {
		return [2]string{r.Name, r.NameLine2}
	})

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(table.OutputWeb)

	// Text2 field returns [2]string: [primary, secondary]
	if deviceInfo, ok := data[0]["device_info"].([2]string); ok {
		fmt.Printf("Text2 field output: %v\n", deviceInfo)
		// Output: [Device 1 Fleet A - Primary]
	}

	// Text2 is typically rendered with two lines:
	// - Primary text (regular weight)
	// - Secondary text (lighter/smaller font)
	// Used for subtitles, secondary info, or descriptions
}

// TestExample15_InputField demonstrates input field (inline editing)
func TestExample15_InputField(t *testing.T) {
	ctx := exampleContext()
	rows := generateDeviceData()

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	builder.IdField("id", "device.id", func(r DeviceTableRow) int64 {
		return r.ID
	})

	builder.TextField("name", "device.name", func(r DeviceTableRow) string {
		return r.Name
	})

	// Input field - enables inline editing with API callback
	// The field itself is just configured with Input type
	builder.InputField("notes", "device.notes", func(r DeviceTableRow) any {
		return r.Notes
	}).WithInputType("text").WithInputRequired(false)

	// The update URL is set at table level via SaveInputUrl
	tbl := builder.
		SetSaveInputUrl("/Portal/Device/UpdateField").
		Build()
	tbl.SetData(rows)

	data := tbl.GetData(table.OutputWeb)

	// Input field returns string value
	if notesValue, ok := data[0]["notes"].(string); ok {
		fmt.Printf("Input field output: %s\n", notesValue)
		// Output: Vehicle requires maintenance
	}

	// The table options include the save URL for all input fields
	options := tbl.Print(exampleTranslator)["data"].(map[string]any)["options"].(map[string]any)
	if saveUrl, ok := options["saveInputUrl"].(string); ok {
		fmt.Printf("Table save input URL: %s\n", saveUrl)
		// Output: /Portal/Device/UpdateField
	}

	// Input fields are useful for:
	// - Quick inline editing without opening a dialog
	// - Updating simple text fields directly in the table
	// - Batch editing multiple rows efficiently
	// - All input fields in a table share the same save URL
}

// TestExample16_HeaderField demonstrates header field (section divider)
func TestExample16_HeaderField(t *testing.T) {
	ctx := exampleContext()
	rows := generateDeviceData()

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	// Header field - creates a visual section divider in the table
	// Often used to group related columns or mark important sections
	builder.HeaderField("basic_info", "device.basic_info", func(r DeviceTableRow) string {
		return "Basic Information" // Text for header
	})

	builder.IdField("id", "device.id", func(r DeviceTableRow) int64 {
		return r.ID
	})

	builder.TextField("name", "device.name", func(r DeviceTableRow) string {
		return r.Name
	})

	builder.HeaderField("tracking_info", "device.tracking_info", func(r DeviceTableRow) string {
		return "Tracking Information"
	})

	builder.DateTimeField("last_seen", "device.last_seen", func(r DeviceTableRow) time.Time {
		return r.LastSeen
	})

	tbl := builder.Build()
	tbl.SetData(rows)
	fields := tbl.Print(exampleTranslator)["data"].(map[string]any)["fields"].([]map[string]any)

	// Header fields appear in the field definitions
	var headerCount int
	for _, field := range fields {
		// Field type is in "format" key (legacy compatibility)
		if field["format"] == "header" {
			headerCount++
			fmt.Printf("Header field: %s (text: %s)\n", field["id"], field["name"])
		}
	}

	if headerCount > 0 {
		fmt.Printf("Total header fields: %d\n", headerCount)
		// Output: 2 header fields for grouping columns
	}

	// Header fields are useful for:
	// - Organizing wide tables with many columns
	// - Creating visual groupings (e.g., "Contact Info", "Location Data")
	// - Improving readability in complex tables
	// - Matching multi-level column headers in reports
}

// TestExample17_Text2IntField demonstrates two-line integer field with locale formatting
func TestExample17_Text2IntField(t *testing.T) {
	ctx := exampleContext() // German locale

	// Create test data with today vs total metrics
	rows := []DeviceTableRow{
		{
			ID: 1, Name: "Device 1",
			TripsToday: 5, TotalTrips: 1234,
		},
	}

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	// Text2IntField - displays two integers side-by-side
	builder.Text2IntField("trips", "device.trips", func(r DeviceTableRow) [2]int {
		return [2]int{r.TripsToday, r.TotalTrips}
	})

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(table.OutputWeb)

	// Web output: [2]string with locale-aware formatting
	trips := data[0]["trips"].([2]string)
	fmt.Printf("Today vs Total Trips: [%s, %s]\n", trips[0], trips[1])
	// Output: ["5", "1.234"] (German locale uses dot for thousands)

	// CSV output: combined with " - " separator
	csvData := tbl.GetData(table.OutputCSV)
	fmt.Printf("CSV Trips: %s\n", csvData[0]["trips"])
	// Output: "5 - 1234"
}

// TestExample18_Text2FloatField demonstrates two-line float field
func TestExample18_Text2FloatField(t *testing.T) {
	ctx := exampleContextEnglish() // English locale

	rows := []DeviceTableRow{
		{
			ID: 1, Name: "Device 1",
			FuelCurrent: 45.67, FuelAverage: 52.34,
		},
	}

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	// Text2FloatField - current vs average with 2 decimals
	builder.Text2FloatField("fuel", "device.fuel", func(r DeviceTableRow) [2]float64 {
		return [2]float64{r.FuelCurrent, r.FuelAverage}
	}).WithDecimals(2)

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(table.OutputWeb)

	fuel := data[0]["fuel"].([2]string)
	fmt.Printf("Current vs Average Fuel: [%s, %s]\n", fuel[0], fuel[1])
	// Output: ["45.67", "52.34"] (English locale uses dot for decimals)
}

// TestExample19_Text2DateTimeField demonstrates two-line datetime field
func TestExample19_Text2DateTimeField(t *testing.T) {
	ctx := exampleContext()

	lastSeen := time.Unix(1640000000, 0) // 2021-12-20
	created := time.Unix(1630000000, 0)  // 2021-08-26

	rows := []DeviceTableRow{
		{
			ID: 1, Name: "Device 1",
			LastSeen: lastSeen, CreatedDate: created,
		},
	}

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	// Text2DateTimeField - last seen vs created timestamps
	builder.Text2DateTimeField("times", "device.times", func(r DeviceTableRow) [2]time.Time {
		return [2]time.Time{r.LastSeen, r.CreatedDate}
	})

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(table.OutputWeb)

	times := data[0]["times"].([2]string)
	fmt.Printf("Last Seen vs Created: [%s, %s]\n", times[0], times[1])
	// Output: Timezone-aware datetime formatting for both timestamps
}

// TestExample20_Text2DistanceField demonstrates two-line distance field with unit conversion
func TestExample20_Text2DistanceField(t *testing.T) {
	ctx := exampleContext() // German locale, km units

	rows := []DeviceTableRow{
		{
			ID: 1, DeviceID: 1, Name: "Device 1",
			TodayKm: 123.45, TotalKm: 9876.54,
		},
	}

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	// Need device_id field for unit conversion
	builder.IntField("device_id", "device.id", func(r DeviceTableRow) int {
		return int(r.DeviceID)
	})

	// Text2DistanceField - today vs total with 1 decimal
	builder.Text2DistanceField("distances", "device.distances", func(r DeviceTableRow) [2]float64 {
		return [2]float64{r.TodayKm, r.TotalKm}
	}).WithDecimals(1)

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(table.OutputWeb)

	distances := data[0]["distances"].([2]string)
	fmt.Printf("Today vs Total Distance: [%s, %s]\n", distances[0], distances[1])
	// Output: ["123,5 km", "9.876,5 km"] (German locale + km units)

	// Test CSV output
	csvData := tbl.GetData(table.OutputCSV)
	distancesCSV := csvData[0]["distances"].(string)
	fmt.Printf("CSV output: %s\n", distancesCSV)
	// Output: "123.5 - 9876.5" (combined with separator)
}

// TestExample21_Text2SpeedField demonstrates two-line speed field with unit conversion
func TestExample21_Text2SpeedField(t *testing.T) {
	ctx := exampleContext()

	rows := []DeviceTableRow{
		{
			ID: 1, DeviceID: 1, Name: "Device 1",
			MaxSpeed: 120.5, AvgSpeed: 85.3,
		},
	}

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	builder.IntField("device_id", "device.id", func(r DeviceTableRow) int {
		return int(r.DeviceID)
	})

	// Text2SpeedField - max vs average speed
	builder.Text2SpeedField("speeds", "device.speeds", func(r DeviceTableRow) [2]float64 {
		return [2]float64{r.MaxSpeed, r.AvgSpeed}
	}).WithDecimals(1)

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(table.OutputWeb)

	speeds := data[0]["speeds"].([2]string)
	fmt.Printf("Max vs Avg Speed: [%s, %s]\n", speeds[0], speeds[1])
	// Output: ["120,5 km/h", "85,3 km/h"]
}

// TestExample22_Text2BoolField demonstrates two-line boolean field
func TestExample22_Text2BoolField(t *testing.T) {
	ctx := exampleContext()

	rows := []DeviceTableRow{
		{
			ID: 1, Name: "Device 1",
			Active: true, Online: false,
		},
	}

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	// Text2BoolField - active vs online status
	builder.Text2BoolField("status", "device.status", func(r DeviceTableRow) [2]bool {
		return [2]bool{r.Active, r.Online}
	})

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(table.OutputWeb)

	status := data[0]["status"].([2]string)
	fmt.Printf("Active vs Online: [%s, %s]\n", status[0], status[1])
	// Output: ["Yes", "No"]

	csvData := tbl.GetData(table.OutputCSV)
	fmt.Printf("CSV Status: %s\n", csvData[0]["status"])
	// Output: "Yes - No"
}

// TestExample23_TimeLengthField demonstrates time duration field
func TestExample23_TimeLengthField(t *testing.T) {
	ctx := exampleContext()

	rows := []TripTableRow{
		{
			ID: 1, DeviceName: "Device 1",
			DurationShort:  2700,   // 45 minutes
			DurationMedium: 19800,  // 5 hours 30 minutes
			DurationLong:   183900, // 2 days 3 hours 5 minutes
		},
	}

	builder := table.NewBuilder[TripTableRow](ctx, exampleTranslator)

	// TimeLengthField - formats seconds as HH:MM or "Xd HH:MM"
	builder.TimeLengthField("short", "trip.short", func(r TripTableRow) int64 {
		return r.DurationShort
	})

	builder.TimeLengthField("medium", "trip.medium", func(r TripTableRow) int64 {
		return r.DurationMedium
	})

	builder.TimeLengthField("long", "trip.long", func(r TripTableRow) int64 {
		return r.DurationLong
	})

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(table.OutputWeb)

	fmt.Printf("Short duration: %s\n", data[0]["short"])   // "00:45"
	fmt.Printf("Medium duration: %s\n", data[0]["medium"]) // "05:30"
	fmt.Printf("Long duration: %s\n", data[0]["long"])     // "2d 03:05"

	// CSV exports as integer minutes
	csvData := tbl.GetData(table.OutputCSV)
	fmt.Printf("CSV short: %s minutes\n", csvData[0]["short"])   // "45"
	fmt.Printf("CSV medium: %s minutes\n", csvData[0]["medium"]) // "330"
	fmt.Printf("CSV long: %s minutes\n", csvData[0]["long"])     // "3065"
}

// TestExample24_Text2TimeLengthField demonstrates two-line time duration field
func TestExample24_Text2TimeLengthField(t *testing.T) {
	ctx := exampleContext()

	rows := []TripTableRow{
		{
			ID: 1, DeviceName: "Device 1",
			Duration: 19800, // 5 hours 30 minutes
			IdleTime: 7200,  // 2 hours
		},
	}

	builder := table.NewBuilder[TripTableRow](ctx, exampleTranslator)

	// Text2TimeLengthField - duration vs idle time
	builder.Text2TimeLengthField("times", "trip.times", func(r TripTableRow) [2]int64 {
		return [2]int64{int64(r.Duration), r.IdleTime}
	})

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(table.OutputWeb)

	times := data[0]["times"].([2]string)
	fmt.Printf("Duration vs Idle: [%s, %s]\n", times[0], times[1])
	// Output: ["05:30", "02:00"]

	csvData := tbl.GetData(table.OutputCSV)
	fmt.Printf("CSV times: %s\n", csvData[0]["times"])
	// Output: "330 - 120" (minutes)
}

// ===== ADVANCED CONFIGURATION EXAMPLES =====
// These examples demonstrate advanced table configuration options

// TestAdvanced01_FieldAlignment demonstrates text alignment options
func TestAdvanced01_FieldAlignment(t *testing.T) {
	ctx := exampleContext()
	rows := generateDeviceData()

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	// Left-aligned field (default for text)
	builder.TextField("name", "device.name", func(r DeviceTableRow) string {
		return r.Name
	}).WithAlign(table.FieldAlignLeft)

	// Center-aligned field (useful for status indicators)
	statusIcons := table.NewIconSet()
	statusIcons.Add("online", "check_circle", table.FieldColorAccent, "Device is online")
	statusIcons.Add("offline", "cancel", table.FieldColorWarning, "Device is offline")

	builder.IconFieldFromSet("status", "device.status",
		func(r DeviceTableRow) *table.IconRef {
			return statusIcons.Resolve(r.Status)
		},
		statusIcons,
	).WithAlign(table.FieldAlignCenter)

	// Right-aligned field (default for numbers)
	builder.DistanceField("distance", "device.distance", func(r DeviceTableRow) float64 {
		return r.DistanceKm
	}).WithAlign(table.FieldAlignRight).WithDecimals(1)

	tbl := builder.Build()
	tbl.SetData(rows)
	fields := tbl.Print(exampleTranslator)["data"].(map[string]any)["fields"].([]map[string]any)

	// Verify alignment settings
	fmt.Printf("Name field alignment: %v\n", fields[0]["align"])     // left
	fmt.Printf("Status field alignment: %v\n", fields[1]["align"])   // center
	fmt.Printf("Distance field alignment: %v\n", fields[2]["align"]) // right

	// Alignment affects:
	// - Text positioning in table cells
	// - Header text positioning
	// - Visual hierarchy and readability
}

// TestAdvanced02_FieldWidths demonstrates width and minimum width configuration
func TestAdvanced02_FieldWidths(t *testing.T) {
	ctx := exampleContext()
	rows := generateDeviceData()

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	// Fixed width field (ID column)
	builder.IdField("id", "device.id", func(r DeviceTableRow) int64 {
		return r.ID
	}).WithWidth("80px")

	// Minimum width field (Name column with responsive sizing)
	builder.TextField("name", "device.name", func(r DeviceTableRow) string {
		return r.Name
	}).WithMinWidth("200px")

	// No width constraint (flexible column)
	builder.TextField("notes", "device.notes", func(r DeviceTableRow) string {
		return r.Notes
	})

	tbl := builder.Build()
	tbl.SetData(rows)
	fields := tbl.Print(exampleTranslator)["data"].(map[string]any)["fields"].([]map[string]any)

	// Verify width settings
	if width, ok := fields[0]["width"].(string); ok {
		fmt.Printf("ID field has fixed width: %s\n", width)
	}
	if minWidth, ok := fields[1]["minWidth"].(string); ok {
		fmt.Printf("Name field has minimum width: %s\n", minWidth)
	}
	if _, hasWidth := fields[2]["width"]; !hasWidth {
		fmt.Printf("Notes field has flexible width\n")
	}

	// Width configuration:
	// - width: Fixed pixel width (e.g., "80px")
	// - minWidth: Minimum responsive width (e.g., "200px")
	// - No setting: Flexible, shares remaining space
}

// TestAdvanced03_StickyColumns demonstrates sticky/frozen columns
func TestAdvanced03_StickyColumns(t *testing.T) {
	ctx := exampleContext()
	rows := generateDeviceData()

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	// Sticky ID column (always visible when scrolling horizontally)
	builder.IdField("id", "device.id", func(r DeviceTableRow) int64 {
		return r.ID
	}).WithWidth("80px").WithSticky(true)

	// Sticky Name column
	builder.TextField("name", "device.name", func(r DeviceTableRow) string {
		return r.Name
	}).WithMinWidth("200px").WithSticky(true)

	// Regular scrollable columns
	builder.DistanceField("distance", "device.distance", func(r DeviceTableRow) float64 {
		return r.DistanceKm
	}).WithDecimals(1)

	builder.SpeedField("speed", "device.speed", func(r DeviceTableRow) float64 {
		return r.SpeedKmh
	}).WithDecimals(1)

	tbl := builder.Build()
	tbl.SetData(rows)
	fields := tbl.Print(exampleTranslator)["data"].(map[string]any)["fields"].([]map[string]any)

	// Count sticky fields
	var stickyCount int
	for i, field := range fields {
		if sticky, ok := field["sticky"].(bool); ok && sticky {
			stickyCount++
			fmt.Printf("Field %d (%s) is sticky\n", i, field["id"])
		}
	}
	fmt.Printf("Total sticky fields: %d\n", stickyCount)

	// Sticky columns:
	// - Remain visible when scrolling horizontally
	// - Useful for ID, name, or action columns
	// - Typically used on first 1-2 columns
	// - Improves usability on wide tables
}

// TestAdvanced04_FieldHints demonstrates tooltips/hints
func TestAdvanced04_FieldHints(t *testing.T) {
	ctx := exampleContext()
	rows := generateDeviceData()

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	builder.IdField("id", "device.id", func(r DeviceTableRow) int64 {
		return r.ID
	}).WithHint("device.hint.id") // Translation key for hint

	builder.TextField("name", "device.name", func(r DeviceTableRow) string {
		return r.Name
	}).WithHint("device.hint.name")

	builder.DistanceField("distance", "device.distance", func(r DeviceTableRow) float64 {
		return r.DistanceKm
	}).WithHint("device.hint.distance").WithDecimals(1)

	tbl := builder.Build()
	tbl.SetData(rows)
	fields := tbl.Print(exampleTranslator)["data"].(map[string]any)["fields"].([]map[string]any)

	// Verify hints are set
	for _, field := range fields {
		if hint, ok := field["hint"].(string); ok {
			fmt.Printf("Field %s has hint: %s\n", field["id"], hint)
		}
	}

	// Hints/tooltips:
	// - Show on column header hover
	// - Provide additional context or explanation
	// - Use translation keys for i18n support
	// - Keep brief and informative
}

// TestAdvanced05_TextPrefixSuffix demonstrates text prefixes and suffixes
func TestAdvanced05_TextPrefixSuffix(t *testing.T) {
	ctx := exampleContext()
	rows := generateDeviceData()

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	// Currency prefix
	builder.IdField("id", "device.id", func(r DeviceTableRow) int64 {
		return r.ID
	}).WithTextPrefix("ID:")

	// Unit suffix (alternative to built-in formatters)
	builder.FloatField("custom_distance", "device.distance", func(r DeviceTableRow) float64 {
		return r.DistanceKm
	}).WithDecimals(2).WithTextSuffix(" km")

	// Both prefix and suffix
	builder.TextField("name", "device.name", func(r DeviceTableRow) string {
		return r.Name
	}).WithTextPrefix("[").WithTextSuffix("]")

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(table.OutputWeb)

	// Prefix/suffix are applied to formatted output
	fmt.Printf("ID field with prefix: %v\n", data[0]["id"])
	fmt.Printf("Distance with suffix: %v\n", data[0]["custom_distance"])
	fmt.Printf("Name with brackets: %v\n", data[0]["name"])

	// Text prefix/suffix:
	// - Applied after value formatting
	// - Useful for units, labels, decorations
	// - Simple alternative to custom formatters
	// - Not translated (use static text)
}

// TestAdvanced06_CustomHeaders demonstrates custom header text and spans
func TestAdvanced06_CustomHeaders(t *testing.T) {
	ctx := exampleContext()
	rows := generateDeviceData()

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	// Fields with custom header (overrides column name)
	builder.IdField("id", "device.id", func(r DeviceTableRow) int64 {
		return r.ID
	}).WithHeader("Device ID")

	builder.TextField("name", "device.name", func(r DeviceTableRow) string {
		return r.Name
	}).WithHeader("Device Name")

	// Field with header span (used for grouped columns)
	builder.DistanceField("distance", "device.distance", func(r DeviceTableRow) float64 {
		return r.DistanceKm
	}).WithHeader("Tracking Data").WithHeaderSpan(2).WithDecimals(1)

	builder.SpeedField("speed", "device.speed", func(r DeviceTableRow) float64 {
		return r.SpeedKmh
	}).WithDecimals(1)

	tbl := builder.Build()
	tbl.SetData(rows)
	fields := tbl.Print(exampleTranslator)["data"].(map[string]any)["fields"].([]map[string]any)

	// Verify custom headers
	for _, field := range fields {
		if header, ok := field["header"].(string); ok {
			fmt.Printf("Field %s has custom header: %s", field["id"], header)
			if span, ok := field["headerSpan"].(int); ok {
				fmt.Printf(" (spans %d columns)", span)
			}
			fmt.Println()
		}
	}

	// Custom headers:
	// - Override the translated field name
	// - headerSpan creates grouped column headers
	// - Useful for complex multi-level headers
	// - Can be dynamic based on context
}

// TestAdvanced07_ColumnOrdering demonstrates explicit column order
func TestAdvanced07_ColumnOrdering(t *testing.T) {
	ctx := exampleContext()
	rows := generateDeviceData()

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	// Explicitly set column order (default is insertion order)
	builder.TextField("name", "device.name", func(r DeviceTableRow) string {
		return r.Name
	}).WithColumnOrder(2) // Third column

	builder.IdField("id", "device.id", func(r DeviceTableRow) int64 {
		return r.ID
	}).WithColumnOrder(1) // Second column

	statusIcons := table.NewIconSet()
	statusIcons.Add("online", "check_circle", table.FieldColorAccent, "Online")
	statusIcons.Add("offline", "cancel", table.FieldColorWarning, "Offline")

	builder.IconFieldFromSet("status", "device.status",
		func(r DeviceTableRow) *table.IconRef {
			return statusIcons.Resolve(r.Status)
		},
		statusIcons,
	).WithColumnOrder(0) // First column

	tbl := builder.Build()
	tbl.SetData(rows)
	fields := tbl.Print(exampleTranslator)["data"].(map[string]any)["fields"].([]map[string]any)

	// Verify column order (fields are sorted by columnOrder)
	fmt.Printf("Column order:\n")
	for i, field := range fields {
		fmt.Printf("  Position %d: %s\n", i, field["id"])
	}
	// Expected: status (0), id (1), name (2)

	// Column ordering:
	// - Fields are sorted by columnOrder before display
	// - Default: insertion order (auto-assigned 0, 1, 2, ...)
	// - Explicit ordering allows reordering without changing code flow
	// - Useful for dynamic column visibility/reordering
}

// TestAdvanced08_AccessControl demonstrates permission-based field visibility
func TestAdvanced08_AccessControl(t *testing.T) {
	ctx := exampleContext()
	rows := generateDeviceData()

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	// Public field (no access control)
	builder.TextField("name", "device.name", func(r DeviceTableRow) string {
		return r.Name
	})

	// Admin-only field
	builder.IdField("id", "device.id", func(r DeviceTableRow) int64 {
		return r.ID
	}).WithAccess([]string{"admin", "superadmin"})

	// Edit permission required
	builder.InputField("notes", "device.notes", func(r DeviceTableRow) any {
		return r.Notes
	}).WithAccess([]string{"device.edit"}).WithInputType("text")

	tbl := builder.Build()
	tbl.SetData(rows)
	fields := tbl.Print(exampleTranslator)["data"].(map[string]any)["fields"].([]map[string]any)

	// Access control is set in field definition
	// Frontend checks user permissions and hides fields accordingly
	for _, field := range fields {
		if access, ok := field["access"].([]string); ok && len(access) > 0 {
			fmt.Printf("Field %s requires permissions: %v\n", field["id"], access)
		} else {
			fmt.Printf("Field %s is public\n", field["id"])
		}
	}

	// Access control:
	// - Field-level permissions checked by frontend
	// - Backend should also enforce permissions on data
	// - Multiple permissions = OR logic (any permission grants access)
	// - No access set = public field
}

// ===== REAL CONTROLLER INTEGRATION EXAMPLES =====
// These examples demonstrate patterns used in actual controller implementations

// TestIntegration01_AJAXTable demonstrates AJAX table with URL endpoint
func TestIntegration01_AJAXTable(t *testing.T) {
	ctx := exampleContext()

	// In controller: no data provided, URL is set for AJAX loading
	url := xurl.NewUrl("/Portal/Device/TableData")

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	builder.IdField("id", "device.id", func(r DeviceTableRow) int64 {
		return r.ID
	}).WithWidth("80px")

	builder.TextField("name", "device.name", func(r DeviceTableRow) string {
		return r.Name
	})

	builder.DistanceField("distance", "device.distance", func(r DeviceTableRow) float64 {
		return r.DistanceKm
	}).WithDecimals(1)

	// Build AJAX table with URL but no data
	tbl := builder.Build()
	tbl.SetURL(url)

	output := tbl.Print(exampleTranslator)
	dataSection := output["data"].(map[string]any)

	// Verify AJAX mode
	if dataSection["url"] != nil {
		fmt.Printf("Table is in AJAX mode with URL: %v\n", dataSection["url"])
	}
	if dataSection["data"] == nil {
		fmt.Printf("Data is nil (will be loaded via AJAX)\n")
	}

	// AJAX tables:
	// - Frontend makes separate request to URL for data
	// - Supports server-side pagination, sorting, filtering
	// - Reduces initial page load time
	// - URL endpoint returns same table.GetData(OutputWeb) format
}

// TestIntegration02_FormGroupFilters demonstrates tables with filter forms
func TestIntegration02_FormGroupFilters(t *testing.T) {
	ctx := exampleContext()
	rows := generateDeviceData()

	// Create filter form (would typically come from formgroup package)
	searchField := field.NewTextField("search", "common.search", false, "")
	searchField.Subtype = "text"

	filters := []field.FormField{
		searchField,
		field.NewBoolField("active_only", "device.active_only", false, false),
		field.NewTimeRangeFieldWithDefault("daterange", "common.daterange", false, 7),
	}

	fg, err := group.NewFormGroupWithContext(filters, ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Build table with filter form
	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	builder.IdField("id", "device.id", func(r DeviceTableRow) int64 {
		return r.ID
	})

	builder.TextField("name", "device.name", func(r DeviceTableRow) string {
		return r.Name
	})

	builder.DateTimeField("last_seen", "device.last_seen", func(r DeviceTableRow) time.Time {
		return r.LastSeen
	})

	tbl := builder.SetFilter(fg).Build()
	tbl.SetData(rows)

	output := tbl.Print(exampleTranslator)
	dataSection := output["data"].(map[string]any)

	// Verify filter is attached
	if hasFilter, ok := dataSection["hasFilter"].(bool); ok && hasFilter {
		fmt.Printf("Table has filter form attached\n")
	}

	// Tables with filters:
	// - FormGroup defines available filter fields
	// - Frontend renders filter UI above table
	// - Filter values sent with AJAX requests
	// - Backend parses filter values from request
	// - Use fg.ParseAndValidate(requestParams) to get typed values
}

// TestIntegration03_FooterAggregations demonstrates sum and count footers
func TestIntegration03_FooterAggregations(t *testing.T) {
	ctx := exampleContext()
	rows := generateDeviceData()

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	// No footer
	builder.TextField("name", "device.name", func(r DeviceTableRow) string {
		return r.Name
	})

	// Count footer (counts non-empty values)
	builder.IdField("id", "device.id", func(r DeviceTableRow) int64 {
		return r.ID
	}).WithFooter(table.FieldFooterCount)

	// Sum footer (adds numeric values)
	builder.DistanceField("distance", "device.distance", func(r DeviceTableRow) float64 {
		return r.DistanceKm
	}).
		WithDecimals(1).
		WithFooter(table.FieldFooterSum)

	tbl := builder.Build()
	tbl.SetData(rows)
	fields := tbl.Print(exampleTranslator)["data"].(map[string]any)["fields"].([]map[string]any)

	// Verify footer settings
	for _, field := range fields {
		if footer, ok := field["footer"].(string); ok {
			fmt.Printf("Field %s has footer: %s\n", field["id"], footer)
		}
	}

	// Footer aggregations:
	// - Displayed at bottom of table
	// - FieldFooterSum: Adds all numeric values
	// - FieldFooterCount: Counts non-null rows
	// - FieldFooterStatic: Custom static text
	// - Frontend calculates from visible (filtered) data
}

// TestIntegration04_EmptyStateHandling demonstrates handling empty data
func TestIntegration04_EmptyStateHandling(t *testing.T) {
	ctx := exampleContext()

	// Empty data slice
	var emptyRows []DeviceTableRow

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	builder.IdField("id", "device.id", func(r DeviceTableRow) int64 {
		return r.ID
	})

	builder.TextField("name", "device.name", func(r DeviceTableRow) string {
		return r.Name
	})

	// Set custom empty message
	tbl := builder.
		SetTextNoData("No devices found. Please adjust your filters or add devices.").
		Build()
	tbl.SetData(emptyRows)

	data := tbl.GetData(table.OutputWeb)
	options := tbl.Print(exampleTranslator)["data"].(map[string]any)["options"].(map[string]any)

	// Verify empty data handling
	fmt.Printf("Table has %d rows\n", len(data))
	if noDataText, ok := options["textNoData"].(string); ok {
		fmt.Printf("Custom empty message: %s\n", noDataText)
	}

	// Empty state best practices:
	// - Always provide helpful "no data" message
	// - Suggest actions (adjust filters, add data, check permissions)
	// - Distinguish between "no results" vs "loading" vs "error"
	// - Use TextNoData option for custom message
}

// TestIntegration05_MultiFormatOutput demonstrates CSV/PDF/Excel export
func TestIntegration05_MultiFormatOutput(t *testing.T) {
	ctx := exampleContext()
	rows := generateDeviceData()

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	builder.IdField("id", "device.id", func(r DeviceTableRow) int64 {
		return r.ID
	})

	builder.TextField("name", "device.name", func(r DeviceTableRow) string {
		return r.Name
	})

	// Distance field (formatted differently for web vs export)
	builder.DistanceField("distance", "device.distance", func(r DeviceTableRow) float64 {
		return r.DistanceKm
	}).WithDecimals(1)

	// Action buttons (excluded from CSV export)
	builder.ButtonsField("actions", "device.actions", func(r DeviceTableRow) map[string]string {
		return map[string]string{
			"0": fmt.Sprintf("/Portal/Device/Edit?id=%d", r.ID),
		}
	}).
		AddButton(0, table.FieldButtonActionLink, "edit", table.FieldColorPrimary, "Edit device").
		HideInCSV() // Exclude from CSV/PDF/Excel

	tbl := builder.Build()
	tbl.SetData(rows)

	// Get data in different formats
	webData := tbl.GetData(table.OutputWeb)
	csvData := tbl.GetData(table.OutputCSV)

	// Web output: distance is array [display, value] for sorting
	fmt.Printf("Web format distance: %v\n", webData[0]["distance"])
	// Output: ["196,5 km", 196.5]

	// CSV output: distance is plain text for export
	fmt.Printf("CSV format distance: %v\n", csvData[0]["distance"])
	// Output: "196,5 km"

	// CSV output: actions field excluded
	if _, hasActions := csvData[0]["actions"]; !hasActions {
		fmt.Printf("Actions field excluded from CSV\n")
	}

	// Multi-format output:
	// - OutputWeb: Arrays for sortable numbers, includes all fields
	// - OutputCSV: Plain text values, respects WithCsv(false)
	// - OutputPDF: Same as CSV
	// - OutputExcel: Same as CSV
	// - Use same table definition for all formats
}

// TestIntegration06_DynamicFieldVisibility demonstrates permission-based columns
func TestIntegration06_DynamicFieldVisibility(t *testing.T) {
	ctx := exampleContext()
	rows := generateDeviceData()

	// Simulate user permissions (would come from auth context)
	userPermissions := map[string]bool{
		"device.view":   true,
		"device.edit":   false, // User cannot edit
		"device.delete": false, // User cannot delete
		"admin":         false, // Not admin
	}

	hasPermission := func(required []string) bool {
		for _, perm := range required {
			if userPermissions[perm] {
				return true // OR logic: any permission grants access
			}
		}
		return false
	}

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	// Public field (always visible)
	builder.TextField("name", "device.name", func(r DeviceTableRow) string {
		return r.Name
	})

	// Admin-only field
	idField := builder.IdField("id", "device.id", func(r DeviceTableRow) int64 {
		return r.ID
	}).WithAccess([]string{"admin"})

	// Conditionally hide based on permissions
	if !hasPermission([]string{"admin"}) {
		idField.Hide()
	}

	// Edit button (only if user has edit permission)
	if hasPermission([]string{"device.edit"}) {
		builder.ButtonsField("edit_button", "device.edit", func(r DeviceTableRow) map[string]string {
			return map[string]string{
				"0": fmt.Sprintf("/Portal/Device/Edit?id=%d", r.ID),
			}
		}).AddButton(0, table.FieldButtonActionLink, "edit", table.FieldColorPrimary, "Edit")
	}

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(table.OutputWeb)

	// Verify hidden field doesn't appear in output
	if _, hasID := data[0]["id"]; !hasID {
		fmt.Printf("ID field hidden (user lacks admin permission)\n")
	}

	if _, hasEdit := data[0]["edit_button"]; !hasEdit {
		fmt.Printf("Edit button hidden (user lacks edit permission)\n")
	}

	// Dynamic field visibility:
	// - Check user permissions in controller
	// - Use WithHide(true) to hide fields without permission
	// - Or conditionally add fields to builder
	// - Both approaches work, choose based on complexity
	// - Always enforce permissions on backend data access too
}

// ===== JSON OUTPUT TESTS =====
// These tests output complete JSON structures for frontend compatibility validation

// TestJSON_ComponentStructure outputs the complete table component JSON structure
// This is what controllers send to the frontend via Table.Print()
func TestJSON_ComponentStructure(t *testing.T) {
	ctx := exampleContext()
	rows := generateDeviceData()

	// Create filter form for demonstration
	searchField := field.NewTextField("search", "common.search", false, "")
	searchField.Subtype = "text"
	filters := []field.FormField{
		searchField,
		field.NewTimeRangeFieldWithDefault("daterange", "common.daterange", false, 7),
	}
	fg, _ := group.NewFormGroupWithContext(filters, ctx)

	// Build comprehensive table with various field types
	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	// Basic fields
	builder.IdField("id", "device.id", func(r DeviceTableRow) int64 {
		return r.ID
	}).WithWidth("80px").WithFooter(table.FieldFooterCount)

	builder.TextField("name", "device.name", func(r DeviceTableRow) string {
		return r.Name
	}).WithMinWidth("200px")

	// Icon field with status mapping
	statusIcons := table.NewIconSet()
	statusIcons.Add("online", "check_circle", table.FieldColorAccent, "Device is online")
	statusIcons.Add("offline", "cancel", table.FieldColorWarning, "Device is offline")

	builder.IconFieldFromSet("status", "device.status",
		func(r DeviceTableRow) *table.IconRef {
			return statusIcons.Resolve(r.Status)
		},
		statusIcons,
	).WithAlign(table.FieldAlignCenter)

	// Distance field with unit conversion
	builder.DistanceField("distance", "device.distance", func(r DeviceTableRow) float64 {
		return r.DistanceKm
	}).WithDecimals(1).WithFooter(table.FieldFooterSum)

	// Buttons field with actions
	builder.ButtonsField("actions", "device.actions", func(r DeviceTableRow) map[string]string {
		// Return button data as map: button index -> URL (empty string to hide)
		return map[string]string{
			"0": fmt.Sprintf("/Portal/Device/Edit?id=%d", r.ID),   // Edit button
			"1": fmt.Sprintf("/Portal/Device/Delete?id=%d", r.ID), // Delete button
		}
	}).
		AddButton(0, table.FieldButtonActionLink, "edit", table.FieldColorPrimary, "Edit device").
		AddButton(1, table.FieldButtonActionDialog, "delete", table.FieldColorWarning, "Delete device").
		HideInCSV()

	// Table options for demonstration
	itemsPerPage := 25
	dense := true
	pagination := true
	search := true

	// ========================================
	// EXAMPLE 1: AJAX Mode (URL set, data nil)
	// ========================================
	fmt.Println("\n=== AJAX MODE: Component Structure (Print output) ===")
	fmt.Println("This is what controllers return for AJAX tables (data loaded separately)")

	ajaxUrl := xurl.NewUrl("/Portal/Device/TableData")
	ajaxTable := builder.
		SetFilter(fg).
		SetItemsPerPage(itemsPerPage).
		SetDense(dense).
		SetPagination(pagination).
		SetSearch(search).
		Build()
	ajaxTable.SetURL(ajaxUrl)

	ajaxOutput := ajaxTable.Print(exampleTranslator)
	ajaxJSON, _ := json.MarshalIndent(ajaxOutput, "", "  ")
	fmt.Println(string(ajaxJSON))

	// ========================================
	// EXAMPLE 2: Static Mode (data embedded)
	// ========================================
	fmt.Println("\n=== STATIC MODE: Component Structure (Print output) ===")
	fmt.Println("This is what controllers return for static tables (data embedded)")

	staticTable := builder.
		SetFilter(fg).
		SetItemsPerPage(itemsPerPage).
		SetDense(dense).
		SetPagination(pagination).
		SetSearch(search).
		Build()
	staticTable.SetData(rows)

	staticOutput := staticTable.Print(exampleTranslator)
	staticJSON, _ := json.MarshalIndent(staticOutput, "", "  ")
	fmt.Println(string(staticJSON))

	// Key differences between AJAX and Static modes:
	// AJAX: data["url"] = string, data["data"] = null
	// Static: data["url"] = null, data["data"] = array of rows
	//
	// Component structure (Print output) includes:
	// - type: "table"
	// - display: CSS class for layout
	// - data:
	//   - fields: Array of field definitions (id, name, format, width, etc.)
	//   - options: Table configuration (pagination, search, dense, etc.)
	//   - hasFilter: Boolean indicating if filter form attached
	//   - url: AJAX endpoint URL (AJAX mode only)
	//   - data: Array of formatted rows (Static mode only)
}

// TestJSON_TableDataStructure outputs the table data array structure
// This is what AJAX endpoints return via Table.GetData(OutputWeb)
func TestJSON_TableDataStructure(t *testing.T) {
	ctx := exampleContext()
	rows := generateDeviceData()

	// Build table with various field types
	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	builder.IdField("id", "device.id", func(r DeviceTableRow) int64 {
		return r.ID
	})

	builder.TextField("name", "device.name", func(r DeviceTableRow) string {
		return r.Name
	})

	// Icon field - returns mapped icon value
	statusIcons := table.NewIconSet()
	statusIcons.Add("online", "check_circle", table.FieldColorAccent, "Device is online")
	statusIcons.Add("offline", "cancel", table.FieldColorWarning, "Device is offline")

	builder.IconFieldFromSet("status", "device.status",
		func(r DeviceTableRow) *table.IconRef {
			return statusIcons.Resolve(r.Status)
		},
		statusIcons,
	)

	// Distance field - returns [display, value] array for sorting
	builder.DistanceField("distance", "device.distance", func(r DeviceTableRow) float64 {
		return r.DistanceKm
	}).WithDecimals(1)

	// Speed field - also returns [display, value] array
	builder.SpeedField("speed", "device.speed", func(r DeviceTableRow) float64 {
		return r.SpeedKmh
	}).WithDecimals(1)

	// DateTime field - returns formatted date/time string
	builder.DateTimeField("last_seen", "device.last_seen", func(r DeviceTableRow) time.Time {
		return r.LastSeen
	})

	// Buttons field - returns map of button index to URL
	builder.ButtonsField("actions", "device.actions", func(r DeviceTableRow) map[string]string {
		// Return button data as map: button index -> URL (empty string to hide)
		return map[string]string{
			"0": fmt.Sprintf("/Portal/Device/Edit?id=%d", r.ID),   // Edit button
			"1": fmt.Sprintf("/Portal/Device/Delete?id=%d", r.ID), // Delete button
		}
	}).
		AddButton(0, table.FieldButtonActionLink, "edit", table.FieldColorPrimary, "Edit").
		AddButton(1, table.FieldButtonActionDialog, "delete", table.FieldColorWarning, "Delete")

	tbl := builder.Build()
	tbl.SetData(rows)

	// ========================================
	// Table Data Array (GetData output)
	// ========================================
	fmt.Println("\n=== TABLE DATA: Array Structure (GetData output) ===")
	fmt.Println("This is what AJAX endpoints return as JSON response")

	data := tbl.GetData(table.OutputWeb)
	dataJSON, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(dataJSON))

	// Critical formatting rules for frontend compatibility:
	//
	// 1. NUMBER FIELDS (Integer, Float, Distance, Speed, Pressure):
	//    - Web output: [displayString, numericValue]
	//    - Example: ["196,5 km", 196.5]
	//    - Frontend sorts by numericValue, displays displayString
	//
	// 2. ICON FIELDS:
	//    - Returns the mapped icon value (string)
	//    - Example: "online" -> "check_circle"
	//
	// 3. BUTTON FIELDS:
	//    - Returns map with button index as key, URL or false as value
	//    - Example: {"0": "/Portal/Device/Edit?id=1", "1": "/Portal/Device/Delete?id=1"}
	//    - Use false to hide a button: {"0": "/url", "1": false}
	//
	// 4. DATETIME FIELDS:
	//    - Returns Unix timestamp in SECONDS (not milliseconds!)
	//    - Example: 1609459200
	//
	// 5. TEXT FIELDS:
	//    - Returns plain string
	//    - Example: "Device Name"
	//
	// This array structure is what the Angular frontend (@xiri/xiri-ui v20.0.3)
	// expects for rendering table rows with proper sorting, formatting, and actions.
}

// ============================================================================
// ICON FIELD FROM SET EXAMPLES
// ============================================================================

// TestExample_IconFieldFromSet demonstrates the type-safe IconFieldFromSet method
// using direct IconRef variables. This approach prevents typos at compile time.
func TestExample_IconFieldFromSet(t *testing.T) {
	ctx := exampleContext()
	rows := generateDeviceData()

	// Step 1: Define all possible icons upfront in an IconSet
	statusIcons := table.NewIconSet()
	iconOnline := statusIcons.Add("online", "check_circle", table.FieldColorAccent, "Device is online")
	iconOffline := statusIcons.Add("offline", "cancel", table.FieldColorWarning, "Device is offline")

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	builder.TextField("name", "device.name", func(r DeviceTableRow) string {
		return r.Name
	})

	// Step 2: Use IconFieldFromSet with direct IconRef variables
	// The accessor returns *IconRef — only registered icons are possible
	builder.IconFieldFromSet("status", "device.status",
		func(r DeviceTableRow) *table.IconRef {
			if r.Online {
				return iconOnline
			}
			return iconOffline
		},
		statusIcons,
	)

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(table.OutputWeb)

	// Output values are identical to IconField + AddIcon
	fmt.Printf("Device 1 (online): status=%v\n", data[0]["status"])
	fmt.Printf("Device 2 (offline): status=%v\n", data[1]["status"])

	// Advantage: A typo like "onlne" is impossible because only
	// iconOnline and iconOffline variables exist — no free strings.
}

// TestExample_IconFieldFromSet_Resolve demonstrates using Resolve() for string data
// from external sources (e.g., database). Resolve returns nil for unknown values.
func TestExample_IconFieldFromSet_Resolve(t *testing.T) {
	ctx := exampleContext()
	rows := generateDeviceData()

	statusIcons := table.NewIconSet()
	statusIcons.Add("online", "check_circle", table.FieldColorAccent, "Device is online")
	statusIcons.Add("offline", "cancel", table.FieldColorWarning, "Device is offline")

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	builder.TextField("name", "device.name", func(r DeviceTableRow) string {
		return r.Name
	})

	// Use Resolve() when the source data is a string field (e.g., from DB)
	builder.IconFieldFromSet("status", "device.status",
		func(r DeviceTableRow) *table.IconRef {
			return statusIcons.Resolve(r.Status) // nil for unknown values
		},
		statusIcons,
	)

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(table.OutputWeb)

	fmt.Printf("Device 1 (status=%s): %v\n", rows[0].Status, data[0]["status"])
	fmt.Printf("Device 2 (status=%s): %v\n", rows[1].Status, data[1]["status"])
}

// ============================================================================
// N-LINE (VARIABLE LENGTH) FIELD EXAMPLES
// ============================================================================

// TestExample_TextNField demonstrates variable-line text field
func TestExample_TextNField(t *testing.T) {
	ctx := exampleContext()

	rows := []DeviceTableRow{
		{ID: 1, Name: "Device 1", NameLine2: "Fleet A", GroupName: "Primary"},
	}

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	builder.TextNField("info", "device.info", func(r DeviceTableRow) []string {
		return []string{r.Name, r.NameLine2, r.GroupName}
	})

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(table.OutputWeb)

	info := data[0]["info"].([]string)
	fmt.Printf("TextNField output: %v\n", info)
	// Output: [Device 1 Fleet A Primary]

	if len(info) != 3 {
		t.Errorf("expected 3 lines, got %d", len(info))
	}
	if info[0] != "Device 1" {
		t.Errorf("expected 'Device 1', got %q", info[0])
	}

	csvData := tbl.GetData(table.OutputCSV)
	csvInfo := csvData[0]["info"].([]string)
	fmt.Printf("TextNField CSV: %v\n", csvInfo)
	// Output: ["Device 1", "Fleet A", "Primary"]
	if len(csvInfo) != 3 {
		t.Errorf("expected 3 values, got %d", len(csvInfo))
	}
	if csvInfo[0] != "Device 1" {
		t.Errorf("expected 'Device 1', got %q", csvInfo[0])
	}
	if csvInfo[1] != "Fleet A" {
		t.Errorf("expected 'Fleet A', got %q", csvInfo[1])
	}
	if csvInfo[2] != "Primary" {
		t.Errorf("expected 'Primary', got %q", csvInfo[2])
	}
}

// TestExample_IntNField demonstrates variable-line integer field
func TestExample_IntNField(t *testing.T) {
	ctx := exampleContext() // German locale

	rows := []DeviceTableRow{
		{ID: 1, TripsToday: 5, TotalTrips: 1234},
	}

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	builder.IntNField("stats", "device.stats", func(r DeviceTableRow) []int {
		return []int{r.TripsToday, r.TotalTrips, 0}
	})

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(table.OutputWeb)

	stats := data[0]["stats"].([]string)
	fmt.Printf("IntNField output: %v\n", stats)

	if len(stats) != 3 {
		t.Errorf("expected 3 lines, got %d", len(stats))
	}

	csvData := tbl.GetData(table.OutputCSV)
	csvStats := csvData[0]["stats"].([]string)
	fmt.Printf("IntNField CSV: %v\n", csvStats)
	// Output: ["5", "1234", "0"]
	if len(csvStats) != 3 {
		t.Errorf("expected 3 values, got %d", len(csvStats))
	}
	if csvStats[0] != "5" {
		t.Errorf("expected '5', got %q", csvStats[0])
	}
	if csvStats[1] != "1234" {
		t.Errorf("expected '1234', got %q", csvStats[1])
	}
}

// TestExample_FloatNField demonstrates variable-line float field
func TestExample_FloatNField(t *testing.T) {
	ctx := exampleContext()

	rows := []DeviceTableRow{
		{ID: 1, FuelCurrent: 45.5, FuelAverage: 52.3},
	}

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	builder.FloatNField("fuel", "device.fuel", func(r DeviceTableRow) []float64 {
		return []float64{r.FuelCurrent, r.FuelAverage}
	}).WithDecimals(1)

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(table.OutputWeb)

	fuel := data[0]["fuel"].([]string)
	fmt.Printf("FloatNField output: %v\n", fuel)

	if len(fuel) != 2 {
		t.Errorf("expected 2 lines, got %d", len(fuel))
	}

	csvData := tbl.GetData(table.OutputCSV)
	csvFuel := csvData[0]["fuel"].([]string)
	fmt.Printf("FloatNField CSV: %v\n", csvFuel)
	if len(csvFuel) != 2 {
		t.Errorf("expected 2 values, got %d", len(csvFuel))
	}
	if csvFuel[0] != "45.5" {
		t.Errorf("expected '45.5', got %q", csvFuel[0])
	}
}

// TestExample_DateTimeNField demonstrates variable-line datetime field
func TestExample_DateTimeNField(t *testing.T) {
	ctx := exampleContext()

	lastSeen := time.Unix(1640000000, 0)
	created := time.Unix(1630000000, 0)
	updated := time.Unix(1635000000, 0)

	rows := []DeviceTableRow{
		{ID: 1, LastSeen: lastSeen, CreatedDate: created},
	}

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	builder.DateTimeNField("times", "device.times", func(r DeviceTableRow) []time.Time {
		return []time.Time{r.LastSeen, r.CreatedDate, updated}
	})

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(table.OutputWeb)

	times := data[0]["times"].([]string)
	fmt.Printf("DateTimeNField output: %v\n", times)

	if len(times) != 3 {
		t.Errorf("expected 3 lines, got %d", len(times))
	}
	for i, v := range times {
		if v == "" {
			t.Errorf("expected non-empty datetime at index %d", i)
		}
	}
}

// TestExample_DateNField demonstrates variable-line date field
func TestExample_DateNField(t *testing.T) {
	ctx := exampleContext()

	date1 := time.Unix(1640000000, 0)
	date2 := time.Unix(1630000000, 0)

	rows := []DeviceTableRow{
		{ID: 1, LastSeen: date1, CreatedDate: date2},
	}

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	builder.DateNField("dates", "device.dates", func(r DeviceTableRow) []time.Time {
		return []time.Time{r.LastSeen, r.CreatedDate}
	})

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(table.OutputWeb)

	dates := data[0]["dates"].([]string)
	fmt.Printf("DateNField output: %v\n", dates)

	if len(dates) != 2 {
		t.Errorf("expected 2 lines, got %d", len(dates))
	}
}

// TestExample_DistanceNField demonstrates variable-line distance field
func TestExample_DistanceNField(t *testing.T) {
	ctx := exampleContext()

	rows := []DeviceTableRow{
		{ID: 1, TodayKm: 123.45, TotalKm: 9876.54},
	}

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	builder.DistanceNField("distances", "device.distances", func(r DeviceTableRow) []float64 {
		return []float64{r.TodayKm, r.TotalKm, 42.0}
	}).WithDecimals(1)

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(table.OutputWeb)

	distances := data[0]["distances"].([]string)
	fmt.Printf("DistanceNField output: %v\n", distances)

	if len(distances) != 3 {
		t.Errorf("expected 3 lines, got %d", len(distances))
	}

	csvData := tbl.GetData(table.OutputCSV)
	csvDistances := csvData[0]["distances"].([]string)
	fmt.Printf("DistanceNField CSV: %v\n", csvDistances)
	if len(csvDistances) != 3 {
		t.Errorf("expected 3 values, got %d", len(csvDistances))
	}
}

// TestExample_SpeedNField demonstrates variable-line speed field
func TestExample_SpeedNField(t *testing.T) {
	ctx := exampleContext()

	rows := []DeviceTableRow{
		{ID: 1, MaxSpeed: 120.5, AvgSpeed: 85.3},
	}

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	builder.SpeedNField("speeds", "device.speeds", func(r DeviceTableRow) []float64 {
		return []float64{r.MaxSpeed, r.AvgSpeed}
	}).WithDecimals(1)

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(table.OutputWeb)

	speeds := data[0]["speeds"].([]string)
	fmt.Printf("SpeedNField output: %v\n", speeds)

	if len(speeds) != 2 {
		t.Errorf("expected 2 lines, got %d", len(speeds))
	}

	csvData := tbl.GetData(table.OutputCSV)
	csvSpeeds := csvData[0]["speeds"].([]string)
	fmt.Printf("SpeedNField CSV: %v\n", csvSpeeds)
	if len(csvSpeeds) != 2 {
		t.Errorf("expected 2 values, got %d", len(csvSpeeds))
	}
}

// TestExample_BoolNField demonstrates variable-line boolean field
func TestExample_BoolNField(t *testing.T) {
	ctx := exampleContext()

	rows := []DeviceTableRow{
		{ID: 1, Active: true, Online: false},
	}

	builder := table.NewBuilder[DeviceTableRow](ctx, exampleTranslator)

	builder.BoolNField("flags", "device.flags", func(r DeviceTableRow) []bool {
		return []bool{r.Active, r.Online, true}
	})

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(table.OutputWeb)

	flags := data[0]["flags"].([]string)
	fmt.Printf("BoolNField output: %v\n", flags)

	if len(flags) != 3 {
		t.Errorf("expected 3 lines, got %d", len(flags))
	}
	if flags[0] != "Yes" {
		t.Errorf("expected 'Yes', got %q", flags[0])
	}
	if flags[1] != "No" {
		t.Errorf("expected 'No', got %q", flags[1])
	}

	csvData := tbl.GetData(table.OutputCSV)
	csvFlags := csvData[0]["flags"].([]string)
	fmt.Printf("BoolNField CSV: %v\n", csvFlags)
	if len(csvFlags) != 3 {
		t.Errorf("expected 3 values, got %d", len(csvFlags))
	}
	if csvFlags[0] != "Yes" {
		t.Errorf("expected 'Yes', got %q", csvFlags[0])
	}
	if csvFlags[1] != "No" {
		t.Errorf("expected 'No', got %q", csvFlags[1])
	}
	if csvFlags[2] != "Yes" {
		t.Errorf("expected 'Yes', got %q", csvFlags[2])
	}
}

// TestExample_TimeLengthNField demonstrates variable-line time duration field
func TestExample_TimeLengthNField(t *testing.T) {
	ctx := exampleContext()

	rows := []TripTableRow{
		{
			ID:             1,
			DurationShort:  2700,   // 45 minutes
			DurationMedium: 19800,  // 5 hours 30 minutes
			DurationLong:   183900, // 2 days 3 hours 5 minutes
		},
	}

	builder := table.NewBuilder[TripTableRow](ctx, exampleTranslator)

	builder.TimeLengthNField("durations", "trip.durations", func(r TripTableRow) []int64 {
		return []int64{r.DurationShort, r.DurationMedium, r.DurationLong}
	})

	tbl := builder.Build()
	tbl.SetData(rows)
	data := tbl.GetData(table.OutputWeb)

	durations := data[0]["durations"].([]string)
	fmt.Printf("TimeLengthNField output: %v\n", durations)

	if len(durations) != 3 {
		t.Errorf("expected 3 lines, got %d", len(durations))
	}
	if durations[0] != "00:45" {
		t.Errorf("expected '00:45', got %q", durations[0])
	}
	if durations[1] != "05:30" {
		t.Errorf("expected '05:30', got %q", durations[1])
	}
	if durations[2] != "2d 03:05" {
		t.Errorf("expected '2d 03:05', got %q", durations[2])
	}

	csvData := tbl.GetData(table.OutputCSV)
	csvDurations := csvData[0]["durations"].([]string)
	fmt.Printf("TimeLengthNField CSV: %v\n", csvDurations)
	// Output: ["45", "330", "3065"]
	if len(csvDurations) != 3 {
		t.Errorf("expected 3 values, got %d", len(csvDurations))
	}
	if csvDurations[0] != "45" {
		t.Errorf("expected '45', got %q", csvDurations[0])
	}
	if csvDurations[1] != "330" {
		t.Errorf("expected '330', got %q", csvDurations[1])
	}
	if csvDurations[2] != "3065" {
		t.Errorf("expected '3065', got %q", csvDurations[2])
	}
}
