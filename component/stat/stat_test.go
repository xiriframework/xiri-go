package stat_test

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/stat"
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
