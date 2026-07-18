package stat_test

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/stat"
	"github.com/xiriframework/xiri-go/component/url"
)

func TestStat_Reference_AddsField(t *testing.T) {
	s := stat.New(1.4, "PF").Reference("Gate ≥ 1,1")

	d := s.PrintData(nil)

	if d["reference"] != "Gate ≥ 1,1" {
		t.Errorf("d[\"reference\"] = %v, want %q", d["reference"], "Gate ≥ 1,1")
	}
}

func TestStat_NoReference_OmitsField(t *testing.T) {
	s := stat.New(1.4, "PF")

	d := s.PrintData(nil)

	if _, has := d["reference"]; has {
		t.Errorf("expected no \"reference\" key, got %v", d["reference"])
	}
}

func TestStat_Link_AddsField(t *testing.T) {
	s := stat.New(12, "Offen").Link(url.NewUrl("/Orders/Table?status=open"))

	d := s.PrintData(nil)

	if d["link"] != "/Orders/Table?status=open" {
		t.Errorf("d[\"link\"] = %v, want %q", d["link"], "/Orders/Table?status=open")
	}
}

func TestStat_NoLink_OmitsField(t *testing.T) {
	d := stat.New(12, "Offen").PrintData(nil)

	if _, has := d["link"]; has {
		t.Errorf("expected no \"link\" key, got %v", d["link"])
	}
}
