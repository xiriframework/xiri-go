package table

import (
	"testing"

	"github.com/xiriframework/xiri-go/component/button"
	"github.com/xiriframework/xiri-go/component/core"
	xurl "github.com/xiriframework/xiri-go/component/url"
	"github.com/xiriframework/xiri-go/types/distance"
	"github.com/xiriframework/xiri-go/types/language"
	"github.com/xiriframework/xiri-go/types/locale"
	"github.com/xiriframework/xiri-go/types/timezone"
	"github.com/xiriframework/xiri-go/uicontext"
)

// Test row struct
type testOptionRow struct {
	ID   int64
	Name string
}

// Test context
func testOptionContext() *uicontext.UiContext {
	return &uicontext.UiContext{
		Timezone: timezone.EuropeVienna,
		Lang:     language.Deutsch,
		Locale:   locale.De,
		Distance: distance.Kilometer,
	}
}

// Test translator
func testOptionTranslator(key string) string {
	return key
}

// TestSetBooleanOptions tests all boolean option setters
func TestSetBooleanOptions(t *testing.T) {
	ctx := testOptionContext()
	builder := NewBuilder[testOptionRow](ctx, testOptionTranslator)

	builder.
		SetDense(true).
		SetPagination(false).
		SetSearch(true).
		SetReload(false).
		SetQuery(true).
		SetCsv(false).
		SetSaveState(true).
		SetBorders(false).
		SetBordersHeader(true).
		SetFooter(false)

	tbl := builder.Build()
	opts := tbl.GetOptions()

	if opts.Dense == nil || *opts.Dense != true {
		t.Error("Expected Dense to be true")
	}
	if opts.Pagination == nil || *opts.Pagination != false {
		t.Error("Expected Pagination to be false")
	}
	if opts.Search == nil || *opts.Search != true {
		t.Error("Expected Search to be true")
	}
	if opts.Reload == nil || *opts.Reload != false {
		t.Error("Expected Reload to be false")
	}
	if opts.Query == nil || *opts.Query != true {
		t.Error("Expected Query to be true")
	}
	if opts.Csv == nil || *opts.Csv != false {
		t.Error("Expected Csv to be false")
	}
	if opts.SaveState == nil || *opts.SaveState != true {
		t.Error("Expected SaveState to be true")
	}
	if opts.Borders == nil || *opts.Borders != false {
		t.Error("Expected Borders to be false")
	}
	if opts.BordersHeader == nil || *opts.BordersHeader != true {
		t.Error("Expected BordersHeader to be true")
	}
	if opts.Footer == nil || *opts.Footer != false {
		t.Error("Expected Footer to be false")
	}
}

// TestSetStringOptions tests all string option setters
func TestSetStringOptions(t *testing.T) {
	ctx := testOptionContext()
	builder := NewBuilder[testOptionRow](ctx, testOptionTranslator)

	builder.
		SetClass("custom-table").
		SetTitle("Test Table").
		SetTextNoData("No data available").
		SetMinWidth("800px").
		SetSaveStateId("my-table-state").
		SetSaveInput("input-field").
		SetSaveInputUrl("/save").
		SetDisplay("xcol-12")

	tbl := builder.Build()
	opts := tbl.GetOptions()

	if opts.Class == nil || *opts.Class != "custom-table" {
		t.Error("Expected Class to be 'custom-table'")
	}
	if opts.Title == nil || *opts.Title != "Test Table" {
		t.Error("Expected Title to be 'Test Table'")
	}
	if opts.TextNoData == nil || *opts.TextNoData != "No data available" {
		t.Error("Expected TextNoData to be 'No data available'")
	}
	if opts.MinWidth == nil || *opts.MinWidth != "800px" {
		t.Error("Expected MinWidth to be '800px'")
	}
	if opts.SaveStateId == nil || *opts.SaveStateId != "my-table-state" {
		t.Error("Expected SaveStateId to be 'my-table-state'")
	}
	if opts.SaveInput == nil || *opts.SaveInput != "input-field" {
		t.Error("Expected SaveInput to be 'input-field'")
	}
	if opts.SaveInputUrl == nil || *opts.SaveInputUrl != "/save" {
		t.Error("Expected SaveInputUrl to be '/save'")
	}
	if opts.Display == nil || *opts.Display != "xcol-12" {
		t.Error("Expected Display to be 'xcol-12'")
	}
}

// TestSetNumericOptions tests numeric option setters
func TestSetNumericOptions(t *testing.T) {
	ctx := testOptionContext()
	builder := NewBuilder[testOptionRow](ctx, testOptionTranslator)

	builder.SetItemsPerPage(50)

	tbl := builder.Build()
	opts := tbl.GetOptions()

	if opts.ItemsPerPage == nil || *opts.ItemsPerPage != 50 {
		t.Error("Expected ItemsPerPage to be 50")
	}
}

// TestSetSliceOptions tests slice option setters
func TestSetSliceOptions(t *testing.T) {
	ctx := testOptionContext()
	builder := NewBuilder[testOptionRow](ctx, testOptionTranslator)

	pageSizes := []int{10, 25, 50, 100}
	builder.SetPageSizes(pageSizes)

	tbl := builder.Build()
	opts := tbl.GetOptions()

	if len(opts.PageSizes) != 4 {
		t.Errorf("Expected PageSizes to have 4 items, got %d", len(opts.PageSizes))
	}
	if opts.PageSizes[0] != 10 || opts.PageSizes[3] != 100 {
		t.Error("Expected PageSizes to be [10, 25, 50, 100]")
	}
}

// TestSetSelectButtonsAutoLogic tests the auto-logic for Select when SetSelectButtons is called
func TestSetSelectButtonsAutoLogic(t *testing.T) {
	ctx := testOptionContext()

	// Test 1: Setting non-empty buttons auto-enables Select
	builder1 := NewBuilder[testOptionRow](ctx, testOptionTranslator)
	btn1 := button.NewTableButton(core.ButtonActionLink, "edit", xurl.NewUrl("/edit"), "Edit tooltip", core.ColorPrimary, false, nil)
	btn2 := button.NewTableButton(core.ButtonActionDialog, "delete", xurl.NewUrl("/delete"), "Delete tooltip", core.ColorWarning, false, nil)

	builder1.SetSelectButtons([]*button.TableButton{btn1, btn2})
	tbl1 := builder1.Build()
	opts1 := tbl1.GetOptions()

	if opts1.Select == nil || *opts1.Select != true {
		t.Error("Expected Select to be auto-enabled when SetSelectButtons called with non-empty slice")
	}
	if len(opts1.SelectButtons) != 2 {
		t.Errorf("Expected 2 select buttons, got %d", len(opts1.SelectButtons))
	}

	// Test 2: Setting empty buttons auto-disables Select
	builder2 := NewBuilder[testOptionRow](ctx, testOptionTranslator)
	builder2.SetSelectButtons([]*button.TableButton{})
	tbl2 := builder2.Build()
	opts2 := tbl2.GetOptions()

	if opts2.Select == nil || *opts2.Select != false {
		t.Error("Expected Select to be auto-disabled when SetSelectButtons called with empty slice")
	}

	// Test 3: Setting nil buttons auto-disables Select
	builder3 := NewBuilder[testOptionRow](ctx, testOptionTranslator)
	builder3.SetSelectButtons(nil)
	tbl3 := builder3.Build()
	opts3 := tbl3.GetOptions()

	if opts3.Select == nil || *opts3.Select != false {
		t.Error("Expected Select to be auto-disabled when SetSelectButtons called with nil")
	}
}

// TestAddSelectButton tests the AddSelectButton helper method
func TestAddSelectButton(t *testing.T) {
	ctx := testOptionContext()
	builder := NewBuilder[testOptionRow](ctx, testOptionTranslator)

	btn1 := button.NewTableButton(core.ButtonActionLink, "edit", xurl.NewUrl("/edit"), "Edit tooltip", core.ColorPrimary, false, nil)
	btn2 := button.NewTableButton(core.ButtonActionDialog, "delete", xurl.NewUrl("/delete"), "Delete tooltip", core.ColorWarning, false, nil)
	btn3 := button.NewTableButton(core.ButtonActionDownload, "download", xurl.NewUrl("/export"), "Export tooltip", core.ColorAccent, false, nil)

	builder.
		AddSelectButton(btn1).
		AddSelectButton(btn2).
		AddSelectButton(btn3)

	tbl := builder.Build()
	opts := tbl.GetOptions()

	if opts.Select == nil || *opts.Select != true {
		t.Error("Expected Select to be auto-enabled when AddSelectButton is called")
	}
	if len(opts.SelectButtons) != 3 {
		t.Errorf("Expected 3 select buttons, got %d", len(opts.SelectButtons))
	}
}

// TestClearSelectButtons tests the ClearSelectButtons method
func TestClearSelectButtons(t *testing.T) {
	ctx := testOptionContext()
	builder := NewBuilder[testOptionRow](ctx, testOptionTranslator)

	// Add buttons first
	btn1 := button.NewTableButton(core.ButtonActionLink, "edit", xurl.NewUrl("/edit"), "Edit tooltip", core.ColorPrimary, false, nil)
	builder.AddSelectButton(btn1)

	// Then clear them
	builder.ClearSelectButtons()

	tbl := builder.Build()
	opts := tbl.GetOptions()

	if opts.Select == nil || *opts.Select != false {
		t.Error("Expected Select to be disabled after ClearSelectButtons")
	}
	if opts.SelectButtons != nil {
		t.Error("Expected SelectButtons to be nil after ClearSelectButtons")
	}
}

// TestSelectButtonsOverrideSelect tests that explicit SetSelect can override auto-logic
func TestSelectButtonsOverrideSelect(t *testing.T) {
	ctx := testOptionContext()
	builder := NewBuilder[testOptionRow](ctx, testOptionTranslator)

	btn1 := button.NewTableButton(core.ButtonActionLink, "edit", xurl.NewUrl("/edit"), "Edit tooltip", core.ColorPrimary, false, nil)

	// Set buttons (auto-enables Select), then explicitly disable it
	builder.
		SetSelectButtons([]*button.TableButton{btn1}).
		SetSelect(false)

	tbl := builder.Build()
	opts := tbl.GetOptions()

	if opts.Select == nil || *opts.Select != false {
		t.Error("Expected explicit SetSelect(false) to override auto-logic")
	}
	if len(opts.SelectButtons) != 1 {
		t.Error("Expected SelectButtons to remain set even when Select is false")
	}
}

// TestSelectButtonsJSONExport tests that SelectButtons are properly exported in JSON
func TestSelectButtonsJSONExport(t *testing.T) {
	ctx := testOptionContext()
	builder := NewBuilder[testOptionRow](ctx, testOptionTranslator)

	btn1 := button.NewTableButton(core.ButtonActionLink, "edit", xurl.NewUrl("/edit"), "Edit tooltip", core.ColorPrimary, false, nil)
	btn2 := button.NewTableButton(core.ButtonActionDialog, "delete", xurl.NewUrl("/delete"), "Delete tooltip", core.ColorWarning, false, nil)

	builder.IdField("id", "test.id", func(r testOptionRow) int64 { return r.ID })
	builder.AddSelectButton(btn1)
	builder.AddSelectButton(btn2)

	tbl := builder.Build()
	output := tbl.Print(testOptionTranslator)

	data, ok := output["data"].(map[string]any)
	if !ok {
		t.Fatal("Expected data to be map[string]any")
	}

	options, ok := data["options"].(map[string]any)
	if !ok {
		t.Fatal("Expected options to be map[string]any")
	}

	// Check that selectButtons are exported
	selectButtons, ok := options["selectButtons"].([]map[string]any)
	if !ok {
		t.Fatal("Expected selectButtons to be []map[string]any")
	}

	if len(selectButtons) != 2 {
		t.Errorf("Expected 2 select buttons in JSON, got %d", len(selectButtons))
	}

	// Check that select is true
	selectValue, ok := options["select"].(bool)
	if !ok || !selectValue {
		t.Error("Expected select to be true in JSON")
	}
}

// TestAddMultiEditButton tests the AddMultiEditButton helper
func TestAddMultiEditButton(t *testing.T) {
	ctx := testOptionContext()
	builder := NewBuilder[testOptionRow](ctx, testOptionTranslator)

	builder.AddMultiEditButton("/Portal/Device/MultiEdit")

	tbl := builder.Build()
	opts := tbl.GetOptions()

	// Check that Select is auto-enabled
	if opts.Select == nil || *opts.Select != true {
		t.Error("Expected Select to be auto-enabled by AddMultiEditButton")
	}

	// Check that one button was added
	if len(opts.SelectButtons) != 1 {
		t.Errorf("Expected 1 select button, got %d", len(opts.SelectButtons))
	}

	// Verify button properties by checking JSON output
	output := tbl.Print(testOptionTranslator)
	data, _ := output["data"].(map[string]any)
	options, _ := data["options"].(map[string]any)
	selectButtons, _ := options["selectButtons"].([]map[string]any)

	if len(selectButtons) != 1 {
		t.Fatal("Expected 1 button in JSON output")
	}

	btn := selectButtons[0]
	if btn["icon"] != "edit" {
		t.Errorf("Expected icon 'edit', got %v", btn["icon"])
	}
	if btn["color"] != "primary" {
		t.Errorf("Expected color 'primary', got %v", btn["color"])
	}
	if btn["action"] != "dialog" {
		t.Errorf("Expected action 'dialog', got %v", btn["action"])
	}
}

// TestAddMultiDeleteButton tests the AddMultiDeleteButton helper
func TestAddMultiDeleteButton(t *testing.T) {
	ctx := testOptionContext()
	builder := NewBuilder[testOptionRow](ctx, testOptionTranslator)

	builder.AddMultiDeleteButton("/Portal/Device/MultiDel")

	tbl := builder.Build()
	opts := tbl.GetOptions()

	// Check that Select is auto-enabled
	if opts.Select == nil || *opts.Select != true {
		t.Error("Expected Select to be auto-enabled by AddMultiDeleteButton")
	}

	// Check that one button was added
	if len(opts.SelectButtons) != 1 {
		t.Errorf("Expected 1 select button, got %d", len(opts.SelectButtons))
	}

	// Verify button properties by checking JSON output
	output := tbl.Print(testOptionTranslator)
	data, _ := output["data"].(map[string]any)
	options, _ := data["options"].(map[string]any)
	selectButtons, _ := options["selectButtons"].([]map[string]any)

	if len(selectButtons) != 1 {
		t.Fatal("Expected 1 button in JSON output")
	}

	btn := selectButtons[0]
	if btn["icon"] != "delete" {
		t.Errorf("Expected icon 'delete', got %v", btn["icon"])
	}
	if btn["color"] != "warn" {
		t.Errorf("Expected color 'warn', got %v", btn["color"])
	}
	if btn["action"] != "dialog" {
		t.Errorf("Expected action 'dialog', got %v", btn["action"])
	}
}

// TestAddMultiEditAndDeleteButtons tests the combined helper
func TestAddMultiEditAndDeleteButtons(t *testing.T) {
	ctx := testOptionContext()
	builder := NewBuilder[testOptionRow](ctx, testOptionTranslator)

	builder.AddMultiEditAndDeleteButtons("/Portal/Device/MultiEdit", "/Portal/Device/MultiDel")

	tbl := builder.Build()
	opts := tbl.GetOptions()

	// Check that Select is auto-enabled
	if opts.Select == nil || *opts.Select != true {
		t.Error("Expected Select to be auto-enabled by AddMultiEditAndDeleteButtons")
	}

	// Check that two buttons were added
	if len(opts.SelectButtons) != 2 {
		t.Errorf("Expected 2 select buttons, got %d", len(opts.SelectButtons))
	}

	// Verify button properties by checking JSON output
	output := tbl.Print(testOptionTranslator)
	data, _ := output["data"].(map[string]any)
	options, _ := data["options"].(map[string]any)
	selectButtons, _ := options["selectButtons"].([]map[string]any)

	if len(selectButtons) != 2 {
		t.Fatal("Expected 2 buttons in JSON output")
	}

	// First button should be edit
	editBtn := selectButtons[0]
	if editBtn["icon"] != "edit" {
		t.Errorf("Expected first button icon 'edit', got %v", editBtn["icon"])
	}
	if editBtn["color"] != "primary" {
		t.Errorf("Expected first button color 'primary', got %v", editBtn["color"])
	}

	// Second button should be delete
	deleteBtn := selectButtons[1]
	if deleteBtn["icon"] != "delete" {
		t.Errorf("Expected second button icon 'delete', got %v", deleteBtn["icon"])
	}
	if deleteBtn["color"] != "warn" {
		t.Errorf("Expected second button color 'warn', got %v", deleteBtn["color"])
	}
}

// TestMultiButtonMethodChaining tests that helpers can be chained with other methods
func TestMultiButtonMethodChaining(t *testing.T) {
	ctx := testOptionContext()
	builder := NewBuilder[testOptionRow](ctx, testOptionTranslator)

	builder.
		SetDense(true).
		SetPagination(true).
		AddMultiEditAndDeleteButtons("/Portal/Device/MultiEdit", "/Portal/Device/MultiDel").
		SetItemsPerPage(50)

	tbl := builder.Build()
	opts := tbl.GetOptions()

	// Verify all options were set
	if opts.Dense == nil || *opts.Dense != true {
		t.Error("Expected Dense to be true")
	}
	if opts.Pagination == nil || *opts.Pagination != true {
		t.Error("Expected Pagination to be true")
	}
	if opts.ItemsPerPage == nil || *opts.ItemsPerPage != 50 {
		t.Error("Expected ItemsPerPage to be 50")
	}
	if len(opts.SelectButtons) != 2 {
		t.Errorf("Expected 2 select buttons, got %d", len(opts.SelectButtons))
	}
}
