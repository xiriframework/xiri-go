package dialog

import (
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
)

// NewDialogDelete creates a delete confirmation dialog
//
// Parameters:
//   - text: Question text to display ("Delete this item?")
//   - u: URL to call when user confirms deletion
//   - extra: Optional extra data (e.g., item IDs for multi-delete)
//   - headerText: Optional custom header (uses translation or "Delete" as fallback)
//   - okText: Optional custom OK button text (uses translation or "Ok" as fallback)
//   - closeText: Optional custom close button text (uses translation or "Back" as fallback)
//   - translator: Translation function for i18n
//
// The dialog shows a warning icon and creates a question-type dialog with
// standard Close + Submit buttons.
func NewDialogDelete(
	text string,
	u *url.Url,
	extra map[string]any,
	headerText *string,
	okText *string,
	closeText *string,
	translator core.TranslateFunc,
) Dialog {
	texts := resolveDialogTexts(
		headerText, "Delete", "Delete",
		okText, "Ok", "Ok",
		closeText, "Back", "Back",
		translator,
	)

	content := DialogQuestionContent{
		Icon:     "warning",
		Question: text,
	}

	buttons := buildStandardButtons(texts.Close, texts.Ok, u)

	return newDialog(
		core.DialogTypeQuestion,
		texts.Header,
		content,
		buttons,
		extra,
		nil,
	)
}

// NewDialogWarning creates a warning confirmation dialog
//
// Parameters:
//   - text: Warning text to display
//   - u: URL to call when user confirms
//   - extra: Optional extra data
//   - headerText: Optional custom header (uses translation or "Warning" as fallback)
//   - okText: Optional custom OK button text (uses translation or "Ok" as fallback)
//   - closeText: Optional custom close button text (uses translation or "Back" as fallback)
//   - translator: Translation function for i18n
//
// Similar to NewDialogDelete but uses "help_outline" icon instead of "warning".
// Creates a question-type dialog with standard Close + Submit buttons.
func NewDialogWarning(
	text string,
	u *url.Url,
	extra map[string]any,
	headerText *string,
	okText *string,
	closeText *string,
	translator core.TranslateFunc,
) Dialog {
	texts := resolveDialogTexts(
		headerText, "Warning", "Warning",
		okText, "Ok", "Ok",
		closeText, "Back", "Back",
		translator,
	)

	content := DialogQuestionContent{
		Icon:     "help_outline",
		Question: text,
	}

	buttons := buildStandardButtons(texts.Close, texts.Ok, u)

	return newDialog(
		core.DialogTypeQuestion,
		texts.Header,
		content,
		buttons,
		extra,
		nil,
	)
}
