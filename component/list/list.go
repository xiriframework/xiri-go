// Package list provides list display components for the Angular frontend.
package list

import (
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
	"github.com/xiriframework/xiri-go/response"
)

// List represents a list component with sections
type List struct {
	sections    []*ListSection
	rawSections []map[string]any
	display     *string
	url         *url.Url
	reload      *bool
}

// NewList creates a new list component from pre-rendered section maps
func NewList(sections []map[string]any, display *string) *List {
	if sections == nil {
		sections = make([]map[string]any, 0)
	}
	return &List{
		rawSections: sections,
		display:     display,
	}
}

// AddSection adds a section to the list. Rendering is deferred to Print().
func (l *List) AddSection(section *ListSection) *List {
	l.sections = append(l.sections, section)
	return l
}

// WithDisplay sets the display/layout class (optional)
func (l *List) WithDisplay(display string) *List {
	l.display = &display
	return l
}

// SetURL sets the AJAX data URL. When set, sections are cleared and the frontend loads data dynamically.
func (l *List) SetURL(u *url.Url) *List {
	l.url = u
	l.sections = nil
	l.rawSections = nil
	return l
}

// WithReload enables periodic reload of the list data when using AJAX mode.
func (l *List) WithReload(reload bool) *List {
	l.reload = &reload
	return l
}

// Print returns the JSON representation of the list
func (l *List) Print(translator core.TranslateFunc) map[string]any {
	if l.url != nil {
		return map[string]any{
			"type":    "list",
			"display": l.display,
			"data": map[string]any{
				"url":    l.url.PrintPrefix(),
				"reload": l.reload,
			},
		}
	}

	return map[string]any{
		"type":    "list",
		"display": l.display,
		"data":    l.printData(translator),
	}
}

// PrintData returns only the data portion of the list (for use in data endpoints).
func (l *List) PrintData(translator core.TranslateFunc) map[string]any {
	return l.printData(translator)
}

// DataResponse returns a DataResult wrapping the list data in {"data": ...} envelope.
func (l *List) DataResponse(translator core.TranslateFunc) response.DataResult {
	return response.NewJSONDataResult(l.PrintData(translator))
}

// printData builds the data map used by both Print and PrintData.
func (l *List) printData(translator core.TranslateFunc) map[string]any {
	allSections := make([]map[string]any, 0, len(l.rawSections)+len(l.sections))
	allSections = append(allSections, l.rawSections...)
	for _, s := range l.sections {
		allSections = append(allSections, s.Print(translator))
	}

	return map[string]any{
		"sections": allSections,
	}
}

// ListSection represents a section in a list
type ListSection struct {
	name    *string
	items   []*ListSectionItem
	rawData []map[string]any
}

// NewListSection creates a new list section
func NewListSection(name *string, data []map[string]any) *ListSection {
	if data == nil {
		data = make([]map[string]any, 0)
	}
	return &ListSection{
		name:    name,
		rawData: data,
	}
}

// AddItem adds an item to the section. Rendering is deferred to Print().
func (ls *ListSection) AddItem(item *ListSectionItem) *ListSection {
	ls.items = append(ls.items, item)
	return ls
}

// WithName sets the section name/header (optional)
func (ls *ListSection) WithName(name string) *ListSection {
	ls.name = &name
	return ls
}

// Print returns the JSON representation of the list section
func (ls *ListSection) Print(translator core.TranslateFunc) map[string]any {
	allData := make([]map[string]any, 0, len(ls.rawData)+len(ls.items))
	allData = append(allData, ls.rawData...)
	for _, item := range ls.items {
		allData = append(allData, item.Print(translator))
	}

	return map[string]any{
		"name": ls.name,
		"data": allData,
	}
}

// ListSectionItem represents an item in a list section
type ListSectionItem struct {
	name         string
	info         string
	url          *url.Url
	icon         string
	iconColor    core.ListItemColor
	iconSet      *string
	hasFavorite  bool
	isFavorite   bool
	favoriteUrl  *url.Url
	favoriteHint *string
}

// NewListSectionItem creates a new list section item
func NewListSectionItem(
	name string,
	info string,
	u *url.Url,
	icon string,
	iconColor core.ListItemColor,
	iconSet *string,
	hasFavorite bool,
	isFavorite bool,
	favoriteUrl *url.Url,
	favoriteHint *string,
) *ListSectionItem {
	return &ListSectionItem{
		name:         name,
		info:         info,
		url:          u,
		icon:         icon,
		iconColor:    iconColor,
		iconSet:      iconSet,
		hasFavorite:  hasFavorite,
		isFavorite:   isFavorite,
		favoriteUrl:  favoriteUrl,
		favoriteHint: favoriteHint,
	}
}

// NewSimpleListSectionItem creates a list section item without favorite support
func NewSimpleListSectionItem(name, info string, u *url.Url, icon string, iconColor core.ListItemColor) *ListSectionItem {
	return &ListSectionItem{
		name:      name,
		info:      info,
		url:       u,
		icon:      icon,
		iconColor: iconColor,
	}
}

// WithIconSet sets the icon set (optional)
func (lsi *ListSectionItem) WithIconSet(iconSet string) *ListSectionItem {
	lsi.iconSet = &iconSet
	return lsi
}

// WithFavorite configures favorite functionality (optional)
func (lsi *ListSectionItem) WithFavorite(isFavorite bool, favoriteUrl *url.Url, favoriteHint string) *ListSectionItem {
	lsi.hasFavorite = true
	lsi.isFavorite = isFavorite
	lsi.favoriteUrl = favoriteUrl
	lsi.favoriteHint = &favoriteHint
	return lsi
}

// WithIsFavorite sets whether this item is currently a favorite
func (lsi *ListSectionItem) WithIsFavorite(isFavorite bool) *ListSectionItem {
	lsi.isFavorite = isFavorite
	return lsi
}

// Print returns the JSON representation of the list section item
func (lsi *ListSectionItem) Print(translator core.TranslateFunc) map[string]any {
	result := map[string]any{
		"name":         lsi.name,
		"url":          lsi.url.PrintPrefix(),
		"info":         lsi.info,
		"icon":         lsi.icon,
		"iconColor":    string(lsi.iconColor),
		"iconSet":      lsi.iconSet,
		"hasFavorite":  lsi.hasFavorite,
		"isFavorite":   lsi.isFavorite,
		"favoriteUrl":  nil,
		"favoriteHint": lsi.favoriteHint,
	}

	if lsi.favoriteUrl != nil {
		result["favoriteUrl"] = lsi.favoriteUrl.PrintPrefix()
	}

	return result
}
