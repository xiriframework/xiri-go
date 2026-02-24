package button

import (
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/response"
)

// ButtonLine represents a container for multiple buttons with layout options
type ButtonLine struct {
	class   string
	display *string
	buttons []*Button
}

// NewButtonLine creates a new button line container
// If class is empty string, defaults to "right"
func NewButtonLine(class string, display *string) *ButtonLine {
	if class == "" {
		class = "right"
	}
	return &ButtonLine{
		class:   class,
		display: display,
		buttons: make([]*Button, 0),
	}
}

// Add adds a button to the button line
func (bl *ButtonLine) Add(button *Button) *ButtonLine {
	bl.buttons = append(bl.buttons, button)
	return bl
}

// Builder methods for optional button line properties

// WithDisplay sets the display/layout class (optional)
// Returns the ButtonLine for method chaining
func (bl *ButtonLine) WithDisplay(display string) *ButtonLine {
	bl.display = &display
	return bl
}

// Print returns the full component JSON representation
func (bl *ButtonLine) Print(translator core.TranslateFunc) map[string]any {
	// Serialize buttons
	buttonData := make([]map[string]any, len(bl.buttons))
	for i, button := range bl.buttons {
		buttonData[i] = button.Print(translator)
	}

	return map[string]any{
		"type":    "buttonline",
		"display": bl.display,
		"data": map[string]any{
			"class":   bl.class,
			"buttons": buttonData,
		},
	}
}

// PrintData returns just the data portion (without type wrapper)
func (bl *ButtonLine) PrintData(translator core.TranslateFunc) map[string]any {
	// Serialize buttons
	buttonData := make([]map[string]any, len(bl.buttons))
	for i, button := range bl.buttons {
		buttonData[i] = button.Print(translator)
	}

	return map[string]any{
		"class":   bl.class,
		"buttons": buttonData,
	}
}

// DataResponse returns a DataResult wrapping the button line data in {"data": ...} envelope.
func (bl *ButtonLine) DataResponse(translator core.TranslateFunc) response.DataResult {
	return response.NewJSONDataResult(bl.PrintData(translator))
}

// PrintButtons returns just the buttons array
func (bl *ButtonLine) PrintButtons(translator core.TranslateFunc) []map[string]any {
	// Serialize buttons
	buttonData := make([]map[string]any, len(bl.buttons))
	for i, button := range bl.buttons {
		buttonData[i] = button.Print(translator)
	}
	return buttonData
}
