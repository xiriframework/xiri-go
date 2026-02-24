package page

import (
	"encoding/json"
	"testing"

	"github.com/xiriframework/xiri-go/component/url"
)

func TestNewPage_Empty(t *testing.T) {
	p := NewPage(nil)
	result := p.Print(nil)

	data, ok := result["data"].([]map[string]any)
	if !ok {
		t.Fatal("expected data to be []map[string]any")
	}
	if len(data) != 0 {
		t.Errorf("expected empty data, got %d items", len(data))
	}

	if _, hasBread := result["bread"]; hasBread {
		t.Error("expected no bread key for empty breadcrumbs")
	}
}

func TestNewPage_WithBreadcrumbs(t *testing.T) {
	p := NewPage(nil)
	u := url.NewUrl("/Portal/Device/Table")
	p.Bread("Home", nil, false)
	p.Bread("Devices", u, false)

	result := p.Print(nil)

	bread, ok := result["bread"].([]map[string]any)
	if !ok {
		t.Fatal("expected bread to be []map[string]any")
	}
	if len(bread) != 2 {
		t.Fatalf("expected 2 breadcrumbs, got %d", len(bread))
	}

	if bread[0]["label"] != "Home" {
		t.Errorf("expected 'Home', got %v", bread[0]["label"])
	}
	if bread[0]["link"] != (*string)(nil) {
		t.Errorf("expected nil link for Home, got %v", bread[0]["link"])
	}

	if bread[1]["label"] != "Devices" {
		t.Errorf("expected 'Devices', got %v", bread[1]["label"])
	}
	if bread[1]["link"] == nil {
		t.Fatal("expected non-nil link for Devices")
	}
}

func TestNewPage_WithExtra(t *testing.T) {
	p := NewPage(nil)
	p.Extra("title", "Test Page")
	p.Extra("version", 2)

	result := p.Print(nil)

	if result["title"] != "Test Page" {
		t.Errorf("expected 'Test Page', got %v", result["title"])
	}
	if result["version"] != 2 {
		t.Errorf("expected 2, got %v", result["version"])
	}
}

func TestBreadcrumbItem_Print(t *testing.T) {
	link := "/test"
	item := NewBreadcrumbItem("Test", &link, true)
	printed := item.print()

	if printed["label"] != "Test" {
		t.Errorf("expected 'Test', got %v", printed["label"])
	}
	if printed["extern"] != true {
		t.Errorf("expected extern=true, got %v", printed["extern"])
	}
}

func TestNewPage_JSONSnapshot(t *testing.T) {
	p := NewPage(nil)
	u := url.NewUrl("/Portal/Home")
	p.Bread("Home", u, false)

	result := p.Print(nil)

	// Verify it serializes to valid JSON
	_, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("failed to marshal page to JSON: %v", err)
	}
}
