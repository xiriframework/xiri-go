# Table Builder (component/table)

Import: `"github.com/xiriframework/xiri-go/component/table"`

Generische `Table[T]` mit typed Field-Accessors, Server-Side Pagination, und Multi-Format Export.

## Builder erstellen

```go
b := table.NewBuilder[Device]()
// Defaults: Pagination, Search, Reload, CSV, Excel aktiviert
```

## Typed Field Accessors

Jedes Field hat: `(id string, name string, accessor func(T) ValueType) *FieldBuilder`

```go
// ID (int64, spezieller "id" Typ)
b.IdField("id", "ID", func(d Device) int64 { return d.ID })

// Text
b.TextField("name", "device.name", func(d Device) string { return d.Name })
b.TextNField("info", "device.info", func(d Device) []string { return []string{d.Line1, d.Line2} })

// Zahlen
b.IntField("count", "device.count", func(d Device) int { return d.Count })
b.Int32Field("status", "device.status", func(d Device) int32 { return d.Status })
b.Int64Field("total", "device.total", func(d Device) int64 { return d.Total })
b.FloatField("price", "device.price", func(d Device) float64 { return d.Price })

// Boolean
b.BoolField("active", "device.active", func(d Device) bool { return d.Active })
b.BoolNField("flags", "device.flags", func(d Device) []bool { return []bool{d.Flag1, d.Flag2} })

// Datum/Zeit (Unix-Timestamp int64, nil-safe via Pointer)
b.DateTimeField("created", "device.created", func(d Device) int64 { return d.CreatedAt })
b.DateTimeNField("times", "device.times", func(d Device) []*int64 { return []*int64{d.T1, d.T2} })
b.DateField("date", "device.date", func(d Device) int64 { return d.Date })
b.DateNField("dates", "device.dates", func(d Device) []*int64 { return []*int64{d.D1, d.D2} })

// Einheiten (auto-konvertiert je UiContext)
b.DistanceField("km", "device.distance", func(d Device) float64 { return d.Km })
b.SpeedField("speed", "device.speed", func(d Device) float64 { return d.Speed })
b.PressureField("pressure", "device.pressure", func(d Device) float64 { return d.Pressure })

// Dauer (HH:MM oder Xd HH:MM)
b.TimeLengthField("duration", "device.duration", func(d Device) int { return d.DurationMin })

// Buttons (Aktionen pro Zeile)
b.ButtonsField("actions", "", func(d Device) map[string]string {
    return map[string]string{
        "id": fmt.Sprintf("%d", d.ID),
    }
})

// Link
b.LinkField("link", "device.link", func(d Device) [2]string {
    return [2]string{d.Name, fmt.Sprintf("/device/%d", d.ID)}
})

// Icon aus IconSet
b.IconFieldFromSet("status", "device.status", func(d Device) *table.IconRef {
    return statusIcons.Resolve(d.Status)
}, statusIcons)

// HTML
b.HtmlField("badge", "device.badge", func(d Device) string {
    return fmt.Sprintf(`<span class="badge">%s</span>`, d.Status)
})

// Input (inline editierbar)
b.InputField("value", "device.value", func(d Device) string { return d.Value })

// Header (Abschnitts-Trenner)
b.HeaderField("section", "Abschnitt")

// Chips (Status-Tags pro Zelle, kein Inline-Edit)
// Accessor liefert []table.Chip{Label, Color}. Pro Zelle werden alle Chips
// als farbige Pillen gerendert. Für Single-Wert-Zellen einfach ein 1-elementiges Slice.
b.ChipsField("warnings", "Warning lights", func(d Device) []table.Chip {
    out := []table.Chip{}
    for _, w := range d.Warnings {
        out = append(out, table.Chip{Label: w.Label, Color: w.Severity})
    }
    return out
})
b.ChipsField("battery", "Battery", func(d Device) []table.Chip {
    if d.Battery == nil { return []table.Chip{{Label: "N/A", Color: core.ColorGray}} }
    color := core.ColorLightGray
    if *d.Battery < 50 { color = core.ColorRed }
    return []table.Chip{{Label: fmt.Sprintf("%d%%", *d.Battery), Color: color}}
})
```

## FieldBuilder Methoden (Method Chaining)

```go
b.TextField("name", "device.name", accessor).
    WithWidth("200px").
    WithMinWidth("100px").
    WithAlign(table.AlignLeft).     // AlignLeft, AlignCenter, AlignRight
    WithSort(true).
    WithSearch(true).
    WithSticky(true).               // Fixierte Spalte
    WithHint("device.name.hint").
    WithDisplay("xcol-md-6").
    WithTextPrefix("€").
    WithTextSuffix("kg").
    WithHeader("Gruppe A").         // Header-Zeile über der Spalte
    WithHeaderSpan(3).              // Header spans N Spalten
    WithColumnOrder(1).             // Spaltenreihenfolge
    Hide().                          // Spalte verstecken
    HideInCSV().                     // Nur in CSV verstecken
    ShowInCSV().                     // Nur in CSV anzeigen
    WithFooterSum().                 // Footer: Summe
    WithFooterCount().               // Footer: Anzahl
    WithDecimals(2).                 // Dezimalstellen (Float/Distance)
    WithBoolText("Ja", "Nein")      // Custom Bool-Text

// Custom Formatter
b.TextField("name", "device.name", accessor).
    WithFormatter(table.FormatterFunc(func(value any, row table.Row, output table.OutputType, ctx *core.UiContext) any {
        // value: Rohwert aus Accessor
        // row: Zugriff auf andere Felder via row.Get("fieldId"), row.GetString("fieldId")
        // output: OutputWeb, OutputCSV, OutputExcel
        return fmt.Sprintf("Prefix: %s", value)
    }))

// Per-Output Formatter
b.TextField(...).WithWebFormatter(webFmt).WithCSVFormatter(csvFmt).WithExcelFormatter(excelFmt)
```

## Inline Editing

Felder können direkt in der Tabelle editiert werden. Das Frontend sendet `POST { id, field, value }` an die `editUrl`.

### Table-Level: EditUrl setzen (Pflicht)

```go
b.SetEditUrl("/api/entity/inline-edit")
```

### Field-Level: Editable markieren

```go
// Text-Input (freie Eingabe)
b.TextField("name", "Name", accessor).WithEditable(true)

// Select (statische Optionen) — map[string]string{value: label}
b.TextField("status", "Status", accessor).
    WithEditableOptions(map[string]string{
        "active":   "Aktiv",
        "inactive": "Inaktiv",
    })

// Select (dynamische Optionen per URL)
// Frontend ruft GET {url}?id={rowId}&field={fieldId} auf
b.TextField("category", "Kategorie", accessor).
    WithEditableOptionsUrl("/api/categories")

// Chips (Multi-Select mit Farben — INLINE-EDIT)
// Der Accessor liefert hier den aktuellen String-Slice der ausgewählten Werte.
// Für reine Anzeige-Chips ohne Edit: siehe ChipsField oben.
b.TextField("tags", "Tags", chipsAccessor).
    WithEditableChipOptions([]table.EditableChipOption{
        {Value: "Frontend", Label: "Frontend", Color: core.ColorPrimary},
        {Value: "Backend", Label: "Backend", Color: core.ColorAccent},
    })
```

### Backend-Handler für Inline Edits

```go
func (ctrl *Controller) InlineEditSave(c echo.Context) error {
    wc := webcontext.GetWebContext(c)
    tbl := buildTable(wc.UiContext())
    req, err := tbl.ParseInlineEdit(c)  // validiert Field existiert + editable
    if err != nil {
        return wc.BadRequest(err.Error())
    }
    // req.ID, req.Field, req.Value — business logic + DB update

    // Response-Varianten:
    return c.JSON(200, response.NewReturnInlineEdit().
        WithUpdates(map[string]any{"status": "Aktiv"}))                             // Nur Felder patchen
    return c.JSON(200, response.NewReturnInlineEdit().WithRefreshTable())            // Tabelle neu laden
    return c.JSON(200, response.NewReturnInlineEdit().WithRefreshPage())             // Seite neu laden
    return c.JSON(200, response.NewReturnInlineEdit().WithGoto("/other/page"))       // Navigation
    return c.JSON(200, response.NewReturnInlineEdit().                               // Kombiniert
        WithUpdates(map[string]any{"status": "Aktiv"}).
        WithRefreshTable().
        WithMessage("Gespeichert", response.MessageSuccess))
}
```

### Typen

```go
// Request (Frontend → Backend)
type InlineEditRequest struct {
    ID    int64  `json:"id"`
    Field string `json:"field"`
    Value any    `json:"value"`
}

// Response (Backend → Frontend) — response.ReturnInlineEdit
// Methoden: WithUpdates, WithRefreshTable, WithRefreshPage, WithGoto, WithMessage

type EditableChipOption struct {
    Value string
    Label string
    Color core.Color
}
```

## Buttons in ButtonsField

```go
b.ButtonsField("actions", "", accessor).
    AddButton("edit", "edit", "/device/{id}/edit", "", core.ButtonActionDialog, core.ColorPrimary, false).
    AddButton("delete", "delete", "/device/{id}/delete", "", core.ButtonActionDialog, core.ColorError, false).
    AddMenuItem("export", "Export CSV", "/device/{id}/export", core.ButtonActionDownload, false)
// URL-Platzhalter {key} werden aus dem accessor-Map ersetzt
```

## Table Options

```go
b.SetTitle("Geräte")
b.SetPagination(true)
b.SetSearch(true)
b.SetReload(true)
b.SetCsv(true)
b.SetExcel(true)
b.SetDense(true)
b.SetBorders(true)
b.SetFooter(true)
b.SetServerSide(true)
b.SetSaveState(true)
b.SetSaveStateId("devices-table")
b.SetItemsPerPage(25)
b.SetPageSizes([]int{10, 25, 50, 100})
b.SetClass("custom-class")
b.SetDisplay("xcol-12")
b.SetMinWidth("800px")
b.SetScrollHeight("400px")
b.SetTextNoData("Keine Geräte gefunden")
b.SetEmptyState("/api/devices/empty")  // URL für EmptyState-Komponente

// Tree-Modus (Einrückung + Expand/Collapse; Zeilen bleiben flach, Baum aus id/parentId)
b.Tree(table.TreeConfig{
    IdField:         "id",
    ParentIdField:   "parentId",      // NICHT via .Hide() verstecken → als IdField/Format 'id' ausgeben
    TreeColumn:      "name",          // optional, Default: erste Spalte
    PersistStateKey: "devices-tree",  // optional: Expand-State in localStorage
    AddSubURL:       xurl.NewUrl("/api/devices/add?parent={id}"), // optional: "+ Sub"-Button
})  // Details: references/tables.md → "Tree-Modus"

// Select-Buttons (Multi-Select Actions)
b.SetSelect(true)
b.AddMultiEditButton("edit", "Bearbeiten", "/api/devices/multi-edit")
b.AddMultiDeleteButton("delete", "Löschen", "/api/devices/multi-delete")
b.AddMultiEditAndDeleteButtons(editBtn, deleteBtn)

// Top Buttons
b.SetButtonsTop([]*button.Button{
    button.NewSimpleDialogButton("add", "/api/devices/add", core.ColorPrimary),
})
```

## Filter

```go
// Filter mit FormGroup
filterFields := []field.FormField{
    field.NewTextField("search", "filter.search", false, nil),
    field.NewSelectField("status", "filter.status", false, statusOptions),
}
filterGroup := group.NewFormGroup(filterFields)
b.SetFilter(filterGroup)
b.SetFlags("search")  // UI-only Filter-Felder (nicht an DB gesendet)
```

## Build & Verwenden

```go
t, err := b.Build()
if err != nil {
    return err
}
```

### Statische Daten

```go
t.SetData(devices)  // []Device
result := t.Print(ctx)
```

### Server-Side Pagination (Standard-Pattern)

```go
func HandleDeviceTable(c echo.Context) error {
    ctx := getUiContext(c)
    t := buildDeviceTable()

    if err := t.LoadFilterData(c); err != nil {
        return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
    }

    pp := t.LoadPaginationParams()
    // pp.Page      — Seite (0-basiert)
    // pp.PageSize  — Einträge pro Seite
    // pp.Sort      — Sort-Feld-ID (validiert gegen definierte Spalten)
    // pp.SortDir   — "asc" oder "desc"
    // pp.Search    — Suchtext

    // DB-Query
    query := db.Model(&Device{})
    if pp.Search != "" {
        query = query.Where("name ILIKE ?", "%"+pp.Search+"%")
    }

    // Filter-Werte auslesen
    if t.HasFilter() {
        filterValues := t.GetFilterValues()
        if status, ok := filterValues["status"]; ok && status != nil {
            query = query.Where("status = ?", status)
        }
    }

    var totalCount int64
    query.Count(&totalCount)

    if pp.Sort != "" {
        query = query.Order(fmt.Sprintf("%s %s", pp.Sort, pp.SortDir))
    }

    var devices []Device
    query.Offset(pp.Page * pp.PageSize).Limit(pp.PageSize).Find(&devices)

    t.SetData(devices)
    return t.ToServerSideResponse(ctx, int(totalCount)).DataResponse().Send(c)
}
```

### Table-Seite mit Card

```go
func HandleDevicePage(c echo.Context) error {
    ctx := getUiContext(c)
    t := buildDeviceTable()

    card := card.NewCard(core.CardTypeTable, t, "devices.list", "", "", "", true, false, "")
    card.SetURL("/api/devices/table")  // AJAX-Modus
    card.ButtonTop(button.NewSimpleDialogButton("add", "/api/devices/add-dialog", core.ColorPrimary))

    p := page.NewPage()
    p.Bread("Geräte", "", false)
    p.Add(card)
    return c.JSON(http.StatusOK, p.Print(ctx))
}
```

## IconSet

```go
var statusIcons = table.NewIconSet()

var (
    IconOnline  = statusIcons.Add("check_circle", core.ColorSuccess, "Online", nil)
    IconOffline = statusIcons.Add("cancel", core.ColorError, "Offline", nil)
    IconUnknown = statusIcons.Add("help", core.ColorGray, "Unbekannt", nil)
)

// Verwendung im Field:
b.IconFieldFromSet("status", "device.status", func(d Device) *table.IconRef {
    switch d.Status {
    case 1: return IconOnline
    case 2: return IconOffline
    default: return IconUnknown
    }
}, statusIcons)
```

## Output-Typen

Die Tabelle unterstützt automatisch Web, CSV und Excel Export. Der OutputType wird über `_csv`/`_excel` Flags im Request-Body gesteuert.

```go
// CSV: Semikolon-Trenner, Formel-Injection-Schutz
// Excel: excelize Library, Auto-Spaltenbreite
// PDF: Benutzerdefiniert (OutputPDF)
```

Wenn `Csv: &true` / `Excel: &true` in den Tabellen-Optionen gesetzt ist und die Tabelle eine `URL` hat, generiert der Builder automatisch passende Download-Buttons im `ButtonsTop`-Bereich. Diese Auto-Buttons setzen das Flag intern via `button.WithData(map[string]any{"_csv": true})` — beim Klick merged das Frontend (`actionDownload`) `button.data` mit der Filter-Data und postet sie an die Tabellen-URL. Wer eigene Download-Buttons baut, sollte denselben Weg nutzen (siehe `components.md` → "Custom-Payload via WithData").

## Response-Methoden

```go
// Server-Side mit Pagination
t.ToServerSideResponse(ctx, totalCount).DataResponse().Send(c)

// Vollständige Response (mit Fields, Footer, Components)
t.ToTableDataResponse(ctx).WithFooter().WithTotalCount(total).DataResponse().Send(c)

// Direkt als Komponente (statisch)
t.Print(ctx)  // map[string]any für Einbettung in Page/Card
```
