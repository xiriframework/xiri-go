package field

import (
	"github.com/xiriframework/xiri-go/uicontext"
)

// DividerField represents a visual divider/separator in a form
type DividerField struct {
	*BaseField
	Content string // Optional label text for the divider
}

func (f *DividerField) Validate(value interface{}) error {
	return nil
}

func (f *DividerField) Parse(raw interface{}) (interface{}, error) {
	return nil, nil
}

// ExportForFrontend exports the field for frontend rendering
func (f *DividerField) ExportForFrontend(ctx *uicontext.UiContext, value interface{}) map[string]interface{} {
	return f.BaseField.GetBaseExport(ctx, f.Content)
}

// NewDividerField creates a visual divider/separator field
func NewDividerField(id string) *DividerField {
	return &DividerField{
		BaseField: &BaseField{
			ID:       id,
			Type:     FieldTypeDivider,
			Name:     "",
			Required: false,
			Default:  nil,
			Form:     true,
		},
	}
}

// SetContent sets the optional label text for the divider
func (f *DividerField) SetContent(content string) *DividerField {
	f.Content = content
	return f
}

// SetClass sets the CSS class for frontend styling
func (f *DividerField) SetClass(class string) *DividerField {
	f.BaseField.SetClass(class)
	return f
}

// SetStep sets the step indicator for multi-step forms
func (f *DividerField) SetStep(step int) *DividerField {
	f.BaseField.SetStep(step)
	return f
}
