package query

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/core"
)

func queryCtx() *core.UiContext {
	return &core.UiContext{Translate: func(k string) string { return k }}
}

func queryData(t *testing.T, q *Query) map[string]any {
	t.Helper()
	out := q.Print(queryCtx())
	if out["type"] != "query" {
		t.Fatalf("type=%v want query", out["type"])
	}
	data, ok := out["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data map, got %T", out["data"])
	}
	return data
}

// By default the new visibility flags must not be serialized at all.
func TestQuery_DefaultsOmitVisibilityFlags(t *testing.T) {
	q := NewQuery(nil, nil, nil)
	data := queryData(t, q)

	if _, ok := data["showActiveFilters"]; ok {
		t.Errorf("did not expect showActiveFilters key by default")
	}
	if _, ok := data["showResultCount"]; ok {
		t.Errorf("did not expect showResultCount key by default")
	}
	if _, ok := data["collapsed"]; ok {
		t.Errorf("did not expect collapsed key by default")
	}
}

// collapsed drives the expansion panel: true = collapsed, false = expanded, absent = no panel.
func TestQuery_Collapsed(t *testing.T) {
	q := NewQuery(nil, nil, nil).Collapsed(true)
	data := queryData(t, q)

	if data["collapsed"] != true {
		t.Errorf("collapsed=%v want true", data["collapsed"])
	}
}

func TestQuery_ShowActiveFilters(t *testing.T) {
	q := NewQuery(nil, nil, nil).ShowActiveFilters(true)
	data := queryData(t, q)

	if data["showActiveFilters"] != true {
		t.Errorf("showActiveFilters=%v want true", data["showActiveFilters"])
	}
	// The other flag stays absent.
	if _, ok := data["showResultCount"]; ok {
		t.Errorf("did not expect showResultCount key")
	}
}

func TestQuery_ShowResultCount(t *testing.T) {
	q := NewQuery(nil, nil, nil).ShowResultCount(true)
	data := queryData(t, q)

	if data["showResultCount"] != true {
		t.Errorf("showResultCount=%v want true", data["showResultCount"])
	}
}

// Explicit false must still serialize (nil vs. false are distinct).
func TestQuery_ShowActiveFiltersFalseIsSerialized(t *testing.T) {
	q := NewQuery(nil, nil, nil).ShowActiveFilters(false)
	data := queryData(t, q)

	val, ok := data["showActiveFilters"]
	if !ok {
		t.Fatalf("expected showActiveFilters key when explicitly set")
	}
	if val != false {
		t.Errorf("showActiveFilters=%v want false", val)
	}
}
