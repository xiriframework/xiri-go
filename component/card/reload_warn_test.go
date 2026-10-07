package card

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
)

func captureWarn(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	t.Cleanup(func() { slog.SetDefault(prev) })
	return &buf
}

// reload is only exported in URL mode — a static card gets no reload button.
func TestCard_ReloadWithoutURL_WarnsOnce(t *testing.T) {
	buf := captureWarn(t)
	c := NewCard(core.CardTypeTable, nil, "Status", nil, nil, nil, false, false, nil).WithReload(true)
	c.Print(cardCtx())
	c.Print(cardCtx())

	if n := strings.Count(buf.String(), "WithReload"); n != 1 {
		t.Errorf("expected exactly one WithReload warning, got %d: %q", n, buf.String())
	}
}

// The data endpoint typically rebuilds the same card without URL and answers via DataResponse.
func TestCard_ReloadWithURLOrDataResponse_DoesNotWarn(t *testing.T) {
	buf := captureWarn(t)
	NewCard(core.CardTypeTable, nil, "Status", nil, nil, nil, false, false, nil).
		WithReload(true).SetURL(url.NewUrlPrefix("/card", "/api")).Print(cardCtx())
	NewCard(core.CardTypeTable, nil, "Status", nil, nil, nil, false, false, nil).
		WithReload(true).DataResponse(cardCtx())
	NewCard(core.CardTypeTable, nil, "Status", nil, nil, nil, false, false, nil).
		WithReload(false).Print(cardCtx())

	if buf.Len() != 0 {
		t.Errorf("expected no warning, got %q", buf.String())
	}
}
