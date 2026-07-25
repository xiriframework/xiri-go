package descriptionlist

import "testing"

// #20: A pointer returned by Add must stay valid after later appends reallocate the slice.
func TestAddPointerStableAfterRealloc(t *testing.T) {
	d := New()

	first := d.Add("label0", "value0")
	// Force several reallocations of the backing array.
	for i := 1; i < 100; i++ {
		d.Add("l", "v")
	}

	// Configure the first item AFTER the reallocations.
	first.Icon("star")

	out := d.Print(nil)
	items := out["data"].(map[string]any)["items"].([]map[string]any)
	if items[0]["icon"] != "star" {
		t.Errorf("first item icon = %v, want %q (stale pointer after realloc)", items[0]["icon"], "star")
	}
}
