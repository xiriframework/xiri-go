package card

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/barchart"
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/stat"
)

func cardCtx() *core.UiContext {
	return &core.UiContext{Translate: func(k string) string { return k }}
}

// TestCard_MultiComponent verifies that Card.Add(...) switches the JSON
// output to multi-component mode (data.components[]) and skips table fields.
func TestCard_MultiComponent(t *testing.T) {
	c := NewCard(core.CardTypeTable, nil, "Activity", nil, nil, nil, false, false, nil)
	c.Add(barchart.New("activity").Mode(barchart.ModeSimple).Bar("M", 3))
	c.Add(stat.New("18h", "Today").Compact())
	c.Add(stat.New("32h", "Last 7 days").Compact())

	out := c.Print(cardCtx())
	if out["type"] != "card" {
		t.Errorf("type=%v want card", out["type"])
	}
	data := out["data"].(map[string]any)

	comps, ok := data["components"].([]map[string]any)
	if !ok {
		t.Fatalf("expected data.components []map[string]any, got %T", data["components"])
	}
	if len(comps) != 3 {
		t.Fatalf("expected 3 components, got %d", len(comps))
	}
	if comps[0]["type"] != "barchart" {
		t.Errorf("comps[0].type=%v want barchart", comps[0]["type"])
	}
	if comps[1]["type"] != "stat" {
		t.Errorf("comps[1].type=%v want stat", comps[1]["type"])
	}

	// The Stat sub-components must carry compact: true via printData.
	stat1Data := comps[1]["data"].(map[string]any)
	if stat1Data["compact"] != true {
		t.Errorf("comps[1].data.compact=%v want true", stat1Data["compact"])
	}

	// Multi-component mode must NOT emit table-fields/content.
	if _, hasFields := data["fields"]; hasFields {
		t.Errorf("did not expect 'fields' in multi-component card data")
	}
	if _, hasContent := data["content"]; hasContent {
		t.Errorf("did not expect 'content' in multi-component card data")
	}
}

// TestCard_TableModeUnchanged verifies the existing CardListContent rendering
// still produces fields/data/dense (no regression).
func TestCard_TableModeUnchanged(t *testing.T) {
	content := NewCardListContent([]CardListContentLine{
		{Name: "Status", Content: "Active"},
		{Name: "Version", Content: "2.1"},
	})
	c := NewCardList("Details", content)

	out := c.Print(cardCtx())
	data := out["data"].(map[string]any)

	if _, hasComponents := data["components"]; hasComponents {
		t.Errorf("expected NO 'components' on table-mode card")
	}
	if data["fields"] == nil {
		t.Errorf("expected 'fields' on table-mode card")
	}
	if data["data"] == nil {
		t.Errorf("expected 'data' on table-mode card")
	}
}

// TestCard_StatCompactDefaultFalse verifies a default stat (no .Compact())
// does not emit compact:false (omitempty semantic).
func TestCard_StatCompactDefaultFalse(t *testing.T) {
	out := stat.New("100", "Total").Print(cardCtx())
	data := out["data"].(map[string]any)
	if _, has := data["compact"]; has {
		t.Errorf("expected 'compact' to be absent on default stat")
	}
}
