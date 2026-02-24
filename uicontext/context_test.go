package uicontext

import "testing"

func TestSafeTranslate_NilContext(t *testing.T) {
	var ctx *UiContext
	result := ctx.SafeTranslate("HELLO")
	if result != "HELLO" {
		t.Errorf("expected 'HELLO', got %q", result)
	}
}

func TestSafeTranslate_NilTranslateFunc(t *testing.T) {
	ctx := &UiContext{}
	result := ctx.SafeTranslate("HELLO")
	if result != "HELLO" {
		t.Errorf("expected 'HELLO', got %q", result)
	}
}

func TestSafeTranslate_WithTranslator(t *testing.T) {
	ctx := &UiContext{
		Translate: func(key string) string {
			if key == "HELLO" {
				return "Hallo"
			}
			return key
		},
	}
	result := ctx.SafeTranslate("HELLO")
	if result != "Hallo" {
		t.Errorf("expected 'Hallo', got %q", result)
	}
}
