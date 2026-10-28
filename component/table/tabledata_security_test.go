package table

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/xuri/excelize/v2"
)

// TestSanitizeExportValue verifies that formula injection prefixes are neutralized
// with a single quote in CSV exports (H2 unit test). Plain numbers are no formula and stay unchanged.
func TestSanitizeExportValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"equals prefix", "=cmd|'/C calc'!A0", "'=cmd|'/C calc'!A0"},
		{"plus prefix", "+A1", "'+A1"},
		{"plus number", "+1234", "+1234"},
		{"negative decimal comma", "-12,50", "-12,50"},
		{"negative decimal dot", "-6.50", "-6.50"},
		{"negative with thousands", "-1.234,56", "-1.234,56"},
		{"arithmetic", "-1+2", "'-1+2"},
		{"exponent", "-1e5", "'-1e5"},
		{"phone number", "+43 664 1234", "'+43 664 1234"},
		{"sign only", "-", "'-"},
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
		{"name": "+A1", "value": "normal"},
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
	if !strings.Contains(csvOutput, "'+A1") {
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

// TestGenerateExcel_FormulaInjection verifies that Excel output writes formula-like text as plain
// string cells: never a formula, and without a visible "'" prefix (H2 integration test).
func TestGenerateExcel_FormulaInjection(t *testing.T) {
	values := []string{"=cmd|'/C calc'!A0", "+A1", "-formula", "@evil"}
	rows := make([]map[string]any, len(values))
	for i, v := range values {
		rows[i] = map[string]any{"name": v}
	}
	td := NewTableDataResponse(rows, OutputExcel)
	td.WithFields([]map[string]any{{"id": "name", "name": "=Name"}})

	excelBytes, err := td.generateExcel(nil)
	if err != nil {
		t.Fatalf("generateExcel returned unexpected error: %v", err)
	}
	f, err := excelize.OpenReader(bytes.NewReader(excelBytes))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	for i, want := range append([]string{"=Name"}, values...) {
		cell := fmt.Sprintf("A%d", i+1)
		typ, _ := f.GetCellType("Sheet1", cell)
		got, _ := f.GetCellValue("Sheet1", cell)
		formula, _ := f.GetCellFormula("Sheet1", cell)
		if typ != excelize.CellTypeSharedString || got != want || formula != "" {
			t.Errorf("%s: type %v value %q formula %q, want string cell %q", cell, typ, got, formula, want)
		}
	}
}
