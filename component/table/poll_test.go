package table

import (
	"testing"

	xurl "github.com/xiriframework/xiri-go/component/url"
)

// TestWithPollEmittedForWeb verifies WithPoll adds the "poll" field for Web output.
func TestWithPollEmittedForWeb(t *testing.T) {
	ctx := testContext()

	rows := []map[string]any{{"id": int64(1), "name": "A"}}
	td := NewTableDataResponse(rows, OutputWeb).WithPoll(2000)

	out := td.Print(ctx)
	poll, ok := out["poll"].(int)
	if !ok {
		t.Fatalf("Expected poll to be int, got %T (%v)", out["poll"], out["poll"])
	}
	if poll != 2000 {
		t.Errorf("Expected poll 2000, got %d", poll)
	}
}

// TestWithoutPollOmitsField verifies the "poll" field is absent when WithPoll is not used.
func TestWithoutPollOmitsField(t *testing.T) {
	ctx := testContext()

	rows := []map[string]any{{"id": int64(1), "name": "A"}}
	out := NewTableDataResponse(rows, OutputWeb).Print(ctx)

	if _, exists := out["poll"]; exists {
		t.Errorf("Expected no poll field, got %v", out["poll"])
	}
}

// TestPollZeroOmitsField verifies WithPoll(0) does not emit the field.
func TestPollZeroOmitsField(t *testing.T) {
	ctx := testContext()

	rows := []map[string]any{{"id": int64(1), "name": "A"}}
	out := NewTableDataResponse(rows, OutputWeb).WithPoll(0).Print(ctx)

	if _, exists := out["poll"]; exists {
		t.Errorf("Expected no poll field for WithPoll(0), got %v", out["poll"])
	}
}

// TestPollNotEmittedForCSV verifies the "poll" field is never emitted for CSV output.
func TestPollNotEmittedForCSV(t *testing.T) {
	ctx := testContext()

	rows := []map[string]any{{"id": int64(1), "name": "A"}}
	out := NewTableDataResponse(rows, OutputCSV).WithPoll(2000).Print(ctx)

	if _, exists := out["poll"]; exists {
		t.Errorf("Expected no poll field for CSV output, got %v", out["poll"])
	}
}

// TestSetPollFlowsThroughToTableDataResponse verifies tbl.SetPoll propagates to the response.
func TestSetPollFlowsThroughToTableDataResponse(t *testing.T) {
	ctx := testContext()

	builder := NewBuilder[testDeviceRow]()
	builder.IdField("id", "device.id", func(r testDeviceRow) int64 { return r.ID })
	builder.TextField("name", "device.name", func(r testDeviceRow) string { return r.Name })
	tbl := builder.Build()
	tbl.SetURL(xurl.NewUrl("/Portal/Device/TableData"))
	tbl.SetData([]testDeviceRow{{ID: 1, Name: "A"}})
	tbl.SetPoll(3000)

	out := tbl.ToTableDataResponse(ctx).Print(ctx)
	poll, ok := out["poll"].(int)
	if !ok {
		t.Fatalf("Expected poll to be int, got %T (%v)", out["poll"], out["poll"])
	}
	if poll != 3000 {
		t.Errorf("Expected poll 3000, got %d", poll)
	}
}
