package timeline

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
)

func testTranslator(key string) string {
	translations := map[string]string{
		"SCHRITT_1":      "Schritt 1",
		"BESCHREIBUNG_1": "Erste Beschreibung",
	}
	if v, ok := translations[key]; ok {
		return v
	}
	return key
}

func TestNewEmpty(t *testing.T) {
	tl := New()
	result := tl.Print(nil)

	if result["type"] != "timeline" {
		t.Errorf("Expected type 'timeline', got %v", result["type"])
	}
	if _, exists := result["display"]; exists {
		t.Error("Expected display to not be present when unset")
	}
	if _, exists := result["orientation"]; exists {
		t.Error("Expected orientation to not be present when unset")
	}

	items, ok := result["items"].([]map[string]any)
	if !ok {
		t.Fatal("Expected items to be []map[string]any")
	}
	if len(items) != 0 {
		t.Errorf("Expected 0 items, got %d", len(items))
	}
}

func TestAddItemMinimal(t *testing.T) {
	tl := New()
	tl.Add("Hello")

	result := tl.Print(nil)
	items := result["items"].([]map[string]any)

	if len(items) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(items))
	}
	if items[0]["title"] != "Hello" {
		t.Errorf("Expected title 'Hello', got %v", items[0]["title"])
	}
	for _, key := range []string{"description", "datetime", "icon", "iconColor"} {
		if _, exists := items[0][key]; exists {
			t.Errorf("Expected %q to not be present when unset", key)
		}
	}
}

func TestAddItemFull(t *testing.T) {
	tl := New()
	tl.Add("Order Placed").
		Description("Your order was received.").
		Datetime("10:00").
		Icon("shopping_cart").
		IconColor("primary")

	result := tl.Print(nil)
	item := result["items"].([]map[string]any)[0]

	if item["title"] != "Order Placed" {
		t.Errorf("Expected title 'Order Placed', got %v", item["title"])
	}
	if item["description"] != "Your order was received." {
		t.Errorf("Expected description, got %v", item["description"])
	}
	if item["datetime"] != "10:00" {
		t.Errorf("Expected datetime '10:00', got %v", item["datetime"])
	}
	if item["icon"] != "shopping_cart" {
		t.Errorf("Expected icon 'shopping_cart', got %v", item["icon"])
	}
	if item["iconColor"] != "primary" {
		t.Errorf("Expected iconColor 'primary', got %v", item["iconColor"])
	}
}

func TestWithOrientationHorizontal(t *testing.T) {
	tl := New().WithOrientation(core.TimelineOrientationHorizontal)
	tl.Add("Step 1")

	result := tl.Print(nil)

	if result["orientation"] != "horizontal" {
		t.Errorf("Expected orientation 'horizontal', got %v", result["orientation"])
	}
}

func TestWithOrientationVerticalExplicit(t *testing.T) {
	tl := New().WithOrientation(core.TimelineOrientationVertical)
	tl.Add("Step 1")

	result := tl.Print(nil)

	if result["orientation"] != "vertical" {
		t.Errorf("Expected orientation 'vertical' when set explicitly, got %v", result["orientation"])
	}
}

func TestWithDisplay(t *testing.T) {
	tl := New().WithDisplay("xcol-md-6")
	tl.Add("Step 1")

	result := tl.Print(nil)

	if result["display"] != "xcol-md-6" {
		t.Errorf("Expected display 'xcol-md-6', got %v", result["display"])
	}
	if _, exists := result["orientation"]; exists {
		t.Error("Expected orientation to not be present when only display is set")
	}
}

func TestMethodChaining(t *testing.T) {
	tl := New().
		WithOrientation(core.TimelineOrientationHorizontal).
		WithDisplay("xcol-12")

	tl.Add("Step 1").Icon("check").IconColor("success")
	tl.Add("Step 2").Datetime("10:00")

	result := tl.Print(nil)

	if result["orientation"] != "horizontal" {
		t.Errorf("Expected orientation 'horizontal', got %v", result["orientation"])
	}
	if result["display"] != "xcol-12" {
		t.Errorf("Expected display 'xcol-12', got %v", result["display"])
	}

	items := result["items"].([]map[string]any)
	if len(items) != 2 {
		t.Fatalf("Expected 2 items, got %d", len(items))
	}
}

func TestTranslation(t *testing.T) {
	tl := New()
	tl.Add("SCHRITT_1").Description("BESCHREIBUNG_1")

	result := tl.Print(&core.UiContext{Translate: testTranslator})
	item := result["items"].([]map[string]any)[0]

	if item["title"] != "Schritt 1" {
		t.Errorf("Expected translated title 'Schritt 1', got %v", item["title"])
	}
	if item["description"] != "Erste Beschreibung" {
		t.Errorf("Expected translated description, got %v", item["description"])
	}

	// Without translator
	result2 := tl.Print(nil)
	item2 := result2["items"].([]map[string]any)[0]

	if item2["title"] != "SCHRITT_1" {
		t.Errorf("Expected raw key 'SCHRITT_1' without translator, got %v", item2["title"])
	}
}

func TestComponentInterface(t *testing.T) {
	var _ core.Component = New()
}
