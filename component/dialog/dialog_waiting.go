package dialog

import (
	"github.com/xiriframework/xiri-go/component/button"
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
)

// Waiting dialog state constants
const (
	// WaitingStateInitial represents the initial waiting dialog shown to the user
	// This displays a dialog with icon, text, and close button while polling begins
	WaitingStateInitial = 0

	// WaitingStateNotDone is returned during polling to indicate the operation is still in progress
	// Frontend continues polling when receiving this state
	WaitingStateNotDone = 1

	// WaitingStateDone is returned when the operation completes
	// Contains URL for redirect and optional blocked element identifier
	WaitingStateDone = 2

	// WaitingStateError is returned when the operation fails
	// Contains error message, stops polling and allows user to close dialog
	WaitingStateError = 3
)

// DialogWaiting represents a waiting dialog with polling capability
//
// State machine:
//  1. Initial state (WaitingStateInitial): Shows dialog with polling configuration
//  2. Polling states (WaitingStateNotDone): Returns {"done": false} to continue polling
//  3. Completion state (WaitingStateDone): Returns {"done": true, "url": ..., "blocked": ...}
//  4. Error state (WaitingStateError): Returns {"done": true, "error": ...} to stop polling with error
type DialogWaiting struct {
	*dialogImpl
	waitingType    int
	waitingUrl     string
	waitingBlocked string
	waitingError   string
}

// NewDialogWaiting creates a waiting dialog with polling
func NewDialogWaiting(
	text string,
	u *url.Url,
	header string,
	checkTime int,
	extra map[string]any,
	closeText *string,
) Dialog {
	closeBtn := resolveText(closeText, "Back")

	content := DialogWaitingContent{
		Icon: "help_outline",
		Text: text,
	}

	buttons := []*button.Button{
		button.NewCloseButton(closeBtn, core.ColorPrimary, core.ButtonTypeStroked, "", false, nil, false, nil),
	}

	options := map[string]any{
		"checkTime": checkTime,
		"url":       u.PrintPrefix(),
	}

	baseDialog := newDialog(
		core.DialogTypeWaiting,
		header,
		content,
		buttons,
		extra,
		options,
	)

	return &DialogWaiting{
		dialogImpl:     baseDialog,
		waitingType:    WaitingStateInitial,
		waitingUrl:     "",
		waitingBlocked: "",
	}
}

// NewDialogWaitingNotDone creates a "not done" polling response
//
// Used during polling to indicate the operation is still in progress.
// Frontend will continue polling when receiving this response.
func NewDialogWaitingNotDone() Dialog {
	return &DialogWaiting{
		dialogImpl:  newDialog(core.DialogTypeWaiting, "Not Done", nil, []*button.Button{}, nil, nil),
		waitingType: WaitingStateNotDone,
	}
}

// NewDialogWaitingDone creates a "done" polling response with redirect
//
// Parameters:
//   - u: URL to navigate to after completion
//   - blocked: Optional identifier for blocked/disabled UI element
//
// Frontend stops polling and navigates to the provided URL.
func NewDialogWaitingDone(u string, blocked string) Dialog {
	return &DialogWaiting{
		dialogImpl:     newDialog(core.DialogTypeWaiting, "Done", nil, []*button.Button{}, nil, nil),
		waitingType:    WaitingStateDone,
		waitingUrl:     u,
		waitingBlocked: blocked,
	}
}

// NewDialogWaitingError creates an error polling response
//
// Parameters:
//   - message: Error message to display to the user
//
// Frontend stops polling, shows the error, and allows the user to close the dialog.
func NewDialogWaitingError(message string) Dialog {
	return &DialogWaiting{
		dialogImpl:   newDialog(core.DialogTypeWaiting, "Error", nil, []*button.Button{}, nil, nil),
		waitingType:  WaitingStateError,
		waitingError: message,
	}
}

// Print returns the JSON representation based on waiting type
//
// Returns different JSON structures depending on the waiting state:
//   - WaitingStateInitial: Full dialog structure with polling config
//   - WaitingStateNotDone: {"done": false} to continue polling
//   - WaitingStateDone: {"done": true, "url": ..., "blocked": ...} to complete
func (dw *DialogWaiting) Print(ctx *core.UiContext) map[string]any {
	switch dw.waitingType {
	case WaitingStateInitial:
		return dw.dialogImpl.Print(ctx)
	case WaitingStateNotDone:
		return map[string]any{
			"done": false,
		}
	case WaitingStateDone:
		return map[string]any{
			"done":    true,
			"url":     dw.waitingUrl,
			"blocked": dw.waitingBlocked,
		}
	case WaitingStateError:
		return map[string]any{
			"done":  true,
			"error": dw.waitingError,
		}
	default:
		return dw.dialogImpl.Print(ctx)
	}
}
