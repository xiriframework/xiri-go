package table

import (
	"testing"
)

// findFieldJSON returns the serialized field with the given id from a Print() output.
func findFieldJSON(t *testing.T, tbl *Table[testDeviceRow], id string) map[string]any {
	t.Helper()
	output := tbl.Print(testContext())
	data := output["data"].(map[string]any)
	fields := data["fields"].([]map[string]any)
	for _, f := range fields {
		if f["id"] == id {
			return f
		}
	}
	t.Fatalf("field %q not found in output", id)
	return nil
}

// TestInlineEditServerSearch verifies WithEditableSearchOptionsUrl serializes the
// search url and enables the search flag.
func TestInlineEditServerSearch(t *testing.T) {
	builder := NewBuilder[testDeviceRow]()
	builder.IdField("id", "device.id", func(r testDeviceRow) int64 { return r.ID })
	builder.TextField("name", "device.name", func(r testDeviceRow) string { return r.Name }).
		WithEditableSearchOptionsUrl("/api/devices/search")
	tbl := builder.Build()

	field := findFieldJSON(t, tbl, "name")

	if field["editable"] != true {
		t.Errorf("expected editable=true, got %v", field["editable"])
	}
	if field["editableSearchUrl"] != "/api/devices/search" {
		t.Errorf("expected editableSearchUrl, got %v", field["editableSearchUrl"])
	}
	if field["editableOptionsSearch"] != true {
		t.Errorf("expected editableOptionsSearch=true, got %v", field["editableOptionsSearch"])
	}
}

// TestInlineEditClientSearch verifies WithEditableOptionsSearch enables the search
// box without a server url, keeping static options intact.
func TestInlineEditClientSearch(t *testing.T) {
	builder := NewBuilder[testDeviceRow]()
	builder.IdField("id", "device.id", func(r testDeviceRow) int64 { return r.ID })
	builder.TextField("name", "device.name", func(r testDeviceRow) string { return r.Name }).
		WithEditableOptions(map[string]string{"a": "Alpha", "b": "Beta"}).
		WithEditableOptionsSearch(true)
	tbl := builder.Build()

	field := findFieldJSON(t, tbl, "name")

	if field["editable"] != true {
		t.Errorf("expected editable=true, got %v", field["editable"])
	}
	if field["editableOptionsSearch"] != true {
		t.Errorf("expected editableOptionsSearch=true, got %v", field["editableOptionsSearch"])
	}
	if _, ok := field["editableSearchUrl"]; ok {
		t.Errorf("expected no editableSearchUrl for client-side search, got %v", field["editableSearchUrl"])
	}
	if opts, ok := field["editableOptions"].([]map[string]any); !ok || len(opts) != 2 {
		t.Errorf("expected 2 static editableOptions, got %v", field["editableOptions"])
	}
}

// TestInlineEditNoSearchBackwardCompat verifies plain editable fields do not emit
// the new search keys (backward compatibility).
func TestInlineEditNoSearchBackwardCompat(t *testing.T) {
	builder := NewBuilder[testDeviceRow]()
	builder.IdField("id", "device.id", func(r testDeviceRow) int64 { return r.ID })
	builder.TextField("name", "device.name", func(r testDeviceRow) string { return r.Name }).
		WithEditable(true)
	tbl := builder.Build()

	field := findFieldJSON(t, tbl, "name")

	if _, ok := field["editableOptionsSearch"]; ok {
		t.Errorf("expected no editableOptionsSearch key, got %v", field["editableOptionsSearch"])
	}
	if _, ok := field["editableSearchUrl"]; ok {
		t.Errorf("expected no editableSearchUrl key, got %v", field["editableSearchUrl"])
	}
}
