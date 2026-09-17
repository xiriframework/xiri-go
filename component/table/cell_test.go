package table_test

import (
	"encoding/json"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/xiriframework/xiri-go/component/table"
	"github.com/xiriframework/xiri-go/formatter"
)

type cellRow struct {
	When time.Time
	Span [2]time.Time
	Many []time.Time
	Secs int64
	N    int
}

func cellTable() *table.Table[cellRow] {
	b := table.NewBuilder[cellRow]()
	b.DateField("d", "d", func(r cellRow) time.Time { return r.When })
	b.DateTimeField("dt", "dt", func(r cellRow) time.Time { return r.When })
	b.Text2DateField("d2", "d2", func(r cellRow) [2]time.Time { return r.Span })
	b.DateNField("dn", "dn", func(r cellRow) []time.Time { return r.Many })
	b.TimeLengthField("tl", "tl", func(r cellRow) int64 { return r.Secs })
	b.IntField("n", "n", func(r cellRow) int { return r.N })
	b.TextField("t", "t", func(r cellRow) string { return "x" })
	return b.Build()
}

func obj(d any, v any) map[string]any { return map[string]any{"d": d, "v": v} }

// Web cells of cell-object fields are {d: display, v: raw value}; v is ISO in the user's timezone.
func TestCellObjectWebCells(t *testing.T) {
	ctx := exampleContext()                               // Europe/Vienna, De
	when := time.Date(2024, 2, 24, 14, 5, 0, 0, time.UTC) // 15:05 in Vienna
	ts := when.Unix()
	tbl := cellTable()
	tbl.SetData([]cellRow{{When: when, Span: [2]time.Time{when, when.Add(time.Hour)}, Many: []time.Time{when}, Secs: 100 * 3600}})
	row := tbl.GetData(ctx, table.OutputWeb)[0]

	want := map[string]any{
		"d":  obj(formatter.FormatTimestampDate(ts, ctx), "2024-02-24"),
		"dt": obj(formatter.FormatTimestampDateTime(ts, ctx), "2024-02-24T15:05:00"),
		"d2": obj([2]string{formatter.FormatTimestampDate(ts, ctx), formatter.FormatTimestampDate(ts+3600, ctx)}, "2024-02-24"),
		"dn": obj([]string{formatter.FormatTimestampDate(ts, ctx)}, "2024-02-24"),
		"tl": obj("4d 04:00", int64(360000)),
		"t":  "x",
	}
	for id, w := range want {
		if got := row[id]; !reflect.DeepEqual(got, w) {
			t.Errorf("%s = %#v, want %#v", id, got, w)
		}
	}
	if n, ok := row["n"].([]any); !ok || len(n) != 2 {
		t.Errorf("number field must keep its legacy [display, value] pair, got %#v", row["n"])
	}
}

// Empty values carry v == nil, never 0 or "": zero time, empty slice, negative duration.
func TestCellObjectEmptyIsNil(t *testing.T) {
	tbl := cellTable()
	tbl.SetData([]cellRow{{Secs: -1}})
	row := tbl.GetData(exampleContext(), table.OutputWeb)[0]

	want := map[string]any{
		"d":  obj("", nil),
		"d2": obj([2]string{"", ""}, nil),
		"dn": obj([]string{}, nil),
		"tl": obj("", nil),
	}
	for id, w := range want {
		if got := row[id]; !reflect.DeepEqual(got, w) {
			t.Errorf("%s = %#v, want %#v", id, got, w)
		}
	}
}

func TestCellObjectZeroDurationKeepsZero(t *testing.T) {
	tbl := cellTable()
	tbl.SetData([]cellRow{{}})
	if got := tbl.GetData(exampleContext(), table.OutputWeb)[0]["tl"]; !reflect.DeepEqual(got, obj("00:00", int64(0))) {
		t.Errorf("tl = %#v, want {00:00 0}", got)
	}
}

func TestCellObjectNotOnCSV(t *testing.T) {
	tbl := cellTable()
	tbl.SetData([]cellRow{{When: time.Unix(1708732800, 0)}})
	if _, ok := tbl.GetData(exampleContext(), table.OutputCSV)[0]["d"].(string); !ok {
		t.Errorf("CSV date cell must stay a plain string")
	}
}

func TestCellObjectFieldJSON(t *testing.T) {
	byID := map[string]map[string]any{}
	for _, f := range cellTable().ExportFields(exampleContext()) {
		byID[f["id"].(string)] = f
	}
	kinds := map[string]string{"d": "string", "dt": "string", "d2": "string", "dn": "string", "tl": "number"}
	for id, kind := range kinds {
		if byID[id]["cellObject"] != kind {
			t.Errorf("field %s: cellObject = %v, want %q", id, byID[id]["cellObject"], kind)
		}
	}
	for _, id := range []string{"n", "t"} {
		if _, ok := byID[id]["cellObject"]; ok {
			t.Errorf("field %s must not carry cellObject", id)
		}
	}
	if byID["d"]["inputType"] != "date" || byID["dt"]["inputType"] != "datetime-local" || byID["tl"]["inputType"] != "number" {
		t.Errorf("date/dateTime/timeLength need inputType defaults, got %v / %v / %v",
			byID["d"]["inputType"], byID["dt"]["inputType"], byID["tl"]["inputType"])
	}
	if _, ok := byID["t"]["inputType"]; ok {
		t.Errorf("text must not get an inputType default")
	}
}

// Cell builds one web cell the way GetData does — for ReturnInlineEdit.Updates after a save.
func TestCellMatchesGetData(t *testing.T) {
	ctx := exampleContext()
	r := cellRow{When: time.Date(2024, 2, 25, 9, 0, 0, 0, time.UTC), Secs: 90}
	tbl := cellTable()
	tbl.SetData([]cellRow{r})
	full := tbl.GetData(ctx, table.OutputWeb)[0]
	for _, id := range []string{"d", "dt", "tl", "n", "t"} {
		if got := tbl.Cell(ctx, id, r); !reflect.DeepEqual(got, full[id]) {
			t.Errorf("Cell(%s) = %#v, GetData gives %#v", id, got, full[id])
		}
	}
	if got := tbl.Cell(ctx, "nope", r); got != nil {
		t.Errorf("unknown field must yield nil, got %#v", got)
	}
}

// Cell refuses fields whose GetData output is more than the formatted value (link split, menu data).
func TestCellRefusesLinkAndButtons(t *testing.T) {
	b := table.NewBuilder[cellRow]()
	b.LinkField("l", "l", func(r cellRow) [2]string { return [2]string{"x", "/x"} })
	b.ButtonsField("b", "b", func(r cellRow) map[string]string { return nil }).
		AddButton(0, table.FieldButtonActionLink, "edit", "", "")
	tbl := b.Build()
	for _, id := range []string{"l", "b"} {
		if got := tbl.Cell(exampleContext(), id, cellRow{}); got != nil {
			t.Errorf("Cell(%s) must be nil, got %#v", id, got)
		}
	}
}

func TestCellObjectFloatsEncode(t *testing.T) {
	type r struct{ F [2]float64 }
	b := table.NewBuilder[r]()
	b.Text2FloatField("f", "f", func(x r) [2]float64 { return x.F })
	tbl := b.Build()
	tbl.SetData([]r{{F: [2]float64{math.NaN(), 1}}, {F: [2]float64{math.Inf(1), 1}}, {F: [2]float64{2.5, 1}}})
	data := tbl.GetData(exampleContext(), table.OutputWeb)
	if _, err := json.Marshal(data); err != nil {
		t.Fatalf("web data must encode: %v", err)
	}
	if v := data[0]["f"].(map[string]any)["v"]; v != nil {
		t.Errorf("NaN must become nil, got %#v", v)
	}
	if v := data[2]["f"].(map[string]any)["v"]; v != 2.5 {
		t.Errorf("finite float must pass through, got %#v", v)
	}
}

// CalculateFooter reuses field.format, so footers of cell-object fields are {d, v} too on web
// output — the client (Task 1, cellDisplay) unpacks server footers the same way as row cells.
func TestCellObjectFooterIsCellObject(t *testing.T) {
	b := table.NewBuilder[cellRow]()
	b.TimeLengthField("tl", "tl", func(r cellRow) int64 { return r.Secs }).WithFooterSum()
	tbl := b.Build()
	tbl.SetData([]cellRow{{Secs: 60}, {Secs: 120}})

	footer := tbl.CalculateFooter(exampleContext(), table.OutputWeb)
	got, ok := footer["tl"].(map[string]any)
	if !ok {
		t.Fatalf("tl footer = %#v, want map[string]any{d, v}", footer["tl"])
	}
	if got["d"] != "00:03" {
		t.Errorf("tl footer d = %#v, want %q", got["d"], "00:03")
	}
	// sumField aggregates via toFloat64, but cellValueFor's toInt64 converts it back for v.
	if got["v"] != int64(180) {
		t.Errorf("tl footer v = %#v, want %#v", got["v"], int64(180))
	}

	if got := tbl.CalculateFooter(exampleContext(), table.OutputCSV)["tl"]; got != "3" {
		t.Errorf("tl footer (CSV) = %#v, want plain string %q", got, "3")
	}
}
