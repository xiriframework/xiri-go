package dialog

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
)

func TestNewDialogDelete(t *testing.T) {
	testUrl := url.NewUrl("/api/delete")

	t.Run("Standard with translation", func(t *testing.T) {
		mockTranslator := func(key string) string {
			translations := map[string]string{
				"Delete": "Löschen",
				"Ok":     "Bestätigen",
				"Back":   "Zurück",
			}
			return translations[key]
		}

		dialog := NewDialogDelete(
			"Delete this item?",
			testUrl,
			nil,
			nil,
			nil,
			nil,
			mockTranslator,
		)

		result := dialog.Print(nil)

		if result["header"] != "Löschen" {
			t.Errorf("header = %v, want Löschen", result["header"])
		}
		if result["type"] != string(core.DialogTypeQuestion) {
			t.Errorf("type = %v, want %v", result["type"], core.DialogTypeQuestion)
		}

		contentMap, ok := result["content"].(map[string]any)
		if !ok {
			t.Fatal("content is not a map")
		}
		if contentMap["icon"] != "warning" {
			t.Errorf("content.icon = %v, want warning", contentMap["icon"])
		}
		if contentMap["question"] != "Delete this item?" {
			t.Errorf("content.question = %v, want 'Delete this item?'", contentMap["question"])
		}

		buttons, ok := result["buttons"].([]map[string]any)
		if !ok {
			t.Fatal("buttons is not an array")
		}
		if len(buttons) != 2 {
			t.Fatalf("buttons length = %d, want 2", len(buttons))
		}

		// Check close button
		if buttons[0]["text"] != "Zurück" {
			t.Errorf("buttons[0].text = %v, want Zurück", buttons[0]["text"])
		}
		if buttons[0]["action"] != string(core.ButtonActionClose) {
			t.Errorf("buttons[0].action = %v, want close", buttons[0]["action"])
		}

		// Check form button
		if buttons[1]["text"] != "Bestätigen" {
			t.Errorf("buttons[1].text = %v, want Bestätigen", buttons[1]["text"])
		}
		if buttons[1]["action"] != string(core.ButtonActionForm) {
			t.Errorf("buttons[1].action = %v, want form", buttons[1]["action"])
		}
		if buttons[1]["default"] != true {
			t.Errorf("buttons[1].default = %v, want true", buttons[1]["default"])
		}
	})

	t.Run("With custom texts", func(t *testing.T) {
		customHeader := "Remove Item"
		customOk := "Confirm Delete"
		customClose := "Cancel"

		dialog := NewDialogDelete(
			"Are you sure?",
			testUrl,
			nil,
			&customHeader,
			&customOk,
			&customClose,
			nil,
		)

		result := dialog.Print(nil)

		if result["header"] != "Remove Item" {
			t.Errorf("header = %v, want Remove Item", result["header"])
		}

		buttons, ok := result["buttons"].([]map[string]any)
		if !ok {
			t.Fatal("buttons is not an array")
		}

		if buttons[0]["text"] != "Cancel" {
			t.Errorf("buttons[0].text = %v, want Cancel", buttons[0]["text"])
		}
		if buttons[1]["text"] != "Confirm Delete" {
			t.Errorf("buttons[1].text = %v, want Confirm Delete", buttons[1]["text"])
		}
	})

	t.Run("With extra data", func(t *testing.T) {
		extra := map[string]any{
			"itemId": int64(123),
			"type":   "user",
		}

		dialog := NewDialogDelete(
			"Delete user?",
			testUrl,
			extra,
			nil,
			nil,
			nil,
			nil,
		)

		result := dialog.Print(nil)

		extraMap, ok := result["extra"].(map[string]any)
		if !ok {
			t.Fatal("extra is not a map")
		}
		if extraMap["itemId"] != int64(123) {
			t.Errorf("extra.itemId = %v, want 123", extraMap["itemId"])
		}
		if extraMap["type"] != "user" {
			t.Errorf("extra.type = %v, want user", extraMap["type"])
		}
	})

	t.Run("Defaults without translator", func(t *testing.T) {
		dialog := NewDialogDelete(
			"Delete?",
			testUrl,
			nil,
			nil,
			nil,
			nil,
			nil,
		)

		result := dialog.Print(nil)

		if result["header"] != "Delete" {
			t.Errorf("header = %v, want Delete", result["header"])
		}

		buttons, ok := result["buttons"].([]map[string]any)
		if !ok {
			t.Fatal("buttons is not an array")
		}

		if buttons[0]["text"] != "Back" {
			t.Errorf("buttons[0].text = %v, want Back", buttons[0]["text"])
		}
		if buttons[1]["text"] != "Ok" {
			t.Errorf("buttons[1].text = %v, want Ok", buttons[1]["text"])
		}
	})
}

func TestNewDialogWarning(t *testing.T) {
	testUrl := url.NewUrl("/api/confirm")

	t.Run("Icon difference from delete", func(t *testing.T) {
		dialog := NewDialogWarning(
			"Are you sure?",
			testUrl,
			nil,
			nil,
			nil,
			nil,
			nil,
		)

		result := dialog.Print(nil)

		contentMap, ok := result["content"].(map[string]any)
		if !ok {
			t.Fatal("content is not a map")
		}

		// Warning uses help_outline, Delete uses warning
		if contentMap["icon"] != "help_outline" {
			t.Errorf("content.icon = %v, want help_outline", contentMap["icon"])
		}

		if result["type"] != string(core.DialogTypeQuestion) {
			t.Errorf("type = %v, want question", result["type"])
		}
	})

	t.Run("With translation", func(t *testing.T) {
		mockTranslator := func(key string) string {
			if key == "Warning" {
				return "Achtung"
			}
			return key
		}

		dialog := NewDialogWarning(
			"This action cannot be undone",
			testUrl,
			nil,
			nil,
			nil,
			nil,
			mockTranslator,
		)

		result := dialog.Print(nil)

		if result["header"] != "Achtung" {
			t.Errorf("header = %v, want Achtung", result["header"])
		}
	})

	t.Run("Custom header", func(t *testing.T) {
		customHeader := "Important Notice"

		dialog := NewDialogWarning(
			"Please review before continuing",
			testUrl,
			nil,
			&customHeader,
			nil,
			nil,
			nil,
		)

		result := dialog.Print(nil)

		if result["header"] != "Important Notice" {
			t.Errorf("header = %v, want Important Notice", result["header"])
		}
	})
}
