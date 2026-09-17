package formatter

import (
	"testing"
	"time"

	"github.com/xiriframework/xiri-go/component/core"
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
