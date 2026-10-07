package field

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	"github.com/xiriframework/xiri-go/component/url"
)

// captureWarn routes slog into a buffer for the duration of the test.
func captureWarn(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	t.Cleanup(func() { slog.SetDefault(prev) })
	return &buf
}

// selectAll is only exported in multi-select mode — without Multiple the toggle silently vanishes.
func TestSelectField_SelectAllWithoutMultiple_WarnsOnce(t *testing.T) {
	buf := captureWarn(t)
	f := NewSelectField("status", "Status", false, selectTestOptions()).SetSelectAll(true)
	f.ExportForFrontend(nil, nil)
	f.ExportForFrontend(nil, nil)

	if n := strings.Count(buf.String(), "selectAll"); n != 1 {
		t.Errorf("expected exactly one selectAll warning, got %d: %q", n, buf.String())
	}
}

func TestSelectField_SelectAllWithMultiple_DoesNotWarn(t *testing.T) {
	buf := captureWarn(t)
	NewSelectField("status", "Status", false, selectTestOptions()).SetSelectAll(true).SetMultiple(true).ExportForFrontend(nil, nil)
	NewSelectField("status", "Status", false, selectTestOptions()).SetSelectAll(false).ExportForFrontend(nil, nil)

	if buf.Len() != 0 {
		t.Errorf("expected no warning, got %q", buf.String())
	}
}

// xiri-ng renders no suggestions on a textarea (subtypes textarea and html).
func TestTextField_TextareaWithSuggestions_WarnsOnce(t *testing.T) {
	for _, subtype := range []string{"textarea", "html"} {
		buf := captureWarn(t)
		f := NewTextField("note", "Notiz", false, "").SetSuggestions("a")
		f.Subtype = subtype
		f.ExportForFrontend(nil, nil)
		f.ExportForFrontend(nil, nil)

		if n := strings.Count(buf.String(), "suggestions"); n != 1 {
			t.Errorf("%s: expected exactly one suggestions warning, got %d: %q", subtype, n, buf.String())
		}
	}
}

func TestTextField_TextareaWithSuggestionsURL_Warns(t *testing.T) {
	buf := captureWarn(t)
	f := NewTextField("note", "Notiz", false, "").SetSuggestionsURL(url.NewUrlPrefix("/suggest", "/api"))
	f.Subtype = "textarea"
	f.ExportForFrontend(nil, nil)

	if !strings.Contains(buf.String(), "suggestions") {
		t.Errorf("expected a suggestions warning, got %q", buf.String())
	}
}

func TestTextField_SuggestionsOnTextOrEmptyOnTextarea_DoesNotWarn(t *testing.T) {
	buf := captureWarn(t)
	NewTextField("city", "Stadt", false, "").SetSuggestions("Wien").ExportForFrontend(nil, nil)
	// An explicitly empty list only clears suggestions via reload patch — nothing is lost.
	f := NewTextField("note", "Notiz", false, "").SetSuggestions()
	f.Subtype = "textarea"
	f.ExportForFrontend(nil, nil)

	if buf.Len() != 0 {
		t.Errorf("expected no warning, got %q", buf.String())
	}
}
