package table

// ExportNumber is returned by CSV/Excel formatters for numeric cells.
type ExportNumber struct {
	Value    float64
	Decimals int
}
