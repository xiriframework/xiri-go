package table

import (
	"testing"

	xurl "github.com/xiriframework/xiri-go/component/url"
	"github.com/xiriframework/xiri-go/types/distance"
	"github.com/xiriframework/xiri-go/types/language"
	"github.com/xiriframework/xiri-go/types/locale"
	"github.com/xiriframework/xiri-go/types/timezone"
	"github.com/xiriframework/xiri-go/uicontext"
)

// Test row struct
type testDeviceRow struct {
	ID     int64
	Name   string
	Active bool
}

// Test translator
func testTranslator(key string) string {
	translations := map[string]string{
		"device.id":     "ID",
		"device.name":   "Device Name",
		"device.active": "Active",
	}
	if val, ok := translations[key]; ok {
		return val
	}
	return key
}

// Test context
func testContext() *uicontext.UiContext {
	return &uicontext.UiContext{
		Timezone: timezone.EuropeVienna,
		Lang:     language.Deutsch,
		Locale:   locale.De,
		Distance: distance.Kilometer,
	}
}

// TestPrintAJAXMode verifies Print() returns correct structure for AJAX tables (url != nil)
func TestPrintAJAXMode(t *testing.T) {
	ctx := testContext()

	url := xurl.NewUrl("/Portal/Device/TableData")

	builder := NewBuilder[testDeviceRow](ctx, testTranslator)
	builder.IdField("id", "device.id", func(r testDeviceRow) int64 { return r.ID })
	builder.TextField("name", "device.name", func(r testDeviceRow) string { return r.Name })
	tbl := builder.Build()
	tbl.SetURL(url)

	// Call Print() - should produce AJAX mode JSON
	output := tbl.Print(testTranslator)

	// Verify structure
	if output["type"] != "table" {
		t.Errorf("Expected type 'table', got %v", output["type"])
	}

	data, ok := output["data"].(map[string]any)
	if !ok {
		t.Fatalf("Expected data to be map[string]any, got %T", output["data"])
	}

	// AJAX mode: url should be set, data should be nil
	if data["url"] == nil {
		t.Error("Expected url to be set in AJAX mode")
	}
	if data["data"] != nil {
		t.Error("Expected data to be nil in AJAX mode")
	}

	// Verify fields array exists
	fields, ok := data["fields"].([]map[string]any)
	if !ok {
		t.Fatalf("Expected fields to be []map[string]any, got %T", data["fields"])
	}
	if len(fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(fields))
	}

	// Verify options exist
	options, ok := data["options"].(map[string]any)
	if !ok {
		t.Fatalf("Expected options to be map[string]any, got %T", data["options"])
	}
	_ = options
}

// TestPrintStaticMode verifies Print() returns correct structure for static tables (url == nil)
func TestPrintStaticMode(t *testing.T) {
	ctx := testContext()

	rows := []testDeviceRow{
		{ID: 1, Name: "Device 1", Active: true},
		{ID: 2, Name: "Device 2", Active: false},
	}

	builder := NewBuilder[testDeviceRow](ctx, testTranslator)
	builder.IdField("id", "device.id", func(r testDeviceRow) int64 { return r.ID })
	builder.TextField("name", "device.name", func(r testDeviceRow) string { return r.Name })
	builder.BoolField("active", "device.active", func(r testDeviceRow) bool { return r.Active }).
		WithBoolText("Yes", "No")
	tbl := builder.Build()
	tbl.SetData(rows)

	// Call Print() - should produce static mode JSON
	output := tbl.Print(testTranslator)

	// Verify structure
	if output["type"] != "table" {
		t.Errorf("Expected type 'table', got %v", output["type"])
	}

	data, ok := output["data"].(map[string]any)
	if !ok {
		t.Fatalf("Expected data to be map[string]any, got %T", output["data"])
	}

	// Static mode: url should be nil, data should be array
	if data["url"] != nil {
		t.Error("Expected url to be nil in static mode")
	}

	rowData, ok := data["data"].([]map[string]any)
	if !ok {
		t.Fatalf("Expected data to be []map[string]any, got %T", data["data"])
	}
	if len(rowData) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(rowData))
	}

	// Verify row structure
	if rowData[0]["id"] == nil {
		t.Error("Expected id field in row 0")
	}
	if rowData[0]["name"] == nil {
		t.Error("Expected name field in row 0")
	}
	if rowData[0]["active"] != "Yes" {
		t.Errorf("Expected active='Yes' for row 0, got %v", rowData[0]["active"])
	}
	if rowData[1]["active"] != "No" {
		t.Errorf("Expected active='No' for row 1, got %v", rowData[1]["active"])
	}
}

// TestPrintWithTranslator verifies translator is properly passed to fields
func TestPrintWithTranslator(t *testing.T) {
	ctx := testContext()

	translator := func(key string) string {
		if key == "device.id" {
			return "Translated ID"
		}
		return key
	}

	builder := NewBuilder[testDeviceRow](ctx, testTranslator)
	builder.IdField("id", "device.id", func(r testDeviceRow) int64 { return r.ID })
	tbl := builder.Build()
	tbl.SetData([]testDeviceRow{{ID: 1}})

	output := tbl.Print(translator)

	data := output["data"].(map[string]any)
	fields := data["fields"].([]map[string]any)

	// Verify translation was applied
	if fields[0]["name"] != "Translated ID" {
		t.Errorf("Expected field name to be 'Translated ID', got %v", fields[0]["name"])
	}
}

// TestPrintFallbackTranslator verifies fallback to table's translator if nil passed
func TestPrintFallbackTranslator(t *testing.T) {
	ctx := testContext()

	builder := NewBuilder[testDeviceRow](ctx, testTranslator)
	builder.IdField("id", "device.id", func(r testDeviceRow) int64 { return r.ID })
	tbl := builder.Build()
	tbl.SetData([]testDeviceRow{{ID: 1}})

	// Pass nil translator - should fall back to table's translator
	output := tbl.Print(nil)

	// Should not panic, should use table's translator
	if output["type"] != "table" {
		t.Errorf("Expected type 'table', got %v", output["type"])
	}
}
