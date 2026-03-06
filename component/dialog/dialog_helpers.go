package dialog

import (
	"github.com/xiriframework/xiri-go/component/button"
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
)

// resolveText returns the custom text if provided, otherwise the default key.
// Translation happens lazily at Print time via core.Translate.
func resolveText(customText *string, defaultKey string) string {
	if customText != nil {
		return *customText
	}
	return defaultKey
}

// dialogTexts holds resolved text values for standard dialog elements
type dialogTexts struct {
	Header string
	Ok     string
	Close  string
}

// resolveDialogTexts returns custom texts or default keys for standard dialog elements.
// Translation happens lazily at Print time via core.Translate.
func resolveDialogTexts(
	headerText *string, headerDefault string,
	okText *string, okDefault string,
	closeText *string, closeDefault string,
) dialogTexts {
	return dialogTexts{
		Header: resolveText(headerText, headerDefault),
		Ok:     resolveText(okText, okDefault),
		Close:  resolveText(closeText, closeDefault),
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
