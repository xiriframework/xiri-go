package table

import (
	"log/slog"

	"github.com/xiriframework/xiri-go/component/core"
)

// IconRef is an opaque reference to a registered icon definition.
// It can only be created by IconSet.Add(), ensuring that only
// valid, registered icon values can be used in IconFieldFromSet accessors.
type IconRef struct {
	value string // unexported: can only be set within this package
}

// IconSet holds icon definitions and produces IconRef values.
// Use NewIconSet() to create an IconSet, then Add() to register
// icon definitions. The resulting IconRefs are used in IconFieldFromSet accessors.
type IconSet struct {
	icons map[string]*iconDef
	order []string
}

// NewIconSet creates a new empty IconSet.
func NewIconSet() *IconSet {
	return &IconSet{
		icons: make(map[string]*iconDef),
	}
}

// Add registers an icon definition and returns an IconRef for use in accessors.
//
// Parameters:
//   - value: The string value that identifies this icon (used as key in JSON)
//   - icon: Material icon name (e.g., "check_circle", "cancel")
//   - color: Icon color theme
//   - hint: Tooltip text (translation key)
//
// Example:
//
//	icons := table.NewIconSet()
//	iconOnline  := icons.Add("online",  "check_circle", core.ColorAccent,  "Online")
//	iconOffline := icons.Add("offline", "cancel",       core.ColorWarning, "Offline")
func (s *IconSet) Add(value, icon string, color core.Color, hint string) *IconRef {
	s.icons[value] = &iconDef{
		icon:    icon,
		color:   color,
		hint:    hint,
		options: make(map[string]any),
	}
	s.order = append(s.order, value)
	return &IconRef{value: value}
}

// AddWithOptions registers an icon definition with additional custom options.
func (s *IconSet) AddWithOptions(value, icon string, color core.Color, hint string, opts map[string]any) *IconRef {
	s.icons[value] = &iconDef{
		icon:    icon,
		color:   color,
		hint:    hint,
		options: opts,
	}
	s.order = append(s.order, value)
	return &IconRef{value: value}
}

// Resolve returns the IconRef for a given string value, or nil if the value
// is not registered in this IconSet. Use this when the source data is a string
// (e.g., from a database) rather than a compile-time constant.
//
// Example:
//
//	builder.IconFieldFromSet("status", "device.status",
//	    func(r Row) *table.IconRef {
//	        return statusIcons.Resolve(r.Status) // nil for unknown values
//	    },
//	    statusIcons,
//	)
func (s *IconSet) Resolve(value string) *IconRef {
	if _, ok := s.icons[value]; ok {
		return &IconRef{value: value}
	}
	slog.Warn("table: IconSet.Resolve: unknown value", "value", value)
	return nil
}

// Len returns the number of icon definitions in the set.
func (s *IconSet) Len() int {
	return len(s.icons)
}
