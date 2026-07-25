package table

import (
	"strings"
	"testing"
)

// #9: N-field expansion must not mutate the input rows, and must be idempotent.
func TestExpandNFieldColumnsDoesNotMutateInput(t *testing.T) {
	data := []map[string]any{
		{"dist": []string{"1", "2", "3"}},
	}
	fields := []map[string]any{
		{"id": "dist", "name": "Distance"},
	}
	td := NewTableDataResponse(data, OutputCSV).withFieldsForCSV(fields)

	first := td.generateCSV(nil)

	// Expansion must actually happen (guards against a no-op passing the test).
	if !strings.Contains(first, "1;2;3") {
		t.Errorf("N-field not expanded into columns; CSV = %q", first)
	}
	// The original input row must be untouched (still a []string).
	if _, ok := data[0]["dist"].([]string); !ok {
		t.Errorf("input row mutated: dist = %v (%T), want []string", data[0]["dist"], data[0]["dist"])
	}
	// The original fields slice must be untouched (still a single field def).
	if len(fields) != 1 {
		t.Errorf("input fields mutated: len = %d, want 1", len(fields))
	}

	// Re-rendering must yield identical output.
	second := td.generateCSV(nil)
	if first != second {
		t.Errorf("CSV not idempotent:\nfirst  = %q\nsecond = %q", first, second)
	}
}

// #10: CSV header cells must be sanitized against formula injection like data cells.
func TestCSVHeaderSanitized(t *testing.T) {
	data := []map[string]any{{"col": "x"}}
	fields := []map[string]any{{"id": "col", "name": "=SUM(A1)"}}
	td := NewTableDataResponse(data, OutputCSV).withFieldsForCSV(fields)

	csv := td.generateCSV(nil)
	if !strings.Contains(csv, "'=SUM(A1)") {
		t.Errorf("header not sanitized against formula injection: %q", csv)
	}
}
