# Seiten-Aufbau — Page, Layout, Header, Sections, Farben

Dieses File beschreibt, wie eine komplette xiri-go-Seite strukturiert wird: `page.Page` als Top-Level-Container, das 12-Spalten-Grid der xiri-ng-Frontend-Renderung, wie man `display`-Klassen pro Komponente setzt, Header/Sections als Strukturelemente, ButtonLines, und die Farbpalette.

## Das Gesamtbild

```
Page (top-level)
├── Bread("Home", …)           ─┐
├── Bread("Devices", …)         ├── Breadcrumbs am Seitenkopf
├── Bread("Detail", nil, …)    ─┘
├── Add(pageheader.New(...))    ── PageHeader (Titel + Buttons)
├── Add(statgrid…)              ┐
├── Add(section.New(...))       ├── Komponenten im 12-Spalten-Grid
├── AddNewRow(table…)           ┘   (jede mit optionalem display)
└── Extra("customKey", value)   ── Extra-Felder am JSON-Root
```

Das Frontend (`xiri-ng` + `xiri-dyncomponent`) rendert jede Komponente innerhalb eines `<div class="xcol …" [class]="obj.display">` — die `display`-Klasse, die du in Go via `WithDisplay(...)` setzt, landet genau dort.

## `page.Page` — Top-Level-Container

Quelle: `component/page/page.go`.

```go
import (
    "github.com/xiriframework/xiri-go/component/page"
    xurl "github.com/xiriframework/xiri-go/component/url"
)

p := page.NewPage()

// Breadcrumbs
p.Bread("Start",   xurl.NewUrl("/"),           false)
p.Bread("Devices", c.pageUrl(),                 false)
p.Bread("Detail",  nil,                         false)   // letztes: kein Link
// 3. Arg: extern — true öffnet in neuem Tab

// Komponenten
p.Add(component)           // nächste freie Grid-Position
p.AddNewRow(component)     // erzwingt neue Grid-Zeile (xcol-start)
p.AddOld(mapData)          // bereits serialisiertes map[string]any

// Extras am JSON-Root (selten)
p.Extra("tenantId", 42)
```

`p.Print(ctx)` liefert:

```json
{
  "type": "page",
  "bread": [{ "name": "Start", "link": "/", "extern": false }, ...],
  "data": [
    { "type": "page-header", ... },
    { "type": "stat-grid",   ..., "newRow": false },
    { "type": "table",        ..., "newRow": true }
  ],
  "tenantId": 42
}
```

Rückgabe meist via `wc.Page(p)` (GEM-WebContext), nicht via `c.JSON(...)`.

### `Add` vs `AddNewRow`

- `Add(cmp)` → hängt die Komponente in die aktuelle Grid-Zeile. Wenn mehrere nebeneinander passen, wird gefüllt.
- `AddNewRow(cmp)` → setzt `newRow: true` auf dem Entry — das Frontend hängt die CSS-Klasse `xcol-start` an den Wrapper, d.h. die Komponente startet in einer neuen Grid-Zeile (Spalte 1).

Faustregel: Die erste Komponente nach dem `page-header` **nicht** als `AddNewRow`, aber z.B. vor einer zweiten, großen Tabelle schon.

## 12-Spalten-Grid im Frontend

Jede Komponente wird vom Frontend in einen Grid-Slot gesteckt. Die Grid-CSS-Regeln (`projects/xiri-ng/styles/grid.scss`):

```scss
.xrow {
  display: grid;
  grid-template-columns: repeat(12, 1fr);
  grid-column-gap: var(--xiri-grid-column-gap, 24px);
  grid-row-gap:    var(--xiri-grid-row-gap, 20px);
}

.xcol               { grid-column-end: span 12; }   // Default: volle Breite
.xcol-{N}           { grid-column-end: span N; }    // N = 1..12
.xcol-{bp}-{N}      { /* wirksam ab Breakpoint bp */ }
.xcol-start         { grid-column-start: 1; }       // erzwungener Zeilenumbruch
.xcol-middle-{bp}-{N} { /* zentriert über N Spalten */ }
.xcol-right-{bp}-{N}  { /* rechtsbündig */ }
.x-{bp}-none        { display: none ab Breakpoint bp; }
.x-valign-center    { align-content: center; }
.x-align-center     { justify-self: center; }
.x-align-right      { justify-self: end; }
```

### Breakpoints

| Abkürzung | `min-width` |
| --------- | ----------- |
| `xs`      | `0` (alle)  |
| `sm`      | `576px`     |
| `md`      | `768px`     |
| `lg`      | `992px`     |
| `xl`      | `1200px`    |
| `xxl`     | `1400px`    |

### Typische Layout-Klassen

- `xcol-md-6` → auf Tablet/Desktop halbe Breite (6 von 12), auf Mobile volle Breite
- `xcol-md-4 xcol-lg-3` → Tablet: ein Drittel, Desktop: ein Viertel
- `xcol-md-12` → immer volle Breite (kein Neben-einander)
- `xcol-middle-md-8` → 8 Spalten breit, auf Tablet+ zentriert
- `xcol-right-md-4` → 4 Spalten breit, rechtsbündig

## `display`-Klasse pro Komponente setzen

Fast jede Komponente hat eine `WithDisplay(class string)`-Methode oder akzeptiert `display *string` im Konstruktor:

```go
statgrid.New().WithDisplay("xcol-md-12")
card.NewCard(...).WithDisplay("xcol-md-6")
section.New().WithDisplay("xcol-md-8 xcol-middle-md-8")
button.NewButtonLine("right", nil).WithDisplay("xcol-md-6")
form.NewForm(fields, url, nil, nil, nil, ctx).WithDisplay("xcol-md-8 xcol-middle-md-8")
```

Komponenten **ohne** `display` landen in einem `.xcol`-Wrapper und nehmen volle Breite (Default). Wird `display` gesetzt, **überschreibt** die Klasse den Default — deshalb muss man bei `display` an alle gewünschten Breakpoints denken (z.B. `xcol-md-6` allein nimmt auf Mobile volle Breite, das ist meist gewollt).

### Drei-Spalten-Layout

```go
p := page.NewPage()
p.Add(stat.New("42", "Online").WithDisplay("xcol-md-4"))
p.Add(stat.New("3",  "Offline").WithDisplay("xcol-md-4"))
p.Add(stat.New("98%", "Uptime").WithDisplay("xcol-md-4"))
```

Auf Desktop drei nebeneinander, auf Mobile untereinander. Für noch feinere Kontrolle mit `statgrid`:

```go
sg := statgrid.New().Columns(3)       // setzt intern die Logik
```

Der `statgrid.Columns` Helper ist spezialisiert — für generische Komponenten nutzt man `xcol-md-N`-Klassen.

### Zentriertes Formular (klassisches Muster)

```go
f := form.NewForm(fields, c.apiUrl("save"), nil, nil, nil, wc.UiContext()).
    WithDisplay("xcol-md-8 xcol-middle-md-8")
p.Add(f)
```

→ auf Tablet+ 2/3 breit, zentriert; auf Mobile volle Breite.

## `pageheader.PageHeader` — Seitentitel mit Buttons

```go
import "github.com/xiriframework/xiri-go/component/pageheader"

ph := pageheader.New("devices.title").
    Subtitle("devices.subtitle").
    Icon("router", core.ColorPrimary).
    Buttons(buttonLine).
    WithDisplay("xcol-md-12")

p.Add(ph)
```

API:

```go
pageheader.New(title string) *PageHeader
  .Subtitle(s string) *PageHeader
  .Icon(icon string, color core.Color) *PageHeader
  .Buttons(b *button.ButtonLine) *PageHeader
  .WithDisplay(display string) *PageHeader
```

Konvention: **Jede Seite bekommt einen PageHeader**, auch wenn leer — er ist der visuelle Anker und enthält Seiten-Actions (Add-Button, Sync-Button, Settings, …).

## `section.Section` — Gruppierung

Für mehrere Komponenten, die logisch zusammengehören (z.B. "Allgemeine Daten", "Administrative Daten"):

```go
import "github.com/xiriframework/xiri-go/component/section"

sec := section.New().
    Title("section.devices").
    Subtitle("section.devices.sub").
    Icon("router", core.ColorPrimary).
    Collapsible(true).                       // klappbar (kollabiert default)
    Buttons(buttonLine).                     // Buttons in der Section-Kopfzeile
    Add(tableCmp).
    Add(buttonsCmp).
    WithDisplay("xcol-md-12")

p.Add(sec)
```

API:

```go
section.New() *Section
  .Title(t string) *Section
  .Subtitle(t string) *Section
  .Icon(icon string, color core.Color) *Section
  .Collapsible(collapsed bool) *Section    // true = initial kollabiert
  .Buttons(b *button.ButtonLine) *Section
  .Add(component core.Component) *Section
  .WithDisplay(display string) *Section
```

**Unterschied zu Card:** Eine `section` ist offener (keine Box), dient zur Strukturierung. Eine `card` hat eine Box mit Rand, Schatten, Header. Für KPI-Panels, Formulare, Details → `section`. Für abgesetzte Inhalte mit visueller Grenze → `card`.

## Layout-Helpers (`component/layout`)

Kleine Bausteine, die meist **innerhalb** anderer Komponenten (Section, Card) oder als Zeilentrenner in `Page` eingesetzt werden.

### `layout.Header` — Überschrift mitten in der Seite

```go
import "github.com/xiriframework/xiri-go/component/layout"

size := "h2"
display := "xcol-md-12"
layout.NewHeader("Allgemeine Daten", core.ColorPrimary, &size, &display)
// oder als Chain:
h := layout.NewHeader("Infos", core.ColorInherit, nil, nil).
    WithSize("h3").
    WithDisplay("xcol-md-12")
```

Renderable Sizes: `h1 … h6`. Default (wenn `nil` übergeben) ist i.d.R. `h2`.

### `layout.Divider` — Trennlinie

```go
div := layout.NewDivider().
    Text("Weitere Optionen").
    Icon("chevron_right").
    Spacing("large").                        // "compact" | "normal" | "large"
    WithDisplay("xcol-md-12")
p.Add(div)
```

### `layout.Spacer` — vertikaler Leerraum

```go
display := "xcol-md-12"
p.Add(layout.NewSpacer(&display))            // oder NewSpacer(nil) für Default
```

Nützlich zwischen zwei Komponenten, wenn der Grid-Row-Gap nicht genug ist.

### `layout.Container` — mehrere Komponenten in einem Grid-Slot

```go
container := layout.NewContainer(nil).
    Add(stat1).
    Add(stat2).
    Add(stat3)

container.(*layout.Container)    // Struct-Name
```

Alles im Container landet in **einem** Grid-Slot, intern aber wieder als 12-Spalten-Subgrid. Praktisch für Kompositionen, die als Einheit skalieren sollen.

### `layout.Html` — Raw-HTML

```go
layout.NewHtml("<p>Beliebiger <b>HTML</b>-Inhalt</p>", nil)
```

**Vorsicht vor XSS** — niemals Nutzereingaben direkt in `NewHtml` pipen.

## `button.ButtonLine` — Button-Leiste

```go
bl := button.NewButtonLine("right", nil)          // class + display
bl.Add(button.NewSimpleCloseButton("Abbrechen"))
bl.Add(button.NewSimpleApiButton("Speichern", c.apiUrl("save").PrintPrefix(), core.ColorPrimary))
bl.WithDisplay("xcol-md-12")

p.Add(bl)
```

API:

```go
button.NewButtonLine(class string, display *string) *ButtonLine
  .Add(button *Button) *ButtonLine
  .WithDisplay(display string) *ButtonLine
  .Print(ctx) map[string]any
  .PrintData(ctx) map[string]any            // nur die data, ohne Type-Wrapper
  .PrintButtons(ctx) []map[string]any        // nur das Button-Array
  .DataResponse(ctx) response.DataResult     // für Ajax-Replace via RefreshButtons
```

### `class`-Werte

- `""` (leer) → defaults zu `"right"` (rechtsbündig)
- `"right"` → rechtsbündig (Standard für Form-Buttons)
- `"left"` / `"center"` → entsprechend ausgerichtet (Frontend-CSS)
- `"small"` → kompakte Button-Leiste (kleinere Icons/Padding), typisch in Table-Toolbars

Die konkret verfügbaren Klassen hängen vom Frontend-Styling ab — im Zweifel im `xiri-ng`-Projekt nach `.buttonline.<class>` greppen.

### `ButtonLine` in PageHeader / Section

Sehr häufig:

```go
buttons := button.NewButtonLine("", nil)    // "right"
buttons.Add(button.NewLinkButton("Neu", c.pageUrl("add"),
    core.ColorPrimary, core.ButtonTypeRaised, "", false, nil, nil))
p.Add(pageheader.New("Geräte").Buttons(buttons))
```

## Farben

Quelle: `component/core/enums.go`. `core.Color` ist ein `string`-Alias.

### Theme-Farben (primär zu nutzen)

| Konstante             | String         | Einsatz                                          |
| --------------------- | -------------- | ------------------------------------------------ |
| `core.ColorPrimary`   | `"primary"`    | Haupt-Actions (Speichern, Submit, Neu)           |
| `core.ColorSecondary` | `"secondary"`  | Sekundär-Actions                                 |
| `core.ColorTertiary`  | `"tertiary"`   | Tertiär-Actions                                  |
| `core.ColorAccent`    | `"accent"`     | Highlights, Links                                |
| `core.ColorWarning`   | `"warn"`       | Warnungen, Delete-Dialoge                        |
| `core.ColorError`     | `"error"`      | Fehler, destruktive Aktionen                     |
| `core.ColorSuccess`   | `"success"`    | Success-Status, Checkmarks                       |

### Extended-Farben (für Visualisierungen, Icons)

| Konstante             | String       |
| --------------------- | ------------ |
| `core.ColorEmerald`   | `"emerald"`  |
| `core.ColorRed`       | `"red"`      |
| `core.ColorYellow`    | `"yellow"`   |
| `core.ColorGreen`     | `"green"`    |
| `core.ColorBlue`      | `"blue"`     |
| `core.ColorPurple`    | `"purple"`   |
| `core.ColorOrange`    | `"orange"`   |
| `core.ColorGray`      | `"gray"`     |
| `core.ColorLightGray` | `"lightgray"`|
| `core.ColorDarkGray`  | `"darkgray"` |
| `core.ColorWhite`     | `"white"`    |
| `core.ColorBlack`     | `"black"`    |
| `core.ColorInherit`   | `"inherit"`  |

**Regel:** Theme-Farben für semantische UI-Bedeutung (Primary = "das macht man", Warn = "Vorsicht"), Extended-Farben nur für rein-visuelle Kategorisierung (Chart-Serien, Status-Icons, Tags).

## Button-Types

```go
core.ButtonTypeRaised    // "raised"     — gefüllt, hebt sich ab — Primary-Actions
core.ButtonTypeBasic     // "basic"      — flach, kein Hintergrund
core.ButtonTypeStroked   // "stroked"    — Outline, ideal für Cancel
core.ButtonTypeFlat      // "flat"       — gefüllt ohne Elevation
core.ButtonTypeFab       // "fab"        — Floating Action Button (rund, groß)
core.ButtonTypeMiniFab   // "minifab"    — FAB, kleiner
core.ButtonTypeIcon      // "icon"       — nur Icon, keine Text
core.ButtonTypeIconText  // "icontext"   — Icon + Text
```

## Vollständiges Layout-Beispiel

```go
func (c *Controller) DevicePage(ctx echo.Context) error {
    wc := webcontext.GetWebContext(ctx)
    uc := wc.UiContext()

    p := page.NewPage()
    p.Bread("Start",   xurl.NewUrl("/"),   false)
    p.Bread("Devices", c.pageUrl(),         false)

    // PageHeader mit Add-Button (rechtsbündig)
    buttons := button.NewButtonLine("", nil)
    buttons.Add(button.NewLinkButton("Neu", c.pageUrl("add"),
        core.ColorPrimary, core.ButtonTypeRaised, "add", false, nil, nil))
    p.Add(pageheader.New("Geräte").
        Icon("router", core.ColorPrimary).
        Buttons(buttons))

    // KPI-Reihe: drei Stats nebeneinander auf Desktop
    p.Add(stat.New(strconv.Itoa(c.svc.CountTotal()), "Gesamt").
        Icon("router").WithDisplay("xcol-md-4"))
    p.Add(stat.New(strconv.Itoa(c.svc.CountOnline()), "Online").
        Icon("check_circle").IconColor(core.ColorSuccess).WithDisplay("xcol-md-4"))
    p.Add(stat.New(strconv.Itoa(c.svc.CountOffline()), "Offline").
        Icon("cancel").IconColor(core.ColorError).WithDisplay("xcol-md-4"))

    // Section mit Tabelle — volle Breite, neue Zeile
    sec := section.New().
        Title("Alle Geräte").
        Icon("list", core.ColorPrimary).
        Collapsible(false).
        Add(c.buildTable(uc)).
        WithDisplay("xcol-md-12")
    p.AddNewRow(sec)

    return wc.Page(p)
}
```

Ergebnis: Breadcrumbs → Header mit Neu-Button → drei Stat-Karten in einer Zeile → Section (volle Breite, neue Zeile) mit Tabelle.

## Häufige Fehler

- **`display` überschreibt Default**: Ein nacktes `WithDisplay("xcol-md-6")` lässt das Default-`xcol` weg. In der Regel nicht schlimm (das Frontend fügt `xcol` immer als Basisklasse hinzu), aber beim Debugging beachten.
- **Breadcrumb-URL-Format**: `p.Bread(name, *xurl.Url, extern)` — der zweite Parameter ist `*xurl.Url`, nicht `string`. Für "kein Link" → `nil`.
- **`AddNewRow` vor der ersten Komponente** hat keinen visuellen Effekt (die erste ist immer in einer neuen Zeile).
- **ButtonLine `class=""`** defaults zu `"right"` — das ist **nicht** "keine Klasse". Wenn du linksbündig willst, explizit `NewButtonLine("left", nil)`.
- **PageHeader ohne Page**: Ein `PageHeader` alleine ist eine Komponente, die man auch ohne `Page` rendern kann — aber das ist selten gewollt (z.B. bei Card-internen Headern, dafür gibt's die `Card`-Header-Parameter).
- **`Collapsible(true)` setzt initial-kollabiert**: `true` = "initial zugeklappt", nicht "aktiviert". Der `collapsible`-State ist immer aktiv, das Argument steuert nur den Start-Zustand.
