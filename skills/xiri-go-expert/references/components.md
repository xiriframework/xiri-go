# Komponenten (component package)

Import-Basis: `"github.com/xiriframework/xiri-go/component/"`

Alle Komponenten implementieren `Component.Print(ctx *core.UiContext) map[string]any`.

## Page (`component/page`)

Container für eine komplette Seite mit Breadcrumbs.

```go
p := page.NewPage()
p.Bread("Home", "/", false)
p.Bread("Devices", "/devices", false)
p.Bread("Detail", "", false)  // letzter ohne Link
p.Add(component)               // Komponente hinzufügen
p.AddNewRow(component)          // Neue Grid-Zeile erzwingen
p.Extra("customKey", value)     // Extra-Feld am Root-Level
p.Print(ctx)
```

## PageHeader (`component/pageheader`)

Seitentitel mit optionalem Icon und Buttons.

```go
ph := pageheader.New("devices.title")
ph.Subtitle("devices.subtitle")
ph.Icon("devices", core.ColorPrimary)
ph.Buttons(buttonLine)
ph.Print(ctx)
```

## Card (`component/card`)

Karte mit Header, Content und Buttons.

```go
// Table-Card
c := card.NewCard(core.CardTypeTable, tableComponent, "devices.list", "", "", "", true, false, "")
c.ButtonTop(addButton)
c.WithCollapsible(true)
c.WithMaxHeight("400px")
c.Print(ctx)

// List-Card
c := card.NewCardList("info.title", cardListContent)
c.Print(ctx)

// AJAX-Card (lädt Daten von URL)
c := card.NewCard(core.CardTypeTable, nil, "devices.list", "", "", "", true, false, "")
c.SetURL("/api/devices/table")
c.WithReload(true)
c.Print(ctx)
```

**Spalten-Header bei Tabellen-Cards** (`CardTypeTable`) standardmäßig aus. Opt-in via `.WithTableHeader()` — gleicher Name wie auf `Dialog`, no-op für Nicht-Tabellen-Cards:

```go
c := card.NewCardList("Übersicht", content).WithTableHeader()
```

### CardListContent

```go
content := card.NewCardListContent([]card.CardListContentLine{
    {Name: "Name", Content: "Device 1"},
    {Name: "Status", Content: "Online"},
})
content.SetDense(true)
```

### CardListField (custom fields)

```go
content := card.NewCardListContentFields(
    []card.CardListField{{Name: "name", Label: "Name"}, {Name: "status", Label: "Status", Type: "badge"}},
    []map[string]interface{}{{"name": "Device 1", "status": "Online"}},
)
```

### Multi-Component Card (Card.Add)

Eine Card kann mehrere beliebige Sub-Components in einem xcol-Grid hosten — Card-Header bleibt darüber. Sobald `Add(...)` mindestens einmal aufgerufen wurde, ignoriert das Frontend `content`/`fields` und rendert stattdessen die Sub-Components via `xiri-dyncomponent`. Layout pro Component via `WithDisplay("xcol-…")`.

```go
c := card.NewCard(core.CardTypeTable, nil, "Activity", nil, ptr("show_chart"), ptr("primary"), false, false, nil)
c.WithPadding("md").                                           // Innen-Padding (Token oder CSS-Wert)
  Add(barchart.New("activity").Mode(barchart.ModeSimple).
        BarNamed("M", "Monday", 3).BarNamed("T", "Tuesday", 9). /* ... */
        Compact().                                              // ohne eigene mat-card-Hülle
        WithDisplay("xcol-12")).
  Add(stat.New("18h", "Today").Compact().WithDisplay("xcol-6")).
  Add(stat.New("32h", "Last 7 days").Compact().WithDisplay("xcol-6"))
```

`Add(component)` akzeptiert jeden `core.Component` (Stat, BarChart, DescriptionList, Timeline, …). Der `cardType`-Parameter ist im Multi-Component-Modus irrelevant — `core.CardTypeTable` als Platzhalter ist OK.

#### `Card.WithPadding(value string)`

Innen-Padding der Sub-Components-Fläche im Multi-Component-Modus.

- **Tokens:** `"xs" | "sm" | "md" | "lg" | "xl"` — gemappt auf `--xiri-spacing-*` (xs=4, sm=8, md=16, lg=24, xl=32 px).
- **Freier CSS-Wert:** `"16px"`, `"1rem"`, `"var(--xiri-spacing-md)"`.
- **Default:** `"md"` (16 px).
- **xs-Viewport (<576 px) IMMER 8 px** — der konfigurierte Wert greift erst ab `sm`. Wirkt nicht im Tabellen-Modus (CardListContent / fields), dort bleibt das Padding wie zuvor.

#### `Compact()` auf Sub-Components

Folgende Components zeichnen *standardmäßig* eine eigene mat-card. Im Multi-Component-Modus erzeugt das einen Card-in-Card-Look. `Compact()` lässt die innere mat-card weg und zeigt den Inhalt flach:

| Komponente | Compact-Methode |
| --- | --- |
| `stat.Stat` | `s.Compact()` |
| `barchart.BarChart` | `bc.Compact()` |
| `imagetext.ImageText` | `it.Compact()` |
| `info.InfoPoint` | `ip.Compact()` |
| `links.Links` | `l.Compact()` |
| `progress.MultiProgress` | `mp.Compact()` |

## Stat (`component/stat`)

Einzelne Statistik-Kachel.

```go
s := stat.New("42", "devices.total")
s.Icon("devices")
s.IconColor(core.ColorPrimary)
s.Suffix("km")
s.SetTrend(5.2, stat.TrendUp)
s.Print(ctx)

// AJAX-Stat
s := stat.New("", "devices.total")
s.SetURL("/api/stats/devices")
s.WithReload(true)

// Compact-Stat (für Multi-Component-Cards)
// Skipt die eigene mat-card-Hülle und nutzt kleinere Schrift, damit Stats
// flach in einer äußeren Card sitzen (kein Card-in-Card-Look).
s := stat.New("18h", "Today").Compact().WithDisplay("xcol-6")
```

## StatGrid (`component/statgrid`)

Grid-Layout für mehrere eigenständige Stat-**Karten**.

> **Für zusammengehörige Zahlen → `MultiStat` (siehe unten), nicht StatGrid.**
> StatGrid ist für separate, je-für-sich-stehende KPIs. Ein Satz verwandter
> Kennzahlen (z. B. „offen / läuft / fertig") gehört in **eine** MultiStat-Karte.

```go
sg := statgrid.New()
sg.Title("dashboard.stats")
sg.Columns(3)
sg.Add(stat.New("42", "devices.total").Icon("devices"))
sg.Add(stat.New("3", "devices.offline").Icon("warning").IconColor(core.ColorError))
sg.Add(stat.New("98%", "devices.uptime").Icon("check_circle").IconColor(core.ColorSuccess))
sg.Print(ctx)
```

## MultiStat (`component/multistat`)

> **Standard für mehrere zusammengehörige Zahlen.** Baue **nicht** mehrere `stat`
> nebeneinander (eigene `xcol`-Spalten) und **kein** `StatGrid` mit N Karten, wenn
> die Zahlen inhaltlich zusammengehören — nimm **ein** `MultiStat`. Das ist
> kompakter, hat einen gemeinsamen Header, ist responsiv (Items brechen auf
> schmalen Breiten sauber um) und liest sich als eine Kennzahlen-Einheit.

Mehrere Kennzahlen in **einer** Karte — je mit eigener Farbe/Icon/Trend, unter
gemeinsamem Header (Titel + Kopf-Icon). Anders als `StatGrid` (N separate Karten)
sitzen alle Items nebeneinander in derselben Karte. Item-Typ ist `*stat.Stat` —
alle stat-Setter (`Icon`, `Color`, `Prefix`, `Suffix`, `SetTrend`, `Link`) gelten pro Zahl.

```go
ms := multistat.New()
ms.Title("dashboard.orders.today")
ms.Icon("shopping_cart")
ms.IconColor("primary")
ms.Add(stat.New(12, "orders.open").Icon("inventory").Color("orange"))
ms.Add(stat.New(8, "orders.running").Icon("sync").Color("blue"))
ms.Add(stat.New(45, "orders.done").Icon("check_circle").Color("green").SetTrend(5, stat.TrendUp))
ms.Print(ctx)
```

Farbe/IconColor sind hier freie Strings (wie bei `stat`), nicht `core.Color`.
Header (`Title`/`Icon`/`IconColor`) ist optional — ohne rendert die Karte nur die
Zahlen. Items werden standardmäßig **horizontal** angeordnet (Icon links neben
Wert/Label); `VerticalItems()` stapelt sie stattdessen (Icon oben).

**AJAX-Nachladen** (`SetURL` + `WithReload`, wie `card`): Items werden vom Backend
geladen, Header bleibt sofort sichtbar, `WithReload(true)` zeigt einen Reload-Button.

```go
ms := multistat.New().Title("Live-Kennzahlen").Icon("monitoring").IconColor("accent")
ms.SetURL(c.apiUrl("multistat", "live"))   // Frontend lädt Items per POST
ms.WithReload(true)
// Endpoint liefert die Items via ms.DataResponse(ctx) bzw. wc.Data(ms)
```

**Klickbare Zahl** (`stat.Link`): macht einen Wert zu einem SPA-Navigations-Link
(inkl. Query-Parametern). Unterscheidet sich von `stat.SetURL` (= AJAX-Datenquelle).

```go
ms.Add(stat.New(12, "Offen").Icon("inventory").Color("orange").
    Link(xurl.NewUrl("/Orders/Table?status=open")))
```

## Tabs (`component/tabs`)

Tab-Container.

```go
t := tabs.NewTabs()

tab1 := tabs.NewTab("tab.overview").WithIcon("info")
tab1.AddContent(overviewComponent)
t.AddTab(tab1)

tab2 := tabs.NewTab("tab.details").WithLazy(true)   // Content on-demand
tab2.AddContent(detailComponent)
t.AddTab(tab2)

// Chain-Methoden auf *Tabs:
t.WithSelectedIndex(0)
t.WithDynamicHeight(true)                                   // Höhe an Content anpassen
t.WithAnimationDuration("300ms")
t.WithLazy(true)                                             // Default-Lazy für alle Tabs
t.WithUnload(true)                                           // Content bei Tab-Wechsel zerstören
t.WithHeaderPosition(core.TabHeaderPositionAbove)            // Above | Below
t.WithAlignTabs(core.TabAlignmentStart)                      // Start | Center | End
t.WithStretchTabs(true)                                      // Volle Breite
t.WithDisplay("xcol-md-12")

t.Print(ctx)
```

Tab-Chain-Methoden: `.WithIcon(string)`, `.WithDisabled(bool)`, `.WithLazy(bool)` (override), `.WithUnload(bool)` (override), `.AddContent(core.Component)`.

## Expansion (`component/expansion`)

Akkordeon-Panels.

```go
e := expansion.NewExpansion()

panel := expansion.NewPanel("section.general").
    WithIcon("settings").
    WithDescription("Allgemeine Einstellungen").
    WithExpanded(true)                                       // Initial offen
panel.AddContent(formComponent)
e.AddPanel(panel)

// Chain-Methoden auf *Expansion:
e.WithMulti(true)                                            // mehrere Panels gleichzeitig offen
e.WithDisplayMode(core.ExpansionDisplayModeDefault)          // Default | Flat
e.WithTogglePosition(core.ExpansionTogglePositionBefore)     // Before | After (Pfeil-Position)
e.WithHideToggle(false)                                      // true versteckt den Expand-Pfeil
e.WithLazy(true)                                             // Panel-Content on-demand
e.WithUnload(true)                                           // Content bei Close zerstören
e.WithDisplay("xcol-md-12")

e.Print(ctx)
```

Panel-Chain-Methoden: `.WithDescription(string)`, `.WithIcon(string)`, `.WithDisabled(bool)`, `.WithExpanded(bool)`, `.WithLazy(bool)` (override), `.WithUnload(bool)` (override), `.AddContent(core.Component)`.

## Section (`component/section`)

Abschnitt mit Titel, Icon und verschachtelten Komponenten.

```go
s := section.New()
s.Title("section.devices")
s.Subtitle("section.devices.sub")
s.Icon("devices", core.ColorPrimary)
s.Collapsible(false)
s.Buttons(buttonLine)
s.Add(tableComponent)
s.Print(ctx)
```

## Form (`component/form`)

Formular-Komponente (rendert FormGroup-Fields).

```go
f := form.NewForm(fields, "/api/vehicle/save", "vehicle.add", buttons, "", ctx)
// fields: []map[string]interface{} aus fg.ExportForFrontendWithValues(defaults)
// buttons: []*button.Button (nil = Standard Back+Save)
f.Print(ctx)
```

## Dialog (`component/dialog`)

Modale Dialoge.

```go
// Frage-Dialog (Löschen bestätigen)
content := &dialog.DialogQuestionContent{Icon: "delete", Question: "Wirklich löschen?"}
d := dialog.NewDialog(core.DialogTypeQuestion, "dialog.delete", content,
    []*button.Button{
        button.NewSimpleCloseButton("cancel"),
        button.NewSimpleApiButton("delete", "/api/device/123/delete", core.ColorError),
    }, nil, nil)
d.Print(ctx)

// Form-Dialog
d := dialog.NewDialog(core.DialogTypeForm, "dialog.edit", formComponent, buttons, nil, nil)
```

## DescriptionList (`component/descriptionlist`)

Schlüssel-Wert-Liste.

```go
dl := descriptionlist.New()
dl.Columns(2)
dl.Add("Name", "Device 1")
dl.Add("Status", "Online").Color(core.ColorSuccess).Type("badge")
dl.Add("IP", "192.168.1.1").Icon("lan")
dl.Print(ctx)
```

## Timeline (`component/timeline`)

Zeitstrahl — vertikal (Default) oder horizontal.

```go
tl := timeline.New()
tl.Add("Erstellt").Description("Von Admin").Datetime("2024-01-15").Icon("add").IconColor(core.ColorSuccess)
tl.Add("Bearbeitet").Description("Name geändert").Datetime("2024-01-16").Icon("edit")
tl.Print(ctx)

// Horizontal (statt vertikal)
tl.WithOrientation(core.TimelineOrientationHorizontal)
// core.TimelineOrientationVertical = "vertical"    (Default)
// core.TimelineOrientationHorizontal = "horizontal"

// Layout im 12-Spalten-Grid
tl.WithDisplay("xcol-md-12")

// Item-Chain-Methoden:
// .Description(string) .Datetime(string) .Icon(string) .IconColor(string)
```

## BarChart (`component/barchart`)

Bar-Chart mit drei Modi (`ModeSimple`, `ModeStacked`, `ModeHeatmap`). Frontend rendert via **echarts v6** (in `xiri-ng` als optionale peerDependency — nur installieren, wenn Bar-Charts verwendet werden).

```go
import "github.com/xiriframework/xiri-go/component/barchart"

// Simple — ein Wert pro Kategorie
bc := barchart.New("weekly").
    Mode(barchart.ModeSimple).
    Title("Weekly activities").
    YAxis(0, 12).
    Color(core.ColorPurple).
    Bar("M", 3).Bar("T", 9).Bar("W", 5)

// Stacked — pro Kategorie mehrere farbige Segmente
bc := barchart.New("strain").
    Mode(barchart.ModeStacked).
    Title("Vehicle strain").
    YAxis(0, 4).
    StackedBar("M",
        barchart.Seg(2, core.ColorGreen),
        barchart.Seg(1, core.ColorYellow),
        barchart.Seg(1, core.ColorRed))

// Heatmap — viele dünne Bars über Zeit (timeMs = unix milliseconds)
bc := barchart.New("engine").
    Mode(barchart.ModeHeatmap).
    Title("Engine system").
    Color(core.ColorPurple)
for _, s := range samples {
    bc.Point(s.TimestampMs, s.Value)
}
```

### Tooltip-Vollnamen (axis-Label kurz, Tooltip ausgeschrieben)

Pro Mode gibt es eine *Named*-Variante, die einen separaten Tooltip-Text annimmt.
Die kurze `label` bleibt auf der X-Achse, `name` erscheint im Tooltip.

```go
// Simple
bc.BarNamed("M", "Monday", 3).BarNamed("T", "Tuesday", 9)

// Stacked — Bar-Name + optional Segment-Namen
bc.StackedBarNamed("M", "Monday",
    barchart.SegNamed(2, "Low strain", core.ColorGreen),
    barchart.SegNamed(1, "Medium strain", core.ColorYellow),
    barchart.SegNamed(1, "High strain", core.ColorRed))

// Heatmap
bc.PointNamed(timestampMs, "Repeat #3", 1.0)
```

`Bar/StackedBar/Point` (ohne Suffix) bleiben gültig — `name` wird nur in der JSON-Ausgabe gesetzt, wenn explizit gegeben.

### Klickbare Balken (Navigation)

`Link(u *url.Url)` macht den **zuletzt hinzugefügten** Balken/Punkt klickbar — das Frontend navigiert beim Klick zur URL. Direkt nach dem zugehörigen `Bar`/`StackedBar`/`Point` aufrufen. Funktioniert in allen drei Modi. Es ist ein Frontend-Link → mit `url.NewUrl(...)` (ohne Prefix) erstellen; serialisiert via `Print()` als `url` im Bar-JSON. Balken ohne `Link` sind nicht klickbar (kein Cursor-Pointer).

```go
bc := barchart.New("weekly").Mode(barchart.ModeSimple).
    Bar("M", 3).Link(xurl.NewUrl("/day/mon")).
    Bar("T", 8).Link(xurl.NewUrl("/day/tue")).
    Bar("W", 6)   // ohne Link → nicht klickbar

// http(s)-URLs öffnen im Frontend in einem neuen Tab, interne Routen via Router.
```

### Wert-/Text-Labels auf den Balken

`ShowValues()` zeigt pro Balken ein Label. Standard ist der Wert; `LabelText(s)` überschreibt das Label des **zuletzt hinzugefügten** Balkens mit eigenem Text (z. B. `"8h"`). Position via `WithLabelPosition(LabelPositionTop|LabelPositionInside)` — Default `top` (über dem Balken), `inside` schreibt weißen Text in den Balken. In `ModeStacked` wird jedes Segment mit seinem Wert beschriftet (immer innen); `LabelText` greift dort am Balken, ist aber primär für simple gedacht. `ModeHeatmap` ignoriert Labels.

```go
bc := barchart.New("weekly").Mode(barchart.ModeSimple).
    ShowValues().                                     // Werte über den Balken
    WithLabelPosition(barchart.LabelPositionInside).  // optional: im Balken
    Bar("M", 3).LabelText("3h").                      // eigener Text statt "3"
    Bar("T", 8)                                       // zeigt "8"
```

### X-Achsen-Labels drehen (lange Beschriftungen)

`XLabelRotate(deg int)` dreht die X-Achsen-Beschriftung um `deg` Grad (gültig **-90..90**; `90` = vertikal, `60`/`45` = schräg). Hilft bei schmalen Bars mit langen Labels — das Frontend erzwingt dabei die Anzeige **aller** Labels (kein Ausdünnen). Kein Effekt in `ModeHeatmap` (Zeitachse).

```go
bc := barchart.New("weekly").
    Bar("Montag", 3).Bar("Dienstag", 8).Bar("Mittwoch", 6).
    XLabelRotate(60)   // schräg; 90 = vertikal
```

Dieselbe Methode gibt es auf `linechart.LineChart` (auch via `LineBuilder`) und `heatmap.Heatmap` (dreht dort die X-/Spalten-Labels).

### Optionen / AJAX

```go
bc.WithDisplay("xcol-md-6")          // Grid-Layout-Klasse
bc.SetURL(xurl.NewUrlPrefix("/api/strain", "/api"))  // Daten kommen lazy via URL
bc.WithReload(true)                  // periodisches Reload im AJAX-Modus
bc.PrintData(ctx)                    // nur das data-Objekt (für AJAX-Endpoints)
bc.DataResponse(ctx)                 // {"data": ...} envelope für DataResult
```

### Konstanten

```go
barchart.ModeSimple | ModeStacked | ModeHeatmap
barchart.Segment{Value, Color, Name}
barchart.Seg(value, color)                       // ohne Tooltip-Name
barchart.SegNamed(value, name, color)            // mit Tooltip-Name
```

## LineChart (`component/linechart`)

Multi-Linien-Chart mit gemeinsamer Kategorie-X-Achse. Tooltip zeigt alle Linien-Werte am gehoverten Punkt. Verwendet die geteilte echarts-Basis (`xiri-echarts-host`).

```go
import "github.com/xiriframework/xiri-go/component/linechart"

lc := linechart.New("revenue").
    Title("Monthly revenue").
    XLabels("Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug").
    Line("Product A", []float64{100, 120, 150, 180, 170, 195, 205, 210}).Color(core.ColorBlue).Done().
    Line("Product B", []float64{80, 95, 110, 130, 125, 145, 150, 160}).Color(core.ColorGreen).Done().
    Line("Forecast",  []float64{90, 108, 130, 155, 148, 170, 178, 185}).Color(core.ColorPurple).Dashed().Done().
    Smooth().                  // alle Linien als geglättete Kurve
    YAxis(0, 220).
    Compact()                  // ohne eigene mat-card-Hülle (für Multi-Component-Cards)
```

Pro-Linien-Optionen via `LineBuilder`-Chain: `.Color(c)`, `.Dashed()`, `.Area()` (Fill unter der Linie). `.Done()` gibt den `LineChart` zurück für weitere Kette.

`XLabelRotate(deg int)` dreht die X-Achsen-Labels um `deg` Grad (-90..90; `90` = vertikal) — hilft bei langen Kategorie-Beschriftungen. Auch via `LineBuilder` verfügbar (`.XLabelRotate(45)`).

AJAX-Mode: `lc.SetURL(u).WithReload(true)`, `lc.PrintData(ctx)`, `lc.DataResponse(ctx)`.

## PieChart (`component/piechart`)

Pie- oder Donut-Chart. Tooltip zeigt Name + Wert + Prozent.

```go
import "github.com/xiriframework/xiri-go/component/piechart"

pc := piechart.New("traffic").
    Title("Traffic sources").
    Slice("Direct", 1234, core.ColorBlue).
    Slice("Search", 856,  core.ColorGreen).
    Slice("Social", 423,  core.ColorPurple).
    Slice("Email",  187,  core.ColorOrange).
    Slice("Other",  92,   core.ColorGray)

donut := piechart.New("storage").
    Title("Storage usage").
    Donut().                   // Ring statt Pie
    Slice("Used", 78, core.ColorWarning).
    Slice("Free", 22, core.ColorLightGray).
    Compact()

// Nightingale (Rose-Style — Slice-Radien skalieren mit dem Wert)
rose := piechart.New("severity").
    Title("Bug severity").
    Nightingale().             // 'radius'-Skalierung (Default)
    // .NightingaleArea()      // alternativ 'area'-Skalierung
    Slice("Critical", 8,  core.ColorRed).
    Slice("High",     22, core.ColorWarning).
    Slice("Medium",   41, core.ColorYellow).
    Slice("Low",      73, core.ColorGreen)
```

## GaugeChart (`component/gaugechart`)

Einzelwert auf einem Bogen.

```go
import "github.com/xiriframework/xiri-go/component/gaugechart"

g := gaugechart.New("cpu").
    Title("CPU").
    Value(72).
    Range(0, 100).             // Default 0..100
    Color(core.ColorWarning).
    Label("%").                // Einheits-Label unter dem Wert
    Compact()                  // ohne Achsen-Beschriftung + ohne eigene Card
```

## Heatmap (`component/heatmap`)

2D-Matrix-Heatmap mit X/Y-Kategorien.

```go
import "github.com/xiriframework/xiri-go/component/heatmap"

h := heatmap.New("activity").
    Title("Activity by hour and weekday").
    XLabels("00", "06", "12", "18").
    YLabels("Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun").
    Range(0, 30).                                 // visualMap min/max (optional, sonst auto)
    ColorRange("#ede9fe", "#7c3aed").             // CSS low → high (optional)
    ShowValues()                                  // Werte in Zellen anzeigen (optional)

for _, m := range measurements {
    h.Cell(m.HourIdx, m.DayIdx, m.Value)
}
```

`XLabelRotate(deg int)` dreht die X-/Spalten-Labels um `deg` Grad (-90..90; `90` = vertikal) — nützlich für lange Spalten-/Bereichsnamen.

## Calendar (`component/calendar`)

GitHub-Style Calendar-Heatmap.

```go
import "github.com/xiriframework/xiri-go/component/calendar"

c := calendar.New("contributions").
    Title("2025 contributions").
    Year("2025").                                 // oder Range("2025-01-01", "2025-06-30")
    MinMax(0, 8).
    ColorRange("#e6f4ea", "#2e7d32")

for _, day := range days {
    c.Cell(day.Date /* "YYYY-MM-DD" */, day.Count)
}
```

## Tree (`component/tree`)

Hierarchischer Baum.

```go
import "github.com/xiriframework/xiri-go/component/tree"

root := tree.NewNode("Company").WithValue(200).AppendChild(
    tree.NewNode("Engineering").WithValue(80).AppendChild(
        tree.NewNode("Frontend").WithValue(30),
        tree.NewNode("Backend").WithValue(30),
        tree.NewNode("DevOps").WithValue(20),
    ),
    tree.NewNode("Sales").WithValue(60),
    tree.NewNode("Operations").WithValue(40).Collapse(),  // initial eingeklappt
)

t := tree.New("org").
    Title("Org chart").
    Root(root).
    Orient("LR").                                 // 'LR' (default), 'RL', 'TB', 'BT'
    Layout("orthogonal")                          // oder "radial"
```

## Sankey (`component/sankey`)

Flussdiagramm mit gerichteten Verbindungen.

```go
import "github.com/xiriframework/xiri-go/component/sankey"

s := sankey.New("energy").Title("Energy flow").
    Node("Coal",  core.ColorGray).
    Node("Gas",   core.ColorBlue).
    Node("Solar", core.ColorYellow).
    Node("Grid",  core.ColorPurple).
    Node("Industry", "")                         // "" = kein expliziter Node-Color

s.Link("Coal",  "Grid", 30).
  Link("Gas",   "Grid", 25).
  Link("Solar", "Grid", 15).
  Link("Grid", "Industry", 38)

s.Vertical()                                      // optional vertikale Ausrichtung
```

## Gantt (`component/gantt`)

Gantt-Chart über echarts custom-series. Zeiten als unix milliseconds.

```go
import "github.com/xiriframework/xiri-go/component/gantt"

t0 := time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC).UnixMilli()
day := int64(86400000)

g := gantt.New("project").Title("Project plan").
    Rows("Design", "Build", "Test", "Release").
    Task(0, "Wireframes", t0,         t0+7*day,  core.ColorBlue).
    Task(0, "Mockups",    t0+5*day,   t0+12*day, core.ColorPurple).
    Task(1, "API",        t0+7*day,   t0+24*day, core.ColorGreen).
    Task(2, "QA",         t0+24*day,  t0+32*day, core.ColorWarning).
    Task(3, "Go-live",    t0+35*day,  t0+36*day, core.ColorRed).
    XRange(t0, t0+40*day)                         // optional: sichtbarer X-Bereich
```

## Eigenes Chart-Component bauen

Alle vier Charts teilen sich `chart.BaseChart` (`component/chart/`). Für einen neuen Chart-Typ:

1. Neues Package `component/<type>chart/`. Struct embedded `*chart.BaseChart`. Builder mit `.Title(t)`, `.Color(c)`, `.Compact()`, `.WithDisplay(d)`, `.SetURL(u)`, `.WithReload(r)` als Forward-Methoden.
2. `Print(ctx)` ruft `base.Envelope("<type>chart", data, nil)`. `data` enthält `base.PrintBaseData(ctx)` plus die typ-spezifischen Felder.
3. Angular-Seite: neuer Component `xiri-<type>chart` mit `[option]` an `xiri-echarts-host` (siehe linechart/piechart als Vorlage).

## EmptyState (`component/emptystate`)

Leerzustand-Anzeige.

```go
es := emptystate.New("devices", core.ColorPrimary, "Keine Geräte")
es.WithDescription("Fügen Sie ein neues Gerät hinzu")
es.WithButton(button.NewSimpleLinkButton("add", "/devices/add", core.ColorPrimary))
es.Print(ctx)
```

## MultiProgress (`component/progress`)

Fortschrittsbalken mit mehreren Zeilen.

```go
mp := progress.NewMultiProgress("Verteilung", 5, false, "")
mp.AddLine("Online", 42, core.ColorSuccess, "42")
mp.AddLine("Offline", 3, core.ColorError, "3")
mp.AddTotal("Gesamt", 45, core.ColorPrimary, "45")
mp.Print(ctx)
```

## List (`component/list`)

Listen-Komponente mit Sektionen.

```go
section := list.NewListSection("Favoriten", nil)
section.AddItem(list.NewSimpleListSectionItem("Device 1", "Online", "/device/1", "devices", core.ColorSuccess))
section.AddItem(list.NewSimpleListSectionItem("Device 2", "Offline", "/device/2", "devices", core.ColorError))

l := list.NewList(nil, "")
l.AddSection(section)
l.Print(ctx)
```

## Links (`component/links`)

Link-Karte.

```go
lk := links.New()
lk.Header("navigation.title")
lk.HeaderIcon("menu", core.ColorPrimary)
lk.Add(button.NewSimpleLinkButton("Geräte", "/devices", core.ColorPrimary))
lk.Add(button.NewSimpleLinkButton("Benutzer", "/users", core.ColorSecondary))
lk.Print(ctx)
```

## Toolbar (`component/toolbar`)

Toolbar mit Suche und Buttons.

```go
tb := toolbar.New()
tb.Title("devices.toolbar")
tb.Icon("devices")
tb.Search("Suchen...")
tb.Buttons(buttonLine)
tb.Print(ctx)
```

## Query (`component/query`)

Filter-Bereich für Tabellen.

```go
q := query.NewQueryWithFormGroup(filterGroup, filterValues, "/api/devices/table", nil, "devices-filter", nil)
q.Collapsed(true)
q.Print(ctx)
```

## Stepper (`component/stepper`)

Mehrstufiger Wizard.

```go
s, err := stepper.NewStepper(
    "/api/wizard/save",
    2,
    []string{"Schritt 1", "Schritt 2"},
    []stepper.StepFields{step1Fields, step2Fields},
    "Zurück", "Weiter", "Fertig", "",
)
s.Print(ctx)
```

## Layout-Helfer (`component/layout`)

```go
layout.NewSpacer("")                                    // Leerraum
layout.NewContainer("")                                 // Container für Verschachtelung
layout.NewHeader("Titel", core.ColorPrimary, nil, "")   // Überschrift
layout.NewDivider().Text("Abschnitt").Spacing("large")  // Trennlinie
layout.NewHtml("<b>HTML</b>", "")                        // Raw HTML
```

## Button (`component/button`)

### Konstruktoren im Überblick

| Konstruktor                 | Action                       | Wann verwenden                                    |
| --------------------------- | ---------------------------- | ------------------------------------------------- |
| `NewButton`                 | beliebig                     | maximale Kontrolle — Low-Level                    |
| `NewLinkButton`             | `ButtonActionLink`           | Angular-Route (`routerLink`)                      |
| `NewHrefButton`             | `ButtonActionHref`           | Externer Link / `<a href>`                        |
| `NewApiButton`              | `ButtonActionApi`            | POST, keine Response-Aktion                       |
| `NewDialogButton`           | `ButtonActionDialog`         | POST → öffnet MatDialog mit Response              |
| `NewFormButton`             | `ButtonActionForm`           | Form-Submit (Save)                                |
| `NewDownloadButton`         | `ButtonActionDownload`       | POST → File-Download                              |
| `NewCloseButton`            | `ButtonActionClose`          | Dialog schließen                                  |
| `NewBackButton`             | `ButtonActionBack`           | Zurück-Navigation                                 |
| `NewTableButton`            | beliebig, `ButtonTypeIcon`   | Icon-only — für Table-Top/Select-Buttons          |

#### Datei anzeigen statt herunterladen

`WithTarget("_blank")` an einem Download-Button lässt das Frontend die Datei in einem neuen Tab
**anzeigen** statt sie zu speichern — der Fall „generiertes PDF ansehen":

```go
button.NewDownloadButton("PDF", url.NewUrl("/report/pdf"), core.ColorPrimary,
    core.ButtonTypeRaised, "", &filename, false, nil, nil).
    WithTarget("_blank")

// Table-Top/Bulk-Button
button.NewTableButton(core.ButtonActionDownload, "pdf", tblUrl, "PDF", core.ColorAccent, false, nil).
    WithTarget("_blank")

// Zellen-Button pro Zeile
builder.ButtonsField("actions", "actions", accessor).
    AddButton(0, table.FieldButtonActionDownload, "picture_as_pdf", core.ColorPrimary, "PDF").
    WithButtonTarget(0, "_blank")
```

Voraussetzung ist ein Content-Type, den der Browser rendern kann (`application/pdf`).
`Content-Disposition` spielt keine Rolle — das Frontend baut den Blob selbst. Ohne
Benutzer-Interaktion (z. B. `WithAutoLoad`) blockt der Browser das Tab; die Datei wird dann
heruntergeladen.

### Simple-Konstruktoren (3 Parameter, sane defaults)

```go
button.NewSimpleLinkButton   (text string, u *url.Url, color core.Color) *Button
button.NewSimpleDialogButton (text string, u *url.Url, color core.Color) *Button
button.NewSimpleApiButton    (text string, u *url.Url, color core.Color) *Button
button.NewSimpleFormButton   (text string, u *url.Url) *Button
button.NewSimpleCloseButton  (text string) *Button        // defaults: Primary, Stroked
button.NewSimpleBackButton   (text string) *Button
```

**Beispiele:**

```go
button.NewSimpleLinkButton("Details",  c.pageUrl("detail"), core.ColorPrimary)
button.NewSimpleDialogButton("Löschen", c.apiUrl("delete"), core.ColorError)
button.NewSimpleApiButton("Aktivieren", c.apiUrl("activate"), core.ColorSuccess)
button.NewSimpleFormButton("Speichern", c.apiUrl("save"))
button.NewSimpleCloseButton("Abbrechen")
button.NewSimpleBackButton("Zurück")
```

### Voll-Konstruktoren (mehr Kontrolle)

Alle teilen eine ähnliche Signatur — Unterschiede liegen in der `action`-Konstante und optional spezifischen Params. Typisches Beispiel:

```go
button.NewLinkButton(
    text       string,
    u          *url.Url,
    color      core.Color,
    buttonType core.ButtonType,
    hint       string,
    disabled   bool,
    tabIndex   *int,
    options    map[string]any,
) *Button
```

### Chain-Methoden auf `*Button`

```go
btn.WithHint("Details anzeigen")
btn.WithDisabled(true)
btn.WithData(map[string]any{"_csv": true})  // Custom-Payload an Frontend (siehe unten)
btn.WithAutoLoad(true)                       // Aktion einmalig automatisch beim Laden auslösen (siehe unten)
btn.WithTarget("_blank")                     // bei ButtonActionDownload: im Tab anzeigen (siehe oben)
// … siehe button.go für weitere Optionen (WithIcon, etc.)
```

`*TableButton` reicht `WithData`, `WithAutoLoad` und `WithTarget` an den darunterliegenden
`*Button` durch; alles andere via `GetButton()`.

### Custom-Payload via `WithData`

Wenn ein Button (typisch `ButtonActionApi`, `ButtonActionDownload`) zusätzliche Daten ans Backend mitsenden soll, die **nicht** im Filter stehen, nutze `WithData(map[string]any)`. Das Frontend liest die Payload unter `XiriButton.data` und merged sie beim Klick mit der Filter-Data in den POST-Body:

```go
// CSV-Download-Trigger (das Backend liest _csv im Filter-Body und schaltet auf CSV-Output)
csvBtn := button.NewTableButton(
    core.ButtonActionDownload, "csv", url, "CSV", core.ColorAccent, false, nil,
).WithData(map[string]any{"_csv": true})

// API-Call mit Kontext
btn := button.NewSimpleApiButton("Aktivieren", url, core.ColorSuccess).
    WithData(map[string]any{"reason": "manual", "auditTrail": true})
```

Verhalten:
- Leere/`nil`-Payload ⇒ kein `data`-Key im JSON.
- `WithData` schlägt ein eventuell vorhandenes `WithOption("data", …)` (gewinnt im Print).
- `TableButton` hat dieselbe Pass-Through-Methode (`*TableButton.WithData`).

> **Nicht** `WithOption("csv", true)` o. ä. verwenden — landet als Top-Level-Feld am Button-JSON, das Frontend liest dort nicht. `WithOption`/`WithOptions` sind für Custom-Payload deprecated.

### Auto-Load via `WithAutoLoad`

Standardmäßig lädt ein Filter mit Button (ButtonLine in der Query) erst beim **Klick**. Mit `WithAutoLoad(true)` triggert das Frontend die Button-Aktion **einmalig automatisch beim Laden** — sobald der Button nicht mehr disabled ist (d. h. der Filter gültig/vorhanden ist). Gedacht für datenladende Aktionen (`ButtonActionApi`, `ButtonActionDownload`), damit die Daten ohne initialen Klick erscheinen:

```go
// Filter-Button, der die Ergebnisse sofort beim Seitenaufruf lädt
reloadBtn := button.NewSimpleApiButton("Suchen", c.apiUrl("query/table"), core.ColorPrimary).
    WithAutoLoad(true)
```

Verhalten:
- Feuert **genau einmal**; spätere Filter-Änderungen lösen kein erneutes Auto-Load aus (der Nutzer klickt dann wie gewohnt).
- `false`/nicht gesetzt ⇒ kein `autoLoad`-Key im JSON (backward-kompatibel).
- `TableButton` hat dieselbe Pass-Through-Methode (`*TableButton.WithAutoLoad`).
- Alternative ohne Button: Query mit gesetzter `url` lädt bereits automatisch beim ersten gültigen Filter-Zustand (siehe `references/table-filtering.md`).

### ButtonLine (Container)

```go
bl := button.NewButtonLine("right", nil)   // "right" = rechtsbündig; "" → default "right"
bl.Add(button1)
bl.Add(button2)
bl.WithDisplay("xcol-md-12")
bl.Print(ctx)

// Partial-Refresh-Varianten:
bl.PrintData(ctx)              // nur data (für AJAX-Refresh)
bl.PrintButtons(ctx)           // nur []map für Button-Array
bl.DataResponse(ctx)           // response.DataResult wrapper
```

Details zu class-Werten (`"right"`, `"left"`, `"center"`, `"small"`) und Einsatz in PageHeader/Section: `references/pages.md`.

## ImageText (`component/imagetext`)

```go
it := imagetext.New("/api/image/device/1", "Device Info Text")
it.Header("device.name")
it.HeaderSub("device.model")
it.Print(ctx)
```

## Icon (`component/icon`)

```go
i := icon.NewIcon("check_circle", "Aktiv", core.ColorSuccess, nil)
i.WithData(map[string]any{"badge": 3})  // Custom-Payload unter Top-Level "data" (analog Button)
```

Custom-Payload-Konvention identisch zum Button: `WithData` für Frontend-Daten, `WithOption`/`WithOptions` für Custom-Payload deprecated.

## InfoText / InfoPoint (`component/info`)

```go
info.NewInfoText("Hinweis: Dieses Feature ist Beta.", "")
info.NewInfoPoint("192.168.1.1", "lan", core.ColorPrimary, nil, nil, nil, nil, nil, "")
info.NewInfoText("Hinweis: <b>Beta</b>", "").WithHtml()  // Text als HTML rendern (opt-in, auch auf InfoPoint)
```

## TachoTime (`component/tachotime`)

Fahrtenschreiber-/Lenkzeit-Diagramm. Siehe `references/tachotime.md` für Konstruktoren und Beispiel.
