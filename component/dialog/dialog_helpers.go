package dialog

import (
	"github.com/xiriframework/xiri-go/component/button"
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
)

// resolveText resolves text with priority: custom > translation > default
//
// Priority order:
//  1. customText (if provided and non-nil)
//  2. translator(translationKey) (if translator is provided)
//  3. defaultText (fallback)
func resolveText(
	customText *string,
	translationKey string,
	defaultText string,
	translator core.TranslateFunc,
) string {
	if customText != nil {
		return *customText
	}
	if translator != nil {
		return translator(translationKey)
	}
	return defaultText
}

// dialogTexts holds resolved text values for standard dialog elements
type dialogTexts struct {
	Header string
	Ok     string
	Close  string
}

// resolveDialogTexts resolves all standard dialog texts using the priority system
//
// For each text field, the priority is: custom > translation > default
func resolveDialogTexts(
	headerText *string, headerKey, headerDefault string,
	okText *string, okKey, okDefault string,
	closeText *string, closeKey, closeDefault string,
	translator core.TranslateFunc,
) dialogTexts {
	return dialogTexts{
		Header: resolveText(headerText, headerKey, headerDefault, translator),
		Ok:     resolveText(okText, okKey, okDefault, translator),
		Close:  resolveText(closeText, closeKey, closeDefault, translator),
	}
}

// buildStandardButtons creates the standard [Close, Submit] button pair
//
// This is used in question and form dialogs to provide consistent button layout:
//   - Close button (stroked, left side)
//   - Submit button (raised, right side, default/autofocus)
func buildStandardButtons(closeText, okText string, u *url.Url) []*button.Button {
	return []*button.Button{
		button.NewCloseButton(closeText, core.ColorPrimary, core.ButtonTypeStroked, "", false, nil, false, nil),
		button.NewFormButton(okText, u, core.ColorPrimary, core.ButtonTypeRaised, "", false, nil, true, nil),
	}
}
