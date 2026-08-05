# Changelog

Alle nennenswerten Änderungen an `github.com/xiriframework/xiri-go` werden hier festgehalten.

Format nach [Keep a Changelog](https://keepachangelog.com/de/1.1.0/),
Versionierung nach [Semantic Versioning](https://semver.org/lang/de/).

## [Unreleased]

Noch keine Änderungen seit `v0.3.4`.

## [0.3.4] - 2026-08-05

### Documentation

- **Falsche Konstruktor-Signaturen in der Skill-Doku korrigiert.** `NewTextField`, `NewIntField`,
  `NewBoolField` und `NewTimeField` nehmen ihren Default als **Wert** (`string`, `int32`, `bool`,
  `int64`) — die Doku zeigte an 17 Stellen `nil` bzw. `&entity.Feld`, was nie kompilierte.
  Besonders irreführend bei `NewTimeField`: das dokumentierte `defaultValue (*int64)` legte ein
  „kein Default = nil" nahe, das die Library so nicht anbietet — real ist `0` ein echter Timestamp
  (1970-01-01) und `Parse(nil)` liefert `0`, nicht `nil`. Betrifft alle Versionen bis `v0.3.3`.

- **`NewDeviceListField` aus der Doku entfernt** — dieser Konstruktor existiert nicht. Korrekt ist
  `NewModelListField(id, name, required, "device", ids)`.

- **TimeField-Subtype-Default korrigiert.** Ein leerer `Subtype` exportiert `type: "datetime"`,
  nicht `"date"`; `"yearmonth"` fehlte in der Liste der gültigen Subtypes.

- **Neuer Guard-Test `form/field/skilldoc_test.go`.** Prüft die ausgelieferte Skill-Doku künftig
  automatisch gegen die echten Signaturen dieses Packages: jeder dokumentierte `field.New…` muss
  existieren, und `nil`/`&x` darf nicht an einen Wert-Parameter übergeben werden. Die Regel wird
  per `go/parser` aus den tatsächlichen Signaturen abgeleitet, nicht aus einer gepflegten Liste —
  `nil` bleibt daher bei Slice/Map-Defaults (`NewModelListField`, `NewArrayField`, `NewJsonField`)
  korrekt zulässig. Läuft in `go test ./...` und damit in `./release.sh` mit.

## [0.3.3] - 2026-08-01

### Added

- **`TableBuilder.SetDensity()` — die vollständige Zeilenhöhen-API.** Bisher gab es nur
  `SetDense(bool)`, das im Frontend als Alias für `compact` gilt — `relaxed` war von Go aus gar
  nicht erreichbar. Neu ist der Typ `Density` mit `DensityCompact`, `DensityRegular` und
  `DensityRelaxed`; gesetzt landet der Wert als `options["density"]` im JSON. `SetDense` bleibt als
  Legacy-Setter unverändert funktionsfähig und emittiert weiter `options["dense"]`; sind beide
  gesetzt, gewinnt im Frontend `density` (der Alias greift nur, wenn keine explizite Density
  ankommt). Braucht ein Frontend, das `density` versteht — das ist seit `xiri-ng` v0.2.49 der Fall.

- **`Button.WithHide()` / `TableButton.WithHide()`.** Das Frontend wertet `XiriButton.hide` seit
  der zugehörigen `xiri-ng`-Version in allen Renderpfaden aus — von Go aus war das Feld bisher nur
  über den deprecateten `WithOption("hide", …)`-Umweg erreichbar. Ein versteckter Button landet
  nicht im DOM und führt seine Aktion auch nicht aus (auch nicht per `WithAutoLoad`). Kein
  Berechtigungs-Contract: der Endpoint hinter dem Button muss ohnehin abgesichert sein.

- **Warnung bei Icon-Buttons ohne `hint`.** `Print()` emittiert für `ButtonTypeIcon`, `Fab` und
  `MiniFab` **kein** `text` — bei leerem Icon wandert `text` sogar in das `icon`-Feld. Für diese
  Typen ist `hint` damit die einzige Quelle, aus der das Frontend einen Accessible Name bauen kann;
  ohne `hint` melden Screenreader einen Button ohne Namen. Gemeldet wird das per `slog.Warn` in
  `Print()` — nicht im Konstruktor, weil der Hint auch später über `WithHint()` kommen darf und
  eine Warnung dort korrekte Fluent-Nutzung anmeckern würde. Pro Button wird höchstens einmal
  gewarnt, damit wiederholtes `Print()` das Log nicht flutet.

### Documentation

- **`xiri-go-expert` skill**: `SetDensity` in der Options-Liste (`tables.md`), Hinweis dass
  `SetDense` legacy ist, und ein neuer Abschnitt zu `hide` in `XiriNavigationField`
  (`url-routing.md`) — inklusive der ausdrücklichen Klarstellung, dass das **kein**
  Berechtigungs-Contract ist und der Server die Routes trotzdem absichern muss. Dazu in
  `components.md` je ein Abschnitt zu `WithHide` und dazu, dass `hint` bei Icon-Button-Typen
  Pflicht ist, weil `Print()` für sie kein `text` emittiert.

## [0.3.2] - 2026-07-31

### Added

- **Download-Buttons können die Datei im Tab anzeigen lassen statt sie zu speichern.**
  `WithTarget("_blank")` an einem Download-Button ist nicht mehr wirkungslos: das Frontend öffnet
  die Datei damit in einem neuen Tab — der Fall „generiertes PDF ansehen" statt „PDF in den
  Download-Ordner legen". Voraussetzung ist ein Content-Type, den der Browser rendern kann
  (`application/pdf`); `Content-Disposition` ist dafür irrelevant, weil das Frontend den Blob aus
  dem Response-Body selbst baut und nur noch den Dateinamen aus dem Header liest. Braucht
  `@xiriframework/xiri-ng` mit der zugehörigen Änderung — ältere Frontend-Versionen ignorieren
  `target` bei `download` weiterhin.

  Neu für die Kontexte, die vorher keinen Weg zu `target` hatten:
  `TableButton.WithTarget()` (Table-Top- und Bulk-Buttons; bisher nur über
  `GetButton().WithTarget()` erreichbar) und `FieldBuilder.WithButtonTarget(key, target)` für
  Zellen-Buttons — letzteres schreibt in die schon vorhandenen Button-Options, die als
  Top-Level-Keys im Button-JSON landen, ein Serializer-Change war nicht nötig. Unbekannte Keys
  ignoriert der Setter, damit er nach einem von `AddButton` verworfenen Out-of-Range-Key nicht
  ins Leere schreibt.

  Ohne Benutzer-Interaktion (`WithAutoLoad`) blockt der Browser das Tab — die Datei wird dann
  heruntergeladen. Das ist Absicht, nicht abgefangen.

### Documentation

- **`xiri-go-expert` skill**: neuer Abschnitt „Datei anzeigen statt herunterladen" in
  `components.md` (alle drei Button-Varianten) und in `tables.md` (Zellen-Buttons);
  `WithTarget` in der Chain-Methoden-Liste ergänzt, samt Hinweis welche Methoden `*TableButton`
  durchreicht.

## [0.3.1] - 2026-07-25

### Security

- **Dependencies aktualisiert, 15 Dependabot-Alerts geschlossen.** Beide direkten Dependencies
  angehoben: `echo/v4` v4.15.1 → v4.15.4 und `excelize/v2` v2.10.1 → v2.11.0 (letzteres schließt
  CVE-2026-54063, High). Damit ziehen die transitiven Module automatisch über die betroffenen
  Schwellen: `golang.org/x/crypto` v0.48.0 → v0.53.0 (13 Alerts, davon 7 kritisch),
  `golang.org/x/net` v0.51.0 → v0.56.0. Zusätzlich `golang.org/x/text` → v0.39.0 (GO-2026-5970).
  `govulncheck` meldet vorher wie nachher **0 aufgerufene** Schwachstellen — betroffen war also nur
  mitkompilierter, von dieser Library nicht erreichter Code; für Consumer, die dieselben Module
  selbst nutzen, war die Exposition dennoch real.

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

[Unreleased]: https://github.com/xiriframework/xiri-go/compare/v0.3.1...HEAD
[0.3.1]: https://github.com/xiriframework/xiri-go/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/xiriframework/xiri-go/compare/v0.2.31...v0.3.0
[0.2.31]: https://github.com/xiriframework/xiri-go/compare/v0.2.30...v0.2.31
