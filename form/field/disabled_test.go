package field

import "testing"

// SetDisabled muss das Frontend erreichen: ohne den Key im Export kann xiri-ng das Feld nicht sperren,
// der Benutzer tippt in ein Feld, dessen Wert das Binding stillschweigend verwirft.
func TestBaseFieldExportDisabled(t *testing.T) {
	ctx := newTestCtx()

	f := NewTextField("locked", "LOCKED", false, "v")
	if got, ok := f.ExportForFrontend(ctx, "v")["disabled"]; !ok || got != false {
		t.Fatalf("expected disabled=false in export by default, got %v (present=%v)", got, ok)
	}

	f.SetDisabled(true)
	if got, ok := f.ExportForFrontend(ctx, "v")["disabled"]; !ok || got != true {
		t.Fatalf("expected disabled=true in export after SetDisabled(true), got %v (present=%v)", got, ok)
	}
}
