package builder

import (
	"testing"

	"github.com/xiriframework/xiri-go/form/field"
	"github.com/xiriframework/xiri-go/form/group"
	"github.com/xiriframework/xiri-go/response"
)

func TestBindAndValidate_PatternViolationReportedPerField(t *testing.T) {
	code := field.NewTextField("code", "CODE", false, "")
	code.Pattern = "[a-z]+"
	fg := group.NewFormGroup([]field.FormField{code})

	err := BindAndValidate(suggestContext(`{"code":"ab1"}`), fg)
	if err == nil {
		t.Fatal("expected pattern error")
	}

	resp := response.NewErrorResponseFromError(err)
	if _, ok := resp.Fields["code"]; !ok {
		t.Fatalf("expected fields.code, got %#v (err %v)", resp.Fields, err)
	}
}

func TestBindFromMap_PatternChecksDefaultOfMissingField(t *testing.T) {
	// Defaults laufen nicht durch Parse und werden trotzdem geprüft: ein Altwert, der zum Pattern
	// nicht passt, blockiert das Speichern (z. B. bei per showWhen verstecktem Feld).
	code := field.NewTextField("code", "CODE", false, "ab1")
	code.Pattern = "[a-z]+"
	fg := group.NewFormGroup([]field.FormField{code})

	if err := BindFromMap(map[string]interface{}{}, fg); err == nil {
		t.Fatal("expected pattern error for default value")
	}
}
