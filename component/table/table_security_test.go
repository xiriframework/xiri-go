package table

import (
	"strings"
	"testing"
)

// buildSecurityTestTable creates a table with fields "id", "name", "status"
// for use in pagination/sort security tests.
func buildSecurityTestTable() *Table[testDeviceRow] {
	builder := NewBuilder[testDeviceRow]()
	builder.IdField("id", "device.id", func(r testDeviceRow) int64 { return r.ID })
	builder.TextField("name", "device.name", func(r testDeviceRow) string { return r.Name })
	builder.BoolField("status", "device.active", func(r testDeviceRow) bool { return r.Active })
	return builder.Build()
}

// TestLoadPaginationParams_SortWhitelist verifies that sort column is validated
// against defined field IDs to prevent SQL injection (H1).
func TestLoadPaginationParams_SortWhitelist(t *testing.T) {
	tests := []struct {
		name     string
		sort     string
		expected string
	}{
		{"valid field", "name", "name"},
		{"valid field id", "id", "id"},
		{"invalid field", "nonexistent", ""},
		{"SQL injection", "name; DROP TABLE users--", ""},
		{"partial match", "nam", ""},
		{"unicode", "名前", ""},
		{"empty", "", ""},
		{"whitespace", " name ", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tbl := buildSecurityTestTable()
			tbl.SetFilterData(map[string]any{
				"_sort": tt.sort,
			})
			params := tbl.LoadPaginationParams()
			if params.Sort != tt.expected {
				t.Errorf("sort=%q: expected %q, got %q", tt.sort, tt.expected, params.Sort)
			}
		})
	}
}

// TestLoadPaginationParams_PageSizeClamping verifies that pageSize is clamped
// to [1, 1000] to prevent memory exhaustion (M2).
func TestLoadPaginationParams_PageSizeClamping(t *testing.T) {
	tests := []struct {
		name     string
		pageSize float64
		expected int
	}{
		{"zero", 0, 50},
		{"negative", -10, 50},
		{"max allowed", 1000, 1000},
		{"above max", 1001, 1000},
		{"way above max", 999999, 1000},
		{"normal", 25, 25},
		{"minimum valid", 1, 1},
		{"large valid", 500, 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tbl := buildSecurityTestTable()
			tbl.SetFilterData(map[string]any{
				"_pageSize": tt.pageSize,
			})
			params := tbl.LoadPaginationParams()
			if params.PageSize != tt.expected {
				t.Errorf("pageSize=%v: expected %d, got %d", tt.pageSize, tt.expected, params.PageSize)
			}
		})
	}
}

// TestLoadPaginationParams_SearchMaxLength verifies that search string is
// truncated to 200 characters to prevent performance issues (M6).
func TestLoadPaginationParams_SearchMaxLength(t *testing.T) {
	tests := []struct {
		name           string
		searchLen      int
		expectedMaxLen int
	}{
		{"short", 10, 10},
		{"at limit", 200, 200},
		{"over limit", 201, 200},
		{"way over limit", 5000, 200},
		{"empty", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			search := strings.Repeat("x", tt.searchLen)
			tbl := buildSecurityTestTable()
			tbl.SetFilterData(map[string]any{
				"_search": search,
			})
			params := tbl.LoadPaginationParams()
			if len(params.Search) != tt.expectedMaxLen {
				t.Errorf("searchLen=%d: expected len %d, got %d", tt.searchLen, tt.expectedMaxLen, len(params.Search))
			}
		})
	}
}

// TestLoadPaginationParams_PageNonNegative verifies that page index cannot be
// negative to prevent unexpected query behavior (M7).
func TestLoadPaginationParams_PageNonNegative(t *testing.T) {
	tests := []struct {
		name     string
		page     float64
		expected int
	}{
		{"negative", -1, 0},
		{"very negative", -9999, 0},
		{"zero", 0, 0},
		{"positive", 5, 5},
		{"large", 100, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tbl := buildSecurityTestTable()
			tbl.SetFilterData(map[string]any{
				"_page": tt.page,
			})
			params := tbl.LoadPaginationParams()
			if params.Page != tt.expected {
				t.Errorf("page=%v: expected %d, got %d", tt.page, tt.expected, params.Page)
			}
		})
	}
}

// TestLoadPaginationParams_SortDirValidation verifies that sortDir only accepts
// "asc" or "desc" and defaults to "asc" for invalid values.
func TestLoadPaginationParams_SortDirValidation(t *testing.T) {
	tests := []struct {
		name     string
		sortDir  string
		expected string
	}{
		{"asc", "asc", "asc"},
		{"desc", "desc", "desc"},
		{"invalid", "INVALID", "asc"},
		{"SQL injection", "asc; DROP TABLE users--", "asc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tbl := buildSecurityTestTable()
			tbl.SetFilterData(map[string]any{
				"_sortDir": tt.sortDir,
			})
			params := tbl.LoadPaginationParams()
			if params.SortDir != tt.expected {
				t.Errorf("sortDir=%q: expected %q, got %q", tt.sortDir, tt.expected, params.SortDir)
			}
		})
	}
}
