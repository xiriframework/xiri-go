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
```

## StatGrid (`component/statgrid`)

Grid-Layout für mehrere Stats.

```go
sg := statgrid.New()
sg.Title("dashboard.stats")
sg.Columns(3)
sg.Add(stat.New("42", "devices.total").Icon("devices"))
sg.Add(stat.New("3", "devices.offline").Icon("warning").IconColor(core.ColorError))
sg.Add(stat.New("98%", "devices.uptime").Icon("check_circle").IconColor(core.ColorSuccess))
sg.Print(ctx)
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
// … siehe button.go für weitere Optionen (WithIcon, WithTarget, etc.)
```

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
```

## InfoText / InfoPoint (`component/info`)

```go
info.NewInfoText("Hinweis: Dieses Feature ist Beta.", "")
info.NewInfoPoint("192.168.1.1", "lan", core.ColorPrimary, nil, nil, nil, nil, nil, "")
```

## TachoTime (`component/tachotime`)

Fahrtenschreiber-/Lenkzeit-Diagramm. Siehe `references/tachotime.md` für Konstruktoren und Beispiel.
