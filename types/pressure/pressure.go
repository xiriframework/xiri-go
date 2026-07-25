package pressure

import "slices"

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

// names maps pressure values to human-readable names for debugging and logging
var names = map[Pressure]string{
	Bar: "Bar",
	Psi: "Psi",
	Kpa: "Kpa",
}

// symbols maps pressure values to their unit symbols
var symbols = map[Pressure]string{
	Bar: "bar",
	Psi: "psi",
	Kpa: "kPa",
}

// All returns every valid Pressure, ordered by its numeric value.
// The returned slice is freshly allocated on each call, so callers can sort or
// filter it without affecting the package's internal lookup tables.
func All() []Pressure {
	out := make([]Pressure, 0, len(names))
	for v := range names {
		out = append(out, v)
	}
	slices.Sort(out)
	return out
}

// String returns the string representation of the pressure value
func (p Pressure) String() string {
	if name, ok := names[p]; ok {
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
	if symbol, ok := symbols[p]; ok {
		return symbol
	}
	return ""
}

// IsValid checks if a pressure value is valid
func IsValid(p Pressure) bool {
	_, ok := names[p]
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
