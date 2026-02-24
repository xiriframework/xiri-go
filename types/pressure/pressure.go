package pressure

// Pressure represents a pressure unit preference
// It's a type-safe enum backed by int for database compatibility
type Pressure int

// Pressure unit constants
// Supported pressure unit constants
const (
	Bar Pressure = 0 // Bar (metric)
	Psi Pressure = 1 // PSI - Pounds per square inch (imperial)
	Kpa Pressure = 2 // kPa - Kilopascal (metric)
)

// Names maps pressure values to human-readable names for debugging and logging
var Names = map[Pressure]string{
	Bar: "Bar",
	Psi: "Psi",
	Kpa: "Kpa",
}

// Symbols maps pressure values to their unit symbols
var Symbols = map[Pressure]string{
	Bar: "bar",
	Psi: "psi",
	Kpa: "kPa",
}

// String returns the string representation of the pressure value
func (p Pressure) String() string {
	if name, ok := Names[p]; ok {
		return name
	}
	return "Unknown"
}

// GetName returns the human-readable name for a pressure value
func GetName(p Pressure) string {
	return p.String()
}

// GetSymbol returns the unit symbol for a pressure value
func (p Pressure) GetSymbol() string {
	if symbol, ok := Symbols[p]; ok {
		return symbol
	}
	return ""
}

// IsValid checks if a pressure value is valid
func IsValid(p Pressure) bool {
	_, ok := Names[p]
	return ok
}

// IsMetric returns true if the pressure unit is metric (bar)
func (p Pressure) IsMetric() bool {
	return p == Bar || p == Kpa
}

// IsImperial returns true if the pressure unit is imperial (psi)
func (p Pressure) IsImperial() bool {
	return p == Psi
}

// ToInt32 converts the Pressure to int32 for database storage
func (p Pressure) ToInt32() int32 {
	return int32(p)
}

// FromInt32 converts an int32 to a Pressure
func FromInt32(i int32) Pressure {
	return Pressure(i)
}
