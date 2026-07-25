package table_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/xiriframework/xiri-go/component/table"
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
