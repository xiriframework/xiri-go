# Dialoge — alle Typen + Workflows

`component/dialog` bietet fünf Dialog-Typen (`question`, `form`, `table`, `waiting`, `component`) plus eine Handvoll Helper für die häufigen Fälle (Delete, Warning, MultiDelete). Jeder Dialog wird vom Backend als JSON-Response geliefert, das Frontend öffnet einen MatDialog daraus.

## Dialog-Interface

```go
type Dialog interface {
    Print(ctx *core.UiContext) map[string]any
    WithExtra(extra map[string]any) Dialog
    WithOptions(options map[string]any) Dialog
    WithOption(key string, value any) Dialog
    WithSize(size Size) Dialog
    WithData(payload map[string]any) Dialog
    WithTableHeader() Dialog
    SetButtons(buttons []*button.Button) Dialog
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

Wurde der Dialog über den „＋“-Button eines Formfelds geöffnet (`field.BaseField.SetAddURL`), antwortet
der POST stattdessen mit der neuen Option, die das Feld übernimmt und selektiert:

```go
return wc.Component(response.NewReturnDone().WithCreated(d.ID, d.Name).WithMessage("Gespeichert", response.MessageSuccess))
```

### `NewDialogFormMultiEdit` — Formular für mehrere selektierte IDs

```go
dialog.NewDialogFormMultiEdit(fields []map[string]any, u *url.Url, ids []int64, header, okText, closeText string) Dialog
```

Wrapper um `NewDialogForm`, der `extra["data"] = ids` und `extra["done"] = true` setzt; das Frontend
hängt `extra` beim Submit wieder an, sodass der POST-Handler `ExtractMultiSelectRequest` nutzen kann.
Header/Texte als Strings (nicht Pointer). Vollständiges Muster mit GET/POST-Handler: `patterns.md` §5b.

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

### Kein Picker-Dialog

Der Table-Dialog ist im Frontend **reine Anzeige**: `xiri-raw-table` rendert keine Buttons und liefert
keinen Wert zurück. Einen „Zeile wählen und ins Parent-Form übernehmen"-Flow gibt es nicht — dafür
`ModelField`/`ModelListField` (Select mit Options/URL) oder `SetAddURL` (siehe form-fields.md) verwenden.

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

`renderAddress` muss jedes Adressfeld mit `html.EscapeString` escapen — Adressen sind Kundeneingaben.

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

> **Achtung — `url` ist kein Navigationsziel.** Bei `done: true` schließt xiri-ng den Dialog und öffnet
> `apiBaseUrl + url` per `window.open(…, '_blank')` in einem neuen Tab (Download/Report).
> Die URL muss also eine **API-URL** sein (wie `c.apiUrl(…).Print()`), keine Angular-Route wie
> `c.pageUrl()` — die würde als `/api/…` aufgerufen. Blockt der Browser das Popup, zeigt der Dialog
> stattdessen einen „Download"-Button. Soll nach dem Job nur die Seite neu geladen werden, ist der
> Waiting-Dialog das falsche Werkzeug. `blocked` ist deprecated, `""` übergeben.

### Flow

```
1. Button klickt → Controller.StartImport
      returns NewDialogWaiting("Importiere...", pollUrl, "Import", 2000, nil, nil)
      Frontend öffnet Dialog + startet Polling

2. Polling → Controller.ImportStatus
      - Solange noch nicht fertig:  returns NewDialogWaitingNotDone()
      - Fertig:                     returns NewDialogWaitingDone(c.apiUrl("import", "result", jobID).Print(), "")
                                    → Frontend öffnet apiBaseUrl + url in neuem Tab
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
        c.apiUrl("import", "result", jobID).Print(), // API-URL (z.B. Report-Download), wird per window.open geöffnet
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

// Handler: beide Schritte kommen als POST auf dieselbe URL.
//   1. Select-Button der Tabelle: POST {"data":[ids]}              → Dialog öffnen
//   2. OK im Dialog:              POST {"data":[ids],"done":true}  → wirklich löschen
// ExtractMultiSelectRequest unterscheidet über den "done"-Key (isDialogOpen).
func (c *Controller) BulkDeleteConfirm(ctx echo.Context) error {
    wc := webcontext.GetWebContext(ctx)

    ids, _, isDialogOpen, err := dialog.ExtractMultiSelectRequest(ctx)
    if err != nil { return wc.BadRequest(err.Error()) }

    if isDialogOpen {
        dlg := dialog.NewDialogFormMultiDelete(
            c.apiUrl("bulk-delete-confirm"),  // POST-Target (selbe URL)
            ids,
            fmt.Sprintf("%d Geräte wirklich löschen?", len(ids)),
            nil, nil, nil,
        )
        return wc.Component(dlg)
    }

    // Submit: Bulk-Delete ausführen
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
    ids          []int64,
    body         map[string]interface{}, // kompletter Request-Body
    isDialogOpen bool,                   // true = kein "done"-Key → Schritt 1 (Dialog öffnen)
    err          error,
)
```

Liest die IDs aus `body["data"]` (Zahlen oder numerische Strings, andere Elemente werden übersprungen).
Fehlt `data`, gibt es einen Fehler. `isDialogOpen` unterscheidet „Dialog öffnen" (Select-Button, ohne `done`)
vom Submit des Dialogs (mit `done: true`).

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
- **MultiDelete ohne `ExtractMultiSelectRequest`**: die IDs kommen immer unter `data`, mal als Zahl, mal als String. Der Extractor normalisiert das und liefert `isDialogOpen`.
- **`NewDialogFormMultiDelete`: Öffnen und Submit sind beide POST auf dieselbe URL.** Der Select-Button schickt `{"data":[ids]}`, das OK im Dialog `{"data":[ids],"done":true}`. Nicht über `ctx.Request().Method` unterscheiden (immer POST), sondern über `isDialogOpen` aus `ExtractMultiSelectRequest`.
- **Größen:** Nutze `WithSize(dialog.SizeLg)` statt `WithOption("size", "lg")` — typsicher und compile-time-geprüft. Raw CSS-Werte (`"750px"`) funktionieren weiterhin als Escape-Hatch.
