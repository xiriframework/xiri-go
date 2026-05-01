package table

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/types/distance"
	"github.com/xiriframework/xiri-go/types/language"
	"github.com/xiriframework/xiri-go/types/locale"
	"github.com/xiriframework/xiri-go/types/timezone"
)

type fleetRow struct {
	Vehicle  string
	Warnings []Chip
}

func chipsTestCtx() *core.UiContext {
	return &core.UiContext{
		Timezone:  timezone.EuropeVienna,
		Lang:      language.Deutsch,
		Locale:    locale.De,
		Distance:  distance.Kilometer,
		Translate: func(k string) string { return k },
	}
}

func TestChipsField_WebOutput(t *testing.T) {
	ctx := chipsTestCtx()
	rows := []fleetRow{
		{
			Vehicle: "GRI-422",
			Warnings: []Chip{
				{Label: "Charging issue", Color: core.ColorRed},
				{Label: "45%", Color: core.ColorRed},
			},
		},
		{
			Vehicle:  "SZH-689",
			Warnings: []Chip{{Label: "Brake system", Color: core.ColorWarning}},
		},
	}

	b := NewBuilder[fleetRow]()
	b.TextField("vehicle", "fleet.vehicle", func(r fleetRow) string { return r.Vehicle })
	b.ChipsField("warnings", "fleet.warnings", func(r fleetRow) []Chip { return r.Warnings })
	tbl := b.Build()
	tbl.SetData(rows)

	output := tbl.Print(ctx)
	data := output["data"].(map[string]any)

	fields := data["fields"].([]map[string]any)
	if len(fields) != 2 {
		t.Fatalf("Expected 2 fields, got %d", len(fields))
	}
	if fields[1]["format"] != "chips" {
		t.Errorf("Expected fields[1].format='chips', got %v", fields[1]["format"])
	}

	rowData := data["data"].([]map[string]any)
	if len(rowData) != 2 {
		t.Fatalf("Expected 2 rows, got %d", len(rowData))
	}

	chips0, ok := rowData[0]["warnings"].([]map[string]any)
	if !ok {
		t.Fatalf("Expected row 0 warnings to be []map[string]any, got %T (%v)", rowData[0]["warnings"], rowData[0]["warnings"])
	}
	if len(chips0) != 2 {
		t.Fatalf("Expected row 0 to have 2 chips, got %d", len(chips0))
	}
	if chips0[0]["label"] != "Charging issue" {
		t.Errorf("chips0[0].label = %v, want 'Charging issue'", chips0[0]["label"])
	}
	if chips0[0]["color"] != "red" {
		t.Errorf("chips0[0].color = %v, want 'red'", chips0[0]["color"])
	}
	if chips0[1]["label"] != "45%" || chips0[1]["color"] != "red" {
		t.Errorf("chips0[1] = %v, want {label:'45%%', color:'red'}", chips0[1])
	}

	chips1, ok := rowData[1]["warnings"].([]map[string]any)
	if !ok {
		t.Fatalf("Expected row 1 warnings to be []map[string]any, got %T", rowData[1]["warnings"])
	}
	if len(chips1) != 1 || chips1[0]["color"] != "warn" {
		t.Errorf("chips1 = %v, want one chip with color 'warn'", chips1)
	}
}

func TestChipsField_EmptyAndNil(t *testing.T) {
	ctx := chipsTestCtx()
	rows := []fleetRow{
		{Vehicle: "A", Warnings: nil},
		{Vehicle: "B", Warnings: []Chip{}},
	}

	b := NewBuilder[fleetRow]()
	b.TextField("vehicle", "fleet.vehicle", func(r fleetRow) string { return r.Vehicle })
	b.ChipsField("warnings", "fleet.warnings", func(r fleetRow) []Chip { return r.Warnings })
	tbl := b.Build()
	tbl.SetData(rows)

	output := tbl.Print(ctx)
	rowData := output["data"].(map[string]any)["data"].([]map[string]any)

	for i, row := range rowData {
		chips, ok := row["warnings"].([]map[string]any)
		if !ok {
			t.Fatalf("Row %d: expected []map[string]any, got %T", i, row["warnings"])
		}
		if len(chips) != 0 {
			t.Errorf("Row %d: expected 0 chips, got %d", i, len(chips))
		}
	}
}

func TestChipsField_CSVOutput(t *testing.T) {
	ctx := chipsTestCtx()
	rows := []fleetRow{
		{
			Vehicle: "GRI-422",
			Warnings: []Chip{
				{Label: "Charging issue", Color: core.ColorRed},
				{Label: "45%", Color: core.ColorRed},
			},
		},
	}

	b := NewBuilder[fleetRow]()
	b.TextField("vehicle", "fleet.vehicle", func(r fleetRow) string { return r.Vehicle })
	b.ChipsField("warnings", "fleet.warnings", func(r fleetRow) []Chip { return r.Warnings })
	tbl := b.Build()
	tbl.SetData(rows)

	csvData := tbl.GetData(ctx, OutputCSV)
	if len(csvData) != 1 {
		t.Fatalf("Expected 1 row, got %d", len(csvData))
	}
	got := csvData[0]["warnings"]
	if got != "Charging issue, 45%" {
		t.Errorf("CSV warnings = %q, want %q", got, "Charging issue, 45%")
	}
}
