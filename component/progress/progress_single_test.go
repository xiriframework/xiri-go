package progress_test

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/progress"
)

func TestProgress_Print_ComputesValue(t *testing.T) {
	out := progress.NewProgress(3, 10, "Symbole").Print(nil)

	if out["type"] != "progress" {
		t.Fatalf("out[\"type\"] = %v, want \"progress\"", out["type"])
	}
	d, ok := out["data"].(map[string]any)
	if !ok {
		t.Fatalf("out[\"data\"] not a map, got %T", out["data"])
	}
	if d["value"] != float64(30) {
		t.Errorf("d[\"value\"] = %v, want 30", d["value"])
	}
	if d["current"] != 3 || d["total"] != 10 {
		t.Errorf("current/total = %v/%v, want 3/10", d["current"], d["total"])
	}
	if d["label"] != "Symbole" {
		t.Errorf("d[\"label\"] = %v, want Symbole", d["label"])
	}
}

func TestProgress_Print_ZeroTotalNoPanic(t *testing.T) {
	out := progress.NewProgress(0, 0, "x").Print(nil)

	d := out["data"].(map[string]any)
	if d["value"] != float64(0) {
		t.Errorf("d[\"value\"] = %v, want 0", d["value"])
	}
}
