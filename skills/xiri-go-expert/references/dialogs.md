# Dialoge — alle Typen + Workflows

`component/dialog` bietet fünf Dialog-Typen (`question`, `form`, `table`, `waiting`, `component`) plus eine Handvoll Helper für die häufigen Fälle (Delete, Warning, MultiDelete). Jeder Dialog wird vom Backend als JSON-Response geliefert, das Frontend öffnet einen MatDialog daraus.

## Dialog-Interface

```go
type Dialog interface {
    core.Component
    WithExtra(extra map[string]any) Dialog
    WithOptions(options map[string]any) Dialog
    WithOption(key string, value any) Dialog
    WithData(payload map[string]any) Dialog
    WithTableHeader() Dialog
    SetButtons(buttons []*button.Button) Dialog
    Print(ctx *core.UiContext) map[string]any
}
```

`WithTableHeader()` aktiviert die Spalten-Header-Zeile (`<thead>`) bei `DialogTypeTable`. Default = aus (Header verborgen), kompatibel mit allen bestehenden Aufrufern. Wirkt nur bei Table-Dialogen, bei anderen Typen no-op.

`WithOption` / `WithOptions` sind für **strukturelle Top-Level-Felder** (`size`, `url`, `filter`, …) gedacht und werden von der Frontend-Logik direkt am Dialog-Root gelesen. Für **Custom-Payload**, das später ans Backend zurückgesendet wird, nutze `WithData(map[string]any)` — das Ergebnis landet unter einem expliziten `data`-Key im JSON (analog Button/Icon).

## `NewDialog` — generischer Konstruktor

```go
func dialog.NewDialog(
    dialogType core.DialogType,         // DialogTypeQuestion | Form | Table | Waiting | Component
    header     string,
    content    any,                      // DialogContent oder Component (wird geprintet)
    buttons    []*button.Button,
    extra      map[string]any,            // zusätzliche Daten im JSON-Root (z.B. selectedIds)
    options    map[string]any,            // Dialog-Behavior: size, url, filter
) Dialog
```

Gehört zur seltenen Kategorie "maximum flexibility" — für die typischen Fälle gibt es spezialisierte Konstruktoren.

### Generische Optionen via `WithOption(key, value)`

| Key      | Typ      | Zweck                                          |
| -------- | -------- | ---------------------------------------------- |
| `size`   | `string` | Dialog-Größe: `"sm"`, `"md"`, `"lg"`, `"xl"`, `"full"` — bevorzugt typsicher via `WithSize(...)` |
| `url`    | `string` | Data-/Submit-URL (bei Form/Table)              |
| `filter` | `any`    | Filter-Data, die ans Ziel weitergereicht wird  |

### Dialog-Größe — `WithSize(dialog.Size)`

Bevorzugter, typsicherer Weg. Frontend mappt das Token auf eine konkrete Breite:

| Token             | Frontend-Breite (Desktop) |
| ----------------- | ------------------------- |
| `dialog.SizeSm`   | `400px`                   |
| `dialog.SizeMd`   | `600px` (Default)         |
| `dialog.SizeLg`   | `900px`                   |
| `dialog.SizeXl`   | `1200px`                  |
| `dialog.SizeFull` | `95vw` (fast Vollbild)    |

Auf XSmall/Small Breakpoints (Mobile) ignoriert das Frontend die Größe und nutzt `90vw`.

```go
dlg := dialog.NewDialog(core.DialogTypeForm, "Bearbeiten", fields, buttons, nil, nil).
    WithSize(dialog.SizeLg)

// Equivalent (alte API, weiterhin unterstützt):
dlg.WithOption("size", "lg")

// Auch raw CSS-Werte sind möglich (Backward-Compat):
dlg.WithOption("size", "750px")
```

## Question-Dialog — Delete & Warning

```go
// Delete (Icon "warning")
dialog.NewDialogDelete(
    text       string,          // "Gerät wirklich löschen?"
    u          *url.Url,         // POST-Ziel bei Bestätigung
    extra      map[string]any,   // nil oder z.B. {"itemId": 42}
    headerText *string,          // nil = "Delete" (vom Translator übersetzt)
    okText     *string,          // nil = "Ok"
    closeText  *string,          // nil = "Back"
) Dialog

// Warning (Icon "help_outline")
dialog.NewDialogWarning(text, u, extra, headerText, okText, closeText)
```

Beide erzeugen einen `DialogTypeQuestion` mit `DialogQuestionContent{Icon, Question}` und zwei Standard-Buttons (Close + Submit).

### Beispiel-Use

```go
func (c *Controller) Delete(ctx echo.Context) error {
    wc := webcontext.GetWebContext(ctx)
    id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)

    if ctx.Request().Method == "GET" {
        d, _ := c.svc.DB.Device.GetByID(id)
        dlg := dialog.NewDialogDelete(
            fmt.Sprintf("Gerät '%s' wirklich löschen?", d.Name),
            c.apiUrl("delete", ctx.Param("id")),
            nil, nil, nil, nil,
        )
        return wc.Component(dlg)
    }

    // POST: Löschen
    _ = c.svc.DB.Device.Delete(ctx.Request().Context(), id)
    return wc.RefreshTable()
}
```

### `HandleDelRequest` — Convenience-Handler

Kürzt den GET/POST-Boilerplate in eine Funktion:

```go
func (c *Controller) Delete(ctx echo.Context) error {
    wc := webcontext.GetWebContext(ctx)
    id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)

    return dialog.HandleDelRequest(
        ctx,
        fmt.Sprintf("Gerät #%d wirklich löschen?", id),
        c.apiUrl("delete", strconv.FormatInt(id, 10)).PrintPrefix(),
        func() error { return c.svc.DB.Device.Delete(ctx.Request().Context(), id) },
        response.NewReturnRefreshTable(),
        wc.UiContext(),
    )
}
```

## Form-Dialog — Formular im Dialog statt auf eigener Seite

```go
func dialog.NewDialogForm(
    fields    []map[string]any,  // aus fb.BuildAddForDisplay() oder fb.BuildEditForDisplay()
    u         *url.Url,          // Submit-Target
    header    *string,           // POINTER — `*header` wird als Dialog-Titel genutzt
    extra     map[string]any,
    okText    *string,           // nil = "Ok"
    closeText *string,           // nil = "Back"
) Dialog
```

### Pattern: Add-via-Dialog

```go
// GET: liefert den Dialog mit leerem Form
func (c *Controller) AddDialog(ctx echo.Context) error {
    wc := webcontext.GetWebContext(ctx)
    uc := wc.UiContext()

    nameF, activeF := buildFormFields(nil)
    fb := formbuilder.NewFormBuilder(uc).AddField(nameF).AddField(activeF)
    fields, err := fb.BuildAddForDisplay()
    if err != nil { return wc.InternalServerError(err.Error()) }

    header := "Neues Gerät"
    dlg := dialog.NewDialogForm(
        fields,
        c.apiUrl("add-dialog"),   // Submit-URL
        &header,
        nil, nil, nil,
    ).WithOption("size", "md")

    return wc.Component(dlg)
}

// POST: validiert + speichert + refresht Tabelle
func (c *Controller) AddDialogSubmit(ctx echo.Context) error {
    wc := webcontext.GetWebContext(ctx)
    nameF, activeF := buildFormFields(nil)
    fb := formbuilder.NewFormBuilder(wc.UiContext()).AddField(nameF).AddField(activeF)

    fg, _, _ := fb.BuildAdd()
    if err := formbuilder.BindAndValidate(ctx, fg); err != nil {
        return wc.BadRequest(err.Error())
    }

    d := &entities.Device{Name: *nameF.Value}
    if activeF.Value != nil { d.Active = *activeF.Value }
    if err := c.svc.DB.Device.Create(ctx.Request().Context(), d); err != nil {
        return wc.InternalServerError(err.Error())
    }
    return wc.RefreshTable()   // schließt Dialog + reloaded Parent-Table
}
```

Der `action: dialog`-Button in der Tabelle (oder PageHeader) zeigt auf `AddDialog`; der Dialog selbst POSTed an `AddDialogSubmit`.

## Table-Dialog — Picker-Dialog oder Info-Anzeige

```go
func dialog.NewDialogTable[T any](header string, tbl *table.Table[T]) Dialog
```

Rendert ein `DialogTypeTable` mit einer `xiri-raw-table` darin. Daten + Fields werden **lazy** beim `Print(ctx)` aus der Table gezogen.

**Spalten-Header standardmäßig aus.** Mit `.WithTableHeader()` opt-in:

```go
dlg := dialog.NewDialogTable("Geräte", tbl).WithTableHeader()
```

### Beispiel: Info-Dialog mit Details

```go
type InfoRow struct {
    Label string
    Value string
}

func (c *Controller) Info(ctx echo.Context) error {
    wc := webcontext.GetWebContext(ctx)
    id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
    d, _ := c.svc.DB.Device.GetByID(id)

    b := table.NewBuilder[InfoRow]()
    b.TextField("label", "Feld", func(r InfoRow) string { return r.Label })
    b.TextField("value", "Wert", func(r InfoRow) string { return r.Value })

    tbl := b.Build()
    tbl.SetData([]InfoRow{
        {Label: "Name",   Value: d.Name},
        {Label: "Status", Value: d.Status},
        {Label: "IP",     Value: d.IP},
    })

    return wc.Component(dialog.NewDialogTable("Device-Details", tbl))
}
```

### Beispiel: Picker-Dialog (User wählt Zeile)

Der Table-Dialog rendert — wenn der Frontend-Button in einer Zeile ein `FieldButtonActionClose` (oder `Save`) hat — beim Klick den Row-Wert zurück in das Parent-Form.

```go
// Picker: User wählt Device aus der Liste
b := table.NewBuilder[Device]()
b.IdField("id", "ID", func(r Device) int64 { return r.ID })
b.TextField("name", "Name", func(r Device) string { return r.Name })
b.ButtonsField("select", "", func(r Device) map[string]string {
    return map[string]string{"0": strconv.FormatInt(r.ID, 10)}
}).AddButton(0, table.FieldButtonActionSave, "check", core.ColorPrimary, "Wählen")

tbl := b.Build()
tbl.SetData(c.svc.DB.Device.AllActive())

dlg := dialog.NewDialogTable("Gerät wählen", tbl).WithOption("size", "lg")
return wc.Component(dlg)
```

Der Frontend-Flow: User klickt in der Zeile, `select`-Button liefert den Wert (`r.ID`) zurück an das Form, das den Dialog geöffnet hat.

## Component-Dialog — beliebige Komponente im Dialog

```go
func dialog.NewDialogComponent(header string, component core.Component) Dialog
```

Rendert eine **beliebige `core.Component`** (z. B. `expansion`, `card`, `stepper`, `tabs`) als Dialog-Inhalt — ideal für read-only Detail-/Info-Ansichten, die mehr als eine Tabelle brauchen. Erzeugt einen `DialogTypeComponent` mit einem Default-`Close`-Button.

Funktioniert **ohne eigene Serialisierung**, weil `core.Component` und `DialogContent` dieselbe Signatur `Print(ctx *core.UiContext) map[string]any` haben: `dialogImpl.Print()` ruft `component.Print(ctx)` und bettet das Ergebnis (`{type, display, data}`) als `content` ein. Das Frontend rendert diesen `content` über den generischen `xiri-dyncomponent` (siehe xiri-ng-expert).

Anpassungen laufen über die fluente `Dialog`-Schnittstelle (`WithSize`, `SetButtons`, …).

### Beispiel: Expansion-Akkordeon im Dialog

```go
func (c *Controller) OrderDetails(ctx echo.Context) error {
    wc := webcontext.GetWebContext(ctx)
    id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
    o, _ := c.svc.DB.Order.GetByID(id)

    exp := expansion.NewExpansion().WithMulti(true).
        AddPanel(expansion.NewPanel("Lieferadresse").WithIcon("local_shipping").WithExpanded(true).
            AddContent(layout.NewHtml(renderAddress(o.Delivery), nil))).
        AddPanel(expansion.NewPanel("Rechnungsadresse").WithIcon("receipt_long").
            AddContent(layout.NewHtml(renderAddress(o.Billing), nil)))

    dlg := dialog.NewDialogComponent(fmt.Sprintf("Order #%d", o.ID), exp).
        WithSize(dialog.SizeLg)
    return wc.Component(dlg)
}
```

GET liefert `{header, type:"component", content:{type:"expansion", display, data:{panels:[…]}}, buttons:[close]}`. Da der Inhalt read-only ist, gibt es i. d. R. keinen POST-Submit — nur der Close-Button.

> **Tipp:** Mehrere aufklappbare Abschnitte als gestapelte, unabhängige Panels gehören in eine `expansion` (echtes Akkordeon). Mehrere collapsible `HeaderField` in *einem Formular* erzeugen ebenfalls gestapelte Sektionen (jeder Header startet seine eigene Section bis zum nächsten Header) — aber ohne nestbare Inhalte. Für reine Detail-Ansichten im Dialog ist `NewDialogComponent` + `expansion` der direkte Weg.

## Waiting-Dialog — Long-Running-Tasks mit Polling

```go
func dialog.NewDialogWaiting(
    text      string,          // "Wird verarbeitet…"
    u         *url.Url,         // Polling-URL (wird wiederholt GET'ed)
    header    string,           // "Import läuft"
    checkTime int,              // Polling-Interval in ms (z.B. 2000)
    extra     map[string]any,
    closeText *string,          // nil = "Back"
) Dialog
```

Zusätzlich gibt es **Response-Helper** für die Polling-Antworten:

```go
dialog.NewDialogWaitingNotDone()          // {"done": false}  — weiter pollen
dialog.NewDialogWaitingDone(url, blocked) // {"done": true, "url": ..., "blocked": ...}
dialog.NewDialogWaitingError(message)     // {"done": true, "error": ...}
```

### Flow

```
1. Button klickt → Controller.StartImport
      returns NewDialogWaiting("Importiere...", /api/import/status, 2000, ...)
      Frontend öffnet Dialog + startet Polling

2. Polling → Controller.ImportStatus
      - Solange noch nicht fertig:  returns NewDialogWaitingNotDone()
      - Fertig:                     returns NewDialogWaitingDone("/devices", "")
      - Fehler:                     returns NewDialogWaitingError("Import fehlgeschlagen: ...")
```

### Beispiel

```go
func (c *Controller) StartImport(ctx echo.Context) error {
    wc := webcontext.GetWebContext(ctx)
    jobID, err := c.svc.ImportService.Start()
    if err != nil { return wc.InternalServerError(err.Error()) }

    pollUrl := c.apiUrl("import", "status", jobID)
    dlg := dialog.NewDialogWaiting(
        "Import läuft...",
        pollUrl,
        "Import",
        2000, // 2s Polling-Interval
        map[string]any{"jobId": jobID},
        nil,
    )
    return wc.Component(dlg)
}

func (c *Controller) ImportStatus(ctx echo.Context) error {
    wc := webcontext.GetWebContext(ctx)
    jobID := ctx.Param("id")

    job, err := c.svc.ImportService.Status(jobID)
    if err != nil {
        return wc.Component(dialog.NewDialogWaitingError(err.Error()))
    }
    if !job.Done {
        return wc.Component(dialog.NewDialogWaitingNotDone())
    }
    if job.Failed {
        return wc.Component(dialog.NewDialogWaitingError(job.ErrorMessage))
    }
    return wc.Component(dialog.NewDialogWaitingDone(
        c.pageUrl().Print(),  // Ziel nach Completion
        "",
    ))
}
```

## MultiDelete-Dialog (via SelectButtons)

Beim Bulk-Delete aus einer Tabelle wird die Selected-IDs-Liste an den Server geschickt. Der Server öffnet einen Confirm-Dialog, der die IDs in `extra.data` mitführt und beim OK an die Delete-URL POSTed.

```go
// Table-Setup
b.SetSelectButtons([]*button.TableButton{
    button.NewTableButton(
        core.ButtonActionDialog,
        "delete",
        c.apiUrl("bulk-delete-confirm"),   // → öffnet Confirm-Dialog
        "Löschen",
        core.ColorWarning,
        false,
        nil,
    ),
})

// Handler: GET = Dialog öffnen, POST = wirklich löschen
func (c *Controller) BulkDeleteConfirm(ctx echo.Context) error {
    wc := webcontext.GetWebContext(ctx)

    if ctx.Request().Method == "GET" {
        // Die Frontend-Buttonline hängt die IDs im POST-Body an.
        // Hier lesen wir sie aus, um sie im Dialog-Extra wieder mitzuschicken.
        ids, _, _, err := dialog.ExtractMultiSelectRequest(ctx)
        if err != nil { return wc.BadRequest(err.Error()) }

        dlg := dialog.NewDialogFormMultiDelete(
            c.apiUrl("bulk-delete-confirm"),  // POST-Target (selbe URL)
            ids,
            fmt.Sprintf("%d Geräte wirklich löschen?", len(ids)),
            nil, nil, nil,
        )
        return wc.Component(dlg)
    }

    // POST: Bulk-Delete ausführen
    ids, _, _, err := dialog.ExtractMultiSelectRequest(ctx)
    if err != nil { return wc.BadRequest(err.Error()) }

    if err := c.svc.DB.Device.DeleteMany(ctx.Request().Context(), ids); err != nil {
        return wc.InternalServerError(err.Error())
    }
    return wc.Component(response.NewReturnRefreshTable().
        WithMessage(fmt.Sprintf("%d gelöscht", len(ids)), response.MessageSuccess))
}
```

### `ExtractMultiSelectRequest`

```go
func dialog.ExtractMultiSelectRequest(c echo.Context) (
    ids      []int64,
    extra    map[string]interface{},
    hasExtra bool,
    err      error,
)
```

Parst die IDs aus dem Request-Body (unterstützt mehrere Payload-Formate, die das Frontend schickt) — bequemer als manuelles `c.Bind(&struct{IDs []int64}{})`.

## Dialog-Extras und Options

Manchmal muss man Daten durch den Dialog durchreichen (z.B. Parent-Row-ID, Filter-Context):

```go
dlg := dialog.NewDialogForm(fields, submitURL, &header, nil, nil, nil).
    WithExtra(map[string]any{
        "parentId": parentID,
        "context":  "edit",
    }).
    WithOption("size", "lg").
    WithOption("filter", filterData)
```

Beim Submit liefert der Frontend-Request diese Extras **zusätzlich** zum Form-Value mit — der Controller kann sie via `c.Bind(...)` einlesen.

## `DialogContent`-Interface

Für Custom-Content wird ein Typ implementiert, der `Print(ctx)` liefert:

```go
type DialogContent interface {
    Print(ctx *core.UiContext) map[string]any
}
```

Eingebaut:

| Typ                     | Zweck                                     |
| ----------------------- | ----------------------------------------- |
| `DialogQuestionContent` | `{icon, question}` — für Confirm-Dialoge  |
| `DialogWaitingContent`  | `{icon, text}`     — für Waiting-Anzeige  |

Custom-Content (z.B. eine Komponente): einfach eine Komponente als `content` an `NewDialog` übergeben. `dialogImpl.Print` ruft `content.Print(ctx)` (wenn es ein `core.Component` oder `DialogContent` ist) automatisch auf. Für genau diesen Fall gibt es den Convenience-Konstruktor `NewDialogComponent(header, component)` mit `DialogTypeComponent` (siehe oben).

## Häufige Fehler

- **Form-Dialog: `header` ist `*string`** — nicht `string`. `&header` oder `nil` übergeben.
- **Waiting-Dialog: URL muss den API-Prefix tragen**, weil Polling-GETs darauf gehen. `xurl.NewUrlPrefix(...)` oder `c.apiUrl(...)`.
- **MultiDelete ohne `ExtractMultiSelectRequest`**: die IDs im Request kommen in verschiedenen Feldern (`ids`, `data`, `selectedIds`, je nach Button-Action). Der Extractor glättet das.
- **`NewDialogFormMultiDelete` POST-URL = selbe GET-URL**: Der Delete-Dialog POSTed beim OK wieder an die URL, die im Button als `action: dialog` definiert wurde. Deshalb dieselbe URL — Controller unterscheidet via `ctx.Request().Method`.
- **Größen:** Nutze `WithSize(dialog.SizeLg)` statt `WithOption("size", "lg")` — typsicher und compile-time-geprüft. Raw CSS-Werte (`"750px"`) funktionieren weiterhin als Escape-Hatch.
