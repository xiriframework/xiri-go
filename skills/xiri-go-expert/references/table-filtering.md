# Table-Filter & Pagination

Der Filter-Workflow ist der kritischste Pfad in jedem xiri-go-Controller, der eine Tabelle liefert. Diese Datei deckt **Filter aufbauen → parsen → in GORM-WHERE-Clauses übersetzen** ab, plus Bulk-Actions und Exports.

## Flow im Überblick

```
Request-Body (JSON)
        │
        ▼
tbl.LoadFilterData(ctx)  ──► map[string]any  (parsed + validated)
        │                     + setzt outputType (Web/CSV/Excel)
        ▼
tbl.LoadPaginationParams()  ──► {Page, PageSize, Sort, SortDir, Search}
        │                        (Sort gegen Field-IDs validiert = SQL-safe)
        ▼
DB-Query mit Filters + Pagination + Sort
        ▼
tbl.SetData(rows)
        ▼
wc.Data(tbl)  ──► JSON-Response
```

## Filter-FormGroup aufbauen

Filter nutzen dieselben Fields wie normale Forms, aber `required=false` und keine Defaults:

```go
import (
    formbuilder "github.com/xiriframework/xiri-go/form/builder"
    "github.com/xiriframework/xiri-go/form/field"
    "github.com/xiriframework/xiri-go/form/group"
)

func (c *Controller) buildFilterGroup(ctx *core.UiContext) *group.FormGroup {
    fb := formbuilder.NewFormBuilder(ctx)

    // Single-Value Filter
    fb.AddField(field.NewTextField("name", "Name", false, ""))
    fb.AddField(field.NewIntField("minCount", "Anzahl ≥", false, 0))
    fb.AddField(field.NewTimeField("from", "Von", false, 0))
    fb.AddField(field.NewSelectField("status", "Status", false, statusOptions))

    // Multi-Value Filter: SetMultiple(true)
    fb.AddField(field.NewSelectField("priority", "Priorität", false,
        priorityOptions).SetMultiple(true))
    fb.AddField(field.NewModelListField("tags", "Tags", false, "tag", nil))

    fg, _, _ := fb.BuildAdd()   // für Query-Komponente reicht BuildAdd
    return fg
}
```

Filter-FormGroup der Tabelle zuordnen:

```go
tbl := b.Build()
tbl.SetFilter(filterFg)   // aktiviert ParseAndValidate in LoadFilterData
```

Oder direkt via `NewQueryWithFormGroup` (siehe unten), das passt die Filter automatisch in die `LoadFilterData`-Pipeline.

Weil Filter dieselben Fields wie Forms rendern, erben sie auch **abhängige Felder**: ein Filter kann
per `SetReloadOn(...)` seine Optionen nachladen, sobald ein anderer Filter sich ändert (Details in
`form-fields.md`, Abschnitt „Abhängige Felder"). Verwirft der Reload einen ungültig gewordenen
Filterwert, läuft die Filter-Pipeline damit ein weiteres Mal — der Handler braucht dafür nichts
Besonderes.

## `LoadFilterData` — was drin landet

Quelle: `component/table/table.go`.

```go
func (tc *tableCore) LoadFilterData(c echo.Context) (map[string]any, error)
```

Ablauf intern:

1. Request-Body → `map[string]any` via `c.Bind(...)`.
2. `_csv` / `_excel` Flags → setzt `outputType` (wenn Options das erlauben).
3. `_csv`, `_excel` und deklarierte `flags` aus der Map entfernen.
4. **Wenn Filter gesetzt:** `filter.ParseAndValidate(data)` — pro Field ein `Parse(raw)` + `Validate(parsed)`.
5. Gibt parsed Filters zurück (ohne `_page/_pageSize/_sort/_sortDir/_search` falls kein Filter gesetzt).

**Wichtig:** Pagination-Params bleiben in `tc.filterData` — `LoadPaginationParams` liest sie **dort**. Du rufst `LoadPaginationParams()` also immer **nach** `LoadFilterData()`.

## Value-Typen nach Parse

| Field-Konstruktor                                  | Go-Typ in parsed map           |
| -------------------------------------------------- | ------------------------------ |
| `NewTextField`                                     | `string`                       |
| `NewIntField`                                      | `int32`                        |
| `NewBoolField`                                     | `bool`                         |
| `NewTimeField`                                     | `int64` (Unix-Seconds)         |
| `NewSelectField(...)` mit `int`-Options            | `int` (wie Option-Werte)       |
| `NewSelectField(...)` mit `int32`-Options          | `int32`                        |
| `NewSelectField(...)` mit `string`-Options         | `string`                       |
| `NewSelectField(...).SetMultiple(true)`            | `field.ModelListValue` (`[]int32`) |
| `NewModelField`                                    | `int32`                        |
| `NewModelListField`                                | `field.ModelListValue`         |
| `NewChipsField`                                    | `[]interface{}` (int64 Options-IDs + string Freitext) |
| `NewArrayField`                                    | `[]interface{}` (pro-Item-Typ wie Inner-Field) |
| `NewTimeRangeField`                                | `map[string]int64` (`from`, `to`) |
| `NewFileField`                                     | `string` (Path/Base64/Key)     |

Unbekannte/fehlende Keys ⇒ nicht in der Map (weil `required=false` → Field liefert `nil` oder Default).

## Safe Type-Assertions (Pattern)

Weil JSON flach serialisiert und der Parse-Typ pro Field streng ist, **immer** mit `, ok` casten und leere Werte ausfiltern:

```go
if v, ok := filters["name"].(string); ok && v != "" {
    q = q.Where("name ILIKE ?", "%"+v+"%")
}
if v, ok := filters["minCount"].(int32); ok && v > 0 {
    q = q.Where("count >= ?", v)
}
if v, ok := filters["from"].(int64); ok && v > 0 {
    q = q.Where("created_at >= ?", time.Unix(v, 0))
}
if v, ok := filters["active"].(bool); ok {
    q = q.Where("active = ?", v)
}
```

## Multi-Value: SelectField, ModelListField, Chips

```go
// SelectField.SetMultiple(true) → field.ModelListValue (alias []int32)
if v, ok := filters["priority"].(field.ModelListValue); ok && len(v) > 0 {
    q = q.Where("priority_id IN ?", []int32(v))
}

// ModelListField → field.ModelListValue
if v, ok := filters["tags"].(field.ModelListValue); ok && len(v) > 0 {
    q = q.Where("id IN (SELECT entity_id FROM entity_tags WHERE tag_id IN ?)",
        []int32(v))
}

// ChipsField → []interface{} (gemischt int64 Options-IDs + string Freitext)
// Bei Bedarf über die Element-Typen trennen:
if v, ok := filters["labels"].([]interface{}); ok && len(v) > 0 {
    var ids []int64
    var texts []string
    for _, item := range v {
        switch x := item.(type) {
        case int64:
            ids = append(ids, x)
        case string:
            texts = append(texts, x)
        }
    }
    if len(ids) > 0 {
        q = q.Where("label_id IN ?", ids)
    }
    if len(texts) > 0 {
        q = q.Where("label IN ?", texts)
    }
}
```

**Tipp:** Weil `ModelListValue` ein Named-Type auf `[]int32` ist, kannst du ihn überall dort benutzen, wo GORM `[]int32` erwartet — Go konvertiert implizit bei direktem Assignment, aber für Variadic-/Interface-Params empfiehlt sich der explizite `[]int32(v)` Cast.

## Pagination + Sort + Search

```go
pg := tbl.LoadPaginationParams()

if pg.Search != "" {
    q = q.Where("name ILIKE ? OR description ILIKE ?",
        "%"+pg.Search+"%", "%"+pg.Search+"%")
}

if pg.Sort != "" {
    q = q.Order(pg.Sort + " " + pg.SortDir)   // Sort ist SQL-safe validiert
} else {
    q = q.Order("id DESC")                    // Default-Sort
}

// Für Server-Side-Tabellen: Total vor der Pagination zählen
var total int64
q.Count(&total)

q = q.Offset(pg.Page * pg.PageSize).Limit(pg.PageSize)

var rows []Device
q.Find(&rows)
```

## Query-Komponente: Filter + Tabelle zusammen

`query.NewQueryWithFormGroup` rendert Filter-Form und darunter beliebige Komponenten (meistens die Tabelle) in einem container:

```go
import "github.com/xiriframework/xiri-go/component/query"

func (c *Controller) Page(ctx echo.Context) error {
    wc := webcontext.GetWebContext(ctx)
    uc := wc.UiContext()

    fg := c.buildFilterGroup(uc)
    tbl := c.buildTable()
    tbl.SetURL(c.apiUrl("data"))
    tbl.SetFilter(fg)   // Tabelle kennt die Filter-Fields für ParseAndValidate

    q := query.NewQueryWithFormGroup(
        fg,
        nil,             // fieldValues (Defaults — meist nil)
        c.apiUrl("data"),// POST-URL für Filter-Änderungen
        nil,             // buttonLine
        nil,             // saveStateId
        nil,             // extra map
    )
    q.WithSaveStateId("devices.filter").Collapsed(false)
    q.Add(tbl, uc)

    p := page.NewPage()
    p.Bread("Devices", c.pageUrl(), false)
    p.Add(pageheader.New("Geräte"))
    p.Add(q)
    return wc.Page(p)
}
```

Die Data-Route (`apiUrl("data")`) bleibt **gleich** wie bei nicht-filtered Tabellen — `LoadFilterData` parst automatisch die Filter-Fields.

### Filter einklappen

`Collapsed` steuert das Expansion-Panel um den Filter — drei Zustände:

| Aufruf | Ergebnis |
|---|---|
| gar nicht gesetzt | kein Panel, Filter immer offen (Default) |
| `Collapsed(false)` | Panel, aufgeklappt |
| `Collapsed(true)` | Panel, eingeklappt |

Wer keine eigene Query baut, sondern die Tabelle über `SetFilter` **automatisch** wrappen lässt,
setzt dasselbe am TableBuilder:

```go
builder.SetFilter(fg).SetFilterCollapsed(true)
```

Klappt der User das Panel selbst auf oder zu, merkt sich das Frontend den Zustand unter demselben
`saveStateId` wie die Filterwerte (Session-Storage, 1 h) — beim nächsten Laden gewinnt der
gemerkte Zustand über den hier gesetzten Wert. Ohne `saveStateId` wird nichts gemerkt.

### Automatisches Laden beim Seitenaufruf

Es gibt zwei Wege, die Daten schon beim Laden (ohne Klick) zu zeigen:

- **Query mit `url`** (wie oben): lädt **automatisch** beim ersten gültigen Filter-Zustand. Kein Button nötig — das ist der Standard-Fall für Filter+Tabelle.
- **ButtonLine-Button mit `WithAutoLoad(true)`**: Wenn der Filter statt der Query-`url` einen expliziten Button zum Nachladen verwendet (z. B. wegen Snackbar/Poll/`WithData`-Payload), triggert `WithAutoLoad(true)` diesen Button **einmalig automatisch beim Laden**, sobald der Filter gültig ist:

  ```go
  bl := button.NewButtonLine("right", nil)
  bl.Add(button.NewSimpleApiButton("Suchen", c.apiUrl("data"), core.ColorPrimary).WithAutoLoad(true))

  q := query.NewQueryWithFormGroup(fg, nil, nil /* keine url */, bl, nil, nil)
  ```

  Feuert genau einmal; spätere Filter-Änderungen erfordern wieder den Klick. Details: `references/components.md` → „Auto-Load via `WithAutoLoad`".

## Flags — UI-Only State, der nicht als Filter zählt

Manchmal schickt das Frontend zusätzliche UI-State-Keys mit (z.B. View-Modus), die **nicht** als Filter in die DB-Query gehören:

```go
tbl.SetFlags("_viewMode", "_showArchived")
```

Diese Keys landen **nicht** in der `parsedFilters`-Map (werden vor `ParseAndValidate` entfernt). Sie bleiben aber in `tbl.GetFilterData()` — für den Fall, dass du sie trotzdem inspizieren willst.

## CSV-/Excel-Export

Wenn die Tabelle `Csv: &true` oder `Excel: &true` in den Options hat, setzt das Frontend `_csv: true` bzw. `_excel: true` im Request — `LoadFilterData` erkennt das und setzt `tbl.GetOutputType()`.

```go
filters, _ := tbl.LoadFilterData(ctx)
rows := c.svc.FindAll(filters, pagination.NewInfinite())  // kein Limit für Export

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

## Bulk-Actions mit Select-Buttons

Zeilen-Mehrfach-Selection:

```go
selectBtn := button.NewTableButton(
    core.ButtonActionApi,
    "delete",
    c.apiUrl("bulk-delete"),
    "Löschen",
    core.ColorError,
    false,
    nil,
)
tbl.SetSelectButtons([]*button.TableButton{selectBtn})
```

Das Frontend aktiviert automatisch Row-Checkboxes und sendet bei Button-Click eine Liste von IDs:

```go
// Controller
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
    return wc.RefreshTable()
}
```

## Häufige Fehler

- **Reihenfolge vertauscht:** `LoadPaginationParams()` **vor** `LoadFilterData()` → leere Werte (weil `tc.filterData == nil`). Immer erst Filter laden.
- **Falscher Type-Assertion:** `filters["priority"].([]int32)` — falsch, korrekt ist `field.ModelListValue`. Entweder `field.ModelListValue`-Cast oder (wenn man den Cast nicht importieren will) `[]int32` funktioniert, weil `ModelListValue` genau `[]int32` ist — aber Go's Type-Assertion schaut auf den Named-Type, also `ModelListValue` nehmen.
- **Leere Filter werden als `nil` weitergegeben:** `q.Where("status = ?", nil)` erzeugt `status IS NULL`. Deshalb `, ok && v != ""` bzw. `len(v) > 0` Checks.
- **Sort-Validierung umgehen:** `LoadPaginationParams` clamped `Sort` auf definierte Field-IDs. Wenn du aber `q.Order(pg.Sort)` mit leerem Sort aufrufst, produziert GORM ein `ORDER BY`-Fragment — immer `if pg.Sort != ""` prüfen.
