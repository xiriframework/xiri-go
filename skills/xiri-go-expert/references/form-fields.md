# Form Fields (form/field package)

Import: `"github.com/xiriframework/xiri-go/form/field"`

Alle Fields implementieren `FormField` Interface und embedden `*BaseField`.

## Gemeinsame BaseField-Methoden (alle Fields)

```go
field.SetClass("xcol-md-6")      // CSS-Klasse (Grid-Breite)
field.SetHint("tooltip.text")    // Tooltip/Hilfetext
field.SetStep(1)                 // Schritt in Multi-Step-Forms
field.SetDisabled(true)          // Deaktiviert
field.SetAccess([]string{"admin"}) // Zugriffskontrolle
field.SetScenario([]string{"add"}) // Nur in bestimmten Szenarien
field.SetForm(false)             // Nicht im Formular anzeigen
```

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
f := field.NewTextField("name", "vehicle.name", true, nil)
// Parameter: id, translationKey, required, currentValue (*string)

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
f := field.NewIntField("count", "device.count", true, nil)
// Parameter: id, translationKey, required, currentValue (*int32)

f := field.NewIntFieldWithBounds("count", "device.count", true, nil, intPtr(0), intPtr(100))
// Mit Min/Max-Grenzen

f := field.NewNumberField("price", "device.price", true, nil)
// Alias mit Subtype "float"

// Optionale Konfiguration:
f.Subtype = "float"   // int (default), pint (positive int), float, bigint, real
f.TextPrefix = "€"
f.TextSuffix = "Stück"

// Nach BindAndValidate:
count := *f.Value  // int32
```

## BoolField

Checkbox. Value: `*bool`

```go
f := field.NewBoolField("active", "device.active", false, nil)
// Parameter: id, translationKey, required, currentValue (*bool)

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

// Spezial-Konstruktoren:
f := field.NewDeviceListField("device_ids", "route.devices", true, currentIDs)

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
f := field.NewTimeField("start", "event.start", true, nil)
// Parameter: id, translationKey, required, defaultValue (*int64)

f.Subtype = "datetime"  // date (default), datetime, time

// Tag-Offsets für Min/Max (z.B. -30 Tage bis +365 Tage):
f.Min = int64Ptr(-30)
f.Max = int64Ptr(365)

// Oder absolute Timestamps:
f.MinDate = &time.Time{...}
f.MaxDate = &time.Time{...}

// Nach BindAndValidate:
timestamp := *f.Value  // int64 (Unix seconds)
```

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

Tag/Chip-Eingabe. Value: `[]string`

```go
f := field.NewChipsField("tags", "device.tags", false)
f.SetList([]field.SelectOption{
    {Value: "indoor", Label: "Indoor"},
    {Value: "outdoor", Label: "Outdoor"},
})
f.SetFreeText(true)  // Erlaubt freie Eingabe

// Nach BindAndValidate:
tags := f.Value  // []string
```

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

active := field.NewBoolField("active", "Aktiv", false, nil)
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
