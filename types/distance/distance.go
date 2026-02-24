package distance

// Distance represents a distance unit preference
// It's a type-safe enum backed by int for database compatibility
type Distance int

// Distance unit constants
// Supported distance unit constants
const (
	Kilometer Distance = 0 // Kilometers (metric)
	Miles     Distance = 1 // Miles (imperial)
	Seemiles  Distance = 2 // Nautical miles (maritime)
)

// Names maps distance values to human-readable names for debugging and logging
var Names = map[Distance]string{
	Kilometer: "Kilometer",
	Miles:     "Miles",
	Seemiles:  "Seemiles",
}

// Symbols maps distance values to their unit symbols
var Symbols = map[Distance]string{
	Kilometer: "km",
	Miles:     "mi",
	Seemiles:  "NM",
}

// String returns the string representation of the distance value
func (d Distance) String() string {
	if name, ok := Names[d]; ok {
		return name
	}
	return "Unknown"
}

// GetName returns the human-readable name for a distance value
func GetName(d Distance) string {
	return d.String()
}

// GetSymbol returns the unit symbol for a distance value
func (d Distance) GetSymbol() string {
	if symbol, ok := Symbols[d]; ok {
		return symbol
	}
	return ""
}

// IsValid checks if a distance value is valid
func IsValid(d Distance) bool {
	_, ok := Names[d]
	return ok
}

// IsMetric returns true if the distance unit is metric (kilometers)
func (d Distance) IsMetric() bool {
	return d == Kilometer
}

// IsImperial returns true if the distance unit is imperial (miles)
func (d Distance) IsImperial() bool {
	return d == Miles
}

// IsMaritime returns true if the distance unit is maritime (nautical miles)
func (d Distance) IsMaritime() bool {
	return d == Seemiles
}

// ToInt32 converts the Distance to int32 for database storage
func (d Distance) ToInt32() int32 {
	return int32(d)
}

// FromInt32 converts an int32 to a Distance
func FromInt32(i int32) Distance {
	return Distance(i)
}
