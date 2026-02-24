package language

// Language represents a user interface language
// It's a type-safe enum backed by int for database compatibility
type Language int

// Language constants - supported UI languages
// Supported language constants
const (
	Deutsch         Language = 0  // German
	Englisch        Language = 1  // English
	Kroatisch       Language = 2  // Croatian
	Spanisch        Language = 3  // Spanish
	Franzoesisch    Language = 4  // French
	Italienisch     Language = 5  // Italian
	Portugiesisch   Language = 6  // Portuguese
	Niederlaendisch Language = 7  // Dutch
	Polnisch        Language = 8  // Polish
	Tschechisch     Language = 9  // Czech
	Ungarisch       Language = 10 // Hungarian
	Rumaenisch      Language = 11 // Romanian
	Tuerkisch       Language = 12 // Turkish
	Schwedisch      Language = 13 // Swedish
	Bulgarisch      Language = 14 // Bulgarian
	Slowenisch      Language = 15 // Slovenian
	Slowakisch      Language = 16 // Slovak
	Serbisch        Language = 17 // Serbian
	Griechisch      Language = 18 // Greek
	Norwegisch      Language = 19 // Norwegian
	Daenisch        Language = 20 // Danish
	Finnisch        Language = 21 // Finnish
	Russisch        Language = 22 // Russian
	Ukrainisch      Language = 23 // Ukrainian
	Japanisch       Language = 24 // Japanese
	Chinesisch      Language = 25 // Chinese
	Arabisch        Language = 26 // Arabic
)

// Names maps language values to human-readable names for debugging and logging
var Names = map[Language]string{
	Deutsch:         "Deutsch",
	Englisch:        "Englisch",
	Kroatisch:       "Kroatisch",
	Spanisch:        "Spanisch",
	Franzoesisch:    "Franzoesisch",
	Italienisch:     "Italienisch",
	Portugiesisch:   "Portugiesisch",
	Niederlaendisch: "Niederlaendisch",
	Polnisch:        "Polnisch",
	Tschechisch:     "Tschechisch",
	Ungarisch:       "Ungarisch",
	Rumaenisch:      "Rumaenisch",
	Tuerkisch:       "Tuerkisch",
	Schwedisch:      "Schwedisch",
	Bulgarisch:      "Bulgarisch",
	Slowenisch:      "Slowenisch",
	Slowakisch:      "Slowakisch",
	Serbisch:        "Serbisch",
	Griechisch:      "Griechisch",
	Norwegisch:      "Norwegisch",
	Daenisch:        "Daenisch",
	Finnisch:        "Finnisch",
	Russisch:        "Russisch",
	Ukrainisch:      "Ukrainisch",
	Japanisch:       "Japanisch",
	Chinesisch:      "Chinesisch",
	Arabisch:        "Arabisch",
}

// LanguageCodes maps language values to ISO 639-1 language codes
var LanguageCodes = map[Language]string{
	Deutsch:         "de",
	Englisch:        "en",
	Kroatisch:       "hr",
	Spanisch:        "es",
	Franzoesisch:    "fr",
	Italienisch:     "it",
	Portugiesisch:   "pt",
	Niederlaendisch: "nl",
	Polnisch:        "pl",
	Tschechisch:     "cs",
	Ungarisch:       "hu",
	Rumaenisch:      "ro",
	Tuerkisch:       "tr",
	Schwedisch:      "sv",
	Bulgarisch:      "bg",
	Slowenisch:      "sl",
	Slowakisch:      "sk",
	Serbisch:        "sr",
	Griechisch:      "el",
	Norwegisch:      "no",
	Daenisch:        "da",
	Finnisch:        "fi",
	Russisch:        "ru",
	Ukrainisch:      "uk",
	Japanisch:       "ja",
	Chinesisch:      "zh",
	Arabisch:        "ar",
}

// String returns the string representation of the language value
func (l Language) String() string {
	if name, ok := Names[l]; ok {
		return name
	}
	return "Unknown"
}

// GetName returns the human-readable name for a language value
func GetName(l Language) string {
	return l.String()
}

// GetCode returns the ISO 639-1 language code for a language value
func (l Language) GetCode() string {
	if code, ok := LanguageCodes[l]; ok {
		return code
	}
	return ""
}

// IsValid checks if a language value is valid
func IsValid(l Language) bool {
	_, ok := Names[l]
	return ok
}

// ToInt32 converts the Language to int32 for database storage
func (l Language) ToInt32() int32 {
	return int32(l)
}

// FromInt32 converts an int32 to a Language
func FromInt32(i int32) Language {
	return Language(i)
}

// FromCode converts a language code to a Language
func FromCode(code string) (Language, bool) {
	for lang, c := range LanguageCodes {
		if c == code {
			return lang, true
		}
	}
	return 0, false
}
