package url

import "testing"

func TestNewUrl(t *testing.T) {
	u := NewUrlPrefix("/Portal/Device/Table", "/api")

	if u.Print() != "/Portal/Device/Table" {
		t.Errorf("expected '/Portal/Device/Table', got %q", u.Print())
	}
	if u.PrintPrefix() != "/api/Portal/Device/Table" {
		t.Errorf("expected '/api/Portal/Device/Table', got %q", u.PrintPrefix())
	}
}

func TestUrl_Add(t *testing.T) {
	u := NewUrl("/Portal/Device")
	result := u.Add("7")

	// Should return same pointer for chaining
	if result != u {
		t.Error("expected Add to return same pointer")
	}
	if u.Print() != "/Portal/Device/7" {
		t.Errorf("expected '/Portal/Device/7', got %q", u.Print())
	}
}

func TestUrl_AddChaining(t *testing.T) {
	u := NewUrl("/Portal").Add("Device").Add("Edit").Add("7")

	if u.Print() != "/Portal/Device/Edit/7" {
		t.Errorf("expected '/Portal/Device/Edit/7', got %q", u.Print())
	}
}

func TestUrl_AddPrefix(t *testing.T) {
	u := NewUrlPrefix("/Table", "api")
	result := u.AddPrefix("v2")

	if result != u {
		t.Error("expected AddPrefix to return same pointer")
	}
	if u.PrintPrefix() != "v2/api/Table" {
		t.Errorf("expected 'v2/api/Table', got %q", u.PrintPrefix())
	}
}

func TestUrl_EmptyPrefix(t *testing.T) {
	u := NewUrl("/test")
	if u.PrintPrefix() != "/test" {
		t.Errorf("expected '/test', got %q", u.PrintPrefix())
	}
}
