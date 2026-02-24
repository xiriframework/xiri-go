package table_test

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/xiriframework/xiri-go/component/table"
)

// ============================================================================
// IconSet Unit Tests
// ============================================================================

func TestIconSet_Add(t *testing.T) {
	icons := table.NewIconSet()

	ref := icons.Add("online", "check_circle", table.FieldColorAccent, "Online")
	if ref == nil {
		t.Fatal("Add() returned nil")
	}
	if icons.Len() != 1 {
		t.Fatalf("expected Len() = 1, got %d", icons.Len())
	}

	ref2 := icons.Add("offline", "cancel", table.FieldColorWarning, "Offline")
	if ref2 == nil {
		t.Fatal("Add() returned nil for second icon")
	}
	if icons.Len() != 2 {
		t.Fatalf("expected Len() = 2, got %d", icons.Len())
	}
}

func TestIconSet_AddWithOptions(t *testing.T) {
	icons := table.NewIconSet()

	opts := map[string]any{"pulse": true}
	ref := icons.AddWithOptions("warning", "warning", table.FieldColorWarning, "Warning", opts)
	if ref == nil {
		t.Fatal("AddWithOptions() returned nil")
	}
	if icons.Len() != 1 {
		t.Fatalf("expected Len() = 1, got %d", icons.Len())
	}
}

func TestIconSet_Resolve_Known(t *testing.T) {
	icons := table.NewIconSet()
	icons.Add("online", "check_circle", table.FieldColorAccent, "Online")
	icons.Add("offline", "cancel", table.FieldColorWarning, "Offline")

	ref := icons.Resolve("online")
	if ref == nil {
		t.Fatal("Resolve('online') returned nil for registered value")
	}

	ref2 := icons.Resolve("offline")
	if ref2 == nil {
		t.Fatal("Resolve('offline') returned nil for registered value")
	}
}

func TestIconSet_Resolve_Unknown(t *testing.T) {
	icons := table.NewIconSet()
	icons.Add("online", "check_circle", table.FieldColorAccent, "Online")

	ref := icons.Resolve("unknown_value")
	if ref != nil {
		t.Fatal("Resolve('unknown_value') should return nil for unregistered value")
	}

	ref2 := icons.Resolve("")
	if ref2 != nil {
		t.Fatal("Resolve('') should return nil for empty string")
	}
}

func TestIconSet_Resolve_WarnsOnUnknownValue(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelWarn})
	logger := slog.New(handler)
	old := slog.Default()
	slog.SetDefault(logger)
	t.Cleanup(func() { slog.SetDefault(old) })

	icons := table.NewIconSet()
	icons.Add("online", "check_circle", table.FieldColorAccent, "Online")

	ref := icons.Resolve("unknown_value")
	if ref != nil {
		t.Fatal("Resolve('unknown_value') should return nil")
	}

	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("unknown value")) {
		t.Errorf("expected slog warning containing 'unknown value', got: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("unknown_value")) {
		t.Errorf("expected slog warning containing the actual value 'unknown_value', got: %s", output)
	}
}

func TestIconSet_Len_Empty(t *testing.T) {
	icons := table.NewIconSet()
	if icons.Len() != 0 {
		t.Fatalf("expected Len() = 0, got %d", icons.Len())
	}
}

// ============================================================================
// IconFieldFromSet Integration Tests
// ============================================================================

type iconTestRow struct {
	IsOnline bool
	Status   string
}

func TestIconFieldFromSet_DirectRefs(t *testing.T) {
	ctx := exampleContext()

	statusIcons := table.NewIconSet()
	iconOnline := statusIcons.Add("online", "check_circle", table.FieldColorAccent, "Online")
	iconOffline := statusIcons.Add("offline", "cancel", table.FieldColorWarning, "Offline")

	builder := table.NewBuilder[iconTestRow](ctx, exampleTranslator)
	builder.IconFieldFromSet("status", "device.status",
		func(r iconTestRow) *table.IconRef {
			if r.IsOnline {
				return iconOnline
			}
			return iconOffline
		},
		statusIcons,
	)

	tbl := builder.Build()
	tbl.SetData([]iconTestRow{
		{IsOnline: true},
		{IsOnline: false},
	})

	data := tbl.GetData(table.OutputWeb)

	if len(data) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(data))
	}

	// Row 0: online
	if data[0]["status"] != "online" {
		t.Errorf("expected 'online', got %v", data[0]["status"])
	}

	// Row 1: offline
	if data[1]["status"] != "offline" {
		t.Errorf("expected 'offline', got %v", data[1]["status"])
	}
}

func TestIconFieldFromSet_Resolve(t *testing.T) {
	ctx := exampleContext()

	statusIcons := table.NewIconSet()
	statusIcons.Add("online", "check_circle", table.FieldColorAccent, "Online")
	statusIcons.Add("offline", "cancel", table.FieldColorWarning, "Offline")

	builder := table.NewBuilder[iconTestRow](ctx, exampleTranslator)
	builder.IconFieldFromSet("status", "device.status",
		func(r iconTestRow) *table.IconRef {
			return statusIcons.Resolve(r.Status)
		},
		statusIcons,
	)

	tbl := builder.Build()
	tbl.SetData([]iconTestRow{
		{Status: "online"},
		{Status: "offline"},
		{Status: "unknown"},
	})

	data := tbl.GetData(table.OutputWeb)

	if data[0]["status"] != "online" {
		t.Errorf("expected 'online', got %v", data[0]["status"])
	}
	if data[1]["status"] != "offline" {
		t.Errorf("expected 'offline', got %v", data[1]["status"])
	}
	// Unknown value resolves to nil → empty string
	if data[2]["status"] != "" {
		t.Errorf("expected empty string for unknown value, got %v", data[2]["status"])
	}
}

func TestIconFieldFromSet_NilRef(t *testing.T) {
	ctx := exampleContext()

	statusIcons := table.NewIconSet()
	statusIcons.Add("online", "check_circle", table.FieldColorAccent, "Online")

	builder := table.NewBuilder[iconTestRow](ctx, exampleTranslator)
	builder.IconFieldFromSet("status", "device.status",
		func(r iconTestRow) *table.IconRef {
			return nil // Always nil
		},
		statusIcons,
	)

	tbl := builder.Build()
	tbl.SetData([]iconTestRow{{IsOnline: true}})

	data := tbl.GetData(table.OutputWeb)

	if data[0]["status"] != "" {
		t.Errorf("expected empty string for nil IconRef, got %v", data[0]["status"])
	}
}

func TestIconFieldFromSet_IconDefinitionsCopied(t *testing.T) {
	ctx := exampleContext()

	statusIcons := table.NewIconSet()
	iconOnline := statusIcons.Add("online", "check_circle", table.FieldColorAccent, "Online")

	builder := table.NewBuilder[iconTestRow](ctx, exampleTranslator)
	builder.IconFieldFromSet("status", "device.status",
		func(r iconTestRow) *table.IconRef {
			return iconOnline
		},
		statusIcons,
	)

	tbl := builder.Build()
	tbl.SetData([]iconTestRow{{IsOnline: true}})

	// Verify icon definitions are present in field JSON output
	fields := tbl.Print(exampleTranslator)["data"].(map[string]any)["fields"].([]map[string]any)

	if len(fields) != 1 {
		t.Fatalf("expected 1 field, got %d", len(fields))
	}

	icons, ok := fields[0]["icons"].(map[string]any)
	if !ok {
		t.Fatal("expected icons map in field definition")
	}

	if _, exists := icons["online"]; !exists {
		t.Error("expected 'online' icon definition in field")
	}
}

func TestIconFieldFromSet_WithOptions(t *testing.T) {
	ctx := exampleContext()

	statusIcons := table.NewIconSet()
	iconOnline := statusIcons.AddWithOptions("online", "check_circle", table.FieldColorAccent, "Online", map[string]any{"pulse": true})

	builder := table.NewBuilder[iconTestRow](ctx, exampleTranslator)
	builder.IconFieldFromSet("status", "device.status",
		func(r iconTestRow) *table.IconRef {
			return iconOnline
		},
		statusIcons,
	)

	tbl := builder.Build()
	tbl.SetData([]iconTestRow{{IsOnline: true}})

	// Verify custom options are present in the exported icon definition
	fields := tbl.Print(exampleTranslator)["data"].(map[string]any)["fields"].([]map[string]any)
	icons := fields[0]["icons"].(map[string]any)
	onlineIcon := icons["online"].(map[string]any)

	if pulse, ok := onlineIcon["pulse"]; !ok || pulse != true {
		t.Error("expected 'pulse: true' in icon options")
	}
}

func TestIconFieldFromSet_BuilderChaining(t *testing.T) {
	ctx := exampleContext()

	statusIcons := table.NewIconSet()
	iconOnline := statusIcons.Add("online", "check_circle", table.FieldColorAccent, "Online")

	builder := table.NewBuilder[iconTestRow](ctx, exampleTranslator)

	// Verify FieldBuilder methods can be chained
	builder.IconFieldFromSet("status", "device.status",
		func(r iconTestRow) *table.IconRef {
			return iconOnline
		},
		statusIcons,
	).AlignCenter().WithWidth("100px").WithHint("Status")

	tbl := builder.Build()
	tbl.SetData([]iconTestRow{{IsOnline: true}})

	fields := tbl.Print(exampleTranslator)["data"].(map[string]any)["fields"].([]map[string]any)
	field := fields[0]

	if field["align"] != "center" {
		t.Errorf("expected align 'center', got %v", field["align"])
	}
	if field["width"] != "100px" {
		t.Errorf("expected width '100px', got %v", field["width"])
	}
	if field["hint"] != "Status" {
		t.Errorf("expected hint 'Status', got %v", field["hint"])
	}
}
