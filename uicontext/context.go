// Package uicontext provides the UiContext struct with user preferences for UI rendering.
package uicontext

import (
	"github.com/xiriframework/xiri-go/types/distance"
	"github.com/xiriframework/xiri-go/types/language"
	"github.com/xiriframework/xiri-go/types/locale"
	"github.com/xiriframework/xiri-go/types/pressure"
	"github.com/xiriframework/xiri-go/types/timezone"
)

// UiContext contains user preferences for UI rendering, translation, and formatting.
// This is a slim, project-independent context without database or access control dependencies.
// Project-specific extensions (UserID, GroupID, DeviceContext, etc.) should embed this struct.
type UiContext struct {
	Timezone  timezone.Timezone
	Lang      language.Language
	Locale    locale.Locale
	Distance  distance.Distance
	Pressure  pressure.Pressure
	Translate func(key string) string // Injected per-project translation function
}

// SafeTranslate returns the translated string for the given key.
// It is nil-safe: if the UiContext or its Translate function is nil,
// it returns the key unchanged.
func (uc *UiContext) SafeTranslate(key string) string {
	if uc != nil && uc.Translate != nil {
		return uc.Translate(key)
	}
	return key
}
