package table

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/types/locale"
	"github.com/xiriframework/xiri-go/types/pressure"
)

// #8: distance/speed/pressure formatters must honor the nil-context contract (Print(nil)).
func TestFormattersNilContextNoPanic(t *testing.T) {
	cases := []struct {
		name  string
		f     OutputFormatter
		value any
	}{
		{"distance", createDistanceFormatter(2), 1.0},
		{"speed", createSpeedFormatter(2), 1.0},
		{"pressure", createPressureFormatter(2), 1.0},
		{"text2distance", createText2DistanceFormatter(2), [2]float64{1, 2}},
		{"text2speed", createText2SpeedFormatter(2), [2]float64{1, 2}},
		{"distanceN", createDistanceNFormatter(2), []float64{1, 2}},
		{"speedN", createSpeedNFormatter(2), []float64{1, 2}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%s panicked with nil ctx: %v", tc.name, r)
				}
			}()
			tc.f.Format(tc.value, nil, OutputWeb, nil)
		})
	}
}

// #13: the pressure table formatter must support kPa, not fall back to the bar value.
func TestPressureFormatterKpa(t *testing.T) {
	f := createPressureFormatter(1)
	ctx := &core.UiContext{Pressure: pressure.Kpa, Locale: locale.EnUS}

	got := f.Format(2.0, nil, OutputWeb, ctx)
	if got != "200.0 kPa" {
		t.Errorf("pressure kPa (web) = %v, want %q", got, "200.0 kPa")
	}
}
