package formatter

import (
	"testing"
	"time"

	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/types/locale"
	"github.com/xiriframework/xiri-go/types/timezone"
)

// v of a dateTime cell is local time without zone; the server parses it back in the user's timezone.
func TestParseLocalDateTimeRoundtrip(t *testing.T) {
	ctx := &core.UiContext{Timezone: timezone.EuropeVienna}
	vienna, _ := time.LoadLocation("Europe/Vienna")
	want := time.Date(2024, 2, 24, 15, 5, 30, 0, vienna)

	cases := map[string]time.Time{
		"2024-02-24T15:05:30": want,
		"2024-02-24T15:05":    want.Truncate(time.Minute), // browsers drop seconds without step="1"
		"2024-02-24":          time.Date(2024, 2, 24, 0, 0, 0, 0, vienna),
	}
	for in, w := range cases {
		got, err := ParseLocalDateTime(in, ctx)
		if err != nil || !got.Equal(w) {
			t.Errorf("ParseLocalDateTime(%q) = %v, %v; want %v", in, got, err, w)
		}
	}
	if _, err := ParseLocalDateTime("24.02.2024", ctx); err == nil {
		t.Errorf("locale strings must be rejected")
	}
	// 2024-03-31 02:30 does not exist in Vienna (clocks jump 02:00 → 03:00): reject instead of shifting.
	if got, err := ParseLocalDateTime("2024-03-31T02:30:00", ctx); err == nil {
		t.Errorf("spring-forward gap must be an error, got %v", got)
	}
}

// One row per locale — a representative per group would miss exactly the misfiled locales this pins.
// Formats are deliberately simplified product formats (numeric, no CLDR spaces/trailing dots, 12h only EnUS).
func TestLayoutsPerLocale(t *testing.T) {
	at := time.Date(2024, 2, 24, 14, 5, 0, 0, time.UTC) // 15:05 in Europe/Vienna

	dotDMY := [3]string{"24.02.2024", "15:05", "24.02.2024 15:05"}
	slashDMY := [3]string{"24/02/2024", "15:05", "24/02/2024 15:05"}
	slashYMD := [3]string{"2024/02/24", "15:05", "2024/02/24 15:05"}
	want := map[locale.Locale][3]string{
		locale.De: dotDMY, locale.DeAT: dotDMY, locale.DeCH: dotDMY,
		locale.Hr: dotDMY, locale.Pl: dotDMY, locale.Cs: dotDMY, locale.Ro: dotDMY, locale.Tr: dotDMY,
		locale.Bg: dotDMY, locale.Sl: dotDMY, locale.Sk: dotDMY, locale.Sr: dotDMY,
		locale.Nb: dotDMY, locale.Da: dotDMY, locale.Fi: dotDMY, locale.Ru: dotDMY, locale.Uk: dotDMY,
		locale.EnGB: slashDMY, locale.Es: slashDMY, locale.Fr: slashDMY, locale.It: slashDMY,
		locale.Pt: slashDMY, locale.PtBR: slashDMY, locale.El: slashDMY, locale.ArAE: slashDMY,
		locale.Nl:   {"24-02-2024", "15:05", "24-02-2024 15:05"},
		locale.Hu:   {"2024.02.24", "15:05", "2024.02.24 15:05"},
		locale.Sv:   {"2024-02-24", "15:05", "2024-02-24 15:05"},
		locale.EnUS: {"02/24/2024", "03:05 PM", "02/24/2024 03:05 PM"},
		locale.Ja:   slashYMD, locale.ZhCN: slashYMD,
	}

	for _, loc := range locale.All() {
		w, ok := want[loc]
		if !ok {
			t.Errorf("%s: no expectation — every locale needs a row", loc)
			continue
		}
		ctx := &core.UiContext{Locale: loc, Timezone: timezone.EuropeVienna}
		if got := FormatDate(at, ctx); got != w[0] {
			t.Errorf("%s FormatDate = %q, want %q", loc, got, w[0])
		}
		if got := FormatTime(at, ctx); got != w[1] {
			t.Errorf("%s FormatTime = %q, want %q", loc, got, w[1])
		}
		if got := FormatDateTime(at, ctx); got != w[2] {
			t.Errorf("%s FormatDateTime = %q, want %q", loc, got, w[2])
		}
	}
	if len(want) != len(locale.All()) {
		t.Errorf("expectation table has %d rows, locale.All() has %d", len(want), len(locale.All()))
	}
}
