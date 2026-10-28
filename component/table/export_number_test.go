package table

import (
	"bytes"
	"math"
	"strings"
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/types/distance"
	"github.com/xiriframework/xiri-go/types/locale"
	"github.com/xuri/excelize/v2"
)

type exportRow struct {
	Amount float64
	Count  int64
	Name   string
}

func exportTable(output OutputType, rows []exportRow) *TableDataResponse {
	b := NewBuilder[exportRow]()
	b.FloatField("amount", "Betrag", func(r exportRow) float64 { return r.Amount }).WithDecimals(2)
	b.Int64Field("count", "Anzahl", func(r exportRow) int64 { return r.Count })
	b.TextField("name", "Name", func(r exportRow) string { return r.Name })
	tbl := b.Build()
	tbl.SetData(rows)
	tbl.SetOutputType(output)
	return tbl.ToTableDataResponse(&core.UiContext{Locale: locale.De})
}

func csvLines(t *testing.T, td *TableDataResponse, loc locale.Locale) []string {
	t.Helper()
	out := td.generateCSV(&core.UiContext{Locale: loc})
	if !strings.HasPrefix(out, "\xEF\xBB\xBF") {
		t.Fatalf("CSV must start with UTF-8 BOM, got %q", out)
	}
	return strings.Split(strings.TrimSuffix(strings.TrimPrefix(out, "\xEF\xBB\xBF"), "\n"), "\n")
}

func TestExportCSV_NumbersFollowLocale(t *testing.T) {
	rows := []exportRow{
		{Amount: 1234.567, Count: 1234567, Name: "=SUM(A1)"},
		{Amount: -12.5, Count: -3, Name: "ok"},
		{Amount: 1e7, Count: 0, Name: "Übung"},
	}

	de := csvLines(t, exportTable(OutputCSV, rows), locale.De)
	want := []string{
		"Betrag;Anzahl;Name",
		"1234,57;1234567;'=SUM(A1)",
		"-12,50;-3;ok",
		"10000000,00;0;Übung",
	}
	if strings.Join(de, "\n") != strings.Join(want, "\n") {
		t.Errorf("CSV de:\n%s\nwant:\n%s", strings.Join(de, "\n"), strings.Join(want, "\n"))
	}

	en := csvLines(t, exportTable(OutputCSV, rows), locale.EnUS)
	if en[1] != "1234.57;1234567;'=SUM(A1)" || en[2] != "-12.50;-3;ok" {
		t.Errorf("CSV en rows = %q", en[1:])
	}
}

func TestExportCSV_NativeAndNilNumbers(t *testing.T) {
	td := NewTableDataResponse([]map[string]any{
		{"f": float64(1e6), "i": -5, "n": nil, "e": ExportNumber{Value: -0.5, Decimals: 1}, "s": int8(-7), "u": uint16(9)},
	}, OutputCSV)
	td.withFieldsForCSV([]map[string]any{
		{"id": "f", "name": "F"}, {"id": "i", "name": "I"}, {"id": "n", "name": "N"}, {"id": "e", "name": "E"},
		{"id": "s", "name": "S"}, {"id": "u", "name": "U"},
	})
	lines := csvLines(t, td, locale.De)
	if lines[1] != "1000000;-5;;-0,5;-7;9" {
		t.Errorf("row = %q, want %q", lines[1], "1000000;-5;;-0,5;-7;9")
	}
}

func TestExportExcel_NumbersAreNumericCells(t *testing.T) {
	rows := []exportRow{{Amount: -12.5, Count: 1234567, Name: "=SUM(A1)"}}
	td := exportTable(OutputExcel, rows)
	data, err := td.generateExcel(&core.UiContext{Locale: locale.De})
	if err != nil {
		t.Fatal(err)
	}
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	for _, c := range []struct{ cell, raw, numFmt string }{
		{"A2", "-12.5", "#,##0.00"},
		{"B2", "1234567", "#,##0"},
	} {
		typ, _ := f.GetCellType("Sheet1", c.cell)
		if typ == excelize.CellTypeSharedString || typ == excelize.CellTypeInlineString {
			t.Errorf("%s: cell type %v is text", c.cell, typ)
		}
		raw, _ := f.GetCellValue("Sheet1", c.cell, excelize.Options{RawCellValue: true})
		if raw != c.raw {
			t.Errorf("%s: raw value %q, want %q", c.cell, raw, c.raw)
		}
		styleID, _ := f.GetCellStyle("Sheet1", c.cell)
		style, err := f.GetStyle(styleID)
		if err != nil || style.CustomNumFmt == nil || *style.CustomNumFmt != c.numFmt {
			t.Errorf("%s: number format %v, want %q", c.cell, style, c.numFmt)
		}
	}

	// xlsx string cells are never formulas: the text goes out unchanged, without a "'" prefix.
	if v, _ := f.GetCellValue("Sheet1", "C2"); v != "=SUM(A1)" {
		t.Errorf("text must stay unchanged, got %q", v)
	}
	if formula, _ := f.GetCellFormula("Sheet1", "C2"); formula != "" {
		t.Errorf("text must not become a formula, got %q", formula)
	}
}

func TestExportExcel_CustomFormatterNativeFloat(t *testing.T) {
	b := NewBuilder[exportRow]()
	b.TextField("amount", "Betrag", func(r exportRow) string { return "€ 1.234,56" }).
		WithExcelFormatter(FormatterFunc(func(value any, row Row, output OutputType, ctx *core.UiContext) any {
			return 1234.56
		}))
	tbl := b.Build()
	tbl.SetData([]exportRow{{}})
	tbl.SetOutputType(OutputExcel)
	data, err := tbl.ToTableDataResponse(nil).generateExcel(nil)
	if err != nil {
		t.Fatal(err)
	}
	f, _ := excelize.OpenReader(bytes.NewReader(data))
	defer f.Close()
	typ, _ := f.GetCellType("Sheet1", "A2")
	raw, _ := f.GetCellValue("Sheet1", "A2", excelize.Options{RawCellValue: true})
	if typ == excelize.CellTypeSharedString || typ == excelize.CellTypeInlineString || raw != "1234.56" {
		t.Errorf("custom float formatter: type %v raw %q, want numeric 1234.56", typ, raw)
	}
}

func TestExportWeb_NumbersUnchanged(t *testing.T) {
	ctx := &core.UiContext{Locale: locale.De}
	if got := createFloatFormatter(2).Format(-12.5, nil, OutputWeb, ctx); got != "-12,50" {
		t.Errorf("web float = %v", got)
	}
	if got := createIntegerFormatter().Format(int64(1234567), nil, OutputPDF, ctx); got != "1.234.567" {
		t.Errorf("pdf int = %v", got)
	}
}

func TestExport_IntegerPrecisionAndNonFinite(t *testing.T) {
	const big = int64(9007199254740993) // 2^53 + 1
	if got := createIntegerFormatter().Format(big, nil, OutputCSV, nil); got != big {
		t.Errorf("big integer must stay exact int64, got %#v", got)
	}

	td := NewTableDataResponse([]map[string]any{
		{"nan": math.NaN(), "inf": ExportNumber{Value: math.Inf(1)}, "big": big},
	}, OutputCSV)
	td.withFieldsForCSV([]map[string]any{{"id": "nan", "name": "A"}, {"id": "inf", "name": "B"}, {"id": "big", "name": "C"}})
	if lines := csvLines(t, td, locale.De); lines[1] != "NaN;'+Inf;9007199254740993" {
		t.Errorf("row = %q", lines[1])
	}

	td.outputType = OutputExcel
	data, err := td.generateExcel(nil)
	if err != nil {
		t.Fatal(err)
	}
	f, _ := excelize.OpenReader(bytes.NewReader(data))
	defer f.Close()
	if v, _ := f.GetCellValue("Sheet1", "B2"); v != "+Inf" {
		t.Errorf("Excel +Inf must be plain text, got %q", v)
	}
}

type multiRow struct {
	FN  []float64
	IN  []int
	T2F [2]float64
	T2I [2]int
	Txt string
}

func multiTable(output OutputType, ctx *core.UiContext) *TableDataResponse {
	b := NewBuilder[multiRow]()
	b.FloatNField("fn", "FN", func(r multiRow) []float64 { return r.FN }).WithDecimals(2)
	b.IntNField("in", "IN", func(r multiRow) []int { return r.IN })
	b.Text2FloatField("t2f", "T2F", func(r multiRow) [2]float64 { return r.T2F }).WithDecimals(2)
	b.Text2IntField("t2i", "T2I", func(r multiRow) [2]int { return r.T2I })
	b.TextField("txt", "Txt", func(r multiRow) string { return r.Txt })
	tbl := b.Build()
	tbl.SetData([]multiRow{
		{FN: []float64{-6.5, 1.25}, IN: []int{-3, 4}, T2F: [2]float64{-12.5, 3}, T2I: [2]int{-5, 0}, Txt: "-12,50"},
		{FN: []float64{7}},
	})
	tbl.SetOutputType(output)
	return tbl.ToTableDataResponse(ctx)
}

// Negative numbers must not get a visible "'": N and Text2 fields become numeric cells, one per value.
func TestExportExcel_MultiValueNumbersWithoutApostrophe(t *testing.T) {
	ctx := &core.UiContext{Locale: locale.De}
	data, err := multiTable(OutputExcel, ctx).generateExcel(ctx)
	if err != nil {
		t.Fatal(err)
	}
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	header, _ := f.GetRows("Sheet1")
	if got := strings.Join(header[0], ";"); got != "FN;FN 2;IN;IN 2;T2F;T2F 2;T2I;T2I 2;Txt" {
		t.Errorf("header = %q", got)
	}

	for _, c := range []struct{ cell, raw, numFmt string }{
		{"A2", "-6.5", "#,##0.00"}, {"B2", "1.25", "#,##0.00"},
		{"C2", "-3", "#,##0"}, {"D2", "4", "#,##0"},
		{"E2", "-12.5", "#,##0.00"}, {"F2", "3", "#,##0.00"},
		{"G2", "-5", "#,##0"}, {"H2", "0", "#,##0"},
		{"A3", "7", "#,##0.00"},
	} {
		typ, _ := f.GetCellType("Sheet1", c.cell)
		raw, _ := f.GetCellValue("Sheet1", c.cell, excelize.Options{RawCellValue: true})
		if typ == excelize.CellTypeSharedString || typ == excelize.CellTypeInlineString || raw != c.raw {
			t.Errorf("%s: type %v raw %q, want numeric %q", c.cell, typ, raw, c.raw)
			continue
		}
		styleID, _ := f.GetCellStyle("Sheet1", c.cell)
		if style, err := f.GetStyle(styleID); err != nil || style.CustomNumFmt == nil || *style.CustomNumFmt != c.numFmt {
			t.Errorf("%s: number format %v, want %q", c.cell, style, c.numFmt)
		}
	}
	for _, cell := range []string{"B3", "C3", "D3"} {
		if v, _ := f.GetCellValue("Sheet1", cell); v != "" {
			t.Errorf("%s: missing value must stay empty, got %q", cell, v)
		}
	}
	if v, _ := f.GetCellValue("Sheet1", "I2"); v != "-12,50" {
		t.Errorf("text cell = %q, want -12,50 without apostrophe", v)
	}
}

func TestExportCSV_MultiValueNumbersFollowLocale(t *testing.T) {
	de := csvLines(t, multiTable(OutputCSV, &core.UiContext{Locale: locale.De}), locale.De)
	want := "FN;FN 2;IN;IN 2;T2F;T2F 2;T2I;T2I 2;Txt\n" +
		"-6,50;1,25;-3;4;-12,50;3,00;-5;0;-12,50\n" +
		"7,00;;;;0,00;0,00;0;0;"
	if got := strings.Join(de, "\n"); got != want {
		t.Errorf("CSV de:\n%s\nwant:\n%s", got, want)
	}
	en := csvLines(t, multiTable(OutputCSV, &core.UiContext{Locale: locale.EnUS}), locale.EnUS)
	if en[1] != "-6.50;1.25;-3;4;-12.50;3.00;-5;0;-12,50" {
		t.Errorf("CSV en row = %q", en[1])
	}
}

func TestExportCSV_DistanceSpeedNConverted(t *testing.T) {
	type r struct{ D, S []float64 }
	b := NewBuilder[r]()
	b.DistanceNField("d", "D", func(x r) []float64 { return x.D }).WithDecimals(1)
	b.SpeedNField("s", "S", func(x r) []float64 { return x.S }).WithDecimals(1)
	tbl := b.Build()
	tbl.SetData([]r{{D: []float64{-10, 100}, S: []float64{-50}}})
	tbl.SetOutputType(OutputCSV)
	ctx := &core.UiContext{Locale: locale.De, Distance: distance.Miles}
	lines := csvLines(t, tbl.ToTableDataResponse(ctx), locale.De)
	if lines[1] != "-6,2;62,1;-31,1" {
		t.Errorf("row = %q, want %q", lines[1], "-6,2;62,1;-31,1")
	}
}
