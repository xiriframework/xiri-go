package dialog

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/button"
	"github.com/xiriframework/xiri-go/component/core"
)

func TestNewDialog(t *testing.T) {
	header := "Test Dialog"
	content := "Test Content"
	buttons := []*button.Button{
		button.NewCloseButton("Close", core.ColorPrimary, core.ButtonTypeStroked, "", false, nil, false, nil),
	}
	extra := map[string]any{"key": "value"}
	options := map[string]any{"size": "lg"}

	dialog := NewDialog(core.DialogTypeQuestion, header, content, buttons, extra, options)

	if dialog == nil {
		t.Fatal("NewDialog() returned nil")
	}

	// Verify Print output
	result := dialog.Print(nil)

	if result["header"] != header {
		t.Errorf("header = %v, want %v", result["header"], header)
	}
	if result["type"] != string(core.DialogTypeQuestion) {
		t.Errorf("type = %v, want %v", result["type"], core.DialogTypeQuestion)
	}
	if result["content"] != content {
		t.Errorf("content = %v, want %v", result["content"], content)
	}

	extraMap, ok := result["extra"].(map[string]any)
	if !ok {
		t.Fatal("extra is not a map")
	}
	if extraMap["key"] != "value" {
		t.Errorf("extra[key] = %v, want value", extraMap["key"])
	}

	if result["size"] != "lg" {
		t.Errorf("size option = %v, want lg", result["size"])
	}

	buttonsArr, ok := result["buttons"].([]map[string]any)
	if !ok {
		t.Fatal("buttons is not an array")
	}
	if len(buttonsArr) != 1 {
		t.Errorf("buttons length = %d, want 1", len(buttonsArr))
	}
}

func TestDialogImpl_WithExtra(t *testing.T) {
	dialog := newDialog(core.DialogTypeQuestion, "Test", nil, nil, nil, nil)

	extra := map[string]any{"data": []int64{1, 2, 3}}
	dialog.WithExtra(extra)

	result := dialog.Print(nil)
	extraMap, ok := result["extra"].(map[string]any)
	if !ok {
		t.Fatal("extra is not a map")
	}

	dataSlice, ok := extraMap["data"].([]int64)
	if !ok {
		t.Fatal("extra.data is not []int64")
	}
	if len(dataSlice) != 3 || dataSlice[0] != 1 {
		t.Errorf("extra.data = %v, want [1, 2, 3]", dataSlice)
	}
}

func TestDialogImpl_WithOptions(t *testing.T) {
	dialog := newDialog(core.DialogTypeForm, "Test", nil, nil, nil, nil)

	options := map[string]any{
		"size":   "xl",
		"filter": "active",
	}
	dialog.WithOptions(options)

	result := dialog.Print(nil)

	if result["size"] != "xl" {
		t.Errorf("size = %v, want xl", result["size"])
	}
	if result["filter"] != "active" {
		t.Errorf("filter = %v, want active", result["filter"])
	}

	// Verify nil options don't overwrite existing
	dialog.WithOptions(nil)
	result = dialog.Print(nil)
	if result["size"] != "xl" {
		t.Errorf("size after nil = %v, want xl (should not be overwritten)", result["size"])
	}
}

func TestDialogImpl_WithOption(t *testing.T) {
	dialog := newDialog(core.DialogTypeTable, "Test", nil, nil, nil, nil)

	dialog.WithOption("url", "/api/data")
	dialog.WithOption("checkTime", 5000)

	result := dialog.Print(nil)

	if result["url"] != "/api/data" {
		t.Errorf("url = %v, want /api/data", result["url"])
	}
	if result["checkTime"] != 5000 {
		t.Errorf("checkTime = %v, want 5000", result["checkTime"])
	}
}

func TestDialogImpl_Print_WithDialogContent(t *testing.T) {
	// Test that DialogContent interface is properly called
	content := DialogQuestionContent{
		Icon:     "warning",
		Question: "Are you sure?",
	}

	dialog := newDialog(core.DialogTypeQuestion, "Confirm", content, nil, nil, nil)
	result := dialog.Print(nil)

	contentMap, ok := result["content"].(map[string]any)
	if !ok {
		t.Fatal("content is not a map (DialogContent.Print() not called)")
	}

	if contentMap["icon"] != "warning" {
		t.Errorf("content.icon = %v, want warning", contentMap["icon"])
	}
	if contentMap["question"] != "Are you sure?" {
		t.Errorf("content.question = %v, want 'Are you sure?'", contentMap["question"])
	}
}

func TestDialogImpl_Print_WithPlainContent(t *testing.T) {
	// Test that plain content (not DialogContent) is passed through
	content := []map[string]any{
		{"field": "name", "label": "Name"},
	}

	dialog := newDialog(core.DialogTypeForm, "Form", content, nil, nil, nil)
	result := dialog.Print(nil)

	contentSlice, ok := result["content"].([]map[string]any)
	if !ok {
		t.Fatal("content is not a slice")
	}

	if len(contentSlice) != 1 {
		t.Fatalf("content length = %d, want 1", len(contentSlice))
	}
	if contentSlice[0]["field"] != "name" {
		t.Errorf("content[0].field = %v, want name", contentSlice[0]["field"])
	}
}

func TestDialogImpl_Print_NilOptions(t *testing.T) {
	// Verify nil options are initialized to empty map
	dialog := newDialog(core.DialogTypeQuestion, "Test", nil, nil, nil, nil)

	result := dialog.Print(nil)

	// Should not panic and should have basic structure
	if result["header"] != "Test" {
		t.Errorf("header = %v, want Test", result["header"])
	}
	if result["type"] != string(core.DialogTypeQuestion) {
		t.Errorf("type = %v, want %v", result["type"], core.DialogTypeQuestion)
	}
}

// Test that DialogQuestionContent implements DialogContent interface
func TestDialogQuestionContent_ImplementsDialogContent(t *testing.T) {
	var _ DialogContent = DialogQuestionContent{}
}

// Test that DialogWaitingContent implements DialogContent interface
func TestDialogWaitingContent_ImplementsDialogContent(t *testing.T) {
	var _ DialogContent = DialogWaitingContent{}
}

// Test DialogQuestionContent.Print()
func TestDialogQuestionContent_Print(t *testing.T) {
	content := DialogQuestionContent{
		Icon:     "help_outline",
		Question: "Continue?",
	}

	result := content.Print(nil)

	if result["icon"] != "help_outline" {
		t.Errorf("icon = %v, want help_outline", result["icon"])
	}
	if result["question"] != "Continue?" {
		t.Errorf("question = %v, want 'Continue?'", result["question"])
	}
}

// Test DialogWaitingContent.Print()
func TestDialogWaitingContent_Print(t *testing.T) {
	content := DialogWaitingContent{
		Icon: "hourglass_empty",
		Text: "Please wait...",
	}

	result := content.Print(nil)

	if result["icon"] != "hourglass_empty" {
		t.Errorf("icon = %v, want hourglass_empty", result["icon"])
	}
	if result["text"] != "Please wait..." {
		t.Errorf("text = %v, want 'Please wait...'", result["text"])
	}
}
