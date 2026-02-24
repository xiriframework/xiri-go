package locale

// Locale represents a locale for formatting preferences (dates, numbers, etc.)
// It's a type-safe enum backed by int for database compatibility
type Locale int

// Locale constants - supported locales
// Supported locale constants
const (
	De   Locale = 0  // German (Germany)
	EnGB Locale = 1  // English (United Kingdom)
	Hr   Locale = 2  // Croatian
	Es   Locale = 3  // Spanish
	Fr   Locale = 4  // French
	It   Locale = 5  // Italian
	EnUS Locale = 6  // English (United States)
	DeAT Locale = 7  // German (Austria)
	DeCH Locale = 8  // German (Switzerland)
	Pt   Locale = 9  // Portuguese (Portugal)
	PtBR Locale = 10 // Portuguese (Brazil)
	Nl   Locale = 11 // Dutch (Netherlands)
	Pl   Locale = 12 // Polish (Poland)
	Cs   Locale = 13 // Czech (Czech Republic)
	Hu   Locale = 14 // Hungarian (Hungary)
	Ro   Locale = 15 // Romanian (Romania)
	Tr   Locale = 16 // Turkish (Turkey)
	Sv   Locale = 17 // Swedish (Sweden)
	Bg   Locale = 18 // Bulgarian (Bulgaria)
	Sl   Locale = 19 // Slovenian (Slovenia)
	Sk   Locale = 20 // Slovak (Slovakia)
	Sr   Locale = 21 // Serbian (Serbia)
	El   Locale = 22 // Greek (Greece)
	Nb   Locale = 23 // Norwegian Bokmål (Norway)
	Da   Locale = 24 // Danish (Denmark)
	Fi   Locale = 25 // Finnish (Finland)
	Ru   Locale = 26 // Russian (Russia)
	Uk   Locale = 27 // Ukrainian (Ukraine)
	Ja   Locale = 28 // Japanese (Japan)
	ZhCN Locale = 29 // Chinese (China)
	ArAE Locale = 30 // Arabic (UAE)
)

// Names maps locale values to human-readable names for debugging and logging
var Names = map[Locale]string{
	De:   "De",
	EnGB: "EnGB",
	Hr:   "Hr",
	Es:   "Es",
	Fr:   "Fr",
	It:   "It",
	EnUS: "EnUS",
	DeAT: "DeAT",
	DeCH: "DeCH",
	Pt:   "Pt",
	PtBR: "PtBR",
	Nl:   "Nl",
	Pl:   "Pl",
	Cs:   "Cs",
	Hu:   "Hu",
	Ro:   "Ro",
	Tr:   "Tr",
	Sv:   "Sv",
	Bg:   "Bg",
	Sl:   "Sl",
	Sk:   "Sk",
	Sr:   "Sr",
	El:   "El",
	Nb:   "Nb",
	Da:   "Da",
	Fi:   "Fi",
	Ru:   "Ru",
	Uk:   "Uk",
	Ja:   "Ja",
	ZhCN: "ZhCN",
	ArAE: "ArAE",
}

// LocaleStrings maps locale values to standard locale strings (e.g., "de-DE", "en-GB")
var LocaleStrings = map[Locale]string{
	De:   "de-DE",
	EnGB: "en-GB",
	Hr:   "hr-HR",
	Es:   "es-ES",
	Fr:   "fr-FR",
	It:   "it-IT",
	EnUS: "en-US",
	DeAT: "de-AT",
	DeCH: "de-CH",
	Pt:   "pt-PT",
	PtBR: "pt-BR",
	Nl:   "nl-NL",
	Pl:   "pl-PL",
	Cs:   "cs-CZ",
	Hu:   "hu-HU",
	Ro:   "ro-RO",
	Tr:   "tr-TR",
	Sv:   "sv-SE",
	Bg:   "bg-BG",
	Sl:   "sl-SI",
	Sk:   "sk-SK",
	Sr:   "sr-RS",
	El:   "el-GR",
	Nb:   "nb-NO",
	Da:   "da-DK",
	Fi:   "fi-FI",
	Ru:   "ru-RU",
	Uk:   "uk-UA",
	Ja:   "ja-JP",
	ZhCN: "zh-CN",
	ArAE: "ar-AE",
}

// String returns the string representation of the locale value
func (l Locale) String() string {
	if name, ok := Names[l]; ok {
		return name
	}
	return "Unknown"
}

// GetName returns the human-readable name for a locale value
func GetName(l Locale) string {
	return l.String()
}

// GetLocaleString returns the standard locale string (e.g., "de-DE", "en-GB")
func (l Locale) GetLocaleString() string {
	if localeStr, ok := LocaleStrings[l]; ok {
		return localeStr
	}
	return ""
}

// IsValid checks if a locale value is valid
func IsValid(l Locale) bool {
	_, ok := Names[l]
	return ok
}

// ToInt32 converts the Locale to int32 for database storage
func (l Locale) ToInt32() int32 {
	return int32(l)
}

// FromInt32 converts an int32 to a Locale
func FromInt32(i int32) Locale {
	return Locale(i)
}

// FromLocaleString converts a locale string to a Locale
func FromLocaleString(localeStr string) (Locale, bool) {
	for locale, str := range LocaleStrings {
		if str == localeStr {
			return locale, true
		}
	}
	return 0, false
}
