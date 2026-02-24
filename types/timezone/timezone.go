package timezone

// Timezone represents a timezone for user preferences
// It's a type-safe enum backed by int for database compatibility
type Timezone int

// Timezone constants - supported timezones
// Supported timezone constants
const (
	EuropeVienna      Timezone = 0  // Europe/Vienna (Austria)
	EuropeBerlin      Timezone = 1  // Europe/Berlin (Germany)
	EuropeZagreb      Timezone = 2  // Europe/Zagreb (Croatia)
	EuropeMadrid      Timezone = 3  // Europe/Madrid (Spain)
	EuropeRome        Timezone = 4  // Europe/Rome (Italy)
	EuropeParis       Timezone = 5  // Europe/Paris (France)
	EuropeLondon      Timezone = 6  // Europe/London (UK)
	EuropeBrussels    Timezone = 7  // Europe/Brussels (Belgium)
	EuropeLisbon      Timezone = 8  // Europe/Lisbon (Portugal)
	EuropeAmsterdam   Timezone = 9  // Europe/Amsterdam (Netherlands)
	EuropeBucharest   Timezone = 10 // Europe/Bucharest (Romania)
	EuropeWarsaw      Timezone = 11 // Europe/Warsaw (Poland)
	EuropeHelsinki    Timezone = 12 // Europe/Helsinki (Finland)
	EuropeAthens      Timezone = 13 // Europe/Athens (Greece)
	EuropePrague      Timezone = 14 // Europe/Prague (Czech Republic)
	EuropeBudapest    Timezone = 15 // Europe/Budapest (Hungary)
	EuropeStockholm   Timezone = 16 // Europe/Stockholm (Sweden)
	EuropeCopenhagen  Timezone = 17 // Europe/Copenhagen (Denmark)
	EuropeOslo        Timezone = 18 // Europe/Oslo (Norway)
	EuropeIstanbul    Timezone = 19 // Europe/Istanbul (Turkey)
	EuropeDublin      Timezone = 20 // Europe/Dublin (Ireland)
	EuropeMoscow      Timezone = 21 // Europe/Moscow (Russia)
	EuropeKyiv        Timezone = 22 // Europe/Kyiv (Ukraine)
	EuropeSofia       Timezone = 23 // Europe/Sofia (Bulgaria)
	EuropeLjubljana   Timezone = 24 // Europe/Ljubljana (Slovenia)
	EuropeBratislava  Timezone = 25 // Europe/Bratislava (Slovakia)
	EuropeBelgrade    Timezone = 26 // Europe/Belgrade (Serbia)
	AmericaNewYork    Timezone = 27 // America/New_York (US Eastern)
	AmericaChicago    Timezone = 28 // America/Chicago (US Central)
	AmericaDenver     Timezone = 29 // America/Denver (US Mountain)
	AmericaLosAngeles Timezone = 30 // America/Los_Angeles (US Pacific)
	AsiaTokyo         Timezone = 31 // Asia/Tokyo (Japan)
	AsiaShanghai      Timezone = 32 // Asia/Shanghai (China)
	AsiaDubai         Timezone = 33 // Asia/Dubai (UAE)
	AustraliaSydney   Timezone = 34 // Australia/Sydney (Australia)
	UTC               Timezone = 35 // UTC
)

// Names maps timezone values to human-readable names for debugging and logging
var Names = map[Timezone]string{
	EuropeVienna:      "Europe/Vienna",
	EuropeBerlin:      "Europe/Berlin",
	EuropeZagreb:      "Europe/Zagreb",
	EuropeMadrid:      "Europe/Madrid",
	EuropeRome:        "Europe/Rome",
	EuropeParis:       "Europe/Paris",
	EuropeLondon:      "Europe/London",
	EuropeBrussels:    "Europe/Brussels",
	EuropeLisbon:      "Europe/Lisbon",
	EuropeAmsterdam:   "Europe/Amsterdam",
	EuropeBucharest:   "Europe/Bucharest",
	EuropeWarsaw:      "Europe/Warsaw",
	EuropeHelsinki:    "Europe/Helsinki",
	EuropeAthens:      "Europe/Athens",
	EuropePrague:      "Europe/Prague",
	EuropeBudapest:    "Europe/Budapest",
	EuropeStockholm:   "Europe/Stockholm",
	EuropeCopenhagen:  "Europe/Copenhagen",
	EuropeOslo:        "Europe/Oslo",
	EuropeIstanbul:    "Europe/Istanbul",
	EuropeDublin:      "Europe/Dublin",
	EuropeMoscow:      "Europe/Moscow",
	EuropeKyiv:        "Europe/Kyiv",
	EuropeSofia:       "Europe/Sofia",
	EuropeLjubljana:   "Europe/Ljubljana",
	EuropeBratislava:  "Europe/Bratislava",
	EuropeBelgrade:    "Europe/Belgrade",
	AmericaNewYork:    "America/New_York",
	AmericaChicago:    "America/Chicago",
	AmericaDenver:     "America/Denver",
	AmericaLosAngeles: "America/Los_Angeles",
	AsiaTokyo:         "Asia/Tokyo",
	AsiaShanghai:      "Asia/Shanghai",
	AsiaDubai:         "Asia/Dubai",
	AustraliaSydney:   "Australia/Sydney",
	UTC:               "UTC",
}

// TimezoneStrings maps timezone values to IANA timezone strings
var TimezoneStrings = map[Timezone]string{
	EuropeVienna:      "Europe/Vienna",
	EuropeBerlin:      "Europe/Berlin",
	EuropeZagreb:      "Europe/Zagreb",
	EuropeMadrid:      "Europe/Madrid",
	EuropeRome:        "Europe/Rome",
	EuropeParis:       "Europe/Paris",
	EuropeLondon:      "Europe/London",
	EuropeBrussels:    "Europe/Brussels",
	EuropeLisbon:      "Europe/Lisbon",
	EuropeAmsterdam:   "Europe/Amsterdam",
	EuropeBucharest:   "Europe/Bucharest",
	EuropeWarsaw:      "Europe/Warsaw",
	EuropeHelsinki:    "Europe/Helsinki",
	EuropeAthens:      "Europe/Athens",
	EuropePrague:      "Europe/Prague",
	EuropeBudapest:    "Europe/Budapest",
	EuropeStockholm:   "Europe/Stockholm",
	EuropeCopenhagen:  "Europe/Copenhagen",
	EuropeOslo:        "Europe/Oslo",
	EuropeIstanbul:    "Europe/Istanbul",
	EuropeDublin:      "Europe/Dublin",
	EuropeMoscow:      "Europe/Moscow",
	EuropeKyiv:        "Europe/Kyiv",
	EuropeSofia:       "Europe/Sofia",
	EuropeLjubljana:   "Europe/Ljubljana",
	EuropeBratislava:  "Europe/Bratislava",
	EuropeBelgrade:    "Europe/Belgrade",
	AmericaNewYork:    "America/New_York",
	AmericaChicago:    "America/Chicago",
	AmericaDenver:     "America/Denver",
	AmericaLosAngeles: "America/Los_Angeles",
	AsiaTokyo:         "Asia/Tokyo",
	AsiaShanghai:      "Asia/Shanghai",
	AsiaDubai:         "Asia/Dubai",
	AustraliaSydney:   "Australia/Sydney",
	UTC:               "UTC",
}

// String returns the string representation of the timezone value
func (tz Timezone) String() string {
	if name, ok := Names[tz]; ok {
		return name
	}
	return "Unknown"
}

// GetName returns the human-readable name for a timezone value
func GetName(tz Timezone) string {
	return tz.String()
}

// GetIANA returns the IANA timezone string (e.g., "Europe/Vienna")
func (tz Timezone) GetIANA() string {
	if tzStr, ok := TimezoneStrings[tz]; ok {
		return tzStr
	}
	return ""
}

// IsValid checks if a timezone value is valid
func IsValid(tz Timezone) bool {
	_, ok := Names[tz]
	return ok
}

// ToInt32 converts the Timezone to int32 for database storage
func (tz Timezone) ToInt32() int32 {
	return int32(tz)
}

// FromInt32 converts an int32 to a Timezone
func FromInt32(i int32) Timezone {
	return Timezone(i)
}

// FromIANA converts an IANA timezone string to a Timezone
func FromIANA(tzStr string) (Timezone, bool) {
	for tz, str := range TimezoneStrings {
		if str == tzStr {
			return tz, true
		}
	}
	return 0, false
}

// Scan implements sql.Scanner interface for reading from database
// Handles both integer and string values from the database
func (tz *Timezone) Scan(value interface{}) error {
	if value == nil {
		*tz = EuropeVienna // Default to Vienna
		return nil
	}

	switch v := value.(type) {
	case int64:
		*tz = Timezone(v)
		return nil
	case int:
		*tz = Timezone(v)
		return nil
	case string:
		// Database stores timezone as IANA string (e.g., "Europe/Vienna")
		parsed, ok := FromIANA(v)
		if !ok {
			// Default to Vienna if unknown
			*tz = EuropeVienna
			return nil
		}
		*tz = parsed
		return nil
	case []byte:
		// Handle byte array as string
		parsed, ok := FromIANA(string(v))
		if !ok {
			*tz = EuropeVienna
			return nil
		}
		*tz = parsed
		return nil
	default:
		*tz = EuropeVienna
		return nil
	}
}

// Value implements driver.Valuer interface for writing to database
// Returns the IANA timezone string for database storage
func (tz Timezone) Value() (interface{}, error) {
	return tz.GetIANA(), nil
}
