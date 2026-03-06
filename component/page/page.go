// Package page provides the Page container component which holds breadcrumbs and nested components.
package page

import (
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
)

// BreadcrumbItem represents a single breadcrumb navigation entry.
type BreadcrumbItem struct {
	name   string
	link   *string
	extern bool
}

// NewBreadcrumbItem creates a new breadcrumb item with name, optional link, and extern flag.
func NewBreadcrumbItem(name string, link *string, extern bool) BreadcrumbItem {
	return BreadcrumbItem{name: name, link: link, extern: extern}
}

// print returns the JSON-compatible map for a breadcrumb item.
func (b BreadcrumbItem) print() map[string]any {
	return map[string]any{
		"label":  b.name,
		"link":   b.link,
		"extern": b.extern,
	}
}

// pageEntry stores either a lazy component or a raw map for later rendering.
type pageEntry struct {
	component core.Component
	raw       map[string]any
	newRow    bool
}

// Page represents a page container with breadcrumbs and components.
type Page struct {
	bread   []BreadcrumbItem
	entries []pageEntry
	extra   map[string]any
}

// NewPage creates a new page container.
func NewPage() *Page {
	return &Page{
		bread:   make([]BreadcrumbItem, 0),
		entries: make([]pageEntry, 0),
		extra:   make(map[string]any),
	}
}

// Add adds a component to the page. The component is rendered lazily when Print is called.
func (p *Page) Add(component core.Component) *Page {
	p.entries = append(p.entries, pageEntry{component: component})
	return p
}

// AddNewRow adds a component to the page and forces it to start a new grid row.
// The component is rendered lazily when Print is called.
func (p *Page) AddNewRow(component core.Component) *Page {
	p.entries = append(p.entries, pageEntry{component: component, newRow: true})
	return p
}

// AddOld adds a raw component map to the page (for backwards compatibility).
func (p *Page) AddOld(component map[string]any) *Page {
	p.entries = append(p.entries, pageEntry{raw: component})
	return p
}

// Extra adds an extra field to the page root.
func (p *Page) Extra(name string, data any) *Page {
	p.extra[name] = data
	return p
}

// Bread adds a breadcrumb item to the page.
func (p *Page) Bread(name string, u *url.Url, extern bool) *Page {
	item := BreadcrumbItem{
		name:   name,
		link:   nil,
		extern: extern,
	}
	if u != nil {
		urlStr := u.PrintPrefix()
		item.link = &urlStr
	}
	p.bread = append(p.bread, item)
	return p
}

// Print returns the JSON representation of the page.
// Components added via Add/AddNewRow are rendered lazily here with the given ctx.
func (p *Page) Print(ctx *core.UiContext) map[string]any {
	result := make(map[string]any)

	// Add breadcrumbs if any
	if len(p.bread) > 0 {
		breadData := make([]map[string]any, len(p.bread))
		for i, b := range p.bread {
			breadData[i] = b.print()
		}
		result["bread"] = breadData
	}

	// Render entries lazily
	data := make([]map[string]any, 0, len(p.entries))
	for _, e := range p.entries {
		var printed map[string]any
		if e.component != nil {
			printed = e.component.Print(ctx)
		} else {
			printed = e.raw
		}
		if e.newRow {
			printed["newRow"] = true
		}
		data = append(data, printed)
	}
	result["data"] = data

	// Merge extra fields at root level
	for key, value := range p.extra {
		result[key] = value
	}

	return result
}
