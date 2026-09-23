package table_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/xiriframework/xiri-go/component/table"
	"github.com/xiriframework/xiri-go/form/field"
	"github.com/xiriframework/xiri-go/form/group"
)

func newJSONContext(body string) echo.Context {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec)
}

// #3: a malformed filter request body must return an error, not silently fall back to empty filters.
func TestLoadFilterDataMalformedBody(t *testing.T) {
	type row struct{}
	tbl := table.NewBuilder[row]().Build()

	c := newJSONContext("{ not valid json ")
	if _, err := tbl.LoadFilterData(c); err == nil {
		t.Fatal("expected error for malformed JSON body, got nil")
	}
}

// #3: an empty body still yields empty filters without error (documented "empty map if no filter").
func TestLoadFilterDataEmptyBody(t *testing.T) {
	type row struct{}
	tbl := table.NewBuilder[row]().Build()

	c := newJSONContext("")
	filters, err := tbl.LoadFilterData(c)
	if err != nil {
		t.Fatalf("empty body should not error, got %v", err)
	}
	if len(filters) != 0 {
		t.Errorf("empty body should yield empty filters, got %v", filters)
	}
}

// Ein leerer Filter darf nicht filtern: fehlende Keys bekommen keinen Field-Default
// (NewSelectField würde sonst die erste Option liefern).
func TestLoadFilterDataMissingKeyHasNoDefault(t *testing.T) {
	type row struct{}
	status := field.NewSelectField("status", "Status", false, []field.SelectOption{
		{Value: "online", Label: "Online"},
		{Value: "offline", Label: "Offline"},
	})
	fg := group.NewFormGroup([]field.FormField{status})

	b := table.NewBuilder[row]()
	b.SetFilter(fg)
	tbl := b.Build()

	filters, err := tbl.LoadFilterData(newJSONContext(`{"page": 1}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := filters["status"]; ok {
		t.Errorf("missing filter key must not be defaulted, got status=%v", v)
	}

	filters, err = tbl.LoadFilterData(newJSONContext(`{"status": "offline"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if filters["status"] != "offline" {
		t.Errorf("present key must be parsed, got %v", filters["status"])
	}
}

// Ein gesperrter (disabled) Filter darf nicht per Request-Body überschrieben werden.
func TestLoadFilterDataDisabledFilterIgnoresClient(t *testing.T) {
	type row struct{}
	owner := field.NewIntField("owner", "Owner", false, 7)
	owner.SetDisabled(true)
	fg := group.NewFormGroup([]field.FormField{owner})

	b := table.NewBuilder[row]()
	b.SetFilter(fg)
	tbl := b.Build()

	filters, err := tbl.LoadFilterData(newJSONContext(`{"owner": 999}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if filters["owner"] != int32(7) {
		t.Errorf("owner = %#v, want int32(7)", filters["owner"])
	}
}
