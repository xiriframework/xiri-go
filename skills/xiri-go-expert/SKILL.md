---
name: xiri-go-expert
description: Experte für die xiri-go Go-Library. Verwende diesen Skill IMMER wenn Go-Code geschrieben wird der xiri-go importiert (github.com/xiriframework/xiri-go), oder wenn der User nach Komponenten, Formularen, Tabellen, Filter-Parsing, URL-Prefixes/Sidebar-Routing, Dialogen, Responses, UiContext, Inline-Edit, abhängigen Formularfeldern (SetReloadOn/ExportPatch), Bulk-Actions/MassEdit/MassDelete oder dem Builder-Pattern der xiri-go Library fragt.
---

# xiri-go Expert

Go-Framework für JSON-driven UIs (`github.com/xiriframework/xiri-go`). Das Angular-Frontend **xiri-ng** rendert die JSON-Ausgabe.

**Arbeitsweise:** Diese Datei ist nur Navigation + minimale Anker. Für jede konkrete Aufgabe die passende Reference-Datei lesen (siehe Tabelle unten) — **nicht** alles hier suchen.

## Architektur

```
FormGroup/Builder → *.Print(ctx *core.UiContext) → map[string]any → JSON → xiri-ng rendert
```

- **HTTP:** Echo v4 · **ORM:** GORM · **UiContext** per Request
- **Typische Projekt-Convention:** Die meisten Projekte wrappen `echo.Context` + `UiContext` in einen eigenen Web-Context-Helper (z. B. `wc := webcontext.GetWebContext(ctx)`), der die xiri-go-`response.*`-Funktionen und Error-Codes dünn kapselt. Erwartetes Interface:
  - Daten/Pages: `Page(p)`, `Data(comp)`, `Dialog(dlg)` — `Data(...)` handled CSV/Excel-Export automatisch (Content-Type + Disposition aus `response.DataResult.Type`).
  - Navigation/Refresh: `Goto(url)`, `RefreshPage()`, `RefreshTable()`, `Done()`
  - Dialoge: `DeleteDialog(name)`
  - Context: `UiContext()`
  - Errors: `BadRequest`, `Unauthorized`, `Forbidden`, `NotFound`, `InternalServerError`, `ServiceUnavailable` — alle `(msg)`.
- **Ohne Wrapper** direkt über xiri-go: `response.NewReturnDone()`, `response.NewReturnGoto(url)`, `response.NewReturnRefreshPage()`, `response.NewReturnRefreshTable()`, `response.NewDataResponse(data)`, `response.NewErrorResponse(msg)`. Siehe `references/responses.md`.

## URL-Handling — die wichtigste Regel

```go
import xurl "github.com/xiriframework/xiri-go/component/url"

xurl.NewUrl("/devices")                      // Frontend-Link (kein Prefix)
xurl.NewUrlPrefix("/devices", "/api/v1")     // API-Endpoint (mit Prefix)
u.Add("edit", "42")                           // Chain-Append
u.Print()        // ohne Prefix — Links, Breadcrumbs
u.PrintPrefix()  // mit Prefix   — API-Calls, form/table Submit-URLs
```

**Regel:** Page-/Link-URLs = `Print()`. API-URLs (`tbl.SetURL`, `form.NewForm`, Dialog-Submit) = `*xurl.Url` direkt (intern korrekt aufgelöst) oder `PrintPrefix()` wo `string` erwartet wird. Controller-Helper + Sidebar-Pattern: **`references/url-routing.md`**.

**Query-String in Navigations-Links ist erlaubt:** `xurl.NewUrl("/Admin/Devices/Table?config=98")` funktioniert für Frontend-Links (Button `action: 'link'`, cardlink, links/list/sidenav, Tabellen-Link-Felder) — das Frontend parst den `?…`-Teil zu Route-Query-Params (via `xiriUrl`-Pipe in xiri-ng). `Add(...)` hängt nur Path-Segmente an (`/…/98`); für Query-Params den String direkt an `NewUrl` übergeben.

## Was wann lesen

| Du baust / brauchst …                                    | Lies                                        |
| -------------------------------------------------------- | ------------------------------------------- |
| **Eine ganze Seite** (Page, Grid-Layout, Header, Farben) | `references/pages.md`                       |
| **Eine Tabelle** (Spalten, Row/Top/Bulk-Buttons)         | `references/tables.md`                      |
| **Einen Dialog** (Delete, Form, Table, Waiting, Component) | `references/dialogs.md`                     |
| **Ein Formular** (Felder, showWhen, Validierung)         | `references/form-fields.md` + `form-builder.md` |
| **Abhängige Felder** (Optionen vom Server nachladen)     | `references/form-fields.md` (Abschnitt „Abhängige Felder“) |
| **Filter + Pagination** mit GORM                         | `references/table-filtering.md`             |
| **End-to-End Pattern** (CRUD, MassEdit, MassDelete, Stepper, Inline-Edit, Dashboard) | `references/patterns.md`                    |
| **Response-Typen** (Refresh, Goto, DataResult, DataResponse) | `references/responses.md`                   |
| **Locale-Formatter** (Datum, Zahlen, Distanz, Pressure)  | `references/formatter.md`                   |
| **UiContext / Translator / Locale-Setup**                | `references/uicontext.md`                   |
| Eine nicht-tabellen Komponente (Tabs, Timeline, BarChart …) | `references/components.md`                  |
| TachoTime (Fahrtenschreiber)                             | `references/tachotime.md`                   |
| Enum-Werte                                               | `references/enums.md`                       |
| Field-Methoden-Details (Table-Builder-Chain)             | `references/table-builder.md`               |

## Absolute Minimal-Anker (Signatur-Gerüst)

Wenn du nichts lesen willst, ist das die Notfall-Zusammenfassung. Für alle echten Aufgaben: Reference-Datei lesen.

### Imports-Block

```go
import (
    "github.com/xiriframework/xiri-go/component/core"
    "github.com/xiriframework/xiri-go/component/page"
    "github.com/xiriframework/xiri-go/component/pageheader"
    "github.com/xiriframework/xiri-go/component/button"
    "github.com/xiriframework/xiri-go/component/table"
    "github.com/xiriframework/xiri-go/component/barchart"
    "github.com/xiriframework/xiri-go/component/form"
    "github.com/xiriframework/xiri-go/component/dialog"
    "github.com/xiriframework/xiri-go/form/field"
    formbuilder "github.com/xiriframework/xiri-go/form/builder"
    "github.com/xiriframework/xiri-go/response"
    xurl "github.com/xiriframework/xiri-go/component/url"
)
```

### Table-Grundgerüst

```go
b := table.NewBuilder[Row]()
b.IdField("id", "ID", func(r Row) int64 { return r.ID })
b.TextField("name", "Name", func(r Row) string { return r.Name })
tbl := b.Build()
tbl.SetURL(c.apiUrl("data"))

// Data-Handler
filters, err := tbl.LoadFilterData(ctx)   // Fehler NICHT verwerfen: ein fehlerhafter
if err != nil {                           // Body würde sonst als "kein Filter" gelten
    return err                            // und einen Full-Table-Export auslösen
}
pg := tbl.LoadPaginationParams()  // entfernen bei client-seitiger Pagination (default)
tbl.SetData(rows)
return wc.Data(tbl)
```

### Form-Field Konstruktoren (id, name, required, default)

Der Default ist ein **Wert**, kein Pointer — `nil` kompiliert nicht. Ohne Default den Nullwert
übergeben (`""`, `0`, `false`). `Value` ist nach `BindAndValidate` dagegen ein Pointer:

```go
field.NewTextField    (id, name, required, "")         // Value: *string
field.NewIntField     (id, name, required, 0)          // Value: *int32
field.NewBoolField    (id, name, required, false)      // Value: *bool
field.NewTimeField    (id, name, required, 0)          // Value: *int64 (Unix)
field.NewSelectField  (id, name, required, opts)       // Value: int32
field.NewSelectField  (id, name, required, opts).SetMultiple(true)  // Values: []int32
field.NewModelField   (id, name, required, "group", 0) // Value: int32
// Subtypes auf TextField: "email" | "tel" | "url" | "password" | "textarea"
```

### Abhängige Felder (Optionen vom Server nachladen)

`SetShowWhen` blendet ein Feld ein/aus, `SetReloadOn` lädt seinen **Inhalt** neu:

```go
reloadURL := xurl.NewUrlPrefix("/Portal/Thing/FormReload", "/api")
tags.BaseField.SetReloadOn(reloadURL, "status")   // beides Pflicht, sonst wird nichts exportiert

// Handler hinter reloadURL:
status, tags, fg := ctrl.buildThingForm(wc.UiContext())
if err := builder.BindReload(c, fg); err != nil {   // nachsichtig: leeres Pflichtfeld ist ok
    return wc.BadRequest(err.Error())
}
tags.Options = ctrl.tagsForStatus(status.Value)
return c.JSON(http.StatusOK, response.NewReturnFields(fg.ExportPatch()))
```

Das Frontend postet **nur die Trigger-Werte**, behält Werte die es in der neuen Liste noch gibt und
verwirft den Rest. Filter erben das. Details + Grenzen in `references/form-fields.md`.

### FormBuilder

```go
fb := formbuilder.NewFormBuilder(wc.UiContext())
fb.AddField(nameField)

fields, _ := fb.BuildAddForDisplay()   // für form.NewForm
fg, _, _  := fb.BuildAdd()              // für BindAndValidate
formbuilder.BindAndValidate(ctx, fg)
// Werte direkt: if nameField.Value != nil { entity.Name = *nameField.Value }
```

### Filter-Typen (Single vs. Multi)

| Field                           | Go-Typ in `filters` map       |
| ------------------------------- | ----------------------------- |
| TextField                       | `string`                      |
| IntField                        | `int32`                       |
| BoolField                       | `bool`                        |
| TimeField                       | `int64`                       |
| SelectField                     | Option-Value-Typ              |
| SelectField + `.SetMultiple(true)` | `field.ModelListValue` (`[]int32`) |
| ModelListField                  | `field.ModelListValue`        |
| ChipsField                      | `[]interface{}` (int64 IDs + string Freitext) |

GORM-Integration (mit `, ok` Guards, leere Werte raus) → `references/table-filtering.md`.

### Responses

```go
response.NewReturnRefreshPage()          // Seite neu laden
response.NewReturnRefreshTable()         // Table reload
response.NewReturnGoto("/path")          // Navigation
response.NewReturnDone()                 // Fertig
response.NewReturnSuccess("Gespeichert") // Snackbar
resp.WithMessage("Text", response.MessageSuccess)
response.NewReturnInlineEdit().WithUpdates(upd).WithRefreshTable()
response.NewReturnPoll(statusUrl, 2000)  // selbst-pollender Button (.WithText / .WithButton)
response.NewButtonPatch().WithText("Erledigt").WithColor("success").Disable() // Button am Ende ändern (.WithButton)
tbl.SetPoll(2000)                        // Tabelle pollt selbsttätig (auto-refresh), solange gesetzt
```

### Farben / Button-Types

```go
core.ColorPrimary | Secondary | Tertiary | Accent | Warning | Error | Success
  // + Extended: Emerald/Red/Yellow/Green/Blue/Purple/Orange/Gray/LightGray/DarkGray/White/Black/Inherit
core.ButtonTypeRaised | Basic | Stroked | Flat | Fab | MiniFab | Icon | IconText
```

## Wichtige Do-nots

- **Kein** `NewTextareaField` — `NewTextField` mit `.Subtype = "textarea"`.
- **Kein** `NewMultiSelectField` — `NewSelectField(...).SetMultiple(true)`.
- **Kein** `NewDeviceListField` — `NewModelListField(id, name, required, "device", ids)`.
- **Kein `nil` als Default-Argument** bei Text/Int/Bool/Time — das sind Wert-Parameter, nicht
  Pointer. Nullwert übergeben (`""`, `0`, `false`). `nil` ist nur bei Slice/Map-Defaults erlaubt
  (`NewModelListField`, `NewArrayField`, `NewJsonField`).
- **Keine** URL-Strings konkatenieren — `*xurl.Url` via Controller-Helper.
- **`LoadPaginationParams` NACH `LoadFilterData`** (liest aus gespeicherten Filtern). Nur bei **Server-Side-Pagination** nötig — client-seitig (xiri-ng default) weglassen.
- **`ButtonsField`-Keys sind Strings** (`"0"`, `"1"`), nicht Ints.
- **Chips in Tabellen-Zellen**: `b.ChipsField(id, name, accessor)` (pure Display, accessor → `[]table.Chip`). Für editierbare Multi-Select-Chips siehe `WithEditableChipOptions(...)` (anderer Mechanismus).
- **Filter-Guards**: `if v, ok := filters[k].(string); ok && v != ""` — leere Werte erzeugen sonst falsche WHERE-Clauses.
- **`SetReloadOn` braucht URL *und* Feld-IDs.** Fehlt eines, wird gar nichts exportiert und das Feld
  verhält sich wie ein normales — kein Fehler, nur wirkungslos.
- **Im Reload-Handler `BindReload`, nicht `BindAndValidate`.** Letzteres scheitert an leeren
  Pflichtfeldern, die mitten im Ausfüllen völlig normal sind.
- **Abhängige `ModelListField`/Treeselects ohne `URL` bauen.** Mit URL lädt das Frontend den Baum
  selbst und ignoriert die gepatchte Liste.
