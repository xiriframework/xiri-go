package field

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/types/distance"
	"github.com/xiriframework/xiri-go/types/language"
	"github.com/xiriframework/xiri-go/types/locale"
	"github.com/xiriframework/xiri-go/types/pressure"
	"github.com/xiriframework/xiri-go/types/timezone"
)

func newTestCtx() *core.UiContext {
	return &core.UiContext{
		Timezone: timezone.EuropeVienna,
		Lang:     language.Deutsch,
		Locale:   locale.De,
		Distance: distance.Kilometer,
		Pressure: pressure.Bar,
	}
}

func TestBaseFieldSetHide(t *testing.T) {
	f := NewTextField("secret", "SECRET", false, "")

	if f.IsHidden() {
		t.Fatal("expected IsHidden() to be false by default")
	}

	ret := f.SetHide(true)
	if ret == nil {
		t.Fatal("SetHide returned nil, expected chainable *BaseField")
	}
	if !f.IsHidden() {
		t.Fatal("expected IsHidden() to be true after SetHide(true)")
	}
	if !f.Hide {
		t.Fatal("expected Hide field to be true")
	}

	f.SetHide(false)
	if f.IsHidden() {
		t.Fatal("expected IsHidden() to be false after SetHide(false)")
	}
}

func TestBaseFieldExportHide_True(t *testing.T) {
	ctx := newTestCtx()
	f := NewTextField("secret", "SECRET", false, "v")
	f.SetHide(true)

	result := f.ExportForFrontend(ctx, "v")

	hide, ok := result["hide"]
	if !ok {
		t.Fatal("expected 'hide' key in export")
	}
	if hide != true {
		t.Errorf("expected hide=true, got %v", hide)
	}
}

func TestBaseFieldExportHide_FalseAlwaysEmitted(t *testing.T) {
	ctx := newTestCtx()
	f := NewTextField("plain", "PLAIN", false, "")

	result := f.ExportForFrontend(ctx, "")

	hide, ok := result["hide"]
	if !ok {
		t.Fatal("expected 'hide' key to always be present in export")
	}
	if hide != false {
		t.Errorf("expected hide=false, got %v", hide)
	}
}

func TestBaseFieldHideIndependentOfForm(t *testing.T) {
	f := NewTextField("both", "BOTH", false, "")
	f.SetHide(true)
	f.SetForm(false)

	if !f.IsHidden() {
		t.Error("Hide should stay true when Form is set false")
	}
	if f.GetForm() {
		t.Error("Form should be false")
	}
}
