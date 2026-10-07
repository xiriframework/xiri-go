package table

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

// saveState is only exported together with saveStateId — without the id the state is never saved.
func TestBuild_SaveStateWithoutId_Warns(t *testing.T) {
	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	defer slog.SetDefault(prev)

	NewBuilder[testOptionRow]().SetSaveState(true).Build()

	if !strings.Contains(buf.String(), "SaveStateId") {
		t.Errorf("expected a warning about the missing SaveStateId, got %q", buf.String())
	}
}

func TestBuild_SaveStateWithIdOrOff_DoesNotWarn(t *testing.T) {
	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	defer slog.SetDefault(prev)

	NewBuilder[testOptionRow]().SetSaveState(true).SetSaveStateId("devices").Build()
	NewBuilder[testOptionRow]().SetSaveState(false).Build()

	if buf.Len() != 0 {
		t.Errorf("expected no warning, got %q", buf.String())
	}
}
