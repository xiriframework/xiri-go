package dialog

import (
	"github.com/xiriframework/xiri-go/component/button"
	"github.com/xiriframework/xiri-go/component/core"
)

/*
Creates a dialog that renders an arbitrary core.Component (e.g. expansion, card,
stepper) as its content, with an Ok button.

Angular Frontend Flow:
  1. Dialog receives response with type="component"
  2. Angular reads res.content (the component's Print() output: {type, display, data})
  3. Renders it through the generic renderer: <xiri-dyncomponent [data]="content">

Because core.Component and DialogContent share the same Print signature
(Print(ctx *core.UiContext) map[string]any), the dialog's Print() serializes the
component automatically — no wrapper is required.

Additional dialog options via the fluent interface:
  - WithSize(dialog.SizeLg): dialog size
  - SetButtons([]*button.Button{…}): replace the default Ok button
  - WithOption(key, value): structural top-level fields (e.g. url)

Example usage:

	exp := expansion.NewExpansion().
	    AddPanel(expansion.NewPanel("Lieferadresse").AddContent(addressCard)).
	    AddPanel(expansion.NewPanel("Rechnungsadresse").AddContent(billingCard))

	dlg := dialog.NewDialogComponent(t("INFO"), exp).WithSize(dialog.SizeLg)
*/

// NewDialogComponent creates a dialog that displays an arbitrary core.Component.
// The component is serialized lazily at Print time with the provided ctx.
func NewDialogComponent(
	header string,
	component core.Component,
) Dialog {

	return NewDialog(
		core.DialogTypeComponent,
		header,
		component,
		[]*button.Button{
			button.NewCloseButton("Ok", core.ColorPrimary, core.ButtonTypeRaised, "", false, nil, false, nil),
		},
		nil,
		nil,
	)
}
