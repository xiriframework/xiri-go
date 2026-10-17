package dialog

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
	"github.com/xiriframework/xiri-go/response"
)

func newDelContext(method string) (echo.Context, *httptest.ResponseRecorder) {
	rec := httptest.NewRecorder()
	return echo.New().NewContext(httptest.NewRequest(method, "/", nil), rec), rec
}

func buttonTexts(t *testing.T, rec *httptest.ResponseRecorder) (closeText, okText any) {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &m); err != nil {
		t.Fatalf("invalid JSON %q: %v", rec.Body.String(), err)
	}
	buttons, ok := m["buttons"].([]any)
	if !ok || len(buttons) != 2 {
		t.Fatalf("buttons = %v, want 2 buttons", m["buttons"])
	}
	return buttons[0].(map[string]any)["text"], buttons[1].(map[string]any)["text"]
}

func TestHandleDelRequestWithDialog_Get(t *testing.T) {
	c, rec := newDelContext(http.MethodGet)
	dlg := NewDialogDelete("Weg?", url.NewUrl("/x"), nil, nil, stringPtr("Löschen"), stringPtr("Abbrechen"))

	err := HandleDelRequestWithDialog(c, dlg, func() error {
		t.Fatal("deleteFunc must not run on GET")
		return nil
	}, response.NewReturnRefreshTable(), &core.UiContext{})
	if err != nil {
		t.Fatal(err)
	}

	closeText, okText := buttonTexts(t, rec)
	if closeText != "Abbrechen" || okText != "Löschen" {
		t.Errorf("buttons = %v/%v, want Abbrechen/Löschen", closeText, okText)
	}
}

func TestHandleDelRequestWithDialog_Post(t *testing.T) {
	dlg := NewDialogWarning("Sicher?", url.NewUrl("/x"), nil, nil, nil, nil)

	t.Run("success", func(t *testing.T) {
		c, rec := newDelContext(http.MethodPost)
		called := false
		err := HandleDelRequestWithDialog(c, dlg, func() error {
			called = true
			return nil
		}, response.NewReturnRefreshTable(), &core.UiContext{})
		if err != nil {
			t.Fatal(err)
		}
		if !called {
			t.Error("deleteFunc not called")
		}
		if rec.Code != http.StatusOK {
			t.Errorf("status = %d, want 200", rec.Code)
		}
	})

	t.Run("error", func(t *testing.T) {
		c, rec := newDelContext(http.MethodPost)
		err := HandleDelRequestWithDialog(c, dlg, func() error {
			return errors.New("boom")
		}, response.NewReturnRefreshTable(), &core.UiContext{})
		if err != nil {
			t.Fatal(err)
		}
		if rec.Code != http.StatusInternalServerError {
			t.Errorf("status = %d, want 500", rec.Code)
		}
	})
}

func TestHandleDelRequest_DefaultTexts(t *testing.T) {
	c, rec := newDelContext(http.MethodGet)

	err := HandleDelRequest(c, "Weg?", "/x", func() error {
		t.Fatal("deleteFunc must not run on GET")
		return nil
	}, response.NewReturnRefreshTable(), &core.UiContext{})
	if err != nil {
		t.Fatal(err)
	}

	closeText, okText := buttonTexts(t, rec)
	if closeText != "Back" || okText != "Ok" {
		t.Errorf("buttons = %v/%v, want Back/Ok", closeText, okText)
	}
}
