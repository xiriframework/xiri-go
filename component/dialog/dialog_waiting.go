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
)

// DialogWaiting represents a waiting dialog with polling capability
//
// State machine:
//  1. Initial state (WaitingStateInitial): Shows dialog with polling configuration
//  2. Polling states (WaitingStateNotDone): Returns {"done": false} to continue polling
//  3. Completion state (WaitingStateDone): Returns {"done": true, "url": ..., "blocked": ...}
type DialogWaiting struct {
	*dialogImpl
	waitingType    int
	waitingUrl     string
	waitingBlocked string
}

// NewDialogWaiting creates a waiting dialog with polling
func NewDialogWaiting(
	text string,
	u *url.Url,
	header string,
	checkTime int,
	extra map[string]any,
	closeText *string,
	translator core.TranslateFunc,
) Dialog {
	closeBtn := "Back"
	if closeText != nil {
		closeBtn = *closeText
	} else if translator != nil {
		closeBtn = translator("Back")
	}

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

// Print returns the JSON representation based on waiting type
//
// Returns different JSON structures depending on the waiting state:
//   - WaitingStateInitial: Full dialog structure with polling config
//   - WaitingStateNotDone: {"done": false} to continue polling
//   - WaitingStateDone: {"done": true, "url": ..., "blocked": ...} to complete
func (dw *DialogWaiting) Print(translator core.TranslateFunc) map[string]any {
	switch dw.waitingType {
	case WaitingStateInitial:
		return dw.dialogImpl.Print(translator)
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
	default:
		return dw.dialogImpl.Print(translator)
	}
}
