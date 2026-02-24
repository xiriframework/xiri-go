package dialog

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
)

func TestNewDialogForm(t *testing.T) {
	testUrl := url.NewUrl("/api/save")
	header := "Edit User"

	t.Run("Standard form dialog", func(t *testing.T) {
		fields := []map[string]any{
			{"field": "name", "label": "Name", "type": "text"},
			{"field": "email", "label": "Email", "type": "text"},
		}

		dialog := NewDialogForm(
			fields,
			testUrl,
			&header,
			nil,
			nil,
			nil,
			nil,
		)

		result := dialog.Print(nil)

		if result["header"] != "Edit User" {
			t.Errorf("header = %v, want Edit User", result["header"])
		}
		if result["type"] != string(core.DialogTypeForm) {
			t.Errorf("type = %v, want form", result["type"])
		}

		contentSlice, ok := result["content"].([]map[string]any)
		if !ok {
			t.Fatal("content is not a slice")
		}
		if len(contentSlice) != 2 {
			t.Fatalf("content length = %d, want 2", len(contentSlice))
		}
		if contentSlice[0]["field"] != "name" {
			t.Errorf("content[0].field = %v, want name", contentSlice[0]["field"])
		}

		buttons, ok := result["buttons"].([]map[string]any)
		if !ok {
			t.Fatal("buttons is not an array")
		}
		if len(buttons) != 2 {
			t.Fatalf("buttons length = %d, want 2", len(buttons))
		}
	})

	t.Run("With translation", func(t *testing.T) {
		mockTranslator := func(key string) string {
			translations := map[string]string{
				"Ok":   "Speichern",
				"Back": "Abbrechen",
			}
			return translations[key]
		}

		fields := []map[string]any{}
		dialog := NewDialogForm(
			fields,
			testUrl,
			&header,
			nil,
			nil,
			nil,
			mockTranslator,
		)

		result := dialog.Print(nil)

		buttons, ok := result["buttons"].([]map[string]any)
		if !ok {
			t.Fatal("buttons is not an array")
		}

		if buttons[0]["text"] != "Abbrechen" {
			t.Errorf("buttons[0].text = %v, want Abbrechen", buttons[0]["text"])
		}
		if buttons[1]["text"] != "Speichern" {
			t.Errorf("buttons[1].text = %v, want Speichern", buttons[1]["text"])
		}
	})

	t.Run("With custom button texts", func(t *testing.T) {
		okText := "Submit"
		closeText := "Discard"

		fields := []map[string]any{}
		dialog := NewDialogForm(
			fields,
			testUrl,
			&header,
			nil,
			&okText,
			&closeText,
			nil,
		)

		result := dialog.Print(nil)

		buttons, ok := result["buttons"].([]map[string]any)
		if !ok {
			t.Fatal("buttons is not an array")
		}

		if buttons[0]["text"] != "Discard" {
			t.Errorf("buttons[0].text = %v, want Discard", buttons[0]["text"])
		}
		if buttons[1]["text"] != "Submit" {
			t.Errorf("buttons[1].text = %v, want Submit", buttons[1]["text"])
		}
	})

	t.Run("With extra data", func(t *testing.T) {
		extra := map[string]any{
			"userId": int64(456),
		}

		fields := []map[string]any{}
		dialog := NewDialogForm(
			fields,
			testUrl,
			&header,
			extra,
			nil,
			nil,
			nil,
		)

		result := dialog.Print(nil)

		extraMap, ok := result["extra"].(map[string]any)
		if !ok {
			t.Fatal("extra is not a map")
		}
		if extraMap["userId"] != int64(456) {
			t.Errorf("extra.userId = %v, want 456", extraMap["userId"])
		}
	})
}

func TestNewDialogFormMultiDelete(t *testing.T) {
	testUrl := url.NewUrl("/api/multi-delete")
	selectedIDs := []int64{10, 20, 30}

	t.Run("Standard multi-delete", func(t *testing.T) {
		dialog := NewDialogFormMultiDelete(
			testUrl,
			selectedIDs,
			"Delete 3 items?",
			nil,
			nil,
			nil,
			nil,
		)

		result := dialog.Print(nil)

		// Should be a question-type dialog (via NewDialogDelete)
		if result["type"] != string(core.DialogTypeQuestion) {
			t.Errorf("type = %v, want question", result["type"])
		}

		// Check extra data
		extraMap, ok := result["extra"].(map[string]any)
		if !ok {
			t.Fatal("extra is not a map")
		}

		dataSlice, ok := extraMap["data"].([]int64)
		if !ok {
			t.Fatal("extra.data is not []int64")
		}
		if len(dataSlice) != 3 || dataSlice[0] != 10 || dataSlice[1] != 20 || dataSlice[2] != 30 {
			t.Errorf("extra.data = %v, want [10, 20, 30]", dataSlice)
		}

		if extraMap["done"] != true {
			t.Errorf("extra.done = %v, want true", extraMap["done"])
		}

		// Check content
		contentMap, ok := result["content"].(map[string]any)
		if !ok {
			t.Fatal("content is not a map")
		}
		if contentMap["question"] != "Delete 3 items?" {
			t.Errorf("content.question = %v, want 'Delete 3 items?'", contentMap["question"])
		}
	})

	t.Run("With custom texts", func(t *testing.T) {
		header := "Bulk Delete"
		okText := "Delete All"
		closeText := "Keep"

		dialog := NewDialogFormMultiDelete(
			testUrl,
			selectedIDs,
			"Remove selected?",
			&header,
			&okText,
			&closeText,
			nil,
		)

		result := dialog.Print(nil)

		if result["header"] != "Bulk Delete" {
			t.Errorf("header = %v, want Bulk Delete", result["header"])
		}

		buttons, ok := result["buttons"].([]map[string]any)
		if !ok {
			t.Fatal("buttons is not an array")
		}

		if buttons[1]["text"] != "Delete All" {
			t.Errorf("buttons[1].text = %v, want Delete All", buttons[1]["text"])
		}
	})
}

func TestNewDialogFormMultiEdit(t *testing.T) {
	testUrl := url.NewUrl("/api/multi-edit")
	selectedIDs := []int64{5, 15, 25}

	t.Run("Standard multi-edit", func(t *testing.T) {
		fields := []map[string]any{
			{"field": "status", "label": "Status", "type": "select"},
		}

		dialog := NewDialogFormMultiEdit(
			fields,
			testUrl,
			selectedIDs,
			"Edit 3 Users",
			"Save Changes",
			"Cancel",
			nil,
		)

		result := dialog.Print(nil)

		// Should be a form-type dialog
		if result["type"] != string(core.DialogTypeForm) {
			t.Errorf("type = %v, want form", result["type"])
		}

		if result["header"] != "Edit 3 Users" {
			t.Errorf("header = %v, want Edit 3 Users", result["header"])
		}

		// Check extra data
		extraMap, ok := result["extra"].(map[string]any)
		if !ok {
			t.Fatal("extra is not a map")
		}

		dataSlice, ok := extraMap["data"].([]int64)
		if !ok {
			t.Fatal("extra.data is not []int64")
		}
		if len(dataSlice) != 3 || dataSlice[0] != 5 || dataSlice[1] != 15 || dataSlice[2] != 25 {
			t.Errorf("extra.data = %v, want [5, 15, 25]", dataSlice)
		}

		if extraMap["done"] != true {
			t.Errorf("extra.done = %v, want true", extraMap["done"])
		}

		// Check buttons
		buttons, ok := result["buttons"].([]map[string]any)
		if !ok {
			t.Fatal("buttons is not an array")
		}

		if buttons[0]["text"] != "Cancel" {
			t.Errorf("buttons[0].text = %v, want Cancel", buttons[0]["text"])
		}
		if buttons[1]["text"] != "Save Changes" {
			t.Errorf("buttons[1].text = %v, want Save Changes", buttons[1]["text"])
		}

		// Check fields
		contentSlice, ok := result["content"].([]map[string]any)
		if !ok {
			t.Fatal("content is not a slice")
		}
		if len(contentSlice) != 1 {
			t.Fatalf("content length = %d, want 1", len(contentSlice))
		}
		if contentSlice[0]["field"] != "status" {
			t.Errorf("content[0].field = %v, want status", contentSlice[0]["field"])
		}
	})
}
