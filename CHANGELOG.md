# Changelog

Alle nennenswerten Änderungen an `github.com/xiriframework/xiri-go` werden hier festgehalten.

Format nach [Keep a Changelog](https://keepachangelog.com/de/1.1.0/),
Versionierung nach [Semantic Versioning](https://semver.org/lang/de/).

## [Unreleased]

### Added

- **Übersetzbare Prüffehler.** Benutzerrelevante Feldfehler sind jetzt `*field.ValidationError`
  mit `Code` und `Params`; `Error()` bleibt der bisherige englische Text. Hat der `UiContext` eine
  `TranslateFunc`, übersetzt die FormGroup den Key `validation.<Code>` und setzt `{field}`
  (übersetzte Feldbeschriftung) sowie die Params ein, z. B. `validation.max_length` →
  `Höchstens {max} Zeichen`. Codes: `required`, `min_length`/`max_length` `{min}`/`{max}`,
  `pattern`, `min`/`max` `{min}`/`{max}`, `min_items`/`max_items` `{min}`/`{max}`, `not_allowed`,
  `not_past`, `not_future`, `min_date`, `max_date`, `range_order`. Die Texte liefert die App;
  fehlt ein Key (oder liefert die Funktion `""`), bleibt der englische Text. Greift in
  `BindAndValidate`/`BindFromMap`, `ParseValues(Sparse)`, `ValidateValues`,
  `ParseAndValidate(Sparse)` und damit auch in Tabellenfiltern. Neu: `FormGroup.FieldErrorMessage`,
  `FormGroup.WrapFieldErrors`.
- **`table.ExportNumber{Value, Decimals}`** für Zahlenzellen in CSV/Excel. Eigene CSV-/Excel-Formatter
  können `ExportNumber` oder eine native Zahl (`int`, `int64`, `float64` …) liefern, dann wird
  die Zelle als Zahl ausgegeben. Damit lässt sich z. B. eine `TextField`-Spalte mit der Anzeige `€ 1.234,56` als Zahl exportieren.

### Changed

- **CSV-/Excel-Export liefert Excel-taugliche Zahlen.** Integer, Float, Distance, Pressure und
  Speed kommen als `ExportNumber` aus dem Formatter.
  - CSV: Das Dezimalzeichen folgt dem Locale (`de`: `1234,57`, `en`: `1234.57`), ohne
    Tausendertrenner und nie in Exponentialschreibweise.
  - Excel: Die Werte sind numerische Zellen mit dem Format `#,##0` bzw. `#,##0.00`.
  - Negative Zahlen bekommen kein `'` mehr. Der Formel-Schutz für Text bleibt.
  - Ist der Wert `nil`, bleibt die Zelle leer.
  - Die CSV beginnt jetzt mit der UTF-8-BOM.
  - Wer maschinell einliest und einen Punkt erwartet, liefert per `WithCSVFormatter` einen String.

- **Banner nennt die Feldbeschriftung statt der ID**, sobald der `UiContext` eine `TranslateFunc`
  hat: `err.Error()` lautet dann `Ort: Höchstens 50 Zeichen` statt `ort: …`. Der konkrete
  Fehlertyp ist in diesem Fall ein Wrapper um `group.FieldErrors`; `errors.As(err, &fe)` und
  `response.NewErrorResponseFromError` funktionieren unverändert, eine Typ-Assertion
  `err.(group.FieldErrors)` nicht mehr — auf `errors.As` umstellen. Ohne `TranslateFunc` ändert
  sich nichts.

## [0.6.0]
### Added

- **`SetAllowedFunc` an `ModelField`/`ModelListField`: serverseitige Autorisierung gewählter IDs.**
  `func(ids []int32) bool` ersetzt die Prüfung gegen die angebotene Liste — gedacht für Felder
  mit `URL` (Server-Suche, Treeselect), deren IDs der Server sonst nicht kennt. Der Hook sieht
  weder IDs aus `Sub` (immer abgelehnt) noch den unveränderten Default noch bei `ModelField` die
  `0`; `ModelListField` ruft ihn einmal mit allen übrigen IDs auf, bei keiner übrigen ID gar nicht.
  `false` lehnt ab (`… is not an allowed option`, bei `ModelListField` ohne konkrete ID). Eigene
  Fehler (DB) loggt die App selbst und gibt `false` zurück — der Validierungstext geht an den
  Client. Den Benutzer bindet die App per Closure, `UiContext` kennt ihn nicht. Empfehlung: bei
  jeder `URL` setzen. Greift in `BindAndValidate`, Tabellenfiltern und `BindReload` (dort bleibt
  der Default).

### Security

- **Doku: HTML-Komponenten escapen Record-Daten.** Die Beispiele für `HtmlField` und
  `layout.NewHtml` (Adress-Dialog) betteten Record-Daten ungeescaped ein. Sie nutzen jetzt
  `html.EscapeString`, und `descriptionlist` `Type("html")` sowie `info.WithHtml()` warnen jetzt
  ebenfalls. xiri-ng sanitisiert HTML ab der nächsten Version selbst: Scripts, Handler, `style`,
  `id` und `data-*` fallen dabei weg.

### Fixed

- **Direkt gesetzte Defaults an `ModelField`/`ModelListField` beim Binden.** Ein per
  `f.Default = 42` (`int` statt `int32`) gesetzter Default wurde von `BindValue` nicht übernommen,
  `Value` blieb still `0` und überschrieb beim Speichern den aktuellen Wert. Jetzt wird er wie ein
  Request-Wert nach `int32` gewandelt. Ein `ModelListField`-Default als `[]int32` (statt
  `ModelListValue`) ließ `Parse` paniken und wird jetzt angenommen; andere Typen (auch beim
  Multi-`SelectField`) liefern einen Fehler statt eines Panics. Das gilt für `BindAndValidate`,
  `BindFromMap` und `BindReload` und ebenso für `FormGroup.ParseValues`/`ParseAndValidate(Sparse)`
  (auch Tabellenfilter): dort läuft der Default jetzt durch `Parse`, wie beim Binden.
  **Verhaltensänderung** in diesen Group-Pfaden: ein `TextField`-Default wird wie beim Binden
  getrimmt, ein ungültiger Default ist ein Feldfehler statt eines durchgereichten Rohwerts.

- **Direkt gesetzte Defaults an `SelectField` (single, auch Radio) und `ChipsField`.** Ein
  `f.Default = 2` (`int` bei `int32`-Optionen) bzw. `[]string{…}` wurde von `ParseAndValidate` und
  Tabellenfiltern als ungültig abgelehnt, obwohl das Binden ihn annahm. `Parse` normalisiert den
  Default jetzt wie einen Request-Wert. **Verhaltensänderung:** numerische Chips-Defaults kommen als
  `int64` und werden gegen die Optionen geprüft; ein Select-Default ohne passende Option ist schon in
  `ParseValues` ein Feldfehler.

- **`TimeRangeField`/`TimeLimitField`/`GeoformField` mit Default ließen sich nicht binden,** wenn
  das Feld im Request fehlte (z. B. per `showWhen` versteckt) — bei disabled oder `Form=false`
  **immer**: `BindAndValidate`/`BindFromMap` schickten den typisierten Default als Request-Wert
  durch `Parse` (`timerange field expects map with start/end`, beim `NewTimeLimitField` schon mit
  dem Konstruktor-Default), `BindReload` ließ `Value` still `nil`. Die Bind-Pfade (auch
  `BindSuggest`) binden den Default jetzt wie `FormGroup.ParseValues` über `Parse(nil)`.
  `IntField` bindet dabei `int`-, `int32`- und `int64`-Defaults. **Verhaltensänderung:** ein
  direkt gesetzter `int`- oder `time.Time`-Default am `TimeField` kommt in `ParseValues` und
  Tabellenfiltern als `int64` (wie ein Request-Wert).

- **`GeoformField`/`TimeLimitField` prüfen Zahlen-Strings vollständig.** `Validate` las per
  `fmt.Sscanf` nur den Anfang: `"45 ,0)) UNION…"`, `"NaN"` oder `"7abc"` gingen durch und kamen
  unverändert bei der App an. Jetzt müssen Koordinaten und Radius ganze, endliche Dezimalzahlen sein,
  Stunden/Minuten 1–2 Ziffern. Timelimit prüft das auch bei `check: false`. **Verhaltensänderung:**
  Koordinaten als String mit Leerzeichen oder Hex werden abgelehnt; ein Timelimit-Wert der App
  außerhalb von 1–2 Ziffern blockiert das Speichern jetzt auch bei `check: false`.

- **`ModelField`/`ModelListField` akzeptieren nur angebotene IDs.** Bisher prüfte `Validate` nur
  den Typ; wer am Frontend vorbei postete (curl, MCP `act`), konnte jede Datensatz-ID binden,
  auch fremde (IDOR), in Formularen wie in Tabellenfiltern. Jetzt muss die ID in `Options`
  (aus `LoaderFunc`), sonst in `List` stehen und darf nicht in `Sub` sein — dieselbe Menge, die
  exportiert wird. Der unveränderte Default (aktueller Wert des Datensatzes, sofern nicht in `Sub`) und bei `ModelField`
  die `0` („nichts gewählt") gehen immer durch — den Default daher nur aus serverseitig
  autorisierten Daten setzen, nie aus Request-Input. Ist ein `LoaderFunc` gesetzt und die Liste leer — auch weil das Formular ohne
  `UiContext` gebaut wurde und `LoadOptions` nie lief —, wird jede andere ID abgelehnt.
  **Nicht geprüft** (nur `Sub`) wird bei gesetzter `URL` (Server-Suche/Treeselect liefert IDs
  außerhalb der Liste) und bei Feldern ohne Loader und ohne `List` (z. B. Optionen erst per
  `SetReloadOn`): dort muss die App die ID selbst autorisieren, am einfachsten per `SetAllowedFunc`
  (siehe Added). Eine `int64`-ID außerhalb des int32-Bereichs wird abgelehnt. Der Fehler lautet
  `model field <id>: id <n> is not an allowed option` (bzw. `modellist field …`).

- **Disabled-Felder in `FormGroup.ParseValues`/`ParseAndValidate(Sparse)` sind nicht mehr vom
  Client setzbar.** Bisher ignorierte nur `BindAndValidate` den Client-Wert eines disabled Felds;
  der Group-Pfad und damit jeder Tabellenfilter (`LoadFilterData`) übernahm ihn. Ein gesperrter
  Filter (z. B. eigener Mandant) ließ sich per curl oder MCP `act` überschreiben. Jetzt gilt wie
  bei `Form=false` immer der Default, auch im Sparse-Pfad. **Verhaltensänderung:** Ein ungültiger
  Default an einem disabled Filter lässt den Filter-Request jetzt scheitern.

- **`TextField.Pattern` wird exportiert und geprüft.** Das Feld stand seit jeher im Struct und in
  der Doku, wurde aber nie gelesen: weder kam `pattern` beim Frontend an, noch prüfte `Validate`.
  Jetzt geht `pattern` ans Frontend, und `Validate` verlangt einen Match des ganzen Werts. Das gilt
  für `BindAndValidate`, `BindReload` und Filter; der Fehler lautet `text field <id> has invalid
  format` und steht unter `fields.<id>`. Verankert wird wie in Angulars `Validators.pattern`
  (`a|b` → `^a|b$`). Leere Werte werden nicht geprüft. Neuer Setter `SetPattern(string)`: er paniert
  bei ungültigem Pattern und bei Syntax, die nur Go versteht (`(?i)`, `\A`, `\p{…}`, `[[:alpha:]]` …).
  Ein direkt gesetztes solches `Pattern` wird nicht exportiert und lässt `Validate` scheitern.
  **Verhaltensänderung:** Apps, die `Pattern` schon setzen, bekommen es jetzt im Browser und auf dem
  Server durchgesetzt. Auch Defaults werden geprüft: Ein Altwert, der nicht passt, blockiert das
  Speichern, selbst an einem per showWhen versteckten Feld. In `BindReload` bleibt `Value` dann `nil`.

- **`IntField` mit `Subtype "pint"` lehnt negative Werte serverseitig ab.** Bisher ging nur `min=0`
  ans Frontend; wer am Frontend vorbei postet (MCP `act`, curl), konnte negative Werte speichern.
  Jetzt meldet `Validate` `int field <id> must be >= 0` (`BindAndValidate` unter `fields.<id>`,
  Filter); `BindReload` verwirft den Wert und behält den Default. Ein negatives
  `Min` wird bei `pint` auf 0 angehoben, im Export wie in der Prüfung. **Verhaltensänderung:** Ein
  negativer Default an einem `pint`-Feld blockiert jetzt das Speichern, auch wenn das Feld disabled ist.

### Removed

- **`IntField.Pattern`.** Das Feld wurde nie gelesen, weder exportiert noch geprüft; eine Zuweisung
  tat still nichts und ist jetzt ein Compile-Fehler. Bereiche über `Min`/`Max`, formatierte
  Kennungen (feste Stellenzahl, führende Nullen) über `TextField` mit `Pattern`, sonstige
  Zahlenregeln in der App nach `BindAndValidate`.

## [0.5.0]
### Changed

- **`TextField` trimmt Werte standardmäßig.** `Trim` stand seit jeher im Struct, wurde aber nie
  ausgewertet. Jetzt entfernt `Parse` führenden und nachgestellten Whitespace (`strings.TrimSpace`).
  Das wirkt in `BindAndValidate`, `BindReload`, für die Kontextfelder von `BindSuggest` (der
  zurückgegebene Suchtext bleibt roh) und für gesendete Werte in `FormGroup.ParseValues`/
  `ParseAndValidate`; `Validate` (Min/Max/Required) prüft dort den getrimmten Wert. Auch Defaults
  werden beim Binden getrimmt (disabled/Form=false-Felder, fehlender Wert). Roh bleiben der
  Export ans Frontend und Defaults, die `FormGroup.ParseValues` bei fehlendem Key oder
  Form=false direkt übernimmt. `NewTextField` und `NewTextFieldWithLength` setzen `Trim: true`,
  neuer Setter `SetTrim(bool)`. Subtype `password` wird nie getrimmt.
  **Verhaltensänderung:** Apps bekommen Werte ab jetzt getrimmt; wer Whitespace braucht, setzt
  `SetTrim(false)`. Per Struct-Literal gebaute Felder bleiben ungetrimmt (Zero-Value).

## [0.4.3]
### Fixed

- **Pflicht-`TextField` lehnt leere Werte ab.** `Validate` prüfte `Required` nur auf `nil`; `""` und
  reiner Whitespace gingen durch. Wer am Frontend vorbei postet (MCP `act`, curl), konnte so
  Datensätze ohne Pflichtwert anlegen. Jetzt meldet `BindAndValidate` das Feld wie ein fehlendes
  (`text field <id> is required`, unter `fields.<id>` via `response.NewErrorResponseFromError`).
  Betrifft alle Subtypes (text, textarea, email, …). Nebenwirkung in `BindReload`: ein Pflicht-
  `TextField` mit leerem Default und leerem Request-Wert hat danach `Value == nil` statt `""` —
  Reload-Handler, die `*f.Value` lesen, müssen auf `nil` prüfen.

## [0.4.2]
### Added

- **Paket `mcp`.** `mcp.Handler(e, opts)` liefert einen MCP-Server (Streamable HTTP, offizielles
  Go-SDK) mit zwei generischen Tools: `read_page(route)` gibt das Seiten-JSON zurück,
  `act(url, method, data)` führt ein Formular oder einen Button mit `action: "api"` aus. Aufrufe
  laufen intern über `e.ServeHTTP` durch dasselbe Echo, also durch Routing und Middleware des Hosts;
  `Cookie` und `Authorization` werden weitergereicht (`Options.ForwardHeaders`). Der Agent handelt
  damit unter der weitergereichten Session und erreicht alles unter `/api/` — das Seiten-JSON ist der
  Katalog, keine Schranke. Der Host mountet den Pfad hinter derselben Auth-Middleware wie `/api`.
  Abgelehnt statt zerstört werden Binärantworten (Downloads), nicht-UTF-8-Bodies und Antworten über
  `Options.MaxResponseBytes` (Default 1 MiB); Status ≠ 200 steht als `HTTP <code>` in der ersten
  Zeile, ab 300 ist das Ergebnis `IsError`. Nicht unterstützt: Streaming, Host-basiertes Routing und
  `HTTPSRedirect`-Middleware am selben Echo; Antwort-Header inklusive `Set-Cookie` gehen verloren.
  `url` bzw. `route` werden mit und ohne `/api`-Prefix akzeptiert — Formulare und `api`-Buttons
  exportieren ihn (`xurl.NewUrlPrefix`), `goto` und Breadcrumbs nicht; der Agent reicht beides
  unverändert weiter. Neue Abhängigkeit `github.com/modelcontextprotocol/go-sdk`.

- **`core.ComponentTypes()`.** Liefert die sortierte Liste aller `type`-Werte, die xiri-go im Komponenten-JSON
  ausgibt (als Kopie). Ein Test scannt die Literale unter `component/` dagegen und gleicht, wenn xiri-ng im
  Workspace liegt, mit dessen `COMPONENT_CATALOG` (`goBuilder`) ab.
- **Textfeld-Vorschläge (mat-autocomplete).** `TextField.SetSuggestions(...string)` exportiert
  `list: [{id, name}]` (id == name; ohne Argumente `list: []`, damit ein Reload-Patch Vorschläge
  abräumen kann). `SetSuggestionsURL(u *url.Url, felder...)` exportiert `url` und `searchWith`; das
  Frontend fragt beim Tippen `POST url {search, <felder>}` ab (aktuelle Werte der genannten, enabled
  Formularfelder) und erwartet `[{id, name}]`. Neuer Handler-Helfer `builder.BindSuggest(c, kontextfelder...)`
  liefert den Suchtext und bindet genau diese Felder lenient (Allowlist, JSON-Body). Freie Eingabe
  bleibt gültig. Braucht xiri-ng ≥ 0.4.15.
- **Feldgenaue Validierungsfehler.** `BindAndValidate`, `BindFromMap`, `ParseValues`, `ValidateValues`
  und `ParseAndValidate[Sparse]` liefern `group.FieldErrors` (`map[string]string`, implementiert `error`)
  mit allen fehlgeschlagenen Feldern statt nur dem ersten. `response.NewErrorResponseFromError(err)`
  gibt sie als `fields` in der 400-Antwort aus; xiri-ng zeigt sie ab dem gleichzeitig veröffentlichten
  Release am Feld. `err.Error()` bleibt ein Text (`id: msg; id2: msg2`), bestehende Aufrufer laufen
  unverändert.

## [0.4.1]
### Fixed

- **`LoadFilterData` liefert den Default von `Form=false`-Feldern wieder mit.** Der Sparse-Pfad aus
  v0.4.0 hat ihn mitgestrichen — ein Feld, das nicht im Formular steht, sendet der Client aber nie,
  sein Default ist also die einzige Quelle. Damit verlor die Einzelobjekt-Ansicht (`SetForm(false)`
  plus `Default = <objectID>`, das Muster hinter `GetFiltersForSingleObject`) ihren Filterwert und
  der Report brach mit „filter 'x' not found" ab. Für Formularfelder bleibt es bei v0.4.0:
  ein nicht gesendeter Key bekommt keinen Default.

## [0.4.0]
**Nächstes Release: v0.4.0** — Wire-Format-Änderung (Zellobjekte), nicht als Patch-Bump veröffentlichen.

### Changed

- **Datumsformate folgen jetzt der Locale.** Deutsch (und DeAT, DeCH, Nb, Da, Fi) rendert `24.02.2024`
  statt ISO; Hr, Pl, Cs, Ro, Tr, Bg, Sl, Sk, Sr, Ru, Uk ebenfalls mit Punkt statt Slash; Nl `24-02-2024`;
  Hu `2024.02.24`; EnGB `FormatTime` 24 h statt 12 h. Unverändert: Sv (ISO), EnUS, Ja/ZhCN und die
  Slash-DMY-Gruppe (EnGB — nur das Datum, die Uhrzeit wechselt auf 24 h —, Es, Fr, It, Pt, PtBR, El,
  ArAE). `FormatDateTime` ist jetzt immer
  `FormatDate + " " + FormatTime`. Eine Zeile pro Locale in `formatter/datetime_test.go`. (E1 aus
  `todo/10-skill-audit-funde.md`)
- **Tabellenzellen von Datums-, Dauer-, `text2*`- und `*N`-Feldern sind Zellobjekte `{"d", "v"}`.**
  `d` ist die bisherige Anzeige, `v` der Rohwert (ISO-Datum, lokale ISO-Zeit `2006-01-02T15:04:05`,
  Sekunden, Zahl; `null` = leer); Feld-JSON bekommt `"cellObject": "string" | "number"`. xiri-ng sortiert nach `v`,
  zeigt `d` und editiert `v` — der Server bekommt beim Inline-Edit `"2024-02-25"` statt `"25.02.2024"`.
  **Braucht xiri-ng ≥ 0.4.14** (ältere Clients zeigen `[object Object]`). Weitere Keys im Objekt sind
  für spätere Erweiterungen reserviert. Integer-`v` > 2^53 ziehen im Browser gleich (Anzeige bleibt exakt).
  Apps, die Rows von Hand bauen (`NewTableDataResponse`) und Feld-Metadaten dieser Typen mitschicken,
  müssen `{d, v}` liefern. Antwortet der Server für die editierte Zelle mit einem Zellobjekt, übernimmt der
  Client es; mit einem nackten Wert bleibt dieser als Anzeige `d`, `v` wird der eingegebene Wert, kein Reload;
  ohne Antwort für die Zelle zeigt der Client `v` als Text und lädt URL-Tabellen neu (statische Daten behalten
  den Text). `Updates` plus `Refresh` ist weiterhin erlaubt (erst Patch, dann Reload); nicht unterstützt
  sind nur widersprüchliche Antworten, die dieselbe Zelle doppelt beschreiben oder ein Einzel-Zellupdate
  (`table: "update"`) mit einem Refresh kombinieren — dort gewinnt der Refresh. Nackte Werte in anderen Zellen werden
  zu `{d: value, v: null}` normalisiert (Konsolenwarnung) und sortieren wie leer.
  `dateTime`-`v` ist lokale Zeit ohne Offset: in der doppelten Stunde der DST-Rückstellung ist die
  Reihenfolge undefiniert (bewusst akzeptiert). Tabellen mit `SetSaveInputUrl`, die Datums-/Dauer-/
  `text2*`-/`*N`-Spalten enthalten, posten für diese Spalten jetzt `{d, v}` statt des Anzeigestrings
  (der Client sendet die ganzen Rows).
- **PDF-Ausgabe** von Datumsfeldern folgt denselben Locale-Formaten (kein Zellobjekt, nur das Format).
- **Inline-Edit für `date`, `dateTime`, `timeLength`** sendet jetzt `v` (`"2024-02-25"`,
  `"2024-02-25T15:05:00"`, Sekunden) statt des Anzeigestrings; `nil` heißt geleert. `text2`/`textn`
  bleiben wie bisher nicht inline-editierbar.

### Added

- **`Table.Cell(ctx, fieldID, row)`** formatiert eine Zelle wie `GetData` — für `ReturnInlineEdit.Updates`
  nach einem Inline-Save. Nicht für Link- und Buttons-Felder (`nil`; `GetData` ergänzt dort Link-Split
  bzw. Menüdaten), die per Refresh aktualisiert werden. Mit Zellobjekt wird es übernommen; mit nacktem Wert
  bleibt dieser als Anzeige, kein Reload; ohne Antwort zeigt der Client den eingegebenen Wert als Text
  (URL-Tabellen laden neu, sofern die Antwort nicht selbst refresht/navigiert; statische Daten behalten den Text).
- **`formatter.ParseLocalDateTime(v, ctx)`**, `CellDateLayout`, `CellDateTimeLayout`: parsen das `v`
  eines Datums-/Zeit-Inline-Edits in der Zeitzone des Users; eine bei der DST-Vorstellung nicht
  existierende Uhrzeit ist ein Fehler, keine stille Verschiebung.
- `DateField`/`DateTimeField`/`TimeLengthField` setzen `inputType` `date`/`datetime-local`/`number`
  als Default (per `WithInputType` überschreibbar).
- **`FormGroup.ParseValuesSparse` / `ParseAndValidateSparse`**: parsen ohne Field-Defaults, validieren nur
  vorhandene Keys.

### Fixed

- **`SetDisabled(true)` erreicht jetzt das Frontend.** `GetBaseExport` schreibt `disabled` ins Feld-JSON (auch im
  `ExportPatch`). Bisher wirkte das Flag nur serverseitig — das Binding verwarf den Wert, das Feld war im UI aber
  bedienbar. xiri-ng sperrt das Control ab 0.4.14 beim Aufbau. Ältere Clients (≤ 0.4.13) werten den Key bereits
  in `reloadOn`-Patches und beim Wiederaktivieren nach globalem Disable aus und sperren zusammengesetzte Felder
  visuell — das Control bleibt dort aber beim Aufbau enabled.
  (`todo/12-disabled-feld-export.md`)
- **`LoadFilterData` setzt keine Field-Defaults mehr** für fehlende Filter-Keys. Ein `NewSelectField`-Filter
  ohne gesendeten Wert lieferte bisher die erste Option und filterte damit ungewollt. Fehlende Keys fehlen
  jetzt in der Map (`if v, ok := filters["status"]; ok`). Formulare (`BindAndValidate`) bekommen weiterhin
  Defaults. (E4 aus `todo/10-skill-audit-funde.md`)
- **`FormatTimeLengthH`** nutzt das Dezimaltrennzeichen der Locale (`1,5 h` für De, `1.5 h` für EnUS).
- **`TimeField.MinDate`/`MaxDate`** werden jetzt auch als `min`/`max` ans Frontend exportiert, sodass
  eine absolute Grenze Picker und Server-Validierung zugleich begrenzt. `Min`/`Max` haben beim Export
  weiterhin Vorrang.

### Deprecated

- `ModelField.Filter` (wirkungslos, `Params` verwenden), `TableBuilder.SetQuery` (xiri-ng liest
  `options.query` nicht), Parameter `blocked` von `NewDialogWaitingDone` (von xiri-ng ignoriert),
  `formatter.FormatTimestampFullDate` (identisch mit `FormatTimestampDateTime`).

## [0.3.11]
### Changed

- **Doku zu `NewReturnRefreshPanel()`:** Fallback-Kette präzisiert. Ohne umschließende Card lädt das Frontend ab
  `xiri-ng >= 0.4.11` die Seite neu (0.4.10: nur Konsolenwarnung). Kein Code geändert.

### Added

- **`BaseField.SetAddURL(u)` zeigt neben dem Feld einen „+“-Button, der eine neue Option per Dialog anlegt.**
  Für `SelectField`, `ModelListField`, `ModelField` und `ChipsField`; exportiert `addUrl`. GET auf die URL
  liefert einen `dialog.NewDialogForm`, der POST antwortet mit `response.NewReturnDone().WithCreated(id, name)`
  — das Frontend hängt `{id, name}` an die Optionsliste und selektiert die Option (bei Mehrfachwerten
  zusätzlich zu den bestehenden). `id` muss den JSON-Typ der vorhandenen Options-IDs haben.
  Wirkung erst mit `@xiriframework/xiri-ng >= 0.4.11`; ältere Frontends ignorieren `addUrl`.

## [0.3.10]
### Added

- **`SelectField.SetSelectAll(bool)` blendet im Multi-Select einen „Alle / Keine“-Toggle ein.** Exportiert
  `selectAll: true`, aber nur zusammen mit `SetMultiple(true)`. Das Frontend zeigt dann über der Optionsliste
  eine Checkbox, die alle aktuell sichtbaren — also bei aktiver Suche nur die gefilterten — Optionen auswählt
  bzw. abwählt; deaktivierte Optionen bleiben unangetastet.

  Wirkung erst mit `@xiriframework/xiri-ng >= 0.4.10`; ältere Frontends ignorieren das Feld stillschweigend.

### Added

- **`response.NewReturnRefreshPanel()` lädt genau eine Card neu.** Liefert `{"done": true, "refresh": "panel"}`
  (mit `WithMessage` wie die anderen Return-Typen). Das Frontend lädt die URL der nächstgelegenen Card
  (`Card.SetURL`) erneut und übernimmt Titel, Buttons und Inhalt aus `card.DataResponse(ctx)` — die
  restliche Seite bleibt stehen. Gedacht für Panels auf Detailseiten (Versicherung, Leasing, Preise …),
  deren Aktionen bisher `NewReturnRefreshPage()` zurückgeben mussten. Frontend: `xiri-ng >= 0.4.10`.

- **`expansion.Panel.SetURL(u)`, `PrintData(ctx)`, `DataResponse(ctx)` — nachladbare Expansion-Panels.** Mit
  `SetURL` druckt `Print` nur die Shell (Header-Felder, `url`, leeres `data`); der Endpoint liefert
  `panel.DataResponse(ctx)` als `{"panel": {…}}` mit Titel, Beschreibung, Icon, Buttons und Inhalt. Ein
  `NewReturnRefreshPanel()` aus einem Header-Button des Panels, aus dem Inhalt oder einer Tabellenaktion lädt genau
  dieses Panel neu. Frontend: `xiri-ng >= 0.4.10`.

### Changed

- **`Card.DataResponse(ctx)` antwortet mit `{"card": {…}}` statt `{"data": {…}}`.** Das Frontend konnte
  die bisherige Form nie sinnvoll darstellen (die Header-Felder landeten als Key/Value-Zeilen im Inhalt);
  das eigene Envelope macht eine komplette Card eindeutig von Zeilen-Antworten (`response.NewDataResponse(rows)`)
  unterscheidbar. Endpoints, die Zeilen liefern, sind nicht betroffen.

## [0.3.9]
### Added

- **`TableBuilder.SetFlat(bool)` rendert eine Tabelle rahmenlos.** Setzt `flat` in den Table-Options.
  Das Frontend lässt Elevation, Hintergrund, Radius und das untere Außen-Margin weg. Gedacht für
  Tabellen, die als Inhalt in einem Expansion-Panel oder in einer Card liegen — bisher stand dort
  ein Rahmen im Rahmen. Pendant zu `Card.WithFlat(bool)`.

  Wirkung erst mit `@xiriframework/xiri-ng >= 0.4.9`; ältere Frontends ignorieren das Feld stillschweigend.


## [0.3.8]
### Added

- **`Panel.Buttons(*button.ButtonLine)` setzt Aktions-Buttons in den Header eines Expansion-Panels.**
  Gleiches Modell wie `Section.Buttons` und `PageHeader.Buttons`; das Feld `buttons` im Panel-JSON
  enthält `{class, buttons}`. Das Frontend rendert sie rechtsbündig im Panel-Header, Klick und Tastatur
  lösen nur die Button-Aktion aus und klappen das Panel nicht um. Bisher brauchte man dafür eine
  eigene Card mit `ButtonTop` im Panel, mit doppeltem Titel und Card-Schatten.

  Wirkung erst mit `@xiriframework/xiri-ng >= 0.4.8`; ältere Frontends ignorieren das Feld stillschweigend.

- **`Card.WithFlat(bool)` rendert eine Card rahmenlos.** Setzt das Feld `flat` im Card-JSON (auch im
  AJAX-Pfad mit `SetURL`). Das Frontend lässt Schatten, Hintergrund und Radius weg und rendert den Header
  nur, wenn er Inhalt hat. Gedacht für Cards, die als Inhalt in einem Expansion-Panel liegen; wer dort
  keinen Doppeltitel will, setzt einfach keinen Header.

  Wirkung erst mit `@xiriframework/xiri-ng >= 0.4.8`.


## [0.3.7]
### Added

- **`Tab.WithNoPadding(bool)` rendert einen Tab ohne Innenabstand.**
  Das Frontend gibt jedem Tab-Body 16px Padding; für randlose Inhalte wie Tabellen ließ sich das aus
  Go bisher nicht abschalten. `WithNoPadding(true)` setzt das Feld `noPadding` im Tab-JSON; nur das
  Padding des eigenen Tab-Bodys entfällt, umgebende Cards und verschachtelte Tabs bleiben unverändert.
  Gedacht für bündige Inhalte wie Tabellen; Cards mit Elevation brauchen das Padding, weil der
  Tab-Body ihre Schatten sonst abschneidet.

  Wirkung erst mit `@xiriframework/xiri-ng >= 0.4.7`; ältere Frontends ignorieren das Feld stillschweigend.


## [0.3.6]
### Added

- **`FieldMeta` liefert jetzt `Header` und `HeaderSpan`.**
  `GetFieldMetas()` ist die einzige Sicht, die externe Renderer (PDF, Excel, eigene Exporte) auf die
  Felder haben — `WithHeader(...)` / `WithHeaderSpan(...)` landeten bisher nur im JSON-Export. Damit
  ließen sich gruppierte Spaltenköpfe außerhalb des Frontends nicht nachbauen.

  Beide Werte sind `nil`, wenn am Field nichts gesetzt wurde. Rein additiv — bestehende Consumer
  bleiben unverändert.

### Security

- **Transitive Abhängigkeiten angehoben:** `golang.org/x/crypto`, `golang.org/x/net`,
  `golang.org/x/sys`, `golang.org/x/text` und `go-isatty`.

## [0.3.5]
### Added

- **Filter einklappen auch bei automatisch gewrappten Tabellen (`SetFilterCollapsed`).**
  `Query.Collapsed(bool)` gab es schon, erreichbar war es aber nur, wenn man die Query von Hand
  mit `query.NewQueryWithFormGroup(...)` gebaut hat. Wer die Tabelle über `SetFilter` ihre Query
  selbst erzeugen ließ, kam an das Expansion-Panel nicht heran — `saveStateId`, `Display` und
  `extraData` wurden durchgereicht, `collapsed` nicht.

  `TableBuilder.SetFilterCollapsed(bool)` schließt die Lücke und behält die Tri-State-Semantik:
  nicht gesetzt = kein Panel (Filter immer offen, unverändertes JSON), `false` = Panel aufgeklappt,
  `true` = Panel eingeklappt.

  ```go
  builder.SetFilter(fg).SetFilterCollapsed(true)
  ```

- **Abhängige Felder: Optionen vom Server nachladen (`SetReloadOn`).** Bisher konnte ein Feld per
  `SetShowWhen` nur abhängig von einem anderen Wert ein- und ausgeblendet werden; sein **Inhalt**
  war statisch, sobald das Formular gerendert war. `BaseField.SetReloadOn(reloadURL, fields...)`
  markiert ein Feld jetzt als inhaltlich abhängig: ändert sich einer der genannten Feldwerte,
  postet das Frontend die Trigger-Werte an die URL und merged den zurückgelieferten Feld-Patch.
  Der klassische Fall ist ein Select `status` auf „aktiv" und ein Multiselect, das danach nur noch
  die für „aktiv" gültigen Einträge anbieten darf — eine Liste, die nur der Server kennt.

  Neu dazu: `FormGroup.ExportPatch()` exportiert genau die Felder mit `ReloadOn` (gekeyed nach
  Feld-ID, `Form=false` übersprungen), `builder.BindReload()` bindet den Reload-Request nachsichtig
  — ein Reload passiert mitten im Ausfüllen, ein leeres Pflichtfeld ist dort normal und kein Fehler
  — und `response.NewReturnFields()` liefert die Antwort als `{"fields": {...}}`. `BindReload`
  bindet jedes Feld zuerst auf seinen Default, damit ein unbrauchbarer Request-Wert nicht einen
  nil-Wert hinterlässt, wo die Business-Logik den Default erwartet; der Overposting-Schutz von
  `BindAndValidate` gilt unverändert, die Extraktion teilen sich beide.

  Beide Richtungen bleiben abwärtskompatibel: ein älteres Frontend ignoriert die unbekannten Keys
  `reloadOn`/`reloadUrl`, und ein Formular ohne `SetReloadOn` verhält sich exakt wie bisher. Das
  Nachladen selbst braucht `@xiriframework/xiri-ng` in einer Version mit `reloadOn`-Unterstützung.
  Filter erben das Verhalten, weil sie dieselben Fields rendern.

  Dokumentiert in `skills/xiri-go-expert/references/form-fields.md` (Abschnitt „Abhängige Felder"),
  inklusive der bewusst gezogenen Grenzen: es gehen nur die Trigger-Werte an den Server, abhängige
  Treeselects dürfen kein eigenes `URL` setzen, und Step-übergreifende Abhängigkeiten in
  Multi-Step-Forms werden nicht unterstützt.

## [0.3.4]
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

## [0.3.3]
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

## [0.3.2]
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

## [0.3.1]
### Security

- **Dependencies aktualisiert, 15 Dependabot-Alerts geschlossen.** Beide direkten Dependencies
  angehoben: `echo/v4` v4.15.1 → v4.15.4 und `excelize/v2` v2.10.1 → v2.11.0 (letzteres schließt
  CVE-2026-54063, High). Damit ziehen die transitiven Module automatisch über die betroffenen
  Schwellen: `golang.org/x/crypto` v0.48.0 → v0.53.0 (13 Alerts, davon 7 kritisch),
  `golang.org/x/net` v0.51.0 → v0.56.0. Zusätzlich `golang.org/x/text` → v0.39.0 (GO-2026-5970).
  `govulncheck` meldet vorher wie nachher **0 aufgerufene** Schwachstellen — betroffen war also nur
  mitkompilierter, von dieser Library nicht erreichter Code; für Consumer, die dieselben Module
  selbst nutzen, war die Exposition dennoch real.

## [0.3.0]
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

## [0.2.31]
### Added

- **`component/progress`**: `Progress` — einzelne determinate Fortschrittsanzeige („current of total"), inkl. `Indeterminate()`-Modus. Für Share-of-Sum weiterhin `MultiProgress`.
- **`component/bulletchart`**: Chart-Builder als kompakte Gauge-Alternative (`value` / `target` / `max` / `label`).
- **`component/stat`**: `Stat.Reference(text)` — muted Benchmark-/Anker-Zeile neben dem Wert (z. B. „Gate ≥ 1,1"), macht die Zahl auf einen Blick einordbar.
- **`component/multistat`**: `MultiStat` — mehrere Kennzahlen in einer KPI-Karte, je mit eigener Farbe/Icon/Trend, unter gemeinsamem Header. Item-Typ ist `*stat.Stat`. Inkl. AJAX-Nachladen via `SetURL`/`WithReload` (Card-Muster). Items standardmäßig horizontal; `VerticalItems()` stapelt sie.
- **`component/stat`**: `Stat.Link(u)` — macht den Wert zu einem SPA-Navigations-Link (mit Query-Parametern); getrennt von `SetURL` (AJAX-Datenquelle). Wirkt u. a. pro Zahl in `multistat`.
- **`form/field`**: `radio`-Feldtyp — Single-Select über `SelectField` mit `type=radio`, für kleine Optionsmengen.

[Unreleased]: https://github.com/xiriframework/xiri-go/compare/v0.3.6...HEAD
[0.3.6]: https://github.com/xiriframework/xiri-go/compare/v0.3.5...v0.3.6
[0.3.5]: https://github.com/xiriframework/xiri-go/compare/v0.3.4...v0.3.5
[0.3.4]: https://github.com/xiriframework/xiri-go/compare/v0.3.3...v0.3.4
[0.3.3]: https://github.com/xiriframework/xiri-go/compare/v0.3.2...v0.3.3
[0.3.2]: https://github.com/xiriframework/xiri-go/compare/v0.3.1...v0.3.2
[0.3.1]: https://github.com/xiriframework/xiri-go/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/xiriframework/xiri-go/compare/v0.2.31...v0.3.0
[0.2.31]: https://github.com/xiriframework/xiri-go/compare/v0.2.30...v0.2.31
