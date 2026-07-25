package core

import (
	"slices"
	"testing"

	"github.com/xiriframework/xiri-go/types/distance"
	"github.com/xiriframework/xiri-go/types/language"
	"github.com/xiriframework/xiri-go/types/locale"
	"github.com/xiriframework/xiri-go/types/pressure"
	"github.com/xiriframework/xiri-go/types/timezone"
)

// This file tests the public API of the types/* enum packages from a consumer's
// point of view — core is the package that imports all five. The lookup tables
// behind these enums are unexported (#16): a writable package-level map lets any
// consumer corrupt them at runtime, so enumeration goes through All() instead,
// which must hand out a fresh slice every call.

// checkEnum asserts All() is complete, stably ordered and isolated from the
// package's internal table.
func checkEnum[T ~int](t *testing.T, pkg string, all func() []T, name func(T) string, valid func(T) bool) {
	t.Helper()

	values := all()
	if len(values) == 0 {
		t.Fatalf("%s: All() returned no values", pkg)
	}

	if !slices.IsSorted(values) {
		t.Errorf("%s: All() must be ordered by numeric value, got %v", pkg, values)
	}

	for _, v := range values {
		if !valid(v) {
			t.Errorf("%s: All() returned %v, which IsValid rejects", pkg, v)
		}
		if n := name(v); n == "" || n == "Unknown" {
			t.Errorf("%s: value %v has no name (%q)", pkg, v, n)
		}
	}

	// A caller must not be able to affect the package by writing to the result.
	original := slices.Clone(values)
	for i := range values {
		values[i] = T(-1)
	}
	if again := all(); !slices.Equal(again, original) {
		t.Errorf("%s: All() is not isolated from its caller: expected %v, got %v", pkg, original, again)
	}
}

func TestEnumAll(t *testing.T) {
	checkEnum(t, "timezone", timezone.All, timezone.GetName, timezone.IsValid)
	checkEnum(t, "locale", locale.All, locale.GetName, locale.IsValid)
	checkEnum(t, "language", language.All, language.GetName, language.IsValid)
	checkEnum(t, "distance", distance.All, distance.GetName, distance.IsValid)
	checkEnum(t, "pressure", pressure.All, pressure.GetName, pressure.IsValid)
}

// TestEnumAccessorsCoverLookupTables verifies the per-value accessors return
// something for every enumerated value, so consumers never need the raw tables.
func TestEnumAccessorsCoverLookupTables(t *testing.T) {
	for _, tz := range timezone.All() {
		if tz.GetIANA() == "" {
			t.Errorf("timezone %v has no IANA string", tz)
		}
	}
	for _, l := range locale.All() {
		if l.GetLocaleString() == "" {
			t.Errorf("locale %v has no locale string", l)
		}
	}
	for _, l := range language.All() {
		if l.GetCode() == "" {
			t.Errorf("language %v has no code", l)
		}
	}
	for _, d := range distance.All() {
		if d.GetSymbol() == "" {
			t.Errorf("distance %v has no symbol", d)
		}
	}
	for _, p := range pressure.All() {
		if p.GetSymbol() == "" {
			t.Errorf("pressure %v has no symbol", p)
		}
	}
}
