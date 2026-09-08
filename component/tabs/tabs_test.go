package tabs_test

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/tabs"
)

func TestTab_WithNoPadding_PrintsNoPadding(t *testing.T) {
	out := tabs.NewTab("x").WithNoPadding(true).Print(nil)

	if out["noPadding"] != true {
		t.Errorf(`out["noPadding"] = %v, want true`, out["noPadding"])
	}
}

func TestTab_WithNoPaddingFalse_PrintsFalse(t *testing.T) {
	out := tabs.NewTab("x").WithNoPadding(false).Print(nil)

	v, has := out["noPadding"]
	if !has {
		t.Fatal(`expected "noPadding" key when explicitly set to false`)
	}
	if v != false {
		t.Errorf(`out["noPadding"] = %v, want false`, v)
	}
}

func TestTab_Default_OmitsNoPadding(t *testing.T) {
	out := tabs.NewTab("x").Print(nil)

	if _, has := out["noPadding"]; has {
		t.Errorf(`expected no "noPadding" key, got %v`, out["noPadding"])
	}
}
