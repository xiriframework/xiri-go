package table

import (
	"strings"
	"testing"
)

// TestSanitizeExportValue verifies that formula injection prefixes are neutralized
// with a single quote in CSV/Excel exports (H2 unit test).
func TestSanitizeExportValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"equals prefix", "=cmd|'/C calc'!A0", "'=cmd|'/C calc'!A0"},
		{"plus prefix", "+1234", "'+1234"},
		{"minus prefix", "-formula", "'-formula"},
		{"at prefix", "@evil", "'@evil"},
		{"tab prefix", "\tevil", "'\tevil"},
		{"cr prefix", "\revil", "'\revil"},
		{"lf prefix", "\nevil", "'\nevil"},
		{"safe string", "hello", "hello"},
		{"empty string", "", ""},
		{"numeric string", "12345", "12345"},
		{"equals in middle", "x=y", "x=y"},
		{"space then equals", " =SUM", " =SUM"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeExportValue(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizeExportValue(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestGenerateCSV_FormulaInjection verifies that CSV output sanitizes dangerous
// values with formula injection prefixes (H2 integration test).
func TestGenerateCSV_FormulaInjection(t *testing.T) {
	td := NewTableDataResponse([]map[string]any{
		{"name": "=cmd|'/C calc'!A0", "value": "safe-value"},
		{"name": "+1234", "value": "normal"},
		{"name": "-formula", "value": "ok"},
		{"name": "@evil", "value": "fine"},
		{"name": "safe", "value": "also-safe"},
	}, OutputCSV)

	td.WithFields([]map[string]any{
		{"id": "name", "name": "Name"},
		{"id": "value", "name": "Value"},
	})

	csvOutput := td.generateCSV(nil)

	// Dangerous values must be prefixed with single quote
	if !strings.Contains(csvOutput, "'=cmd|") {
		t.Error("expected '=' prefix to be sanitized in CSV output")
	}
	if !strings.Contains(csvOutput, "'+1234") {
		t.Error("expected '+' prefix to be sanitized in CSV output")
	}
	if !strings.Contains(csvOutput, "'-formula") {
		t.Error("expected '-' prefix to be sanitized in CSV output")
	}
	if !strings.Contains(csvOutput, "'@evil") {
		t.Error("expected '@' prefix to be sanitized in CSV output")
	}

	// Safe values must NOT be prefixed
	if strings.Contains(csvOutput, "'safe-value") {
		t.Error("safe value should not be prefixed")
	}
	if strings.Contains(csvOutput, "'normal") {
		t.Error("safe value should not be prefixed")
	}
	if strings.Contains(csvOutput, "'also-safe") {
		t.Error("safe value should not be prefixed")
	}
}

// TestGenerateExcel_FormulaInjection verifies that Excel output handles
// formula injection values without panicking and produces valid output (H2 integration test).
func TestGenerateExcel_FormulaInjection(t *testing.T) {
	td := NewTableDataResponse([]map[string]any{
		{"name": "=cmd|'/C calc'!A0", "value": "safe"},
		{"name": "+1234", "value": "normal"},
		{"name": "-formula", "value": "ok"},
		{"name": "@evil", "value": "fine"},
	}, OutputExcel)

	td.WithFields([]map[string]any{
		{"id": "name", "name": "Name"},
		{"id": "value", "name": "Value"},
	})

	excelBytes, err := td.generateExcel(nil)
	if err != nil {
		t.Fatalf("generateExcel returned unexpected error: %v", err)
	}
	if len(excelBytes) == 0 {
		t.Fatal("generateExcel returned empty bytes")
	}
}
