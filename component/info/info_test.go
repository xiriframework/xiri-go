package info

import "testing"

func TestInfoTextPrint(t *testing.T) {
	it := NewInfoText("plain", nil)
	data := it.Print(nil)["data"].(map[string]any)
	if data["text"] != "plain" {
		t.Errorf("expected text 'plain', got %v", data["text"])
	}
	if _, ok := data["html"]; ok {
		t.Errorf("expected no html key without WithHtml(), got %v", data["html"])
	}
}

func TestInfoTextPrintWithHtml(t *testing.T) {
	it := NewInfoText("<b>bold</b>", nil).WithHtml()
	data := it.Print(nil)["data"].(map[string]any)
	if data["text"] != "<b>bold</b>" {
		t.Errorf("expected text '<b>bold</b>', got %v", data["text"])
	}
	if data["html"] != true {
		t.Errorf("expected html true with WithHtml(), got %v", data["html"])
	}
}

func TestInfoPointPrint(t *testing.T) {
	ip := NewInfoPoint("192.168.1.1", "lan", "primary", nil, nil, nil, nil, nil, nil)
	data := ip.Print(nil)["data"].(map[string]any)
	if data["info"] != "192.168.1.1" {
		t.Errorf("expected info '192.168.1.1', got %v", data["info"])
	}
	if _, ok := data["html"]; ok {
		t.Errorf("expected no html key without WithHtml(), got %v", data["html"])
	}
}

func TestInfoPointPrintWithHtml(t *testing.T) {
	ip := NewInfoPoint("<b>bold</b>", "lan", "primary", nil, nil, nil, nil, nil, nil).WithHtml()
	data := ip.Print(nil)["data"].(map[string]any)
	if data["html"] != true {
		t.Errorf("expected html true with WithHtml(), got %v", data["html"])
	}
}
