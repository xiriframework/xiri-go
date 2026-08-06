# Form Fields (form/field package)

Import: `"github.com/xiriframework/xiri-go/form/field"`

Alle Fields implementieren `FormField` Interface und embedden `*BaseField`.

## Gemeinsame BaseField-Methoden (alle Fields)

```go
field.SetClass("xcol-md-6")      // CSS-Klasse (Grid-Breite)
field.SetHint("tooltip.text")    // Tooltip/Hilfetext
field.SetStep(1)                 // Schritt in Multi-Step-Forms
field.SetDisabled(true)          // Deaktiviert
field.SetAccess([]string{"admin"}) // Rollen-Metadaten (KEIN Zugriffsschutz, siehe unten)
field.SetScenario([]string{"add"}) // Szenario-Metadaten (KEIN Zugriffsschutz, siehe unten)
field.SetForm(false)             // Nicht im Formular anzeigen
```

> ⚠️ **`SetAccess`/`SetScenario` erzwingen nichts.** Weder der Export noch das Binding werten sie
> aus — ein manipulierter Request kann ein so markiertes Feld trotzdem setzen. Wer ein Feld
> wirklich schützen will, darf es rollenabhängig **nicht in die `FormGroup` aufnehmen**; der
> Form-Binder schützt korrekt gegen Overposting und akzeptiert nur deklarierte, aktive Felder.
>
> (Finding #2 des Audits; echte Durchsetzung steht noch aus.)

### ShowWhen (Bedingte Sichtbarkeit)

```go
// Feld nur anzeigen wenn "type" den Wert "vehicle" hat
f.ShowWhen = []field.Condition{
    field.NewCondition("type", field.CondEquals, "vehicle"),
}
// Feld nur anzeigen wenn "name" nicht leer ist
f.ShowWhen = []field.Condition{
    field.NewConditionNotEmpty("name"),
}
```

---

## TextField

Freitext-Eingabe. Value: `*string`

```go
f := field.NewTextField("name", "vehicle.name", true, "")
// Parameter: id, translationKey, required, currentValue (string — Wert, nicht *string)

// Optionale Konfiguration:
f.Subtype = "textarea"  // text (default), textarea, html, email, url, tel, password
f.MinLength = 3
f.MaxLength = 100
f.Pattern = `^[a-zA-Z]+$`
f.TextPrefix = "€"
f.TextSuffix = "kg"
f.IconPrefix = "euro"
f.IconSuffix = "weight"
f.Trim = true

// Nach BindAndValidate:
name := *f.Value  // string
```

## IntField

Zahleneingabe. Value: `*int32`

```go
f := field.NewIntField("count", "device.count", true, 0)
// Parameter: id, translationKey, required, currentValue (int32)

f := field.NewIntFieldWithBounds("count", "device.count", true, 0, 0, 100)
// Mit Min/Max-Grenzen (min, max sind int, nicht *int)

f, err := field.NewNumberField("price", "device.price", true, 42.0)
// Alias für IntField. Gibt einen Fehler zurück, wenn der Default nicht verlustfrei
// in int32 passt (Bruchzahl oder außerhalb des int32-Bereichs).

// Optionale Konfiguration:
f.Subtype = "float"   // int (default), pint (positive int), float, bigint, real
f.TextPrefix = "€"
f.TextSuffix = "Stück"

// Nach BindAndValidate:
count := *f.Value  // int32
```

Parsing ist verlustfrei: `1.9`, `"1.9"` und `"3000000000"` werden abgelehnt, nicht trunkiert oder
gewrappt. Gleiches gilt für Model-IDs (`ModelField`, `ModelListField`) und `SelectField`-Optionen.

## BoolField

Checkbox. Value: `*bool`

```go
f := field.NewBoolField("active", "device.active", false, false)
// Parameter: id, translationKey, required, currentValue (bool — Wert, nicht *bool)

// Nach BindAndValidate:
active := *f.Value  // bool
```

## SelectField

Dropdown-Auswahl. Value: `int32`

```go
options := []field.SelectOption{
    {Value: 1, Label: "status.online"},
    {Value: 2, Label: "status.offline"},
    {Value: 3, Label: "status.maintenance"},
}
f := field.NewSelectField("status", "device.status", true, options)
// Parameter: id, translationKey, required, options

f.Search = boolPtr(true)  // Suchfeld aktivieren (auto bei 20+ Optionen)

// Nach BindAndValidate:
status := f.Value  // int32
```

## ModelField

Einzelne Modell-Auswahl (z.B. Gruppe, Benutzer). Value: `int32`

```go
f := field.NewModelField("group_id", "device.group", true, "group", currentGroupID)
// Parameter: id, translationKey, required, modelType, currentValue (int32)

// Dynamische Optionen laden via LoaderFunc:
f.SetLoaderFunc(func(ctx *core.UiContext) ([]field.ModelOption, error) {
    groups, _ := db.GetGroups()
    opts := make([]field.ModelOption, len(groups))
    for i, g := range groups {
        opts[i] = field.ModelOption{ID: int32(g.ID), Name: g.Name}
    }
    return opts, nil
})

// Oder statische Liste:
f.List = []field.ModelOption{
    {ID: 1, Name: "Gruppe A"},
    {ID: 2, Name: "Gruppe B"},
}

f.AllowSearch = true
f.URL = "/api/groups/search"  // Live-Suche via API
f.Filter = map[string]interface{}{"active": true}  // Vorfilter

// Nach BindAndValidate:
groupID := f.Value  // int32
```

## ModelListField

Mehrfach-Modell-Auswahl. Value: `ModelListValue` (= `[]int32`)

```go
f := field.NewModelListField("device_ids", "route.devices", true, "device", currentIDs)
// Parameter: id, translationKey, required, modelType, currentValue ([]int32)
// nil als currentValue ist hier erlaubt und wird zu []int32{} normalisiert.

f.MinItems = intPtr(1)
f.MaxItems = intPtr(10)
f.AllowEmpty = true
f.SingleOnly = true  // Nur ein Element erlaubt
f.SetLoaderFunc(loaderFunc)

// Nach BindAndValidate:
deviceIDs := f.Value  // []int32
```

## TimeField

Datum/Zeit-Auswahl. Value: `*int64` (Unix-Timestamp in Sekunden)

```go
f := field.NewTimeField("start", "event.start", true, 0)
// Parameter: id, translationKey, required, defaultValue (int64 Unix-Sekunden — Wert, nicht *int64)
// Kein Default → 0 übergeben. Achtung: 0 ist ein echter Timestamp (1970-01-01), kein "leer" —
// das Feld kommt damit vorbelegt im Frontend an und Parse(nil) liefert 0, nicht nil.

f.Subtype = "datetime"  // datetime (default), date, time, yearmonth

// Tag-Offsets für Min/Max (z.B. -30 Tage bis +365 Tage):
f.Min = int64Ptr(-30)
f.Max = int64Ptr(365)

// Oder absolute Timestamps:
f.MinDate = &time.Time{...}
f.MaxDate = &time.Time{...}

// Nach BindAndValidate:
timestamp := *f.Value  // int64 (Unix seconds)
```

## YearMonthField (TimeField-Subtype "yearmonth")

Monats-Auswahl (Jahr+Monat). Value: `*int64` (Unix-Timestamp = 1. des Monats 00:00 in User-Timezone)

```go
f := field.NewYearMonthField("month", "report.month", true, defaultUnix)
// Parameter: id, translationKey, required, defaultValue (int64 Unix-Sekunden)
// Intern: liefert ein *TimeField mit Subtype="yearmonth" → frontend type="yearmonth"

// Erbt alle TimeField-Optionen:
f.Min = int64Ptr(-365)  // 365 Tage zurück (Day-Offset, |val| < 10000 = Tage ab heute)
f.Max = int64Ptr(0)     // bis heute
// Oder absolute Grenzen:
f.MinDate = &time.Time{...}
f.MaxDate = &time.Time{...}
f.AllowPast   = true
f.AllowFuture = false

// Nach BindAndValidate:
ts := *f.Value             // int64 Unix-Sekunden, 1. des gewählten Monats 00:00 (User-TZ)
t  := time.Unix(ts, 0)     // → time.Time
```

Hinweis: Subtype `"yearmonth"` wird vom Frontend (xiri-ng) als `XiriYearMonthComponent`
mit Multi-Year-Datepicker gerendert. Der Eingabewert wird automatisch auf den 1. des
Monats normalisiert. Parsing akzeptiert ISO-Format `"2006-01"` zusätzlich zu Date/RFC3339.

## TimeRangeField

Zeitbereich-Auswahl. Value: `*TimeRangeValue`

```go
f := field.NewTimeRangeField("range", "report.range", true)
f := field.NewTimeRangeFieldWithDefault("range", "report.range", true, 7) // Letzte 7 Tage

f.Subtype = "daterange"  // daterange (default), time
f.AllowSingleDay = true

// Nach BindAndValidate:
start := f.Value.Start  // time.Time
end := f.Value.End      // time.Time
```

## ChipsField

Tag/Chip-Eingabe im **gemischten Modus**. Value: `[]interface{}` — jedes Element ist entweder
ein numerisches Options-`int64` (Chip aus der `List` gewählt → ID) oder ein `string` (frei
eingegebener Text). Die JSON-/Go-Typunterscheidung trägt die Semantik: `number/int64` = Options-ID,
`string` = Freitext.

Die `List`-Optionen brauchen **numerische `Value`** (int64), damit sie als IDs funktionieren — die
exportierte `id` wird im Frontend (`XiriFormFieldSelectOption.id: number`) erkannt und beim Auswählen
als ID zurückgesendet. String-Option-Values würden als Freitext interpretiert.

```go
f := field.NewChipsField("tags", "device.tags", false)
f.SetList([]field.SelectOption{
    {Value: int64(1), Label: "Indoor"},
    {Value: int64(2), Label: "Outdoor"},
})
f.SetFreeText(true)  // erlaubt zusätzlich freie Texteingabe (ohne ID)

// Nach BindAndValidate — bequemer getrennter Zugriff über Helper:
ids := f.IDs()      // []int64  — aus der Liste gewählte Chips (Options-IDs)
texts := f.Texts()  // []string — frei eingegebene Chips
all := f.Value      // []interface{} — int64 + string in Eingabereihenfolge
```

Ein ChipsField **ohne `List`** enthält schlicht nur Freitext-Strings (`IDs()` ist leer).

## FileField

Datei-Upload.

```go
f := field.NewFileField("document", "device.document", true, 10*1024*1024) // 10MB
f.AllowedTypes = []string{"application/pdf", "image/png"}
f.AllowedExtensions = []string{".pdf", ".png"}
f.Multiple = true
```

## TimeLimitField

Zeitbeschränkung mit Wochentagen. Mappt auf 5 DB-Spalten.

```go
f := field.NewTimeLimitField("timelimit", "device.timelimit", false)
// Exportiert: time_check, time_weekdays, time_from, time_to, time_in
// Value: TimeLimitValue{Check, Weekdays [7]bool, FromHour, FromMin, ToHour, ToMin, In}
```

## GeoformField

Geometrie-Feld für Geofencing.

```go
f := field.NewGeoformField("geofence", "zone.geofence", true)
// Value: GeoformValue{Type (1=Polygon, 2=Kreis), Path}
```

## ArrayField

Array/Liste.

```go
f := field.NewArrayField("items", "order.items", true, "string", nil)
f.MinItems = intPtr(1)
f.MaxItems = intPtr(10)
f.UniqueItems = true
```

## JsonField

JSON-Objekt.

```go
f := field.NewJsonField("config", "device.config", false, nil)
// Value: map[string]interface{}
```

## Display-Only Fields (kein Value)

```go
field.NewHeaderField("section1", "Allgemein")       // Abschnitts-Überschrift
field.NewHeaderField("section1", "Allgemein").SetCollapsible(true).SetCollapsed(false)
field.NewHtmlField("notice", "<b>Hinweis:</b> ...")  // HTML-Inhalt
field.NewInfoField("help", "Hilfetext hier")         // Info-Text
field.NewDividerField("div1")                        // Trennlinie
field.NewSerialField("id", "ID")                     // Auto-ID (read-only)
```

## Conditional Visibility — `SetShowWhen` (Go-API)

Jedes Field erbt von `BaseField` und hat Chain-Methoden, um Sichtbarkeits-Bedingungen zu setzen — das ist der **bevorzugte Weg**, statt manuell `Condition`-Structs zusammenzubauen.

```go
// BaseField-Methoden (via Feld-Embedding auf allen Field-Typen verfügbar)
field.BaseField.SetShowWhen(field string, operator field.ConditionOperator, value interface{}) *BaseField
field.BaseField.SetShowWhenNotEmpty(field string) *BaseField
```

Mehrfach aufrufen = UND-Verknüpfung aller Bedingungen.

### Operatoren

Quelle: `form/field/condition.go`

| Konstante              | String          | Wert-Typ                   |
| ---------------------- | --------------- | -------------------------- |
| `field.CondEquals`     | `"equals"`      | beliebig (string/int/bool) |
| `field.CondNotEquals`  | `"notEquals"`   | beliebig                   |
| `field.CondContains`   | `"contains"`    | string                     |
| `field.CondGreater`    | `"greaterThan"` | Zahl                       |
| `field.CondLess`       | `"lessThan"`    | Zahl                       |
| `field.CondIn`         | `"in"`          | **`[]any`** (Liste)        |
| `field.CondNotEmpty`   | `"notEmpty"`    | — (kein value)             |

### Beispiele

```go
// Einfach: nur wenn "active" = false
reasonField := field.NewTextField("reason", "Grund", false, "")
reasonField.BaseField.SetShowWhen("active", field.CondEquals, false)

// Shortcut "notEmpty"
notesField := field.NewTextField("notes", "Notizen", false, "")
notesField.BaseField.SetShowWhenNotEmpty("name")

// UND-Verknüpfung (alle Bedingungen müssen wahr sein)
critField := field.NewTextField("criticalReason", "Kritisch-Grund", true, "")
critField.BaseField.
    SetShowWhen("priority", field.CondIn, []any{"high", "critical"}).
    SetShowWhen("active", field.CondEquals, true)
```

### Pattern im Form-Builder

```go
fb := formbuilder.NewFormBuilder(uc)

active := field.NewBoolField("active", "Aktiv", false, false)
reason := field.NewTextField("reason", "Abschalt-Grund", false, "")
reason.BaseField.SetShowWhen("active", field.CondEquals, false)

prio := field.NewSelectField("priority", "Priorität", true, priorityOpts)
critNote := field.NewTextField("critNote", "Kritisch-Notiz", false, "")
critNote.BaseField.SetShowWhen("priority", field.CondIn, []any{int32(3), int32(4)})

fb.AddField(active).AddField(reason).AddField(prio).AddField(critNote)
```

### Runtime-Verhalten

- Das Frontend wertet die Bedingungen live aus (reactive) — Ein-/Ausblenden ohne Roundtrip.
- Beim Submit werden **nur** Werte sichtbarer Felder mitgeschickt — versteckte Felder sind bewusst leer.
- Deshalb nach `BindAndValidate`: wenn ein Feld ausgeblendet war, ist `field.Value == nil` — prüfen, bevor man dereferenziert.

### Low-Level — direkt `Condition` bauen (selten nötig)

```go
field.NewCondition("fieldId", field.CondEquals, "value")
field.NewCondition("fieldId", field.CondIn, []any{"a", "b"})
field.NewConditionNotEmpty("fieldId")
```

Direkt `Condition`-Structs baust du nur dann, wenn du eine Liste von Bedingungen programmatisch zusammenstellst und dann `BaseField.ShowWhen = conditions` setzt.

---

## Abhängige Felder — `SetReloadOn` (Optionen vom Server nachladen)

`SetShowWhen` blendet ein Feld abhängig von einem anderen Wert ein und aus. `SetReloadOn` geht
einen Schritt weiter: es lädt den **Inhalt** eines Felds neu, wenn sich ein anderes Feld ändert.

Typischer Fall: ein Select `status` steht auf „aktiv", und ein Multiselect darf danach nur noch
die für „aktiv" gültigen Einträge anbieten. Diese Liste kennt nur der Server.

```go
func (f *BaseField) SetReloadOn(reloadURL *xurl.Url, fields ...string) *BaseField
```

Beide Angaben sind Pflicht. Ein `nil`-URL oder eine leere Feldliste ergibt eine Abhängigkeit,
die das Frontend nicht auflösen kann — dann wird **nichts** exportiert und das Feld verhält sich
wie ein normales.

### Ablauf

1. Das abhängige Feld exportiert `reloadOn` (Feld-IDs) und `reloadUrl`.
2. Ändert sich einer dieser Werte, postet das Frontend **nur die Trigger-Werte** an die URL
   (200 ms entprellt, ein Request pro URL, der neueste gewinnt).
3. Der Handler baut dieselbe `FormGroup` auf, bindet nachsichtig, setzt die neuen Optionen und
   antwortet mit `ExportPatch()`.
4. Das Frontend merged den Patch und verwirft Werte, die die neue Liste nicht mehr anbietet.

Ein Reload läuft **immer auch einmal direkt nach dem Aufbau des Formulars**. Die Listen im
ersten Render sind damit optional — sie verhindern nur, dass die Felder kurz leer aussehen.

### Feld deklarieren

```go
status := field.NewSelectField("status", "STATUS", true, statusOptions())
tags := field.NewModelListField("tags", "TAGS", false, "Tag", nil)
tags.BaseField.SetReloadOn(xurl.NewUrlPrefix("/Portal/Thing/FormReload", "/api"), "status")
```

Wie bei allen API-Zielen: `SetReloadOn` exportiert die URL über `PrintPrefix()`.

### Handler

```go
func (ctrl *Controller) FormReload(c echo.Context) error {
    wc := webcontext.GetWebContext(c)

    status, tags, fg := ctrl.buildThingForm(wc.UiContext())

    // Nachsichtig binden: mitten im Ausfüllen ist ein leeres Pflichtfeld normal.
    if err := builder.BindReload(c, fg); err != nil {
        return wc.BadRequest(err.Error())
    }

    tags.Options = ctrl.tagsForStatus(status.Value)

    return c.JSON(http.StatusOK, response.NewReturnFields(fg.ExportPatch()))
}
```

`buildThingForm` muss **dieselbe** Funktion sein, die auch die normale Form-Action benutzt, und
die typisierten Feld-Pointer mit zurückgeben — sonst laufen Formular und Reload auseinander.

- `builder.BindReload(c, fg)` — wie `BindAndValidate`, aber ohne Validierungsfehler. Jedes Feld
  wird zuerst auf seinen Default gebunden; scheitert der Request-Wert, bleibt der Default stehen.
  Nur ein unlesbarer Request-Body ist ein Fehler. Der Overposting-Schutz gilt unverändert.

  Zwei Feinheiten: Das gilt, weil alle mitgelieferten Felder in `BindValue` erst validieren und
  dann zuweisen — ein eigenes Feld, das umgekehrt vorgeht, ließe den ungültigen Request-Wert
  stehen. Und ein Default, der selbst die Validierung nicht besteht (z. B. ein required
  Multi-Select ohne Vorauswahl), wird auch nicht gesetzt; dort bleibt der Nullwert.
- `fg.ExportPatch()` — exportiert genau die Felder mit `ReloadOn` **und** `ReloadURL`, gekeyed nach
  Feld-ID. `Form=false`-Felder werden übersprungen.
- `response.NewReturnFields(...)` — `{"fields": {...}}`, bewusst ohne `done`. Mit
  `.WithMessage(text, response.MessageInfo)` gibt es zusätzlich eine Snackbar.

### Was das Frontend übernimmt — und was nicht

Übernommen werden nur bekannte Properties mit passendem Typ: `list`, `name`, `hint`, `class`,
`required`, `disabled`, `hide`, `search`, `min`, `max`, `params`.

Was davon aus xiri-go tatsächlich ankommt, hängt am Feldexport: `disabled` steht **nicht** im
Basis-Export, ein Go-Backend kann es also nicht patchen. Und ein Patch ist additiv — eine Property,
die gar nicht im Patch steht, bleibt unverändert. Zurücksetzen geht nur, wo der Export einen
expliziten Leerwert kennt: `hint` wird als `null` exportiert und räumt den Hinweis damit ab,
`min`/`max` haben keinen solchen Leerwert.

Bewusst **nicht** übernommen:

| Property | Grund |
| --- | --- |
| `value` | Den Wert behält der Client. Er verwirft nur, was die neue `list` nicht mehr enthält. |
| `type`, `subtype`, `id` | Werden beim Aufbau der Controls einmalig normalisiert. |
| `url` | Ob ein Select Server-Suche macht, entscheidet sich beim Aufbau — eine Änderung käme nie an. |
| `showWhen` | Wird ohnehin clientseitig aus den aktuellen Werten ausgewertet. |

Ein Patch für ein Feld ohne `ReloadOn` oder mit einer anderen `reloadUrl` wird ignoriert.

### Grenzen

- **Nur die Trigger-Werte gehen an den Server**, nicht das ganze Formular — und pro URL nur die,
  von denen deren eigene Felder abhängen. Alles Weitere (Parent-ID o. ä.) gehört in die `reloadUrl`.
  Bei mehreren URLs wird trotzdem bei jeder Trigger-Änderung jede URL angefragt.
- **Abhängige `ModelListField`/Treeselects dürfen kein `URL` setzen.** Mit gesetzter URL lädt das
  Frontend den Baum selbst per GET und ignoriert die gepatchte Liste.
- **Keine Step-übergreifenden Abhängigkeiten:** in einem Multi-Step-Form kann ein Feld in Schritt 2
  nicht auf ein Feld aus Schritt 1 reagieren.
- **Ketten funktionieren, Zyklen laufen leer.** Hängt C an B und wird B durch einen Patch geleert,
  lädt C nach. Das terminiert, weil ein Patch nur Werte verwirft.
