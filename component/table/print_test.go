package table

import (
	"testing"

	xurl "github.com/xiriframework/xiri-go/component/url"
	"github.com/xiriframework/xiri-go/types/distance"
	"github.com/xiriframework/xiri-go/types/language"
	"github.com/xiriframework/xiri-go/types/locale"
	"github.com/xiriframework/xiri-go/types/timezone"
	"github.com/xiriframework/xiri-go/component/core"
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
func testContext() *core.UiContext {
	return &core.UiContext{
		Timezone:  timezone.EuropeVienna,
		Lang:      language.Deutsch,
		Locale:    locale.De,
		Distance:  distance.Kilometer,
		Translate: testTranslator,
	}
}

// TestPrintAJAXMode verifies Print() returns correct structure for AJAX tables (url != nil)
func TestPrintAJAXMode(t *testing.T) {
	ctx := testContext()

	url := xurl.NewUrl("/Portal/Device/TableData")

	builder := NewBuilder[testDeviceRow]()
	builder.IdField("id", "device.id", func(r testDeviceRow) int64 { return r.ID })
	builder.TextField("name", "device.name", func(r testDeviceRow) string { return r.Name })
	tbl := builder.Build()
	tbl.SetURL(url)

	// Call Print() - should produce AJAX mode JSON
	output := tbl.Print(ctx)

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

	builder := NewBuilder[testDeviceRow]()
	builder.IdField("id", "device.id", func(r testDeviceRow) int64 { return r.ID })
	builder.TextField("name", "device.name", func(r testDeviceRow) string { return r.Name })
	builder.BoolField("active", "device.active", func(r testDeviceRow) bool { return r.Active }).
		WithBoolText("Yes", "No")
	tbl := builder.Build()
	tbl.SetData(rows)

	// Call Print() - should produce static mode JSON
	output := tbl.Print(ctx)

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
	ctx := &core.UiContext{
		Timezone:  timezone.EuropeVienna,
		Lang:      language.Deutsch,
		Locale:    locale.De,
		Distance:  distance.Kilometer,
		Translate: func(key string) string {
			if key == "device.id" {
				return "Translated ID"
			}
			return key
		},
	}

	builder := NewBuilder[testDeviceRow]()
	builder.IdField("id", "device.id", func(r testDeviceRow) int64 { return r.ID })
	tbl := builder.Build()
	tbl.SetData([]testDeviceRow{{ID: 1}})

	output := tbl.Print(ctx)

	data := output["data"].(map[string]any)
	fields := data["fields"].([]map[string]any)

	// Verify translation was applied
	if fields[0]["name"] != "Translated ID" {
		t.Errorf("Expected field name to be 'Translated ID', got %v", fields[0]["name"])
	}
}

// Tree test row struct (adds ParentID to the device row)
type testTreeRow struct {
	ID       int64
	ParentID int64
	Name     string
}

// TestPrintTreeMode verifies that .Tree(...) serializes the nested tree settings object.
func TestPrintTreeMode(t *testing.T) {
	ctx := testContext()

	builder := NewBuilder[testTreeRow]()
	builder.IdField("id", "device.id", func(r testTreeRow) int64 { return r.ID })
	builder.IdField("parentId", "", func(r testTreeRow) int64 { return r.ParentID })
	builder.TextField("name", "device.name", func(r testTreeRow) string { return r.Name })
	tbl := builder.Tree(TreeConfig{
		IdField:         "id",
		ParentIdField:   "parentId",
		TreeColumn:      "name",
		PersistStateKey: "portal-groups",
		AddSubURL:       xurl.NewUrl("/Portal/Group/Add?parent={id}"),
	}).Build()
	tbl.SetData([]testTreeRow{
		{ID: 1, ParentID: 0, Name: "Wien"},
		{ID: 2, ParentID: 1, Name: "Favoriten"},
	})

	output := tbl.Print(ctx)
	data := output["data"].(map[string]any)
	options := data["options"].(map[string]any)

	tree, ok := options["tree"].(map[string]any)
	if !ok {
		t.Fatalf("Expected options[tree] to be map[string]any, got %T", options["tree"])
	}

	if tree["idField"] != "id" {
		t.Errorf("Expected idField='id', got %v", tree["idField"])
	}
	if tree["parentIdField"] != "parentId" {
		t.Errorf("Expected parentIdField='parentId', got %v", tree["parentIdField"])
	}
	if tree["treeColumn"] != "name" {
		t.Errorf("Expected treeColumn='name', got %v", tree["treeColumn"])
	}
	if tree["persistStateKey"] != "portal-groups" {
		t.Errorf("Expected persistStateKey='portal-groups', got %v", tree["persistStateKey"])
	}
	if tree["addSubUrl"] != "/Portal/Group/Add?parent={id}" {
		t.Errorf("Expected addSubUrl placeholder URL, got %v", tree["addSubUrl"])
	}
	// HideCounts defaults to false → showCounts should be true.
	if tree["showCounts"] != true {
		t.Errorf("Expected showCounts=true (default), got %v", tree["showCounts"])
	}
	if tree["collapseAllByDefault"] != false {
		t.Errorf("Expected collapseAllByDefault=false (default → expanded), got %v", tree["collapseAllByDefault"])
	}

	// parentId must be present in the row data (id-format field, not hidden).
	rowData := data["data"].([]map[string]any)
	if rowData[1]["parentId"] == nil {
		t.Error("Expected parentId to be present in row data")
	}
}

// TestPrintNoTreeByDefault verifies backwards-compatibility: tables without .Tree(...)
// must NOT emit a tree key in options.
func TestPrintNoTreeByDefault(t *testing.T) {
	ctx := testContext()

	builder := NewBuilder[testDeviceRow]()
	builder.IdField("id", "device.id", func(r testDeviceRow) int64 { return r.ID })
	builder.TextField("name", "device.name", func(r testDeviceRow) string { return r.Name })
	tbl := builder.Build()
	tbl.SetData([]testDeviceRow{{ID: 1, Name: "Device 1"}})

	output := tbl.Print(ctx)
	data := output["data"].(map[string]any)
	options := data["options"].(map[string]any)

	if _, exists := options["tree"]; exists {
		t.Error("Expected no 'tree' key in options for a flat table")
	}
}

// TestPrintTreeShowCountsDisabled verifies HideCounts=true yields showCounts=false.
func TestPrintTreeShowCountsDisabled(t *testing.T) {
	ctx := testContext()

	builder := NewBuilder[testTreeRow]()
	builder.IdField("id", "device.id", func(r testTreeRow) int64 { return r.ID })
	builder.IdField("parentId", "", func(r testTreeRow) int64 { return r.ParentID })
	builder.TextField("name", "device.name", func(r testTreeRow) string { return r.Name })
	tbl := builder.Tree(TreeConfig{
		IdField:            "id",
		ParentIdField:      "parentId",
		CollapseAllDefault: true,
		HideCounts:         true,
	}).Build()
	tbl.SetData([]testTreeRow{{ID: 1, Name: "Wien"}})

	output := tbl.Print(ctx)
	tree := output["data"].(map[string]any)["options"].(map[string]any)["tree"].(map[string]any)

	if tree["showCounts"] != false {
		t.Errorf("Expected showCounts=false when HideCounts set, got %v", tree["showCounts"])
	}
	if tree["collapseAllByDefault"] != true {
		t.Errorf("Expected collapseAllByDefault=true, got %v", tree["collapseAllByDefault"])
	}
	// Optional keys must be omitted when unset.
	if _, exists := tree["treeColumn"]; exists {
		t.Error("Expected treeColumn omitted when empty")
	}
	if _, exists := tree["addSubUrl"]; exists {
		t.Error("Expected addSubUrl omitted when nil")
	}
}

// TestPrintNilCtx verifies that Print(nil) does not panic
func TestPrintNilCtx(t *testing.T) {
	builder := NewBuilder[testDeviceRow]()
	builder.IdField("id", "device.id", func(r testDeviceRow) int64 { return r.ID })
	tbl := builder.Build()
	tbl.SetData([]testDeviceRow{{ID: 1}})

	// Pass nil ctx - should not panic, fields returned with untranslated keys
	output := tbl.Print(nil)

	if output["type"] != "table" {
		t.Errorf("Expected type 'table', got %v", output["type"])
	}
}
