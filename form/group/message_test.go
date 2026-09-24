package group

import (
	"errors"
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/form/field"
)

func deCtx(texts map[string]string) *core.UiContext {
	return &core.UiContext{Translate: func(key string) string {
		if t, ok := texts[key]; ok {
			return t
		}
		return key
	}}
}

func groupWithCtx(t *testing.T, ctx *core.UiContext, fields ...field.FormField) *FormGroup {
	t.Helper()
	fg, err := NewFormGroupWithContext(fields, ctx)
	if err != nil {
		t.Fatal(err)
	}
	return fg
}

var deTexts = map[string]string{
	"ORT":                    "Ort",
	"ART":                    "Art",
	"validation.required":    "{field} ist ein Pflichtfeld",
	"validation.max_length":  "Höchstens {max} Zeichen",
	"validation.not_allowed": "Ungültige Auswahl",
}

func fieldErrors(t *testing.T, err error) FieldErrors {
	t.Helper()
	var fe FieldErrors
	if !errors.As(err, &fe) {
		t.Fatalf("expected FieldErrors, got %T: %v", err, err)
	}
	return fe
}

func TestParseAndValidate_TranslatesCodeWithParamsAndLabel(t *testing.T) {
	ort := field.NewTextFieldWithLength("ort", "ORT", false, "", 0, 5)
	fg := groupWithCtx(t, deCtx(deTexts), ort)

	_, err := fg.ParseAndValidate(map[string]interface{}{"ort": "Wien-Floridsdorf"})

	if got := fieldErrors(t, err)["ort"]; got != "Höchstens 5 Zeichen" {
		t.Fatalf("field message: %q", got)
	}
	if got := err.Error(); got != "Ort: Höchstens 5 Zeichen" {
		t.Fatalf("banner: %q", got)
	}
}

func TestParseValues_MissingRequired_Translated(t *testing.T) {
	ort := field.NewTextField("ort", "ORT", true, "")
	fg := groupWithCtx(t, deCtx(deTexts), ort)

	_, err := fg.ParseValues(map[string]interface{}{})

	if got := fieldErrors(t, err)["ort"]; got != "Ort ist ein Pflichtfeld" {
		t.Fatalf("field message: %q", got)
	}
}

func TestParseAndValidateSparse_SelectParseNotAllowed_Translated(t *testing.T) {
	art := field.NewSelectField("art", "ART", false, []field.SelectOption{{Value: int32(1), Label: "A"}})
	fg := groupWithCtx(t, deCtx(deTexts), art)

	_, err := fg.ParseAndValidateSparse(map[string]interface{}{"art": float64(9)})

	if got := fieldErrors(t, err)["art"]; got != "Ungültige Auswahl" {
		t.Fatalf("field message: %q", got)
	}
}

func TestTranslateWithoutValidationKeys_KeepsEnglishMessage_LabelsBanner(t *testing.T) {
	ort := field.NewTextFieldWithLength("ort", "ORT", false, "", 0, 5)
	fg := groupWithCtx(t, deCtx(map[string]string{"ORT": "Ort"}), ort)

	_, err := fg.ParseAndValidate(map[string]interface{}{"ort": "Wien-Floridsdorf"})

	const english = "text field ort must be at most 5 characters"
	if got := fieldErrors(t, err)["ort"]; got != english {
		t.Fatalf("field message: %q", got)
	}
	if got := err.Error(); got != "Ort: "+english {
		t.Fatalf("banner: %q", got)
	}
}

func TestTranslateReturnsEmpty_FallsBackToEnglishAndID(t *testing.T) {
	ort := field.NewTextFieldWithLength("ort", "ORT", false, "", 0, 5)
	ctx := &core.UiContext{Translate: func(string) string { return "" }}
	fg := groupWithCtx(t, ctx, ort)

	_, err := fg.ParseAndValidate(map[string]interface{}{"ort": "Wien-Floridsdorf"})

	if got := err.Error(); got != "ort: text field ort must be at most 5 characters" {
		t.Fatalf("banner: %q", got)
	}
}

func TestNoContext_UnchangedFieldErrors(t *testing.T) {
	ort := field.NewTextFieldWithLength("ort", "ORT", false, "", 0, 5)
	fg := NewFormGroup([]field.FormField{ort})

	_, err := fg.ParseAndValidate(map[string]interface{}{"ort": "Wien-Floridsdorf"})

	if _, ok := err.(FieldErrors); !ok {
		t.Fatalf("without context the error stays FieldErrors, got %T", err)
	}
	if got := err.Error(); got != "ort: text field ort must be at most 5 characters" {
		t.Fatalf("banner: %q", got)
	}
}

func TestWrappedFieldErrors_ExposesMap(t *testing.T) {
	ort := field.NewTextField("ort", "ORT", true, "")
	fg := groupWithCtx(t, deCtx(deTexts), ort)

	_, err := fg.ParseValues(map[string]interface{}{})

	var fe interface{ FieldErrors() map[string]string }
	if !errors.As(err, &fe) || fe.FieldErrors()["ort"] != "Ort ist ein Pflichtfeld" {
		t.Fatalf("FieldErrors() via interface: %T %v", err, err)
	}
}
