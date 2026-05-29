# xiri-go

Go library for type-safe UI components, form builders, and API response types. Generates JSON structures for the [xiri-ng](https://github.com/xiriframework/xiri-ng) Angular frontend.

## Installation

```bash
go get github.com/xiriframework/xiri-go
```

## Overview

xiri-go provides Go types and builders for generating JSON that drives the xiri-ng Angular UI library. Every UI component implements the `Component` interface:

```go
type Component interface {
    Print(ctx *core.UiContext) map[string]any
}
```

`*core.UiContext` carries per-request locale, timezone and an optional `Translate` function; it is nil-safe (pass `nil` for defaults).

### Packages

- **component/** - UI component builders (table, form, card, dialog, stepper, tabs, etc.)
- **form/** - Form field and group builders with struct binding
- **formatter/** - Number, date, and time formatting utilities
- **response/** - HTTP response helpers for Echo framework
- **types/** - Shared type definitions
- **uicontext/** - Request context with locale and timezone

## Quick Start

A complete, runnable Echo server that serves a full page (header + data table) as JSON
for the xiri-ng frontend to render. Copy, `go mod tidy`, run, done.

```go
package main

import (
    "net/http"

    "github.com/labstack/echo/v4"
    "github.com/xiriframework/xiri-go/component/core"
    "github.com/xiriframework/xiri-go/component/page"
    "github.com/xiriframework/xiri-go/component/pageheader"
    "github.com/xiriframework/xiri-go/component/table"
)

type Device struct {
    ID   int64
    Name string
}

func main() {
    e := echo.New()
    e.GET("/api/devices", devicesPage)
    e.Logger.Fatal(e.Start(":8080")) // xiri-ng dev server proxies /api here
}

func devicesPage(c echo.Context) error {
    ctx := &core.UiContext{} // per-request locale/timezone/translator; nil-safe defaults

    // Typed table over the Device struct (inline data — no extra endpoint needed).
    b := table.NewBuilder[Device]()
    b.IdField("id", "ID", func(r Device) int64 { return r.ID })
    b.TextField("name", "Name", func(r Device) string { return r.Name }).WithSort(true)
    tbl := b.Build()
    tbl.SetData([]Device{{ID: 1, Name: "Alpha"}, {ID: 2, Name: "Beta"}})

    // Assemble the page and return its JSON.
    p := page.NewPage()
    p.Add(pageheader.New("Devices").Icon("devices", core.ColorPrimary))
    p.Add(tbl)
    return c.JSON(http.StatusOK, p.Print(ctx))
}
```

The matching frontend is a single Angular "DynPage" route that fetches this URL and renders
the JSON via `xiri-dyncomponent` — see the [xiri-ng Quick Start](https://github.com/xiriframework/xiri-ng#quick-start).
For larger apps, give a table its own data endpoint via `tbl.SetURL(...)` and return
`wc.Data(tbl)` / `tbl.DataResponse(ctx)`; details in the bundled skill (below).

## Start a new project (and let Claude help)

The fastest path to a new Xiri app is to install **both** Claude Code skills, then describe
what you want to build:

- **`xiri-go-expert`** — bundled in this repo (see below). Knows every builder, response type,
  filter/pagination pattern, dialog and URL-routing convention.
- **`xiri-ng-expert`** — bundled in [xiri-ng](https://github.com/xiriframework/xiri-ng#claude-code-integration--xiri-ng-expert-skill).
  Knows the Angular side: `provideXiriServices`, `xiri-dyncomponent`, tables, forms, theming.

With both installed, Claude can scaffold the backend handlers **and** the matching Angular
shell from a single prompt — no need to read the full API by hand.

## Requirements

- Go 1.25+
- [Echo v4](https://github.com/labstack/echo) (for HTTP response helpers)

## Claude Code Integration — `xiri-go-expert` Skill

Dieses Repo enthält einen bundled [Claude Code](https://claude.com/claude-code) Skill unter `skills/xiri-go-expert/`, der Claude beim Schreiben von xiri-go-Code unterstützt. Der Skill wird **mit jedem Library-Release mit-versioniert**, sodass die Skill-Inhalte (API-Signaturen, Patterns, Konventionen) zum Code deiner Library-Version passen.

### Was der Skill kann

Sobald aktiviert, triggert der Skill automatisch, wenn du Go-Code schreibst oder Fragen stellst zu:

- **Komponenten**: Page-Layout, Tabellen, Dialoge, Formulare, Cards, Stepper, Tabs, Timeline, Stats, …
- **Tabellen**: Spalten-Typen, Top/Row/Bulk-Buttons, Inline-Edit, Server-Side-Pagination, CSV/Excel-Export
- **Formulare**: 20+ Feldtypen, `SetShowWhen`-Conditional-Visibility, `FormBuilder` für Add/Edit
- **Filter-Parsing**: Single- vs. Multi-Value aus `LoadFilterData`, GORM-Integration
- **Dialoge**: Question, Form, Table, Waiting (inkl. Polling), MassDelete/MassEdit
- **URL-Routing**: `*xurl.Url` mit Prefix-Pattern für Sidebar + API-Endpoints
- **Patterns**: Vollständiges CRUD, Inline-Edit-Flow, Bulk-Actions, Dashboard
- **UiContext, Formatter, Responses, Enums** — alles was die Library exportiert

Der Skill ist modular aufgebaut: **SKILL.md** (~160 Zeilen) als Navigation, plus 15 fokussierte Reference-Dateien, die nur bei Bedarf geladen werden (Progressive Disclosure) — sparsam mit Context-Tokens.

### Installation — Variante A: via `skills-lock.json`

Wenn dein Projekt den standardisierten `skills-lock.json`-Mechanismus nutzt (z.B. über `claude-granit` oder ähnliches Tooling), trage einen Eintrag ein, der auf einen Release-Tag verweist:

```json
{
  "version": 1,
  "skills": {
    "xiri-go-expert": {
      "source": "xiriframework/xiri-go",
      "sourceType": "github",
      "ref": "v1.2.3"
    }
  }
}
```

Der Loader materialisiert den Skill nach `.claude/skills/xiri-go-expert/` in deinem Projekt. Ersetze `v1.2.3` durch den Tag, der zu deiner installierten `xiri-go`-Version passt (`go list -m github.com/xiriframework/xiri-go`).

### Installation — Variante B: Direkt aus dem Modul-Cache

Weil der Skill bereits im geklonten/herunter­geladenen Modul liegt, kannst du ihn ohne zusätzliches Tooling als Projekt-Skill verlinken oder kopieren:

```bash
# Pfad zum gedownloadeten Modul ermitteln
XIRIGO_DIR=$(go list -m -f '{{.Dir}}' github.com/xiriframework/xiri-go)

# Als Symlink in dein Projekt (empfohlen — folgt Updates via go get)
mkdir -p .claude/skills
ln -s "$XIRIGO_DIR/skills/xiri-go-expert" .claude/skills/xiri-go-expert

# ODER: Kopieren (statisch — wird bei go get -u nicht aktualisiert)
cp -r "$XIRIGO_DIR/skills/xiri-go-expert" .claude/skills/
```

### Installation — Variante C: Global als User-Skill

Wenn du xiri-go in mehreren Projekten nutzt, lohnt sich eine globale Installation. Der Skill ist dann in jeder Claude-Code-Session verfügbar:

```bash
# Claude-Code-Plugin-Verzeichnis (Pfad je nach Installation: ~/.claude/skills/ oder ~/.claude-granit/skills/)
XIRIGO_DIR=$(go list -m -f '{{.Dir}}' github.com/xiriframework/xiri-go)
cp -r "$XIRIGO_DIR/skills/xiri-go-expert" ~/.claude/skills/
```

Aktualisieren nach einem `go get -u`: den globalen Skill-Ordner ersetzen.

### Skill-Struktur

```
skills/xiri-go-expert/
├── SKILL.md                      # Navigation + Minimal-Anker (always-loaded)
├── references/                    # On-demand (nur wenn Claude liest)
│   ├── pages.md                   # Page-Aufbau, 12-Spalten-Grid, PageHeader, Section, Farben
│   ├── tables.md                  # Alle Spalten-Typen, Top/Row/Bulk-Buttons, Inline-Edit
│   ├── table-builder.md           # FieldBuilder-Chaining Deep-Dive
│   ├── table-filtering.md         # Filter-FormGroup, Single/Multi-Parse, Pagination
│   ├── dialogs.md                 # Delete, Form, Table, Waiting + MassDelete/MassEdit
│   ├── form-fields.md             # Alle Feldtypen + SetShowWhen
│   ├── form-builder.md            # BuildAdd/BuildEdit, OnEditValueCheck, BindAndValidate
│   ├── responses.md               # Return-Types, DataResult, DataResponse
│   ├── patterns.md                # End-to-End: CRUD, Inline-Edit, Dashboard, Stepper, …
│   ├── url-routing.md             # *xurl.Url, Sidebar-Prefix-Pattern
│   ├── formatter.md               # Datum/Zahlen/Distanz/Pressure-Formatter
│   ├── uicontext.md               # UiContext, TranslateFunc, Enum-Werte
│   ├── components.md              # Nicht-Tabellen-Komponenten (Tabs, Timeline, …)
│   ├── tachotime.md               # TachoTime (Fahrtenschreiber)
│   └── enums.md                   # Alle Enum-Werte
└── evals/
    └── evals.json                 # Test-Prompts + Assertions (für skill-creator)
```

### Kompatibilität

Der Skill ist an den Source-Code **dieses Tags** gekoppelt. Ein `xiri-go-expert` von Tag `v1.2.0` kennt keine APIs, die erst in `v1.3.0` hinzugefügt wurden — immer die Tags synchron halten.

### Feedback / Verbesserungen

Der Skill lebt im selben Repo wie die Library: Wenn du Lücken findest, kann ein PR sowohl Code als auch Skill in einem Commit ergänzen. Siehe [CONTRIBUTING](CONTRIBUTING.md) (falls vorhanden) für Details.

## License

Apache-2.0
