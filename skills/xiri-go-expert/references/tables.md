# Tabellen — das eine Referenzdokument

Alles was man für eine xiri-go-Tabelle braucht: Aufbau, Spalten-Typen, Buttons (Top / Bulk / Row), Inline-Edit, Filter, Pagination, Exports. Für reine Filter-Parsing-Details siehe `table-filtering.md` (dieses File verweist an den richtigen Stellen hin).

## Import-Header

```go
import (
    "github.com/xiriframework/xiri-go/component/table"
    "github.com/xiriframework/xiri-go/component/button"
    "github.com/xiriframework/xiri-go/component/core"
    xurl "github.com/xiriframework/xiri-go/component/url"
)
```

## Grundgerüst

```go
b := table.NewBuilder[Device]()    // generisch über Row-Typ

b.SetTitle("Geräte")               // optional — Title in der Toolbar
b.IdField("id", "ID", func(r Device) int64 { return r.ID })
b.TextField("name", "Name", func(r Device) string { return r.Name })

tbl := b.Build()                   // *Table[Device]
tbl.SetURL(c.apiUrl("data"))       // POST-Target für Filter/Pagination-Requests
```

Der Builder (`TableBuilder[T]`) konfiguriert Struktur und Options; `Build()` gibt eine `*Table[T]`, die man im Handler dann mit `SetData(rows)` füttert und via `wc.Data(tbl)` ausliefert.

## Spalten-Typen — vollständige Liste

Alle Field-Konstruktoren haben die **gleiche Signatur**: `(id, name string, accessor func(T) <Typ>) *FieldBuilder`. Der zurückgegebene `*FieldBuilder` erlaubt Chaining für Modifier (siehe unten).

### Skalar-Typen

| Methode                                | Accessor-Rückgabe | Rendering                              |
| -------------------------------------- | ----------------- | -------------------------------------- |
| `IdField`                              | `int64`           | ID-Spalte (klein, grau)                |
| `IntField`                             | `int`             | Zahl, locale-formatiert                |
| `Int32Field`                           | `int32`           | wie IntField                           |
| `Int64Field`                           | `int64`           | wie IntField                           |
| `FloatField`                           | `float64`         | Zahl mit Dezimalen                     |
| `TextField`                            | `string`          | Text                                   |
| `BoolField`                            | `bool`            | ✓/✗ oder `WithBoolText`-Override       |
| `DateTimeField`                        | `time.Time`       | Datum + Zeit (via UiContext-Locale)    |
| `DateField`                            | `time.Time`       | Nur Datum                              |
| `TimeLengthField`                      | `int64` (Sekunden)| Dauer z.B. "2h 15min"                  |
| `DistanceField`                        | `float64` (m)     | Distanz (km/mi je nach UiContext)      |
| `SpeedField`                           | `float64` (m/s)   | Geschwindigkeit                        |
| `PressureField`                        | `float64` (bar)   | Druck, konvertiert nach bar/psi/kPa je UiContext |
| `HtmlField`                            | `string`          | RAW HTML (keine Escape-Logik!)         |
| `HeaderField`                          | `string`          | Gruppierungs-Header in der Zeile       |

### Text2-Typen (zwei Werte übereinander in einer Zelle)

Für "main value + sub value" in einer Spalte (typischer UI-Pattern: grosser Wert + kleiner Untertitel):

| Methode                        | Accessor-Rückgabe         |
| ------------------------------ | ------------------------- |
| `Text2Field`                   | `[2]string`               |
| `Text2IntField`                | `[2]int`                  |
| `Text2FloatField`              | `[2]float64`              |
| `Text2DateTimeField`           | `[2]time.Time`            |
| `Text2DateField`               | `[2]time.Time`            |
| `Text2BoolField`               | `[2]bool`                 |
| `Text2DistanceField`           | `[2]float64`              |
| `Text2SpeedField`              | `[2]float64`              |
| `Text2TimeLengthField`         | `[2]int64`                |

### N-Typen (Array in einer Zelle)

Für "Liste von Werten in einer Zelle" (z.B. mehrere Messungen, Tags):

| Methode                  | Accessor-Rückgabe |
| ------------------------ | ----------------- |
| `TextNField`             | `[]string`        |
| `IntNField`              | `[]int`           |
| `FloatNField`            | `[]float64`       |
| `DateTimeNField`         | `[]time.Time`     |
| `DateNField`             | `[]time.Time`     |
| `DistanceNField`         | `[]float64`       |
| `SpeedNField`            | `[]float64`       |
| `BoolNField`             | `[]bool`          |
| `TimeLengthNField`       | `[]int64`         |

### Spezial-Typen

#### `LinkField` — anklickbarer Link

```go
b.LinkField("detail", "Link", func(r Device) [2]string {
    return [2]string{r.Name, c.pageUrl("detail", strconv.FormatInt(r.ID, 10)).Print()}
})
// [2]string{text, url}  — text wird gerendert, url ist das Link-Target
```

#### `InputField` — generische Edit-Zelle

Nur sinnvoll mit `WithEditable(true)` — gibt dem Frontend einen raw-input mit dem aktuellen Wert (beliebiger Typ). Für typisierte Inline-Edit lieber einen der normalen Skalarfelder + `WithEditable(true)`.

```go
b.InputField("custom", "Custom", func(r Device) any { return r.Custom })
```

#### `IconFieldFromSet` — Status-Icon aus einem Icon-Set

```go
icons := table.NewIconSet()
iconOn := icons.Add("on", "check_circle", core.ColorSuccess, "Aktiv")
iconOff := icons.Add("off", "cancel", core.ColorError, "Inaktiv")

b.IconFieldFromSet("status", "Status", func(r Device) *table.IconRef {
    if r.Active { return iconOn }
    return iconOff
}, icons)
```

Der Icon-Set wird einmal gebaut, die Zeilen referenzieren die Icons per Pointer → kompaktere JSON, weniger Duplizierung.

#### `ButtonsField` — Row-Action-Buttons (siehe Abschnitt „Row Buttons")

## FieldBuilder-Modifier (für jede Spalte verfügbar)

Jeder Field-Konstruktor gibt ein `*FieldBuilder` zurück, an dem man fluent beliebig viele Modifier chained:

### Layout / Darstellung

```go
b.TextField("name", "Name", accessor).
    WithWidth("200px").           // fixe Breite
    WithMinWidth("120px").
    WithAlign(table.FieldAlignRight).  // oder AlignLeft()/AlignCenter()/AlignRight() Shortcuts
    WithSticky(true).              // bleibt beim horizontalen Scrollen
    WithDisplay("no-wrap").        // CSS-Klasse am <td>
    WithHint("Gerätename").        // Tooltip
    WithHeader("NAME").            // Custom Header (z.B. abgekürzt)
    WithHeaderSpan(2).             // Header-colspan
    WithColumnOrder(5).            // Reihenfolge überschreiben
    WithTextPrefix("€ ").
    WithTextSuffix(" km")
```

### Sortierung / Suche

```go
b.TextField("name", "Name", accessor).
    WithSort(true).       // Server-Side-Sort aktivieren (Default: true wenn indexbar)
    WithSearch(true)      // ist im pg.Search-Scope
```

### Sichtbarkeit

```go
b.TextField("debug", "Debug", accessor).
    Hide()                // komplett versteckt
    
b.TextField("notes", "Notizen", accessor).
    HideInCSV()           // in Web sichtbar, in CSV nicht
    
b.TextField("barcode", "Barcode", accessor).
    ShowInCSV()           // nur beim CSV-Export
```

### Formatter (Custom Rendering)

```go
b.FloatField("rate", "Rate", accessor).
    WithDecimals(3).                                          // Float-Nachkommastellen
    WithFormatter(func(v any) string { return fmt.Sprintf("%.2f %%", v) })

b.TextField("status", "Status", accessor).
    WithWebFormatter(myWebFormatter).
    WithCSVFormatter(myCSVFormatter).                         // getrennt pro Output-Typ
    WithExcelFormatter(myExcelFormatter).
    WithPDFFormatter(myPDFFormatter)

b.BoolField("active", "Aktiv", accessor).
    WithBoolText("Ja", "Nein")                                // Default: ✓/✗
```

### Footer-Aggregation

Jede numerische Spalte kann in der Footer-Row einen Summenwert tragen, wenn die Tabelle `SetFooter(true)` hat:

```go
b.FloatField("price", "Preis", accessor).WithFooterSum()
b.TextField("name", "Name", accessor).WithFooterCount()    // Zählt Zeilen
b.TextField("status", "Status", accessor).WithFooter(table.FieldFooterNo)
```

### Inline-Edit

```go
b.TextField("status", "Status", accessor).
    WithEditable(true).                    // Simple-Text-Edit
    WithInputType("text").                 // "text" | "email" | "tel" | "url" | "number"
    WithInputRequired(true).
    WithInputLang("de").
    WithInputPaste(true)

b.TextField("prio", "Prio", accessor).
    WithEditableOptions(map[string]string{
        "low":  "Niedrig",
        "high": "Hoch",
    })                                     // Dropdown mit statischen Options

b.TextField("owner", "Owner", accessor).
    WithEditableOptionsUrl("/api/owners")  // Dropdown: Options vom Backend (GET, einmalig)

b.TextField("status", "Status", accessor).
    WithEditableOptions(opts).
    WithEditableOptionsSearch(true)        // Suchfeld, client-seitig (lokal filtern)

b.TextField("owner", "Owner", accessor).
    WithEditableSearchOptionsUrl("/api/owners/search") // Suchfeld, SERVER-seitig
                                           // POST {id, field, search} → [{value,label,color?}]

b.TextField("tags", "Tags", accessor).
    WithEditableChipOptions([]table.EditableChipOption{
        {Value: "fe", Label: "Frontend", Color: core.ColorPrimary},
        {Value: "be", Label: "Backend",  Color: core.ColorAccent},
    })                                     // Chip-MultiSelect (Suche via WithEditableOptionsSearch/WithEditableSearchOptionsUrl kombinierbar)
```

Die Tabelle selbst braucht dafür `b.SetEditUrl(c.apiUrl("inline-edit").PrintPrefix())` — das POST-Target für `tbl.ParseInlineEdit(ctx)`.

### Zugriff / Berechtigung

```go
b.TextField("salary", "Gehalt", accessor).
    WithAccess([]string{"admin", "hr"})    // Rollen-Metadaten für das Frontend
```

> ⚠️ **`WithAccess` ist kein Zugriffsschutz.** Die Library wertet die Rollen **nirgends** aus — es
> gibt keine Backend-Prüfung. Die Werte landen als Metadaten im JSON-Modell; ob und wie das Frontend
> sie berücksichtigt, ist nicht Sache der Library. Wer eine Spalte wirklich schützen will, darf sie
> **gar nicht erst in die Tabelle aufnehmen** (bzw. den Accessor rollenabhängig leer liefern lassen)
> und muss den Zugriff im eigenen Handler prüfen. Ein Klient sieht sonst den Wert im Response.
>
> (Bekannt als Finding #2 des Audits; eine echte Durchsetzung steht noch aus.)

## Row Buttons (Aktionen pro Zeile)

`ButtonsField` liefert Aktionen pro Zeile (Edit-Pen, Delete-Icon, Menü, etc.). Der Accessor gibt eine `map[string]string` zurück: Key = Button-Index als String (`"0"`, `"1"`, …), Value = URL für diesen Button in dieser Zeile.

```go
b.ButtonsField("actions", "", func(r Device) map[string]string {
    return map[string]string{
        "0": c.apiUrl("edit",   strconv.FormatInt(r.ID, 10)).PrintPrefix(),
        "1": c.apiUrl("delete", strconv.FormatInt(r.ID, 10)).PrintPrefix(),
        "2": c.pageUrl("detail", strconv.FormatInt(r.ID, 10)).Print(),
    }
}).
    AddButton(0, table.FieldButtonActionDialog,   "edit",   core.ColorPrimary, "Bearbeiten").
    AddButton(1, table.FieldButtonActionDialog,   "delete", core.ColorWarning, "Löschen").
    AddButton(2, table.FieldButtonActionLink,     "visibility", core.ColorAccent, "Ansehen")
```

Der Key ist ein Positions-Index und muss zwischen `0` und `1000` liegen — er wird bei der
Serialisierung als Slice-Index benutzt. Keys außerhalb werden ignoriert und per `slog.Warn`
gemeldet (früher: Panic bei negativem, Riesenallokation bei sehr großem Key).

### Datei anzeigen statt herunterladen

`WithButtonTarget(key, "_blank")` an einem Download-Zellen-Button lässt das Frontend die Datei in
einem neuen Tab **anzeigen** statt sie zu speichern — der Fall „PDF dieser Zeile ansehen":

```go
b.ButtonsField("actions", "", func(r Device) map[string]string {
    return map[string]string{"0": c.apiUrl("pdf", strconv.FormatInt(r.ID, 10)).PrintPrefix()}
}).
    AddButton(0, table.FieldButtonActionDownload, "picture_as_pdf", core.ColorPrimary, "PDF").
    WithButtonTarget(0, "_blank")
```

Voraussetzung ist ein Content-Type, den der Browser rendern kann (`application/pdf`);
`Content-Disposition` spielt keine Rolle, weil das Frontend den Blob selbst baut. Unbekannte Keys
ignoriert der Setter, er ist also nach einem verworfenen Out-of-Range-Key gefahrlos.
Für Table-Top-/Bulk-Buttons heißt das Pendant `TableButton.WithTarget("_blank")`.

### `FieldButtonAction`-Werte

| Enum                              | Verhalten                                              |
| --------------------------------- | ------------------------------------------------------ |
| `FieldButtonActionLink`           | Angular-Route (`routerLink`)                           |
| `FieldButtonActionHref`           | `<a href>` (externer Link)                             |
| `FieldButtonActionApi`            | POST, keine Response-Action                            |
| `FieldButtonActionDialog`         | POST → MatDialog öffnen mit Response                   |
| `FieldButtonActionForm`           | POST → öffnet Form-Dialog                              |
| `FieldButtonActionDownload`       | POST → Blob-Download (mit `WithButtonTarget(key, "_blank")` im Tab anzeigen) |
| `FieldButtonActionGet/Post/Put/Delete` | Entsprechende HTTP-Methode                         |
| `FieldButtonActionSave` / `Close` / `Back` | Spezielle Flow-Actions                         |
| `FieldButtonActionMenu`           | Öffnet Menü mit weiteren Items (siehe unten)           |

### Mit Row-Hint (Tooltip pro Zeile)

```go
b.ButtonsField("actions", "", accessor).
    AddButton(0, table.FieldButtonActionDialog, "delete", core.ColorWarning, "Löschen")

// Generisch (spezialisiert an Table-Typ):
table.WithRowHint[Device](fieldBuilder, func(r Device) string {
    return "Gerät " + r.Name + " löschen"
})
```

### Menu-Button (Dropdown in Row)

Wenn mehr als 2-3 Actions pro Zeile: Menü einklappen.

```go
fb := b.ButtonsField("menu", "", func(r Device) map[string]string {
    return map[string]string{
        "0": c.pageUrl("detail", strconv.FormatInt(r.ID, 10)).Print(),
    }
}).AddButton(0, table.FieldButtonActionMenu, "more_vert", core.ColorPrimary, "Aktionen")

fb.AddMenuItem(table.FieldButtonActionLink,   "edit",   core.ColorPrimary, "Bearbeiten")
fb.AddMenuItem(table.FieldButtonActionDialog, "delete", core.ColorWarning, "Löschen")
fb.AddMenuItem(table.FieldButtonActionApi,    "archive",core.ColorAccent,  "Archivieren")
```

### AddMenu mit Row-abhängigen Items

```go
// Items sind pro Zeile unterschiedlich (z.B. nur bei Online-Geräten zeigen)
table.AddMenu[Device](fb, 0, "more_vert", core.ColorPrimary, "Aktionen",
    func(r Device) []string {
        if r.Online {
            return []string{"restart", "shutdown"}
        }
        return []string{"start"}
    })
```

## Top Buttons (Toolbar über der Tabelle)

`SetButtonsTop` nimmt eine Liste von `*button.TableButton`. Typischer Use-Case: „Neu", „Importieren", „Export", „Filter reset".

```go
b.SetButtonsTop([]*button.TableButton{
    button.NewTableButton(
        core.ButtonActionLink,
        "add",                       // Material-Icon-Name
        c.pageUrl("add"),
        "Neues Gerät",
        core.ColorPrimary,
        false,                        // disabled?
        nil,                          // options
    ),
    button.NewTableButton(
        core.ButtonActionDialog,
        "upload",
        c.apiUrl("import"),
        "CSV importieren",
        core.ColorAccent,
        false,
        nil,
    ),
})
```

Im Gegensatz zu `pageheader.New(...).Buttons(...)` erscheinen Top-Buttons **direkt in der Tabellen-Toolbar** (neben Suche/Pagination) statt im Seitenkopf. Welcher Platz richtig ist, hängt vom UI-Design: globale Seiten-Actions → PageHeader; tabellen-spezifische Actions (nur wenn diese Tabelle sichtbar ist) → Top-Buttons.

## Select Buttons (Bulk-Actions über Row-Selection)

Bulk-Actions werden **unten** (je nach Frontend-Config auch oben) angezeigt, sobald mindestens eine Zeile selektiert wurde. `SetSelectButtons` aktiviert **automatisch** die Row-Checkbox-Spalte — man muss `SetSelect(true)` nicht extra aufrufen.

```go
b.SetSelectButtons([]*button.TableButton{
    button.NewTableButton(
        core.ButtonActionDialog,
        "delete",
        c.apiUrl("bulk-delete"),
        "Löschen",
        core.ColorWarning,
        false,
        nil,
    ),
    button.NewTableButton(
        core.ButtonActionDialog,
        "edit",
        c.apiUrl("bulk-edit"),
        "Gemeinsam bearbeiten",
        core.ColorPrimary,
        false,
        nil,
    ),
})
```

### Standard-Kombis (Convenience-Methoden)

```go
b.AddMultiDeleteButton(c.apiUrl("bulk-delete").PrintPrefix())
// → dialog+"delete"+warn+"LOESCHEN"

b.AddMultiEditButton(c.apiUrl("bulk-edit").PrintPrefix())
// → dialog+"edit"+primary+"BEARBEITEN"

b.AddMultiEditAndDeleteButtons(editUrl, deleteUrl)
// = AddMultiEditButton + AddMultiDeleteButton
```

### Bulk-Handler im Controller

Das Frontend sendet bei Button-Click ein Array von IDs. Beispiel-Handler:

```go
func (c *Controller) BulkDelete(ctx echo.Context) error {
    wc := webcontext.GetWebContext(ctx)
    var req struct {
        IDs []int64 `json:"ids"`
    }
    if err := ctx.Bind(&req); err != nil {
        return wc.BadRequest(err.Error())
    }
    if err := c.svc.DB.Device.DeleteMany(ctx.Request().Context(), req.IDs); err != nil {
        return wc.InternalServerError(err.Error())
    }
    return wc.Component(response.NewReturnRefreshTable().
        WithMessage(fmt.Sprintf("%d gelöscht", len(req.IDs)), response.MessageSuccess))
}
```

### `ClearSelectButtons` + `AddSelectButton`

```go
b.AddSelectButton(btn)       // einzeln append + Auto-Select true
b.ClearSelectButtons()       // leert + Auto-Select false
```

## Table-Options (Übersicht aller Setter)

```go
b.SetTitle("Geräte")
b.SetClass("dense-rows")                       // CSS-Klasse auf die Tabelle
b.SetDisplay("xcol-md-12")                     // Grid-Layout-Klasse
b.SetTextNoData("Keine Geräte")
b.SetEmptyState(emptystateComponent)           // ganze EmptyState-Komponente

// Features aktivieren/deaktivieren (bool-Pointer, Default: meist true)
b.SetReload(true)           // manueller Reload-Button (Icon in der Toolbar)
b.SetDense(true)            // enger Zeilen-Spacing
b.SetPagination(true)
b.SetSearch(true)
b.SetQuery(false)           // Filter-Panel oben automatisch
b.SetCsv(true)              // CSV-Export-Button verfügbar
b.SetExcel(true)            // Excel-Export
b.SetSaveState(true)        // Filter/Sort/Page persistieren im LocalStorage
b.SetSaveStateId("device-table")
b.SetBorders(true)
b.SetBordersHeader(true)
b.SetSelect(true)           // Row-Checkboxen (SetSelectButtons macht das auto)
b.SetFooter(true)           // Summenzeile anzeigen
b.SetScrollHeight("600px")
b.SetMinWidth("1200px")

// Server-Side
b.SetServerSide(true)
b.SetItemsPerPage(50)
b.SetPageSizes([]int{10, 25, 50, 100, 500})

// Inline-Edit
b.SetEditUrl("/api/devices/inline-edit")

// Zusatz-Inputs
b.SetSaveInput("note")      // zusätzliches Input-Feld in der Toolbar
b.SetSaveInputUrl("/api/devices/note")
```

Fast alle Setter geben `*TableBuilder[T]` zurück, d.h. man kann alles in einer Chain hängen:

```go
tbl := table.NewBuilder[Device]().
    SetTitle("Geräte").
    SetServerSide(true).
    SetSaveState(true).
    SetSaveStateId("devices").
    IdField("id", "ID", func(r Device) int64 { return r.ID }).
    TextField("name", "Name", func(r Device) string { return r.Name }).WithSticky(true).
    SetButtonsTop(topButtons).
    SetSelectButtons(bulkButtons).
    Build()
```

## Auto-Refresh (Polling während Background-Worker)

Solange ein Background-Worker für Zeilen der Tabelle läuft, kann die Tabelle sich
**selbsttätig** im Intervall neu laden — ohne dass der User manuell aktualisiert. Der
Status steht dabei in den normalen Spalten; das Frontend zeigt tabellen-weit im Header
„Auto-Refresh aktiv · nächste in Xs" mit Countdown. Sobald die Response **kein** `poll`
mehr enthält, stoppt das Polling automatisch.

Gesetzt wird das **pro Request** (zustandsabhängig), nicht am Builder:

```go
// im Daten-Handler, nachdem die Rows geladen sind:
tbl.SetData(rows)
if anyJobRunning {
    tbl.SetPoll(2000)   // ms — fließt über wc.Data(tbl) als "poll":2000 in die Response
}
return wc.Data(tbl)
```

`SetPoll(ms)` liegt auf der gebauten `*Table[T]` und wird über
`ToTableDataResponse` / `DataResponse` (Web-Output) durchgereicht. Alternativ direkt auf
der Daten-Response: `tbl.ToTableDataResponse(ctx).WithPoll(2000)`. Für CSV/Excel wird `poll`
nie ausgegeben. Worker fertig → `SetPoll` einfach **nicht** setzen → Response ohne `poll`
→ Frontend stoppt.

## Tree-Modus (Hierarchie mit Einrückung + Expand/Collapse)

Opt-in über `b.Tree(table.TreeConfig{...})`. Die Zeilen bleiben **flach** — das Frontend
baut den Baum aus `IdField`/`ParentIdField` jeder Zeile und rendert Einrückung,
Expand/Collapse-Pfeile, Expand-All/Collapse-All, hierarchie-bewusste Suche (Treffer + deren
Subtree, Vorfahren als gedimmter Kontext) und optional einen „+ Sub"-Button pro Zeile.
Ohne `Tree(...)` verhält sich die Tabelle unverändert (identisches JSON).

```go
b := table.NewBuilder[Region]()
b.IdField("id", "ID", func(r Region) int64 { return r.ID })
b.IdField("parentId", "Parent", func(r Region) int64 { return r.ParentID }) // s. Falle unten
b.TextField("name", "Region", func(r Region) string { return r.Name })

b.Tree(table.TreeConfig{
    IdField:            "id",        // Pflicht: Zeilen-Feld mit Knoten-ID
    ParentIdField:      "parentId",  // Pflicht: Parent-ID; null/0 → Root
    TreeColumn:         "name",      // optional, Default: erste Spalte (rendert die Einrückung)
    CollapseAllDefault: false,       // optional, Default: false → Tree startet voll ausgeklappt
    HideCounts:         false,       // optional, Default: false → Kind-Count "(5)" bei collapsed
    PersistStateKey:    "regions",   // optional: Expand-State im localStorage persistieren
    AddSubURL:          xurl.NewUrl("/Portal/Region/Add?parent={id}"), // optional: "+ Sub"-Button; {id} wird ersetzt
})
tbl := b.Build()
tbl.SetData(regions)
```

- **„+ Sub"-Button pro Zeile steuern:** Standardmäßig erscheint der „+" (bei gesetztem `AddSubURL`)
  auf jeder Zeile. Mit `b.TreeAddSubWhen(func(r Row) bool {...})` nur dort, wo der Accessor `true`
  liefert. Intern wird das Flag pro Zeile als `_addSub` emittiert (`Tree.AddSubField`); `.Tree(...)`
  muss vorher aufgerufen sein.
- Multi-Root wird unterstützt; Knoten mit fehlendem Parent werden als Root behandelt, Zyklen
  abgefangen (Knoten wird Root, Warn-Log).
- Sortierung anderer Spalten ist im Tree-Modus deaktiviert; Geschwister werden alphabetisch
  nach der Tree-Spalte sortiert.
- **Falle:** `ParentIdField` darf **nicht** über `.Hide()` versteckt werden — versteckte
  Felder fallen aus den Row-Daten (`GetData` überspringt sie), und der Baum bricht. Stattdessen
  als `IdField`/Format `id` ausgeben: aus der Anzeige gefiltert, aber in den Daten vorhanden.

## Handler-Flow (Complete)

### Page-Route (GET /devices)

```go
func (c *Controller) Page(ctx echo.Context) error {
    wc := webcontext.GetWebContext(ctx)
    uc := wc.UiContext()

    tbl := c.buildTable(uc)
    tbl.SetURL(c.apiUrl("data"))

    p := page.NewPage()
    p.Bread("Devices", c.pageUrl(), false)
    p.Add(pageheader.New("Geräte"))
    p.Add(tbl)
    return wc.Page(p)
}
```

### Data-Route (POST /api/devices/data)

```go
func (c *Controller) Data(ctx echo.Context) error {
    wc := webcontext.GetWebContext(ctx)
    uc := wc.UiContext()

    tbl := c.buildTable(uc)
    tbl.SetFilter(c.buildFilterGroup(uc))      // Filter-FormGroup attachen

    filters, err := tbl.LoadFilterData(ctx)
    if err != nil { return wc.BadRequest(err.Error()) }
    pg := tbl.LoadPaginationParams()

    rows, total, err := c.svc.FindDevices(filters, pg)
    if err != nil { return wc.InternalServerError(err.Error()) }

    tbl.SetData(rows)
    tbl.SetTotal(int(total))
    return wc.Data(tbl)
}
```

Details zu Filter-Parsing (welcher Go-Typ pro Field, GORM-Integration, Multi-Select): **`references/table-filtering.md`**.

### Inline-Edit-Route

```go
func (c *Controller) InlineEdit(ctx echo.Context) error {
    wc := webcontext.GetWebContext(ctx)
    tbl := c.buildTable(wc.UiContext())

    req, err := tbl.ParseInlineEdit(ctx)
    if err != nil { return wc.BadRequest(err.Error()) }

    d, err := c.svc.DB.Device.GetByID(req.ID)
    if err != nil { return wc.NotFound(err.Error()) }

    switch req.Field {
    case "status":   d.Status   = req.Value.(string)
    case "priority": d.Priority = req.Value.(string)
    default:
        return wc.BadRequest("unknown field")
    }
    if err := c.svc.DB.Device.Update(ctx.Request().Context(), d); err != nil {
        return wc.InternalServerError(err.Error())
    }

    return wc.Component(response.NewReturnInlineEdit().
        WithUpdates(map[string]any{req.Field: req.Value}).
        WithMessage("Gespeichert", response.MessageSuccess))
}
```

## Dynamische Spalten ein-/ausblenden (runtime)

```go
tbl.HideField("secret")
tbl.ShowField("debug")
tbl.HideFields("foo", "bar")
tbl.ShowFields("foo", "bar")
```

Wenn die Field-Liste abhängig vom Request variieren muss (z.B. Role-Based-Columns), kann man `b.SetFieldsCanChange()` setzen — dann wird das Field-Export nicht gecached.

## Flags — UI-State, der nicht als Filter zählt

```go
b.SetFlags("_viewMode", "_showArchived")
```

Keys die im Request-Body als Flags erkannt werden, landen **nicht** in der Filter-Map (`LoadFilterData`-Return). Sie bleiben aber in `tbl.GetFilterData()` — du kannst sie separat abfragen:

```go
raw := tbl.GetFilterData()
if v, ok := raw["_viewMode"].(string); ok && v == "archived" {
    q = q.Unscoped().Where("deleted_at IS NOT NULL")
}
```

## Exports — CSV & Excel

Wenn `SetCsv(true)` / `SetExcel(true)`, hängt das Frontend `_csv: true` bzw. `_excel: true` an den Data-Request. `LoadFilterData` setzt dann `tbl.GetOutputType()` entsprechend.

```go
tbl.SetData(rows)
switch tbl.GetOutputType() {
case table.OutputCSV:
    return wc.CsvFromTable(tbl, "devices.csv")
case table.OutputExcel:
    return wc.ExcelFromTable(tbl, "devices.xlsx")
default:
    return wc.Data(tbl)
}
```

Für Exports macht man typischerweise **keine** Pagination (alle Zeilen) — die App muss selbst entscheiden, ob das erlaubt ist oder ob ein Limit nötig ist.

## Row-Typ mit FK-Auflösung (Lookup-Maps)

Wenn die Tabelle FKs anzeigt, lohnt sich ein separater Row-Typ **mit** aufgelösten Namen statt lazy-Loading pro Zeile:

```go
type Row struct {
    Device
    GroupName string
    OwnerName string
}

func (c *Controller) Data(ctx echo.Context) error {
    // ... Filter/Pagination ...
    devices := c.svc.DB.Device.FindMany(filters, pg)
    groups  := c.svc.DB.Group.AllAsMap()      // map[int64]string
    owners  := c.svc.DB.User.AllAsMap()       // map[int64]string

    rows := make([]Row, 0, len(devices))
    for _, d := range devices {
        rows = append(rows, Row{
            Device:    d,
            GroupName: groups[d.GroupID],
            OwnerName: owners[d.OwnerID],
        })
    }
    tbl.SetData(rows)
    return wc.Data(tbl)
}
```

Table-Builder nutzt dann `TextField` für `GroupName` / `OwnerName` — keine IntField/ModelField, weil der Name bereits aufgelöst ist.

## Häufige Fallen

- **`ButtonsField`-Keys müssen Strings sein:** `"0"`, `"1"` — **nicht** `0`, `1`. Go's `map[string]string` ist zwingend.
- **`SetURL` erwartet `*xurl.Url`, nicht `string`.** Für ein manuelles URL-Konstrukt: `xurl.NewUrlPrefix("/data", "/api/v1")`.
- **Inline-Edit-URL:** `SetEditUrl` nimmt einen `string` (nicht `*xurl.Url`) — hier explizit `c.apiUrl("inline").PrintPrefix()` übergeben.
- **Filter + Multi-Select:** `NewSelectField(..., options).SetMultiple(true)`. Es gibt **kein** `NewMultiSelectField`.
- **SelectButtons aktivieren Select automatisch** — `SetSelect(true)` ist redundant nach `SetSelectButtons(...)`.
- **`SetServerSide(true)` ohne `SetTotal` ergibt falsche Pagination.** Immer Total aus der DB-Query mitschicken.
- **`WithEditableOptionsUrl` erwartet eine JSON-Response-Struktur** der Form `[{value, label, color?}]` — das Backend muss die Options als Liste liefern, nicht als Map.
