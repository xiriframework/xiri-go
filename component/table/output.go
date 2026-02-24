package table

// OutputType defines different table output formats
type OutputType int

const (
	// OutputWeb is for HTML table in web UI (xiri-ui Angular frontend)
	OutputWeb OutputType = iota

	// OutputCSV is for CSV export
	OutputCSV

	// OutputPDF is for PDF report generation
	OutputPDF

	// OutputExcel is for Excel (XLSX) export
	OutputExcel
)

// String returns the string representation of the output type
func (o OutputType) String() string {
	switch o {
	case OutputWeb:
		return "web"
	case OutputCSV:
		return "csv"
	case OutputPDF:
		return "pdf"
	case OutputExcel:
		return "excel"
	default:
		return "unknown"
	}
}
