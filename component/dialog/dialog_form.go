package dialog

import (
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
)

// NewDialogForm creates a form dialog
//
// Parameters:
//   - fields: Array of form field configurations (from form/group.ExportForFrontendWithValues)
//   - u: URL to submit form data
//   - header: Dialog header text (required, passed as pointer for consistency with old API)
//   - extra: Optional extra data passed to frontend
//   - okText: Optional custom OK button text (uses translation or "Ok" as fallback)
//   - closeText: Optional custom close button text (uses translation or "Back" as fallback)
//   - translator: Translation function for i18n
//
// Creates a form-type dialog with standard Close + Submit buttons.
// The fields are rendered by Angular's form component.
func NewDialogForm(
	fields []map[string]any,
	u *url.Url,
	header *string,
	extra map[string]any,
	okText *string,
	closeText *string,
	translator core.TranslateFunc,
) Dialog {
	ok := resolveText(okText, "Ok", "Ok", translator)
	closeBtn := resolveText(closeText, "Back", "Back", translator)

	buttons := buildStandardButtons(closeBtn, ok, u)

	return newDialog(
		core.DialogTypeForm,
		*header,
		fields,
		buttons,
		extra,
		nil,
	)
}

// NewDialogFormMultiDelete creates a delete confirmation dialog for multi-selected items
//
// This is a convenience wrapper around NewDialogDelete that automatically:
//   - Sets extra["data"] to the selected IDs
//   - Sets extra["done"] to true
//
// Parameters:
//   - u: URL to call for deletion
//   - selectedIDs: Array of item IDs to delete
//   - content: Question text to display
//   - headerText: Optional custom header (uses translation or "Delete" as fallback)
//   - okText: Optional custom OK button text (uses translation or "Ok" as fallback)
//   - closeText: Optional custom close button text (uses translation or "Back" as fallback)
//   - translator: Translation function for i18n
func NewDialogFormMultiDelete(
	u *url.Url,
	selectedIDs []int64,
	content string,
	headerText *string,
	okText *string,
	closeText *string,
	translator core.TranslateFunc,
) Dialog {
	extra := map[string]interface{}{
		"data": selectedIDs,
		"done": true,
	}

	return NewDialogDelete(
		content,
		u,
		extra,
		headerText,
		okText,
		closeText,
		translator,
	)
}

// NewDialogFormMultiEdit creates a multi-edit form dialog with selected IDs
//
// This is a convenience wrapper around NewDialogForm that automatically:
//   - Sets extra["data"] to the selected IDs
//   - Sets extra["done"] to true
//
// Parameters:
//   - fields: Array of form field configurations
//   - u: URL to submit form data
//   - selectedIDs: Array of item IDs being edited
//   - header: Dialog header text
//   - okText: OK button text
//   - closeText: Close button text
//   - translator: Translation function for i18n
//
// Note: Unlike NewDialogForm, this takes header/okText/closeText as strings
// (not pointers) for convenience.
func NewDialogFormMultiEdit(
	fields []map[string]any,
	u *url.Url,
	selectedIDs []int64,
	header string,
	okText string,
	closeText string,
	translator core.TranslateFunc,
) Dialog {
	extra := map[string]interface{}{
		"data": selectedIDs,
		"done": true,
	}

	return NewDialogForm(
		fields,
		u,
		&header,
		extra,
		&okText,
		&closeText,
		translator,
	)
}
