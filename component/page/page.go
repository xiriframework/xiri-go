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

// Page represents a page container with breadcrumbs and components.
type Page struct {
	translator core.TranslateFunc
	bread      []BreadcrumbItem
	data       []map[string]any
	extra      map[string]any
}

// NewPage creates a new page container with the given translator.
func NewPage(translator core.TranslateFunc) *Page {
	return &Page{
		translator: translator,
		bread:      make([]BreadcrumbItem, 0),
		data:       make([]map[string]any, 0),
		extra:      make(map[string]any),
	}
}

// Add adds a component to the page.
func (p *Page) Add(component core.Component) *Page {
	printed := component.Print(p.translator)
	p.data = append(p.data, printed)
	return p
}

// AddNewRow adds a component to the page and forces it to start a new grid row.
func (p *Page) AddNewRow(component core.Component) *Page {
	printed := component.Print(p.translator)
	printed["newRow"] = true
	p.data = append(p.data, printed)
	return p
}

// AddOld adds a raw component map to the page (for backwards compatibility).
func (p *Page) AddOld(component map[string]any) *Page {
	p.data = append(p.data, component)
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
func (p *Page) Print(translator core.TranslateFunc) map[string]any {
	result := make(map[string]any)

	// Add breadcrumbs if any
	if len(p.bread) > 0 {
		breadData := make([]map[string]any, len(p.bread))
		for i, b := range p.bread {
			breadData[i] = b.print()
		}
		result["bread"] = breadData
	}

	// Add data
	result["data"] = p.data

	// Merge extra fields at root level
	for key, value := range p.extra {
		result[key] = value
	}

	return result
}
