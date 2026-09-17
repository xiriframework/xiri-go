package formatter

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/types/locale"
)

func TestFormatTimeLengthHUsesLocaleDecimalSeparator(t *testing.T) {
	de := &core.UiContext{Locale: locale.De}
	en := &core.UiContext{Locale: locale.EnUS}

	if got := FormatTimeLengthH(5400, de); got != "1,5 h" {
		t.Errorf("De: want %q, got %q", "1,5 h", got)
	}
	if got := FormatTimeLengthH(5400, en); got != "1.5 h" {
		t.Errorf("EnUS: want %q, got %q", "1.5 h", got)
	}
	if got := FormatTimeLengthH(-1, nil); got != "0,0 h" {
		t.Errorf("negative/nil ctx: want %q, got %q", "0,0 h", got)
	}
}
