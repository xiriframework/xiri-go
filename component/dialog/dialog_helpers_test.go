package dialog

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
)

func TestResolveText(t *testing.T) {
	mockTranslator := func(key string) string {
		translations := map[string]string{
			"Delete":  "Löschen",
			"Ok":      "Bestätigen",
			"Back":    "Zurück",
			"Warning": "Warnung",
		}
		if val, ok := translations[key]; ok {
			return val
		}
		return key
	}

	tests := []struct {
		name           string
		customText     *string
		translationKey string
		defaultText    string
		translator     core.TranslateFunc
		expected       string
	}{
		{
			name:           "Custom text takes priority",
			customText:     stringPtr("Custom Delete"),
			translationKey: "Delete",
			defaultText:    "Delete",
			translator:     mockTranslator,
			expected:       "Custom Delete",
		},
		{
			name:           "Translation when no custom text",
			customText:     nil,
			translationKey: "Delete",
			defaultText:    "Delete",
			translator:     mockTranslator,
			expected:       "Löschen",
		},
		{
			name:           "Default when no custom text or translator",
			customText:     nil,
			translationKey: "Delete",
			defaultText:    "Delete",
			translator:     nil,
			expected:       "Delete",
		},
		{
			name:           "Empty custom text is used",
			customText:     stringPtr(""),
			translationKey: "Delete",
			defaultText:    "Delete",
			translator:     mockTranslator,
			expected:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolveText(tt.customText, tt.translationKey, tt.defaultText, tt.translator)
			if result != tt.expected {
				t.Errorf("resolveText() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestResolveDialogTexts(t *testing.T) {
	mockTranslator := func(key string) string {
		translations := map[string]string{
			"Delete":  "Löschen",
			"Ok":      "Bestätigen",
			"Back":    "Zurück",
			"Warning": "Warnung",
		}
		if val, ok := translations[key]; ok {
			return val
		}
		return key
	}

	tests := []struct {
		name       string
		headerText *string
		okText     *string
		closeText  *string
		translator core.TranslateFunc
		expected   dialogTexts
	}{
		{
			name:       "All custom texts",
			headerText: stringPtr("Custom Header"),
			okText:     stringPtr("Custom OK"),
			closeText:  stringPtr("Custom Close"),
			translator: mockTranslator,
			expected: dialogTexts{
				Header: "Custom Header",
				Ok:     "Custom OK",
				Close:  "Custom Close",
			},
		},
		{
			name:       "All translations",
			headerText: nil,
			okText:     nil,
			closeText:  nil,
			translator: mockTranslator,
			expected: dialogTexts{
				Header: "Löschen",
				Ok:     "Bestätigen",
				Close:  "Zurück",
			},
		},
		{
			name:       "All defaults (no translator)",
			headerText: nil,
			okText:     nil,
			closeText:  nil,
			translator: nil,
			expected: dialogTexts{
				Header: "Delete",
				Ok:     "Ok",
				Close:  "Back",
			},
		},
		{
			name:       "Mixed custom and translation",
			headerText: stringPtr("My Header"),
			okText:     nil,
			closeText:  nil,
			translator: mockTranslator,
			expected: dialogTexts{
				Header: "My Header",
				Ok:     "Bestätigen",
				Close:  "Zurück",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolveDialogTexts(
				tt.headerText, "Delete", "Delete",
				tt.okText, "Ok", "Ok",
				tt.closeText, "Back", "Back",
				tt.translator,
			)
			if result != tt.expected {
				t.Errorf("resolveDialogTexts() = %+v, want %+v", result, tt.expected)
			}
		})
	}
}

func TestBuildStandardButtons(t *testing.T) {
	testUrl := url.NewUrl("/test/action")

	buttons := buildStandardButtons("Zurück", "Bestätigen", testUrl)

	if len(buttons) != 2 {
		t.Fatalf("Expected 2 buttons, got %d", len(buttons))
	}

	// Check Close button (first button)
	closeBtn := buttons[0].Print(nil)
	if closeBtn["action"] != string(core.ButtonActionClose) {
		t.Errorf("First button action = %v, want %v", closeBtn["action"], core.ButtonActionClose)
	}
	if closeBtn["text"] != "Zurück" {
		t.Errorf("First button text = %v, want %v", closeBtn["text"], "Zurück")
	}
	if closeBtn["type"] != string(core.ButtonTypeStroked) {
		t.Errorf("First button type = %v, want %v", closeBtn["type"], core.ButtonTypeStroked)
	}
	if closeBtn["default"] != false {
		t.Errorf("First button default = %v, want false", closeBtn["default"])
	}

	// Check Form button (second button)
	formBtn := buttons[1].Print(nil)
	if formBtn["action"] != string(core.ButtonActionForm) {
		t.Errorf("Second button action = %v, want %v", formBtn["action"], core.ButtonActionForm)
	}
	if formBtn["text"] != "Bestätigen" {
		t.Errorf("Second button text = %v, want %v", formBtn["text"], "Bestätigen")
	}
	if formBtn["type"] != string(core.ButtonTypeRaised) {
		t.Errorf("Second button type = %v, want %v", formBtn["type"], core.ButtonTypeRaised)
	}
	if formBtn["default"] != true {
		t.Errorf("Second button default = %v, want true", formBtn["default"])
	}
	if formBtn["url"] != testUrl.PrintPrefix() {
		t.Errorf("Second button url = %v, want %v", formBtn["url"], testUrl.PrintPrefix())
	}
}

// Helper function to create string pointers for tests
func stringPtr(s string) *string {
	return &s
}
