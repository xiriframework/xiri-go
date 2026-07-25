package timeline

import "testing"

// #20 (sibling): a pointer returned by Add must stay valid after later appends reallocate the slice.
func TestAddPointerStableAfterRealloc(t *testing.T) {
	tl := New()

	first := tl.Add("title0")
	for i := 1; i < 100; i++ {
		tl.Add("t")
	}
	first.Icon("star")

	out := tl.Print(nil)
	items := out["items"].([]map[string]any)
	if items[0]["icon"] != "star" {
		t.Errorf("first item icon = %v, want %q (stale pointer after realloc)", items[0]["icon"], "star")
	}
}
