package dialog

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
)

func TestNewDialogWaiting(t *testing.T) {
	testUrl := url.NewUrl("/api/status")

	t.Run("Standard waiting dialog", func(t *testing.T) {
		dialog := NewDialogWaiting(
			"Processing your request...",
			testUrl,
			"Please Wait",
			3000,
			nil,
			nil,
			nil,
		)

		result := dialog.Print(nil)

		if result["header"] != "Please Wait" {
			t.Errorf("header = %v, want Please Wait", result["header"])
		}
		if result["type"] != string(core.DialogTypeWaiting) {
			t.Errorf("type = %v, want waiting", result["type"])
		}

		// Check options
		if result["checkTime"] != 3000 {
			t.Errorf("checkTime = %v, want 3000", result["checkTime"])
		}
		if result["url"] != testUrl.PrintPrefix() {
			t.Errorf("url = %v, want %v", result["url"], testUrl.PrintPrefix())
		}

		// Check content
		contentMap, ok := result["content"].(map[string]any)
		if !ok {
			t.Fatal("content is not a map")
		}
		if contentMap["icon"] != "help_outline" {
			t.Errorf("content.icon = %v, want help_outline", contentMap["icon"])
		}
		if contentMap["text"] != "Processing your request..." {
			t.Errorf("content.text = %v, want 'Processing your request...'", contentMap["text"])
		}

		// Check buttons
		buttons, ok := result["buttons"].([]map[string]any)
		if !ok {
			t.Fatal("buttons is not an array")
		}
		if len(buttons) != 1 {
			t.Fatalf("buttons length = %d, want 1", len(buttons))
		}
		if buttons[0]["action"] != string(core.ButtonActionClose) {
			t.Errorf("buttons[0].action = %v, want close", buttons[0]["action"])
		}
	})

	t.Run("With translation", func(t *testing.T) {
		mockTranslator := func(key string) string {
			if key == "Back" {
				return "Zurück"
			}
			return key
		}

		dialog := NewDialogWaiting(
			"Loading...",
			testUrl,
			"Waiting",
			5000,
			nil,
			nil,
			mockTranslator,
		)

		result := dialog.Print(nil)

		buttons, ok := result["buttons"].([]map[string]any)
		if !ok {
			t.Fatal("buttons is not an array")
		}

		if buttons[0]["text"] != "Zurück" {
			t.Errorf("buttons[0].text = %v, want Zurück", buttons[0]["text"])
		}
	})

	t.Run("With custom close text", func(t *testing.T) {
		closeText := "Cancel"

		dialog := NewDialogWaiting(
			"Processing...",
			testUrl,
			"Wait",
			2000,
			nil,
			&closeText,
			nil,
		)

		result := dialog.Print(nil)

		buttons, ok := result["buttons"].([]map[string]any)
		if !ok {
			t.Fatal("buttons is not an array")
		}

		if buttons[0]["text"] != "Cancel" {
			t.Errorf("buttons[0].text = %v, want Cancel", buttons[0]["text"])
		}
	})
}

func TestNewDialogWaitingNotDone(t *testing.T) {
	dialog := NewDialogWaitingNotDone()

	result := dialog.Print(nil)

	// Should only have "done: false"
	if len(result) != 1 {
		t.Errorf("result length = %d, want 1", len(result))
	}

	done, ok := result["done"]
	if !ok {
		t.Fatal("done field not present")
	}

	if done != false {
		t.Errorf("done = %v, want false", done)
	}
}

func TestNewDialogWaitingDone(t *testing.T) {
	t.Run("With URL and blocked", func(t *testing.T) {
		dialog := NewDialogWaitingDone("/success", "submit-button")

		result := dialog.Print(nil)

		if len(result) != 3 {
			t.Errorf("result length = %d, want 3", len(result))
		}

		if result["done"] != true {
			t.Errorf("done = %v, want true", result["done"])
		}
		if result["url"] != "/success" {
			t.Errorf("url = %v, want /success", result["url"])
		}
		if result["blocked"] != "submit-button" {
			t.Errorf("blocked = %v, want submit-button", result["blocked"])
		}
	})

	t.Run("With empty blocked", func(t *testing.T) {
		dialog := NewDialogWaitingDone("/complete", "")

		result := dialog.Print(nil)

		if result["blocked"] != "" {
			t.Errorf("blocked = %v, want empty string", result["blocked"])
		}
	})
}

func TestDialogWaiting_Print_States(t *testing.T) {
	testUrl := url.NewUrl("/api/check")

	t.Run("State machine transitions", func(t *testing.T) {
		// State 0: Initial dialog
		initialDialog := NewDialogWaiting(
			"Working...",
			testUrl,
			"Processing",
			1000,
			nil,
			nil,
			nil,
		)
		initialResult := initialDialog.Print(nil)

		if initialResult["type"] != string(core.DialogTypeWaiting) {
			t.Errorf("Initial state type = %v, want waiting", initialResult["type"])
		}
		if initialResult["header"] == nil {
			t.Error("Initial state should have full dialog structure")
		}

		// State 1: Not done (polling)
		notDoneDialog := NewDialogWaitingNotDone()
		notDoneResult := notDoneDialog.Print(nil)

		if notDoneResult["done"] != false {
			t.Errorf("NotDone state done = %v, want false", notDoneResult["done"])
		}
		if _, hasHeader := notDoneResult["header"]; hasHeader {
			t.Error("NotDone state should not have header (minimal response)")
		}

		// State 2: Done (completion)
		doneDialog := NewDialogWaitingDone("/result", "")
		doneResult := doneDialog.Print(nil)

		if doneResult["done"] != true {
			t.Errorf("Done state done = %v, want true", doneResult["done"])
		}
		if doneResult["url"] != "/result" {
			t.Errorf("Done state url = %v, want /result", doneResult["url"])
		}
		if _, hasHeader := doneResult["header"]; hasHeader {
			t.Error("Done state should not have header (minimal response)")
		}
	})
}

func TestWaitingStateConstants(t *testing.T) {
	// Verify constant values match documentation and implementation
	if WaitingStateInitial != 0 {
		t.Errorf("WaitingStateInitial = %d, want 0", WaitingStateInitial)
	}
	if WaitingStateNotDone != 1 {
		t.Errorf("WaitingStateNotDone = %d, want 1", WaitingStateNotDone)
	}
	if WaitingStateDone != 2 {
		t.Errorf("WaitingStateDone = %d, want 2", WaitingStateDone)
	}
}
