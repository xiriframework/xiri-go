package multistat

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/stat"
	"github.com/xiriframework/xiri-go/component/url"
)

func TestPrintDataItems(t *testing.T) {
	m := New().
		Add(stat.New(12, "Offen").Icon("inventory").Color("orange")).
		Add(stat.New(45, "Fertig").Color("green"))

	d := m.PrintData(nil)

	items, ok := d["items"].([]map[string]any)
	if !ok {
		t.Fatalf("items missing or wrong type: %T", d["items"])
	}
	if len(items) != 2 {
		t.Fatalf("want 2 items, got %d", len(items))
	}
	if items[0]["value"] != 12 {
		t.Errorf("item[0] value: want 12, got %v", items[0]["value"])
	}
	if items[0]["color"] != "orange" {
		t.Errorf("item[0] color: want orange, got %v", items[0]["color"])
	}
	if items[0]["icon"] != "inventory" {
		t.Errorf("item[0] icon: want inventory, got %v", items[0]["icon"])
	}
}

func TestHeaderOptional(t *testing.T) {
	d := New().PrintData(nil)
	if _, ok := d["title"]; ok {
		t.Error("title should be absent when not set")
	}
	if _, ok := d["icon"]; ok {
		t.Error("icon should be absent when not set")
	}
	if _, ok := d["iconColor"]; ok {
		t.Error("iconColor should be absent when not set")
	}

	d = New().Title("Bestellungen").Icon("shopping_cart").IconColor("primary").PrintData(nil)
	if d["title"] != "Bestellungen" {
		t.Errorf("title: want Bestellungen, got %v", d["title"])
	}
	if d["icon"] != "shopping_cart" {
		t.Errorf("icon: want shopping_cart, got %v", d["icon"])
	}
	if d["iconColor"] != "primary" {
		t.Errorf("iconColor: want primary, got %v", d["iconColor"])
	}
}

func TestURLMode(t *testing.T) {
	m := New().
		Add(stat.New(1, "a")).
		Title("Bestellungen").
		Icon("shopping_cart").
		SetURL(url.NewUrlPrefix("/multistat", "/api")).
		WithReload(true)

	data := m.Print(nil)["data"].(map[string]any)

	if data["url"] != "/api/multistat" {
		t.Errorf("url = %v, want /api/multistat", data["url"])
	}
	if data["reload"] != true {
		t.Errorf("reload = %v, want true", data["reload"])
	}
	// Items werden im AJAX-Modus vom Frontend nachgeladen → nil.
	if data["items"] != nil {
		t.Errorf("items = %v, want nil in url mode", data["items"])
	}
	// Header bleibt inline sichtbar, während Items laden.
	if data["title"] != "Bestellungen" {
		t.Errorf("title = %v, want Bestellungen", data["title"])
	}
	if data["icon"] != "shopping_cart" {
		t.Errorf("icon = %v, want shopping_cart", data["icon"])
	}
}

func TestVerticalItems(t *testing.T) {
	// Ohne Aufruf: horizontal (Standard) → kein Key.
	d := New().Add(stat.New(1, "a")).PrintData(nil)
	if _, has := d["verticalItems"]; has {
		t.Errorf("expected no verticalItems key by default, got %v", d["verticalItems"])
	}

	// Mit VerticalItems(): opt-in vertikal → Key true.
	d = New().Add(stat.New(1, "a")).VerticalItems().PrintData(nil)
	if d["verticalItems"] != true {
		t.Errorf("verticalItems = %v, want true", d["verticalItems"])
	}
}

func TestPrintEnvelope(t *testing.T) {
	r := New().WithDisplay("xcol-6").Print(nil)
	if r["type"] != "multi-stat" {
		t.Errorf("type: want multi-stat, got %v", r["type"])
	}
	if r["display"] != "xcol-6" {
		t.Errorf("display: want xcol-6, got %v", r["display"])
	}
	if _, ok := r["data"].(map[string]any); !ok {
		t.Errorf("data missing or wrong type: %T", r["data"])
	}
}
