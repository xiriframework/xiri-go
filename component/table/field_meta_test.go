package table

import "testing"

// GetFieldMetas is the only view external renderers (PDF/Excel) get on a field.
// Without header/headerSpan they cannot reconstruct grouped column headers that
// WithHeader/WithHeaderSpan produce in the JSON export.
func TestGetFieldMetas_ExposesHeaderAndSpan(t *testing.T) {
	b := NewBuilder[testOptionRow]()
	b.IdField("id", "row.id", func(r testOptionRow) int64 { return r.ID }).
		WithHeader("Gruppe A").WithHeaderSpan(2)
	b.TextField("name", "row.name", func(r testOptionRow) string { return r.Name })

	metas := b.Build().GetFieldMetas()
	if len(metas) != 2 {
		t.Fatalf("len(metas) = %d, want 2", len(metas))
	}

	if metas[0].Header == nil || *metas[0].Header != "Gruppe A" {
		t.Errorf("metas[0].Header = %v, want %q", metas[0].Header, "Gruppe A")
	}
	if metas[0].HeaderSpan == nil || *metas[0].HeaderSpan != 2 {
		t.Errorf("metas[0].HeaderSpan = %v, want 2", metas[0].HeaderSpan)
	}

	if metas[1].Header != nil {
		t.Errorf("metas[1].Header = %v, want nil", *metas[1].Header)
	}
	if metas[1].HeaderSpan != nil {
		t.Errorf("metas[1].HeaderSpan = %v, want nil", *metas[1].HeaderSpan)
	}
}
