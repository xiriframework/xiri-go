package builder

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/xiriframework/xiri-go/form/field"
)

func suggestContext(body string) echo.Context {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return echo.New().NewContext(req, httptest.NewRecorder())
}

func TestBindSuggest_ReturnsSearchAndBindsOnlyGivenFields(t *testing.T) {
	dept := field.NewIntField("dept", "DEPT", false, 0)
	other := field.NewIntField("other", "OTHER", false, 1) // im Formular, aber kein Kontextfeld
	search, err := BindSuggest(suggestContext(`{"search":"gr","dept":3,"other":9}`), dept)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if search != "gr" {
		t.Errorf("expected search gr, got %q", search)
	}
	if dept.Value == nil || *dept.Value != 3 {
		t.Errorf("expected dept 3, got %v", dept.Value)
	}
	if other.Value != nil {
		t.Errorf("other must stay untouched, got %v", other.Value)
	}
}

func TestBindSuggest_LenientAndDefaults(t *testing.T) {
	dept := field.NewIntField("dept", "DEPT", false, 5)
	search, err := BindSuggest(suggestContext(`{"dept":"kaputt"}`), dept)
	if err != nil || search != "" {
		t.Errorf("expected empty search and no error, got %q / %v", search, err)
	}
	if dept.Value == nil || *dept.Value != 5 {
		t.Errorf("expected default 5 after failed bind, got %v", dept.Value)
	}
}

func TestBindSuggest_SkipsDisabledAndReservedSearch(t *testing.T) {
	disabled := field.NewIntField("dept", "DEPT", false, 5).SetDisabled(true)
	named := field.NewTextField("search", "S", false, "keep")
	_, err := BindSuggest(suggestContext(`{"search":"gr","dept":3}`), disabled, named)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if disabled.Value == nil || *disabled.Value != 5 {
		t.Errorf("disabled field must keep default, got %v", disabled.Value)
	}
	if named.Value == nil || *named.Value != "keep" {
		t.Errorf("field with id search must not receive the search text, got %v", named.Value)
	}
}

func TestBindSuggest_BadJSONIsError(t *testing.T) {
	if _, err := BindSuggest(suggestContext(`{"search":`)); err == nil {
		t.Error("expected error for malformed JSON")
	}
}
