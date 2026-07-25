# Changelog

Alle nennenswerten Änderungen an `github.com/xiriframework/xiri-go` werden hier festgehalten.

Format nach [Keep a Changelog](https://keepachangelog.com/de/1.1.0/),
Versionierung nach [Semantic Versioning](https://semver.org/lang/de/).

## [Unreleased]

Noch keine Änderungen seit `v0.3.0`.

## [0.3.0] - 2026-07-25

Aufarbeitung eines Sicherheits- und Korrektheits-Audits: 13 bestätigte Findings behoben.
Enthält **vier Breaking Changes** — alle sind Verhaltensänderungen an Stellen, die vorher still
falsche Ergebnisse geliefert haben. Migration jeweils unter dem Eintrag.

### Added

- **`component/core`**: `UiContext.SafeDistance()` und `UiContext.SafePressure()` — nil-sichere
  Accessor-Methoden mit Defaults (Kilometer bzw. Bar), analog zu `SafeLocale`/`SafeTimezone`.
- **`types/*`**: `All()` in `timezone`, `locale`, `language`, `distance` und `pressure` — liefert
  alle gültigen Enum-Werte, nach numerischem Wert sortiert. Ersetzt das Iterieren über die
  bisherigen öffentlichen Maps; die zurückgegebene Slice ist pro Aufruf frisch alloziert.

### Changed

- ⚠️ **BREAKING — `types/*`: Die Enum-Lookup-Maps sind nicht mehr öffentlich.** `Names`, `Symbols`,
  `LocaleStrings`, `LanguageCodes` und `TimezoneStrings` waren schreibbare Package-Level-Maps —
  jeder Consumer konnte sie global überschreiben, mit latentem Data Race bei gleichzeitigem Lesen.
  Migration: für einzelne Werte `GetName()`, `GetSymbol()`, `GetIANA()`, `GetLocaleString()` bzw.
  `GetCode()` verwenden, zum Aufzählen das neue `All()`.
- ⚠️ **BREAKING — `component/table`: `LoadFilterData` gibt Bind-Fehler zurück.** Ein fehlerhafter
  Request-Body wurde bisher geloggt und als leere Filter-Map behandelt — ein ungültiger Body konnte
  damit einen ungefilterten Full-Table-Export auslösen (Kosten- und Datenrisiko). Leere Bodies sind
  unverändert erlaubt (echos Binder liefert dafür `nil`). Aufrufer mit `filters, _ := …` sehen den
  Fehler weiterhin nicht und sollten angepasst werden.
- ⚠️ **BREAKING — `form/field`: Ganzzahl-/ID-Konvertierungen sind verlustfrei.** Bruchzahlen,
  NaN/Inf und Werte außerhalb des `int32`-Bereichs werden abgelehnt statt still trunkiert oder
  gewrappt. Betrifft `IntField`, `ModelField`, `ModelListField` und `SelectField`. Bisher wurde
  `1.9` zu ID `1` und `"3000000000"` zu `-1294967296` — beides wählte still den falschen
  Datensatz. Requests mit solchen Werten liefern jetzt einen Fehler.
- ⚠️ **BREAKING — `form/field`: `NewNumberField` gibt `(*IntField, error)` zurück.** Der
  `float64`-Default wird nicht mehr nach `int32` trunkiert, sondern als Fehler gemeldet.
- **`form/field`: `SelectField`** matcht Optionen nur noch bei exaktem Wert — `2.9` trifft nicht
  mehr Option `2`. Nebeneffekt: numerische Optionen akzeptieren jetzt einheitlich `int`, `int32`,
  `int64` und `float64` als Eingabetyp (vorher je Optionstyp unterschiedlich).
- **`component/table`: Button-Keys werden validiert.** `AddButton`/`AddMenu` akzeptieren nur Keys
  von `0` bis `1000`; der Key wird bei der Serialisierung als Slice-Index benutzt, ein negativer
  Key führte zum Panic, ein sehr großer zu einer Riesenallokation. Abgelehnte Keys werden per
  `slog.Warn` gemeldet und ignoriert; `AddMenu` trägt sie auch nicht mehr in die parallele
  Menü-Zustandsverwaltung ein.

### Fixed

- **`component/table`: nil-Context-Panics in Formattern behoben.** `core.Component.Print` erlaubt
  `ctx == nil`, die Distance-/Speed-/Pressure-/Text2-/N-Formatter dereferenzierten aber
  `ctx.Distance`/`ctx.Pressure` direkt. Nutzen jetzt `SafeDistance()`/`SafePressure()`.
- **`component/table`: CSV-/Excel-Export mutiert die Eingabedaten nicht mehr.**
  `expandNFieldColumns` ist eine reine Funktion und baut neue Slices/Row-Maps, statt `td.data` und
  die Felddefinitionen in-place zu überschreiben. Ein späteres `Print` liefert damit dieselben
  Daten, und parallele Exports auf derselben `TableDataResponse` sind race-frei.
- **`component/table`: CSV-/Excel-Header werden gegen Spreadsheet-Formeln geschützt.**
  `sanitizeExportValue` lief nur über Datenzellen; ein `=HYPERLINK(...)` im übersetzten Feldnamen
  blieb in Excel ausführbar.
- **`component/table`: Nullzeiten in typed Tables werden leer dargestellt.** `DateField`/
  `DateTimeField` riefen `.Unix()` auch auf `time.Time{}` und lieferten `-62135596800` (Jahr 0001);
  jetzt `IsZero()`-Prüfung mit Timestamp `0`, wie bei den N-Varianten.
- **`component/table`: Pressure-Felder unterstützen kPa.** Enum und `FormatPressureLocale` kannten
  kPa längst, der Tabellenformatter behandelte nur PSI und fiel bei kPa auf den unveränderten
  Bar-Wert mit Einheit `" bar"` zurück.
- **`formatter`: `FormatInteger` verliert oberhalb von 2⁵³ keine Präzision mehr.** Der `int64` wurde
  vor dem Formatieren nach `float64` gecastet. Betraf auch die Tabellen-Integer-Formatter
  (`createIntegerFormatter`, `createText2IntFormatter`, `createIntegerNFormatter`).
- **`form/group`: `FormGroup.FormatNumber` bei nil Context.** Gab die Formatvorlage `"%.2f"` statt
  des formatierten Werts zurück.
- **`form/field`: `TimeRangeField` berechnet relative Min/Max-Tagesgrenzen lokal und DST-sicher.**
  Statt `time.Now() + n*86400` wird jetzt — wie bei `TimeField` — die lokale Mitternacht in der
  Zeitzone des Nutzers als Anker verwendet und mit `AddDate` gerechnet. Vorher hing die Grenze an
  der aktuellen Uhrzeit statt am Tagesbeginn und lag über einem DST-Sprung eine Stunde falsch.
  Beide Grenzen nutzen dasselbe `now` und können nicht mehr über Mitternacht auseinanderfallen.
- **`form/field`: `TimeRangeField.ExportForFrontend` berücksichtigt den Field-Default.** Ohne
  gebundenen Wert wurde der Default ignoriert und beide Enden auf „jetzt" gesetzt; Reihenfolge ist
  jetzt Wert → Default → jetzt.
- **`component/descriptionlist` / `component/timeline`: `Add()` gibt keinen veralteten Pointer
  mehr zurück.** Die Items lagen als `[]Item` im Slice, `Add` gab `&items[len-1]` zurück — nach
  einer `append`-Reallokation zeigte dieser Pointer auf den alten Speicher. „Erst alle `Add`, dann
  konfigurieren" verlor damit Änderungen. Jetzt `[]*Item`.
- **Doku**: `UiContext.SafeTimezone` dokumentierte UTC als Default, liefert aber `Europe/Vienna`
  (konsistent mit `SafeLocale` → `locale.De`); veralteter `uicontext/`-Verweis in `README.md`
  korrigiert — `UiContext` liegt in `component/core/context.go`.

## [0.2.31] - 2026-07-18

### Added

- **`component/progress`**: `Progress` — einzelne determinate Fortschrittsanzeige („current of total"), inkl. `Indeterminate()`-Modus. Für Share-of-Sum weiterhin `MultiProgress`.
- **`component/bulletchart`**: Chart-Builder als kompakte Gauge-Alternative (`value` / `target` / `max` / `label`).
- **`component/stat`**: `Stat.Reference(text)` — muted Benchmark-/Anker-Zeile neben dem Wert (z. B. „Gate ≥ 1,1"), macht die Zahl auf einen Blick einordbar.
- **`component/multistat`**: `MultiStat` — mehrere Kennzahlen in einer KPI-Karte, je mit eigener Farbe/Icon/Trend, unter gemeinsamem Header. Item-Typ ist `*stat.Stat`. Inkl. AJAX-Nachladen via `SetURL`/`WithReload` (Card-Muster). Items standardmäßig horizontal; `VerticalItems()` stapelt sie.
- **`component/stat`**: `Stat.Link(u)` — macht den Wert zu einem SPA-Navigations-Link (mit Query-Parametern); getrennt von `SetURL` (AJAX-Datenquelle). Wirkt u. a. pro Zahl in `multistat`.
- **`form/field`**: `radio`-Feldtyp — Single-Select über `SelectField` mit `type=radio`, für kleine Optionsmengen.

[Unreleased]: https://github.com/xiriframework/xiri-go/compare/v0.3.0...HEAD
[0.3.0]: https://github.com/xiriframework/xiri-go/compare/v0.2.31...v0.3.0
[0.2.31]: https://github.com/xiriframework/xiri-go/compare/v0.2.30...v0.2.31
