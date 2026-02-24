package emptystate

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/button"
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
)

func testTranslator(key string) string {
	translations := map[string]string{
		"KEINE_EINTRAEGE": "Keine Einträge",
		"BESCHREIBUNG":    "Es wurden keine Daten gefunden.",
		"NEUER_EINTRAG":   "Neuer Eintrag",
	}
	if v, ok := translations[key]; ok {
		return v
	}
	return key
}

func TestNewMinimal(t *testing.T) {
	es := New("inbox", core.ColorGray, "Keine Daten")
	result := es.Print(nil)

	if result["type"] != "empty-state" {
		t.Errorf("Expected type 'empty-state', got %v", result["type"])
	}
	if result["display"] != (*string)(nil) {
		t.Errorf("Expected display nil, got %v", result["display"])
	}

	data, ok := result["data"].(map[string]any)
	if !ok {
		t.Fatal("Expected data to be map[string]any")
	}
	if data["icon"] != "inbox" {
		t.Errorf("Expected icon 'inbox', got %v", data["icon"])
	}
	if data["iconColor"] != "gray" {
		t.Errorf("Expected iconColor 'gray', got %v", data["iconColor"])
	}
	if data["title"] != "Keine Daten" {
		t.Errorf("Expected title 'Keine Daten', got %v", data["title"])
	}
	if _, exists := data["description"]; exists {
		t.Error("Expected description to not be present")
	}
	if _, exists := data["button"]; exists {
		t.Error("Expected button to not be present")
	}
}

func TestWithDescription(t *testing.T) {
	es := New("search", core.ColorPrimary, "Nichts gefunden").
		WithDescription("Versuche es mit anderen Suchbegriffen.")

	result := es.Print(nil)
	data := result["data"].(map[string]any)

	if data["description"] != "Versuche es mit anderen Suchbegriffen." {
		t.Errorf("Expected description text, got %v", data["description"])
	}
}

func TestWithButton(t *testing.T) {
	btn := button.NewSimpleLinkButton("Erstellen", url.NewUrl("/create"), core.ColorPrimary)
	es := New("add_circle", core.ColorPrimary, "Keine Einträge").
		WithButton(btn)

	result := es.Print(nil)
	data := result["data"].(map[string]any)

	btnData, ok := data["button"].(map[string]any)
	if !ok {
		t.Fatal("Expected button to be map[string]any")
	}
	if btnData["action"] != "link" {
		t.Errorf("Expected button action 'link', got %v", btnData["action"])
	}
	if btnData["text"] != "Erstellen" {
		t.Errorf("Expected button text 'Erstellen', got %v", btnData["text"])
	}
}

func TestFullEmptyState(t *testing.T) {
	btn := button.NewSimpleLinkButton("Erstellen", url.NewUrl("/create"), core.ColorPrimary)
	es := New("inbox", core.ColorGray, "Keine Einträge").
		WithDescription("Es sind noch keine Einträge vorhanden.").
		WithButton(btn).
		WithDisplay("xcol-md-6")

	result := es.Print(nil)

	if result["type"] != "empty-state" {
		t.Errorf("Expected type 'empty-state', got %v", result["type"])
	}

	display := result["display"].(*string)
	if *display != "xcol-md-6" {
		t.Errorf("Expected display 'xcol-md-6', got %v", *display)
	}

	data := result["data"].(map[string]any)
	if data["icon"] != "inbox" {
		t.Errorf("Expected icon 'inbox', got %v", data["icon"])
	}
	if data["iconColor"] != "gray" {
		t.Errorf("Expected iconColor 'gray', got %v", data["iconColor"])
	}
	if data["title"] != "Keine Einträge" {
		t.Errorf("Expected title, got %v", data["title"])
	}
	if data["description"] != "Es sind noch keine Einträge vorhanden." {
		t.Errorf("Expected description, got %v", data["description"])
	}
	if _, ok := data["button"].(map[string]any); !ok {
		t.Error("Expected button to be present")
	}
}

func TestPrintData(t *testing.T) {
	es := New("inbox", core.ColorGray, "Keine Daten").
		WithDescription("Beschreibungstext").
		WithDisplay("xcol-md-6")

	data := es.PrintData(nil)

	// PrintData should NOT contain type or display
	if _, exists := data["type"]; exists {
		t.Error("PrintData should not contain 'type'")
	}
	if _, exists := data["display"]; exists {
		t.Error("PrintData should not contain 'display'")
	}

	// Should contain data fields
	if data["icon"] != "inbox" {
		t.Errorf("Expected icon 'inbox', got %v", data["icon"])
	}
	if data["iconColor"] != "gray" {
		t.Errorf("Expected iconColor 'gray', got %v", data["iconColor"])
	}
	if data["title"] != "Keine Daten" {
		t.Errorf("Expected title 'Keine Daten', got %v", data["title"])
	}
	if data["description"] != "Beschreibungstext" {
		t.Errorf("Expected description 'Beschreibungstext', got %v", data["description"])
	}
}

func TestTranslation(t *testing.T) {
	es := New("inbox", core.ColorGray, "KEINE_EINTRAEGE").
		WithDescription("BESCHREIBUNG")

	// With translator
	result := es.Print(testTranslator)
	data := result["data"].(map[string]any)

	if data["title"] != "Keine Einträge" {
		t.Errorf("Expected translated title 'Keine Einträge', got %v", data["title"])
	}
	if data["description"] != "Es wurden keine Daten gefunden." {
		t.Errorf("Expected translated description, got %v", data["description"])
	}

	// Without translator (nil)
	result2 := es.Print(nil)
	data2 := result2["data"].(map[string]any)

	if data2["title"] != "KEINE_EINTRAEGE" {
		t.Errorf("Expected raw key 'KEINE_EINTRAEGE', got %v", data2["title"])
	}
}

func TestPrintDataWithTranslation(t *testing.T) {
	es := New("inbox", core.ColorGray, "KEINE_EINTRAEGE").
		WithDescription("BESCHREIBUNG")

	data := es.PrintData(testTranslator)

	if data["title"] != "Keine Einträge" {
		t.Errorf("Expected translated title, got %v", data["title"])
	}
	if data["description"] != "Es wurden keine Daten gefunden." {
		t.Errorf("Expected translated description, got %v", data["description"])
	}
}

func TestComponentInterface(t *testing.T) {
	// Verify EmptyState implements core.Component
	var _ core.Component = New("inbox", core.ColorGray, "Test")
}

func TestMethodChaining(t *testing.T) {
	btn := button.NewSimpleLinkButton("Action", url.NewUrl("/action"), core.ColorPrimary)

	es := New("inbox", core.ColorGray, "Title").
		WithDescription("Desc").
		WithButton(btn).
		WithDisplay("xcol-12")

	// Should not panic and return valid result
	result := es.Print(nil)
	if result["type"] != "empty-state" {
		t.Error("Method chaining should produce valid result")
	}
}
