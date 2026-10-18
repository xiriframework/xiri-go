package table

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
	xurl "github.com/xiriframework/xiri-go/component/url"
)

func TestSetReorderURL_ExportsPrefixedURL(t *testing.T) {
	opts := NewBuilder[testOptionRow]().
		SetReorderURL(xurl.NewUrlPrefix("/rules/reorder", "/api")).
		Build().exportOptions(&core.UiContext{})

	if opts["reorderUrl"] != "/api/rules/reorder" {
		t.Errorf("reorderUrl = %v, want /api/rules/reorder", opts["reorderUrl"])
	}
}

func TestSetReorderURL_NilOrUnsetOmitsOption(t *testing.T) {
	for name, b := range map[string]*TableBuilder[testOptionRow]{
		"unset": NewBuilder[testOptionRow](),
		"nil":   NewBuilder[testOptionRow]().SetReorderURL(xurl.NewUrl("/x")).SetReorderURL(nil),
		"empty": NewBuilder[testOptionRow]().SetReorderURL(xurl.NewUrl("")),
	} {
		if _, ok := b.Build().exportOptions(&core.UiContext{})["reorderUrl"]; ok {
			t.Errorf("%s: reorderUrl must not be exported", name)
		}
	}
}

func TestReorderRequest_BindsIdAndBefore(t *testing.T) {
	var req ReorderRequest
	if err := json.Unmarshal([]byte(`{"id": 7, "before": 3}`), &req); err != nil {
		t.Fatal(err)
	}
	if req.ID != 7 || req.Before == nil || *req.Before != 3 {
		t.Errorf("got %+v", req)
	}

	req = ReorderRequest{}
	if err := json.Unmarshal([]byte(`{"id": 7, "before": null}`), &req); err != nil {
		t.Fatal(err)
	}
	if req.ID != 7 || req.Before != nil {
		t.Errorf("before:null must be nil (= end), got %+v", req)
	}
}

func captureWarnings(t *testing.T, build func()) string {
	t.Helper()
	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	defer slog.SetDefault(prev)
	build()
	return buf.String()
}

func TestBuild_ReorderWarnings(t *testing.T) {
	idField := func(b *TableBuilder[testOptionRow]) *FieldBuilder {
		return b.IdField("id", "ID", func(r testOptionRow) int64 { return r.ID })
	}
	u := xurl.NewUrl("/reorder")

	cases := map[string]struct {
		build func()
		want  string
	}{
		"no id field": {func() {
			NewBuilder[testOptionRow]().SetReorderURL(u).Build()
		}, "id"},
		"hidden id field": {func() {
			b := NewBuilder[testOptionRow]().SetReorderURL(u)
			idField(b).Hide()
			b.Build()
		}, "id"},
		"server side": {func() {
			b := NewBuilder[testOptionRow]().SetReorderURL(u).SetServerSide(true)
			idField(b)
			b.Build()
		}, "ServerSide"},
		"tree": {func() {
			b := NewBuilder[testOptionRow]().SetReorderURL(u).Tree(TreeConfig{IdField: "id", ParentIdField: "id"})
			idField(b)
			b.Build()
		}, "Tree"},
	}
	for name, c := range cases {
		got := captureWarnings(t, c.build)
		if !strings.Contains(got, "ReorderURL") || !strings.Contains(got, c.want) {
			t.Errorf("%s: expected a ReorderURL warning mentioning %q, got %q", name, c.want, got)
		}
	}

	ok := captureWarnings(t, func() {
		b := NewBuilder[testOptionRow]().SetReorderURL(u)
		idField(b)
		b.Build()
		NewBuilder[testOptionRow]().Build()
	})
	if ok != "" {
		t.Errorf("expected no warning, got %q", ok)
	}
}
